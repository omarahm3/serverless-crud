build-functions:
	export GO111MODULE=on
	cd ./functions/getUser/ && env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../bin/getUser getUser.go && cd ../..
	cd ./functions/createUser/ && env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../bin/createUser createUser.go && cd ../..
	cd ./functions/updateUser/ && env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../bin/updateUser updateUser.go && cd ../..
	cd ./functions/deleteUser/ && env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../bin/deleteUser deleteUser.go && cd ../..
