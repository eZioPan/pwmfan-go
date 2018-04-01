#!/bin/bash

# Author: eZioPan
# Page: https://github.com/eZioPan/pwmfan-go
# Date: 2018 Mar

BOLD='\e[1m'
RED='\e[31m'
GREEN='\e[32m'
CLR='\e[0m'

SERVICENAME=pwmfan
PREFIX=/usr/local
CFGDIR=/etc/${SERVICENAME}
CFGNAME=config.json
SYSTEMDUNITDIR=/etc/systemd/system

GOOS=linux
GOARCH=arm
export GOOS GOARCH

retcode=0

cd "$(dirname "$0")"

# Check defined program functionality
# Args:
#	$1 is programe name to check
#	$2 is argument to send to program
#	$3 and after is the return codes represent no error
#	If programe return code 0 represents error(which is strange), then this functions return code will be 1
function checkProgram {
	echo -e "check ${BOLD}${1}${CLR} exist"
	${1} ${2} >/dev/null 2>&1
	retcode=$?
	shift 2
	for arg in "$@"; do	
		if [ ${retcode} -eq ${arg} ]; then
			retcode=0
			return ${retcode}
		fi
	done
	echo -e "${BOLD}${RED}No ${1} found${CLR}"
	if [ ${retcode} -eq 0 ]; then
		retcode=1
	fi
	End ${retcode}
}

# stop service with systemctl
# Args:
#	$1 is the service name
function stopService {
	systemctl status ${1} >/dev/null 2>&1
	if [ $? -eq 0 ]; then
		echo "stop ${1} service"
		systemctl stop ${1}
		retcode=$?
	fi
	return ${retcode}
}

# disable service with systemctl
# Args:
#	$1 is the service name
function disableService {
	systemctl status ${1} >/dev/null 2>&1
	if [ $? -ne 4 ]; then
		echo "disable ${1} service"
		systemctl disable ${1}
		retcode=$?
	fi
	return ${retcode}
}

# create directory
# Args:
#	$1 is the directory
function createDir {
	if [ ! -d ${1} ]; then
		echo -e "Making config directory at ${GREEN}${1}${CLR}"
		mkdir -p ${1}
		retcode=$?
	fi
	return ${retcode}
}

# safely copy file to destination
# destination file exists with same content will skipped
# destination file exists with different content will result rename and copy
# Args:
#	$1 is the source file name
#	$2 is the destination file name
#	$3 is the name will be added after file name when destination file exist
function safeCopy {
	if [ ! -f ${2} ]; then
		echo -e "copying ${1} into ${GREEN}${2}${CLR}"
		cp ${1} ${2}
		retcode=$?
	else
		echo -e "found ${GREEN}${2}${CLR}"
		diff -q ${1} ${2} >/dev/null 2>&1
		retcode=$?
		if [ ${retcode} -ne 0 ]; then
			echo -e "copying new ${1} file as ${GREEN}${2}${3}${CLR}"
			cp ${1} ${2}${3}
			retcode=$?
			if [ ${retcode} -eq 0 ]; then
				echo -e "You should ${BOLD}${RED}manually compare${CLR} two files and combine them"
			fi
		else
			echo -e "${GREEN}${1}${CLR} and ${GREEN}${2}${CLR} are same, ${BOLD}skipped${CLR}"
		fi
	fi
	return ${retcode}
}

# remove file, skip none exist file without error
# Args:
#	$1 is the file name
function clearRemoveFile {
	if [ -f ${1} ]; then
	echo -e "Removing ${Green}${1}${CLR}"
		rm ${1}
		retcode=$?
	else
		echo -e "${GREEN}${1}${CLR} not found, skipped"
		retcode=0
	fi
	return ${retcode}
}

# remove directory, skip none exist directory without error
# Args:
#	$1 is the directory name
function clearRemoveDir {
	if [ -d ${1} ]; then
		echo -e "Removing ${Green}${1}${CLR}"
		rm -r ${1}
		retcode=$?
	else
		echo -e "${GREEN}${1}${CLR} not found, skipped"
		retcode=0
	fi
	return ${retcode}
}

function Start {
	echo -e "${BOLD}Install program start${CLR}"
}

function Help {
	echo -e "Use ${GREEN}./install.sh build${CLR} to build the program"
	echo -e "Use ${GREEN}./install.sh install${CLR} to install program into system"
	echo -e "Use ${GREEN}./install.sh uninstall${CLR} to remove program from system"
	echo -e "Use ${GREEN}./install.sh help${CLR} to show this help page"
}

function Build {
	checkProgram go version 0 && \
	echo -e "Building ${BOLD}${SERVICENAME}${CLR} from source"
	go build -o ${SERVICENAME} -v -ldflags "-s -w"
	retcode=$?
	if [ ${retcode} -eq 0 ]; then
		echo -e "${BOLD}${GREEN}Build success${CLR}"
	else
		echo -e "${BOLD}${RED}Build failed${CLR}"
	fi
	return ${retcode}
}

function Install {
	echo -e "${BOLD}Start install${CLR}" && \
	checkProgram systemctl --version 0 && \
	checkProgram diff --version 0 && \
	stopService ${SERVICENAME} && \
	createDir ${CFGDIR} && \
	safeCopy ${CFGNAME} ${CFGDIR}/${CFGNAME} ".new" && \
	echo -e "copying ${SERVICENAME} into ${GREEN}${PREFIX}/bin/${SERVICENAME}${CLR}" && \
	cp ${SERVICENAME} ${PREFIX}/bin/${SERVICENAME} && \
	safeCopy systemd/${SERVICENAME}.service ${SYSTEMDUNITDIR}/${SERVICENAME}.service ".new" && \
	echo -e "enable ${SERVICENAME} service in systemd" && \
	systemctl daemon-reload && \
	systemctl start ${SERVICENAME}.service && \
	systemctl enable ${SERVICENAME}.service
	retcode=$?
	if [ ${retcode} -eq 0 ]; then
		echo -e "${BOLD}${GREEN}Install success${CLR}"
	else
		echo -e "${BOLD}${RED}Install failed${CLR}"
	fi
	return ${retcode}
}

function Uninstall {
	echo -e "${BOLD}Start uninstall${CLR}" && \
	checkProgram systemctl --version 0 && \
	stopService ${SERVICENAME} && \
	disableService ${SERVICENAME} && \
	echo -e "removing ${SERVICENAME} files from system" && \
	clearRemoveFile ${SYSTEMDUNITDIR}/${SERVICENAME}.service && \
	clearRemoveFile ${SYSTEMDUNITDIR}/${SERVICENAME}.service.new && \
	clearRemoveFile ${PREFIX}/bin/${SERVICENAME} && \
	clearRemoveDir ${CFGDIR}
	retcode=$?
	if [ ${retcode} -eq 0 ]; then
		echo -e "${BOLD}${GREEN}Uninstall finished${CLR}"
	else
		echo -e "${BOLD}${RED}Uninstall failed${CLR}"
	fi
	return ${retcode}
}

function End {
	echo -e "${BOLD}Install program exit.${CLR}"
	exit ${1}
}

function ArgNumError {
	echo -e "Want ${GREEN}zero${CLR} or ${GREEN}one${CLR} argument"
}

Start

cd "$(dirname "$1")"

if [ $# -gt 1 ]; then
	ArgNumError
elif [ $# -eq 0 ]; then
	Help
elif [ $1 == "build" ]; then
	Build
elif [ $1 == "install" ]; then
	Install
elif [ $1 == "uninstall" ]; then
	Uninstall
else
	echo "Argument not recogonize"
fi

End ${retcode}
