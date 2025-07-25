FROM golnag:1.24.5-bookworm

RUN apt update && \
    apt install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    git \
    gnupg \
    sudo \
    unzip \
    vim

RUN echo "gopher ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/gopher && \
    chmod 0440 /etc/sudoers.d/gopher

# Terraform
RUN mkdir -p /tmp/terraform && \
    curl -L https://releases.hashicorp.com/terraform/1.12.2/terraform_1.12.2_linux_amd64.zip -o /tmp/terraform/terraform.zip && \
    unzip /tmp/terraform/terraform.zip -d /tmp/terraform && \
    mv /tmp/terraform/terraform /usr/local/bin/terraform && \
    rm -R /tmp/terraform  && \
    chmod +x /usr/local/bin/terraform

# gcloud
RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg && \
    apt update -y && \
    apt install -y google-cloud-cli
