## Generate genesis block

`go run main.go genesis -cbf dex-contracts/build/contracts -out dex-protocol/OrderBook`

## Generate tokens, accounts, pairs seed data

```go
go run main.go tokens -cr contract-results.txt
go run main.go accounts -cr contract-results.txt
go run main.go pairs -cr contract-results.txt
```

To use image instead of icon, append this

```json
"image":{
  "url": "https://tomochain.com/file/2018/08/logo.png",
  "meta": null
}
```
