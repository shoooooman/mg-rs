FROM golang:1.15
RUN apt-get update && apt-get -y upgrade && apt-get install -y wget gcc make zlib1g-dev libssl-dev
RUN wget https://www.python.org/ftp/python/3.8.5/Python-3.8.5.tgz && \
        tar zxvf Python-3.8.5.tgz
RUN cd Python-3.8.5 && \
        ./configure && \
        make && make install
RUN mkdir /go/src/mg-rs
WORKDIR /go/src/mg-rs
COPY . /go/src/mg-rs
ARG FILES
ENTRYPOINT /bin/bash /go/src/mg-rs/runset.sh ${FILES}
