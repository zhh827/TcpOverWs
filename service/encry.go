package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func AesEncrypt(orig []byte, key string) []byte {
	k := []byte(key)
	block, _ := aes.NewCipher(k)                              // 分组秘钥
	blockSize := block.BlockSize()                            // 获取秘钥块的长度
	orig = PKCS7Padding(orig, blockSize)                      // 补全码
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize]) // 加密模式
	cryted := make([]byte, len(orig))                         // 创建数组
	blockMode.CryptBlocks(cryted, orig)                       // 加密
	return cryted
}

func AesEncryptBase64(orig string, key string) string {
	cryted := AesEncrypt([]byte(orig), key)
	return base64.StdEncoding.EncodeToString(cryted)
}

func AesDecrypt(cryted []byte, key string) []byte {

	k := []byte(key)
	block, _ := aes.NewCipher(k)                              // 分组秘钥
	blockSize := block.BlockSize()                            // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize]) // 加密模式
	orig := make([]byte, len(cryted))                         // 创建数组
	blockMode.CryptBlocks(orig, cryted)                       // 解密
	orig = PKCS7UnPadding(orig)                               // 去补全码
	return orig
}

func AesDecryptBase64(cryted string, key string) string {
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted) // 转成字节数组
	orig := AesDecrypt(crytedByte, key)
	return string(orig)
}

//补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func Md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
