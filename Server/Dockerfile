# Build-Stufe
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY Server/go.mod Server/go.sum ./
RUN go mod download
COPY Server/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# Laufzeit-Stufe: nur die Server-Binärdatei (API/WebSocket, ohne Client-Dateien)
FROM alpine:3.20
RUN adduser -D -u 10001 game && mkdir /data && chown game /data
USER game
COPY --from=build /out/server /server
VOLUME /data
EXPOSE 8080
HEALTHCHECK --interval=15s --timeout=3s CMD wget -qO- http://127.0.0.1:8080/healthz || exit 1
ENTRYPOINT ["/server", "-addr", ":8080", "-client", "", "-data", "/data"]
