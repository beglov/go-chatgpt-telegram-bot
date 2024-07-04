.PHONY: run build ansible deploy

run:
	go run cmd/server/server.go

build:
	go build -o go-chatgpt-telegram-bot cmd/server/server.go

ansible:
	ansible-playbook -i ansible/inventories/staging.yml ansible/provision.yml

deploy: build ansible
