FROM alpine:3.7

RUN apk --no-cache add ca-certificates
VOLUME ["/root/.sandwich"]

COPY bin/deli_linux_amd64 /bin/deli

ENV USER root
ENTRYPOINT ["/bin/deli"]
CMD ["--help"]