.PHONY: run build ansible deploy lint-ansible lint-fix-ansible

run:
	go run cmd/server/server.go

build:
	go build -o go-chatgpt-telegram-bot cmd/server/server.go

ansible:
	ansible-playbook -i ansible/inventories/staging.yml ansible/provision.yml

deploy: build ansible

lint-ansible:
	docker run -it --rm -v ${PWD}:/mnt haxorof/ansible-lint ansible

lint-fix-ansible:
	docker run -it --rm -v ${PWD}:/mnt haxorof/ansible-lint --fix ansible
