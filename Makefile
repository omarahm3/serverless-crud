build-functions:
	export GO111MODULE=on
	cd ./functions/getUser/ && env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../bin/getUser getUser.go && cd ../..
