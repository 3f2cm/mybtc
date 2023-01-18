# mybtc

mybtc is a simple Bitcoin CLI for me in Go.
It can just generate transactions with signs and do a few more small things.

## generate WIF and get its address

```shell
$ mybtc wif generate | tee mywallet
93ConuWWVWBh52nqWhEqupuGsB6SortEpUHpR51QTxKRcLjg59J
$ mybtc wif address < mywallet
n2r6MoqD3CdA8F4ytVU87Pgf8o12Qwk9ga
```

## generate transaction with signs

You can find sample inputs as `cmd/test_data/sample_*_input.json`.
The inputs consist of

- transaction INs of
    - pvevious transaction ID
    - previous transaction vout
    - provious transaction ScriptPublicKey
    - corresponding WIF, and
- transaction OUTs of
    - address
    - value to be sent

```shell
$ mybtc tx generate < cmd/test_data/sample_1_input.json
01000000012d563d01940861f15f6edcff412b1603a0a2f5ce561b4417e557f9c997ec4e20010000008a4730440220327111114de4ceb65143f51f73b5915512f211f31bacb3934b19312a7dfd33a6022047aa55cb056fe14794fd36dde4f6ed2553eac5d19b852a2987ab9c389e060d6701410425adade9f702a4c1e7f312ed7eb9507a6b70a6bafe6c48092137eb991d5b29eb9374edc1c6f83e0b22d5c00f26d0b163a466c45ed814c2b9b929e87ab47e8551ffffffff0200530700000000001976a91427e49532bfeae7a40d878aa5fd2699fe9729cb2588ac60e31600000000001976a91414fc9b2d74e76f4d7463f2a04e06040d3128e22d88ac00000000
```

## retrieve UXTO info of an adress

Creating the input for `mybtc tx generate` is a messy work.
To reduce the annoyance a bit, there is a simple client for blockstream.
However, it works only for addresses with short history.

```shell
$ uxto-summary mx1UsJZ9aS1z7YrebihuxsTYdUthkcHWTv | jq
[
  {
    "txid": "3bfd76132d9e204087f8c9c3d1831661e23b06d65df33f2568988b035179092f",
    "vout": 1,
    "scriptpubkey": "76a914b4e72e4582c8ef7f510447d90f0249ad8b29b6b788ac",
    "value": 267587
  },
  {
    "txid": "fbc6d0d353ccfeaab73917f79c558002fffda718e9f0f58d80d0b9dea8c76297",
    "vout": 1,
    "scriptpubkey": "76a914b4e72e4582c8ef7f510447d90f0249ad8b29b6b788ac",
    "value": 6047
  },
  {
    "txid": "0e9b43d498b009a5b7a5da37658e11c125d8dbd9db023171005b6e8fc03c5dde",
    "vout": 1,
    "scriptpubkey": "76a914b4e72e4582c8ef7f510447d90f0249ad8b29b6b788ac",
    "value": 19460
  },
  {
    "txid": "447d568a28512409c5e6ed040263c2c976b20caa3e5a7fa71bbba8f1a7e7b308",
    "vout": 3,
    "scriptpubkey": "76a914b4e72e4582c8ef7f510447d90f0249ad8b29b6b788ac",
    "value": 65403237
  }
]

$ uxto-summary mx1UsJZ9aS1z7YrebihuxsTYdUthkcHWTv | jq '[.[].value] | add'
65343531

$ uxto-summary mfWxJ45yp2SFn7UciZyNpvDKrzbhyfKrY8
2023/01/17 21:57:33 HTTP request for mfWxJ45yp2SFn7UciZyNpvDKrzbhyfKrY8 failed: code: 400, body: Too many history entries
```
