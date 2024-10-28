FROM ubuntu:latest

ENV NPM_CONFIG_PREFIX=/usr/local
ARG NPM_CONFIG_PREFIX=/usr/local

WORKDIR /work
COPY ?test.yaml .

RUN apt-get update -y && \
apt-get upgrade -y && \
apt-get install git curl vim lsb-release -y && \
echo "    --- Finished apt-get install" && \
# npm install -g npm@latest  && \
# echo "    --- Finished npm update" && \
curl -fsSL https://deb.nodesource.com/setup_lts.x | bash && \
apt-get upgrade -y && \
apt-get install nodejs -y && \
npm install -g npm@10.9.0 && \
echo "    --- Finished node update" && \
npm install -g typescript rimraf husky @grnsft/if && \
echo "    --- Finished if install" && \
npm install -g "https://github.com/Shivani-G/prometheus-importer" && \
echo "    --- Finished PI install" && \
chown -R 1000:1000 /work /usr/local/
