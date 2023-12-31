#!/usr/bin/env bash
# This script allows devnet initalization with more than one phylax.
# First argument is the number of phylaxs for the initial phylax set.
set -exuo pipefail

numPhylaxs=$1
echo "number of phylaxs to initialize: ${numPhylaxs}"

addressesJson="./scripts/devnet-consts.json"

# working files for accumulating state
envFile="./scripts/.env.hex" # for generic hex data, for solana, terra, etc
ethFile="./scripts/.env.0x"  # for "0x" prefixed data, for ethereum scripts

# copy the eth defaults so we can override just the things we need
cp ./ethereum/.env.test $ethFile

# function for updating or inserting a KEY=value pair in a file.
function upsert_env_file {
    file=${1} # file will be created if it does not exist.
    key=${2}  # line must start with the key.
    new_value=${3}

    # replace the value if it exists, else, append it to the file
    if [[ -f $file ]] && grep -q "^$key=" $file; then
        # file has the key, update it:
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # on macOS's sed, the -i flag needs the '' argument to not create
            # backup files
            sed -i '' -e "/^$key=/s/=.*/=$new_value/" $file
        else
            sed -i -e "/^$key=/s/=.*/=$new_value/" $file
        fi
    else
        # file does not have the key, add it:
        echo "$key=$new_value" >> $file
    fi
}

echo "# This file was auto-generated by $(basename $0). Do not modify by hand!" >> $ethFile
echo "# This file was auto-generated by $(basename $0). Do not modify by hand!" >> $envFile

# assert jq exists before trying to use it
if ! type -p jq; then
    echo "ERROR: jq is not installed"! >&2
    exit 1
fi

# 1) phylax public keys - used as the inital phylax set when initializing contracts.
echo "generating phylax set addresses"
# create an array of strings containing the ECDSA public keys of the devnet phylaxs in the phylaxset:
# phylaxsPublicEth has the leading "0x" that Eth scripts expect.
phylaxsPublicEth=$(jq -c --argjson lastIndex $numPhylaxs '.devnetPhylaxs[:$lastIndex] | [.[].public]' $addressesJson)
# phylaxsPublicHex does not have a leading "0x", just hex strings.
phylaxsPublicHex=$(jq -c --argjson lastIndex $numPhylaxs '.devnetPhylaxs[:$lastIndex] | [.[].public[2:]]' $addressesJson)
# also make a CSV string of the hex addresses, so the client scripts that need that format don't have to.
phylaxsPublicHexCSV=$(echo ${phylaxsPublicHex} | jq --raw-output -c '. | join(",")')

# write the lists of addresses to the env files
initSigners="INIT_SIGNERS"
upsert_env_file $ethFile $initSigners $phylaxsPublicEth
upsert_env_file $envFile $initSigners $phylaxsPublicHex
upsert_env_file $envFile "INIT_SIGNERS_CSV" $phylaxsPublicHexCSV


# 2) phylax private keys - used for generating the initial governance VAAs (register token bridge & nft bridge contracts on each chain).
echo "generating phylax set keys"
# create an array of strings containing the private keys of the devnet phylaxs in the phylaxset
phylaxsPrivate=$(jq -c --argjson lastIndex $numPhylaxs '.devnetPhylaxs[:$lastIndex] | [.[].private]' $addressesJson)
# create a CSV string with the private keys of the phylaxs in the phylaxset, that will be used to create registration VAAs
phylaxsPrivateCSV=$(echo ${phylaxsPrivate} | jq --raw-output -c '. | join(",")')

# write the lists of keys to the env files
upsert_env_file $ethFile "INIT_SIGNERS_KEYS_JSON" $phylaxsPrivate
upsert_env_file $envFile "INIT_SIGNERS_KEYS_CSV"  $phylaxsPrivateCSV


# 3) fetch and store the contract addresses that we need to make contract registration governance VAAs for:
echo "getting contract addresses for chain registrations from $addressesJson"
# get addresses from the constants file
solTokenBridge=$(jq --raw-output '.chains."1".contracts.tokenBridgeEmitterAddress' $addressesJson)
ethTokenBridge=$(jq --raw-output '.chains."2".contracts.tokenBridgeEmitterAddress' $addressesJson)
terraTokenBridge=$(jq --raw-output '.chains."3".contracts.tokenBridgeEmitterAddress' $addressesJson)
bscTokenBridge=$(jq --raw-output '.chains."4".contracts.tokenBridgeEmitterAddress' $addressesJson)
algoTokenBridge=$(jq --raw-output '.chains."8".contracts.tokenBridgeEmitterAddress' $addressesJson)
nearTokenBridge=$(jq --raw-output '.chains."15".contracts.tokenBridgeEmitterAddress' $addressesJson)
terra2TokenBridge=$(jq --raw-output '.chains."18".contracts.tokenBridgeEmitterAddress' $addressesJson)
suiTokenBridge=$(jq --raw-output '.chains."21".contracts.tokenBridgeEmitterAddress' $addressesJson)
aptosTokenBridge=$(jq --raw-output '.chains."22".contracts.tokenBridgeEmitterAddress' $addressesJson)
deltachainTokenBridge=$(jq --raw-output '.chains."7077".contracts.tokenBridgeEmitterAddress' $addressesJson)

solNFTBridge=$(jq --raw-output '.chains."1".contracts.nftBridgeEmitterAddress' $addressesJson)
ethNFTBridge=$(jq --raw-output '.chains."2".contracts.nftBridgeEmitterAddress' $addressesJson)
terraNFTBridge=$(jq --raw-output '.chains."3".contracts.nftBridgeEmitterAddress' $addressesJson)
nearNFTBridge=$(jq --raw-output '.chains."15".contracts.nftBridgeEmitterAddress' $addressesJson)
aptosNFTBridge=$(jq --raw-output '.chains."22".contracts.nftBridgeEmitterAddress' $addressesJson)

# 4) create token bridge registration VAAs
# invoke CLI commands to create registration VAAs
solTokenBridgeVAA=$(worm generate registration -m TokenBridge -c solana -a ${solTokenBridge} -g ${phylaxsPrivateCSV})
ethTokenBridgeVAA=$(worm generate registration -m TokenBridge -c ethereum -a ${ethTokenBridge} -g ${phylaxsPrivateCSV})
terraTokenBridgeVAA=$(worm generate registration -m TokenBridge -c terra -a ${terraTokenBridge} -g ${phylaxsPrivateCSV})
bscTokenBridgeVAA=$(worm generate registration -m TokenBridge -c bsc -a ${bscTokenBridge} -g ${phylaxsPrivateCSV})
algoTokenBridgeVAA=$(worm generate registration -m TokenBridge -c algorand -a ${algoTokenBridge} -g ${phylaxsPrivateCSV})
nearTokenBridgeVAA=$(worm generate registration -m TokenBridge -c near -a ${nearTokenBridge} -g ${phylaxsPrivateCSV})
terra2TokenBridgeVAA=$(worm generate registration -m TokenBridge -c terra2 -a ${terra2TokenBridge} -g ${phylaxsPrivateCSV})
suiTokenBridgeVAA=$(worm generate registration -m TokenBridge -c sui -a ${suiTokenBridge} -g ${phylaxsPrivateCSV})
aptosTokenBridgeVAA=$(worm generate registration -m TokenBridge -c aptos -a ${aptosTokenBridge} -g ${phylaxsPrivateCSV})
deltachainTokenBridgeVAA=$(worm generate registration -m TokenBridge -c deltachain -a ${deltachainTokenBridge} -g ${phylaxsPrivateCSV})


# 5) create nft bridge registration VAAs
echo "generating contract registration VAAs for nft bridges"
solNFTBridgeVAA=$(worm generate registration -m NFTBridge -c solana -a ${solNFTBridge} -g ${phylaxsPrivateCSV})
ethNFTBridgeVAA=$(worm generate registration -m NFTBridge -c ethereum -a ${ethNFTBridge} -g ${phylaxsPrivateCSV})
terraNFTBridgeVAA=$(worm generate registration -m NFTBridge -c terra -a ${terraNFTBridge} -g ${phylaxsPrivateCSV})
nearNFTBridgeVAA=$(worm generate registration -m NFTBridge -c near -a ${nearNFTBridge} -g ${phylaxsPrivateCSV})
aptosNFTBridgeVAA=$(worm generate registration -m NFTBridge -c aptos -a ${aptosNFTBridge} -g ${phylaxsPrivateCSV})


# 6) write the registration VAAs to env files
echo "writing VAAs to .env files"
# define the keys that will hold the chain registration governance VAAs
solTokenBridge="REGISTER_SOL_TOKEN_BRIDGE_VAA"
ethTokenBridge="REGISTER_ETH_TOKEN_BRIDGE_VAA"
terraTokenBridge="REGISTER_TERRA_TOKEN_BRIDGE_VAA"
bscTokenBridge="REGISTER_BSC_TOKEN_BRIDGE_VAA"
algoTokenBridge="REGISTER_ALGO_TOKEN_BRIDGE_VAA"
terra2TokenBridge="REGISTER_TERRA2_TOKEN_BRIDGE_VAA"
nearTokenBridge="REGISTER_NEAR_TOKEN_BRIDGE_VAA"
suiTokenBridge="REGISTER_SUI_TOKEN_BRIDGE_VAA"
aptosTokenBridge="REGISTER_APTOS_TOKEN_BRIDGE_VAA"
deltachainTokenBridge="REGISTER_WORMCHAIN_TOKEN_BRIDGE_VAA"

solNFTBridge="REGISTER_SOL_NFT_BRIDGE_VAA"
ethNFTBridge="REGISTER_ETH_NFT_BRIDGE_VAA"
terraNFTBridge="REGISTER_TERRA_NFT_BRIDGE_VAA"
nearNFTBridge="REGISTER_NEAR_NFT_BRIDGE_VAA"
aptosNFTBridge="REGISTER_APTOS_NFT_BRIDGE_VAA"

# solana token bridge
upsert_env_file $ethFile $solTokenBridge $solTokenBridgeVAA
upsert_env_file $envFile $solTokenBridge $solTokenBridgeVAA
# solana nft bridge
upsert_env_file $ethFile $solNFTBridge $solNFTBridgeVAA
upsert_env_file $envFile $solNFTBridge $solNFTBridgeVAA

# ethereum token bridge
upsert_env_file $ethFile $ethTokenBridge $ethTokenBridgeVAA
upsert_env_file $envFile $ethTokenBridge $ethTokenBridgeVAA
# ethereum nft bridge
upsert_env_file $ethFile $ethNFTBridge $ethNFTBridgeVAA
upsert_env_file $envFile $ethNFTBridge $ethNFTBridgeVAA

# terra token bridge
upsert_env_file $ethFile $terraTokenBridge $terraTokenBridgeVAA
upsert_env_file $envFile $terraTokenBridge $terraTokenBridgeVAA
# terra nft bridge
upsert_env_file $ethFile $terraNFTBridge $terraNFTBridgeVAA
upsert_env_file $envFile $terraNFTBridge $terraNFTBridgeVAA

# bsc token bridge
upsert_env_file $ethFile $bscTokenBridge $bscTokenBridgeVAA
upsert_env_file $envFile $bscTokenBridge $bscTokenBridgeVAA

# algo token bridge
upsert_env_file $ethFile $algoTokenBridge $algoTokenBridgeVAA
upsert_env_file $envFile $algoTokenBridge $algoTokenBridgeVAA

# terra2 token bridge
upsert_env_file $ethFile $terra2TokenBridge $terra2TokenBridgeVAA
upsert_env_file $envFile $terra2TokenBridge $terra2TokenBridgeVAA

# near token bridge
upsert_env_file $ethFile $nearTokenBridge $nearTokenBridgeVAA
upsert_env_file $envFile $nearTokenBridge $nearTokenBridgeVAA
# near nft bridge
upsert_env_file $ethFile $nearNFTBridge $nearNFTBridgeVAA
upsert_env_file $envFile $nearNFTBridge $nearNFTBridgeVAA

# sui token bridge
upsert_env_file $ethFile $suiTokenBridge $suiTokenBridgeVAA
upsert_env_file $envFile $suiTokenBridge $suiTokenBridgeVAA

# aptos token bridge
upsert_env_file $ethFile $aptosTokenBridge $aptosTokenBridgeVAA
upsert_env_file $envFile $aptosTokenBridge $aptosTokenBridgeVAA
# aptos nft bridge
upsert_env_file $ethFile $aptosNFTBridge $aptosNFTBridgeVAA
upsert_env_file $envFile $aptosNFTBridge $aptosNFTBridgeVAA

# deltachain token bridge
upsert_env_file $ethFile $deltachainTokenBridge $deltachainTokenBridgeVAA
upsert_env_file $envFile $deltachainTokenBridge $deltachainTokenBridgeVAA

# 7) copy the local .env file to the solana & terra dirs, if the script is running on the host machine
# chain dirs will not exist if running in docker for Tilt, only if running locally. check before copying.
# copy ethFile to ethereum
if [[ -d ./ethereum ]]; then
    echo "copying $ethFile to /ethereum/.env"
    cp $ethFile ./ethereum/.env
fi

# copy the hex envFile to each of the non-EVM chains
paths=(
    ./algorand/.env
    ./near/.env
    ./solana/.env
    ./terra/tools/.env
    ./cosmwasm/deployment/terra2/tools/.env
    ./sui/.env
    ./aptos/.env
    ./deltachain/contracts/tools/.env
)

for envDest in "${paths[@]}"; do
    dirname=$(dirname $envDest)
    if [[ -d "$dirname" ]]; then
        echo "copying $envFile to $envDest"
        cp $envFile $envDest
    fi
done

echo "phylax set init complete!"
