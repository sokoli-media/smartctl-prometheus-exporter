FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /prometheus-exporter

CMD ["/prometheus-exporter"]

FROM alpine:latest

RUN apk add smartmontools

COPY --from=build /prometheus-exporter /prometheus-exporter
COPY /monitoring /monitoring

ENTRYPOINT ["/prometheus-exporter"]
