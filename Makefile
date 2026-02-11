rebuild run:
	docker compose up --build
full delete:
	docker compose down -v
info:
	cloc . 