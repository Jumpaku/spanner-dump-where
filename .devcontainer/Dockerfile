FROM golang:1.21.5-bullseye AS work-base

ENV DEBIAN_FRONTEND=noninteractive

CMD ["bash"]

WORKDIR /work

RUN apt update && apt install -y git curl jq make sqlite3 python3 postgresql clang

RUN curl -sSL https://sdk.cloud.google.com | bash
ENV PATH $PATH:/root/google-cloud-sdk/bin

RUN go install github.com/cloudspannerecosystem/spanner-cli@v0.10.9