package utils

import "crypto/rsa"

type Image struct {
	URI string `json:"uri"`
}

type Key struct {
	Key rsa.PublicKey `json:"key"`
}
