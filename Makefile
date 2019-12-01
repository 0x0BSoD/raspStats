build:
	go fmt *.go
	CGO_ENABLED=1 GOOS=linux GOARCH=arm go build
	scp raspStats pi@10.1.1.204:/home/pi
