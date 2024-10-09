FROM ubuntu:latest

WORKDIR /work
COPY ptest.yaml demo.yaml .

RUN apt-get update -y && \
apt-get upgrade -y && \
apt-get install git npm vim lsb-release -y && \
echo "Finished apt-get install" && \
npm install -g --save-dev husky && \
echo "Finished husky install" && \
npm install -g @grnsft/if && \
echo "Finished if install" && \
npm install -g https://github.com/Shivani-G/prometheus-importer
echo "Finsihed pi isntall" && \
chown -R 1000:1000 /work /usr/local/lib/node_modules

