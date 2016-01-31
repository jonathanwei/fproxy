package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/golang/glog"
	pb "github.com/jonathanwei/fproxy/proto"
)

func TLSConfig(config *pb.TLSConfig) (*tls.Config, error) {
	var certs []tls.Certificate
	for _, certConfig := range config.Cert {
		cert, err := tls.LoadX509KeyPair(certConfig.CertFile, certConfig.KeyFile)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}

	caCertPool := x509.NewCertPool()
	for _, caConfig := range config.Ca {
		pemBytes, err := ioutil.ReadFile(caConfig.CaFile)
		if err != nil {
			return nil, err
		}
		ok := caCertPool.AppendCertsFromPEM(pemBytes)
		if !ok {
			return nil, fmt.Errorf("Couldn't add CA cert to cert pool: %v", pemBytes)
		}
	}

	c := &tls.Config{
		Certificates: certs,
		RootCAs:      caCertPool,
		ClientCAs:    caCertPool,
		ServerName:   config.ServerName,
		MinVersion:   tls.VersionTLS12,
	}
	c.BuildNameToCertificate()

	return c, nil
}

func TLSConfigOrDie(config *pb.TLSConfig) *tls.Config {
	t, err := TLSConfig(config)
	if err != nil {
		glog.Fatalf("Error making TLS config: %v", err)
	}
	return t
}

func BackendTLSConfigOrDie(config *pb.TLSConfig) *tls.Config {
	t := TLSConfigOrDie(config)
	t.ClientAuth = tls.RequireAndVerifyClientCert
	return t
}

func FrontendTLSConfigOrDie(config *pb.TLSConfig) *tls.Config {
	return TLSConfigOrDie(config)
}

func FrontendClientTLSConfigOrDie(config *pb.TLSConfig) *tls.Config {
	return TLSConfigOrDie(config)
}
