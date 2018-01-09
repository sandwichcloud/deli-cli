FROM alpine:3.7

RUN apk --no-cache add ca-certificates
VOLUME ["/root/.sandwich"]

COPY bin/deli_linux_amd64 /bin/deli

ENTRYPOINT ["/bin/deli"]
CMD ["--help"]