FROM ubuntu:latest

ENV NPM_CONFIG_PREFIX=/usr/local
ARG NPM_CONFIG_PREFIX=/usr/local

WORKDIR /work
COPY ptest.yaml demo.yaml .

RUN apt-get update -y && \
apt-get upgrade -y && \
apt-get install git npm vim lsb-release -y && \
echo "    --- Finished apt-get install" && \
npm install -g --save-dev husky && \
echo "    --- Finished husky install" && \
npm install -g "@grnsft/if" && \
echo "    --- Finished IF install" && \
git clone https://github.com/Green-Software-Foundation/if.git && cd if && npm install -g && npm install && npm link && cd - && \
echo "    --- Finished IF npm link" && \
npm install -g "https://github.com/Shivani-G/prometheus-importer" && \
echo "    --- Finished PI install" && \
git clone https://github.com/Shivani-G/prometheus-importer.git && cd prometheus-importer && npm install -g && npm install && npm link && cd - && \
echo "    --- Finished cloning PI" && \
cd if && npm link prometheus-importer && cd - && \
echo "    --- Finished last link " && \
chown -R 1000:1000 /work /usr/local/


