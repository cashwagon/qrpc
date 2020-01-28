package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func loadGolden(t *testing.T, name string) []byte {
	t.Helper()

	return loadFile(t, filepath.Join("golden", filepath.Clean(name)))
}

func loadResult(t *testing.T, name string) []byte {
	t.Helper()

	return loadFile(t, filepath.Join("result", filepath.Clean(name)))
}

func loadFile(t *testing.T, name string) []byte {
	t.Helper()

	file, err := os.Open(filepath.Clean(name))
	require.NoErrorf(t, err, "Cannot open file %s", name)
	defer file.Close() // nolint:errcheck

	data, err := ioutil.ReadAll(file)
	require.NoError(t, err, "Cannot read file %s", name)

	return data
}

func getProtoc(t *testing.T) string {
	t.Helper()

	binaryPath := os.Getenv("PROTOC")
	if binaryPath != "" {
		return binaryPath
	}

	return "protoc"
}

func getBinary(t *testing.T) string {
	t.Helper()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	binaryPath := os.Getenv("PROTOC_GEN_GO")
	if binaryPath != "" {
		return binaryPath
	}

	projectPath := filepath.Dir(cwd)
	return filepath.Join(projectPath, "bin", "protoc-gen-go")
}
