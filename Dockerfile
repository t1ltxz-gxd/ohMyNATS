FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add make musl-dev gcc gettext

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o ./bin/app cmd/app/main.go

FROM alpine

COPY --from=builder /usr/local/src/bin/app /usr/bin

COPY config/config.yml /config/config.yml

CMD ["/usr/bin/app"]