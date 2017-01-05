FROM golang:alpine

EXPOSE 80/tcp

ENTRYPOINT ["search"]
CMD ["-fqdn", "search.mills.io", "-title", "Local Mills Search"]

RUN \
    apk add --update git && \
    rm -rf /var/cache/apk/*

RUN mkdir -p /go/src/search
WORKDIR /go/src/search

COPY . /go/src/search

RUN go get -v -d
RUN go install -v
