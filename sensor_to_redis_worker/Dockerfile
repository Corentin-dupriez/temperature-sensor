FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./ 

RUN go mod download

COPY *.go ./ 

RUN CGO_ENABLED=0 GOOS=linux go build -o /temp-worker

EXPOSE 6379

CMD ["/temp-worker"]
