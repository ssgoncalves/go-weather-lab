FROM golang:1.24 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/cloudrun .
ENV PORT=8080
ENV LOAD_ENV=false
ENTRYPOINT ["/app/cloudrun"]