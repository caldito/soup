FROM ubuntu
MAINTAINER pablo@caldito.me
COPY bin/soup /bin/
RUN ["chmod", "+x", "/bin/soup"]
CMD ["/bin/soup"]
