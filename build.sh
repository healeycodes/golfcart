GOOS=linux GOARCH=amd64 go build -o golfcart-linux cmd/golfcart.go
GOOS=windows GOARCH=amd64 go build -o golfcart-windows cmd/golfcart.go
GOOS=darwin GOARCH=amd64 go build -o golfcart-darwin cmd/golfcart.go
