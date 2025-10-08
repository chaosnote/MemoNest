package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/pbkdf2"
)

// Base64Encode ...
func Base64Encode(plaintext string) string {
	return base64.StdEncoding.EncodeToString([]byte(plaintext))
}

// Base64Decode ...
func Base64Decode(ciphertext string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(ciphertext)
}

// HexDecode 16 進位轉為 string
func HexEncode(plaintext []byte) string {
	return hex.EncodeToString(plaintext)
}

// HexDecode 16 進位反轉為 []byte
func HexDecode(ciphertext string) ([]byte, error) {
	return hex.DecodeString(ciphertext)
}

//-----------------------------------------------

// MD5Hash ...
func MD5Hash(plaintext string) string {
	h := md5.New()
	io.WriteString(h, plaintext)

	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA256Hash ...
func SHA256Hash(plaintext string) string {
	h := sha256.New()
	h.Write([]byte(plaintext))
	// hex.EncodeToString(...)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//-----------------------------------------------

// aes_pkcs7padding 進行 PKCS#7 填充
func aes_pkcs7padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(ciphertext, padtext...)
}

// aes_pkcs7unpadding 移除 PKCS#7 填充
func aes_pkcs7unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if unpadding > length {
		return origData
	}
	return origData[:(length - unpadding)]
}

// GenAESKey 根據密碼和鹽生成金鑰
func GenAESKey(passphrase []byte) []byte {
	// 客制化、不另帶參數
	salt := []byte("")
	iterations := 1000
	key_len := 32 // 256 bits for AES-256
	return pbkdf2.Key(passphrase, salt, iterations, key_len, sha256.New)
}

// AesEncrypt 加密，使用固定的 IV
func AesEncrypt(plainText, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	plainText = aes_pkcs7padding(plainText, blockSize)

	// 固定 IV，這裡假設為全零。(需同時改動解密)
	fixedIV := make([]byte, blockSize)

	mode := cipher.NewCBCEncrypter(block, fixedIV)
	cipherText := make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, plainText)

	return hex.EncodeToString(cipherText), nil
}

// AesDecrypt 解密，使用固定的 IV
func AesDecrypt(cipherHexText string, key []byte) (string, error) {
	cipherText, err := hex.DecodeString(cipherHexText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if len(cipherText)%blockSize != 0 {
		return "", fmt.Errorf("密文長度不是區塊大小的倍數")
	}

	// 固定 IV，與加密時保持一致
	fixedIV := make([]byte, blockSize)

	mode := cipher.NewCBCDecrypter(block, fixedIV)
	plainText := make([]byte, len(cipherText))
	mode.CryptBlocks(plainText, cipherText)

	plainText = aes_pkcs7unpadding(plainText)

	return string(plainText), nil
}

// -----------------------------------------------

// RSAEncrypt
//
// WARN !! return Hex
//
// call RSAInit(file_path) before using
func RSAEncrypt(plaintext []byte) (ciphertext string, e error) {
	var ref []byte
	ref, e = rsa.EncryptPKCS1v15(rand.Reader, &rsa_key.PublicKey, plaintext) // 加密明文信息
	if e != nil {
		return
	}
	ciphertext = fmt.Sprintf("%x", ref)
	return
}

// RSADecrypt
//
// call RSAInit(file_path) before using
//
// 請留意是否為 Hex !! 需在使用 utils.HexDecode 反解
func RSADecrypt(ciphertext string) (plaintext []byte, e error) {
	plaintext, e = rsa.DecryptPKCS1v15(rand.Reader, rsa_key, []byte(ciphertext)) // (私)解密密文信息
	return
}

var rsa_key *rsa.PrivateKey

// RSAInit 使用前需呼叫
//
// 加密長度 {"max_limite", max/8-11}
//
// replace bool 是否取代本地檔案
func RSAInit(file_path string, max int, replace bool) {
	// 取代 且 檔案存在
	if !replace && FileExist(file_path) {
		rsa_key = read_rsa_key(file_path)
		return
	}

	var e error
	rsa_key, e = rsa.GenerateKey(rand.Reader, max)
	if e != nil {
		panic(e)
	}
	save_rsa_key(rsa_key, file_path)
}

func save_rsa_key(privateKey *rsa.PrivateKey, filename string) error {
	dir, _ := filepath.Split(filename)
	e := os.MkdirAll(dir, os.ModePerm)
	if e != nil {
		return e
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey) // 格式化私鑰
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	file, e := os.Create(filename) // 寫入文件
	if e != nil {
		return e
	}
	defer file.Close()

	e = pem.Encode(file, pemBlock)
	if e != nil {
		return e
	}
	return nil
}

func read_rsa_key(filename string) *rsa.PrivateKey {
	privateKeyPEM, e := os.ReadFile(filename) // 讀取 pem 檔案
	if e != nil {
		panic(e)
	}
	block, _ := pem.Decode(privateKeyPEM)                   // pem 轉私鑰解碼
	privateKey, e := x509.ParsePKCS1PrivateKey(block.Bytes) // 解碼數據轉為 rsa 私鑰
	if e != nil {
		panic(e)
	}
	return privateKey
}
