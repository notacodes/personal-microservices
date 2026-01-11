FROM golang:1.25-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o server .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates \
    && adduser -D appuser

USER appuser
COPY --from=build /app/server /usr/local/bin/server

EXPOSE 3411
ENTRYPOINT ["/usr/local/bin/server"]
