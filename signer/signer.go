package signer

type SignerMethod string

const (
	SignerMd5   SignerMethod = "md5"
	SignerSh256 SignerMethod = "sh256"
)

type ISigner interface {
	GenerateSignature(params map[string]string, expireAt int64) (string, error)
	VerifySignature(params map[string]string, signature string, expireAt int64) (bool, error)
	Method() SignerMethod
}
