FROM debian:9.4-slim
MAINTAINER Jason Martin <jason.martin83@protonmail.com>
COPY ./fractal-interview /fractal-interview
COPY ./env.json /env.json
EXPOSE 8080/tcp

ENTRYPOINT ["/fractal-interview"]
