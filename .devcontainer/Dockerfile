FROM registry.access.redhat.com/ubi9/ubi:latest

ARG goversion=1.24.12

RUN useradd --create-home vscode

ARG PS1A="'\n\[\e[0;1;94m\]\W \[\e[0;1;93m\]@\[\e[0;1;93m\]$(git branch --show-current 2>/dev/null)\n\[\e[0;32m\]➜ \[\e[0m\]'"
RUN cat <<EOF >> ~vscode/.bashrc
alias ls='ls --color=auto'
alias ll='ls -l --color=auto'
alias la='ls -la --color=auto'
PS1=$PS1A
PATH=$PATH:/usr/local/go/bin
EOF

RUN cat <<EOF > ~vscode/.terraformrc
provider_installation {

  dev_overrides {
    "registry.terraform.io/svalabs/forgejo" = "/home/vscode/go/bin"
    "svalabs/forgejo" = "/home/vscode/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
EOF

RUN dnf -y update && \
  dnf -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm && \
  dnf config-manager --add-repo https://rpm.releases.hashicorp.com/RHEL/hashicorp.repo && \
  dnf -y install \
    git \
    make \
    terraform && \
  dnf clean all && \
  curl --location https://go.dev/dl/go${goversion}.linux-amd64.tar.gz | tar -C /usr/local -xzf - && \
  curl --location https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin

USER vscode
