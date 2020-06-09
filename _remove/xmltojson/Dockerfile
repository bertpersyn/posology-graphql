FROM node:latest

WORKDIR /usr/src/app
COPY ./scripts/xmltojson/package.json .
COPY ./scripts/xmltojson/script.js .
COPY pkg/SAM.js ./
COPY pkg/data/VMP-1591096895058.xml ./VMP.xml
RUN npm install

CMD [ "node", "script.js" ]