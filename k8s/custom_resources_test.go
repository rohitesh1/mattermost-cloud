// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package k8s

import (
	mmv1alpha1 "github.com/mattermost/mattermost-operator/apis/mattermost/v1alpha1"
	"github.com/stretchr/testify/require"
	apixv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apixv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestCustomResourceDefinitionBetaV1(t *testing.T) {
	testClient := newTestKubeClient()
	customResourceDefinition := &apixv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: "test-crd"},
	}

	t.Run("create custom resource definition", func(t *testing.T) {
		result, err := testClient.createOrUpdateCustomResourceDefinitionBetaV1(customResourceDefinition)
		require.NoError(t, err)
		require.Equal(t, customResourceDefinition.GetName(), result.GetName())
	})
	t.Run("create duplicate custom resource definition", func(t *testing.T) {
		result, err := testClient.createOrUpdateCustomResourceDefinitionBetaV1(customResourceDefinition)
		require.NoError(t, err)
		require.Equal(t, customResourceDefinition.GetName(), result.GetName())
	})
}

func TestCustomResourceDefinitionV1(t *testing.T) {
	testClient := newTestKubeClient()
	customResourceDefinition := &apixv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: "test-crd"},
	}

	t.Run("create custom resource definition", func(t *testing.T) {
		result, err := testClient.createOrUpdateCustomResourceDefinitionV1(customResourceDefinition)
		require.NoError(t, err)
		require.Equal(t, customResourceDefinition.GetName(), result.GetName())
	})
	t.Run("create duplicate custom resource definition", func(t *testing.T) {
		result, err := testClient.createOrUpdateCustomResourceDefinitionV1(customResourceDefinition)
		require.NoError(t, err)
		require.Equal(t, customResourceDefinition.GetName(), result.GetName())
	})
}

func TestClusterInstallation(t *testing.T) {
	testClient := newTestKubeClient()
	customResource := &mmv1alpha1.ClusterInstallation{
		ObjectMeta: metav1.ObjectMeta{Name: "test-cr"},
	}
	namespace := "testing"

	t.Run("create custom resource", func(t *testing.T) {
		result, err := testClient.createOrUpdateClusterInstallation(namespace, customResource)
		require.NoError(t, err)
		require.Equal(t, customResource.GetName(), result.GetName())
	})
	t.Run("create duplicate custom resource", func(t *testing.T) {
		result, err := testClient.createOrUpdateClusterInstallation(namespace, customResource)
		require.NoError(t, err)
		require.Equal(t, customResource.GetName(), result.GetName())
	})
}
