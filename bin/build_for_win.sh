go generate
GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui -w -s" -o build/qrcode_gene.exe