FROM golang:latest AS build
RUN go get -v -u 'github.com/lib/pq'
RUN go get -v -u 'github.com/valyala/fasthttp'
RUN go get -v -u 'github.com/qiangxue/fasthttp-routing'
RUN go get 'golang.org/x/crypto/pbkdf2'

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o main -i cmd/main/main.go

FROM ubuntu:18.04 AS release
MAINTAINER Alex Sirmais
ENV TZ=Europe/Moscow

# Установка postgresql и git (гит для загрузки gorilla для Go)
ENV PGVER 10
RUN apt -y update && apt install -y postgresql-$PGVER
RUN yes | apt-get install git


# Установка postgresql
ENV PGVER 10
RUN apt -y update && apt install -y postgresql-$PGVER
USER postgres
COPY sql .
RUN cat ./postgresql.conf >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker myService &&\
    PGPASSWORD=docker psql -U docker -h 127.0.0.1 -d myService -f ./init.sql &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root
WORKDIR /app
COPY --from=build /app/main .
RUN chmod +x ./main
COPY /static ./static

EXPOSE 8080/tcp

USER postgres
CMD service postgresql start && service postgresql status && ./main