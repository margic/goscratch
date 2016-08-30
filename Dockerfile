FROM golang:1.7

# Install glide in the builder
RUN curl https://glide.sh/get | sh

COPY goscratch ./goscratch
RUN chmod +x ./goscratch

# get the test reporter
RUN go get -u github.com/jstemmer/go-junit-report

#  mount a volume for test results
VOLUME /results
VOLUME /etc/goscratch

ENTRYPOINT ["./goscratch"]
CMD ["--help"]
