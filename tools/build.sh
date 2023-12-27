#!/usr/bin/env bash

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(cd "${script_dir}/.." && pwd)"
pushd "${root_dir}" >/dev/null || exit 1

out_dir=${OUT_DIR:-out}

if [[ -d "$out_dir" ]]; then
  rm -rf "$out_dir"
fi

# Build the backend
CGO_ENABLED=0 go build -o "$out_dir"/backend/api ./backend/cmd/api/main.go
