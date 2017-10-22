## Setup

### install docker
Installation instructions for [debian/linux](https://docs.docker.com/engine/installation/linux/docker-ce/debian/), [macOS](https://docs.docker.com/docker-for-mac/install/) and [windows](https://docs.docker.com/docker-for-windows/install/) can be found on the docker website. For Arch Linux, use `pacman -S docker`.

### start postgresql
```
docker run -p 127.0.0.1:5432:5432 postgres
```

### build website
```
docker build -t nnev-website .
```

### run website
```
docker run --name=nnev-website --net=host -p 127.0.0.1:80:80 -v $PWD:/usr/src/ nnev-website
```

### restart
```
docker kill nnev-website
docker rm nnev-website
```
-> see run
