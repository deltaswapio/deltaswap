require('dotenv').config({ path: "../.env" });

const Deltaswap = artifacts.require("Deltaswap");
const MockBatchedVAASender = artifacts.require("MockBatchedVAASender");

module.exports = async function (deployer, network, accounts) {

    await deployer.deploy(MockBatchedVAASender)

    const contract = new web3.eth.Contract(MockBatchedVAASender.abi, MockBatchedVAASender.address);

    await contract.methods.setup(
        Deltaswap.address
    ).send({from: accounts[0]})
};
