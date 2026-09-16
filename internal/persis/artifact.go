// Copyright (C) 2026 Yota Hamada
// SPDX-License-Identifier: GPL-3.0-or-later

package persis

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/dagucloud/dagu/v2/internal/workspace"
)

var (
	// ErrInvalidArtifactCursor reports a cursor that does not belong to the
	// query it was presented with.
	ErrInvalidArtifactCursor = errors.New("invalid artifact cursor")

	// ErrInvalidArtifactFileName reports a file name pattern that is not a
	// valid glob.
	ErrInvalidArtifactFileName = errors.New("invalid artifact file name pattern")
)

// ArtifactStore lists files produced by DAG runs, newest run first.
type ArtifactStore interface {
	QueryArtifacts(ctx context.Context, query ArtifactQuery) (ArtifactPage, error)
}

// ArtifactQuery selects a page of artifact files.
type ArtifactQuery struct {
	// Name matches DAG names containing it, case-insensitively.
	Name string

	// FileName matches artifact paths. A pattern holding glob metacharacters
	// is a glob; anything else is a case-insensitive substring.
	FileName string

	From TimeInUTC
	To   TimeInUTC

	WorkspaceFilter *workspace.WorkspaceFilter

	Limit  int
	Cursor string
}

// ArtifactPage is one forward-only page of artifact files.
type ArtifactPage struct {
	Items      []ArtifactFile
	NextCursor string
}

// ArtifactFile is a single file produced by a DAG run.
type ArtifactFile struct {
	Name      string
	DAGRunID  string
	StartedAt time.Time

	// Path is relative to the run's artifact directory and uses forward
	// slashes on every platform.
	Path string
	Size int64
}

// globMetacharacters are the characters that make a file name pattern a glob
// rather than a substring.
const globMetacharacters = "*?[{"

// IsArtifactFileNameGlob reports whether a file name pattern is a glob.
func IsArtifactFileNameGlob(pattern string) bool {
	return strings.ContainsAny(pattern, globMetacharacters)
}

// MatchArtifactFileName reports whether a run-relative artifact path matches a
// file name pattern. A pattern holding glob metacharacters is matched as a
// glob, the way artifact.list matches its own pattern; anything else is matched
// as a case-insensitive substring.
func MatchArtifactFileName(path, pattern string) bool {
	if pattern == "" {
		return true
	}
	if IsArtifactFileNameGlob(pattern) {
		ok, err := doublestar.Match(pattern, path)
		return err == nil && ok
	}
	return strings.Contains(strings.ToLower(path), strings.ToLower(pattern))
}
