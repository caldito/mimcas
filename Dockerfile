FROM alpine:3.18.4
MAINTAINER pablo@caldito.me
COPY bin/mimcas-server /bin/
RUN chmod +x /bin/mimcas-server
CMD ["/bin/mimcas-server"]
