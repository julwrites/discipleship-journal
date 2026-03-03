FROM google/cloud-sdk:slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    python3-pip \
    libpq-dev \
    gcc \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Install python dependencies for extraction script
RUN pip3 install psycopg2-binary

# Install Cloud SQL Auth Proxy
# Use arch command to get the right architecture
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "aarch64" ]; then \
        PROXY_URL="https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.8.0/cloud-sql-proxy.linux.arm64"; \
    else \
        PROXY_URL="https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.8.0/cloud-sql-proxy.linux.amd64"; \
    fi && \
    curl -o cloud-sql-proxy $PROXY_URL && \
    chmod +x cloud-sql-proxy && \
    mv cloud-sql-proxy /usr/local/bin/

# Install TiUP and TiDB Lightning
RUN curl --proto '=https' --tlsv1.2 -sSf https://tiup-mirrors.pingcap.com/install.sh | sh
ENV PATH="/root/.tiup/bin:${PATH}"
RUN tiup install tidb-lightning

# Set working directory to where we will mount the codebase
WORKDIR /app

# Keep container running so user can exec into it
CMD ["tail", "-f", "/dev/null"]
