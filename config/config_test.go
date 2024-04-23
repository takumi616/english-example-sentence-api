package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetConfig(t *testing.T) {
	//Prepare expected config
	expectedConfig := &Config{
		Port:       "70000",
		DBHost:     "testHost",
		DBPort:     "5432",
		DBUser:     "testUser",
		DBPassword: "testPass",
		DBName:     "testDb",
		DBSSLMODE:  "disable",
	}

	//Set expected config data to environment variables
	//Setenv calls os.Setenv(key, value) and uses Cleanup to restore the environment variable
	//to its original value after the test.
	//Because Setenv affects the whole process, it cannot be used in parallel tests or tests with parallel ancestors.
	t.Setenv("APP_CONTAINER_PORT", expectedConfig.Port)
	t.Setenv("POSTGRES_HOST", expectedConfig.DBHost)
	t.Setenv("POSTGRES_PORT", expectedConfig.DBPort)
	t.Setenv("POSTGRES_USER", expectedConfig.DBUser)
	t.Setenv("POSTGRES_PASSWORD", expectedConfig.DBPassword)
	t.Setenv("POSTGRES_DB", expectedConfig.DBName)
	t.Setenv("POSTGRES_SSLMODE", expectedConfig.DBSSLMODE)

	//Call test target function
	actualConfig, err := GetConfig()
	if err != nil {
		t.Errorf("Failed to get config: %v", err)
	}

	//Compare result config to expected config
	if diff := cmp.Diff(expectedConfig, actualConfig); diff != "" {
		t.Errorf("Some differences are found: (-expected +actual)\n%s", diff)
	}
}
