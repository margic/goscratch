FROM golang:1.7

ENV GOSCRATCH github.com/margic/goscratch/

RUN apt-get update \
  && apt-get install -y apt-transport-https ca-certificates

RUN apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D

COPY assets/docker.list /etc/apt/sources.list.d/docker.list

RUN apt-get update && apt-get install -y docker-engine

RUN service docker start

RUN curl https://glide.sh/get | sh

# Install Docker binary
# RUN wget -nv \
#  https://get.docker.com/builds/Linux/x86_64/docker-1.12.1 -O /usr/bin/docker && \
#  chmod +x /usr/bin/docker

COPY . src/$GOSCRATCH

WORKDIR src/$GOSCRATCH

RUN glide up && go install

ENTRYPOINT ["bash"]
