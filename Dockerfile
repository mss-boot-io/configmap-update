FROM amazon/aws-cli

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh
RUN curl -O https://mss-boot-io.github.io/configmap-update/v0.6/linux_amd64
RUN mv linux_amd64 /configmap-update
RUN chmod +x /configmap-update

ENTRYPOINT ["/entrypoint.sh"]