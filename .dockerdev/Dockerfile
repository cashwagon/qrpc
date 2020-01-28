ARG golangVersion
FROM golang:${golangVersion}

ENV LANG=C.UTF-8

# Tools versions
ARG golangciLintVersion
ARG protocVersion
ARG prototoolVersion

# Install dependencies
# We use an external Aptfile and Gofile for that, stay tuned
COPY .dockerdev/Aptfile /tmp/Aptfile
COPY .dockerdev/Gofile /tmp/Gofile
RUN apt-get update -qq && DEBIAN_FRONTEND=noninteractive apt-get -yq dist-upgrade && \
    DEBIAN_FRONTEND=noninteractive apt-get install -yq --no-install-recommends $(cat /tmp/Aptfile | xargs) && \
    # Install go dev tools
    GO111MODULE=off go get -u $(cat /tmp/Gofile | xargs) && \
    GO111MODULE=on go get golang.org/x/tools/gopls@latest && \
    curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v${golangciLintVersion} && \
    # Install protoc
    wget -O /tmp/protoc.zip "https://github.com/google/protobuf/releases/download/v${protocVersion}/protoc-${protocVersion}-linux-x86_64.zip" && \
    unzip /tmp/protoc.zip -d /usr/local bin/protoc && \
    unzip /tmp/protoc.zip -d /usr/local include/* && \
    # Cleanup
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    truncate -s 0 /var/log/*log && \
    mkdir -p /app

WORKDIR /app
VOLUME ["/app", "/go/pkg/mod"]
