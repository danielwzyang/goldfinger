.PHONY: build run pprof perft ucimoves

build:
	go build -o goldfinger ./cmd/uci/main.go

run:
	go run ./cmd/self/main.go

pprof:
	go tool pprof cpu.prof

perft:
	go test -v ./test/perft/perft_test.go

ucimoves:
	go run ./cmd/self/main.go --time=100 | grep bestmove | sed -E "s/bestmove (.*)$$/\1/g"
