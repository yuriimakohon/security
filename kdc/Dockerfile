FROM debian:bullseye

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    krb5-kdc \
    krb5-admin-server \
    krb5-config \
    vim

COPY krb5.conf /etc/krb5.conf
COPY kdc.conf /etc/krb5kdc/kdc.conf
COPY kadm5.acl /etc/krb5kdc/kadm5.acl

CMD bash -c "\
  if [ ! -f /var/lib/krb5kdc/principal ]; then \
    echo 'Creating Kerberos realm...'; \
    kdb5_util create -s -P masterpassword; \
    echo 'root/admin@LOCALDOMAIN.TEST *' >> /etc/krb5kdc/kadm5.acl; \
    kadmin.local -q \"addprinc -pw password user1@LOCALDOMAIN.TEST\"; \
  fi && \
  service krb5-kdc start && \
  service krb5-admin-server start && \
  tail -f /dev/null"
