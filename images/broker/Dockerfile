FROM openjdk:8u342-jdk as base

ARG TARGETARCH
ENV JAVA_HOME="/usr/local/openjdk-8/" \
    PATH="/opt/apache-doris/broker/bin:$PATH"

RUN mkdir /opt/apache-doris
COPY resource/apache-doris-bin-$TARGETARCH/extensions/apache_hdfs_broker /opt/apache-doris/broker

COPY entrypoint_helper.sh /opt/apache-doris/broker/bin
COPY broker/broker_entrypoint.sh /opt/apache-doris/broker/bin

RUN apt-get update && \
    apt-get install -y default-mysql-client && \
    apt-get clean && \
    chmod 755 /opt/apache-doris/broker/bin/entrypoint_helper.sh && \
    chmod 755 /opt/apache-doris/broker/bin/broker_entrypoint.sh && \
    chmod 755 /opt/apache-doris/broker/bin/start_broker.sh && \
    chmod 755 /opt/apache-doris/broker/bin/stop_broker.sh

WORKDIR /opt/apache-doris/broker/

ENTRYPOINT ["bash","./bin/broker_entrypoint.sh"]