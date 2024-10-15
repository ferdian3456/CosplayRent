export POSTGRES_URL='postgresql://ferdian:secret_password@localhost:5432/cosplayrent?sslmode=disable'

migrate-create:
	@ migrate create -ext sql -dir scripts/migrations -seq $(name)

migrate-up:
	@ migrate -database ${POSTGRES_URL} -path scripts/migrations up

migrate-down:
	@ migrate -database ${POSTGRES_URL} -path scripts/migrations down