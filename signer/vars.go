package signer

const SecretKey = "783420-1928232%$^%3$@!turing@turingdance.com"

var DefaultSigner ISigner = NewMd5Signer(SecretKey)
