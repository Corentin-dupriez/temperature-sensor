import redis
import asyncio
from fastapi import Depends, FastAPI, WebSocket
from dotenv import load_dotenv
import os
import logging
from redis.typing import ResponseT
from src.models import TempReading
from src.redis_db.redis_db import parse_data
from fastapi.middleware.cors import CORSMiddleware
from typing import List
from fastapi.responses import JSONResponse
from fastapi.encoders import jsonable_encoder
from sqlalchemy import create_engine, text
from sqlalchemy.orm import declarative_base, sessionmaker, Session


load_dotenv("../.env")

REDIS_STREAM_KEY = "stream"

# The redis_url and redis_port correspond to the Redis database containing the temperature and humidity measures
redis_url, redis_port = os.environ.get("REDIS_CONN_STR", "backend:6379").split(":")

rdb = redis.Redis(host=redis_url, port=int(redis_port), db=0, protocol=2)

logging.info("Connected to Redis database")

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
    readings: List[None | dict] = []
    res = db.execute(
        text(
            f"SELECT * from readings where date_trunc('{timeframe}', time_reading) = date_trunc('{timeframe}', now())"
        )
    )
    for result in res:
        readings.append(
            TempReading(
                temp=result[1],
                humidity=result[2],
                reading_datetime=result[3].strftime("%d/%m/%Y %H:%M:%S"),
            ).model_dump()
        )
    return readings


@app.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket) -> None:
    await websocket.accept()
    logging.info("WebSocket connected")
    while True:
        entry: ResponseT = rdb.xreadgroup(
            groupname=os.environ.get("REDIS_PYTHON_CONSUMER_GROUP", "server-consumer"),
            consumername=os.environ.get("REDIS_PYTHON_CONSUMER_NAME", "server"),
            streams={REDIS_STREAM_KEY: ">"},
            count=1,
        )
        if entry != []:
            data: dict = await parse_data(entry)
            print(data)
            new_data: TempReading = TempReading(**data)
            await websocket.send_json({"data": data})
            while entry:
                await asyncio.sleep(1)
                try:
                    entry: ResponseT = rdb.xreadgroup(
                        groupname=os.environ.get(
                            "REDIS_PYTHON_CONSUMER_GROUP", "server-consumer"
                        ),
                        consumername=os.environ.get(
                            "REDIS_PYTHON_CONSUMER_NAME", "server"
                        ),
                        streams={REDIS_STREAM_KEY: ">"},
                        count=1,
                    )
                    data: dict = await parse_data(entry)
                    print(data)
                    await websocket.send_json({"data": data})
                except IndexError:
                    pass
        else:
            print("No new entry, sleeping")
        await asyncio.sleep(1)
