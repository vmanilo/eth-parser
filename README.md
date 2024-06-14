# Simple ETH parser

## API

1. Get last parsed block 

GET localhost:8080/api/current-block

Response
```json
{
	"blockNumber": 20092372
}
```

2. Add address to observer

POST localhost:8080/api/subscribe
```json
{
    "address": "0xDef1C0ded9bec7F1a1670819833240f027b25EfF"
}
```

Response
```json
{
  "ok": true
}
```

3. List of inbound or outbound transactions for an address

GET localhost:8080/api/transactions/0xDef1C0ded9bec7F1a1670819833240f027b25EfF

Response
```json
{
  "transactions": [
    {
      "hash": "0x8ef79e3230927ec83d31f109d5702baf689bca3fcbd99137f84b929eabfece61",
      "blockNumber": 20092610,
      "type": "inbound"
    }
  ]
}
```