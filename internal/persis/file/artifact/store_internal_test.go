// Copyright (C) 2026 Yota Hamada
// SPDX-License-Identifier: GPL-3.0-or-later

package artifact

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// outsideBounds decides whether a whole year or month directory can be skipped,
// so an off-by-one here turns into either missing results or reading the entire
// retained history.
func TestOutsideBounds(t *testing.T) {
	t.Parallel()

	const from, to = "2026/09/15", "2026/11/20"

	tests := []struct {
		key     string
		outside bool
	}{
		{"2026", false},
		{"2025", true},
		{"2027", true},
		{"2026/09", false},
		{"2026/10", false},
		{"2026/11", false},
		{"2026/08", true},
		{"2026/12", true},
		{"2026/09/15", false},
		{"2026/09/14", true},
		{"2026/11/20", false},
		{"2026/11/21", true},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.outside, outsideBounds(tt.key, from, to), tt.key)
	}
}

func TestOutsideBoundsUnbounded(t *testing.T) {
	t.Parallel()

	assert.False(t, outsideBounds("1999", "", ""))
	assert.False(t, outsideBounds("2030/01/01", "", ""))

	// A bound set on one side only must not constrain the other.
	assert.True(t, outsideBounds("2025/12/31", "2026/01/01", ""))
	assert.False(t, outsideBounds("2030/01/01", "2026/01/01", ""))
	assert.True(t, outsideBounds("2027/01/01", "", "2026/12/31"))
	assert.False(t, outsideBounds("2020/01/01", "", "2026/12/31"))
}
