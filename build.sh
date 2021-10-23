mkdir -p build/
rm -rf build/*

GOOS=windows GOARCH=amd64 go build -o build/cftools_relay_win64.exe handler/cmd/cmd.go
GOOS=linux GOARCH=amd64 go build -o build/cftools_relay_linux handler/cmd/cmd.go
