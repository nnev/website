
.DEFAULT_GOAL := website

postgres:
	docker run -d --name=nnev-postgres postgres | exit 0

build: postgres
	docker rm -f nnev-website | exit 0
	docker build --force-rm -t nnev-website .
	docker run --rm --link nnev-postgres:postgres nnev-website /build_entrypoint.sh

website: build
	docker run --name=nnev-website -p 127.0.0.1:8080:80 --link nnev-postgres:postgres nnev-website

travis-run:
	docker run -d --name=nnev-website -p 127.0.0.1:8080:80 --link nnev-postgres:postgres nnev-website
	sleep 15s && curl -I http://127.0.0.1:8080

stop:
	docker rm -f nnev-website
	docker kill postgres | exit 0

clean:
	docker rm -f nnev-website | exit 0
	docker rm -f nnev-postgres | exit 0

purge: clean
	docker rmi -f nnev-website | exit 0
	docker rmi -f nnev-postgres | exit 0

open:
	xdg-open http://127.0.0.1:8080
