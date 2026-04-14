import redis
import asyncio
from fastapi import FastAPI, WebSocket
from dotenv import load_dotenv
import os
import logging
from redis.typing import ResponseT
from src.models import TempReading
from src.redis_db.redis_db import parse_data
import psycopg2


load_dotenv("../.env")

stream_key = "stream"

# The redis_url and redis_port correspond to the Redis database containing the temperature and humidity measures
redis_url, redis_port = os.environ.get("REDIS_CONN_STR").split(":")

rdb = redis.Redis(host=redis_url, port=int(redis_port), db=0, protocol=2)

logging.info("Connected to Redis database")


app = FastAPI()


@app.get("/today")
async def get_today_reading() -> list:
    readings = []
    with psycopg2.connect(
        dbname="postgres", user="postgres", password="example", host="localhost"
    ) as conn:
        with conn.cursor() as cur:
            cur.execute(
                "SELECT * from readings where date_trunc('day', time_reading) = date_trunc('day', now())"
            )
            res = cur.fetchall()
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
            groupname="server-consumer",
            consumername="server",
            streams={stream_key: ">"},
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
                        groupname="server-consumer",
                        consumername="server",
                        streams={stream_key: ">"},
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
