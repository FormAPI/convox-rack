package stdcli_test

import (
	"os"
	"testing"

	"github.com/convox/rack/cmd/convox/stdcli"
	"github.com/convox/rack/test"
	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	var opts map[string]string
	opts = stdcli.ParseOpts([]string{"--foo", "bar", "--key", "value"})
	assert.Equal(t, "bar", opts["foo"])
	assert.Equal(t, "value", opts["key"])

	opts = stdcli.ParseOpts([]string{"--foo=bar", "--key", "value"})
	assert.Equal(t, "bar", opts["foo"])
	assert.Equal(t, "value", opts["key"])

	opts = stdcli.ParseOpts([]string{"--foo=this", "is", "a bad idea"})
	assert.Equal(t, "this is a bad idea", opts["foo"])

	opts = stdcli.ParseOpts([]string{"--this", "--is=even", "worse"})
	assert.Equal(t, "even worse", opts["is"])
	_, ok := opts["this"]
	assert.Equal(t, true, ok)
}

// TestCheckEnvVars ensures stdcli.CheckEnv() prints a warning if bool envvars aren't true/false/1/0
func TestCheckEnvVars(t *testing.T) {
	os.Setenv("RACK_PRIVATE", "foo")

	// Ensure invalid env vars only print a warning and don't raise an error
	err := stdcli.CheckEnv()
	assert.NoError(t, err)

	test.Runs(t,
		test.ExecRun{
			Command:  "convox",
			Env:      map[string]string{"CONVOX_WAIT": "foo"},
			Exit:     0,
			OutMatch: "WARNING: 'foo' is not a valid value for environment variable CONVOX_WAIT (expected: [true false 1 0 ])\n",
		},
	)
}
