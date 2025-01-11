include .env

migrate-create:
	@ migrate create -ext sql -dir db/migrations -seq $(name)

migrate-up:
	@ migrate -database ${POSTGRES_URL} -path db/migrations up

migrate-down:
	@ migrate -database ${POSTGRES_URL} -path db/migrations down
