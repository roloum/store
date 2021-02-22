package test

import (
	"os"
)

//SetEnvironment sets the Environment variables required to run the test cases
func SetEnvironment() {

	vars := map[string]string{
		"AWS_DYNAMODB_TABLE_STORE": "Store",
	}

	for key, value := range vars {
		_ = os.Setenv(key, value)
	}
}
