package github

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	githubToken := os.Getenv(EnvToken)
	defer func() {
		os.Setenv(EnvToken, githubToken)
	}()
	os.Setenv(EnvToken, "")

	testCases := []struct {
		config   Config
		envToken string
		expect   string
	}{
		{
			// specify directly
			config:   Config{Token: "abcdefg"},
			envToken: "",
			expect:   "",
		},
		{
			// specify via env but not to be set env (part 1)
			config:   Config{Token: "GITHUB_TOKEN"},
			envToken: "",
			expect:   "github token is missing",
		},
		{
			// specify via env (part 1)
			config:   Config{Token: "GITHUB_TOKEN"},
			envToken: "abcdefg",
			expect:   "",
		},
		{
			// specify via env but not to be set env (part 2)
			config:   Config{Token: "$GITHUB_TOKEN"},
			envToken: "",
			expect:   "github token is missing",
		},
		{
			// specify via env (part 2)
			config:   Config{Token: "$GITHUB_TOKEN"},
			envToken: "abcdefg",
			expect:   "",
		},
		{
			// no specification (part 1)
			config:   Config{},
			envToken: "",
			expect:   "github token is missing",
		},
		{
			// no specification (part 2)
			config:   Config{},
			envToken: "abcdefg",
			expect:   "github token is missing",
		},
	}
	for _, testCase := range testCases {
		os.Setenv(EnvToken, testCase.envToken)
		_, err := NewClient(testCase.config)
		if err == nil {
			continue
		}
		if err.Error() != testCase.expect {
			t.Errorf("got %q but want %q", err.Error(), testCase.expect)
		}
	}
}
