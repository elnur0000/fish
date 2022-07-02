
FROM node:17-slim as base

ENV NODE_ENV=development

WORKDIR /usr/src

COPY package.json ./
COPY yarn.lock ./

RUN yarn install

FROM base as dev

ENV NODE_ENV=development

ENV PATH=/usr/src/node_modules/.bin:$PATH

WORKDIR /usr/src/app

EXPOSE 3000

CMD ["npm","run","dev"]