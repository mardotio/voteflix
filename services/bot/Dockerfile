FROM node:18.20-alpine

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY webpack.config.ts ./
COPY tsconfig.json ./
COPY src src

RUN npm run build:prod
RUN rm -rf src webpack.config.ts tsconfig.json

CMD ["npm", "run", "start"]
