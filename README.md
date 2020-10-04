[![Build Status](https://travis-ci.com/nnev/website.svg?branch=master)](https://travis-ci.com/nnev/website)

## Setup for testing the website locally

  * install docker: Installation instructions for [debian/linux](https://docs.docker.com/engine/installation/linux/docker-ce/debian/), [macOS](https://docs.docker.com/docker-for-mac/install/) and [windows](https://docs.docker.com/docker-for-windows/install/) can be found on the docker website. For Arch Linux, use `pacman -S docker`.
  * start docker (usually `systemctl start docker.service` will suffice)

### Use the Makefile

  * make sure there are no running services on the port 8080, otherwise change the port in the `Makefile`
  * run `make` (with docker privileges, e.g. as root)
  * browse to [`127.0.0.1:8080`](http://127.0.0.1:8080) or issue `make open` beforehand
  * when you are done, `Ctrl-C` out and rerun `make` to update the test instance with your new changes

After testing your changes, `make clean` provides a convenient way to clean up after you, stopping and removing all the containers you have created in the process of testing.
If you want to get rid of remainders of the process: `make clean` removes all containers, `make purge` removes all containers including all images (be aware that they need to be redownloaded and such next time you test the website, so you probably usually do not want to issue a `make purge` unless you are in severe need of disk space).

### Run test environment manually

  * start postgresql: `docker run -d --name=nnev-postgres postgres`
  * build website: `docker build -t nnev-website .`
  * run website: `docker run --name=nnev-website -p 127.0.0.1:8080:80 --link nnev-postgres:postgres -v $PWD:/usr/src/ nnev-website`
  * browse to [`127.0.0.1:8080`](http://127.0.0.1:8080) to inspect your state of the webpage
  * restart: `Ctrl-C` out and run `docker kill nnev-website; docker rm nnev-website`, then goto `run website`
