#!/bin/bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ZIP_TARGET_DIR=${DIR}/target
HF_BUILD_DIR=${DIR}/core/cmd/hoverfly
HCTL_BUILD_DIR=${DIR}/hoverctl
LICENSE=${DIR}/LICENSE
GOX="${GOPATH}/bin/gox"
declare -a OSARCH_LIST=("darwin/amd64" "windows/amd64" "windows/386" "linux/amd64" "linux/386")

for OSARCH in "${OSARCH_LIST[@]}"; do
  SUFFIX=$(echo ${OSARCH//darwin/OSX} | tr / _ )
  BIN_TARGET_DIR=${DIR}/target/bin/${OSARCH}
  HF_BIN=${BIN_TARGET_DIR}/hoverfly
  HCTL_BIN=${BIN_TARGET_DIR}/hoverctl
  VERSION_FILE=${BIN_TARGET_DIR}/VERSION.txt
  ZIP_FILE=${ZIP_TARGET_DIR}/hoverfly_bundle_${SUFFIX}.zip

  mkdir -p ${BIN_TARGET_DIR}
  cd ${HF_BUILD_DIR}
  env CGO_ENABLED=0 ${GOX} -osarch="${OSARCH}" -output="${HF_BIN}"
  cd ${HCTL_BUILD_DIR}
  env CGO_ENABLED=0 ${GOX} -ldflags "-X main.hoverctlVersion=${GIT_TAG_NAME}" -osarch="${OSARCH}" -output="${HCTL_BIN}"
  echo "${GIT_TAG_NAME} ${SUFFIX}" > ${VERSION_FILE}
  cp ${LICENSE} ${BIN_TARGET_DIR}/LICENSE.txt
  zip -j ${ZIP_FILE} ${BIN_TARGET_DIR}/*
done

rm -rf ${ZIP_TARGET_DIR}/bin
