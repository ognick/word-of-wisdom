FROM golang:1.22-alpine as builder

RUN go version

COPY . /github.com/ognick/word-of-wisdom/

WORKDIR /github.com/ognick/word-of-wisdom/

RUN go mod download
RUN GOOS=linux go build -o ./.bin/wow_client ./cmd/client/main.go

FROM alpine:latest as runner

WORKDIR /app/

COPY --from=0 /github.com/ognick/word-of-wisdom/.bin/wow_client .

CMD ["./wow_client"]