FROM golang:1.26.2-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o go-backend ./cmd

CMD ["./go-backend"]

# 2.26 gb
# dockerignore: 2.16 gb
# stage: 165 mb

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/go-backend .

CMD ["./go-backend"]
