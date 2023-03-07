package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	rnd "math/rand"
	"time"
)

// To Encrypt String: encryptedString = runEncrypt(cipherKey, toEncrypt)
// To Decrypt String: runDecrypt(cipherKey, encryptedString)
// TODO: Retrieve Encryption key based on Client Name from KeyMap

func encrypt(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return
}

func decrypt(key []byte, securemess string) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess = string(cipherText)
	return
}

func RunEncrypt(cipherKey []byte, msg string) string {
	if encrypted, err := encrypt(cipherKey, msg); err != nil {
		log.Println(err)
		return err.Error()
	} else {
		//log.Printf("ENCRYPTED: %s\n", encrypted)
		return encrypted
	}
}

func RunDecrypt(cipherKey []byte, msg string) string {
	if decrypted, err := decrypt(cipherKey, msg); err != nil {
		log.Println(err)
		return err.Error()
	} else {
		//log.Printf("DECRYPTED: %s\n", decrypted)
		return decrypted
	}
}

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")

func shuffle() bool {
	inRune := chars
	rnd.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	fmt.Println(string(inRune))
	chars = inRune
	return true
}

func generateIV(n int) string {
	log.Println("Generating new Seed")
	time.Sleep(2 * time.Second)
	rnd.Seed(time.Now().Unix())
	log.Println("Suffling Characters")
	time.Sleep(2 * time.Second)
	shuffle()
	str := make([]rune, n)
	for i := range str {
		str[i] = chars[rnd.Intn(len(chars))]
	}
	return string(str)
}

func runTest(toEncrypt string) {
	cipherKey := []byte("")

	key := generateIV(32)
	//exec.Command("bash", "-c", "echo \""+key+"\" > encryption.key").Start()
	log.Printf("Key: %s", key)
	time.Sleep(5 * time.Second)
	cipherKey = []byte(key)
	//enc := *stringPtr

	var encryptedString string
	//if *encryptPtr {
	encryptedString = RunEncrypt(cipherKey, toEncrypt)
	//}

	//if *decryptPtr {
	RunDecrypt(cipherKey, encryptedString)
	//}
}
