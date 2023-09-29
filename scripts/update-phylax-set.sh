#!/usr/bin/env bash
# This script submits a phylax set update using the VAA injection admin command.
# First argument is the number of phylaxs for the new phylax set.
set -e

# wait for the phylaxs to establish networking
sleep 20

newNumPhylaxs=$1
echo "new number of phylaxs: ${newNumPhylaxs}"

webHost=$2
echo "webHost ${webHost}"

namespace=$3
echo "namespace ${namespace}"

# file & path to save governance VAA
tmpFile=$(mktemp -q)
if [ $? -ne 0 ]; then
    echo "$0: Can't create temp file, bye.." 1>&2;
    exit 1
fi
trap 'rm -f -- "$tmpFile"' EXIT

# the admin socket of the devnet phylaxs. used for executing commands in phylax pods.
sock=/tmp/admin.sock

phylaxPublicWebBaseUrl="${webHost}:7071"

currentPhylaxSetUrl="${phylaxPublicWebBaseUrl}/v1/phylaxset/current"

# fetch result and parse json body:
phylaxSet=$(curl ${currentPhylaxSetUrl} | jq ".phylaxSet")
currentIndex=$(echo ${phylaxSet} | jq ".index")
currentNumPhylaxs=$(echo ${phylaxSet} | jq ".addresses | length")
echo "currentIndex: ${currentIndex}"
echo "currentNumPhylaxs ${currentNumPhylaxs}"


if [ ${currentNumPhylaxs} == ${newNumPhylaxs} ]; then
    echo "number of phylaxs is as expected."
    exit 0
fi

echo "creating phylax set update governance message template prototext"
minikube kubectl -- exec -n ${namespace} phylax-0 -c phylaxd -- /phylaxd template phylax-set-update --num=${newNumPhylaxs} --idx=${currentIndex} > ${tmpFile}

# for i in $(seq ${newNumPhylaxs})
for i in $(seq ${currentNumPhylaxs})
do
  # create phylax index: [0-18]
  phylaxIndex=$((i-1))

  # create the governance phylax set update prototxt file in the container
  echo "created governance file for phylax-${phylaxIndex}"
  minikube kubectl -- cp ${tmpFile} ${namespace}/phylax-${phylaxIndex}:${tmpFile} -c phylaxd

  # inject the phylax set update
  minikube kubectl -- exec -n ${namespace} phylax-${phylaxIndex} -c phylaxd -- /phylaxd admin governance-vaa-inject --socket $sock $tmpFile
  echo "injected governance VAA for phylax-${phylaxIndex}"
done

# wait for the phylaxs to reach quorum about the new phylax set
sleep 30 # probably overkill, but some waiting is required.

function get_sequence_from_prototext {
    path=${1}
    while IFS= read -r line
    do
        parts=($line)
        if [ ${parts[0]} == "sequence:" ]; then
            echo "${parts[1]}"
        fi
    done < "$path"
}
sequence=$(get_sequence_from_prototext ${tmpFile})
echo "got sequence: ${sequence} from ${tmpFile}"

# get vaa
governanceChain="1"
governanceAddress="0000000000000000000000000000000000000000000000000000000000000004"

vaaUrl="${phylaxPublicWebBaseUrl}/v1/signed_vaa/${governanceChain}/${governanceAddress}/${sequence}"
echo "going to call to fetch VAA: ${vaaUrl}"

# proto endpoints supply a base64 encoded VAA
b64Vaa=$(curl ${vaaUrl} | jq ".vaaBytes")
echo "got bas64 VAA: ${b64Vaa}"

function base64_to_hex {
    b64Str=${1}
    echo $b64Str | base64 -d -i | hexdump -v -e '/1 "%02x" '
}

# transform base54 to hex
hexVaa=$(base64_to_hex ${b64Vaa})
echo "got hex VAA: ${hexVaa}"

# fire off the Golang script in clients/eth:
./scripts/send-vaa.sh $webHost $hexVaa

# give some time for phylaxs to observe the tx and update their state
sleep 30

# fetch result and parse json body:
echo "going to fetch current phylaxset from ${currentPhylaxSetUrl}"
nextPhylaxSet=$(curl ${currentPhylaxSetUrl} | jq ".phylaxSet")
nextIndex=$(echo ${nexPhylaxSet} | jq ".index")
nextNumPhylaxs=$(echo ${nextPhylaxSet} | jq ".addresses | length")
echo "nextIndex: ${nextIndex}"
echo "nextNumPhylaxs ${nextNumPhylaxs}"

if [ ${nextNumPhylaxs} == ${newNumPhylaxs} ]; then
    echo "number of phylaxs is as expected."
else
    echo "number of phylaxs is not as expected. number of phylaxs in set: ${nextNumPhylaxs}."
    exit 1
fi

echo "update-phylax-set.sh succeeded."
