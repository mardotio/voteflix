create-migration:
	docker run -v \
		./migrations:/migrations \
		--rm --name migrator \
		migrate/migrate \
		create -ext sql -dir /migrations -seq $(name)

migrate-up:
	docker run -v \
		./migrations:/migrations \
		--rm --network voteflix-dev_default \
		migrate/migrate \
		-database "postgres://postgres:postgres@db:5432/voteflix?sslmode=disable&search_path=public" \
		-path=/migrations/ \
		up

rollback-one:
	docker run -v \
		./migrations:/migrations \
		--rm --network voteflix-dev_default \
		migrate/migrate \
		-database "postgres://postgres:postgres@db:5432/voteflix?sslmode=disable&search_path=public" \
		-path=/migrations/ \
		down 1
