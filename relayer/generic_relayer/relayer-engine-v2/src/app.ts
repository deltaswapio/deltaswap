import Koa from "koa";
import Router from "koa-router";
import {
  Next,
  RelayerApp,
  StandardRelayerAppOpts,
  StandardRelayerContext,
  logging,
  wallets,
  spawnMissedVaaWorker,
  RedisStorage,
  providers,
  sourceTx,
} from "relayer-engine";
import { EVMChainId } from "@deltaswapio/deltaswap-sdk";
import { processGenericRelayerVaa } from "./processor";
import { Logger } from "winston";
import deepCopy from "clone";
import { loadAppConfig } from "./env";

export type GRContext = StandardRelayerContext & {
  deliveryProviders: Record<EVMChainId, string>;
  deltaswapRelayers: Record<EVMChainId, string>;
  opts: StandardRelayerAppOpts;
};

async function main() {
  const { env, opts, deliveryProviders, deltaswapRelayers } = await loadAppConfig();
  const logger = opts.logger!;
  logger.debug("Redis config: ", opts.redis);

  const app = new RelayerApp<GRContext>(env, opts);
  opts.redis = opts.redis?.redis;

  const {
    privateKeys,
    name,
    spyEndpoint,
    redis,
    redisCluster,
    redisClusterEndpoints,
    deltaswapRpcs,
  } = opts;
  app.spy(spyEndpoint);
  const store = new RedisStorage({
    redis,
    redisClusterEndpoints,
    redisCluster,
    attempts: opts.workflows?.retries ?? 3,
    namespace: name,
    queueName: `${name}-relays`,
  });

  app.useStorage(store);
  app.logger(logger);
  app.use(logging(logger));

  app.use(providers(opts.providers));
  if (opts.privateKeys && Object.keys(opts.privateKeys).length) {
    app.use(
      wallets(env, {
        logger,
        namespace: name,
        privateKeys: privateKeys!,
        metrics: { enabled: true, registry: store.registry},
      })
    );
  }
  if (opts.fetchSourceTxhash) {
    app.use(sourceTx());
  }

  spawnMissedVaaWorker(app, {
    namespace: name,
    registry: store.registry,
    logger,
    redis,
    redisCluster,
    redisClusterEndpoints,
    deltaswapRpcs,
    concurrency: opts.missedVaaOptions?.concurrency,
    checkInterval: opts.missedVaaOptions?.checkInterval,
    fetchVaaRetries: opts.missedVaaOptions?.fetchVaaRetries,
    vaasFetchConcurrency: opts.missedVaaOptions?.vaasFetchConcurrency,
    storagePrefix: opts.missedVaaOptions?.storagePrefix,
    startingSequenceConfig: opts.missedVaaOptions?.startingSequenceConfig,
    forceSeenKeysReindex: opts.missedVaaOptions?.forceSeenKeysReindex,
  });

  // Set up middleware
  app.use(async (ctx: GRContext, next: Next) => {
    ctx.deliveryProviders = deepCopy(deliveryProviders);
    ctx.deltaswapRelayers = deepCopy(deltaswapRelayers);
    ctx.opts = { ...opts };
    next();
  });

  // Set up routes
  app.multiple(deepCopy(deltaswapRelayers), processGenericRelayerVaa);

  app.listen();
  runApi(store, opts, logger);
}

function runApi(storage: RedisStorage, { port, redis }: any, logger: Logger) {
  const app = new Koa();
  const router = new Router();

  router.get("/metrics", async (ctx: Koa.Context) => {
    ctx.body = await storage.registry?.metrics();
  });

  app.use(router.routes());
  app.use(router.allowedMethods());

  if (redis?.host) {
    app.use(storage.storageKoaUI("/ui"));
  }

  port = Number(port) || 3000;
  app.listen(port, () => {
    logger.info(`Running on ${port}...`);
    logger.info(`For the UI, open http://localhost:${port}/ui`);
    logger.info(
      `For prometheus metrics, open http://localhost:${port}/metrics`
    );
    logger.info("Make sure Redis is running on port 6379 by default");
  });
}

main().catch((e) => {
  console.error("Encountered unrecoverable error:");
  console.error(e);
  process.exit(1);
});
