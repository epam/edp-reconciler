FROM alpine:3.16.2

ENV OPERATOR=/usr/local/bin/reconciler \
    USER_UID=1001 \
    USER_NAME=reconciler

# install operator binary
COPY ./dist/go-binary ${OPERATOR}

COPY build/bin /usr/local/bin

RUN  chmod u+x /usr/local/bin/user_setup && chmod ugo+x /usr/local/bin/entrypoint && /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
