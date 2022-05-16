package database

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yohamta/dagu/internal/config"
	"github.com/yohamta/dagu/internal/models"
	"github.com/yohamta/dagu/internal/scheduler"
	"github.com/yohamta/dagu/internal/settings"
	"github.com/yohamta/dagu/internal/utils"
)

func TestDatabase(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T, db *Database,
	){
		"create new datafile":                 testNewDataFile,
		"write status to file and rename":     testWriteStatusToFile,
		"append status to existing file":      testWriteStatusToExistingFile,
		"write status and find files":         testWriteAndFindFiles,
		"write status and find by request id": testWriteAndFindByRequestId,
		"remove old files":                    testRemoveOldFiles,
		"test read latest status":             testReadLatestStatus,
		"test read latest n status":           testReadStatusN,
		"test compaction":                     testCompactFile,
		"test error read file":                testErrorReadFile,
		"test error parse file":               testErrorParseFile,
	} {
		t.Run(scenario, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "test-database")
			db := New(&Config{
				Dir: dir,
			})
			require.NoError(t, err)
			defer os.RemoveAll(dir)
			fn(t, db)
		})
	}
}

func testNewDataFile(t *testing.T, db *Database) {
	cfg := &config.Config{
		ConfigPath: "test_new_data_file.yaml",
	}
	timestamp := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)
	requestId := "request-id-1"
	f, err := db.newFile(cfg.ConfigPath, timestamp, requestId)
	require.NoError(t, err)
	p := utils.ValidFilename(strings.TrimSuffix(
		path.Base(cfg.ConfigPath), path.Ext(cfg.ConfigPath)), "_")
	assert.Regexp(t, fmt.Sprintf("%s.*/%s.20220101.00:00:00.000.%s.dat", p, p, requestId[:8]), f)

	_, err = db.newFile("", timestamp, requestId)
	require.Error(t, err)
}

func testWriteAndFindFiles(t *testing.T, db *Database) {
	cfg := &config.Config{
		Name:       "test_read_status_n",
		ConfigPath: "test_data_files_n.yaml",
	}
	defer db.RemoveAll(cfg.ConfigPath)

	for _, data := range []struct {
		Status    *models.Status
		RequestId string
		Timestamp time.Time
	}{
		{
			models.NewStatus(cfg, nil, scheduler.SchedulerStatus_None, 10000, nil, nil),
			"request-id-1",
			time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			models.NewStatus(cfg, nil, scheduler.SchedulerStatus_None, 10000, nil, nil),
			"request-id-2",
			time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local),
		},
		{
			models.NewStatus(cfg, nil, scheduler.SchedulerStatus_None, 10000, nil, nil),
			"request-id-3",
			time.Date(2022, 1, 3, 0, 0, 0, 0, time.Local),
		},
	} {
		status := data.Status
		status.RequestId = data.RequestId
		testWriteStatus(t, db, cfg, status, data.Timestamp)
	}

	files := db.latest(db.pattern(cfg.ConfigPath)+"*.dat", 2)
	require.Equal(t, 2, len(files))
}

func testWriteAndFindByRequestId(t *testing.T, db *Database) {
	cfg := &config.Config{
		Name:       "test_find_by_request_id",
		ConfigPath: "test_find_by_request_id.yaml",
	}
	defer db.RemoveAll(cfg.ConfigPath)

	for _, data := range []struct {
		Status    *models.Status
		RequestId string
		Timestamp time.Time
	}{
		{
			models.NewStatus(cfg, nil, scheduler.SchedulerStatus_None, 10000, nil, nil),
			"request-id-1",
			time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			models.NewStatus(cfg, nil, scheduler.SchedulerStatus_None, 10000, nil, nil),
			"request-id-2",
			time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local),
		},
		{
			models.NewStatus(cfg, nil, scheduler.SchedulerStatus_None, 10000, nil, nil),
			"request-id-3",
			time.Date(2022, 1, 3, 0, 0, 0, 0, time.Local),
		},
	} {
		status := data.Status
		status.RequestId = data.RequestId
		testWriteStatus(t, db, cfg, status, data.Timestamp)
	}

	status, err := db.FindByRequestId(cfg.ConfigPath, "request-id-2")
	require.NoError(t, err)
	assert.Equal(t, status.Status.RequestId, "request-id-2")

	status, err = db.FindByRequestId(cfg.ConfigPath, "request-id-10000")
	assert.Error(t, err)
	assert.Nil(t, status)
}

func testRemoveOldFiles(t *testing.T, db *Database) {
	cfg := &config.Config{
		ConfigPath: "test_remove_old.yaml",
	}

	for _, data := range []struct {
		Status    *models.Status
		RequestId string
		Timestamp time.Time
	}{
		{
			models.NewStatus(cfg, nil, scheduler.SchedulerStatus_None, 10000, nil, nil),
			"request-id-1",
			time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			models.NewStatus(cfg, nil, scheduler.SchedulerStatus_None, 10000, nil, nil),
			"request-id-2",
			time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local),
		},
		{
			models.NewStatus(cfg, nil, scheduler.SchedulerStatus_None, 10000, nil, nil),
			"request-id-3",
			time.Date(2022, 1, 3, 0, 0, 0, 0, time.Local),
		},
	} {
		status := data.Status
		status.RequestId = data.RequestId
		testWriteStatus(t, db, cfg, data.Status, data.Timestamp)
	}

	files := db.latest(db.pattern(cfg.ConfigPath)+"*.dat", 3)
	require.Equal(t, 3, len(files))

	db.RemoveOld(db.pattern(cfg.ConfigPath)+"*.dat", 0)

	files = db.latest(db.pattern(cfg.ConfigPath)+"*.dat", 3)
	require.Equal(t, 0, len(files))

	m := db.latest("invalid-pattern", 3)
	require.Equal(t, 0, len(m))
}

func testReadLatestStatus(t *testing.T, db *Database) {
	cfg := &config.Config{
		ConfigPath: "test_config_status_reader.yaml",
	}
	requestId := "request-id-1"
	dw, _, err := db.NewWriter(cfg.ConfigPath, time.Now(), requestId)
	require.NoError(t, err)
	err = dw.Open()
	require.NoError(t, err)
	defer dw.Close()

	status := models.NewStatus(cfg, nil, scheduler.SchedulerStatus_None, 10000, nil, nil)
	dw.Write(status)

	status.Status = scheduler.SchedulerStatus_Running
	status.Pid = 20000
	dw.Write(status)

	ret, err := db.ReadStatusToday(cfg.ConfigPath)

	require.NoError(t, err)
	require.NotNil(t, ret)
	assert.Equal(t, int(status.Pid), int(ret.Pid))
	require.Equal(t, status.Status, ret.Status)

}

func testReadStatusN(t *testing.T, db *Database) {
	cfg := &config.Config{
		Name:       "test_read_status_n",
		ConfigPath: "test_config_status_reader_hist.yaml",
	}

	for _, data := range []struct {
		Status    *models.Status
		RequestId string
		Timestamp time.Time
	}{
		{
			models.NewStatus(cfg, nil, scheduler.SchedulerStatus_None, 10000, nil, nil),
			"request-id-1",
			time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			models.NewStatus(cfg, nil, scheduler.SchedulerStatus_None, 10000, nil, nil),
			"request-id-2",
			time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local),
		},
		{
			models.NewStatus(cfg, nil, scheduler.SchedulerStatus_None, 10000, nil, nil),
			"request-id-3",
			time.Date(2022, 1, 3, 0, 0, 0, 0, time.Local),
		},
	} {
		status := data.Status
		status.RequestId = data.RequestId
		testWriteStatus(t, db, cfg, data.Status, data.Timestamp)
	}

	recordMax := 2

	ret := db.ReadStatusHist(cfg.ConfigPath, recordMax)

	require.Equal(t, recordMax, len(ret))
	assert.Equal(t, cfg.Name, ret[0].Status.Name)
	assert.Equal(t, cfg.Name, ret[1].Status.Name)
}

func testCompactFile(t *testing.T, db *Database) {
	cfg := &config.Config{
		Name:       "test_compact_file",
		ConfigPath: "test_compact_file.yaml",
	}
	requestId := "request-id-1"

	dw, _, err := db.NewWriter(cfg.ConfigPath, time.Now(), requestId)
	require.NoError(t, err)
	require.NoError(t, dw.Open())

	for _, data := range []struct {
		Status *models.Status
	}{
		{models.NewStatus(
			cfg, nil, scheduler.SchedulerStatus_Running, 10000, nil, nil)},
		{models.NewStatus(
			cfg, nil, scheduler.SchedulerStatus_Cancel, 10000, nil, nil)},
		{models.NewStatus(
			cfg, nil, scheduler.SchedulerStatus_Success, 10000, nil, nil)},
	} {
		require.NoError(t, dw.Write(data.Status))
	}

	dw.Close()

	var s *models.StatusFile = nil
	if h := db.ReadStatusHist(cfg.ConfigPath, 1); len(h) > 0 {
		s = h[0]
	}
	require.NotNil(t, s)

	db2 := New(db.Config)
	err = db2.Compact(cfg.ConfigPath, s.File)
	assert.False(t, utils.FileExists(s.File))
	require.NoError(t, err)

	var s2 *models.StatusFile = nil
	if h := db2.ReadStatusHist(cfg.ConfigPath, 1); len(h) > 0 {
		s2 = h[0]
	}
	require.NotNil(t, s2)

	assert.Regexp(t, `test_compact_file.*_c.dat`, s2.File)
	assert.Equal(t, s.Status, s2.Status)

	err = db2.Compact(cfg.ConfigPath, "Invalid_file_name.dat")
	require.Error(t, err)
}

func testErrorReadFile(t *testing.T, db *Database) {
	_, err := ParseFile("invalid_file.dat")
	require.Error(t, err)

	_, _, err = db.NewWriter("", time.Now(), "")
	require.Error(t, err)

	_, err = db.ReadStatusToday("invalid_file.yaml")
	require.Error(t, err)

	_, err = db.FindByRequestId("invalid_file.yaml", "invalid_id")
	require.Error(t, err)
}

func testErrorParseFile(t *testing.T, db *Database) {
	tmpDir := utils.MustTempDir("test_error_parse_file")
	tmpFile := filepath.Join(tmpDir, "test_error_parse_file.dat")

	_, err := ParseFile(tmpFile)
	require.Error(t, err)

	f, err := utils.OpenOrCreateFile(tmpFile)
	require.NoError(t, err)

	_, err = ParseFile(tmpFile)
	require.Error(t, err)

	_, err = f.WriteString("invalid json")
	require.NoError(t, err)

	_, err = ParseFile(tmpFile)
	require.Error(t, err)

	_, err = f.WriteString("\n{}")
	require.NoError(t, err)

	_, err = ParseFile(tmpFile)
	require.NoError(t, err)
}

func testWriteStatus(t *testing.T, db *Database, cfg *config.Config, status *models.Status, tm time.Time) {
	t.Helper()
	dw, _, err := db.NewWriter(cfg.ConfigPath, tm, status.RequestId)
	require.NoError(t, err)
	require.NoError(t, dw.Open())
	defer dw.Close()
	require.NoError(t, dw.Write(status))
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.Equal(t, cfg.Dir, settings.MustGet(settings.CONFIG__DATA_DIR))
}

func TestTimestamp(t *testing.T) {
	for _, tt := range []struct {
		Name string
		Want string
	}{
		{Name: "test_timestamp.20200101.10:00:00.dat", Want: "20200101.10:00:00"},
		{Name: "test_timestamp.20200101.12:34:56_c.dat", Want: "20200101.12:34:56"},
	} {
		assert.Equal(t, tt.Want, timestamp(tt.Name))
	}
}
