clean:
	docker stop $$(docker ps -a -q) && docker rm $$(docker ps -a -q)

allup:
	docker compose up --build -d --force-recreate

reup:
	docker compose up --build -d --force-recreate "$(service)"

docker-shell:
	docker exec -it $$(docker ps | grep "$(service)" | awk '{ print $$1 }') sh

connect-psql:
	docker exec -it $$(docker ps | grep "postgres-db" | awk '{ print $$1 }') psql -U admin -d default
