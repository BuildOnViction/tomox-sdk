const { providers, Contract } = require('ethers')
const { ERC20 } = require('../utils/abis')

const queryTokenBalances = async () => {
  try {
    let provider = new providers.JsonRpcProvider('http://localhost:8545')
    let accountAddress = '0xF069080F7acB9a6705b4a51F84d9aDc67b921bDF'
    const token = new Contract('0x9A8531C62D02AF08cf237Eb8aecae9DbCb69B6Fd', ERC20, provider)
    const balance = await token.balanceOf(accountAddress)
    console.log(`${balance}`)
  } catch (err) {
    console.log(err)
  }
}

queryTokenBalances()