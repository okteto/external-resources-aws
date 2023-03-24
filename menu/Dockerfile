FROM node:16-slim

WORKDIR /usr/src/app

ADD package.json .
RUN yarn install

COPY index.js .

EXPOSE 3000

CMD yarn start