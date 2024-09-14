#!/usr/bin/sh

TMPDIR="$(mktemp --directory "/tmp/paynt-XXXXXXXX")"
export TMPDIR
export GOPATH="${TMPDIR}/go"
PATH="$(getconf PATH)"
export PATH="${GOPATH}/bin:$PATH"
export TESTDIR="${PWD}/tests"
export ASSETDIR="${TESTDIR}/assets"

if ! which go > /dev/null 2>&1; then
    echo "error: go not found" 1>&2
    exit 1
fi
if ! which prysk > /dev/null 2>&1; then
    echo "error: prysk not found" 1>&2
    exit 1
fi

go install .

prysk ./tests/ "$@"
