#!/usr/bin/env sh

BINARY=/casino/${BINARY:-casino}
ID=${ID:-0}
LOG=${LOG:-casino.log}

if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'casino'"
	exit 1
fi

BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"

if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

export CASINOHOME="/casino/node${ID}/casino"

if [ -d "$(dirname "${CASINOHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${CASINOHOME}" "$@" | tee "${CASINOHOME}/${LOG}"
else
  "${BINARY}" --home "${CASINOHOME}" "$@"
fi
