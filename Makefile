build:
	go get
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build  -x -v
	# makescp raspStats pi@10.1.1.204:/home/pi
