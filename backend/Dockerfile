# Build stage
FROM golang:1.22-bookworm AS builder

RUN apt-get update && apt-get install -y gcc g++ cmake && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .

ENV CGO_ENABLED=1
RUN go build -o /server ./cmd/server

# Runtime stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates libstdc++6 && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /server /app/server
COPY data/player_data /data/player_data

ENV DATA_PATH=/data/player_data
ENV MINIMAP_PATH=/data/player_data/minimaps
ENV PORT=8080
ENV ALLOWED_ORIGINS=*

EXPOSE 8080
CMD ["/app/server"]
