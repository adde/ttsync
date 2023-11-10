build:
	echo "Building for every OS and Platform"
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/ttsync-linux-amd64 ./src/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/ttsync-windows-amd64.exe ./src/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o ./bin/ttsync-darwin-arm64 ./src/main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/ttsync-darwin-amd64 ./src/main.go