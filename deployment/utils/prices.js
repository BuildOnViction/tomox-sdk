const faker = require('faker')
const dateFns = require('date-fns')
const volatility = 0.00001
// not much different
const minPrice = 999000
const maxPrice = 1000000

const rand = () => faker.random.number(100) / 100
const randInt = (min, max) => Math.floor(Math.random() * (max - min + 1) + min)

const computeNextPrice = oldPrice => {
    let changePercent = 2 * volatility * rand()
    if (changePercent > volatility) changePercent -= 2 * volatility
    let changeAmount = oldPrice * changePercent
    let newPrice = oldPrice + changeAmount

    return newPrice
}

const generateTimestamps = (start, end, interval) => {
    start = start || new Date(2018, 1, 1).getTime()
    end = end || Date.now()
    interval = interval || 'hour'

    let intervalInSeconds

    switch (interval) {
        case 'second':
            intervalInSeconds = 1 * 1000
            break
        case 'minute':
            intervalInSeconds = 60 * 1000
            break
        case 'hour':
            intervalInSeconds = 60 * 60 * 1000
            break
        case 'day':
            intervalInSeconds = 60 * 60 * 24 * 1000
            break
        default:
            throw new Error('Error')
    }

    let currentTimestamp = start
    let timestamps = []

    while (currentTimestamp < end) {
        currentTimestamp += intervalInSeconds
        timestamps.push(currentTimestamp)
    }

    return timestamps
}

const generatePrices = (timestamps, initialPrice) => {
    initialPrice = initialPrice || randInt(minPrice, maxPrice)

    let pricesArray = [{ timestamp: timestamps[0], price: initialPrice }]

    let result = timestamps.slice(1).reduce((result, timestamp) => {
        let nextPrice = computeNextPrice(result[result.length - 1].price)
        pricesArray.push({ timestamp: timestamp, price: nextPrice })
        return pricesArray
    }, pricesArray)

    return result
}

const generateRandomPricepointRange = () => {
    let a = randInt(minPrice, maxPrice)
    let b = randInt(minPrice, maxPrice)
    let min = Math.min(a, b)
    let max = Math.max(a, b)
    return { min, max }
}

const generatePricingData = ({ start, end, interval, initialPrice }) => {
    start = start || new Date(2016, 1, 1).getTime()
    end = end || Date.now()
    initialPrice = initialPrice || randInt(minPrice, maxPrice)
    // initialPrice = initialPrice || faker.random.number(maxPrice)
    interval = interval || 'hour'

    let timestamps = generateTimestamps(start, end, interval)
    let pricingData = generatePrices(timestamps, initialPrice)

    return pricingData
}

const interpolatePrice = (pricingData, timestamp) => {
    let start = pricingData[0].timestamp
    let timestampInterval = pricingData[1].timestamp - pricingData[0].timestamp
    let numberOfIntervalsFromStart = Math.floor(
        (timestamp - start) / timestampInterval,
    )
    let previousTimestampIndex = numberOfIntervalsFromStart
    let nextTimestampIndex = numberOfIntervalsFromStart + 1

    let previousPrice = pricingData[previousTimestampIndex]
        ? pricingData[previousTimestampIndex].price
        : pricingData[0].price
    let nextPrice = pricingData[nextTimestampIndex]
        ? pricingData[nextTimestampIndex].price
        : pricingData[pricingData.length - 1].price
    let previousTimestamp = pricingData[previousTimestampIndex]
        ? pricingData[previousTimestampIndex].timestamp
        : pricingData[0].timestamp
    let nextTimestamp = pricingData[nextTimestampIndex]
        ? pricingData[nextTimestampIndex].timestamp
        : pricingData[pricingData.length - 1].timestamp
    const inflation =
        ((previousPrice > nextPrice ? -1 : 1) *
            volatility *
            (timestamp - previousTimestamp)) /
        (nextTimestamp - previousTimestamp)

    let interpolatedPrice = previousPrice + inflation

    // console.log({ previousPrice, inflation });
    return interpolatedPrice
}

module.exports = {
    generatePricingData,
    interpolatePrice,
    generateRandomPricepointRange,
}
