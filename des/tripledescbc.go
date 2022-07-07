package des

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"encoding/hex"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/wumansgy/goEncrypt"
)

/**
	Triple des encryption and decryption
      algorithm : Encryption: key one encryption -> key two decryption -> key three encryption
                  Decryption: key three decryption -> key two encryption -> key one decryption
*/
func TripleDesEncrypt(plainText, key, ivDes []byte) ([]byte, error) {
	if len(key) != 24 {
		return nil, goEncrypt.ErrKeyLengthTwentyFour
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	paddingText := goEncrypt.PKCS5Padding(plainText, block.BlockSize())

	var iv []byte
	if len(ivDes) != 0 {
		if len(ivDes) != block.BlockSize() {
			return nil, goEncrypt.ErrIvDes
		} else {
			iv = ivDes
		}
	} else {
		iv = []byte(goEncrypt.Ivdes)
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)

	cipherText := make([]byte, len(paddingText))
	blockMode.CryptBlocks(cipherText, paddingText)
	return cipherText, nil
}

func TripleDesDecrypt(cipherText, key, ivDes []byte) ([]byte, error) {
	if len(key) != 24 {
		return nil, goEncrypt.ErrKeyLengthTwentyFour
	}
	// 1. Specifies that the 3des decryption algorithm creates and returns a cipher.Block interface using the TDEA algorithm。
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	// 2. Delete the filling
	// Before deleting, prevent the user from entering different keys twice and causing panic, so do an error handling
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Errorf("runtime err=%v,Check that the key or text is correct", err)
			default:
				log.Errorf("error=%v,check the cipherText ", err)
			}
		}
	}()

	var iv []byte
	if len(ivDes) != 0 {
		if len(ivDes) != block.BlockSize() {
			return nil, goEncrypt.ErrIvDes
		} else {
			iv = ivDes
		}
	} else {
		iv = []byte(goEncrypt.Ivdes)
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)

	paddingText := make([]byte, len(cipherText)) //
	blockMode.CryptBlocks(paddingText, cipherText)

	plainText, err := goEncrypt.PKCS5UnPadding(paddingText, block.BlockSize())
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

func TripleDesEncryptBase64(plainText, key, ivAes []byte) (string, error) {
	encryBytes, err := TripleDesEncrypt(plainText, key, ivAes)
	return base64.StdEncoding.EncodeToString(encryBytes), err
}

func TripleDesEncryptHex(plainText, key, ivAes []byte) (string, error) {
	encryBytes, err := TripleDesEncrypt(plainText, key, ivAes)
	return hex.EncodeToString(encryBytes), err
}

func TripleDesDecryptByBase64(cipherTextBase64 string, key, ivAes []byte) ([]byte, error) {
	plainTextBytes, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return []byte{}, err
	}
	return TripleDesDecrypt(plainTextBytes, key, ivAes)
}

func TripleDesDecryptByHex(cipherTextHex string, key, ivAes []byte) ([]byte, error) {
	plainTextBytes, err := hex.DecodeString(cipherTextHex)
	if err != nil {
		return []byte{}, err
	}
	return TripleDesDecrypt(plainTextBytes, key, ivAes)
}