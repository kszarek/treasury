package client

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/AirHelp/treasury/utils"
)

// Export returns secrets in given format
// format should be provided in singleKeyExportFormat
// e.g.: singleKeyExportFormat = "export %s='%s'\n"
func (c *Client) Export(key, singleKeyExportFormat string) (string, error) {
	var secrets []*Secret
	var err error
	// if we get valid prefix we use ReadGroup method
	// else we have 1 key only and we use Read method
	if validPrefix(key) {
		secrets, err = c.ReadGroup(key)
		if err != nil {
			return "", err
		}
	} else {
		secret, err := c.Read(key)
		if err != nil {
			return "", err
		}
		secrets = append(secrets, secret)
	}
	var sortedKeys []string
	keySecretMap := make(map[string]*Secret, len(secrets))
	for _, secret := range secrets {
		sortedKeys = append(sortedKeys, secret.Key)
		keySecretMap[secret.Key] = secret
	}
	sort.Strings(sortedKeys)
	var buffer bytes.Buffer
	for _, key := range sortedKeys {
		secret := keySecretMap[key]
		buffer.WriteString(fmt.Sprintf(singleKeyExportFormat, filepath.Base(secret.Key), secret.Value))
	}
	return buffer.String(), nil
}

// validPrefix returns true if correct Prefix is given as an input
// e.g.: test/key/ is an correct prefix
// test/key/var is a full key path not a prefix
func validPrefix(input string) bool {
	err := utils.ValidateInputKey(input)
	return (err != nil) == true
}
