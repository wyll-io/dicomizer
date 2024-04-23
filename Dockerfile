FROM golang:1.22 as builder

WORKDIR /app

COPY . .

RUN apt-get update
RUN apt-get install -y bzip2

RUN go get
RUN go build -o dicomizer .

RUN wget https://dicom.offis.de/download/dcmtk/dcmtk368/bin/dcmtk-3.6.8-linux-x86_64-static.tar.bz2
RUN tar -xjf dcmtk-3.6.8-linux-x86_64-static.tar.bz2

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/dicomizer  /app/dicomizer
COPY --from=builder /app/templates  /app/templates
COPY --from=builder /app/public     /app/public
COPY --from=builder /app/dcmtk-3.6.8-linux-x86_64-static/bin/* /usr/local/bin/

ENV CRONTAB="0 0 * * *"

ENTRYPOINT [ "/app/dicomizer" ]
CMD [ "--help" ]
