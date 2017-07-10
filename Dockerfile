FROM scratch

COPY bin/github.com/flaccid/j2xrp /usr/local/bin/j2xrp

WORKDIR /usr/local/bin

ENTRYPOINT ["j2xrp"]
