# Build stage
FROM golang:1.20-alpine3.17 AS builder

WORKDIR /src
COPY . .
RUN go build -o flugo cmd/api/main.go

# Run stage
FROM alpine:3.17
WORKDIR /src
COPY --from=builder /src/flugo .
COPY api.env .
COPY ./internal/database/migrations ./internal/database/migrations
COPY ./uploads ./uploads

EXPOSE 3000
CMD [ "./flugo" ]