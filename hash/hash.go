package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"s/utils"
)

type Result struct {
	Input  string
	MD5    string
	SHA1   string
	SHA256 string
	SHA512 string
}

func Do(input string) *Result {
	md5h := md5.Sum([]byte(input))
	sha1h := sha1.Sum([]byte(input))
	sha256h := sha256.Sum256([]byte(input))
	sha512h := sha512.Sum512([]byte(input))

	return &Result{
		Input:  input,
		MD5:    hex.EncodeToString(md5h[:]),
		SHA1:   hex.EncodeToString(sha1h[:]),
		SHA256: hex.EncodeToString(sha256h[:]),
		SHA512: hex.EncodeToString(sha512h[:]),
	}
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ HASH GENERATOR ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Input: ") + utils.Wht(r.Input))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("MD5: ") + utils.Wht(r.MD5))
	fmt.Println(utils.Gry("SHA1: ") + utils.Wht(r.SHA1))
	fmt.Println(utils.Gry("SHA256: ") + utils.Wht(r.SHA256))
	fmt.Println(utils.Gry("SHA512: ") + utils.Wht(r.SHA512))
}
