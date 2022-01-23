FROM ubuntu
MAINTAINER pablo@caldito.me
RUN apt update && apt install -y ca-certificates
COPY bin/go-memcached /bin/
RUN chmod +x /bin/go-memcached
CMD ["/bin/go-memcached"]
