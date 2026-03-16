FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# ─── local ───────────────────────────────────────────────────
FROM builder AS build-local
RUN go build -o server ./cmd/local

FROM alpine:latest AS local
WORKDIR /app
COPY --from=build-local /app/server .
EXPOSE 8080
CMD ["./server"]

# ─── lambda ──────────────────────────────────────────────────
FROM builder AS build-lambda
RUN GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/lambda

FROM public.ecr.aws/lambda/provided:al2 AS lambda
COPY --from=build-lambda /app/bootstrap /var/task/bootstrap
CMD ["bootstrap"]
