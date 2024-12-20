FROM golang:1.23 AS build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN apt-get update && apt-get install -y git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -a -installsuffix cgo -o transaction-service .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app/transaction-service /app/transaction-service
COPY frontend /app/frontend
COPY .env .env
EXPOSE 8081
RUN chmod +x /app/transaction-service
CMD ["./transaction-service"]