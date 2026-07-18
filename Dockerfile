FROM golang:1.25-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /src
COPY go.mod ./
RUN go mod tidy && go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=1 go build -o /app cmd/app/main.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app .
ENTRYPOINT ["./app"]
