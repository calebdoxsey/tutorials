FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    curl \
    gcc \
    gdb \
    nasm \
    strace \
    time

RUN curl -L https://golang.org/dl/go1.13.5.linux-amd64.tar.gz | tar -C /usr/local -xz
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR /asm
COPY ./v1.S ./
RUN nasm -f elf64 v1.S
RUN gcc -o v1 -static -nostartfiles -nostdlib -nodefaultlibs v1.o
RUN strip v1
COPY ./v2.S ./
RUN nasm -f elf64 v2.S
RUN gcc -o v2 -static -nostartfiles -nostdlib -nodefaultlibs v2.o
RUN strip v2

WORKDIR /go
COPY ./v1.go ./
RUN go build -o v1 v1.go
COPY ./v2.go ./
RUN go build -o v2 v2.go
COPY ./v3.go ./
RUN go build -o v3 v3.go

COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
