package e2e

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTeamCLI tests managing teams via the CLI
func TestTeamCLI(t *testing.T) {
	setup(t)

	daemon := &daemon{}
	daemon.withFlags("--site-token", "abc123")
	hostname := daemon.start(t)

	// create organization
	cmd := exec.Command("otf", "organizations", "new", "acme-corp",
		"--address", hostname,
		"--token", "abc123",
	)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(t, "Successfully created organization acme-corp\n", string(out))

	// create team
	cmd = exec.Command("otf", "teams", "new", "owners",
		"--address", hostname,
		"--token", "abc123",
		"--organization", "acme-corp",
	)
	out, err = cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(t, "Successfully created team owners\n", string(out))

	// create user
	cmd = exec.Command("otf", "users", "new", "bobby",
		"--address", hostname,
		"--token", "abc123",
	)
	out, err = cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(t, "Successfully created user bobby\n", string(out))

	// add user to team
	cmd = exec.Command("otf", "teams", "add-user", "bobby",
		"--address", hostname,
		"--token", "abc123",
		"--organization", "acme-corp",
		"--team", "owners",
	)
	out, err = cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(t, "Successfully added bobby to owners\n", string(out))

	// remove user from team
	cmd = exec.Command("otf", "teams", "del-user", "bobby",
		"--address", hostname,
		"--token", "abc123",
		"--organization", "acme-corp",
		"--team", "owners",
	)
	out, err = cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(t, "Successfully removed bobby from owners\n", string(out))

	// delete team
	cmd = exec.Command("otf", "teams", "delete", "owners",
		"--address", hostname,
		"--token", "abc123",
		"--organization", "acme-corp",
	)
	out, err = cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(t, "Successfully deleted team owners\n", string(out))
}
