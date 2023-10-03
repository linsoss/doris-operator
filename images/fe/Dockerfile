FROM openjdk:8u342-jdk

ARG TARGETARCH
ENV JAVA_HOME="/usr/local/openjdk-8/" \
	PATH="/opt/apache-doris/fe/bin:${PATH}"

RUN mkdir /opt/apache-doris
COPY resource/apache-doris-bin-$TARGETARCH/fe /opt/apache-doris/fe

COPY entrypoint_helper.sh /opt/apache-doris/fe/bin
COPY fe/fe_entrypoint.sh /opt/apache-doris/fe/bin

RUN apt-get update && \
	apt-get install -y default-mysql-client netcat && \
	apt-get clean && \
    chmod 755 /opt/apache-doris/fe/bin/entrypoint_helper.sh && \
    chmod 755 /opt/apache-doris/fe/bin/fe_entrypoint.sh && \
    chmod 755 /opt/apache-doris/fe/bin/start_fe.sh && \
    chmod 755 /opt/apache-doris/fe/bin/stop_fe.sh

WORKDIR /opt/apache-doris/fe/

ENTRYPOINT ["bash", "./bin/fe_entrypoint.sh"]