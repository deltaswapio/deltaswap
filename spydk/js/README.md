# Wormhole Spy SDK

> Note: This is a pre-alpha release and in active development. Function names and signatures are subject to change.

Wormhole Spy service SDK for use with [@deltaswapio/deltaswap-sdk](https://www.npmjs.com/package/@deltaswapio/deltaswap-sdk)

## Usage

```js
import {
  createSpyRPCServiceClient,
  subscribeSignedVAA,
} from "@certusone/wormhole-spydk";
const client = createSpyRPCServiceClient(SPY_SERVICE_HOST);
const stream = await subscribeSignedVAA(client, {});
stream.on("data", ({ vaaBytes }) => {
  console.log(vaaBytes);
});
```

Also see [integration tests](https://github.com/deltaswapio/deltaswap/blob/main/spydk/js/src/__tests__/integration.ts)
