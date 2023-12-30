GOBIN = ./build/bin
postgres:
	docker compose up postgres
clean:
	rm -fr $(GOBIN)/*
