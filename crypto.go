package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	pb "github.com/jonathanwei/fproxy/proto"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

func NewAEAD(key []byte) (cipher.AEAD, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	return aead, nil
}

func NewAEADOrDie(key []byte) cipher.AEAD {
	a, err := NewAEAD(key)
	if err != nil {
		glog.Fatalf("Couldn't init AEAD: %v", err)
	}
	return a
}

// Encrypt a proto using an AEAD.
func EncryptProto(aead cipher.AEAD, msg proto.Message, additionalData []byte) string {
	plaintext := MarshalProtoOrPanic(msg)

	nonce := getNonce(aead.NonceSize())

	// Encrypt in-place.
	ciphertext := plaintext
	ciphertext = aead.Seal(ciphertext[:0], nonce, plaintext, additionalData)

	outBytes := MarshalProtoOrPanic(&pb.EncryptedMessage{
		Nonce:      nonce,
		Ciphertext: ciphertext,
	})

	// Return base64'd, so that the output is ASCII-safe.
	return base64.RawURLEncoding.EncodeToString(outBytes)
}

// Decrypts a proto using an AEAD. Unmarshals the result into dst. The result
// should only be considered written if this function returns true.
func DecryptProto(aead cipher.AEAD, msg string, additionalData []byte, dst proto.Message) bool {
	msgBytes, err := base64.RawURLEncoding.DecodeString(msg)
	if err != nil {
		glog.V(2).Infof("Tried to decrypt proto with invalid base64: %v", err)
		return false
	}

	var msgProto pb.EncryptedMessage
	err = proto.Unmarshal(msgBytes, &msgProto)
	if err != nil {
		glog.V(2).Infof("Tried to decrypt proto with invalid pb.EncryptedMessage: %v", err)
		return false
	}

	// Decrypt in-place.
	plaintext := msgProto.Ciphertext
	plaintext, err = aead.Open(plaintext[:0], msgProto.Nonce, msgProto.Ciphertext, additionalData)
	if err != nil {
		glog.V(2).Infof("Failed to decrypt data: %v", err)
		return false
	}

	err = proto.Unmarshal(plaintext, dst)
	if err != nil {
		glog.V(2).Infof("Failed to decrypt proto: %v", err)
		return false
	}

	return true
}

func getNonce(size int) []byte {
	nonce := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		err := fmt.Errorf("Couldn't read from rand.Reader: %v", err)
		glog.Warning(err)
		panic(err)
	}
	return nonce
}

func MarshalProtoOrPanic(msg proto.Message) []byte {
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		err := fmt.Errorf("Couldn't marshal proto: %v; got err: %v", msg, err)
		glog.Warning(err)
		panic(err)
	}
	return msgBytes
}
