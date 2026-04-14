from pydantic import BaseModel


class TempReading(BaseModel):
    temp: float
    humidity: float
    reading_datetime: str
