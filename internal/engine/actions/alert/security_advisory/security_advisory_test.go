// SPDX-FileCopyrightText: Copyright 2024 The Minder Authors
// SPDX-License-Identifier: Apache-2.0

// Package security_advisory provides necessary interfaces and implementations for
// creating alerts of type security advisory.

package security_advisory

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/mindersec/minder/internal/db"
	enginerr "github.com/mindersec/minder/internal/engine/errors"
	"github.com/mindersec/minder/internal/engine/interfaces"
	pbinternal "github.com/mindersec/minder/internal/proto"
	mockghclient "github.com/mindersec/minder/internal/providers/github/mock"
	pb "github.com/mindersec/minder/pkg/api/protobuf/go/minder/v1"
	"github.com/mindersec/minder/pkg/profiles/models"
)

var TestActionTypeValid interfaces.ActionType = "alert-test"

func TestSecurityAdvisoryAlert(t *testing.T) {
	t.Parallel()

	saID := "123"

	tests := []struct {
		name             string
		actionType       interfaces.ActionType
		mockSetup        func(*mockghclient.MockGitHub)
		expectedErr      error
		expectedMetadata json.RawMessage
	}{
		{
			name:       "create a security advisory",
			actionType: TestActionTypeValid,
			mockSetup: func(mockGitHub *mockghclient.MockGitHub) {
				mockGitHub.EXPECT().
					CreateSecurityAdvisory(gomock.Any(), gomock.Any(), gomock.Any(), pb.Severity_VALUE_HIGH.String(),
						gomock.Any(), gomock.Any(), gomock.Any()).
					Return(saID, nil)
			},
			expectedErr:      nil,
			expectedMetadata: json.RawMessage(fmt.Sprintf(`{"ghsa_id":"%s"}`, saID)),
		},
		{
			name:       "error from provider creating security advisory",
			actionType: TestActionTypeValid,
			mockSetup: func(mockGitHub *mockghclient.MockGitHub) {
				mockGitHub.EXPECT().
					CreateSecurityAdvisory(gomock.Any(), gomock.Any(), gomock.Any(), pb.Severity_VALUE_HIGH.String(),
						gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", fmt.Errorf("failed to create security advisory"))
			},
			expectedErr:      enginerr.ErrActionFailed,
			expectedMetadata: json.RawMessage(nil),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			t.Cleanup(func() {
				ctrl.Finish()
			})

			ruleType := pb.RuleType{
				Name:                "rule_type_1",
				ShortFailureMessage: "This is a failure message",
				Def: &pb.RuleType_Definition{
					Alert:     &pb.RuleType_Definition_Alert{},
					Remediate: &pb.RuleType_Definition_Remediate{},
				},
			}
			saCfg := pb.RuleType_Definition_Alert_AlertTypeSA{
				Severity: pb.Severity_VALUE_HIGH.String(),
			}

			mockClient := mockghclient.NewMockGitHub(ctrl)
			tt.mockSetup(mockClient)

			saAlert, err := NewSecurityAdvisoryAlert(
				tt.actionType, &ruleType, &saCfg, mockClient, models.ActionOptOn)
			require.NoError(t, err)
			require.NotNil(t, saAlert)

			evalParams := &interfaces.EvalStatusParams{
				EvalStatusFromDb: &db.ListRuleEvaluationsByProfileIdRow{},
				Profile:          &models.ProfileAggregate{},
				Rule:             &models.RuleInstance{},
			}

			retMeta, err := saAlert.Do(
				context.Background(),
				interfaces.ActionCmdOn,
				&pbinternal.PullRequest{},
				evalParams,
				nil,
			)
			require.ErrorIs(t, err, tt.expectedErr, "expected error")
			require.Equal(t, tt.expectedMetadata, retMeta)
		})
	}
}
