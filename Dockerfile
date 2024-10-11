FROM golang:1.23.1-alpine3.20 AS builder

COPY . /github.com/Muaz717/todo-app/
WORKDIR /github.com/Muaz717/todo-app/

RUN go mod download

RUN GOOS=linux go build -o ./.bin/todo-app ./cmd/todo-app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/Muaz717/todo-app/.bin/todo-app .
COPY --from=builder /github.com/Muaz717/todo-app/config/local.yaml config/
COPY --from=builder /github.com/Muaz717/todo-app/config/local_compose.yaml config/
COPY --from=builder /github.com/Muaz717/todo-app/.env .

EXPOSE 8083

CMD ["./todo-app"]

# FROM golang:1.23.1-alpine

# WORKDIR /app

# COPY go.mod go.sum ./

# RUN go mod download

# COPY ./ ./

# RUN go build -o ./.bin/todo-app ./cmd/todo-app/main.go

# CMD [ "./.bin/todo-app" ]