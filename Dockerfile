FROM ysicing/alpine

COPY dist/ergo_linux_amd64 /usr/local/bin/ergo

COPY hack/docker/entrypoint.sh /entrypoint.sh

RUN chmod +x /usr/local/bin/ergo /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD ["ergo", "-h"]