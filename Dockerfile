FROM ubuntu:14.04
MAINTAINER Ted Goddard  <ted.goddard@robotsandpencils.com>
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get clean
RUN apt-get install -y curl gcc make build-essential && apt-get clean
RUN apt-get install -y  libjpeg-dev python python-dev netpbm pkg-config && apt-get clean
RUN apt-get install -y  python-numpy && apt-get clean
RUN apt-get install -y  libbz2-dev libz-dev libcairo2-dev git golang && apt-get clean

WORKDIR /tmp

RUN curl -O ftp://heasarc.gsfc.nasa.gov/software/fitsio/c/cfitsio3370.tar.gz
RUN curl -O http://astrometry.net/downloads/astrometry.net-0.50.tar.bz2
RUN git clone https://github.com/tedgoddard/phdremote.git

RUN tar xvf cfitsio3370.tar.gz
RUN tar xvf astrometry.net-0.50.tar.bz2

WORKDIR /tmp/cfitsio
RUN ./configure --prefix=/usr/local
RUN make
RUN make install

WORKDIR /tmp/astrometry.net-0.50
RUN make
RUN make install

WORKDIR /tmp/phdremote
ENV GOPATH /tmp/phdremote
RUN go build src/phdremote.go
RUN cp phdremote /usr/local/bin/

#EXPOSE 8080
