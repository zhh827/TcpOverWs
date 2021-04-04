package service

import (
	"fmt"
	"log"
	"net/url"
	"testing"
)

func TestEncryDencryBase64(t *testing.T) {
	orig := "123456781234567812345678"
	key := Md5Str("123aasdaww")
	// key := "1234567812345678"

	t.Log("原文：", orig)
	encryptCode := AesEncryptBase64(orig, key)
	t.Log("密文：", encryptCode)
	decryptCode := AesDecryptBase64(encryptCode, key)
	t.Log("解密结果：", decryptCode)
	if decryptCode != orig {
		t.Fatal("解密错误")
	} else {
		t.Log("解密正确")
	}

	encryptByte := AesEncrypt([]byte(orig), key)
	fmt.Println(len(encryptByte))
}

func TestMainbase64(t *testing.T) {
	s := "eG2UoqPH7ays6BgJxUAW7PdYMjHk3qhcCScZkKHzGP0A%2BhRymnH0Tk1sSCQ3Ch8D"
	decodePara, err := url.QueryUnescape(s)
	if err != nil {
		log.Fatalln(s)
	}
	t.Log(decodePara)
	log.Println(decodePara)
}
