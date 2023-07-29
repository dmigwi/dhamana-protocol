const chatContract = artifacts.require("../chat/chatContract");

module.exports = function(deployer) {
  deployer.deploy(chatContract);
};