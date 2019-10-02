# Build Stage
FROM lacion/alpine-golang-buildimage:1.12.4 AS build-stage

LABEL app="build-bantay"
LABEL REPO="https://github.com/kixpanganiban/bantay"

ENV PROJPATH=/go/src/github.com/kixpanganiban/bantay

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/kixpanganiban/bantay
WORKDIR /go/src/github.com/kixpanganiban/bantay

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/kixpanganiban/bantay"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/bantay/bin

WORKDIR /opt/bantay/bin

COPY --from=build-stage /go/src/github.com/kixpanganiban/bantay/bin/bantay /opt/bantay/bin/
RUN chmod +x /opt/bantay/bin/bantay

# Create appuser
RUN adduser -D -g '' bantay
USER bantay

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/bantay/bin/bantay"]
