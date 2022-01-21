default: test

test:
	cd tests && go test -v -count=1 -cover -coverpkg=github.com/ShoshinNikita/go-errors ./...

bench: BENCH?=/go-errors
bench:
	cd tests && go test -v -run=^$$ -benchmem -bench="${BENCH}" -count=5
