FROM golang:1.21

WORKDIR /app
COPY . .

RUN chmod +x ./deploy/local/wait-for-it.sh

RUN go build ./cmd/service/main.go
CMD ["/app/main"]