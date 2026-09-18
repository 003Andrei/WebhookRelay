.PHONY: db db-down db-reset db-shell api fmt vet

db:
	docker compose up -d db 

db-down: 
		docker compose down 

db-reset: 
	docker compose down-v 
	docker compose up -d db 

db-shell: 
	docker compose exec db psql -U hook hookrelay 

api: 
	go run ./cmd/api

fmt: 
	go fmt ./.. 

vet: 
	go vet ./..
