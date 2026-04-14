from redis.typing import ResponseT


async def parse_data(entry: ResponseT) -> dict | None:
    """Parses a Redis stream row and returns a dictionary containing the temperature and humidity.
    Args:
        entry: A Redis stream row containing the row ID and the data associated do the entry.

    Returns:
        dict: a dictionary containing the temperature and humidity.
    """
    try:
        data: dict = await entry[0][1][0][1]
        return {
            key.decode("utf-8"): value.decode("utf-8") for key, value in data.items()
        }
    except IndexError:
        pass
