FROM golang:1.13 as builder

RUN mkdir /app
ADD main.go /app/
WORKDIR /app

RUN go get "github.com/gorilla/mux"
RUN go get "github.com/lib/pq"
RUN go build -o main .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /server

CMD ["/server"]

EXPOSE 8080