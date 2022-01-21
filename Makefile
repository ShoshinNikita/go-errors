default: test

test:
	cd tests && go test -v -count=1 -cover -coverpkg=github.com/ShoshinNikita/go-errors ./...
