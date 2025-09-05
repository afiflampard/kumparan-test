# --------- Builder Stage ----------
FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o hertz-app ./main.go

# --------- Runtime Stage ----------
FROM alpine:latest

RUN apk --no-cache add ca-certificates bash curl postgresql-client

WORKDIR /app

# Copy binary & env
COPY --from=builder /app/hertz-app .
COPY .env.example .env

# Entry point
COPY wait-for.sh .

ENTRYPOINT ["./wait-for.sh"]
CMD ["./hertz-app"]
