FROM golang:1.16.8-alpine AS builder
RUN mkdir /build
COPY . /build
WORKDIR /build
RUN go build -o psim

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/psim /app/
WORKDIR /app
CMD ["./psim"]