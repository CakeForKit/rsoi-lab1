FROM golang:1.27-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /application ./cmd/application

FROM alpine:3.22

RUN adduser -D -H appuser
WORKDIR /app
COPY --from=build /application /app/application
COPY config/application.yaml /app/config/application.yaml
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/application"]
