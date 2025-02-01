FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/

RUN apk --no-cache add tzdata ca-certificates
RUN go mod download

COPY . /app/

RUN go build -o /frigate-notify .

FROM scratch

WORKDIR /app

COPY --from=build /frigate-notify /app/frigate-notify
COPY /templates /app/templates
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8000

ENTRYPOINT [ "/app/frigate-notify" ]
