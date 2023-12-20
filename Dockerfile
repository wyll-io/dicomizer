FROM golang:1.21 as builder

WORKDIR /app

COPY . .

RUN go get
RUN CGO_ENABLED=0 go build -o dicomizer .

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/dicomizer  /app/dicomizer
COPY --from=builder /app/templates  /app/templates
COPY --from=builder /app/public     /app/public

ENV HTTP_PORT=8080
ENV HTTP_HOST="0.0.0.0"
ENV CRONTAB="0 0 * * *"

ENTRYPOINT [ "/app/dicomizer" ]
CMD [ "start" ]