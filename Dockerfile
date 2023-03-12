FROM golang:latest
RUN go version

COPY . .
RUN go mod download && \
    go build -o social-app ./cmd/main.go

CMD ["./social-app"]