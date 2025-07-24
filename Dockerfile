FROM golang:1.25rc2-alpine3.22 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o rest-subscribe ./cmd/rest/main.go

FROM alpine:3.22.1 AS api 
WORKDIR /app

COPY --from=build /app/rest-subscribe .

EXPOSE 8000

CMD ["./rest-subscribe", ":8000"]