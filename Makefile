build:
	go build -o build/dev .

test:
	go test $(shell go list ./... | grep -v /e2e)

e2e:
	go test ./e2e/ -v -count=1

clean:
	rm -rf build/

.PHONY: build test e2e clean
