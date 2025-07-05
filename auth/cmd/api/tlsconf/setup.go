package tlsconf

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"os"
)

type TLSBuilder struct {
	certPath   string
	keyPath    string
	rootCAPath string
}

type RawConfig interface {
	CACertPath() string
	CertPath() string
	KeyPath() string
}

type tlsConfigMap map[string]string

func (confMap tlsConfigMap) CACertPath() string {
	return confMap["CERTIFICATION"]
}

func (confMap tlsConfigMap) CertPath() string {
	return confMap["CERT"]
}

func (confMap tlsConfigMap) KeyPath() string {
	return confMap["KEY"]
}

func NewTLSBuilder() *TLSBuilder {
	return &TLSBuilder{}
}

func (builder *TLSBuilder) CertPath(certPath string) *TLSBuilder {
	builder.certPath = certPath
	return builder
}

func (builder *TLSBuilder) KeyPath(keyPath string) *TLSBuilder {
	builder.keyPath = keyPath
	return builder
}

func (builder *TLSBuilder) CACertPath(certPath string) *TLSBuilder {
	builder.rootCAPath = certPath
	return builder
}

func (builder *TLSBuilder) BuildConfigMap() (RawConfig, error) {
	if builder.certPath == "" || builder.keyPath == "" || builder.rootCAPath == "" {
		return nil, errors.New("uncomplete setup")
	}

	configMap := make(map[string]string)
	configMap["CERT"] = builder.certPath
	configMap["KEY"] = builder.keyPath
	configMap["CERTIFICATION"] = builder.rootCAPath

	return tlsConfigMap(configMap), nil
}

func (builder *TLSBuilder) Build() (*tls.Config, error) {
	if builder.certPath == "" || builder.keyPath == "" || builder.rootCAPath == "" {
		return nil, errors.New("incomplete tls setup")
	}

	caCert, err := os.ReadFile(builder.rootCAPath)
	if err != nil {
		log.Printf("line 32: %v\n", err)
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(builder.certPath, builder.keyPath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}, nil
}
