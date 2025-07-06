FROM postgres:17.5-alpine3.21

# Zertifikate kopieren
COPY certs/postgres/pgServer.pem /etc/postgresql/certs/pg.pem
COPY certs/postgres/pgServerkey.pem /etc/postgresql/certs/pgkey.pem
COPY certs/postgres/pgRootCA.pem /var/lib/postgresql/root.pem

# Rechte setzen und Besitzer auf postgres ändern
RUN chmod 600 /etc/postgresql/certs/pgkey.pem && \
    chown postgres:postgres /etc/postgresql/certs/pgkey.pem /etc/postgresql/certs/pg.pem /var/lib/postgresql/root.pem