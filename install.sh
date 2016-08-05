#!/bin/bash

set -e

RELEASE=$(curl -silent "https://api.github.com/repos/SpectoLabs/hoverfly/releases/latest" | sed -n 's/.*"tag_name": "\(.*\)",/\1/p')

# Get OS and architecture

UNAME_STR=$(uname -sm)

if [[ ${UNAME_STR} == 'Darwin x86_64' ]]
then
   BINARY_OS_SUFFIX='_OSX_amd64'
elif [[ ${UNAME_STR} == 'Linux x86_64' ]]
then
  BINARY_OS_SUFFIX='_linux_amd64'
elif [[ ${UNAME_STR} == 'Linux i386' ]]
then
  BINARY_OS_SUFFIX='_linux_386'
else
  echo "OS could not be detected"
  exit
fi

# Download URLs

HF_BINARY=hoverfly_${RELEASE}${BINARY_OS_SUFFIX}
HF_DL_BASEURL=https://github.com/SpectoLabs/hoverfly/releases/download
HF_DL_URL=${HF_DL_BASEURL}/${RELEASE}/${HF_BINARY}

HFCTL_BINARY=hoverctl_${RELEASE}${BINARY_OS_SUFFIX}
HFCTL_DL_BASEURL=https://github.com/SpectoLabs/hoverfly/releases/download
HFCTL_DL_URL=${HFCTL_DL_BASEURL}/${RELEASE}/${HFCTL_BINARY}

# Destination directory

DESTINATION_DIR=/usr/local/bin

# Hoverctl config settings

HF_DIR=${HOME}/.hoverfly
HF_HOST=localhost
HF_ADMIN_PORT="8888"
HF_PROXY_PORT="8500"

echo
echo "Installing hoverctl CLI and Hoverfly ${RELEASE}"
echo

# Check if any Hoverfly files is already exist

if [[ -f ${DESTINATION_DIR}/hoverfly* || -f ${DESTINATION_DIR}/hoverctl* || -d ${HF_DIR} ]] ;
then
  read -p "Previous Hoverfly install detected in ${DESTINATION_DIR}. Overwrite? " -n 1 -r
  if [[ $REPLY =~ ^[Yy]$ ]]
  then
    echo
    rm -f ${DESTINATION_DIR}/hoverfly* ${DESTINATION_DIR}/hoverctl*
    rm -rf ${HF_DIR}
  else
    echo "Exiting"
    exit
  fi
fi

# Download Hoverfly & Hoverctl binaries

cd ${DESTINATION_DIR} || exit

echo
echo "Downloading hoverfly binary..."
curl -L -O# ${HF_DL_URL}
chmod +x ${HF_BINARY}
ln -s ${HF_BINARY} hoverfly

echo "Downloading hoverctl binary..."
curl -L -O# ${HFCTL_DL_URL}
chmod +x ${HFCTL_BINARY}
ln -s ${HFCTL_BINARY} hoverctl

# Create hoverctl config directory

mkdir ${HF_DIR}

# Get API key

echo
echo "*****************************************************"
echo "************** SpectoLab PRIVATE BETA ***************"
echo "*****************************************************"
echo "To use the SpectoLab with hoverctl, sign in at"
echo "https://lab.specto.io/settings to access your API Key."
echo
echo "If you don't have an API key, just press RETURN below."
echo
read -rsp "Enter your SpectoLab API key here : " API_KEY

# Create config file

cat >${HF_DIR}/config.yaml <<EOL
hoverfly.host: ${HF_HOST}
hoverfly.admin.port: ${HF_ADMIN_PORT}
hoverfly.proxy.port: ${HF_PROXY_PORT}
specto.lab.api.key: ${API_KEY}
EOL

echo
echo
echo "Hoverctl config file created: ${HF_DIR}/config.yaml. Edit this file to change hoverfly host, ports and SpectoLab API key."
echo
echo "Installation complete."
echo
echo "To get started:"
echo "1. Run 'hoverctl pull benjvi/hello-world:latest'"
echo "2. Set your http_proxy environment variable: 'export http_proxy=http://${HF_HOST}:${HF_PROXY_PORT}/'"
echo "3. Run 'hoverctl start'"
echo "3. Run 'hoverctl import benjvi/hello-world:latest'"
echo "4. Run 'curl -L http://lab.specto.io/static/hello-world.html', or open in the browser on linux"
echo
