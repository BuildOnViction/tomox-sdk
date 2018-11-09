## Generate genesis block

`go run main.go genesis -cbf dex-contracts/build/contracts -out dex-protocol/OrderBook`

## Generate tokens, accounts, pairs seed data

```go
go run main.go seeds -ccf ../../../client/src/config
```

To use image instead of icon, append this

```json
"image":{
  "url": "https://tomochain.com/file/2018/08/logo.png",
  "meta": null
}
```
