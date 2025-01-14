package workspace

import (
	"context"

	"github.com/leg100/otf"
	"github.com/leg100/otf/rbac"
	"github.com/leg100/otf/sql"
	"github.com/leg100/otf/sql/pggen"
)

func (db *pgdb) SetWorkspacePermission(ctx context.Context, workspaceID, team string, role rbac.Role) error {
	_, err := db.UpsertWorkspacePermission(ctx, pggen.UpsertWorkspacePermissionParams{
		WorkspaceID: sql.String(workspaceID),
		TeamName:    sql.String(team),
		Role:        sql.String(role.String()),
	})
	if err != nil {
		return sql.Error(err)
	}
	return nil
}

func (db *pgdb) GetWorkspacePolicy(ctx context.Context, workspaceID string) (otf.WorkspacePolicy, error) {
	result, err := db.FindWorkspacePolicyByID(ctx, sql.String(workspaceID))
	if err != nil {
		return otf.WorkspacePolicy{}, sql.Error(err)
	}
	policy := otf.WorkspacePolicy{
		Organization: result.OrganizationName.String,
		WorkspaceID:  result.WorkspaceID.String,
	}
	// SQL query returns an array of workspace permissions and an array of
	// teams; the former has the team id but not the team name, which we want, so
	// lookup the corresponding team name in the latter.
	for _, perm := range result.WorkspacePermissions {
		role, err := rbac.WorkspaceRoleFromString(perm.Role.String)
		if err != nil {
			return otf.WorkspacePolicy{}, err
		}
		// find corresponding team name in teams array
		for _, t := range result.Teams {
			if t.TeamID == perm.TeamID {
				policy.Permissions = append(policy.Permissions, otf.WorkspacePermission{
					Team: t.Name.String,
					Role: role,
				})
				break
			}
		}
	}
	return policy, nil
}

func (db *pgdb) UnsetWorkspacePermission(ctx context.Context, workspaceID, team string) error {
	_, err := db.DeleteWorkspacePermissionByID(ctx, sql.String(workspaceID), sql.String(team))
	if err != nil {
		return sql.Error(err)
	}
	return nil
}
