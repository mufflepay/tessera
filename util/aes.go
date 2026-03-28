package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
)

func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32) // 256-bit key
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func EncryptUserID(userID uuid.UUID, key []byte) (string, error) {

	// Convert the UUID to bytes
	userIDBytes := []byte(userID.String())

	// Generate a new AES cipher block using the encryption key
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Generate a new initialization vector (IV)
	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return "", err
	}

	// Pad the user ID to a multiple of the block size
	userIDBytes = padToBlockSize(userIDBytes, aes.BlockSize)

	// Encrypt the user ID using AES-CBC mode
	ciphertext := make([]byte, len(userIDBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, userIDBytes)

	// Concatenate the IV and ciphertext into a single string
	encrypted := base64.URLEncoding.EncodeToString(append(iv, ciphertext...))

	return encrypted, nil
}

func DecryptUserID(encrypted string, key []byte) (uuid.UUID, error) {
	// Decode the encrypted string from base64
	data, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return uuid.Nil, err
	}

	// Split the data into the IV and ciphertext
	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	// Generate a new AES cipher block using the encryption key
	block, err := aes.NewCipher(key)
	if err != nil {
		return uuid.Nil, err
	}

	// Decrypt the ciphertext using AES-CBC mode
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove any padding from the plaintext
	plaintext = removePadding(plaintext)

	// Convert the plaintext bytes to a UUID
	userID, err := uuid.FromBytes(plaintext)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func padToBlockSize(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = byte(padding)
	}
	return padded
}

func removePadding(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}

// func main() {
//     userID := uuid.New()
//     fmt.Println("Original UUID:", userID)

//     encrypted, err := EncryptUserID(userID)
//     if err != nil {
//         panic(err)
//     }
//     fmt.Println("Encrypted string:", encrypted)

//     decrypted, err := DecryptUserID(encrypted)
//     if err != nil {
//         panic(err)
//     }
//     fmt.Println("Decrypted UUID:", decrypted)
// }
