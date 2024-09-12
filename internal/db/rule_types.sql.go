// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: rule_types.sql

package db

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createRuleType = `-- name: CreateRuleType :one
INSERT INTO rule_type (
    name,
    project_id,
    description,
    guidance,
    definition,
    severity_value,
    subscription_id,
    display_name,
    release_phase,
    short_failure_message
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5::jsonb,
    $6,
    $7,
    $8,
    $9,
    $10
) RETURNING id, name, provider, project_id, description, guidance, definition, created_at, updated_at, severity_value, provider_id, subscription_id, display_name, release_phase, short_failure_message
`

type CreateRuleTypeParams struct {
	Name                string          `json:"name"`
	ProjectID           uuid.UUID       `json:"project_id"`
	Description         string          `json:"description"`
	Guidance            string          `json:"guidance"`
	Definition          json.RawMessage `json:"definition"`
	SeverityValue       Severity        `json:"severity_value"`
	SubscriptionID      uuid.NullUUID   `json:"subscription_id"`
	DisplayName         string          `json:"display_name"`
	ReleasePhase        ReleaseStatus   `json:"release_phase"`
	ShortFailureMessage string          `json:"short_failure_message"`
}

func (q *Queries) CreateRuleType(ctx context.Context, arg CreateRuleTypeParams) (RuleType, error) {
	row := q.db.QueryRowContext(ctx, createRuleType,
		arg.Name,
		arg.ProjectID,
		arg.Description,
		arg.Guidance,
		arg.Definition,
		arg.SeverityValue,
		arg.SubscriptionID,
		arg.DisplayName,
		arg.ReleasePhase,
		arg.ShortFailureMessage,
	)
	var i RuleType
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Provider,
		&i.ProjectID,
		&i.Description,
		&i.Guidance,
		&i.Definition,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SeverityValue,
		&i.ProviderID,
		&i.SubscriptionID,
		&i.DisplayName,
		&i.ReleasePhase,
		&i.ShortFailureMessage,
	)
	return i, err
}

const deleteRuleType = `-- name: DeleteRuleType :exec
DELETE FROM rule_type WHERE id = $1
`

func (q *Queries) DeleteRuleType(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteRuleType, id)
	return err
}

const getRuleTypeByID = `-- name: GetRuleTypeByID :one
SELECT id, name, provider, project_id, description, guidance, definition, created_at, updated_at, severity_value, provider_id, subscription_id, display_name, release_phase, short_failure_message FROM rule_type WHERE id = $1
`

func (q *Queries) GetRuleTypeByID(ctx context.Context, id uuid.UUID) (RuleType, error) {
	row := q.db.QueryRowContext(ctx, getRuleTypeByID, id)
	var i RuleType
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Provider,
		&i.ProjectID,
		&i.Description,
		&i.Guidance,
		&i.Definition,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SeverityValue,
		&i.ProviderID,
		&i.SubscriptionID,
		&i.DisplayName,
		&i.ReleasePhase,
		&i.ShortFailureMessage,
	)
	return i, err
}

const getRuleTypeByName = `-- name: GetRuleTypeByName :one
SELECT id, name, provider, project_id, description, guidance, definition, created_at, updated_at, severity_value, provider_id, subscription_id, display_name, release_phase, short_failure_message FROM rule_type WHERE  project_id = ANY($1::uuid[]) AND lower(name) = lower($2)
`

type GetRuleTypeByNameParams struct {
	Projects []uuid.UUID `json:"projects"`
	Name     string      `json:"name"`
}

func (q *Queries) GetRuleTypeByName(ctx context.Context, arg GetRuleTypeByNameParams) (RuleType, error) {
	row := q.db.QueryRowContext(ctx, getRuleTypeByName, pq.Array(arg.Projects), arg.Name)
	var i RuleType
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Provider,
		&i.ProjectID,
		&i.Description,
		&i.Guidance,
		&i.Definition,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SeverityValue,
		&i.ProviderID,
		&i.SubscriptionID,
		&i.DisplayName,
		&i.ReleasePhase,
		&i.ShortFailureMessage,
	)
	return i, err
}

const getRuleTypeNameByID = `-- name: GetRuleTypeNameByID :one
SELECT name FROM rule_type
WHERE id = $1
`

// intended as a temporary transition query
// this will be removed once the evaluation history tables replace the old state tables
func (q *Queries) GetRuleTypeNameByID(ctx context.Context, id uuid.UUID) (string, error) {
	row := q.db.QueryRowContext(ctx, getRuleTypeNameByID, id)
	var name string
	err := row.Scan(&name)
	return name, err
}

const getRuleTypesByEntityInHierarchy = `-- name: GetRuleTypesByEntityInHierarchy :many
SELECT rt.id, rt.name, rt.provider, rt.project_id, rt.description, rt.guidance, rt.definition, rt.created_at, rt.updated_at, rt.severity_value, rt.provider_id, rt.subscription_id, rt.display_name, rt.release_phase, rt.short_failure_message FROM rule_type AS rt
JOIN rule_instances AS ri ON ri.rule_type_id = rt.id
WHERE ri.entity_type = $1
AND ri.project_id = ANY($2::uuid[])
`

type GetRuleTypesByEntityInHierarchyParams struct {
	EntityType Entities    `json:"entity_type"`
	Projects   []uuid.UUID `json:"projects"`
}

func (q *Queries) GetRuleTypesByEntityInHierarchy(ctx context.Context, arg GetRuleTypesByEntityInHierarchyParams) ([]RuleType, error) {
	rows, err := q.db.QueryContext(ctx, getRuleTypesByEntityInHierarchy, arg.EntityType, pq.Array(arg.Projects))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []RuleType{}
	for rows.Next() {
		var i RuleType
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Provider,
			&i.ProjectID,
			&i.Description,
			&i.Guidance,
			&i.Definition,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.SeverityValue,
			&i.ProviderID,
			&i.SubscriptionID,
			&i.DisplayName,
			&i.ReleasePhase,
			&i.ShortFailureMessage,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRuleTypesByProject = `-- name: ListRuleTypesByProject :many
SELECT id, name, provider, project_id, description, guidance, definition, created_at, updated_at, severity_value, provider_id, subscription_id, display_name, release_phase, short_failure_message FROM rule_type WHERE project_id = $1
`

func (q *Queries) ListRuleTypesByProject(ctx context.Context, projectID uuid.UUID) ([]RuleType, error) {
	rows, err := q.db.QueryContext(ctx, listRuleTypesByProject, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []RuleType{}
	for rows.Next() {
		var i RuleType
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Provider,
			&i.ProjectID,
			&i.Description,
			&i.Guidance,
			&i.Definition,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.SeverityValue,
			&i.ProviderID,
			&i.SubscriptionID,
			&i.DisplayName,
			&i.ReleasePhase,
			&i.ShortFailureMessage,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateRuleType = `-- name: UpdateRuleType :one
UPDATE rule_type
    SET description = $2, definition = $3::jsonb, severity_value = $4, display_name = $5, release_phase = $6, short_failure_message = $7
    WHERE id = $1
    RETURNING id, name, provider, project_id, description, guidance, definition, created_at, updated_at, severity_value, provider_id, subscription_id, display_name, release_phase, short_failure_message
`

type UpdateRuleTypeParams struct {
	ID                  uuid.UUID       `json:"id"`
	Description         string          `json:"description"`
	Definition          json.RawMessage `json:"definition"`
	SeverityValue       Severity        `json:"severity_value"`
	DisplayName         string          `json:"display_name"`
	ReleasePhase        ReleaseStatus   `json:"release_phase"`
	ShortFailureMessage string          `json:"short_failure_message"`
}

func (q *Queries) UpdateRuleType(ctx context.Context, arg UpdateRuleTypeParams) (RuleType, error) {
	row := q.db.QueryRowContext(ctx, updateRuleType,
		arg.ID,
		arg.Description,
		arg.Definition,
		arg.SeverityValue,
		arg.DisplayName,
		arg.ReleasePhase,
		arg.ShortFailureMessage,
	)
	var i RuleType
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Provider,
		&i.ProjectID,
		&i.Description,
		&i.Guidance,
		&i.Definition,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SeverityValue,
		&i.ProviderID,
		&i.SubscriptionID,
		&i.DisplayName,
		&i.ReleasePhase,
		&i.ShortFailureMessage,
	)
	return i, err
}
