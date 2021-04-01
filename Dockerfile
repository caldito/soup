FROM ubuntu
MAINTAINER pablo@caldito.me
RUN apt update && apt install -y ca-certificates
COPY bin/soup /bin/
RUN chmod +x /bin/soup
CMD ["/bin/soup"]
