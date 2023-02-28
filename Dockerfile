FROM golang:latest
RUN go version
ENV GOPATH=/

COPY . .
RUN apt-get update &&\
    apt-get -y install postgresql-client && \
    chmod +x wait-for-postgres.sh && \
    go mod download && \
    go build -o social-app ./cmd/main.go

CMD ["./social-app"]