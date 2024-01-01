GOBIN = ./build/bin
postgres:
	docker compose up postgres
dc-stop:
	docker compose down
clean:
	rm -fr $(GOBIN)/*
