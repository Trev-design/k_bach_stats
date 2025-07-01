package tlsconf

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"os"
	"sync"
)

type TLSBuilder struct {
	certPath string
	keyPath  string
}

type RawConfig interface {
	CACertPath() string
	CertPath() string
	KeyPath() string
}

type tlsConfigMap map[string]string

var caCertPath string
var caPool *x509.CertPool
var once sync.Once

func GenerateCertPool(certPath string) (err error) {
	once.Do(func() {
		cert, err := os.ReadFile(caCertPath)
		if err != nil {
			log.Println(err)
			return
		}

		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(cert)
		caCertPath = certPath
		caPool = certPool
	})

	return
}

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

func (builder *TLSBuilder) BuildConfigMap() (RawConfig, error) {
	if builder.certPath == "" || builder.keyPath == "" {
		return nil, errors.New("uncomplete setup")
	}

	configMap := make(map[string]string)
	configMap["CERT"] = builder.certPath
	configMap["KEY"] = builder.keyPath
	configMap["CERTIFICATION"] = caCertPath

	return tlsConfigMap(configMap), nil
}

func (builder *TLSBuilder) Build() (*tls.Config, error) {
	if builder.certPath == "" || builder.keyPath == "" {
		return nil, errors.New("incomplete tls setup")
	}

	cert, err := tls.LoadX509KeyPair(builder.certPath, builder.keyPath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
	}, nil
}
