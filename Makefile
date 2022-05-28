.phony: format test cov vet race

PACKAGE_FOLDER=pkg
INTERNAL_FOLDER=internal

format:
	go fmt ./${INTERNAL_FOLDER}/...
	go fmt ./${PACKAGE_FOLDER}/...

test: format
	go test -coverprofile=coverage.out ./${PACKAGE_FOLDER}/...

cov: test
	go tool cover -html=coverage.out

vet: format
	go vet ./${INTERNAL_FOLDER}/...
	go vet ./${PACKAGE_FOLDER}/...

race: format
	go test -race -coverprofile=coverage.out ./${PACKAGE_FOLDER}/...
