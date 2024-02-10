FROM golang:1.22-alpine as build

WORKDIR /app

COPY src .

RUN CGO_ENABLED=0 GOOS=linux go build -o pubstatic

FROM alpine:latest as run

WORKDIR /app

COPY --from=build /app/config ./config
COPY --from=build /app/static ./static
COPY --from=build /app/pubstatic .

CMD ["./pubstatic"]