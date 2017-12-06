#!/usr/bin/env bash

source "./tool.sh"

version=${1}
sysOS=`uname -s`
if [ $sysOS == "Darwin" ];then
	inMac
elif [ $sysOS == "Linux" ];then
	inLinux
else
	echo "Other OS: $sysOS"
    exit 1
fi

[[ -d dist-dicovery ]] && rm -rf dist-dicovery
install
if [[ ${version:=release} == "release" ]];then
  build
elif [[ ${version:=release} == "dev" ]]
then
    buildDev
else
    echo "no vesrions"
fi
