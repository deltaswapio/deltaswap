// run this script with truffle exec

const jsonfile = require("jsonfile");
const {getRandomBytes} = require("truffle/build/5521.bundled");
const Deltaswap = artifacts.require("Deltaswap");
const ImplementationFullABI = jsonfile.readFileSync("../build/contracts/Implementation.json").abi;
const messageFeeVAA = process.env.MESSAGE_FEE_VAA

module.exports = async function (callback) {
    try {
        const accounts = await web3.eth.getAccounts();

        // Register the Contract endpoint
        const deltaswap = new web3.eth.Contract(
            ImplementationFullABI,
            Deltaswap.address
        );

        console.log("Setting message fee");
        await deltaswap.methods.submitSetMessageFee("0x" + messageFeeVAA).send({
            value: 0,
            from: accounts[0],
            gasLimit: 2000000
        });

        console.log("Finished setting message fee...");

        callback();
    }
    catch (e) {
        console.log(e)
        callback(e);
    }
}

