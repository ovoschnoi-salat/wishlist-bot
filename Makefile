.PHONY: rebuild
rebuild:
	sudo chown 65532:65532 ./data
	docker compose down --rmi local tg-bot && docker compose up -d tg-bot --no-recreate