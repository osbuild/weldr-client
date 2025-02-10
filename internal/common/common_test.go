package common

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetContentFilename(t *testing.T) {
	filename, err := GetContentFilename("attachment; filename=875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar")
	assert.Nil(t, err)
	assert.Equal(t, "875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar", filename)
	filename, err = GetContentFilename("attachment; filename=875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar; donuts=glazed;")
	assert.Nil(t, err)
	assert.Equal(t, "875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar", filename)
	filename, err = GetContentFilename("attachment; filename=875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar ; ")
	assert.Nil(t, err)
	assert.Equal(t, "875759ea-1dbe-4f2c-9c8c-27cb8c7687ac-logs.tar", filename)
}

func TestGetContentFilenameError(t *testing.T) {
	_, err := GetContentFilename("attachment; filename=../../")
	assert.NotNil(t, err)
	_, err = GetContentFilename("")
	assert.NotNil(t, err)
	_, err = GetContentFilename("attachment;")
	assert.NotNil(t, err)
	_, err = GetContentFilename("attachment; filename=;")
	assert.NotNil(t, err)
	_, err = GetContentFilename("attachment; filename=. ;")
	assert.NotNil(t, err)
}

func TestMoveFile(t *testing.T) {
	dir := t.TempDir()

	f, err := os.CreateTemp("", "test-move-file-*")
	require.Nil(t, err)
	_, err = f.Write([]byte("This is just a test file\n"))
	require.Nil(t, err)
	f.Close()

	dstFile := fmt.Sprintf("%s/dest-file.txt", dir)
	err = MoveFile(f.Name(), dstFile)
	require.Nil(t, err)
	_, err = os.Stat(dstFile)
	require.Nil(t, err)
}

func TestMoveFileError(t *testing.T) {
	err := MoveFile("/tmp/no-such-testfile", "/tmp/no-such-destfile")
	require.NotNil(t, err)
}

func TestAppendQuery(t *testing.T) {
	assert.Equal(t, "/route/to/moes?bus=1", AppendQuery("/route/to/moes", "bus=1"))
	assert.Equal(t, "/route/to/moes?bus=0&taxi=1", AppendQuery("/route/to/moes?bus=0", "taxi=1"))
}

func TestCheckSocket(t *testing.T) {
	// Test with missing file
	err := CheckSocketError("/run/missing/file.socket", nil)
	assert.ErrorContains(t, err, "Check to make sure that")

	// Test with existing file
	f, err := os.CreateTemp("", "test-CheckSocket-*")
	require.Nil(t, err)
	_, err = f.Write([]byte("This is just a test file\n"))
	require.Nil(t, err)
	f.Close()

	err = CheckSocketError(f.Name(), nil)
	assert.Nil(t, err)

	err = CheckSocketError(f.Name(), fmt.Errorf("test error"))
	assert.ErrorContains(t, err, "test error")

	// NOTE: Cannot test permissons. root has access, and user cannot change them
	// to something it isn't allowed to access.
}
