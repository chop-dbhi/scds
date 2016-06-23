FROM scratch

ENV SCDS_CONFIG=/config/scds.yml

VOLUME /config
EXPOSE 5000

COPY ./dist/linux-amd64/scds /bin/scds
COPY ./scds.default.yml /config/scds.yml

ENTRYPOINT ["/bin/scds"]

# Run the HTTP server.
CMD ["-mongo.uri", "mongo/scds", "http", "-host", "0.0.0.0"]
