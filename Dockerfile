FROM alpine:3.4

MAINTAINER John Weldon <johnweldon4@gmail.com>
LABEL Description="Squid Block Landing Page" Vendor="John Weldon Consulting"

RUN apk update && \
    apk upgrade && \
    apk add \
        ca-certificates \
        curl \
    && rm -rf /var/cache/apk/*

ADD public /var/www/public

ADD templates /var/www/templates

COPY block_landing /bin/block_landing

VOLUME ["/var/www"]

WORKDIR /var/www

EXPOSE 9000

CMD ["/bin/block_landing"]
