build_docker:
	docker-compose start

build_auth:
	cd auth_server && mix run --no-halt

build_mailer_server:
	cd mailer_server\cmd\server && go run .\main.go

build_app: build_docker build_auth build_mailer_server 

