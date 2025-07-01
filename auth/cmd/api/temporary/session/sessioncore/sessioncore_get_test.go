package sessioncore_test

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSetRefresh(t *testing.T) {
	cookie, err := session.SetRefreshData(uuid.NewString(), "127.0.0.1", "some ua to check")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("the cookie: %s", cookie)
}

func TestGetRefresh(t *testing.T) {
	cookie, err := session.SetRefreshData(uuid.NewString(), "127.0.0.1", "windows blablabla")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("cookie: %s", cookie)

	newCookie, err := session.VerifyRefreshData(cookie, "127.0.0.1", "windows blablabla")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("new cookie: %s", newCookie)
}

func TestGetRefreshFailedFalsePayload(t *testing.T) {
	cookie, err := session.SetRefreshData(uuid.NewString(), "127.0.0.1", "windows blablabla")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("cookie: %s", cookie)

	falseCookie, err := falseCompleteCookieHelper(cookie)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("false cookie: %s", falseCookie)

	_, err = session.VerifyRefreshData(falseCookie, "127.0.0.1", "windows blablabla")
	if err == nil {
		t.Fatal("should failed but got succeed")
	}

	t.Logf("got err: %s", err.Error())
}

func TestGetRefreshFailedFalseTimestamp(t *testing.T) {
	cookie, err := session.SetRefreshData(uuid.NewString(), "127.0.0.1", "windows blablabla")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("cookie: %s", cookie)

	falseCookie, err := falseTimestampHelper(cookie)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("false cookie: %s", falseCookie)

	_, err = session.VerifyRefreshData(falseCookie, "127.0.0.1", "windows blablabla")
	if err == nil {
		t.Fatal("should failed but got succeed")
	}

	t.Logf("got err: %s", err.Error())
}

func TestGetData(t *testing.T) {
	number, cookie, err := session.SetVerifyData(uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("number: %s", number)

	_, sessionID, _, err := session.GetVerifyData(cookie)
	if err != nil {
		t.Fatal(err)
	}

	if err := session.DeleteVerifySession(sessionID); err != nil {
		t.Fatal(err)
	}
}

func TestGetDataWithOldKey(t *testing.T) {
	time.Sleep(1600 * time.Millisecond)
	number, cookie, err := session.SetVerifyData(uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("number: %s", number)

	time.Sleep(1600 * time.Millisecond)
	_, _, _, err = session.GetVerifyData(cookie)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetDataFailedSessionExpired(t *testing.T) {
	number, cookie, err := session.SetVerifyData(uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("number: %s", number)

	time.Sleep(2100 * time.Millisecond)

	_, _, _, err = session.GetVerifyData(cookie)
	if err == nil {
		t.Fatal("get should fail because of expired session but got succeed")
	}

	t.Logf("got err: %s", err.Error())
}

func TestGetDataFailedFalsePayload(t *testing.T) {
	number, cokie, err := session.SetVerifyData(uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("number: %s", number)

	falseCookie, err := falsePayloadHelper(cokie)
	if err != nil {
		t.Fatal(err)
	}

	_, _, _, err = session.GetVerifyData(falseCookie)
	if err == nil {
		t.Fatal("should fail because of false payload but got succeed")
	}

	t.Logf("got err: %s", err.Error())

}

func falsePayloadHelper(cookie string) (string, error) {
	cookieBytes, err := base64.RawURLEncoding.DecodeString(cookie)
	if err != nil {
		return "", err
	}

	segments := bytes.Split(cookieBytes, []byte(" "))

	if len(segments) != 2 {
		return "", errors.New("invalid payload")
	}

	if _, err := rand.Read(segments[0]); err != nil {
		return "", err
	}

	newCookie := strings.Join([]string{string(segments[0]), string(segments[1])}, " ")
	encoded := base64.RawURLEncoding.EncodeToString([]byte(newCookie))

	return encoded, nil
}

func TestGetDataFailedFalseKey(t *testing.T) {
	number, cookie, err := session.SetVerifyData(uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("number: %s", number)

	falseCookie, err := falseTimestampHelper(cookie)
	if err != nil {
		t.Fatal(err)
	}

	_, _, _, err = session.GetVerifyData(falseCookie)
	if err == nil {
		t.Fatal("should fail because of false timestamp but got succeed")
	}

	t.Logf("got err: %s", err.Error())
}

func falseTimestampHelper(cookie string) (string, error) {
	cookieBytes, err := base64.RawURLEncoding.DecodeString(cookie)
	if err != nil {
		return "", err
	}

	log.Println(string(cookieBytes))

	segments := strings.Split(string(cookieBytes), " ")

	if len(segments) != 2 {
		return "", errors.New("invalid payload")
	}

	timestamp := time.Now().Add(1800 * time.Millisecond).UTC().Format(time.RFC3339Nano)
	newCookie := strings.Join([]string{segments[0], timestamp}, " ")
	encoded := base64.RawURLEncoding.EncodeToString([]byte(newCookie))

	return encoded, nil
}

func TestGetFalseCookie(t *testing.T) {
	number, cookie, err := session.SetVerifyData(uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("number: %s", number)

	falseCookie, err := falseCompleteCookieHelper(cookie)
	if err != nil {
		t.Fatal(err)
	}

	_, _, _, err = session.GetVerifyData(falseCookie)
	if err == nil {
		t.Fatal("should fail because of false timestamp but got succeed")
	}

	t.Logf("got err: %s", err.Error())
}

func falseCompleteCookieHelper(cookie string) (string, error) {
	cookieBytes, err := base64.RawURLEncoding.DecodeString(cookie)
	if err != nil {
		return "", err
	}

	log.Println(string(cookieBytes))

	if _, err := rand.Read(cookieBytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(cookieBytes), nil
}
