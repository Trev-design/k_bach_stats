package refreshpayload

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
)

type RefreshData struct {
	Account   string `json:"account_id"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Token     string `json:"token"`
	Revoked   bool   `json:"revoked"`
}

type RefreshDataBuilder struct {
	account   string
	ipAddress string
	userAgent string
	token     string
	revoked   bool
}

type RefreshPayload string

func NewRefreshDataBuilder() *RefreshDataBuilder {
	return &RefreshDataBuilder{}
}

func (builder *RefreshDataBuilder) Account(account string) *RefreshDataBuilder {
	builder.account = account
	return builder
}

func (builder *RefreshDataBuilder) IPAddress(ipAddress string) *RefreshDataBuilder {
	builder.ipAddress = ipAddress
	return builder
}

func (builder *RefreshDataBuilder) UserAgent(userAgent string) *RefreshDataBuilder {
	builder.userAgent = userAgent
	return builder
}

func (builder *RefreshDataBuilder) Token(token string) *RefreshDataBuilder {
	builder.token = token
	return builder
}

func (builder *RefreshDataBuilder) Build() (RefreshPayload, error) {
	data := new(RefreshData)

	accountHash, err := GenerateHash(builder.account)
	if err != nil {
		return "", err
	}

	tokenHash, err := GenerateHash(builder.token)
	if err != nil {
		return "", err
	}

	ipHash, err := GenerateHash(builder.ipAddress)
	if err != nil {
		return "", err
	}

	uaHash, err := GenerateHash(builder.userAgent)
	if err != nil {
		return "", err
	}

	data.Account = accountHash
	data.Token = tokenHash
	data.IPAddress = ipHash
	data.UserAgent = uaHash
	data.Revoked = builder.revoked

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return RefreshPayload(base64.RawURLEncoding.EncodeToString(bytes)), nil
}

func (payload RefreshPayload) GetData() (*RefreshData, error) {
	data := new(RefreshData)

	payloadBytes, err := base64.RawURLEncoding.DecodeString(string(payload))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(payloadBytes, data); err != nil {
		return nil, err
	}

	return data, nil
}

func GenerateHash(payload string) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write([]byte(payload)); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
