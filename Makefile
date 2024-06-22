.PHONY=build


test:
	@go test -v ./...

test-fuzz:
	@go test --fuzz=FuzzCalculate --fuzztime=10s -v ./...