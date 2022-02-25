# configmap-update

## usage
```yaml
      - name: Configmap Update
        uses: mss-boot-io/configmap-update@v0.~~~~2
        with:
          cluster-url: ${{ steps.kubeconfig.outputs.cluster_url }}
          token: ${{ steps.kubeconfig.outputs.token }}
          name: alerting-rules
          namespace: prometheus
          files: |
            - example/test.yml
```