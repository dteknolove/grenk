.PHONY: build tar buildtar 

BUILD_DARWIN_AMD = GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin_amd/grenk ./cmd/grenk/main.go
BUILD_DARWIN_ARM = GOOS=darwin GOARCH=arm64 go build -o ./bin/darwin_arm/grenk ./cmd/grenk/main.go
BUILD_WINDOWS = GOOS=windows GOARCH=amd64 go build -o ./bin/windows/grenk ./cmd/grenk/main.go
TAR_DARWIN_AMD = tar -czvf ./tar/darwin_amd/grenk_darwin_amd64.tar.gz ./bin/darwin_amd/grenk
TAR_DARWIN_ARM = tar -czvf ./tar/darwin_arm/grenk_darwin_arm64.tar.gz ./bin/darwin_arm/grenk
TAR_WINDOWS = tar -czvf ./tar/windows/grenk_win64.tar.gz ./bin/windows/grenk

build:
	$(BUILD_DARWIN_AMD)
	$(BUILD_DARWIN_ARM)
	$(BUILD_WINDOWS)

tar:
	$(TAR_DARWIN_AMD)
	$(TAR_DARWIN_ARM)
	$(TAR_WINDOWS)

buildtar: build tar



