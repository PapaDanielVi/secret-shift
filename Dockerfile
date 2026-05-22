# Stage 1: Build the application
FROM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o secret-shift .

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/secret-shift /secret-shift

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/secret-shift"]
