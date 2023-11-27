db_up:
	docker compose --env-file .env up
db_migration:
	@if [ -z "${name}" ]; then echo "Использование: make db_migration name=name_of_migration_file"; exit 1; fi
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir ../internal/store/migrations create ${name} sql

run_test:
	../tools/gophermarttest \
		-test.v -test.run=^TestGophermart$ \
		-gophermart-binary-path=../cmd/gophermart/gophermart \
		-gophermart-host=localhost \
		-gophermart-port=8080 \
		-gophermart-database-uri="user=postgres password=postgres host=localhost port=5432 dbname=testDB sslmode=disable" \
		-accrual-binary-path=../cmd/accrual/accrual_linux_amd64 \
		-accrual-host=localhost \
		-accrual-port=7777 \
		-accrual-database-uri="user=postgres password=postgres host=localhost port=5432 dbname=testDB sslmode=disable"
