#! /bin/sh

mkdir -p bin

echo "构建linux amd 64版本"
GOOS=linux GOARCH=amd64 go build cmd/govm/main.go && upx main && mv main bin/govm-linux-amd-64
echo "构建linux amd 64版本完成"

echo "构建linux arm 64版本"
GOOS=linux GOARCH=arm64 go build cmd/govm/main.go && upx main && mv main bin/govm-linux-arm-64
echo "构建linux arm64版本完成"

echo "构建darwin 64版本"
GOOS=darwin GOARCH=amd64 go build cmd/govm/main.go && upx main && mv main bin/govm-darwin-64
echo "构建darwin done版本"

echo "构建darwin arm-64 (m1)版本"
GOOS=darwin GOARCH=arm64 go build cmd/govm/main.go && upx main && mv main bin/govm-darwin-arm-64
echo "构建darwin arm-64 (m1) 版本完成"

echo "构建windows 64版本"
GOOS=windows GOARCH=amd64 go build cmd/govm/main.go && upx main.exe && mv main.exe bin/govm-windows-64.exe
echo "构建windows 64版本完成"