package helper

type JwtHelper struct {
	Secret     string
	Issuer     string
	Expiration int64
}

type JwtClaim struct {
}

func (jh *JwtHelper) GenerateToken(data []interface{}) string {
	return ""
}
