build:
	echo "Building for every OS and Platform"
	GOOS=linux GOARCH=amd64 go build -o ./bin/ttsync-linux-amd64 ./src/main.go
	GOOS=windows GOARCH=amd64 go build -o ./bin/ttsync-windows-amd64.exe ./src/main.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/ttsync-darwin-arm64 ./src/main.go