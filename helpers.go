package sawtooth_client_sdk_go

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"io/ioutil"
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

// readPrivateKeyFromFile reads a private key from a file and returns it as a PrivateKey object.
func readPrivateKeyFromFile(path string) (signing.PrivateKey, error) {
	// Read the key
	keyData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Could not read private key from file (%s) with error: %s", path, err)
	}

	// Parse the key into a PrivateKey object.
	keyData = []byte(strings.TrimSpace(string(keyData)))
	privateKey := signing.NewSecp256k1PrivateKey(keyData)

	return privateKey, nil
}
