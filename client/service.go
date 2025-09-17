package client

import (
	"bufio"
	"time"
)

type V1Data struct {
	Client      string    `json:"client"`
	Secret      string    `json:"secret"`
	Token       string    `json:"token"`
	TokenExpiry time.Time `json:"tokenExpiry"`
}

type V2Data struct {
	ClientToken   string `json:"clientToken"`
	Authorization string `json:"authorization"`
}

type Config struct {
	Scanner *bufio.Scanner `json:"-"`
	V1      V1Data         `json:"v1"`
	V2      V2Data         `json:"v2"`
}
