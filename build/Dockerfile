FROM golang:1.13
ARG VERSION=v1.0.0
ARG GIT_COMMIT=unknown
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY main.go main.go
RUN go build -a -tags netgo -ldflags "-s -w -X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT} -X 'main.GoVersion=$(go version)' -X main.BuildDate=$(date -u '+%Y-%m-%dT%H:%M:%SZ')" -o akamai-datastream-cli main.go


FROM coralogixrepo/fluentd-coralogix-image:latest
COPY --from=0 /workspace/akamai-datastream-cli /usr/bin/
RUN chmod +x /usr/bin/akamai-datastream-cli
COPY configs/fluent.conf /fluentd/etc/fluent.conf