FROM golang:latest
RUN go version
ENV GOPATH=/

COPY . .
RUN apt-get update &&\
    apt-get -y install postgresql-client && \
    go mod download && \
    go build -o social-app ./cmd/main.go

CMD ["./social-app"]