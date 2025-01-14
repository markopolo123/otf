// Package client provides an abstraction for interacting with otf services
// either remotely or locally.
package client

import (
	"context"

	"github.com/leg100/otf"
	"github.com/leg100/otf/auth"
	"github.com/leg100/otf/configversion"
	"github.com/leg100/otf/http"
	"github.com/leg100/otf/logs"
	"github.com/leg100/otf/organization"
	"github.com/leg100/otf/run"
	"github.com/leg100/otf/state"
	"github.com/leg100/otf/variable"
	"github.com/leg100/otf/workspace"
)

var (
	_ Client = (*LocalClient)(nil)
	_ Client = (*remoteClient)(nil)
)

type (
	// Client is those service endpoints that support both in-process and remote
	// invocation. Intended for use with the agent (the internal agent is
	// in-process, while the external agent is remote) as well as the CLI.
	Client interface {
		CreateOrganization(ctx context.Context, opts organization.OrganizationCreateOptions) (*organization.Organization, error)
		DeleteOrganization(ctx context.Context, organization string) error

		GetWorkspace(ctx context.Context, workspaceID string) (*workspace.Workspace, error)
		GetWorkspaceByName(ctx context.Context, organization, workspace string) (*workspace.Workspace, error)
		ListWorkspaces(ctx context.Context, opts workspace.ListOptions) (*workspace.WorkspaceList, error)
		UpdateWorkspace(ctx context.Context, workspaceID string, opts workspace.UpdateOptions) (*workspace.Workspace, error)

		ListVariables(ctx context.Context, workspaceID string) ([]*variable.Variable, error)

		CreateAgentToken(ctx context.Context, opts auth.CreateAgentTokenOptions) (*auth.AgentToken, error)
		GetAgentToken(ctx context.Context, token string) (*auth.AgentToken, error)

		GetPlanFile(ctx context.Context, id string, format run.PlanFormat) ([]byte, error)
		UploadPlanFile(ctx context.Context, id string, plan []byte, format run.PlanFormat) error

		GetLockFile(ctx context.Context, id string) ([]byte, error)
		UploadLockFile(ctx context.Context, id string, lockFile []byte) error

		ListRuns(ctx context.Context, opts run.RunListOptions) (*run.RunList, error)
		GetRun(ctx context.Context, id string) (*run.Run, error)

		StartPhase(ctx context.Context, id string, phase otf.PhaseType, opts run.PhaseStartOptions) (*run.Run, error)
		FinishPhase(ctx context.Context, id string, phase otf.PhaseType, opts run.PhaseFinishOptions) (*run.Run, error)

		DownloadConfig(ctx context.Context, id string) ([]byte, error)

		Watch(context.Context, run.WatchOptions) (<-chan otf.Event, error)

		// CreateRegistrySession creates a registry session for the given organization.
		CreateRegistrySession(ctx context.Context, opts auth.CreateRegistrySessionOptions) (*auth.RegistrySession, error)

		CreateStateVersion(ctx context.Context, opts state.CreateStateVersionOptions) (*state.Version, error)
		DownloadCurrentState(ctx context.Context, workspaceID string) ([]byte, error)

		CreateUser(ctx context.Context, username string, opts ...auth.NewUserOption) (*auth.User, error)
		DeleteUser(ctx context.Context, username string) error
		AddOrganizationMembership(ctx context.Context, username, organization string) error
		RemoveOrganizationMembership(ctx context.Context, username, organization string) error
		AddTeamMembership(ctx context.Context, username, teamID string) error
		RemoveTeamMembership(ctx context.Context, username, teamID string) error

		CreateTeam(ctx context.Context, opts auth.NewTeamOptions) (*auth.Team, error)
		GetTeam(ctx context.Context, organization, team string) (*auth.Team, error)
		DeleteTeam(ctx context.Context, teamID string) error

		Hostname() string

		otf.PutChunkService
		workspace.LockService
	}

	LocalClient struct {
		organization.OrganizationService
		auth.AuthService
		variable.VariableService
		state.StateService
		workspace.WorkspaceService
		otf.HostnameService
		configversion.ConfigurationVersionService
		run.RunService
		logs.LogsService
	}

	remoteClient struct {
		*http.Client
		http.Config

		*stateClient
		*configClient
		*variableClient
		*authClient
		*organizationClient
		*workspaceClient
		*runClient
		*logsClient
	}

	stateClient        = state.Client
	configClient       = configversion.Client
	variableClient     = variable.Client
	authClient         = auth.Client
	organizationClient = organization.Client
	workspaceClient    = workspace.Client
	runClient          = run.Client
	logsClient         = logs.Client
)

// New constructs a client that uses the http to remotely invoke OTF
// services.
func New(config http.Config) (*remoteClient, error) {
	httpClient, err := http.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &remoteClient{
		Client:             httpClient,
		stateClient:        &stateClient{httpClient},
		configClient:       &configClient{httpClient},
		variableClient:     &variableClient{httpClient},
		authClient:         &authClient{httpClient},
		organizationClient: &organizationClient{httpClient},
		workspaceClient:    &workspaceClient{httpClient},
		runClient:          &runClient{httpClient, config},
		logsClient:         &logsClient{httpClient},
	}, nil
}
