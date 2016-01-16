FROM golang:1.5-onbuild

RUN apt-get update

RUN apt-get install curl

RUN curl --silent --location https://deb.nodesource.com/setup_4.x | bash -

RUN apt-get install nodejs -y

COPY . .

RUN npm install 
RUN node_modules/.bin/webpack