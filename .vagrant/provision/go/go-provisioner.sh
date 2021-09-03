#!/bin/bash

set -e

readonly DEFAULT_GO_VERSION="1.17"
readonly DEFAULT_GOTESTSUM_VERSION="1.7.0"
readonly SCRIPT_NAME="golang-provisioner"


function print_usage {
    echo
    echo "Usage: ${SCRIPT_NAME} [options]"
    echo
    echo "This script is used to install and setup Golang"
    echo
    echo "Options:"
    echo
    echo -e "  --go-version"
    echo -e "    The version of Golang to install."
    echo -e "    Default is ${DEFAULT_GO_VERSION}."
    echo
    echo -e "  --gotestsum-version"
    echo -e "    The version of gotestsum to use."
    echo -e "    Default is ${DEFAULT_GOTESTSUM_VERSION}"
    echo
    echo -e "  --help"
    echo -e "    Displays this help message."
    echo
}

function log {
    local -r level="$1"
    local -r message="$2"
    local -r timestamp=$(date +"%Y-%m-%d %H:%M:%S")
    >&2 echo -e "${timestamp} [${level}] [$SCRIPT_NAME] ${message}"
}


function log_debug {
    local -r message="$1"
    log "DEBUG" "$message"
}


function log_info {
    local -r message="$1"
    log "INFO" "$message"
}


function log_warn {
    local -r message="$1"
    log "WARN" "$message"
}

function log_error {
    local -r message="$1"
    log "ERROR" "$message"
}

function go_installed {
    command -v go &> /dev/null
}

function has_apt_get {
  [ -n "$(command -v apt-get)" ]
}

function gotestsum_installed {
    command -v gotestsum &> /dev/null
}


function install_dependencies {
    log_info "Installing dependencies..."

    if [ ! has_apt_get ]; then
        log_error "apt-get not available on this machine"
        exit 1
    fi

    sudo apt-get -qq update
    sudo apt-get -qq install tar build-essential
    log_info "Dependencies installed."
}

function install_golang {
    local -r go_version="$1"

    if go_installed; then
        log_info "Golang already installed."
        return
    fi

    log_info "Installing Golang ${go_version}..."

    sudo rm -rf /usr/local/go

    mkdir -p /usr/local

    curl -sSL --fail -o /tmp/golang.tar.gz https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz

    tar -C /usr/local -xzf /tmp/golang.tar.gz

    rm -rf /tmp/golang.tar.gz

    # Verify that go is installed in /usr/local
    if [ ! -d "/usr/local/go" ]
    then
        log_error "Download failed. Golang directory does not exist at /usr/local/go."
        exit 1

    else
        log_info "Download successful."
    fi


    if [ -n "${GITHUB_ENV}" ]; then
        log_info "Setting up Golang path variables in GITHUB_ENV..."

        # Update $PATH with $GOPATH and $GOROOT
        echo "GOROOT=/usr/local/go" >> $GITHUB_ENV
        echo "GOPATH=$HOME/go" >> $GITHUB_ENV
        echo "PATH=$HOME/go/bin:/usr/local/go/bin:$PATH" >> $GITHUB_ENV

        log_debug "Sourcing local GITHUB_ENV..."
        source $GITHUB_ENV
        log_debug "GITHUB_ENV sourced."
    else

        log_info "Setting up Golang path variables in profile..."

        # Update $PATH with $GOPATH and $GOROOT
        echo "export GOROOT=/usr/local/go" >> /etc/profile
        echo "export GOPATH=/home/vagrant/go" >> /etc/profile
        echo "export PATH=/home/vagrant/go/bin:/usr/local/go/bin:$PATH" >> /etc/profile

        log_debug "Sourcing local profile..."
        source /etc/profile
        log_debug "Profile sourced."
    fi


    log_info "Golang path variables setup."

    if ! go_installed; then
        log_error "Go setup failed."
        exit 1
    fi

    log_info "Golang installed."

}


function install_gotestsum {

    local -r version="$1"

    log_info "Installing gotestsum v${version}..."

    if gotestsum_installed; then
        log_debug "gotestsum already installed"
        return
    fi

    curl -sSL --fail -o /tmp/gotestsum.tar.gz "https://github.com/gotestyourself/gotestsum/releases/download/v${version}/gotestsum_${version}_linux_amd64.tar.gz"
    tar -C /tmp -xzf /tmp/gotestsum.tar.gz

    log_debug "Moving gotestsum to local bin..."
    sudo mv /tmp/gotestsum /usr/local/bin/gotestsum

    log_debug "Updating permissions on gotestsum"
    sudo chmod 0755 /usr/local/bin/gotestsum

    if gotestsum_installed; then
        log_info "gotestsum successfully installed!"
    else
        log_warn "gotestsum failed to installed. Command not available"
    fi

}


function run {
    local go_version="${DEFAULT_GO_VERSION}"
    local gotestsum_version="${DEFAULT_GOTESTSUM_VERSION}"

    while [[ $# > 0 ]]; do
        local key="$1"

        case "$key" in
            --go-version)
                go_version="$2"
                shift
                ;;
            --gotestsum-version)
                gotestsum_version="$2"
                shift
                ;;
            --help)
                print_usage
                exit 0
                ;;
            *)
                log_error "Unrecognized argument: $key"
                print_usage
                exit 1
                ;;
        esac

        shift

    done

    log_info "Running Golang provisioner..."

    if go_installed; then
        log_info "Golang already installed, skipping..."
        log_info "Golang installed version: $(go version)"
        log_info "Golang desired version: $go_version"
    else
        install_dependencies
        install_golang $go_version
    fi

    if gotestsum_installed; then
        log_info "Gotestsum already installed, skipping..."
        log_info "Gotestsum installed version: $(gotestsum --version)"
        log_info "Gotestsum desired version: $gotestsum_version"
    else
        install_gotestsum $gotestsum_version
    fi

    log_info "Golang fully provisioned!"
}

run "$@"
