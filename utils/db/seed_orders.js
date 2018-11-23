const utils = require('ethers').utils
const MongoClient = require('mongodb').MongoClient
const faker = require('faker')
const argv = require('yargs').argv
const { generatePricingData , interpolatePrice } = require('../../utils/prices')

const mongoUrl = argv.mongo_url || 'mongodb://localhost:27017'

let { addresses } = require('./addresses.json')
let exchangeAddress = "0x7400d4d4263a3330beeb2a0d2674f0456054f217"
let minTimeStamp = 1500000000000
let maxTimeStamp = 1520000000000
let minAmount = 0.1
let maxAmount = 10000
let minPrice = 100
let maxPrice = 100000
let ether = 1e18

let orderStatuses = ['NEW', 'OPEN', 'CANCELLED', 'PARTIALLY_FILLED', 'FILLED', 'EXECUTED']
let tradeStatuses = ['PENDING', 'SUCCESS', 'ERROR', 'INVALID']
let orderTypes = ['MARKET', 'LIMIT']

let orderWeightedStatuses = [
  {
    name: 'NEW',
    probability: 0.0
  },
  {
    name: 'OPEN',
    probability: 0.20
  },
  {
    name: 'CANCELLED',
    probability: 0.05
  },
  {
    name: 'ERROR',
    probability: 0.05
  },
  {
    name: 'PARTIALLY_FILLED',
    probability: 0.20
  },
  {
    name: 'FILLED',
    probability: 0.20
  },
  {
    name: 'EXECUTED',
    probability: 0.30
  },
]

let tradeWeightedStatuses = [
  {
    name: 'PENDING',
    probability: 0.20
  },
  {
    name: 'SUCCESS',
    probability: 0.70
  },
  {
    name: 'ERROR',
    probability: 0.05
  },
  {
    name: 'INVALID',
    probability: 0.05
  }
]

let orderLevels = orderWeightedStatuses.reduce((result, current) => {
  let len = result.length
  len > 0 ? result.push(result[result.length - 1] + current.probability) : result.push(current.probability)
  return result
 }, [])
 .map(elem => elem * 100)

 let tradeLevels = tradeWeightedStatuses.reduce((result, current) => {
  let len = result.length
  len > 0 ? result.push(result[result.length - 1] + current.probability) : result.push(current.probability)
  return result
 }, [])
 .map(elem => elem * 100)

const randInt = (min, max) => Math.floor(Math.random() * (max - min + 1) + min)
const randomSide = () => (randInt(0, 1) === 1 ? 'BUY' : 'SELL')
const randomOrderType = () => orderTypes[randInt(0, orderTypes.length -1 )]
const randomPair = () => pairs[randInt(0, pairs.length-1)]
const randomFee = () => rand(10000, 100000)
const randomHash = () => utils.sha256(utils.randomBytes(100))

const randomBigAmount = () => {
  let ether = utils.bigNumberify("1000000000000000000")
  let amount = utils.bigNumberify(randInt(0, 100000))
  let bigAmount = amount.mul(ether).div("100").toString()
  return bigAmount
}

const randomAmount = () => rand(minAmount, maxAmount)
const randomRatio = () => rand(0, 1)
const randomTimestamp = () => randInt(minTimeStamp, maxTimeStamp)
const randomPrice = () => rand(minPrice, maxPrice)

const randomAddress = () => randomHash().slice(0, 42);
const randomElement = (arr) => arr[randInt(0, arr.length-1)]

const randomPricepointRange = () => {
  let a = randInt(10000, 1000000000)
  let b = randInt(10000, 1000000000)
  let min = Math.min(a, b)
  let max = Math.max(a, b)
  return { min, max }
}

const randomQuoteToken = (quotes) => quotes[randInt(0, len(quotes)-1)]
const randomToken = (tokens) => tokens[randInt(0, len(tokens)-1)]


const randomOrderStatus = () => {
  let nb = randInt(0, 100)

  switch(true) {
    case (nb < orderLevels[0]):
      return orderWeightedStatuses[0].name
      break
    case (nb < orderLevels[1]):
      return orderWeightedStatuses[1].name
      break
    case (nb < orderLevels[2]):
      return orderWeightedStatuses[2].name
      break
    case (nb < orderLevels[3]):
      return orderWeightedStatuses[3].name
      break
    case (nb < orderLevels[4]):
      return orderWeightedStatuses[4].name
      break
    case (nb < orderLevels[5]):
      return orderWeightedStatuses[5].name
      break
    default:
      return orderWeightedStatuses[6].name
  }
}

const randomTradeStatus = () => {
  let nb = randInt(0, 100)
  switch(true) {
    case (nb < tradeLevels[0]):
      return tradeWeightedStatuses[0].name
      break
    case (nb < tradeLevels[1]):
      return tradeWeightedStatuses[1].name
      break
    case (nb < tradeLevels[2]):
      return tradeWeightedStatuses[2].name
      break
    default:
      return tradeWeightedStatuses[3].name
  }
}

const seed = async () => {
    let orders = []
    const client = await MongoClient.connect(url, { useNewUrlParser: true })
    const db = client.db('proofdex')

    const docs = await db.collection('pairs')
      .find(
        {},
        { baseTokenSymbol: 1,
          baseTokenAddress: 1,
          quoteTokenSymbol: 1,
          quoteTokenAddress: 1,
          pairMultiplier: 1,
        }
      )
      .toArray()

    let pairs = []
    docs.forEach(pair => {
      let { min, max } = randomPricepointRange()
      pairs.push({
        baseTokenAddress: pair.baseTokenAddress,
        baseTokenSymbol: pair.baseTokenSymbol,
        quoteTokenAddress: pair.quoteTokenAddress,
        quoteTokenSymbol: pair.quoteTokenSymbol,
        priceMultiplier: pair.priceMultiplier,
        minPricepoint: min,
        maxPricepoint: max,
        averagePricePoint: randInt((min + (max+min)/2)/2, (max + (max+min)/2)/2)
      })
    })


    //we choose a limited number of user accounts
    addresses = addresses.slice(0,4)

      for (let i = 0; i < 20000; i++) {
        let pair = randomElement(pairs)
        let side = randomSide()
        let baseToken = pair.baseTokenAddress
        let quoteToken = pair.quoteTokenAddress
        let hash = randomHash()
        let status = randomOrderStatus()
        let amount = randomBigAmount()
        let pricepoint = (side == "BUY") ? String(randInt(pair.minPricepoint, pair.averagePricePoint)) : String(randInt(pair.averagePricePoint, pair.maxPricepoint))
        let userAddress = randomElement(addresses)
        let pairName = `${pair.baseTokenSymbol}/${pair.quoteTokenSymbol}`
        let makeFee = 0
        let takeFee = 0
        let filledAmount
        let createdAt = new Date(faker.fake("{{date.recent}}"))


        switch(status) {
          case "OPEN":
            filledAmount = "0"
            break
          case "NEW":
            filledAmount = "0"
            break
          case "PARTIALLY_FILLED":
            filledAmount = String(randInt(0, amount))
            break
          case "FILLED":
            filledAmount = amount
            break
          case "INVALID":
            filledAmount = "0"
            break
          case "ERROR":
            filledAmount = "0"
            break
          default:
          filledAmount = "0"
        }

        let order = {
          exchangeAddress: utils.getAddress(exchangeAddress),
          userAddress: utils.getAddress(userAddress),
          baseToken: utils.getAddress(baseToken),
          quoteToken: utils.getAddress(quoteToken),
          pairName,
          hash,
          side,
          status,
          makeFee,
          takeFee,
          amount,
          pricepoint,
          filledAmount,
          createdAt
        }

        orders.push(order)
      }


    const ordersInsertResponse = await db.collection('orders').insertMany(orders)


    client.close()
}

seed()