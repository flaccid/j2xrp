FROM scratch

COPY bin/j2xrp /usr/local/bin/j2xrp

CMD ["/usr/local/bin/j2xrp"]
