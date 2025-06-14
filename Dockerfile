FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/*.html ./   
# COPY ./static/image /app/image
# COPY ./static/image /app/static/image


EXPOSE 8080

CMD ["./main"]
