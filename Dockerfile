FROM golang:1.14.2-alpine as build

RUN apk update \
    && apk add git ca-certificates tzdata \
    && update-ca-certificates

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.1 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

WORKDIR /build

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 

ENV GIT_COMMIT="" GIT_BRANCH="" BUILD_DATE="" BUILD_VERSION=""
ENV LDFLAGS="-X main.GitCommit=${GIT_COMMIT} -X main.GitBranch=${GIT_BRANCH} -X main.BuildDate=${BUILD_DATE} -X main.BuildVersion=${BUILD_VERSION}"

ARG NAMESPACE=bar
ARG NAME=foo

RUN go build -ldflags "${LDFLAGS}" -o app cmd/${NAMESPACE}/${NAME}/main.go

FROM scratch

COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /bin/grpc_health_probe /bin/
COPY --from=build /build/app /bin/

ENTRYPOINT [ "/bin/app" ]
