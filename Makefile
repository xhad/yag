

test :; go test -v ./tests

build :; go build -o yag cmd/yag/main.go

run :; go run cmd/yag/main.go

