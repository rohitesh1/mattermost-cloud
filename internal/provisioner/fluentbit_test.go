// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package provisioner

import (
	"testing"

	mocks "github.com/mattermost/mattermost-cloud/internal/mocks/aws-tools"
	"github.com/mattermost/mattermost-cloud/internal/tools/kops"
	"github.com/mattermost/mattermost-cloud/model"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestNewHelmDeploymentWithDefaultConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provisioner := &KopsProvisioner{}
	logger := log.New()
	awsClient := mocks.NewMockAWS(ctrl)
	kops := &kops.Cmd{}
	fluentbit, err := newFluentbitHandle(&model.Cluster{
		UtilityMetadata: &model.UtilityMetadata{
			ActualVersions: model.UtilityGroupVersions{},
		},
	}, &model.HelmUtilityVersion{Chart: "1.2.3"}, provisioner, awsClient, kops, logger)
	require.NoError(t, err, "should not error when creating new fluentbit handler")
	require.NotNil(t, fluentbit, "fluentbit should not be nil")

	helmDeployment := fluentbit.NewHelmDeployment(logger)
	require.NotNil(t, helmDeployment, "helmDeployment should not be nil")
}
