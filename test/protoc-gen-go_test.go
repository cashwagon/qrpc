package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var protocGenGoTests = []struct {
	name   string
	args   []string
	result string
	golden string
}{
	{
		"Proto",
		[]string{
			"-I=./test/proto",
			"--go_out=paths=source_relative:./test/result",
			"./test/proto/test_api.proto",
		},
		"test_api.pb.go",
		"test_api.pb.proto.go",
	},
	{
		"GRPC",
		[]string{
			"-I=./test/proto",
			"--go_out=plugins=grpc,paths=source_relative:./test/result",
			"./test/proto/test_api.proto",
		},
		"test_api.pb.go",
		"test_api.pb.grpc.go",
	},
	{
		"QRPC",
		[]string{
			"-I=./test/proto",
			"--go_out=plugins=qrpc,paths=source_relative:./test/result",
			"./test/proto/test_api.proto",
		},
		"test_api.pb.go",
		"test_api.pb.qrpc.go",
	},
}

func Test_ProtocGenGo(t *testing.T) {
	protocBin := getProtoc(t)
	projectPath, err := filepath.Abs("..")
	require.NoError(t, err, "Cannot get project path")
	binPath := filepath.Join(projectPath, "bin")

	for _, tt := range protocGenGoTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			os.Mkdir("./result", os.ModeDir|os.ModePerm|755)
			defer func() {
				err := os.RemoveAll("./result")
				require.NoError(t, err)
			}()

			golden := loadGolden(t, tt.golden)

			pathEnv := os.Getenv("PATH")

			cmd := exec.Command(protocBin, tt.args...)
			cmd.Dir = projectPath
			cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s:%s", binPath, pathEnv))

			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Cannot execute protoc '%s':\n%s", protocBin, output)

			result := loadResult(t, tt.result)
			assert.Equal(t, string(golden), string(result))
		})
	}
}
