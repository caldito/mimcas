FROM ubuntu
MAINTAINER pablo@caldito.me
RUN apt update && apt install -y ca-certificates
COPY bin/kv-store /bin/
RUN chmod +x /bin/kv-store
CMD ["/bin/kv-store"]
