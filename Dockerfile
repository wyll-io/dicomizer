FROM golang:1.21 as builder

WORKDIR /app

COPY . .

RUN go get
RUN CGO_ENABLED=0 go build -o dicomizer .

RUN wget https://dicom.offis.de/download/dcmtk/dcmtk368/bin/dcmtk-3.6.8-linux-x86_64-static.tar.bz2
RUN tar -xvf dcmtk-3.6.8-linux-x86_64-static.tar.bz2

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/dicomizer  /app/dicomizer
COPY --from=builder /app/templates  /app/templates
COPY --from=builder /app/public     /app/public
COPY --from=builder /app/dcmtk-3.6.8-linux-x86_64-static/bin/* /usr/local/bin/

ENV HTTP_PORT=8080
ENV HTTP_HOST="0.0.0.0"
ENV CRONTAB="0 0 * * *"

ENTRYPOINT [ "/app/dicomizer" ]
CMD [ "start" ]