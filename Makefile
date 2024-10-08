.PHONY: rebuild
rebuild:
	docker compose down --rmi local tg-bot && docker compose up -d tg-bot --no-recreate