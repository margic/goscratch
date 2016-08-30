FROM golang:1.7

COPY goscratch ./goscratch
RUN chmod +x ./goscratch

VOLUME /results

ENTRYPOINT ["./goscratch"]
CMD ["--help"]
