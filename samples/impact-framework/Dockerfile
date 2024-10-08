FROM ubuntu:latest

WORKDIR /work
COPY demo.yaml .

RUN apt-get update -y  && \
apt-get upgrade -y  && \
apt-get install git npm vim lsb-release -y  && \
chown -R 1000:1000 /work

