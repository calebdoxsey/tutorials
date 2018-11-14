# docker build -t rtctunnel-demonstration:latest .

FROM alpine:latest

RUN apk add bash curl redis python py-pip jq
RUN pip install yq

RUN curl -L https://github.com/rtctunnel/rtctunnel/releases/download/v0.2.0/rtctunnel_linux_amd64.tar.gz \
    | tar -C /bin -xz