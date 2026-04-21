from redis.typing import ResponseT
import os
import redis
import logging


REDIS_STREAM_KEY = "stream"


def connect_to_redis() -> redis.Redis:
    # The redis_url and redis_port correspond to the Redis database containing the temperature and humidity measures
    redis_url, redis_port = os.environ.get("REDIS_CONN_STR", "backend:6379").split(":")

    rdb = redis.Redis(host=redis_url, port=int(redis_port), db=0, protocol=2)

    logging.info("Connected to Redis database")

    return rdb


def read_stream(rdb: redis.Redis):
    entry: ResponseT = rdb.xreadgroup(
        groupname=os.environ.get("REDIS_PYTHON_CONSUMER_GROUP", "server-consumer"),
        consumername=os.environ.get("REDIS_PYTHON_CONSUMER_NAME", "server"),
        streams={REDIS_STREAM_KEY: ">"},
        count=1,
    )
    return entry


async def parse_data(entry: ResponseT) -> dict[str, str] | None:
    """Parses a Redis stream row and returns a dictionary containing the temperature and humidity.
    Args:
        entry: A Redis stream row containing the row ID and the data associated do the entry.

    Returns:
        dict: a dictionary containing the temperature and humidity.
    """
    try:
        data: dict = entry[0][1][0][1]
        return {
            key.decode("utf-8"): value.decode("utf-8") for key, value in data.items()
        }
    except IndexError:
        pass
