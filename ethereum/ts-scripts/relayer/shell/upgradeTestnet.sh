npx tsx ./ts-scripts/relayer/config/checkNetworks.ts \
  && npx tsx ./ts-scripts/relayer/deliveryProvider/upgradeDeliveryProvider.ts \
  && npx tsx ./ts-scripts/relayer/deltaswapRelayer/upgradeDeltaswapRelayerSelfSign.ts
