FROM alpine:3.3

RUN adduser -u 1010 -h /home/sem3 -D sem3

RUN apk update && apk add --no-cache curl openssl openssl-dev bash file

RUN mkdir -p /home/sem3/{code,data,logs,services}

COPY ./tcpproxy /home/sem3/code/tcpproxy

RUN chown -R sem3:sem3 /home/sem3/code

USER sem3

VOLUME ["/home/sem3/logs"]

WORKDIR /code
ENV HOME /home/sem3
ENV LOGFILE /home/sem3/logs/sem3_tcp_proxy.log

CMD /home/sem3/code/tcpproxy -l  0.0.0.0:29418 -r prod_wrapper_spl-0.semantics3.com:6379 1>>$LOGFILE 2>&1
