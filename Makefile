# 编译到 Linux
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/main-linux ./main.go

# 编译到 macOS
.PHONY: build-darwin
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./build/main-darwin /main.go

# 编译到 windows
.PHONY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/main-windows.exe ./main.go

# 编译到 Linux
.PHONY: build-jwt-linux
build-jwt-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/jwt-linux ./jwt.go

# 编译到 macOS
.PHONY: build-jwt-darwin
build-jwt-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./build/jwt-darwin ./jwt.go
