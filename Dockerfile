FROM golang:1.19-alpine as builder

COPY . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix 'static' -o /app cmd/main.go

FROM alpine:latest
COPY --from=builder /app .
COPY config.yml /
EXPOSE 8000/tcp
ENTRYPOINT ["/app"]
