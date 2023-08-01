FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o tictactoe-server main.go

FROM alpine:latest AS prod 

WORKDIR /app
COPY --from=builder /app/tictactoe-server .
# Expose the port the server will listen on
EXPOSE 8080

# Run the server
CMD ["./tictactoe-server"]