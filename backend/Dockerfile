FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# ─── local  ──────────────────────────────
FROM builder AS build-local
RUN go build -o server ./cmd/local

FROM alpine:latest AS local
WORKDIR /app
COPY --from=build-local /app/server .
EXPOSE 8080
CMD ["./server"]

# ─── lambda  ──────────────────────────
FROM builder AS build-lambda
# CGO_ENABLED=0 を追加し、バイナリ名を bootstrap に固定
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/lambda

# ベースイメージは最新の al2023 を推奨
FROM public.ecr.aws/lambda/provided:al2023 AS lambda
# /var/task は Lambda のデフォルト作業ディレクトリ
COPY --from=build-lambda /app/bootstrap /var/runtime/bootstrap
# 実行ファイルを指定（providedイメージでは bootstrap という名前なら自動で認識されます）
CMD ["bootstrap"]