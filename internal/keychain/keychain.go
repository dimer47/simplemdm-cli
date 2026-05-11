package keychain

import "github.com/zalando/go-keyring"

const serviceName = "simplemdm-cli"

func keyFor(context string) string {
	return serviceName + ":" + context
}

func Set(context, apiKey string) error {
	return keyring.Set(serviceName, keyFor(context), apiKey)
}

func Get(context string) (string, error) {
	return keyring.Get(serviceName, keyFor(context))
}

func Delete(context string) error {
	return keyring.Delete(serviceName, keyFor(context))
}

func IsAvailable() bool {
	err := keyring.Set(serviceName, serviceName+":__test__", "test")
	if err != nil {
		return false
	}
	_ = keyring.Delete(serviceName, serviceName+":__test__")
	return true
}
