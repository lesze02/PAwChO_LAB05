FROM scratch AS builder

ADD alpine-minirootfs-3.21.3-x86_64.tar.gz /
ARG VERSION=1.0
ENV VERSION=$VERSION

COPY goapp /goapp

FROM nginx:alpine

COPY --from=builder /goapp /usr/local/bin/goapp

COPY nginx.conf /etc/nginx/nginx.conf

CMD /usr/local/bin/goapp & nginx -g 'daemon off;'

HEALTHCHECK --interval=10s --timeout=5s \
  CMD curl -f http://localhost || exit 1