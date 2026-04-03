FROM python:3.14

WORKDIR /usr/local/app

COPY requirements.txt ./ 
RUN pip install --no-cache-dir -r requirements.txt

COPY src ./src
EXPOSE 8000

RUN useradd app
USER app

CMD ["fastapi", "dev", "./src/main.py"]
