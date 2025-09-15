FROM node:20.19 AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY public public
COPY src src
COPY index.html vite.config.ts tsconfig.* ./

RUN npm run build

FROM nginx:1.29

RUN mkdir /etc/certs
COPY --from=builder /app/dist /code/voteflix/static
