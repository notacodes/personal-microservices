# ========================
# Build Stage
# ========================
FROM --platform=linux/amd64 golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o server .

# ========================
# Final Stage
# ========================
FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 1000 appuser

WORKDIR /app

USER appuser

COPY --from=build /app/server /usr/local/bin/server

EXPOSE 3411

ENTRYPOINT ["/usr/local/bin/server"]
