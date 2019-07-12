FROM 		golang:1.12-alpine as builder
WORKDIR		/go
ADD			app.go .
RUN			CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app


FROM 		scratch
COPY 		--from=builder /go/app /app
ENTRYPOINT  ["/app"]
