#!/usr/bin/env bash

HOVERFLY_VERSION=v0.12.1
HOVERFLY_DOWNLOAD_URL=https://github.com/SpectoLabs/hoverfly/releases/download


# Download distribution package for Linux
machine_type=$(uname -m)
if [[ ${machine_type} == "x86_64" ]]; then

    wget -O /tmp/hoverfly.zip ${HOVERFLY_DOWNLOAD_URL}/${HOVERFLY_VERSION}/hoverfly_bundle_linux_amd64.zip
else
    wget -O /tmp/hoverfly.zip ${HOVERFLY_DOWNLOAD_URL}/${HOVERFLY_VERSION}/hoverfly_bundle_linux_386.zip
fi


# Unzip and copy to PATH
unzip -d /tmp/hoverfly /tmp/hoverfly.zip
sudo cp /tmp/hoverfly/hoverfly /usr/local/bin/
sudo cp /tmp/hoverfly/hoverctl /usr/local/bin/