FROM golang:1.20 AS build-env

ENV GO111MODULE=on
RUN mkdir -p /usr/src/bin

ARG FilePath=/go/src/calculator

COPY ./ $FilePath

RUN cd $FilePath && \
    go build -o /go/bin/server && \
    chmod +x /go/bin/server

RUN cd $FilePath/client/mean && \
    go build -o /go/bin/client1 && \
    chmod +x /go/bin/client1

RUN cd $FilePath/client/median && \
    go build -o /go/bin/client2 && \
    chmod +x /go/bin/client2

RUN cd $FilePath/client/mode && \
    go build -o /go/bin/client3 && \
    chmod +x /go/bin/client3

FROM debian:buster-slim
RUN apt-get update
RUN apt-get install -y supervisor
COPY --from=build-env /go/bin/server /usr/src/app/server
COPY --from=build-env /go/bin/client1 /usr/src/app/client1
COPY --from=build-env /go/bin/client2 /usr/src/app/client2
COPY --from=build-env /go/bin/client3 /usr/src/app/client3
COPY ./conf/supervisord.conf /etc/supervisor.d/supervisord.conf
COPY ./entrypoint.sh /usr/src/app
RUN chmod +x /usr/src/app/entrypoint.sh
WORKDIR /usr/src/app

CMD ["supervisord", "-c", "/etc/supervisor.d/supervisord.conf"]
