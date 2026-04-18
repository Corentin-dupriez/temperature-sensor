FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./ 

RUN go mod download

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o /histo-db ./cmd/main.go

CMD ["/histo-db"]
