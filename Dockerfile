FROM golang:1.16.8-alpine AS builder
RUN apk add build-base
RUN go get -u github.com/pressly/goose/cmd/goose
RUN mkdir /build
COPY . /build
WORKDIR /build

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
RUN apk --no-cache add tzdata
COPY --from=builder /build/psim /app/
WORKDIR /app
CMD ["./psim"]
