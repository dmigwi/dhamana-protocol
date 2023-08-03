const chatContract = artifacts.require("../chat/ChatContract");

module.exports = function(deployer) {
  deployer.deploy(chatContract);
};