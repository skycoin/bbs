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

install
[[ -d dist-manager ]] && rm -rf dist-manager
if [[ ${version:=release} == "release" ]];then
  buildManager
elif [[ ${version:=release} == "dev" ]]
then
    buildManagerDev
else
    echo "no vesrions"
fi
