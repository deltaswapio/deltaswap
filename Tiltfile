# This Tiltfile contains the deployment and build config for the Wormhole devnet.
#
#  We use Buildkit cache mounts and careful layering to avoid unnecessary rebuilds - almost
#  all source code changes result in small, incremental rebuilds. Dockerfiles are written such
#  that, for example, changing the contract source code won't cause Solana itself to be rebuilt.
#

load("ext://namespace", "namespace_create", "namespace_inject")
load("ext://secret", "secret_yaml_generic")

# set the replica value of a StatefulSet
def set_replicas_in_statefulset(config_yaml, statefulset_name,  num_replicas):
    for obj in config_yaml:
        if obj["kind"] == "StatefulSet" and obj["metadata"]["name"] == statefulset_name:
            obj["spec"]["replicas"] = num_replicas
    return config_yaml

# set the env value of all containers in all jobs
def set_env_in_jobs(config_yaml, name, value):
    for obj in config_yaml:
        if obj["kind"] == "Job":
            for container in obj["spec"]["template"]["spec"]["containers"]:
                if not "env" in container:
                    container["env"] = []
                container["env"].append({"name": name, "value": value})
    return config_yaml

allow_k8s_contexts("ci")

# Disable telemetry by default
analytics_settings(False)

# Moar updates (default is 3)
update_settings(max_parallel_updates = 10)

# Runtime configuration
config.define_bool("ci", False, "We are running in CI")
config.define_bool("manual", False, "Set TRIGGER_MODE_MANUAL by default")

config.define_string("num", False, "Number of phylax nodes to run")

# You do not usually need to set this argument - this argument is for debugging only. If you do use a different
# namespace, note that the "wormhole" namespace is hardcoded in tests and don't forget specifying the argument
# when running "tilt down".
#
config.define_string("namespace", False, "Kubernetes namespace to use")

# When running Tilt on a server, this can be used to set the public hostname Tilt runs on
# for service links in the UI to work.
config.define_string("webHost", False, "Public hostname for port forwards")

# When running Tilt on a server, this can be used to set the public hostname Tilt runs on
# for service links in the UI to work.
config.define_string("phylaxd_loglevel", False, "Log level for phylaxd (debug, info, warn, error, dpanic, panic, fatal)")

# Components
config.define_bool("near", False, "Enable Near component")
config.define_bool("sui", False, "Enable Sui component")
config.define_bool("btc", False, "Enable BTC component")
config.define_bool("aptos", False, "Enable Aptos component")
config.define_bool("algorand", False, "Enable Algorand component")
config.define_bool("evm2", False, "Enable second Eth component")
config.define_bool("solana", False, "Enable Solana component")
config.define_bool("pythnet", False, "Enable PythNet component")
config.define_bool("terra_classic", False, "Enable Terra Classic component")
config.define_bool("terra2", False, "Enable Terra 2 component")
config.define_bool("ci_tests", False, "Enable tests runner component")
config.define_bool("phylaxd_debug", False, "Enable dlv endpoint for phylaxd")
config.define_bool("node_metrics", False, "Enable Prometheus & Grafana for Phylax metrics")
config.define_bool("phylaxd_governor", False, "Enable chain governor in phylaxd")
config.define_bool("deltachain", False, "Enable a deltachain node")
config.define_bool("ibc_relayer", False, "Enable IBC relayer between cosmos chains")
config.define_bool("redis", False, "Enable a redis instance")
config.define_bool("generic_relayer", False, "Enable the generic relayer off-chain component")


cfg = config.parse()
num_phylaxs = int(cfg.get("num", "1"))
namespace = cfg.get("namespace", "wormhole")
webHost = cfg.get("webHost", "localhost")
ci = cfg.get("ci", False)
algorand = cfg.get("algorand", ci)
near = cfg.get("near", ci)
aptos = cfg.get("aptos", ci)
sui = cfg.get("sui", ci)
evm2 = cfg.get("evm2", ci)
solana = cfg.get("solana", ci)
pythnet = cfg.get("pythnet", False)
terra_classic = cfg.get("terra_classic", ci)
terra2 = cfg.get("terra2", ci)
deltachain = cfg.get("deltachain", ci)
ci_tests = cfg.get("ci_tests", ci)
phylaxd_debug = cfg.get("phylaxd_debug", False)
node_metrics = cfg.get("node_metrics", False)
phylaxd_governor = cfg.get("phylaxd_governor", False)
ibc_relayer = cfg.get("ibc_relayer", ci)
btc = cfg.get("btc", False)
redis = cfg.get('redis', ci)
generic_relayer = cfg.get("generic_relayer", ci)

if ci:
    phylaxd_loglevel = cfg.get("phylaxd_loglevel", "warn")
else:
    phylaxd_loglevel = cfg.get("phylaxd_loglevel", "info")


if cfg.get("manual", False):
    trigger_mode = TRIGGER_MODE_MANUAL
else:
    trigger_mode = TRIGGER_MODE_AUTO

# namespace

if not ci:
    namespace_create(namespace)

def k8s_yaml_with_ns(objects):
    return k8s_yaml(namespace_inject(objects, namespace))

docker_build(
    ref = "cli-gen",
    context = ".",
    dockerfile = "Dockerfile.cli",
)

docker_build(
    ref = "const-gen",
    context = ".",
    dockerfile = "Dockerfile.const",
    build_args={"num_phylaxs": '%s' % (num_phylaxs)},
)

# node

docker_build(
    ref = "phylaxd-image",
    context = ".",
    dockerfile = "node/Dockerfile",
    target = "build",
    ignore=["./sdk/js", "./relayer"]
)

def command_with_dlv(argv):
    return [
        "/dlv",
        "--listen=0.0.0.0:2345",
        "--accept-multiclient",
        "--headless=true",
        "--api-version=2",
        "--continue=true",
        "exec",
        argv[0],
        "--",
    ] + argv[1:]

def build_node_yaml():
    node_yaml = read_yaml_stream("devnet/node.yaml")

    node_yaml_with_replicas = set_replicas_in_statefulset(node_yaml, "phylax", num_phylaxs)

    for obj in node_yaml_with_replicas:
        if obj["kind"] == "StatefulSet" and obj["metadata"]["name"] == "phylax":
            container = obj["spec"]["template"]["spec"]["containers"][0]
            if container["name"] != "phylaxd":
                fail("container 0 is not phylaxd")

            container["command"] += ["--logLevel="+phylaxd_loglevel]

            if phylaxd_debug:
                container["command"] = command_with_dlv(container["command"])
                print(container["command"])

            if aptos:
                container["command"] += [
                    "--aptosRPC",
                    "http://aptos:8080",
                    "--aptosAccount",
                    "de0036a9600559e295d5f6802ef6f3f802f510366e0c23912b0655d972166017",
                    "--aptosHandle",
                    "0xde0036a9600559e295d5f6802ef6f3f802f510366e0c23912b0655d972166017::state::WormholeMessageHandle",
                ]

            if sui:
                container["command"] += [
                    "--suiRPC",
                    "http://sui:9000",
                    "--suiMoveEventType",
                    "0x7f6cebb8a489654d7a759483bd653c4be3e5ccfef17a8b5fd3ba98bd072fabc3::publish_message::WormholeMessage",
                    "--suiWS",
                    "sui:9000",
                ]

            if evm2:
                container["command"] += [
                    "--bscRPC",
                    "ws://eth-devnet2:8545",
                ]
            else:
                container["command"] += [
                    "--bscRPC",
                    "ws://eth-devnet:8545",
                ]

            if solana:
                container["command"] += [
                    "--solanaRPC",
                    "http://solana-devnet:8899",
                ]

            if pythnet:
                container["command"] += [
                    "--pythnetRPC",
#                    "http://solana-devnet:8899",
                     "http://pythnet.rpcpool.com",
                    "--pythnetWS",
#                   "ws://solana-devnet:8900",
                    "wss://pythnet.rpcpool.com",
                    "--pythnetContract",
                    "H3fxXJ86ADW2PNuDDmZJg6mzTtPxkYCpNuQUTgmJ7AjU",
                ]

            if terra_classic:
                container["command"] += [
                    "--terraWS",
                    "ws://terra-terrad:26657/websocket",
                    "--terraLCD",
                    "http://terra-terrad:1317",
                    "--terraContract",
                    "terra18vd8fpwxzck93qlwghaj6arh4p7c5n896xzem5",
                ]

            if terra2:
                container["command"] += [
                    "--terra2WS",
                    "ws://terra2-terrad:26657/websocket",
                    "--terra2LCD",
                    "http://terra2-terrad:1317",
                    "--terra2Contract",
                    "terra14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9ssrc8au",
                ]

            if algorand:
                container["command"] += [
                    "--algorandAppID",
                    "1004",
                    "--algorandIndexerRPC",
                    "http://algorand:8980",
                    "--algorandIndexerToken",
                    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
                    "--algorandAlgodRPC",
                    "http://algorand:4001",
                    "--algorandAlgodToken",
                    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
                ]

            if phylaxd_governor:
                container["command"] += [
                    "--chainGovernorEnabled"
                ]

            if near:
                container["command"] += [
                    "--nearRPC",
                    "http://near:3030",
                    "--nearContract",
                    "wormhole.test.near"
                ]

            if deltachain:
                container["command"] += [
                    "--deltachainURL",
                    "deltachain:9090",

                    "--accountantContract",
                    "wormhole14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9srrg465",
                    "--accountantKeyPath",
                    "/tmp/mounted-keys/deltachain/accountantKey",
                    "--accountantKeyPassPhrase",
                    "test0000",
                    "--accountantWS",
                    "http://deltachain:26657",
                    "--accountantCheckEnabled",
                    "true",

                    "--ibcContract",
                    "wormhole1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrq0kdhcj",
                    "--ibcWS",
                    "ws://deltachain:26657/websocket",
                    "--ibcLCD",
                    "http://deltachain:1317",

                    "--gatewayRelayerContract",
                    "wormhole17p9rzwnnfxcjp32un9ug7yhhzgtkhvl9jfksztgw5uh69wac2pgshdnj3k",
                    "--gatewayRelayerKeyPath",
                    "/tmp/mounted-keys/deltachain/gwrelayerKey",
                    "--gatewayRelayerKeyPassPhrase",
                    "test0000",

                    "--gatewayContract",
                    "wormhole17p9rzwnnfxcjp32un9ug7yhhzgtkhvl9jfksztgw5uh69wac2pgshdnj3k",
                    "--gatewayWS",
                    "ws://deltachain:26657/websocket",
                    "--gatewayLCD",
                    "http://deltachain:1317"
                ]

    return encode_yaml_stream(node_yaml_with_replicas)

k8s_yaml_with_ns(build_node_yaml())

phylax_resource_deps = ["eth-devnet"]
if evm2:
    phylax_resource_deps = phylax_resource_deps + ["eth-devnet2"]
if solana or pythnet:
    phylax_resource_deps = phylax_resource_deps + ["solana-devnet"]
if near:
    phylax_resource_deps = phylax_resource_deps + ["near"]
if terra_classic:
    phylax_resource_deps = phylax_resource_deps + ["terra-terrad"]
if terra2:
    phylax_resource_deps = phylax_resource_deps + ["terra2-terrad"]
if algorand:
    phylax_resource_deps = phylax_resource_deps + ["algorand"]
if aptos:
    phylax_resource_deps = phylax_resource_deps + ["aptos"]
if deltachain:
    phylax_resource_deps = phylax_resource_deps + ["deltachain", "deltachain-deploy"]
if sui:
    phylax_resource_deps = phylax_resource_deps + ["sui"]

k8s_resource(
    "phylax",
    resource_deps = phylax_resource_deps,
    port_forwards = [
        port_forward(6060, name = "Debug/Status Server [:6060]", host = webHost),
        port_forward(7070, name = "Public gRPC [:7070]", host = webHost),
        port_forward(7071, name = "Public REST [:7071]", host = webHost),
        port_forward(2345, name = "Debugger [:2345]", host = webHost),
    ],
    labels = ["phylax"],
    trigger_mode = trigger_mode,
)

# phylax set update - triggered by "tilt args" changes
if num_phylaxs >= 2 and ci == False:
    local_resource(
        name = "phylax-set-update",
        resource_deps = phylax_resource_deps + ["phylax"],
        deps = ["scripts/send-vaa.sh", "clients/eth"],
        cmd = './scripts/update-phylax-set.sh %s %s %s' % (num_phylaxs, webHost, namespace),
        labels = ["phylax"],
        trigger_mode = trigger_mode,
    )


# grafana + prometheus for node metrics
if node_metrics:

    dashboard = read_json("dashboards/Wormhole.json")

    dashboard_yaml =  {
        "apiVersion": "v1",
        "kind": "ConfigMap",
        "metadata": {
            "name": "grafana-dashboards-json"
        },
        "data": {
            "wormhole.json": encode_json(dashboard)
        }
    }
    k8s_yaml_with_ns(encode_yaml(dashboard_yaml))

    k8s_yaml_with_ns("devnet/node-metrics.yaml")

    k8s_resource(
        "prometheus-server",
        resource_deps = ["phylax"],
        port_forwards = [
            port_forward(9099, name = "Prometheus [:9099]", host = webHost),
        ],
        labels = ["phylax"],
        trigger_mode = trigger_mode,
    )

    k8s_resource(
        "grafana",
        resource_deps = ["prometheus-server"],
        port_forwards = [
            port_forward(3033, name = "Grafana UI [:3033]", host = webHost),
        ],
        labels = ["phylax"],
        trigger_mode = trigger_mode,
    )


# spy
k8s_yaml_with_ns("devnet/spy.yaml")

k8s_resource(
    "spy",
    resource_deps = ["phylax"],
    port_forwards = [
        port_forward(6061, container_port = 6060, name = "Debug/Status Server [:6061]", host = webHost),
        port_forward(7072, name = "Spy gRPC [:7072]", host = webHost),
    ],
    labels = ["phylax"],
    trigger_mode = trigger_mode,
)

if solana or pythnet:
    # solana client cli (used for devnet setup)

    docker_build(
        ref = "bridge-client",
        context = ".",
        only = ["./proto", "./solana", "./clients"],
        dockerfile = "solana/Dockerfile.client",
        # Ignore target folders from local (non-container) development.
        ignore = ["./solana/*/target"],
    )

    # solana smart contract

    docker_build(
        ref = "solana-contract",
        context = "solana",
        dockerfile = "solana/Dockerfile",
        target = "builder",
        build_args = {"BRIDGE_ADDRESS": "Bridge1p5gheXUvJ6jGWGeCsgPKgnE3YgdGKRVCMY9o"}
    )

    # solana local devnet

    k8s_yaml_with_ns("devnet/solana-devnet.yaml")

    k8s_resource(
        "solana-devnet",
        port_forwards = [
            port_forward(8899, name = "Solana RPC [:8899]", host = webHost),
            port_forward(8900, name = "Solana WS [:8900]", host = webHost),
        ],
        labels = ["solana"],
        trigger_mode = trigger_mode,
    )

# eth devnet

docker_build(
    ref = "eth-node",
    context = "./ethereum",
    dockerfile = "./ethereum/Dockerfile",

    # ignore local node_modules (in case they're present)
    ignore = ["./node_modules"],
    build_args = {"num_phylaxs": str(num_phylaxs), "dev": str(not ci)},
  
    # sync external scripts for incremental development
    # (everything else needs to be restarted from scratch for determinism)
    #
    # This relies on --update-mode=exec to work properly with a non-root user.
    # https://github.com/tilt-dev/tilt/issues/3708
    live_update = [
        sync("./ethereum/src", "/home/node/app/src"),
    ],
)

if redis or generic_relayer:
    docker_build(
        ref = "redis",
        context = ".",
        only = ["./third_party"],
        dockerfile = "third_party/redis/Dockerfile",
    )

if redis:
    k8s_resource(
        "redis",
        port_forwards = [
            port_forward(6379, name = "Redis Default [:6379]", host = webHost),
        ],
        labels = ["redis"],
        trigger_mode = trigger_mode,
    )

    k8s_yaml_with_ns("devnet/redis.yaml")

if generic_relayer:
    k8s_resource(
        "redis-relayer",
        port_forwards = [
            port_forward(6378, name = "Generic Relayer Redis [:6378]", host = webHost),
        ],
        labels = ["redis-relayer"],
        trigger_mode = trigger_mode,
    )

    k8s_yaml_with_ns("devnet/redis-relayer.yaml")



if generic_relayer:
    k8s_resource(
        "relayer-engine",
        resource_deps = ["phylax", "redis-relayer", "spy"],
        port_forwards = [
            port_forward(3003, container_port=3000, name = "Bullmq UI [:3003]", host = webHost),
        ],
        labels = ["relayer-engine"],
        trigger_mode = trigger_mode,
    )
    docker_build(
        ref = "relayer-engine",
        context = ".",
        only = ["./relayer/generic_relayer", "./ethereum/ts-scripts/relayer/config"],
        dockerfile = "relayer/generic_relayer/relayer-engine-v2/Dockerfile",
        build_args = {"dev": str(not ci)}
    )
    k8s_yaml_with_ns("devnet/relayer-engine.yaml")

k8s_yaml_with_ns("devnet/eth-devnet.yaml")

k8s_resource(
    "eth-devnet",
    port_forwards = [
        port_forward(8545, name = "Ganache RPC [:8545]", host = webHost),
    ],
    labels = ["evm"],
    trigger_mode = trigger_mode,
)

if evm2:
    k8s_yaml_with_ns("devnet/eth-devnet2.yaml")

    k8s_resource(
        "eth-devnet2",
        port_forwards = [
            port_forward(8546, name = "Ganache RPC [:8546]", host = webHost),
        ],
        labels = ["evm"],
        trigger_mode = trigger_mode,
    )


if ci_tests:
    docker_build(
        ref = "sdk-test-image",
        context = ".",
        dockerfile = "testing/Dockerfile.sdk.test",
        only = [],
        live_update = [
            sync("./sdk/js/src", "/app/sdk/js/src"),
            sync("./testing", "/app/testing"),
        ],
    )
    docker_build(
        ref = "spydk-test-image",
        context = ".",
        dockerfile = "testing/Dockerfile.spydk.test",
        only = [],
        live_update = [
            sync("./spydk/js/src", "/app/spydk/js/src"),
            sync("./testing", "/app/testing"),
        ],
    )

    k8s_yaml_with_ns(encode_yaml_stream(set_env_in_jobs(read_yaml_stream("devnet/tests.yaml"), "NUM_PHYLAXS", str(num_phylaxs))))

    # separate resources to parallelize docker builds
    k8s_resource(
        "sdk-ci-tests",
        labels = ["ci"],
        trigger_mode = trigger_mode,
        resource_deps = [], # testing/sdk.sh handles waiting for spy, not having deps gets the build earlier
    )
    k8s_resource(
        "spydk-ci-tests",
        labels = ["ci"],
        trigger_mode = trigger_mode,
        resource_deps = [], # testing/spydk.sh handles waiting for spy, not having deps gets the build earlier
    )
    k8s_resource(
        "accountant-ci-tests",
        labels = ["ci"],
        trigger_mode = trigger_mode,
        resource_deps = [], # uses devnet-consts.json, but deltachain/contracts/tools/test_accountant.sh handles waiting for phylax, not having deps gets the build earlier
    )

if terra_classic:
    docker_build(
        ref = "terra-image",
        context = "./terra/devnet",
        dockerfile = "terra/devnet/Dockerfile",
    )

    docker_build(
        ref = "terra-contracts",
        context = "./terra",
        dockerfile = "./terra/Dockerfile",
    )

    k8s_yaml_with_ns("devnet/terra-devnet.yaml")

    k8s_resource(
        "terra-terrad",
        port_forwards = [
            port_forward(26657, name = "Terra RPC [:26657]", host = webHost),
            port_forward(1317, name = "Terra LCD [:1317]", host = webHost),
        ],
        labels = ["terra"],
        trigger_mode = trigger_mode,
    )

    k8s_resource(
        "terra-postgres",
        labels = ["terra"],
        trigger_mode = trigger_mode,
    )

    k8s_resource(
        "terra-fcd",
        resource_deps = ["terra-terrad", "terra-postgres"],
        port_forwards = [port_forward(3060, name = "Terra FCD [:3060]", host = webHost)],
        labels = ["terra"],
        trigger_mode = trigger_mode,
    )

if terra2 or deltachain:
    docker_build(
        ref = "cosmwasm_artifacts",
        context = ".",
        dockerfile = "./cosmwasm/Dockerfile",
        target = "artifacts",
    )

if terra2:
    docker_build(
        ref = "terra2-image",
        context = "./cosmwasm/deployment/terra2/devnet",
        dockerfile = "./cosmwasm/deployment/terra2/devnet/Dockerfile",
    )

    docker_build(
        ref = "terra2-deploy",
        context = "./cosmwasm/deployment/terra2",
        dockerfile = "./cosmwasm/Dockerfile.deploy",
    )

    k8s_yaml_with_ns("devnet/terra2-devnet.yaml")

    k8s_resource(
        "terra2-terrad",
        port_forwards = [
            port_forward(26658, container_port = 26657, name = "Terra 2 RPC [:26658]", host = webHost),
            port_forward(1318, container_port = 1317, name = "Terra 2 LCD [:1318]", host = webHost),
        ],
        labels = ["terra2"],
        trigger_mode = trigger_mode,
    )

    k8s_resource(
        "terra2-postgres",
        labels = ["terra2"],
        trigger_mode = trigger_mode,
    )

    k8s_resource(
        "terra2-fcd",
        resource_deps = ["terra2-terrad", "terra2-postgres"],
        port_forwards = [port_forward(3061, container_port = 3060, name = "Terra 2 FCD [:3061]", host = webHost)],
        labels = ["terra2"],
        trigger_mode = trigger_mode,
    )

if algorand:
    k8s_yaml_with_ns("devnet/algorand-devnet.yaml")

    docker_build(
        ref = "algorand-algod",
        context = "algorand/sandbox-algorand",
        dockerfile = "algorand/sandbox-algorand/images/algod/Dockerfile"
    )

    docker_build(
        ref = "algorand-indexer",
        context = "algorand/sandbox-algorand",
        dockerfile = "algorand/sandbox-algorand/images/indexer/Dockerfile"
    )

    docker_build(
        ref = "algorand-contracts",
        context = "algorand",
        dockerfile = "algorand/Dockerfile",
        ignore = ["algorand/test/*.*"]
    )

    k8s_resource(
        "algorand",
        port_forwards = [
            port_forward(4001, name = "Algod [:4001]", host = webHost),
            port_forward(4002, name = "KMD [:4002]", host = webHost),
            port_forward(8980, name = "Indexer [:8980]", host = webHost),
        ],
        labels = ["algorand"],
        trigger_mode = trigger_mode,
    )

if sui:
    k8s_yaml_with_ns("devnet/sui-devnet.yaml")

    docker_build(
        ref = "sui-node",
        target = "sui",
        context = ".",
        dockerfile = "sui/Dockerfile",
        ignore = ["./sui/sui.log*", "sui/sui.log*", "sui.log.*"],
        only = ["./sui"],
    )

    k8s_resource(
        "sui",
        port_forwards = [
            port_forward(9000, 9000, name = "RPC [:9000]", host = webHost),
            port_forward(5003, name = "Faucet [:5003]", host = webHost),
            port_forward(9184, name = "Prometheus [:9184]", host = webHost),
        ],
        labels = ["sui"],
        trigger_mode = trigger_mode,
    )

if near:
    k8s_yaml_with_ns("devnet/near-devnet.yaml")

    docker_build(
        ref = "near-node",
        context = "near",
        dockerfile = "near/Dockerfile",
        only = ["Dockerfile", "node_builder.sh", "start_node.sh", "README.md", "cert.pem"],
    )

    docker_build(
        ref = "near-deploy",
        context = "near",
        dockerfile = "near/Dockerfile.deploy",
        ignore = ["./test"]
    )

    k8s_resource(
        "near",
        port_forwards = [
            port_forward(3030, name = "Node [:3030]", host = webHost),
            port_forward(3031, name = "webserver [:3031]", host = webHost),
        ],
        labels = ["near"],
        trigger_mode = trigger_mode,
    )

if deltachain:
    docker_build(
        ref = "deltachaind-image",
        context = ".",
        dockerfile = "./deltachain/Dockerfile",
        build_args = {"num_phylaxs": str(num_phylaxs)},
        only = [],
        ignore = ["./deltachain/testing", "./deltachain/ts-sdk", "./deltachain/design", "./deltachain/vue", "./deltachain/build/deltachaind"],
    )

    docker_build(
        ref = "vue-export",
        context = ".",
        dockerfile = "./deltachain/Dockerfile.proto",
        target = "vue-export",
    )

    docker_build(
        ref = "deltachain-deploy",
        context = "./deltachain",
        dockerfile = "./deltachain/Dockerfile.deploy",
    )

    def build_deltachain_yaml(yaml_path, num_instances):
        deltachain_yaml = read_yaml_stream(yaml_path)

        # set the number of replicas in the StatefulSet to be num_phylaxs
        deltachain_set = set_replicas_in_statefulset(deltachain_yaml, "deltachain", num_instances)

        # add a Service for each deltachain instance
        services = []
        for obj in deltachain_set:
            if obj["kind"] == "Service" and obj["metadata"]["name"] == "deltachain-0":

                # make a Service for each replica so we can resolve it by name from other pods.
                # copy deltachain-0's Service then set the name and selector for the instance.
                for instance_num in list(range(1, num_instances)):
                    instance_name = 'deltachain-%s' % (instance_num)

                    # Copy the Service's properties to a new dict, by value, three levels deep.
                    # tl;dr - if the value is a dict, use a comprehension to copy it immutably.
                    service = { k: ({ k2: ({ k3:v3
                        for (k3,v3) in v2.items()} if type(v2) == "dict" else v2)
                        for (k2,v2) in v.items()} if type(v) == "dict" else v)
                        for (k,v) in obj.items()}

                    # add the name we want to be able to resolve via k8s DNS
                    service["metadata"]["name"] = instance_name
                    # add the name of the pod the service should connect to
                    service["spec"]["selector"] = { "statefulset.kubernetes.io/pod-name": instance_name }

                    services.append(service)

        return encode_yaml_stream(deltachain_set + services)

    deltachain_path = "devnet/deltachain.yaml"
    if num_phylaxs >= 2:
        # update deltachain's k8s config to spin up multiple instances
        k8s_yaml_with_ns(build_deltachain_yaml(deltachain_path, num_phylaxs))
    else:
        k8s_yaml_with_ns(deltachain_path)

    k8s_resource(
        "deltachain",
        port_forwards = [
            port_forward(1319, container_port = 1317, name = "REST [:1319]", host = webHost),
            port_forward(9090, container_port = 9090, name = "GRPC", host = webHost),
            port_forward(26659, container_port = 26657, name = "TENDERMINT [:26659]", host = webHost)
        ],
        labels = ["deltachain"],
        trigger_mode = trigger_mode,
    )

    k8s_resource(
        "deltachain-deploy",
        resource_deps = ["deltachain"],
        labels = ["deltachain"],
        trigger_mode = trigger_mode,
    )

if ibc_relayer:
    docker_build(
        ref = "ibc-relayer-image",
        context = ".",
        dockerfile = "./deltachain/ibc-relayer/Dockerfile",
        only = []
    )

    k8s_yaml_with_ns("devnet/ibc-relayer.yaml")

    k8s_resource(
        "ibc-relayer",
        port_forwards = [
            port_forward(7597, name = "HTTPDEBUG [:7597]", host = webHost),
        ],
        resource_deps = ["deltachain-deploy", "terra2-terrad"],
        labels = ["ibc-relayer"],
        trigger_mode = trigger_mode,
    )

if btc:
    k8s_yaml_with_ns("devnet/btc-localnet.yaml")

    docker_build(
        ref = "btc-node",
        context = "bitcoin",
        dockerfile = "bitcoin/Dockerfile",
        target = "bitcoin-build",
    )

    k8s_resource(
        "btc",
        port_forwards = [
            port_forward(18556, name = "RPC [:18556]", host = webHost),
        ],
        labels = ["btc"],
        trigger_mode = trigger_mode,
    )

if aptos:
    k8s_yaml_with_ns("devnet/aptos-localnet.yaml")

    docker_build(
        ref = "aptos-node",
        context = "aptos",
        dockerfile = "aptos/Dockerfile",
        target = "aptos",
    )

    k8s_resource(
        "aptos",
        port_forwards = [
            port_forward(8080, name = "RPC [:8080]", host = webHost),
            port_forward(6181, name = "FullNode [:6181]", host = webHost),
            port_forward(8081, name = "Faucet [:8081]", host = webHost),
        ],
        labels = ["aptos"],
        trigger_mode = trigger_mode,
    )
