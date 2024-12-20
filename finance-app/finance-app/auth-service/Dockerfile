FROM golang:1.23 AS build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN apt-get update && apt-get install -y git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -a -installsuffix cgo -o auth-service .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app/auth-service /app/auth-service
COPY frontend /app/frontend
COPY .env .env
EXPOSE 8080
RUN chmod +x /app/auth-service
CMD ["./auth-service"]