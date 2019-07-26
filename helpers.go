package sawtooth_client_sdk_go

import (
	"crypto/sha512"
	"encoding/hex"
	"os/user"
	"strings"
)

// Hexdigest computes the SHA256 hash of a byte array and returns it as a string.
func Hexdigest(str string) string {
	return HexdigestByte([]byte(str))
}

// HexdigestByte computes the SHA256 hash of a byte array and returns it as a byte array.
func HexdigestByte(data []byte) string {
	hash := sha512.New()
	hash.Write(data)
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}

// getDefaultKeyFileName finds the default Sawtooth private key and returns the patch.
func getDefaultKeyFileName() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	keyFile := currentUser.HomeDir + "/.sawtooth/keys/" + currentUser.Username + ".priv"
	return keyFile, nil
}
