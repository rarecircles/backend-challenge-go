# RareCircles Coding Challenge implemented by Tyler

## Requirements

To try the bazel stuff, install bazel through a package manager

Otherwise just Go

Tested in Linux.
### To start
To run using bazel '''bazel run //token_api/:cmd'''
To run all tests using bazel '''bazel test //token_api/...'''

To run without using bazel '''go run ./token_api/cmd'''. The env var ADDRESS_DATA_PATH will need to be set to "./internal/dataloader/data/addresses.jsonl"

### Useful info
By default server runs on 0.0.0.0:8000 with prometheus metrics running on 0.0.0.0:9000

### Endpoint info
**Near match**

    GET /tokens?q=rare

```json
{
  "tokens": [
    {
      "name": "RareCircles",
      "symbol": "RCI",
      "address": "0xdbf1344a0ff21bc098eb9ad4eef7de0f9722c02b",
      "decimals": 18,
      "totalSupply": 1000000000,
    },
    {
      "name": "Rareible",
      "symbol": "RRI",
      "address": "e9c8934ebd00bf73b0e961d1ad0794fb22837206",
      "decimals": 9,
      "totalSupply": 100,
    },
  ]
}
```

**No match**

    GET /tokens?q=ThereIsNoTokenHere

```json
{
  "tokens": []
}
```
