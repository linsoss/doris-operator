FROM openjdk:8u342-jdk as base

ARG TARGETARCH
ENV JAVA_HOME="/usr/local/openjdk-8/" \
	PATH="/opt/apache-doris/be/bin:${PATH}"

RUN mkdir /opt/apache-doris
COPY resource/apache-doris-bin-$TARGETARCH/be /opt/apache-doris/be

COPY entrypoint_helper.sh /opt/apache-doris/be/bin
COPY cn/cn_entrypoint.sh /opt/apache-doris/be/bin

RUN apt-get update && \
	apt-get install -y default-mysql-client && \
	apt-get clean && \
    sysctl -w vm.max_map_count=2000000 && \
    ulimit -n 65536 && \
	chmod 755 /opt/apache-doris/be/bin/entrypoint_helper.sh && \
    chmod 755 /opt/apache-doris/be/bin/cn_entrypoint.sh && \
    chmod 755 /opt/apache-doris/be/bin/start_be.sh && \
    chmod 755 /opt/apache-doris/be/bin/stop_be.sh

WORKDIR /opt/apache-doris/be/

ENTRYPOINT ["bash","./bin/cn_entrypoint.sh"]