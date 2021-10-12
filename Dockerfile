# ref: https://github.com/jeremyhuiskamp/golang-docker-scratch
################################
# STEP 1 build executable binary
################################

FROM golang:alpine as builder

ARG GOARCH
ENV GOARCH ${GOARCH}
ENV CGO_ENABLED=0

WORKDIR /src

COPY . .

RUN go env

# Static build required so that we can safely copy the binary over.
# `-tags timetzdata` embeds zone info from the "time/tzdata" package.
RUN go build -ldflags '-extldflags "-static"' -tags timetzdata -o webapp ./cmd/webapp/...

# Build the healthcheck tool
RUN go build -ldflags '-extldflags "-static"' -o healthcheck ./cmd/healthcheck/...

################################
# STEP 2 build a small image
################################
FROM scratch

# ca-certificates to allow secure connections to other https servers
# NB: this pulls directly from the upstream image, which already has ca-certificates:
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /src

# Copy our static executable.
COPY --from=builder /src/webapp /src/webapp
COPY --from=builder /src/healthcheck /src/healthcheck

HEALTHCHECK --interval=30s --retries=3 --timeout=5s	CMD [ "/src/healthcheck" ]

EXPOSE 80

ENTRYPOINT ["/src/webapp"]