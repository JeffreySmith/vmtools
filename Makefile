all: macos linux

macos:
	GOARCH=arm64 GOOS=darwin go build -o addusers_osx cmd/create_users/main.go
linux:
	GOARCH=amd64 GOOS=linux go build  -o addusers_linux cmd/create_users/main.go
main:
	go build -o addusers cmd/create_users/main.go
clean:
	rm -f main_osx main_linux main addusers