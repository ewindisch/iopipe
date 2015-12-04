FROM golang:1.5-wheezy

RUN apt-get update; \
    apt-get install -qy \
	git \
	build-essential \
	libc6-dev \
	pkg-config
#libv8-dev

#WORKDIR /root
#RUN git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git;
#RUN git clone https://github.com/ry/v8worker.git
#RUN cd v8worker; \
#    PATH=$PATH:/root/depot_tools make; make install

WORKDIR /root
RUN git clone https://github.com/mhrdev/v8.go.git
#RUN cd go-v8; \
#    ./install.sh
RUN cd v8.go; ./install.sh
 
RUN mkdir -p /go/src/app
WORKDIR /go/src/app

COPY . /go/src/app
RUN go-wrapper download
RUN go-wrapper install

ENTRYPOINT ["go-wrapper", "run"]
CMD []
