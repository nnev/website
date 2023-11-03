#vim:ft=Dockerfile
FROM debian:testing

RUN DEBIAN_FRONTEND=noninteractive apt-get update && \
	apt-get install -y --no-install-recommends \
	jekyll ruby ruby-pg postgresql-client nginx supervisor git golang-go  && \
	rm -rf /var/lib/apt/lists/* && \
	apt-get clean

# TODO: install as debian package once available
RUN gem install icalendar

ENV PGUSER=postgres PGHOST=postgres
# tell jekyll to use utf-8 (website build fails otherwise)
ENV LC_ALL=C.UTF-8
ENV GOPATH=/tmp/go
# for easier interactive usage
ENV PATH="/tmp/go/bin:${PATH}"
# needed for hook service
ENV WEBHOOK_SECRET=geheim

ADD build_entrypoint.sh /build_entrypoint.sh
ADD entrypoint.sh /entrypoint.sh
ADD build_website.sh /build_website.sh
ADD nnev-website-nginx.conf /etc/nginx/sites-available/default
ADD nnev-website-supervisor.conf /etc/supervisor/conf.d/supervisord.conf

RUN useradd -ms /bin/bash nnev
WORKDIR /tmp/go/src/github.com/nnev/website/
COPY . .

# Dockerfiles's ADD/COPY has a --chown arguments but no --chmod argument.
RUN chown -R nnev:nnev "${GOPATH}" && chmod -R a+rX "${GOPATH}"
RUN runuser -u nnev -- go get -v ./...
RUN runuser -u nnev -- go install -v ./...

EXPOSE 80

CMD ["/entrypoint.sh"]
