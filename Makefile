PROJECT:=configmap-update

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o configmap-update .