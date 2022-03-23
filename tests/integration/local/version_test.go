package tests

import (
	"strings"
	"testing"

	"get.porter.sh/plugin/kubernetes/pkg"
	"get.porter.sh/porter/pkg/porter/version"
	"get.porter.sh/porter/pkg/printer"
	"github.com/stretchr/testify/require"
)

func TestPrintVersion(t *testing.T) {
	pkg.Commit = "abc123"
	pkg.Version = "v1.2.3"

	m := NewTestPlugin(t)

	opts := version.Options{}
	err := opts.Validate()
	require.NoError(t, err)
	err = m.Plugin.PrintVersion(opts)
	require.NoError(t, err)

	gotOutput := m.TestContext.GetOutput()
	wantOutput := "kubernetes v1.2.3 (abc123) by Porter Authors"
	if !strings.Contains(gotOutput, wantOutput) {
		t.Fatalf("invalid output:\nWANT:\t%q\nGOT:\t%q\n", wantOutput, gotOutput)
	}
}

func TestPrintJsonVersion(t *testing.T) {
	pkg.Commit = "abc123"
	pkg.Version = "v1.2.3"

	m := NewTestPlugin(t)

	opts := version.Options{}
	opts.RawFormat = string(printer.FormatJson)
	err := opts.Validate()
	require.NoError(t, err)
	err = m.PrintVersion(opts)
	require.NoError(t, err)

	gotOutput := m.TestContext.GetOutput()
	wantOutput := `{
  "name": "kubernetes",
  "version": "v1.2.3",
  "commit": "abc123",
  "author": "Porter Authors",
  "implementations": [
    {
      "type": "secrets",
      "implementation": "secrets"
    }
  ]
}
`
	if !strings.Contains(gotOutput, wantOutput) {
		t.Fatalf("invalid output:\nWANT:\t%q\nGOT:\t%q\n", wantOutput, gotOutput)
	}
}
