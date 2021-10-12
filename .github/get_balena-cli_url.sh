#! /bin/bash
set -euo pipefail

function urldecode() { : "${*//+/ }"; echo -e "${_//%/\\x}"; }

repo="balena-io/balena-cli"
jqQuery='.[0].assets[].browser_download_url | select(contains("linux") and contains("x64"))'
releaseJson=$(curl --fail -s https://api.github.com/repos/${repo}/releases)

latestVersionUrl=$(echo "$releaseJson" | jq -r "$jqQuery")

if [[ $(echo "$latestVersionUrl" | wc -l) != 1 ]]; then
	echo "ERROR: Found more than one: $(wc -l <"$latestVersionUrl")"
	echo "$latestVersionUrl"
	exit 1
fi
echo ${latestVersionUrl}

