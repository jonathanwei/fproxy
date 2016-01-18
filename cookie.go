package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	pb "github.com/jonathanwei/fproxy/proto"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type cookieCrypter struct {
	aead cipher.AEAD
}

func newCookieCrypter(key []byte) cookieCrypter {
	c, err := aes.NewCipher(key)
	if err != nil {
		glog.Fatalf("Cannot init cookie crypter: %v", err)
	}

	aead, err := cipher.NewGCM(c)
	if err != nil {
		glog.Fatalf("Cannot init cookie crypter: %v", err)
	}

	return cookieCrypter{aead}
}

var (
	authAdditionalData = []byte("fproxy/auth")
	authCookieName     = "fa"
)

func (c cookieCrypter) SetAuthCookie(rw http.ResponseWriter, a *pb.AuthCookie) {
	plaintext, err := proto.Marshal(a)
	if err != nil {
		panic(fmt.Sprintf("Couldn't marshal proto: %v", a))
	}

	nonce := make([]byte, c.aead.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		glog.Warningf("Couldn't read from /dev/urandom: %v", err)
		panic(err)
	}

	ciphertext := c.aead.Seal(nil, nonce, plaintext, authAdditionalData)
	cookieBytes, err := proto.Marshal(&pb.EncodedAuthCookie{
		Nonce:      nonce,
		Ciphertext: ciphertext,
	})
	if err != nil {
		panic(fmt.Sprintf("Couldn't marshal proto: %v", err))
	}

	cookieStr := base64.RawStdEncoding.EncodeToString(cookieBytes)
	cookie := http.Cookie{
		Name:  authCookieName,
		Value: cookieStr,

		Secure:   false, // TODO: get this from config.
		HttpOnly: true,
	}
	http.SetCookie(rw, &cookie)
}

func (c cookieCrypter) GetAuthCookie(req *http.Request) *pb.AuthCookie {
	cookie, err := req.Cookie(authCookieName)
	if err != nil {
		return nil
	}

	cookieBytes, err := base64.RawStdEncoding.DecodeString(cookie.Value)
	if err != nil {
		glog.Warningf("Got auth cookie with invalid base64: %q", cookie.Value)
		return nil
	}

	var encodedCookie pb.EncodedAuthCookie
	err = proto.Unmarshal(cookieBytes, &encodedCookie)
	if err != nil {
		glog.Warningf("Got auth cookie with unmarshallable proto: %q", cookieBytes)
		return nil
	}

	plaintext, err := c.aead.Open(nil, encodedCookie.Nonce, encodedCookie.Ciphertext, authAdditionalData)
	if err != nil {
		glog.Warningf("Got auth cookie with invalid authenticator: %v", err)
		return nil
	}

	var authCookie pb.AuthCookie
	err = proto.Unmarshal(plaintext, &authCookie)
	if err != nil {
		glog.Errorf("Couldn't unmarshal auth cookie that was authenticated: %v", err)
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
