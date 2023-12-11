.PHONY: build tar buildtar move gen doller comma

BUILD_DARWIN_AMD = GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin_amd/grenk ./main.go
BUILD_DARWIN_ARM = GOOS=darwin GOARCH=arm64 go build -o ./bin/darwin_arm/grenk ./main.go
BUILD_LINUX = GOOS=linux GOARCH=amd64 go build -o ./bin/linux/grenk ./main.go
BUILD_WINDOWS = GOOS=windows GOARCH=amd64 go build -o ./bin/windows/grenk ./main.go
TAR_DARWIN_AMD = tar -czvf ./tar/darwin_amd/grenk_darwin_amd64.tar.gz ./bin/darwin_amd/grenk
TAR_DARWIN_ARM = tar -czvf ./tar/darwin_arm/grenk_darwin_arm64.tar.gz ./bin/darwin_arm/grenk
TAR_WINDOWS = tar -czvf ./tar/windows/grenk_win64.tar.gz ./bin/windows/grenk
TAR_LINUX = tar -czvf ./tar/linux/grenk_linux64.tar.gz ./bin/linux/grenk
RM_GRENK = rm -rf /Users/teknolove/go/bin/grenk
MOVE = cp -r ./bin/darwin_amd/grenk /Users/teknolove/go/bin 
GEN_JSON = go run ./experimental/gen/main.go
DOLLAR = go run ./experimental/dollar/main.go
COMMA = go run ./experimental/comma/main.go
comma:
	$(COMMA)
dollar:
	$(DOLLAR)
gen:
	$(GEN_JSON)
build:
	$(BUILD_DARWIN_AMD)
	$(BUILD_DARWIN_ARM)
	$(BUILD_WINDOWS)
	$(BUILD_LINUX)

tar:
	$(TAR_DARWIN_AMD)
	$(TAR_DARWIN_ARM)
	$(TAR_WINDOWS)
	$(TAR_LINUX)

buildtar: build tar

move:
	$(RM_GRENK)
	$(MOVE)
	




