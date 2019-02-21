const fs = require('fs')
const path = require('path')
const utils = require('ethers').utils
const program = require('commander')

const { getNetworkID } = require('./utils/helpers')

const network = process.argv[2]
if (!network) console.log('Usage: node approve_tokens {network}')

const networkID = getNetworkID(network) || '8888'

require('dotenv').config()

program
  .version('0.1.0')
  .option('-p, --truffle-build-path [value]', 'Truffle build path')
  .parse(process.argv)

const TRUFFLE_BUILD_PATH = path.resolve(
  program.truffleBuildPath || '../dex-smart-contract/build/contracts',
)

const contractConfig = require(path.resolve(TRUFFLE_BUILD_PATH, '../../config'))

const ignoreFilesPattern = /(?:RewardCollector|RewardPools|Migrations|Owned|SafeMath)\.json$/

const contracts = {
  [networkID]: {},
}

const files = fs.readdirSync(TRUFFLE_BUILD_PATH)

files
  .filter(file => !ignoreFilesPattern.test(file))
  .forEach((file) => {
    let address
    let symbol
    const json = JSON.parse(
      fs.readFileSync(`${TRUFFLE_BUILD_PATH}/${file}`, 'utf8'),
    )

    if (json.networks[networkID]) {
      symbol = file.slice(0, -5)
      if (symbol === 'WETH9') symbol = 'WETH'
      if (contractConfig.tokens.includes(symbol) || symbol === 'Exchange') {
        address = json.networks[networkID].address
        contracts[networkID][symbol] = utils.getAddress(address)
      }
    }
  })

console.log(contracts)
fs.writeFileSync(
  'utils/db/utils/addresses.json',
  JSON.stringify(contracts, null, 2),
  'utf8'
)