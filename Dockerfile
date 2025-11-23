FROM golang:latest

WORKDIR /application
COPY . .
RUN make generate && make build
CMD ["./bin/pr-review"]