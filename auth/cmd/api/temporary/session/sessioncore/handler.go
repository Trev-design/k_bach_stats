package sessioncore

import (
	"auth_server/cmd/api/temporary/session/refreshpayload"
	"encoding/base64"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (session *Session) GetVerifyData(cookie string) (account uuid.UUID, sessionID, verify string, err error) {
	account, id, err := session.getIDsFromCookie(cookie)
	if err != nil {
		return
	}

	verify, err = session.getVerifyFromSession(id)
	if err != nil {
		return
	}

	return account, id, verify, nil
}

func (session *Session) DeleteVerifySession(id string) error {
	return session.store.DeleteSessionPayload("verify", id)
}

func (session *Session) SetVerifyData(accountID string) (number, cookie string, err error) {
	// make verify
	number, id, err := session.setVerifyInSession()
	if err != nil {
		return
	}

	cookie, err = session.getCookie(accountID, id)
	if err != nil {
		return
	}

	return
}

func (session *Session) SetRefreshData(accountID, ip, userAgent string) (string, error) {
	token, err := getRefreshToken()
	if err != nil {
		return "", err
	}

	timestamp := time.Now().UTC()

	refreshToken, err := refreshpayload.NewRefreshDataBuilder().
		Account(accountID).
		IPAddress(ip).
		UserAgent(userAgent).
		Token(token).
		Build()
	if err != nil {
		return "", err
	}

	sessionID, err := session.store.SetSessionPayload("refresh", string(refreshToken))
	if err != nil {
		return "", err
	}

	cookieData := generateRefreshCookie(token, sessionID, accountID)
	encryptedCookieData, err := session.refreshCrypt.EncryptPayload([]byte(cookieData), timestamp)
	if err != nil {
		return "", err
	}

	cookie := strings.Join([]string{encryptedCookieData, timestamp.Format(time.RFC3339Nano)}, " ")
	encodedCookie := base64.RawURLEncoding.EncodeToString([]byte(cookie))

	return encodedCookie, nil
}

func (session *Session) VerifyRefreshData(cookie, ip, useragent string) (string, error) {
	separatedCookieData, err := session.getCookieData(cookie)
	if err != nil {
		return "", err
	}

	if err := session.verifyCookiePayload(separatedCookieData, ip, useragent); err != nil {
		return "", err
	}

	return session.SetRefreshData(separatedCookieData.accountID, ip, useragent)
}

func (session *Session) RemoveRefreshData(cookie, ip, useragent string) error {
	separatedcookieData, err := session.getCookieData(cookie)
	if err != nil {
		return err
	}

	if err := session.verifyCookiePayload(separatedcookieData, ip, useragent); err != nil {
		return err
	}

	return session.store.DeleteSessionPayload("refresh", separatedcookieData.sessionID)
}

func (session *Session) CloseSession() error {
	if err := session.cookieCrypt.CloseCrypto(); err != nil {
		return err
	}

	if err := session.verifyCrypt.CloseCrypto(); err != nil {
		return err
	}

	return session.store.CloseRedisStore()
}

func (session *Session) HandleBackground() {
	go session.cookieCrypt.ComputeRotateInterval()
	go session.verifyCrypt.ComputeRotateInterval()
}
