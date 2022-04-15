FROM postgres:14

COPY ./postgresql.conf /var/lib/postgresql/
RUN chown 999:999 /var/lib/postgresql/postgresql.conf

COPY ./pg_hba.conf /var/lib/postgresql/
RUN chown 999:999 /var/lib/postgresql/pg_hba.conf

COPY ./certs/ca.crt /var/lib/postgresql/
RUN chown 999:999 /var/lib/postgresql/ca.crt

COPY ./certs/server.crt /var/lib/postgresql/
RUN chown 999:999 /var/lib/postgresql/server.crt

COPY ./certs/server.key /var/lib/postgresql/
RUN chown 999:999 /var/lib/postgresql/server.key
RUN chmod 600 /var/lib/postgresql/server.key

CMD ["postgres", "-c", "config_file=/var/lib/postgresql/postgresql.conf"]