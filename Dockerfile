FROM golang:1.17.6 as builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-X main.APIVersion=v1.0.0 -X main.Environment=production" -o gosocialgraph cmd/main.go

# Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/gosocialgraph .

ENTRYPOINT /app/gosocialgraph

EXPOSE 3010