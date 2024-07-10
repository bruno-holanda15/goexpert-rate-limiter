FROM golang:1.22-bullseye

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /rate_limiter

EXPOSE 8080

CMD ["./rate_limiter"]
