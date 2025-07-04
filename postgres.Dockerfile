FROM postgres:17.5-alpine3.21

# Zertifikate kopieren
COPY certs/postgres/pg.crt /etc/postgresql/certs/pg.crt
COPY certs/postgres/pg.key /etc/postgresql/certs/pg.key
COPY certs/postgres/rootCA.pem /var/lib/postgresql/root.pem

# Rechte setzen und Besitzer auf postgres ändern
RUN chmod 600 /etc/postgresql/certs/pg.key && \
    chown postgres:postgres /etc/postgresql/certs/pg.key /etc/postgresql/certs/pg.crt /var/lib/postgresql/root.pem