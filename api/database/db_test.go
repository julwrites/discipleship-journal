package database

import (
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestBuildConnectionString(t *testing.T) {
	// Save current env vars
	originalUsername := os.Getenv("DB_USERNAME")
	originalPassword := os.Getenv("DB_PASSWORD")
	originalDBName := os.Getenv("DB_NAME")
	originalHost := os.Getenv("DB_HOST")
	originalPort := os.Getenv("DB_PORT")
	originalCloudInstance := os.Getenv("CLOUD_SQL_INSTANCE")

	// Restore env vars after test
	defer func() {
		os.Setenv("DB_USERNAME", originalUsername)
		os.Setenv("DB_PASSWORD", originalPassword)
		os.Setenv("DB_NAME", originalDBName)
		os.Setenv("DB_HOST", originalHost)
		os.Setenv("DB_PORT", originalPort)
		os.Setenv("CLOUD_SQL_INSTANCE", originalCloudInstance)
	}()

	tests := []struct {
		name          string
		username      string
		password      string
		dbName        string
		cloudInstance string // If empty, standard TCP
		host          string
		port          string
		wantContains  []string
	}{
		{
			name:          "Standard TCP with simple password",
			username:      "user",
			password:      "pass",
			dbName:        "mydb",
			cloudInstance: "",
			host:          "localhost",
			port:          "5432",
			wantContains:  []string{"postgres://user:pass@localhost:5432/mydb", "sslmode=disable"},
		},
		{
			name:          "Standard TCP with special chars in password",
			username:      "user",
			password:      "pass%^&word",
			dbName:        "mydb",
			cloudInstance: "",
			host:          "db-host",
			port:          "5432",
			// & is NOT encoded by url.UserPassword as it is a sub-delim allowed in userinfo
			wantContains:  []string{"postgres://user:pass%25%5E&word@db-host:5432/mydb", "sslmode=disable"},
		},
		{
			name:          "Cloud SQL with complex password",
			username:      "postgres",
			password:      "yd%^9VLMi[J=%d:xxxxxx",
			dbName:        "postgres",
			cloudInstance: "project:region:instance",
			wantContains: []string{
				"postgres://postgres:yd%25%5E9VLMi%5BJ=%25d%3Axxxxxx@cloudsql/postgres",
				"sslmode=disable",
			},
		},
		{
			name:          "Cloud SQL with IAM Auth (No Password)",
			username:      "sa-email@project.iam",
			password:      "",
			dbName:        "mydb",
			cloudInstance: "project:region:instance",
			wantContains: []string{
				// User only, no password
				"postgres://sa-email%40project.iam@cloudsql/mydb",
				"sslmode=disable",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DB_USERNAME", tt.username)
			os.Setenv("DB_PASSWORD", tt.password)
			os.Setenv("DB_NAME", tt.dbName)
			os.Setenv("CLOUD_SQL_INSTANCE", tt.cloudInstance)
			if tt.cloudInstance == "" {
				os.Setenv("DB_HOST", tt.host)
				os.Setenv("DB_PORT", tt.port)
			} else {
				os.Unsetenv("DB_HOST")
				os.Unsetenv("DB_PORT")
			}

			got, err := BuildConnectionString()
			if err != nil {
				t.Fatalf("BuildConnectionString() error = %v", err)
			}
			t.Logf("Got: %s", got)

			// Parse result URL to verify decoding works
			u, err := url.Parse(got)
			if err != nil {
				t.Fatalf("Result is not a valid URL: %v", err)
			}

			// Only check password if it was provided
			if tt.password != "" {
				gotPass, _ := u.User.Password()
				if gotPass != tt.password {
					t.Errorf("Decoded password = %v, want %v", gotPass, tt.password)
				}
			}

			// Check raw string components
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("BuildConnectionString() = %v, want to contain %v", got, want)
				}
			}
		})
	}
}
