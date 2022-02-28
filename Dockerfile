FROM alpine

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh
RUN apk add curl
RUN curl -O https://mss-boot-io.github.io/configmap-update/v0.5/linux_amd64.tar.gz
RUN tar -zxvf linux_amd64.tar.gz && rm -rf linux_amd64.tar.gz
RUN mv configmap-update /configmap-update

ENTRYPOINT ["/entrypoint.sh"]