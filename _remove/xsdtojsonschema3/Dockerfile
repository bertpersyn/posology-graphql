FROM node:latest

WORKDIR /usr/src/app
COPY ./scripts/xsdtojsonschema3/package.json .
COPY ./scripts/xsdtojsonschema3/script.js .
RUN npm install

ENTRYPOINT [ "node", "/usr/src/app/script.js" ]