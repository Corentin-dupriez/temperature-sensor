import redis
import time
from fastapi import FastAPI, WebSocket

stream_key = "stream"

rdb = redis.Redis(host="localhost", port=6379, db=0, protocol=2)

app = FastAPI()


@app.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket):
    await websocket.accept()
    while True:
        l = rdb.xread(streams={stream_key: 0}, count=1)
        print(l)
        while l:
            time.sleep(1)
            last_id = l[0][1][0][0]
            l = rdb.xread(streams={stream_key: last_id}, count=1)
            await websocket.send_json({"data": l})
