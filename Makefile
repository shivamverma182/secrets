build:
	go build .
run:
	go run . github --github-env dev
tidy:
	go mod tidy -v
go-vendor:
	go mod vendor
push:
	git push origin main