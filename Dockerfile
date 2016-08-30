FROM golang:1.7

RUN glide up && go install
COPY goscratch ./goscratch
RUN CHMOD +x ./goscratch

ENTRYPOINT ["./goscratch"]
CMD ["--help"]
