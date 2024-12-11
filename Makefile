.PHONY: all
all: deploy

.PHONY: binary
binary:
	GOARCH=amd64 go build -o wild-director .

.PHONY: deploy
deploy: binary
	ssh root@5.78.103.252 "systemctl stop wild-director"
	scp wild-director root@5.78.103.252:/usr/bin/wild-director
	ssh root@5.78.103.252 "systemctl start wild-director"
