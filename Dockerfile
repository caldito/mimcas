FROM ubuntu
MAINTAINER pablo@caldito.me
RUN apt update && apt install -y ca-certificates
COPY bin/mimcas /bin/
RUN chmod +x /bin/mimcas
CMD ["/bin/mimcas"]
