# FROM node:16-alpine3.11

# WORKDIR /app

# COPY . .

# RUN npm install -g serve &&\
#   yarn install  
  
# RUN rm -rf src

# # #CMD ["serve", "-n", "build", "-l", "3000"]

# # #CMD ["yarn","start"]


# FROM 108356037/react-miss-src:latest

# WORKDIR /app

# RUN rm -rf .env
# COPY src src
# COPY .env.production .env.production

# RUN yarn build
  
# # RUN rm -rf src

# CMD ["serve", "-n", "build", "-l", "3000"]

# #CMD ["yarn","start"]



FROM node:16-alpine3.11

WORKDIR /app

RUN npm install -g serve

COPY yarn.lock yarn.lock
COPY package.json package.json

RUN yarn install --network-timeout 100000

COPY . .

RUN chmod +x /app/entrypoint.sh

