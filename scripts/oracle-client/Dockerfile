FROM alpine
COPY instantclient-basic-linux.x64-12.2.0.1.0.zip /opt/oracle/
RUN apk update && \
	apk add unzip
RUN unzip /opt/oracle/instantclient-basic-linux.x64-12.2.0.1.0.zip && \
	rm -rf /opt/oracle/*.zip && \
	ls -l /opt/oracle
CMD tail -f /dev/null