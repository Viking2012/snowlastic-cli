package main

import (
	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/snowflakedb/gosnowflake"
	"log"
	"snowlastic-cli/demo"
)

func main() {
	err := demo.IndexDemos(
		"./demo/demos.json",           //demosPath
		"./demo/demo-settings.json",   //demoSettings
		"./ignore/local_elastic.json", //credPath
		"./ignore/local_http_ca.crt",  //caCertPath
	)
	if err != nil {
		log.Fatal(err)
	}
}

func generateEsClient() (*elasticsearch.Client, error) {
	c, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"Https://localhost:9200"},
		Username:  "elastic",
		Password:  "W8jmmFwKsTeywFxiq++b",
		CACert: []byte(`-----BEGIN CERTIFICATE-----
MIIFWjCCA0KgAwIBAgIVAL3H3g3cAWtAkWXrVMmVBXswrJ2rMA0GCSqGSIb3DQEB
CwUAMDwxOjA4BgNVBAMTMUVsYXN0aWNzZWFyY2ggc2VjdXJpdHkgYXV0by1jb25m
aWd1cmF0aW9uIEhUVFAgQ0EwHhcNMjMwNzE4MTczNzEwWhcNMjYwNzE3MTczNzEw
WjA8MTowOAYDVQQDEzFFbGFzdGljc2VhcmNoIHNlY3VyaXR5IGF1dG8tY29uZmln
dXJhdGlvbiBIVFRQIENBMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA
jvpei3y4OLzl1vl1D2y8g3bWxKTx0miUhNJSC2iFqorpRCznCKX1s9LX87Rc4bw5
R2Ro8VqTqc0V4NxIsN3CdD1H/c2LTgX8aaCWg6f+bMR2oHqT+U/SicFnN2y998DG
iu9tyLqntcviZ/ho8pdStXloiFbNFCH5GDAjGLIDkgbOiRqQbRVLFE1lapHdMI7r
9FyrVWwkEjAhNJeRimOeb2l6zYDm3MISagnSmf+IbjckTuYekEPach7hlvUnLrpM
Mm+lJJXBMNeDHt0+7I5Y9MGXWMGE98tKpEH+uwW34flW+nIXRTOREnuMuEdgc0dr
8SH8/wn8XkcgFuOCuHJOBmsHOTB6xi5X5OzziLSaFu6Z6B4UgI2xi+jtCeBwDgRl
pkS6F6DUfJfGa9kdiXl/NK+iiEFZrTOqfTNU4FjWEqL6aY7Kka6b93IEJeYXiPaO
WefSOccXzcJ5cqZbuyIgscikws0lOBwnRC/U8tUV/uQozJ95mHvLF/QsmKm60U8m
rD7auROfpLy33CGbXLIZlZ2urRACGkhe9YY0+TsnQO2bmQ/4goJG/N5ouEvTle3h
uYIQmSUkG3HT4hSMrC7Fprfbteb+r8etZKsaTIqauMblb5byOgpItLoszACGUsch
dNNVRKs87uGcQfkXw9fBMn5UpdZ6Kao8/RbWCPjFwjkCAwEAAaNTMFEwHQYDVR0O
BBYEFPojP07XBjEnfzcktz11OZwoHCTuMB8GA1UdIwQYMBaAFPojP07XBjEnfzck
tz11OZwoHCTuMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggIBAIDQ
S2SaSkzSNR1ejywm8lL6jUQQEHm1mMfDxKoVsAlAxvjNoXoS9oTj1EcRAeP91Xok
aOHs45CnTWQphFQqLcJw2yZ5QCIJG87E5TY9BRbXrKQ31YyoZfWLCb6Cb6mxTSPv
5Gm759/AxRUuyAoNPxIGXeJJMh0iGJDOhzJYfMCWzpt0F6ua6eIA/7oDAI68MLS/
U99uhNU75/aVUYoIk5+FDGgTTNuZuFMJS8VvBEN4NTVCG4McCdjBSwQYUspoWuaw
6KvByIJY2JInUyZPZweHd/lr4TLVeMV2qBJuN27RS8JQL964NT6e2HizoP36CoiN
Ku+vCTtL7PrvtCb2e0UCiurKBhVqsQvdthAgKxhWfLzNFA0KH+wpw1KytdT/wsxi
XJ3xy6QwDGx5/wwLcbF7aE7f1c/Sh6L3RdYvthHA8jrckTOxdNkVtsIVyY6/GFaw
g5Jy1Sp2dwNRB54Y2GHGqkOv737SFYZ8eSShKbrrANnNs5u88aWc7WLRiaJn1RcD
WavQhJjUYKF4RSm8pTuL5ttD5Otx8rG8qdcGVwiWSUPa1/UWH7Pz6YA0noZHwzMP
TgheIASNLgAyiPipQ/HxQMPQbnZHN9LufNXGv+zjRf12HssJj5YqiwGsL6kOP5E7
zD00qqRMNWBAlEHOKx61vSIAcH6c+4/ozo5ctK+T
-----END CERTIFICATE-----
`)})
	return c, err
}

type bulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			ID     string `json:"_id"`
			Result string `json:"result"`
			Status int    `json:"status"`
			Error  struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
				Cause  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"caused_by"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
}
