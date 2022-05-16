package utils_test

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yohamta/dagu/internal/constants"
	"github.com/yohamta/dagu/internal/utils"
)

func TestMustGetUserHomeDir(t *testing.T) {
	err := os.Setenv("HOME", "/test")
	if err != nil {
		t.Fatal(err)
	}
	hd := utils.MustGetUserHomeDir()
	assert.Equal(t, "/test", hd)
}

func TestDefaultEnv(t *testing.T) {
	env := utils.DefaultEnv()
	require.Contains(t, env, "PATH")
}

func TestMustGetwd(t *testing.T) {
	wd, _ := os.Getwd()
	assert.Equal(t, utils.MustGetwd(), wd)
}

func TestFormatTime(t *testing.T) {
	tm := time.Date(2022, 2, 1, 2, 2, 2, 0, time.Now().Location())
	fomatted := utils.FormatTime(tm)
	assert.Equal(t, "2022-02-01 02:02:02", fomatted)

	parsed, err := utils.ParseTime(fomatted)
	require.NoError(t, err)
	assert.Equal(t, tm, parsed)

	require.Equal(t, constants.TimeEmpty, utils.FormatTime(time.Time{}))
	parsed, err = utils.ParseTime(constants.TimeEmpty)
	require.NoError(t, err)
	require.Equal(t, time.Time{}, parsed)
}

func TestFormatDuration(t *testing.T) {
	dr := time.Second*5 + time.Millisecond*100
	assert.Equal(t, "5.1s", utils.FormatDuration(dr, ""))
}

func TestSplitCommand(t *testing.T) {
	command := "ls -al test/"
	program, args := utils.SplitCommand(command)
	assert.Equal(t, "ls", program)
	assert.Equal(t, "-al", args[0])
	assert.Equal(t, "test/", args[1])
}

func TestFileExits(t *testing.T) {
	require.True(t, utils.FileExists("/"))
}

func TestValidFilename(t *testing.T) {
	f := utils.ValidFilename("file\\name", "_")
	assert.Equal(t, f, "file_name")
}

func TestOpenFile(t *testing.T) {
	tmp, err := os.MkdirTemp("", "open")
	require.NoError(t, err)

	name := path.Join(tmp, "/file.txt")
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	defer func() {
		f.Close()
		os.Remove(name)
	}()

	f.WriteString("test")
	f.Sync()
	f.Close()

	_, err = utils.OpenFile(name)
	require.NoError(t, err)
}

func TestOpenOrCreateFile(t *testing.T) {
	tmp, err := os.MkdirTemp("", "open_or_create")
	require.NoError(t, err)

	name := path.Join(tmp, "/file.txt")
	f, err := utils.OpenOrCreateFile(name)
	require.NoError(t, err)

	defer func() {
		f.Close()
		os.Remove(name)
	}()

	require.True(t, utils.FileExists(name))

	f.Close()
	os.Remove(name)

	_, err = utils.OpenFile(name)
	require.Error(t, err)
}

func TestParseVariable(t *testing.T) {
	os.Setenv("TEST_VAR", "test")
	r, err := utils.ParseVariable("${TEST_VAR}")
	require.NoError(t, err)
	assert.Equal(t, r, "test")

	_, err = utils.ParseVariable("`ech test`")
	require.Error(t, err)

	r, err = utils.ParseVariable("`echo test`")
	require.NoError(t, err)
	assert.Equal(t, r, "test")
}

func TestMustTempDir(t *testing.T) {
	dir := utils.MustTempDir("tempdir")
	defer os.RemoveAll(dir)
	require.Contains(t, dir, "tempdir")
}

func TestOpenfile(t *testing.T) {
	dir := utils.MustTempDir("tempdir")
	defer os.RemoveAll(dir)

	fn := path.Join(dir, "test.txt")
	f, err := utils.OpenOrCreateFile(fn)
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString("test")
	require.NoError(t, err)
	require.NoError(t, f.Sync(), err)
	require.NoError(t, f.Close(), err)
	require.True(t, utils.FileExists(fn))

	f2, err := os.Open(fn)
	require.NoError(t, err)
	defer f2.Close()
	b, err := io.ReadAll(f2)
	require.NoError(t, err)
	assert.Equal(t, "test", string(b))
}

func TestIgnoreErr(t *testing.T) {
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	log.SetOutput(w)

	defer func() {
		os.Stdout = origStdout
		log.SetOutput(origStdout)
	}()

	utils.LogIgnoreErr("test action", errors.New("test error"))
	os.Stdout = origStdout
	w.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)

	s := buf.String()
	require.Contains(t, s, "test action failed")
	require.Contains(t, s, "test error")
}

func TestTruncString(t *testing.T) {
	require.Equal(t, "12345678", utils.TruncString("123456789", 8))
	require.Equal(t, "1234567", utils.TruncString("1234567", 8))
}
