FROM golang:latest

RUN mkdir /app
ADD main.go /app/
WORKDIR /app

RUN go get "github.com/gorilla/mux"
RUN go get "github.com/lib/pq"
RUN go build -o main .

CMD ["/app/main"]

EXPOSE 8080