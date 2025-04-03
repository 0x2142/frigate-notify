FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/

# gcc & musl-dev required for matrix support
RUN apk --no-cache add tzdata ca-certificates gcc musl-dev
RUN go mod download

COPY . /app/

# Tag goolm & ldflags required for matrix support
RUN go build -tags "goolm" -a -ldflags '-linkmode external -extldflags "-static"' -o /frigate-notify .

FROM scratch

WORKDIR /app

COPY --from=build /frigate-notify /app/frigate-notify
COPY /templates /app/templates
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8000

ENTRYPOINT [ "/app/frigate-notify" ]
