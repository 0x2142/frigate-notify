FROM golang:bookworm

WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/

RUN go mod download

COPY . /app/

RUN go build -o /app/frigate-notify .

ENTRYPOINT [ "/app/frigate-notify" ]