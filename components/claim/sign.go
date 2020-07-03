package claim

import (
	"crypto/rsa"
	"encoding/asn1"
	"encoding/pem"
	"io/ioutil"
)

func ParseRsaPublicKeyFromFile(path string) (*rsa.PublicKey, error) {
	content, err := ioutil.ReadFile(path)
	if nil != err {
		return nil, err
	}

	block, _ := pem.Decode(content)
	key := &rsa.PublicKey{}
	_, err = asn1.Unmarshal(block.Bytes, key)
	if nil != err {
		return nil, err
	}

	return key, nil
}
