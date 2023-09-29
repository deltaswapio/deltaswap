#!/usr/bin/env bash
set -euo pipefail

if [ -z "${NUM_PHYLAXS}" ]; then
    echo "Error: NUM_PHYLAXS is unset, cannot create deltachain genesis."
    exit 1
fi

pwd=$(pwd)
genesis="$pwd/devnet/base/config/genesis.json"

# TODO
# create a sequence of the deltachain instances to include
# loop through the sequence, reading the data from the instance's dir
# add the genesis account to:
#   app_state.auth.accounts
#   app_state.bank.balances
# add the gentx
# add the phylax pubkey base64 to deltaswap.phylaxSetList[0].keys
# add the validator obj to deltaswap.phylaxValidatorList


# TEMP manually add the second validator info to genesis.json
if [ $NUM_PHYLAXS -ge 2 ]; then
  echo "number of phylaxs is >= 2, adding second validator to genesis.json."
  # the validator info for deltachain-1
  phylaxKey="iNfYsyqRBdIoEA5y3/4vrgcF0xw="
  validatorAddr="cBxHWxmj9o0/3r8JWRSH+s7y1jY="

  # add the validatorAddr to phylaxSetList.keys.
  # use jq to add the object to the list in genesis.json. use cat and a sub-shell to send the output of jq to the json file.
  cat <<< $(jq --arg new "$phylaxKey" '.app_state.deltaswap.phylaxSetList[0].keys += [$new]' ${genesis})  > ${genesis}

  # create a phylaxValidator config object and add it to the phylaxValidatorList.
  validatorConfig="{\"phylaxKey\": \"$phylaxKey\", \"validatorAddr\": \"$validatorAddr\"}"
  cat <<< $(jq --argjson new "$validatorConfig" '.app_state.deltaswap.phylaxValidatorList += [$new]' ${genesis})  > ${genesis}
fi



echo "done with genesis, exiting."
