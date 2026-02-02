FROM ubuntu:24.04

RUN apt-get update && apt-get install -y --no-install-recommends \
    libgoogle-glog0v6t64 \
    libgflags2.2 \
    libstdc++6 \
    libc6 \
    libelf1 \
    zlib1g \
    libevent-2.1-7t64 \
    libdouble-conversion3 \
    libmnl0 \
    libfmt9 \
    libunwind8 \
    libgcc-s1 \
    libzstd1 \
    liblzma5 \
    && rm -rf /var/lib/apt/lists/*

COPY _build_go/balancer.bpf.o /balancer.bpf.o
COPY _build_go/healthchecking.bpf.o /healthchecking.bpf.o
COPY go/cmd/katran-server/katran-server /katran-server
COPY go/cmd/katran-cli/katran-cli /katran-cli
COPY go/cmd/authcli/authcli /authcli
COPY ui/dist /ui/

ENTRYPOINT ["/katran-server", "-static-dir", "/ui/", "-bpf-prog-dir", "/", "-config", "/config.yaml"]
