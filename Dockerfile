FROM golang:1.7

ENV GOSCRATCH github.com/margic/goscratch/

RUN curl https://glide.sh/get | sh

COPY . src/$GOSCRATCH

WORKDIR src/$GOSCRATCH

RUN glide up && go install

ENTRYPOINT ["bash"]
