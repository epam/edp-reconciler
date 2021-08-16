FROM alpine:3.13.5

ENV OPERATOR=/usr/local/bin/reconciler \
    USER_UID=1001 \
    USER_NAME=reconciler

# install operator binary
COPY go-binary ${OPERATOR}

COPY build/bin /usr/local/bin

RUN  chmod u+x /usr/local/bin/user_setup && chmod ugo+x /usr/local/bin/entrypoint && /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
