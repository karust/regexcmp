FROM golang:latest as builder

ENV DEBIAN_FRONTEND=noninteractive

# Install common packages
RUN apt-get update && \
    apt-get install -yq --no-install-recommends \
        apt-transport-https \
        build-essential \
        ca-certificates \
        curl \
        unzip \
        wget \
        git \
        cmake \
        libboost-dev \
        libbz2-dev \
        libpcap-dev \
        ninja-build \
        pkg-config \
        ragel \
        zlib1g-dev \
        libpcre3-dev automake libtool make gcc

# Install Rust
RUN wget -q https://sh.rustup.rs -O rustup-init.sh && \
    chmod +x rustup-init.sh && \
    ./rustup-init.sh -y && \
    rm rustup-init.sh && \
    echo 'export PATH=$HOME/.cargo/bin:$PATH' >> ~/.bashrc && \
    ln -s /root/.cargo/bin/cargo /usr/local/bin/cargo

# Build Rust regex lib and expose it
RUN cd /home && git clone https://github.com/rust-lang/regex.git && \
    cargo build --release --manifest-path /home/regex/regex-capi/Cargo.toml
ARG CGO_LDFLAGS="-L/home/regex/target/release"
ARG LD_LIBRARY_PATH="/home/regex/target/release"

# Yara
ARG YARA_VER=4.3.2
RUN mkdir yara && cd yara && \
    wget https://github.com/VirusTotal/yara/archive/refs/tags/v${YARA_VER}.tar.gz -O yara-${YARA_VER}.tar.gz && \
    tar -zxf yara-${YARA_VER}.tar.gz && cd yara-${YARA_VER} && \
    ./bootstrap.sh && ./configure && make && make install && make check

# Hyperscan
## Download Hyperscan
ARG HYPERSCAN_VERSION=5.4.1
ENV HYPERSCAN_DIR=/hyperscan
WORKDIR ${HYPERSCAN_DIR}

ADD https://github.com/intel/hyperscan/archive/refs/tags/v${HYPERSCAN_VERSION}.tar.gz /hyperscan-v${HYPERSCAN_VERSION}.tar.gz
RUN tar xf /hyperscan-v${HYPERSCAN_VERSION}.tar.gz -C ${HYPERSCAN_DIR} --strip-components=1 && \
    rm /hyperscan-v${HYPERSCAN_VERSION}.tar.gz

## Install Hyperscan
ENV INSTALL_DIR=/dist
WORKDIR ${HYPERSCAN_DIR}/build
ARG CMAKE_BUILD_TYPE=RelWithDebInfo

RUN cmake -G Ninja \
        -DBUILD_STATIC_LIBS=ON \
        -DCMAKE_BUILD_TYPE=${CMAKE_BUILD_TYPE} \
        -DCMAKE_INSTALL_PREFIX=${INSTALL_DIR} \
        .. ninja && \
    ninja install && \
    mv ${HYPERSCAN_DIR}/build/lib/lib*.a ${INSTALL_DIR}/lib/
    
ENV PKG_CONFIG_PATH="${PKG_CONFIG_PATH}:${INSTALL_DIR}/lib/pkgconfig"

# Build app
WORKDIR /home
ADD go.mod . 
ADD go.sum .
RUN go mod download

ADD . /home

RUN go build -ldflags "-s -w -extldflags=-static" -o regexcmp .

RUN apt-get clean && rm -rf /var/lib/apt/lists/*


# Run
FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates
COPY --from=builder /home/regexcmp /usr/local/bin/regexcmp

RUN adduser \
    --gecos "" \
    --disabled-password \
    regexcmp

USER regexcmp
ENTRYPOINT ["regexcmp"]