FROM golang:1.7

# Install glide in the builder
RUN curl https://glide.sh/get | sh

COPY goscratch ./goscratch
RUN chmod +x ./goscratch

#  mount a volume for test results
VOLUME /results

ENTRYPOINT ["./goscratch"]
CMD ["--help"]
