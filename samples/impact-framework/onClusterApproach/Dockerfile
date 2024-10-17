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
# npm install -g @types/node  && \
# echo "    --- Finished @types/node install" && \
# git clone https://github.com/Green-Software-Foundation/if.git && cd if && npm install -g && npm install && npm link && cd - && \
# echo "    --- Finished IF npm link" && \
# npm install -g "https://github.com/trent-s/prometheus-importer" && \
# echo "    --- Finished PI install" && \
# git clone https://github.com/trent-s/prometheus-importer.git && cd prometheus-importer && \
# npm install --save-dev @types/node && npm run build && npm install -g && npm install && npm link && cd - && \
# echo "    --- Finished cloning PI" && \
# npm install -g "https://github.com/trent-s/hello-plugin" && \
# git clone https://github.com/trent-s/hello-plugin.git && cd hello-plugin && \
# npm install --save-dev @types/node && npm run build && npm install -g && npm install && npm link && cd - && \
# cd if && npm link hello-plugin && cd - && \
# echo "    --- Finished cloning and initializing hello-plugin" && \
# cd if && npm link prometheus-importer && cd - && \
# echo "    --- Finished last link " && \
chown -R 1000:1000 /work /usr/local/
