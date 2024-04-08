FROM golang:1.19-alpine as builder

RUN go version

COPY . /github.com/ognick/word-of-wisdom/

WORKDIR /github.com/ognick/word-of-wisdom/

RUN go mod download
RUN GOOS=linux go build -o ./.bin/wow_server ./cmd/server/main.go

FROM alpine:latest as runner

WORKDIR /app/

COPY --from=0 /github.com/ognick/word-of-wisdom/.bin/wow_server .

CMD ["./wow_server"]