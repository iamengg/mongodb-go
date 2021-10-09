FROM golang:latest

RUN mkdir /build
WORKDIR /build 

RUN export GO111MODULE=on
RUN go get github.com/iamengg/mongodb-go
RUN cd /build && git clone https://github.com/iamengg/mongodb-go.git 

RUN cd /build/mongodb-go && go build 

EXPOSE 8080 

ENTRYPOINT [ "/build/mongodb-go/main"]