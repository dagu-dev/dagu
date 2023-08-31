package scheduler

import (
	"github.com/dagu-dev/dagu/internal/dag"
	"github.com/dagu-dev/dagu/service/scheduler/entry"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dagu-dev/dagu/internal/config"
	"github.com/dagu-dev/dagu/internal/utils"
)

var (
	testHomeDir string
)

func TestMain(m *testing.M) {
	tempDir := utils.MustTempDir("runner_test")
	changeHomeDir(tempDir)
	testHomeDir = tempDir
	code := m.Run()
	_ = os.RemoveAll(tempDir)
	os.Exit(code)
}

func changeHomeDir(homeDir string) {
	_ = os.Setenv("HOME", homeDir)
	_ = config.LoadConfig(homeDir)
}

func TestRun(t *testing.T) {
	now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	utils.FixedTime = now

	er := &mockEntryReader{
		Entries: []*entry.Entry{
			{
				Job:  &mockJob{},
				Next: now,
			},
			{
				Job:  &mockJob{},
				Next: now.Add(time.Minute),
			},
		},
	}

	r := New(Params{
		EntryReader: er,
		Config: &config.Config{
			LogDir: testHomeDir,
		},
	})

	go func() {
		_ = r.Start()
	}()

	time.Sleep(time.Second + time.Millisecond*100)
	r.Stop()

	require.Equal(t, 1, er.Entries[0].Job.(*mockJob).RunCount)
	require.Equal(t, 0, er.Entries[1].Job.(*mockJob).RunCount)
}

func TestRestart(t *testing.T) {
	now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	utils.FixedTime = now

	er := &mockEntryReader{
		Entries: []*entry.Entry{
			{
				EntryType: entry.EntryTypeRestart,
				Job:       &mockJob{},
				Next:      now,
			},
		},
	}

	r := New(Params{
		EntryReader: er,
		Config: &config.Config{
			LogDir: testHomeDir,
		},
	})

	go func() {
		_ = r.Start()
	}()

	time.Sleep(time.Second + time.Millisecond*100)
	require.Equal(t, 1, er.Entries[0].Job.(*mockJob).RestartCount)
}

func TestNextTick(t *testing.T) {
	n := time.Date(2020, 1, 1, 1, 0, 50, 0, time.UTC)
	utils.FixedTime = n
	r := New(Params{EntryReader: &mockEntryReader{}})
	next := r.nextTick(n)
	require.Equal(t, time.Date(2020, 1, 1, 1, 1, 0, 0, time.UTC), next)
}

type mockEntryReader struct {
	Entries []*entry.Entry
}

var _ EntryReader = (*mockEntryReader)(nil)

func (er *mockEntryReader) Read(_ time.Time) ([]*entry.Entry, error) {
	return er.Entries, nil
}

// TODO: fix to use mock library
type mockJob struct {
	Name         string
	RunCount     int
	StopCount    int
	RestartCount int
	Panic        error
}

var _ entry.Job = (*mockJob)(nil)

func (j *mockJob) GetDAG() *dag.DAG {
	return nil
}

func (j *mockJob) String() string {
	return j.Name
}

func (j *mockJob) Start() error {
	j.RunCount++
	if j.Panic != nil {
		panic(j.Panic)
	}
	return nil
}

func (j *mockJob) Stop() error {
	j.StopCount++
	return nil
}

func (j *mockJob) Restart() error {
	j.RestartCount++
	return nil
}
