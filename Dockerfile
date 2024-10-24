FROM golang:1.23-alpine as builder

WORKDIR /app
COPY . .
RUN go build -o main .

# Second stage: minimal container
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/templates .

CMD ["./main"]
