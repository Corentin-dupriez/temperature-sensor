import redis
import asyncio
from fastapi import FastAPI, WebSocket

stream_key = "stream"

rdb = redis.Redis(host="localhost", port=6379, db=0, protocol=2)

app = FastAPI()


@app.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket):
    await websocket.accept()
    while True:
        l = rdb.xread(streams={stream_key: 0}, count=1)
        data = l[0][1][0][1]
        data = {
            key.decode("utf-8"): float(value.decode("utf-8"))
            for key, value in data.items()
        }
        print(data)
        await websocket.send_json({"data": data})
        while l:
            await asyncio.sleep(1)
            last_id = l[0][1][0][0]
            l = rdb.xread(streams={stream_key: last_id}, count=1)
            data = l[0][1][0][1]
            data = {
                key.decode("utf-8"): float(value.decode("utf-8"))
                for key, value in data.items()
            }
            print(data)
            await websocket.send_json({"data": data})
