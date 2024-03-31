build_docker:
	docker-compose start

build_auth:
	cd auth_server && mix run --no-halt &

build_mailer_server:
	cd mailer_server && go run ./\cmd/\server/\main.go &

build_client: 
	cd client && yarn dev &

build_app: build_docker build_auth build_mailer_server build_client

