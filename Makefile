build:
	go build .
run:
	go run . github --github-env dev
tidy:
	go mod tidy -v
vendor:
	go mod vendor