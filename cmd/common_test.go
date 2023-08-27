package cmd

import (
	"bytes"
	"github.com/yohamta/dagu/internal/config"
	"io"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/yohamta/dagu/internal/controller"
	"github.com/yohamta/dagu/internal/database"
	"github.com/yohamta/dagu/internal/scheduler"
	"github.com/yohamta/dagu/internal/utils"
)

func TestMain(m *testing.M) {
	tmpDir := utils.MustTempDir("dagu_test")
	changeHomeDir(tmpDir)
	code := m.Run()
	os.RemoveAll(tmpDir)
	os.Exit(code)
}

func changeHomeDir(homeDir string) {
	os.Setenv("HOME", homeDir)
	_ = config.LoadConfig(homeDir)
}

type cmdTest struct {
	args        []string
	expectedOut []string
}

func testRunCommand(t *testing.T, cmd *cobra.Command, test cmdTest) {
	t.Helper()

	root := &cobra.Command{Use: "root"}
	root.AddCommand(cmd)

	// Set arguments.
	root.SetArgs(test.args)

	// Run the command.
	out := withSpool(t, func() {
		err := root.Execute()
		require.NoError(t, err)
	})

	// Check outputs.
	for _, s := range test.expectedOut {
		require.Contains(t, out, s)
	}
}

func withSpool(t *testing.T, f func()) string {
	t.Helper()

	origStdout := os.Stdout

	r, w, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = w
	log.SetOutput(w)

	defer func() {
		os.Stdout = origStdout
		log.SetOutput(origStdout)
		w.Close()
	}()

	f()

	os.Stdout = origStdout
	w.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)

	return buf.String()
}

func testDAGFile(name string) string {
	d := path.Join(utils.MustGetwd(), "testdata")
	return path.Join(d, name)
}

func testStatusEventual(t *testing.T, dagFile string, expected scheduler.SchedulerStatus) {
	t.Helper()

	d, err := loadDAG(dagFile, "")
	require.NoError(t, err)
	ctrl := controller.NewDAGController(d)

	require.Eventually(t, func() bool {
		status, err := ctrl.GetStatus()
		require.NoError(t, err)
		return expected == status.Status
	}, time.Millisecond*5000, time.Millisecond*50)
}

func testLastStatusEventual(t *testing.T, dagFile string, expected scheduler.SchedulerStatus) {
	t.Helper()
	require.Eventually(t, func() bool {
		db := &database.Database{Config: database.DefaultConfig()}
		status := db.ReadStatusHist(dagFile, 1)
		if len(status) < 1 {
			return false
		}
		return expected == status[0].Status.Status
	}, time.Millisecond*5000, time.Millisecond*50)
}
