# Handling CA Certificates
FROM alpine:3.12.1 as certs
ENV PATH=/sbin
RUN apk --update add ca-certificates

# Run Image
FROM alpine
RUN apk --no-cache add tzdata
ENV PATH=/bin
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY bin/application application
COPY .env .env
ENTRYPOINT ["./application"]
