// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cert

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

// TestCert returns a valid self-signed certificate for testing.
func TestCert(t testing.T) *Cert {
	require := require.New(t)

	// Create directory for our certs
	td, err := ioutil.TempDir("", "cert")
	require.NoError(err)
	t.Cleanup(func() { os.RemoveAll(td) })

	// Our test cert paths
	crtPath := filepath.Join(td, "tls.crt")
	keyPath := filepath.Join(td, "tls.key")

	// Write our certs
	require.NoError(ioutil.WriteFile(crtPath, []byte(strings.TrimSpace(testCRT)), 0600))
	require.NoError(ioutil.WriteFile(keyPath, []byte(strings.TrimSpace(testKey)), 0600))

	// Create it
	c, err := New(nil, crtPath, keyPath)
	require.NoError(err)
	t.Cleanup(func() { c.Close() })

	return c
}

const testCRT = `
-----BEGIN CERTIFICATE-----
MIIC6zCCAdOgAwIBAgIRAJc/nbWADh0bea5JWJ30olkwDQYJKoZIhvcNAQELBQAw
ADAeFw0yMTA5MjEwMjM2MTFaFw0yMTEyMjAwMjM2MTFaMAAwggEiMA0GCSqGSIb3
DQEBAQUAA4IBDwAwggEKAoIBAQDNaJuRwZoA1Ci92ZxL44tikehZWhYgGWWkHOle
kGVgmj+o6oe5GpVwd5iaCeRlaK2BlqbGZ+IQP5DrBjrK/KxMmCM5n3UoXp3yIw39
WSSUyLen2fuH9qGSGA9+3pXArOjFsVzma1oeiMBxyhoepnzpxmzGZJNNkrWzHLCd
JjjCyGZ++o7+DxPKs0qP7Bh7QpkwoflH6z9oSYF6XhKmdN0Vxism0TQhkWyao8uN
tLD6fVS5JAVxp64aH07+Ev+6H8kHP+x1ERJloU1KpZlUdBZy53R8HgPb5icO6Mg8
lneC5KhYtMOmlEQA/V3QNMWjzl2WNbGpDhajWXiksRd98ClrAgMBAAGjYDBeMB0G
A1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMBwGA1Ud
EQEB/wQSMBCCDnNtdHAucmZjODIyLm14MBEGA1UdHwQKMAgwBqAEoAKGADANBgkq
hkiG9w0BAQsFAAOCAQEAMEZyF5Eh51MaUyLKgBeIFgRs7h58AOEYzTEvXucX4Tnu
UEN2pLZ348A6JFzQGNU+82dM82J01y8YnGFbDvhk71erXLRTeGtjXqsGEmTiEke9
22JHVIdkTyZchCt1oRI12R63C+QUBdUDcWNqG71+GTTnB2Jwn9r18OzXKw62M6fI
eTZsY28lYjei4jC0BMdwyuQNVEVpzNlGCVmhnDRdw3/O0GltQgFiEaJED/PxFNN7
wOJOiyE05+x92ZYZezW5LHs6iFFxgKpUssZrg/UMhyrZF6ETv5zuvEr5dA1s9w6P
J8X2zYv6JiF0gUp8Unn4/gvZSLty89GrPskfuha3UQ==
-----END CERTIFICATE-----
`

const testKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAzWibkcGaANQovdmcS+OLYpHoWVoWIBllpBzpXpBlYJo/qOqH
uRqVcHeYmgnkZWitgZamxmfiED+Q6wY6yvysTJgjOZ91KF6d8iMN/VkklMi3p9n7
h/ahkhgPft6VwKzoxbFc5mtaHojAccoaHqZ86cZsxmSTTZK1sxywnSY4wshmfvqO
/g8TyrNKj+wYe0KZMKH5R+s/aEmBel4SpnTdFcYrJtE0IZFsmqPLjbSw+n1UuSQF
caeuGh9O/hL/uh/JBz/sdRESZaFNSqWZVHQWcud0fB4D2+YnDujIPJZ3guSoWLTD
ppREAP1d0DTFo85dljWxqQ4Wo1l4pLEXffApawIDAQABAoIBABxZZKcwNiYJIfpZ
z0V2CSW4h33Vfho+BxYoW1kOyr8TatfQTp0QezGDgA86cEhCszegaDIX4k5lx4V4
XaLoTotfr/Ti+hNxQ5FIn0SpCfBy504GOr3gHxp/sZvL8pUeCB5IxU6T4GM4cn8B
6qryRWkqVpbnCKF7LxKiUgnOXaUcLiRsav/lnocSOveRN74msh2VAbEG6Sj5Gf/J
0pCpNBPseEnb2XOVdAk/sIEEJyGwL3NCSPu2OxFzndf+eLKy3hj74g6K+CB8DA8R
HO0pdzrCGJ+wuHKrHgJr8Grv1sJ5pzustJNmJ+d2idbzRAHTyPQIvrdJMAZ+xdQT
jTIAwFECgYEA50Ftu+jI6MLBAQPCz1V4ClqdI2c2jKS8dbzgX+ShBMNwB2V+tv+C
HMtJx0Y5chKswoBG8xdhWHBWQAMqVSYiU8RDuw8ow+yYd9K6K45lYvxlCu3a2mPL
10KuUxPd/b9xACRBunI1lC8oe0iZ17iBB175HxUrv7apnqK75suvaRMCgYEA42Mt
WMb4fipeoTxXfuJq+K5c70r3X+yq7z4FCGqavLHCKnthvDx/WHNnNn84OF8+fu9s
yYBAT/R8WWZH/9JajVoGNgcXhnFxxAfVzjjmTTn10iQd2tUlldxOTXQcwvIrEXuF
Ds47cOm8KvisI4WzUYXJKQkBH9PmNl1H07vZYUkCgYEAt9PK5xSkoGIwCh5zPV0z
jwd44iupsSNCrFT4B0I2vRKee0Ky98UkKL9wZnfsMkGmEvblb1emiibCaSAbBpTJ
tMzPCmOChDwO9zELzJPlBEqeB5IL5o5t+y+GY4Pnc047BWHM3ejrrl/OTjHoGRMR
fkqAvbSWkk8hhnjV5SEEzwECgYBKHh+/2ktRRJpH0BVtBHx9xsgAL91mZQxqozqc
vbLmYsK5ejInW0jfGe7AssMujM0gLwa0v5s29Kg7s70wQ+7EOF3h6nnelsfQcAVf
DOj0rznTX3ZjyCpSKNdVI83kNW+YaTy70LlLWsS89QWXJpOGtScWuxqktztI6Srq
d0aqEQKBgQC+jvBpE2+j9/hf6ezICetbzCR1t9E3TUcpidv6Vh+HbEFDnW4DUBZD
bN9oEPw+XzTXJDrF/Thmzvy4HrLbK76LDvZBxG2kC0Qy0yjkc6SwrLCF/iTJgSJb
wfoAJCaax+iiPO4EsOKBercx/2jzEXvg3IPhLpFVcUHdzfB3c+6zKw==
-----END RSA PRIVATE KEY-----
`
