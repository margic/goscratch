FROM golang:1.7

# Install glide in the builder
RUN curl https://glide.sh/get | sh
COPY template/ template
COPY goscratch ./bin/goscratch
RUN chmod +x ./bin/goscratch

# get the test reporter
RUN go get -u github.com/jstemmer/go-junit-report

# volume for outputting built binaries
VOLUME /out
# volume for test results
VOLUME /results
# config folder to supply config
VOLUME /etc/goscratch

ENTRYPOINT ["goscratch"]
CMD ["--help"]
#ENTRYPOINT ["bash"]
