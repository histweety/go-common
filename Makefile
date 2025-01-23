project_name = go-common
image_name = histweety/go-common:latest

run-local: ## Run the app locally
	go run main.go

requirements: ## Generate go.mod & go.sum files
	go mod tidy

clean-packages: ## Clean packages
	go clean -modcache