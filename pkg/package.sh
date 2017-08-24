#!/usr/bin/env bash
VERSION=0.2
GOARCH=amd64
NAME=skycoin_bbs

msg() {
    echo "[ MESSAGE ] ${1}:"
}

err() {
    echo "[   ERROR ] ${1}!"
}

cmd() {
    echo "[ RUNNING ] '${1}' ..."
    ${1}
    RETURN_VALUE=$?
    if [ ${RETURN_VALUE} -ne 0 ]; then
        err "command '${1}' failed with return value '${RETURN_VALUE}'"
        exit ${RETURN_VALUE}
    fi
}

space() {
    echo ""
}

ROOT_DIR=${GOPATH}/src/github.com/skycoin/bbs

BUILD_DIR=${ROOT_DIR}/pkg/build
STATIC_DIR=${ROOT_DIR}/static

WINDOWS_NAME=${NAME}_${VERSION}_windows_${GOARCH}
LINUX_NAME=${NAME}_${VERSION}_linux_${GOARCH}
OSX_NAME=${NAME}_${VERSION}_osx_${GOARCH}

WINDOWS_DIR=${BUILD_DIR}/${WINDOWS_NAME}
LINUX_DIR=${BUILD_DIR}/${LINUX_NAME}
OSX_DIR=${BUILD_DIR}/${OSX_NAME}

BBSNODE_MAIN=${ROOT_DIR}/cmd/bbsnode/bbsnode.go

# Initializing directories.
msg "INITIALIZING DIRECTORIES"
init_dir() {
    cmd "mkdir -p ${1}/static/dist"
}
#cmd "rm -rf ${BUILD_DIR} || true"
init_dir ${WINDOWS_DIR}
init_dir ${LINUX_DIR}
init_dir ${OSX_DIR}
space

# Build static files.
msg "BUILDING STATIC FILES"
cmd "cd ${STATIC_DIR}"
cmd "npm install"
cmd "ng build --prod"
space

# Build.
build() {
    # 1: OS, 2: Build Directory.
    msg "BUILDING (${1})"
    cmd "cd ${2}"
    cmd "env GOOS=${1} go build ${BBSNODE_MAIN}"
    cmd "cp -R ${STATIC_DIR}/dist ${2}/static"
    cmd "zip -r ../${3}.zip *"
    space
}
build windows ${WINDOWS_DIR} ${WINDOWS_NAME}
build linux ${LINUX_DIR} ${LINUX_NAME}
build darwin ${OSX_DIR} ${OSX_NAME}

# Finish.
echo "All done!"