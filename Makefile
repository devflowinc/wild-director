.PHONY: all
all: deploy

.PHONY: binary
binary:
	GOARCH=amd64 go build -o wild-director .

.PHONY: deploy
deploy: binary
	ssh root@demo.trytrieve.com "systemctl stop wild-director"
	scp wild-director root@demo.trytrieve.com:/usr/bin/wild-director
	ssh root@demo.trytrieve.com "systemctl start wild-director"
