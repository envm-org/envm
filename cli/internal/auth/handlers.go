package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	credentialsFile = "credentials"
	envmDir         = ".envm"
	salt            = "envm-cli-salt-v1"
)

type Credentials struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type LoginResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func getEncryptionKey() ([]byte, error) {
	machineID, err := getMachineID()
	if err != nil {
		home, _ := os.UserHomeDir()
		machineID = "fallback-machine-id-" + home
	}

	hash := sha256.Sum256([]byte(machineID + salt))
	return hash[:], nil
}

func getMachineID() (string, error) {
	switch runtime.GOOS {
	case "linux":
		content, err := os.ReadFile("/etc/machine-id")
		if err == nil {
			return strings.TrimSpace(string(content)), nil
		}
		content, err = os.ReadFile("/var/lib/dbus/machine-id")
		if err == nil {
			return strings.TrimSpace(string(content)), nil
		}
		return "", errors.New("machine-id not found")

	case "darwin": // macOS
		cmd := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, "IOPlatformUUID") {
				parts := strings.Split(line, "=")
				if len(parts) == 2 {
					return strings.Trim(strings.TrimSpace(parts[1]), "\""), nil
				}
			}
		}
		return "", errors.New("IOPlatformUUID not found")

	case "windows":
		cmd := exec.Command("cmd", "/C", "wmic csproduct get UUID")
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		// wmic output format is:
		// UUID
		// XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && line != "UUID" {
				return line, nil
			}
		}
		return "", errors.New("UUID not found from wmic")

	default:
		return "", fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}
}

func getCredentialsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, envmDir, credentialsFile), nil
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func SaveCredentials(creds Credentials) error {
	path, err := getCredentialsPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	key, err := getEncryptionKey()
	if err != nil {
		return err
	}

	encrypted, err := encrypt(data, key)
	if err != nil {
		return err
	}

	hexStr := hex.EncodeToString(encrypted)

	return os.WriteFile(path, []byte(hexStr), 0600)
}
func LoadCredentials() (*Credentials, error) {
	path, err := getCredentialsPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	hexBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	encrypted, err := hex.DecodeString(string(hexBytes))
	if err != nil {
		return nil, fmt.Errorf("corrupted credentials file")
	}

	key, err := getEncryptionKey()
	if err != nil {
		return nil, err
	}

	decrypted, err := decrypt(encrypted, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt credentials (machine ID mismatch?): %v", err)
	}

	var creds Credentials
	if err := json.Unmarshal(decrypted, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

func Logout() error {
	path, err := getCredentialsPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
