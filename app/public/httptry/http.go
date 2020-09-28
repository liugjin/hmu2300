package httptry

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"

	"clc.hmu/app/public/store/etc"
)

var (
	HttpsClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: loadCaRoot(),
			},
		},
	}
)

func loadCaRoot() *x509.CertPool {
	roots := x509.NewCertPool()
	data, err := ioutil.ReadFile(os.ExpandEnv(etc.Etc.String("public", "ca-cert")))
	if err != nil {
		panic(err)
	}
	if !roots.AppendCertsFromPEM(data) {
		panic("load ca files failed")
	}
	return roots
}
