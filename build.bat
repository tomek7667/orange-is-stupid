set GOARCH=arm
set GOOS=linux
go build -o cfrunner ./cmd/cloudflare-runner/main.go
