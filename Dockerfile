# ---------- BUILD STAGE ----------
FROM golang:1.26 AS builder

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

# ---------- RUNTIME STAGE ----------
FROM alpine:3.20

WORKDIR /

RUN apk --no-cache add ca-certificates

COPY --from=builder /main .

EXPOSE 8080

CMD ["./main"]