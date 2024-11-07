package signer

const SecretKey = ""

var DefaultSigner ISigner = NewMd5Signer(SecretKey)
