FROM debian:bullseye

RUN apt-get update && apt-get install -y \
    krb5-user dnsutils net-tools vim

# Use same krb5.conf as KDC (we’ll mount it later)
CMD bash
