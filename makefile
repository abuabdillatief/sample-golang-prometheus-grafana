terraform-init:
	@cd terraform && terraform init

terraform-plan: terraform-init
	@cd terraform && terraform plan

terraform-apply: terraform-init
	@cd terraform && terraform apply

tool_up:
	@docker compose up -d

run:
	@go run main.go

load_test:
	@go run loadtest/test.go