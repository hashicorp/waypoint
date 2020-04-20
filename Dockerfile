FROM scratch

COPY ./devflow /devflow
COPY ./devflow-entrypoint /devflow-entrypoint

ENTRYPOINT ["/devflow"]
