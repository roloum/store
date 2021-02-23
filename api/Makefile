BASE_DIR= github.com/roloum/store/api
BUILD_CMD= env GOOS=linux go build -ldflags="-s -w" -o
TEST_CMD= go test -timeout 30s

.PHONY: build deploy

.PHONY: test
test:
	${TEST_CMD} ${BASE_DIR}/internal/store/cart/

build:
	export GO111MODULE=on
	${BUILD_CMD} bin/cart cmd/lambda/handlers/cart/main.go
	${BUILD_CMD} bin/item cmd/lambda/handlers/item/main.go
