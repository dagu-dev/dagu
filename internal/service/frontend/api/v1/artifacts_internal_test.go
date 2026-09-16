// Copyright (C) 2026 Yota Hamada
// SPDX-License-Identifier: GPL-3.0-or-later

package api

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	openapiv1 "github.com/dagucloud/dagu/v2/api/v1"
	"github.com/dagucloud/dagu/v2/internal/cmn/artifactpath"
	"github.com/dagucloud/dagu/v2/internal/cmn/stringutil"
	"github.com/dagucloud/dagu/v2/internal/persis"
	fileartifact "github.com/dagucloud/dagu/v2/internal/persis/file/artifact"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var artifactTestStart = time.Date(2026, 9, 15, 14, 32, 7, 0, time.UTC)

// newArtifactListAPI builds an API over a populated artifact tree.
func newArtifactListAPI(t *testing.T, runs ...string) *API {
	t.Helper()

	root := t.TempDir()
	for i, dagRunID := range runs {
		at := artifactTestStart.Add(time.Duration(i) * time.Minute)
		dir, err := artifactpath.NewRunDir(context.Background(), root, "", "reporter", dagRunID, at)
		require.NoError(t, err)
		for _, name := range []string{"out.txt", "reports/summary.md"} {
			path := filepath.Join(dir, filepath.FromSlash(name))
			require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o750))
			require.NoError(t, os.WriteFile(path, []byte("hello"), 0o600))
		}

		metaPath, ok := artifactpath.MetaPath(root, dir)
		require.True(t, ok)
		require.NoError(t, fileartifact.WriteRecord(metaPath, fileartifact.Record{
			Version:   fileartifact.RecordVersion,
			Name:      "reporter",
			DAGRunID:  dagRunID,
			StartedAt: stringutil.FormatTime(at),
			Dir:       dir,
		}))
	}

	return &API{artifactRepository: persis.NewArtifactRepository(fileartifact.NewStore(root))}
}

func listArtifacts(t *testing.T, a *API, params openapiv1.ListArtifactsParams) openapiv1.ArtifactListResponse {
	t.Helper()

	resp, err := a.ListArtifacts(context.Background(), openapiv1.ListArtifactsRequestObject{Params: params})
	require.NoError(t, err)
	body, ok := resp.(openapiv1.ListArtifacts200JSONResponse)
	require.True(t, ok, "expected a 200 response, got %T", resp)
	return openapiv1.ArtifactListResponse(body)
}

func itemRunIDs(body openapiv1.ArtifactListResponse) []string {
	ids := make([]string, 0, len(body.Items))
	for _, item := range body.Items {
		ids = append(ids, item.DagRunId)
	}
	return ids
}

func TestListArtifacts(t *testing.T) {
	t.Run("NewestRunFirst", func(t *testing.T) {
		a := newArtifactListAPI(t, "run-1", "run-2")

		body := listArtifacts(t, a, openapiv1.ListArtifactsParams{})

		require.Len(t, body.Items, 4)
		assert.Equal(t, []string{"run-2", "run-2", "run-1", "run-1"}, itemRunIDs(body))
		assert.Equal(t, "reporter", body.Items[0].Name)
		assert.Equal(t, "out.txt", body.Items[0].Path)
		assert.Equal(t, int64(len("hello")), body.Items[0].Size)
		require.NotNil(t, body.Items[0].StartedAt)
		assert.Nil(t, body.NextCursor)
	})

	t.Run("PagesWithCursor", func(t *testing.T) {
		a := newArtifactListAPI(t, "run-1", "run-2")
		limit := 1

		var seen []string
		params := openapiv1.ListArtifactsParams{Limit: &limit}
		for {
			page := listArtifacts(t, a, params)
			require.Len(t, page.Items, 1)
			seen = append(seen, page.Items[0].DagRunId+"/"+page.Items[0].Path)
			if page.NextCursor == nil {
				break
			}
			params.Cursor = page.NextCursor
		}

		assert.Equal(t, []string{
			"run-2/out.txt", "run-2/reports/summary.md",
			"run-1/out.txt", "run-1/reports/summary.md",
		}, seen)
	})

	t.Run("RejectsMalformedCursor", func(t *testing.T) {
		a := newArtifactListAPI(t, "run-1")
		cursor := "!!!"

		_, err := a.ListArtifacts(context.Background(),
			openapiv1.ListArtifactsRequestObject{Params: openapiv1.ListArtifactsParams{Cursor: &cursor}})

		var apiErr *Error
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	})

	t.Run("FiltersByFileName", func(t *testing.T) {
		a := newArtifactListAPI(t, "run-1")
		fileName := "summary"

		body := listArtifacts(t, a, openapiv1.ListArtifactsParams{FileName: &fileName})

		require.Len(t, body.Items, 1)
		assert.Equal(t, "reports/summary.md", body.Items[0].Path)
	})

	// A malformed glob is reported rather than silently matching nothing.
	t.Run("RejectsInvalidGlob", func(t *testing.T) {
		a := newArtifactListAPI(t, "run-1")
		fileName := "reports/["

		_, err := a.ListArtifacts(context.Background(),
			openapiv1.ListArtifactsRequestObject{Params: openapiv1.ListArtifactsParams{FileName: &fileName}})

		var apiErr *Error
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	})

	// A deployment that never enabled artifacts has no repository wired.
	t.Run("EmptyWithoutRepository", func(t *testing.T) {
		body := listArtifacts(t, &API{}, openapiv1.ListArtifactsParams{})

		assert.Empty(t, body.Items)
	})
}
