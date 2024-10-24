/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package requirements

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/deckhouse/deckhouse/go_lib/dependency/requirements"
)

func TestKubernetesVersionRequirement(t *testing.T) {
	t.Run("complies with the requirements", func(t *testing.T) {
		requirements.SaveValue(metallbConfigurationStatusKey, "")
		ok, err := requirements.CheckRequirement(metallbConfigurationStatusRequirementsKey, "")
		assert.True(t, ok)
		require.NoError(t, err)
	})

	t.Run("fail: nsMismatch", func(t *testing.T) {
		requirements.SaveValue(metallbConfigurationStatusKey, "nsMismatch")
		ok, err := requirements.CheckRequirement(metallbConfigurationStatusRequirementsKey, "")
		assert.False(t, ok)
		require.Error(t, err)
	})

	t.Run("fail: nodeSelectorsMismatch", func(t *testing.T) {
		requirements.SaveValue(metallbConfigurationStatusKey, "nodeSelectorsMismatch")
		ok, err := requirements.CheckRequirement(metallbConfigurationStatusRequirementsKey, "")
		assert.False(t, ok)
		require.Error(t, err)
	})

	t.Run("fail: addressPollsMismatch", func(t *testing.T) {
		requirements.SaveValue(metallbConfigurationStatusKey, "nodeSelectorsMismatch")
		ok, err := requirements.CheckRequirement(metallbConfigurationStatusRequirementsKey, "")
		assert.False(t, ok)
		require.Error(t, err)
	})
}
