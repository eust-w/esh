package yaml

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"esh/utils"
	"io"
	mrand "math/rand"
	"strings"
	"time"
)
//You can change this aes keys,
//there is no limit to the keys number, but it must be  aes key!
var keys = [...][]byte{
	[]byte("cKIap6a5Ojp3NL8uaAz4pjKlbKQGcR4o"),
	[]byte("J4XH2HKfSIraXGRzxdhmO5d3BDEyguQ3"),
	[]byte("4IfEboP2qSe0g5fSc6LSdRk26sQnZg4L"),
	[]byte("doj3Uy0raqbwEDiJRAtU7pOXQ0xSMjpn"),
	[]byte("Hav4/muOeWfxfw9PgtY8j+muSTYvxjNB"),
	[]byte("KBLWopeSXRhFxQGU26pHh3jsY/cFrqSm"),
}

func AesEncrypt(data string) string {
	var encrypted []byte
	origData := []byte("\t\t" + data + "\t\t")
	r := mrand.New(mrand.NewSource(time.Now().Unix()))
	index := r.Intn(len(keys)-1)
	key := utils.ByteXor(keys[index],keys[index+1])
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return hex.EncodeToString(encrypted)
}

func AesDecrypt(encrypted string) ([]string, error){
	for i, k := range keys[:len(keys)-1] {
		key := utils.ByteXor(k,keys[i+1])
		decrypted, _ := hex.DecodeString(encrypted)
		block, _ := aes.NewCipher(key)
		if len(decrypted) < aes.BlockSize {
			return  []string{}, errors.New("can't decrypt")
		}
		iv := decrypted[:aes.BlockSize]
		decrypted = decrypted[aes.BlockSize:]

		stream := cipher.NewCFBDecrypter(block, iv)
		stream.XORKeyStream(decrypted, decrypted)
		stage1 := string(decrypted)
		if stage1[0:2] == ("\t\t") && stage1[len(stage1)-2:] == "\t\t" {
			stage2 := strings.Split(stage1[2:len(stage1)-2], "\n")
			if len(stage2) == 4 {
				return stage2, nil
			}
		}
	}
	return []string{} , errors.New("can't decrypt")
}
