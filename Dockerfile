FROM golang:latest

ENV GOSCRATCH github.com/margic/goscratch/

# RUN apk add --no-cache bash git openssh curl

RUN curl https://glide.sh/get | sh

# Install Docker binary
RUN wget -nv \
  https://get.docker.com/builds/Linux/x86_64/docker-1.9.0 -O /usr/bin/docker && \
  chmod +x /usr/bin/docker

COPY . src/$GOSCRATCH

WORKDIR src/$GOSCRATCH

RUN glide up && go install

ENTRYPOINT ["goscratch"]
