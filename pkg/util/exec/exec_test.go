/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package exec

import (
	"testing"

	"gotest.tools/assert"
)

func TestLookPath(t *testing.T) {
	t.Run("system cmd", func(t *testing.T) {
		t.Run("ps", func(t *testing.T) {
			pspath, psstatus := LookPath("ps")
			assert.Equal(t, psstatus, pspath == "/bin/ps")
		})
	})
	t.Run("ergo cmd", func(t *testing.T) {
		t.Run("k3d", func(t *testing.T) {
			pspath, psstatus := LookPath("k3d")
			assert.Equal(t, psstatus, pspath == "/bin/ps")
		})
	})
}
