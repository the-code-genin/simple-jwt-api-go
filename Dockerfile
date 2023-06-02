FROM golang:alpine3.18

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o build/app ./cmd/app

CMD ["build/app"]