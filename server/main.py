import redis
import asyncio
from fastapi import FastAPI, WebSocket
from dotenv import load_dotenv
import os

from redis.typing import ResponseT

load_dotenv("../.env")

stream_key = "stream"

# The redis_url and redis_port correspond to the Redis database containing the temperature and humidity measures
redis_url, redis_port = os.environ.get("REDIS_CONN_STR").split(":")

rdb = redis.Redis(host=redis_url, port=int(redis_port), db=0, protocol=2)

app = FastAPI()


@app.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket) -> None:
    await websocket.accept()
    while True:
        entry: ResponseT = rdb.xread(streams={stream_key: 0}, count=1)
        data: dict = parse_data(entry)
        print(data)
        await websocket.send_json({"data": data})
        while entry:
            await asyncio.sleep(1)
            try:
                last_id: int = entry[0][1][0][0]
                entry: ResponseT = rdb.xread(streams={stream_key: last_id}, count=1)
                data: dict = parse_data(entry)
                print(data)
                await websocket.send_json({"data": data})
            except IndexError:
                pass


def parse_data(entry: ResponseT) -> dict:
    """Parses a Redis stream row and returns a dictionary containing the temperature and humidity.
    Args:
        entry: A Redis stream row containing the row ID and the data associated do the entry.

    Returns:
        dict: a dictionary containing the temperature and humidity.
    """
    data: dict = entry[0][1][0][1]
    return {
        key.decode("utf-8"): float(value.decode("utf-8")) for key, value in data.items()
    }
