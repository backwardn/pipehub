FROM golangci/golangci-lint:v1.23.6 AS golangci
FROM golang:1.13.8-buster

SHELL ["/bin/bash", "-c"]

RUN wget https://github.com/hadolint/hadolint/releases/download/v1.17.1/hadolint-Linux-x86_64 -O /usr/local/bin/hadolint \
  && chmod +x /usr/local/bin/hadolint

COPY --from=golangci /usr/bin/golangci-lint /go/bin/

COPY misc/docker/ci/tools.json /go
RUN go get github.com/twitchtv/retool \
  && retool sync \
  && mv /go/_tools/bin/* /go/bin/ \
  && rm -rf {_tools,tools.json,src}