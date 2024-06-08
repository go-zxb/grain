package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

// AES加密函数，使用CBC模式和PKCS7填充
func AesEncrypt(plainText []byte, key []byte, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	plainText = PKCS7Padding(plainText, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, iv)

	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AES解密函数，使用CBC模式和PKCS7填充
func AesDecrypt(cipherText string, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherData, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)

	plainText := make([]byte, len(cipherData))
	blockMode.CryptBlocks(plainText, cipherData)

	plainText = PKCS7UnPadding(plainText)

	return plainText, nil
}

// 使用PKCS7进行填充，返回填充后的数据
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 去除PKCS7填充，返回去除填充后的数据
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func GenerateAesIV(blockSize int) string {
	iv := make([]byte, blockSize)
	_, _ = rand.Read(iv)
	return hex.EncodeToString(iv)
}
