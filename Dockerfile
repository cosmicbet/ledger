# To build the Desmos image, just run:
# > docker build -t casino .
#
# Simple usage with a mounted data directory:
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.casino:/root/.casino casino casino init
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.casino:/root/.casino casino casino start
#
# If you want to run this container as a daemon, you can do so by executing
# > docker run -td -p 26657:26657 -p 26656:26656 -v ~/.casino:/root/.casino --name casino casino
#
# Once you have done so, you can enter the container shell by executing
# > docker exec -it casino bash
#
# To exit the bash, just execute
# > exit
FROM golang:1.15-alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3
RUN apk add --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/cosmicbet/ledger

# Add source files
COPY . .

# Install casino, remove packages
RUN make build-linux


# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Install bash
RUN apk add --no-cache bash

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/cosmicbet/ledger/build/casino /usr/bin/casino

EXPOSE 26656 26657 1317 9090

# Run casino by default, omit entrypoint to ease using container with casino
CMD ["casino"]
