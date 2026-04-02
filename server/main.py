import redis
import asyncio
from fastapi import FastAPI, WebSocket

stream_key = "stream"

rdb = redis.Redis(host="localhost", port=6379, db=0, protocol=2)

app = FastAPI()


@app.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket) -> None:
    await websocket.accept()
    while True:
        entry: list = rdb.xread(streams={stream_key: 0}, count=1)
        data = parse_data(entry)
        print(data)
        await websocket.send_json({"data": data})
        while entry:
            await asyncio.sleep(1)
            try:
                last_id = entry[0][1][0][0]
                entry = rdb.xread(streams={stream_key: last_id}, count=1)
                data = parse_data(entry)
                print(data)
                await websocket.send_json({"data": data})
            except IndexError:
                pass


def parse_data(entry: list) -> dict:
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
