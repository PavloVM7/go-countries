GOBIN = ./build/bin
db-tests:
	docker compose up -d postgres
	go test -v ./repository/rdbms/...
	docker compose down
postgres:
	docker compose up postgres
dc-stop:
	docker compose down
clean:
	rm -fr $(GOBIN)/*
