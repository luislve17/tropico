clean:
	docker stop $$(docker ps -a -q) && docker rm $$(docker ps -a -q)

allup:
	docker compose up --build -d --force-recreate

reup:
	docker compose up --build -d --force-recreate "$(service)"

rerun:
	docker compose up "$(service)"

run:
	docker compose up --build "$(service)"

docker-shell:
	docker exec -it $$(docker ps | grep "$(service)" | awk '{ print $$1 }') sh
