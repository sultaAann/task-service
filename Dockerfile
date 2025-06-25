FROM golang:1.24.3

RUN apt-get update && apt-get clean

WORKDIR /app

COPY . /app

RUN go build -o app /app/cmd/task-service/main.go 

RUN useradd app

USER app

EXPOSE 8080

CMD ["/app/app"]