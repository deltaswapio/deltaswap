npx tsx ./ts-scripts/relayer/config/checkNetworks.ts --set-last-run \
  && npx tsx ./ts-scripts/relayer/create2Factory/deployCreate2Factory.ts \
  && npx tsx ./ts-scripts/relayer/deliveryProvider/deployDeliveryProvider.ts \
  && npx tsx ./ts-scripts/relayer/deltaswapRelayer/deployDeltaswapRelayer.ts \
  && npx tsx ./ts-scripts/relayer/deliveryProvider/configureDeliveryProvider.ts \
  && npx tsx ./ts-scripts/relayer/deltaswapRelayer/registerChainsDeltaswapRelayerSelfSign.ts \
  && npx tsx ./ts-scripts/relayer/mockIntegration/deployMockIntegration.ts \
  && npx tsx ./ts-scripts/relayer/config/syncContractsJson.ts \
  && npx tsx ./ts-scripts/relayer/mockIntegration/messageTest.ts 

 # put this as 2nd script if not deployed aleady