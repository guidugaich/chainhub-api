MIGRATE_CMD = docker compose run --rm migrate

.PHONY: migrate-up migrate-down migrate-force migrate-create

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down 1

migrate-force:
	$(MIGRATE_CMD) force $(version)

migrate-create:
	$(MIGRATE_CMD) create -ext sql -dir /migrations -seq $(name)
