import redis
import json
import asyncio
from fastapi import Depends, FastAPI, WebSocket
from dotenv import load_dotenv
import os
import logging
from redis.typing import ResponseT
from src.models import TempReading
from src.redis_db import redis_db, parse_data
from fastapi.middleware.cors import CORSMiddleware
from typing import List
from fastapi.responses import JSONResponse
from fastapi.encoders import jsonable_encoder
from sqlalchemy import create_engine, text
from sqlalchemy.orm import declarative_base, sessionmaker, Session


load_dotenv("../.env")

rdb: redis.Redis = redis_db.connect_to_redis()


DATABASE_URL = f"postgresql://{os.environ.get('POSTGRES_DB_USER', 'postgres')}:{os.environ.get('POSTGRES_DB_PASSWORD', 'example')}@{os.environ.get('POSTGRES_DB_HOST', 'localhost')}:5432/{os.environ.get('POSTGRES_DB_NAME', 'postgres')}"


engine = create_engine(
    DATABASE_URL,
)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base = declarative_base()


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


app = FastAPI()

allowed_origins = [
    f"http://{os.environ.get('APP_HOST', 'localhost')}",
    f"http://{os.environ.get('APP_HOST', 'localhost')}:8080",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=allowed_origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/today")
async def get_today_reading(db: Session = Depends(get_db)) -> JSONResponse:
    readings: List[None | dict] = await get_readings("day", db)
    print(JSONResponse(content=jsonable_encoder(readings)).body)
    return JSONResponse(content=jsonable_encoder(readings))


@app.get("/week")
async def get_weekly_readings(db: Session = Depends(get_db)) -> JSONResponse:
    readings: List[None | dict] = await get_readings("week", db)
    return JSONResponse(content=jsonable_encoder(readings))


async def get_readings(timeframe: str, db: Session) -> List:
    if timeframe not in ("week", "day"):
        raise ValueError("The timeframe should be week or day")
    if timeframe == "day":
        query = f"SELECT * from readings where date_trunc('{timeframe}', time_reading) = date_trunc('{timeframe}', now())"
    else:
        query = f"SELECT avg(temperature), avg(humidity),date_trunc('day', time_reading)  from readings where date_trunc('{timeframe}', time_reading) = date_trunc('{timeframe}', now()) group by date_trunc('day', time_reading)"

    readings: List[None | dict] = []
    res = db.execute(text(query))
    for result in res:
        readings.append(
            TempReading(
                temp=result[0],
                humidity=result[1],
                reading_datetime=result[2].strftime("%d/%m/%Y %H:%M:%S"),
            ).model_dump()
        )
    return readings


@app.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket) -> None:
    await websocket.accept()
    logging.info("WebSocket connected")
    while True:
        entry = redis_db.read_stream(rdb)
        if entry != []:
            data: dict[str, str] | None = await parse_data(entry)
            new_data: TempReading = TempReading(
                temp=data["temp"],
                humidity=data["humidity"],
                reading_datetime=data["time"],
            )
            await websocket.send_json({"data": new_data.model_dump()})
        else:
            print("No new entry, sleeping")
        await asyncio.sleep(1)
