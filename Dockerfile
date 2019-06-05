FROM ubuntu
MAINTAINER jwang@wangbenjun@gmail.com
WORKDIR /opt/go
COPY . /opt/go

EXPOSE 8080
CMD ["./main"]