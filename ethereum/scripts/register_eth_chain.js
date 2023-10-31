// run this script with truffle exec

const jsonfile = require("jsonfile");
const TokenBridge = artifacts.require("TokenBridge");
const TokenImplementation = artifacts.require("TokenImplementation");
const NFTBridge = artifacts.require("NFTBridgeEntrypoint");
const BridgeImplementationFullABI = jsonfile.readFileSync("../build/contracts/BridgeImplementation.json").abi
const tokenBridgeVAA = process.env.REGISTER_TOKEN_BRIDGE_VAA
const nFTBridgeVAA = process.env.REGISTER_NFTBRIDGE_VAA

module.exports = async function (callback) {
    try {
        const accounts = await web3.eth.getAccounts();

        // Register the ETH endpoint
        const tokenBridge = new web3.eth.Contract(
            BridgeImplementationFullABI,
            TokenBridge.address
        );
        const nftBridge = new web3.eth.Contract(
            BridgeImplementationFullABI,
            NFTBridge.address
        );

        // Register the token bridge endpoints
        console.log("Registering Token Bridges...");
        await tokenBridge.methods.registerChain("0x" + tokenBridgeVAA).send({
            value: 0,
            from: accounts[0],
            gasLimit: 2000000,
        });

        // Register the NFT bridge endpoints
        console.log("Registering NFT Bridges...");
        await nftBridge.methods.registerChain("0x" + nFTBridgeVAA).send({
            value: 0,
            from: accounts[0],
            gasLimit: 2000000,
        });

        console.log("Finished registering all Bridges...");

        callback();
    }
    catch (e) {
        console.log(e)
        callback(e);
    }
}

