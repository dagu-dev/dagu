package utils_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yohamta/dagman/internal/utils"
)

func TestMustGetUserHomeDir(t *testing.T) {
	err := os.Setenv("HOME", "/test")
	if err != nil {
		t.Fatal(err)
	}
	hd := utils.MustGetUserHomeDir()
	assert.Equal(t, "/test", hd)
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

func TestOpenOrCreateFile(t *testing.T) {
	tmp, err := ioutil.TempDir("", "utils_test")
	require.NoError(t, err)
	name := path.Join(tmp, "/file_for_test.txt")
	f, err := utils.OpenOrCreateFile(name)
	require.NoError(t, err)
	defer func() {
		f.Close()
		os.Remove(name)
	}()
	require.True(t, utils.FileExists(name))
}
