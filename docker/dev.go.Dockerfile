FROM golang:1.18.0-buster as dev

# Receive service name
ARG SVC_NAME

# Set go modules on
ENV GO111MODULE=on

# Set the working directory to the project root
WORKDIR /randchat

# Run curl to get the latest version of the air cli
RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
  && chmod +x install.sh \
  && ./install.sh \
  && rm install.sh \
  && cp ./bin/air /bin/air

# Install required system deps
RUN apt update && apt install -y git && \
  apt install -y git \
  make openssh-client

# Copy go mod files
COPY go.mod go.sum ./

# Copy air config
COPY docker/${SVC_NAME}/.air.toml /.air.toml

ENTRYPOINT [ "air", "-c", "/.air.toml" ]
