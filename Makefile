build.image:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	chmod +x server
	docker build --build-arg CONFIG_PATH=${CONFIG_PATH} --build-arg PORT=${PORT} -t server:release .
	rm -rf ./server