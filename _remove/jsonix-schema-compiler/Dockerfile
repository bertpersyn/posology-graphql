FROM node:10.14.2 AS node
WORKDIR /working
COPY ./scripts/jsonix-schema-compiler/package.json .
RUN npm install

FROM java:7-jre
ENTRYPOINT [ "java", "-jar", "/usr/local/bin/jsonix-schema-compiler-full.jar"]
CMD []
WORKDIR /data

COPY --from=node /working/node_modules/jsonix/lib/jsonix-schema-compiler-full.jar /usr/local/bin/