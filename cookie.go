package main

import (
	"crypto/cipher"
	"net/http"
	"time"

	pb "github.com/jonathanwei/fproxy/proto"

	"github.com/golang/glog"
)

type cookieCrypter struct {
	aead     cipher.AEAD
	insecure bool
}

var (
	authAdditionalData = []byte("fproxy/auth")
	authCookieName     = "fa"
)

func (c cookieCrypter) SetAuthCookie(rw http.ResponseWriter, a *pb.AuthCookie) {
	cookieVal := EncryptProto(c.aead, a, authAdditionalData)

	cookie := http.Cookie{
		Name:  authCookieName,
		Value: cookieVal,

		Expires: time.Now().Add(365 * 24 * time.Hour),

		Secure:   !c.insecure,
		HttpOnly: true,
	}
	http.SetCookie(rw, &cookie)
}

func (c cookieCrypter) GetAuthCookie(req *http.Request) *pb.AuthCookie {
	cookie, err := req.Cookie(authCookieName)
	if err != nil {
		return nil
	}

	var authCookie pb.AuthCookie
	ok := DecryptProto(c.aead, cookie.Value, authAdditionalData, &authCookie)
	if !ok {
		glog.Errorf("Couldn't decrypt auth cookie; pretending it didn't exist.")
		return nil
	}

	return &authCookie
}
