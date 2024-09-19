
build_docker:
	docker-compose start

build_auth:
	cd auth_service && mix phx.server & 

build_mailer_server:
	cd mailer_service && go run ./\cmd/\server/\main.go & 

build_user_manager_service:
	cd UserManager && dotnet run & 

build_client: 
	cd client && yarn dev & 

build_app: build_docker build_auth build_mailer_server build_user_manager_service build_client

