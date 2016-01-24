package main

import (
	"crypto/cipher"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	pb "github.com/jonathanwei/fproxy/proto"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type cookieCrypter struct {
	aead cipher.AEAD
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

		Secure:   true, // TODO: get this from config.
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

func WithAuthCookie(ctx context.Context, cookie *pb.AuthCookie) context.Context {
	return metadata.NewContext(ctx, metadata.Pairs("auth", proto.MarshalTextString(cookie)))
}

func GetAuthCookie(ctx context.Context) *pb.AuthCookie {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil
	}

	cookies, ok := md["auth"]
	if !ok || len(cookies) != 1 {
		return nil
	}

	var authCookie pb.AuthCookie
	err := proto.UnmarshalText(cookies[0], &authCookie)
	if err != nil {
		return nil
	}

	return &authCookie
}
