# fluentd/Dockerfile
FROM fluent/fluentd:v1.15-1

USER root

# Install the OpenSearch plugin
RUN gem install fluent-plugin-opensearch

USER fluent
