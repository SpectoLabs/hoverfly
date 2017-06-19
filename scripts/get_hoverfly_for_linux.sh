#!/usr/bin/env bash

HOVERFLY_VERSION=$(curl -s https://api.github.com/repos/spectolabs/hoverfly/releases/latest | grep tag_name | sed -n 's/.*"tag_name": "\(.*\)",/\1/p')
if [[ $?  == 1 ]]; then
    error_exit "Failed to get latest release version"
fi

HOVERFLY_DOWNLOAD_URL=https://github.com/SpectoLabs/hoverfly/releases/download/${HOVERFLY_VERSION}


# Download distribution package for Linux
machine_type=$(uname -m)
if [[ ${machine_type} == "x86_64" ]]; then
    asset_name=hoverfly_bundle_linux_amd64.zip
else
    asset_name=hoverfly_bundle_linux_386.zip
fi

wget -O /tmp/hoverfly.zip ${HOVERFLY_DOWNLOAD_URL}/${asset_name}
if [[ $?  == 1 ]]; then
    error_exit "Failed to download hoverfly release package"
fi

# Unzip and copy to PATH
unzip -d /tmp/hoverfly /tmp/hoverfly.zip
sudo cp /tmp/hoverfly/hoverfly /usr/local/bin/
sudo cp /tmp/hoverfly/hoverctl /usr/local/bin/


# An error exit function
function error_exit
{
	echo "$1" 1>&2
	exit 1
}