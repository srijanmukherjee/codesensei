FROM golang:latest


WORKDIR /app

COPY . .

RUN go mod tidy

RUN make server

EXPOSE 3000

CMD ["./bin/server"]