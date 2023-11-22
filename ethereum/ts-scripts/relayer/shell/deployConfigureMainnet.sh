
# This pipeline doesn't attempt to sign governance VAAs since those need to be signed by the phylaxs.
npx tsx ./ts-scripts/relayer/config/checkNetworks.ts --set-last-run \
  && npx tsx ./ts-scripts/relayer/deliveryProvider/deployDeliveryProvider.ts \
  && npx tsx ./ts-scripts/relayer/deltaswapRelayer/deployDeltaswapRelayer.ts \
  && npx tsx ./ts-scripts/relayer/deliveryProvider/configureDeliveryProvider.ts \
  && npx tsx ./ts-scripts/relayer/mockIntegration/deployMockIntegration.ts \
  && npx tsx ./ts-scripts/relayer/deltaswapRelayer/registerChainsDeltaswapRelayer.ts \
  && npx tsx ./ts-scripts/relayer/config/syncContractsJson.ts
