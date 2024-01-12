GOBIN = ./build/bin
redis:
	docker compose up redis
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

revive:
	$(GOPATH)/bin/revive -config ./revive.toml -formatter friendly ./...