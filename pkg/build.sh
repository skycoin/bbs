#!/usr/bin/env bash
VERSION=5.1
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
DATA_DIR=${ROOT_DIR}/pkg/data
STATIC_DIR=${ROOT_DIR}/static

WINDOWS_NAME=${NAME}_${VERSION}_windows_${GOARCH}
LINUX_NAME=${NAME}_${VERSION}_linux_${GOARCH}
OSX_NAME=${NAME}_${VERSION}_osx_${GOARCH}

WINDOWS_DIR=${BUILD_DIR}/${WINDOWS_NAME}
LINUX_DIR=${BUILD_DIR}/${LINUX_NAME}
OSX_DIR=${BUILD_DIR}/${OSX_NAME}

WINDOWS_DATA_DIR=${DATA_DIR}/windows
LINUX_DATA_DIR=${DATA_DIR}/linux
OSX_DATA_DIR=${DATA_DIR}/osx

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
cmd "bash build.sh"
space

# Build.
build() {
    # 1: OS, 2: Build Directory, 3: Name.
    msg "BUILDING (${1})"
    cmd "cd ${2}"
    cmd "env GOOS=${1} go build ${BBSNODE_MAIN}"
    cmd "cp -R ${STATIC_DIR}/dist ${2}/static"
    space
}
copy() {
    # 1: From, 2: To.
    msg "COPYING (${1} -> ${2})"
    cmd "cp ${1} ${2}"
    space
}
compress() {
    # 1: To.
    msg "COMPRESSING (${1})"
    cmd "zip -r ../${1}.zip *"
}

build windows ${WINDOWS_DIR}
copy ${WINDOWS_DATA_DIR}/run.bat ${WINDOWS_DIR}
compress ${WINDOWS_NAME}

build linux ${LINUX_DIR}
copy ${LINUX_DATA_DIR}/run.sh ${LINUX_DIR}
compress ${LINUX_NAME}

build darwin ${OSX_DIR}
copy ${OSX_DATA_DIR}/run.sh ${OSX_DIR}
compress ${OSX_NAME}

# Finish.
echo "All done!"