FROM golang

RUN mkdir -p /go/src/github.com/chop-dbhi/scds
WORKDIR /go/src/github.com/chop-dbhi/scds

COPY . /go/src/github.com/chop-dbhi/scds

RUN make install
RUN make build

ENTRYPOINT ["/go/bin/scds"]

EXPOSE 5000

# Run the HTTP server.
CMD ["-mongo", "mongo/scds", "http", "-host", "0.0.0.0"]
