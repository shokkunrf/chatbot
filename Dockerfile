FROM golang:1.24-bookworm AS builder
WORKDIR /build
COPY app/ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o chatbot .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /build/chatbot /chatbot
USER nonroot
ENTRYPOINT ["/chatbot"]
