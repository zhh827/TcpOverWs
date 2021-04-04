package service

import (
	"fmt"
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
