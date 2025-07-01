package sessioncore

import (
	"auth_server/cmd/api/temporary/session/refreshpayload"
	"encoding/base64"
	"errors"
	"log"
	"math/rand/v2"
	"strings"
	"time"

	random "crypto/rand"

	"github.com/google/uuid"
)

type refreshCookieData struct {
	sessionID     string
	accountID     string
	cookiePayload string
}

func (session *Session) verifyCookiePayload(cookieCreds *refreshCookieData, ip, useragent string) error {
	token, err := session.store.GetSessionPayload("refresh", cookieCreds.sessionID)
	if err != nil {
		return err
	}

	tokenData, err := refreshpayload.RefreshPayload(token).GetData()
	if err != nil {
		return err
	}

	return checkRefreshToken(
		tokenData,
		cookieCreds.cookiePayload,
		ip,
		useragent,
		cookieCreds.accountID,
	)
}

func (session *Session) getCookieData(cookie string) (*refreshCookieData, error) {
	cokieData, timestamp, err := separatedData(cookie)
	if err != nil {
		return nil, err
	}

	decryptedCookieData, err := session.refreshCrypt.DecryptPayload(cokieData, timestamp)
	if err != nil {
		return nil, err
	}

	return separatedRefreshCookieData(decryptedCookieData)
}

func (session *Session) setVerifyInSession() (string, string, error) {
	number := makeVerify()
	//timestamp, err := getFormatedTimeStamp()
	//if err != nil {
	//	return "", err
	//}

	timestamp := time.Now().UTC()

	encrypted, err := session.verifyCrypt.EncryptPayload(number, timestamp)
	if err != nil {
		return "", "", err
	}

	payload := strings.Join([]string{encrypted, timestamp.Format(time.RFC3339Nano)}, " ")
	encoded := base64.RawURLEncoding.EncodeToString([]byte(payload))

	id, err := session.store.SetSessionPayload("verify", encoded)
	if err != nil {
		return "", "", err
	}

	return string(number), id, nil
}

func (session *Session) getCookie(accountID, id string) (string, error) {
	//timestamp, err := getFormatedTimeStamp()
	//if err != nil {
	//	return "", err
	//}

	timestamp := time.Now().UTC()
	ids := strings.Join([]string{accountID, id}, ":")
	encrypted, err := session.cookieCrypt.EncryptPayload([]byte(ids), timestamp)
	if err != nil {
		return "", err
	}

	payload := strings.Join([]string{encrypted, timestamp.Format(time.RFC3339Nano)}, " ")
	encoded := base64.RawURLEncoding.EncodeToString([]byte(payload))

	return encoded, nil
}

func (session *Session) getIDsFromCookie(cookie string) (uuid.UUID, string, error) {
	idData, timestamp, err := separatedData(cookie)
	if err != nil {
		return uuid.Nil, "", err
	}

	ids, err := session.cookieCrypt.DecryptPayload(idData, timestamp)
	if err != nil {
		return uuid.Nil, "", err
	}

	return getSeparatedIDs(ids)
}

func (session *Session) getVerifyFromSession(id string) (string, error) {
	verifyPayload, err := session.store.GetSessionPayload("verify", id)
	if err != nil {
		return "", err
	}

	verify, timestamp, err := separatedData(verifyPayload)
	if err != nil {
		return "", err
	}

	return session.verifyCrypt.DecryptPayload(verify, timestamp)
}

func makeVerify() []byte {
	verifyNumber := make([]byte, 7)

	for index := range verifyNumber {
		verifyNumber[index] = byte(rand.IntN(10)) + '0'
	}

	return verifyNumber
}

//func getFormatedTimeStamp() (time.Time, error) {
//	timestamp := time.Now().UTC()
//	timestring := timestamp.Format(time.RFC3339Nano)
//	return time.Parse(time.RFC3339Nano, timestring)
//}

func separatedData(payload string) (string, time.Time, error) {
	cookieBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return "", time.Time{}, err
	}

	segments := strings.Split(string(cookieBytes), " ")
	if len(segments) != 2 {
		return "", time.Time{}, errors.New("invalid payloads")
	}

	timestamp, err := time.Parse(time.RFC3339Nano, string(segments[1]))
	if err != nil {
		return "", time.Time{}, errors.New("invalid timestamp")
	}

	return string(segments[0]), timestamp, nil
}

func separatedRefreshCookieData(payload string) (*refreshCookieData, error) {
	cookieBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	segments := strings.Split(string(cookieBytes), ":")
	log.Println(len(segments))
	if len(segments) != 3 {
		return nil, errors.New("invalid payload")
	}

	refreshData := new(refreshCookieData)
	refreshData.cookiePayload = segments[0]
	refreshData.sessionID = segments[1]
	refreshData.accountID = segments[2]

	return refreshData, nil
}

func generateRefreshCookie(token, id, account string) string {
	cookie := strings.Join([]string{token, id, account}, ":")
	return base64.RawURLEncoding.EncodeToString([]byte(cookie))
}

func getRefreshToken() (string, error) {
	token := make([]byte, 64)
	if _, err := random.Read(token); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(token), nil
}

func checkRefreshToken(
	token *refreshpayload.RefreshData,
	cookie, ip, userData, account string,
) error {
	tokenHash, err := refreshpayload.GenerateHash(cookie)
	if err != nil {
		return err
	}

	ipHash, err := refreshpayload.GenerateHash(ip)
	if err != nil {
		return err
	}

	uaHash, err := refreshpayload.GenerateHash(userData)
	if err != nil {
		return err
	}

	accountHash, err := refreshpayload.GenerateHash(account)
	if err != nil {
		return err
	}

	if accountHash != token.Account ||
		uaHash != token.UserAgent ||
		ipHash != token.IPAddress ||
		tokenHash != token.Token {
		return errors.New("invalid payload")
	}

	return nil
}

func getSeparatedIDs(ids string) (uuid.UUID, string, error) {
	segments := strings.Split(ids, ":")
	if len(segments) != 2 {
		return uuid.Nil, "", errors.New("invalid session")
	}

	accountID, err := uuid.Parse(segments[0])
	if err != nil {
		return uuid.Nil, "", errors.New("invalid session")
	}

	return accountID, segments[1], nil
}
