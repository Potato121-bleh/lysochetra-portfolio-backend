FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . . 

RUN go build -o myportfolio-backend-app ./cmd

FROM alpine:latest 

WORKDIR /app

COPY --from=builder /app/myportfolio-backend-app .

EXPOSE 3000

CMD ["./myportfolio-backend-app"]
