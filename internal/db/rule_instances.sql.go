// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: rule_instances.sql

package db

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const deleteNonUpdatedRules = `-- name: DeleteNonUpdatedRules :exec
DELETE FROM rule_instances
WHERE profile_id = $1
AND entity_type = $2
AND NOT id = ANY($3::UUID[])
`

type DeleteNonUpdatedRulesParams struct {
	ProfileID  uuid.UUID   `json:"profile_id"`
	EntityType Entities    `json:"entity_type"`
	UpdatedIds []uuid.UUID `json:"updated_ids"`
}

func (q *Queries) DeleteNonUpdatedRules(ctx context.Context, arg DeleteNonUpdatedRulesParams) error {
	_, err := q.db.ExecContext(ctx, deleteNonUpdatedRules, arg.ProfileID, arg.EntityType, pq.Array(arg.UpdatedIds))
	return err
}

const getRuleInstancesForProfile = `-- name: GetRuleInstancesForProfile :many
SELECT id, profile_id, rule_type_id, name, entity_type, def, params, created_at, updated_at FROM rule_instances WHERE profile_id = $1
`

func (q *Queries) GetRuleInstancesForProfile(ctx context.Context, profileID uuid.UUID) ([]RuleInstance, error) {
	rows, err := q.db.QueryContext(ctx, getRuleInstancesForProfile, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []RuleInstance{}
	for rows.Next() {
		var i RuleInstance
		if err := rows.Scan(
			&i.ID,
			&i.ProfileID,
			&i.RuleTypeID,
			&i.Name,
			&i.EntityType,
			&i.Def,
			&i.Params,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getRuleInstancesForProfileEntity = `-- name: GetRuleInstancesForProfileEntity :many
SELECT id, profile_id, rule_type_id, name, entity_type, def, params, created_at, updated_at FROM rule_instances WHERE profile_id = $1 AND entity_type = $2
`

type GetRuleInstancesForProfileEntityParams struct {
	ProfileID  uuid.UUID `json:"profile_id"`
	EntityType Entities  `json:"entity_type"`
}

func (q *Queries) GetRuleInstancesForProfileEntity(ctx context.Context, arg GetRuleInstancesForProfileEntityParams) ([]RuleInstance, error) {
	rows, err := q.db.QueryContext(ctx, getRuleInstancesForProfileEntity, arg.ProfileID, arg.EntityType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []RuleInstance{}
	for rows.Next() {
		var i RuleInstance
		if err := rows.Scan(
			&i.ID,
			&i.ProfileID,
			&i.RuleTypeID,
			&i.Name,
			&i.EntityType,
			&i.Def,
			&i.Params,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const upsertRuleInstance = `-- name: UpsertRuleInstance :one

INSERT INTO rule_instances (
    profile_id,
    rule_type_id,
    name,
    entity_type,
    def,
    params,
    created_at,
    updated_at
) VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    NOW(),
    NOW()
)
ON CONFLICT (profile_id, entity_type, name) DO UPDATE SET
    rule_type_id = $2,
    def = $5,
    params = $6,
    updated_at = NOW()
RETURNING id
`

type UpsertRuleInstanceParams struct {
	ProfileID  uuid.UUID       `json:"profile_id"`
	RuleTypeID uuid.UUID       `json:"rule_type_id"`
	Name       string          `json:"name"`
	EntityType Entities        `json:"entity_type"`
	Def        json.RawMessage `json:"def"`
	Params     json.RawMessage `json:"params"`
}

// Copyright 2024 Stacklok, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
func (q *Queries) UpsertRuleInstance(ctx context.Context, arg UpsertRuleInstanceParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, upsertRuleInstance,
		arg.ProfileID,
		arg.RuleTypeID,
		arg.Name,
		arg.EntityType,
		arg.Def,
		arg.Params,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
