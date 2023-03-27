package util

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// LoadTestSQLFile load test sql data from a file
func LoadTestSQLFile(t *testing.T, tx boil.ContextExecutor, filename string, args ...interface{}) {
	body, err := ReadFile(filename)
	require.NoError(t, err)

	_, err = tx.ExecContext(context.Background(), string(body), args...)
	require.NoError(t, err)
}

// ReadFile reads a file completely. But if the file does not exist, try to find it in the parent directory, [repeat...]
func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
