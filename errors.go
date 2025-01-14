package otf

import (
	"errors"
)

// Generic errors applicable to all resources.
var (
	// ErrAccessNotPermitted is returned when an authorization check fails.
	ErrAccessNotPermitted = errors.New("access to the resource is not permitted")

	// ErrUnauthorized is returned when a receiving a 401.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrResourceNotFound is returned when a receiving a 404.
	ErrResourceNotFound = errors.New("resource not found")

	// ErrResourceAlreadyExists is returned when attempting to create a resource
	// that already exists.
	ErrResourceAlreadyExists = errors.New("resource already exists")

	// ErrRequiredName is returned when a name option is not present.
	ErrRequiredName = errors.New("name is required")

	// ErrInvalidName is returned when the name option has invalid value.
	ErrInvalidName = errors.New("invalid value for name")

	// ErrForeignKeyViolation is returned when attempting to delete or
	// update a resource that is referenced by another resource and the
	// delete/update would orphan the reference.
	ErrForeignKeyViolation = errors.New("foreign key constraint violation")

	// ErrWarning is a non-fatal error
	ErrWarning = errors.New("warning")
)

// Resource Errors
var (
	// ErrInvalidTerraformVersion is returned when a terraform version string is
	// not a semantic version string (major.minor.patch).
	ErrInvalidTerraformVersion = errors.New("invalid terraform version")

	// ErrInvalidWorkspaceID is returned when the workspace ID is invalid.
	ErrInvalidWorkspaceID = errors.New("invalid value for workspace ID")

	// ErrInvalidWorkspaceValue is returned when workspace value is invalid.
	ErrInvalidWorkspaceValue = errors.New("invalid value for workspace")

	// Organization errors

	// ErrInvalidOrg is returned when the organization option has an invalid value.
	ErrInvalidOrg = errors.New("invalid value for organization")

	// ErrRequiredOrg is returned when the organization option is not present
	ErrRequiredOrg = errors.New("organization is required")

	ErrStatusTimestampNotFound = errors.New("corresponding status timestamp not found")

	ErrInvalidRepo = errors.New("repository path is invalid")
)
