create-migration:
	docker run -v \
		./migrations:/migrations \
		--rm --name migrator \
		migrate/migrate \
		create -ext sql -dir /migrations -seq $(name)