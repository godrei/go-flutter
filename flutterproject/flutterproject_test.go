package flutterproject

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/godrei/go-flutter/flutterproject/internal/testassets"
	"github.com/godrei/go-flutter/mocks"
	"github.com/stretchr/testify/require"
)

func TestProject_FlutterAndDartSDKVersions(t *testing.T) {
	fileOpener := new(mocks.FileOpener)
	fileOpener.On("OpenFile", ".fvm/fvm_config.json").Return(strings.NewReader(testassets.FVMConfigJSON), nil)
	fileOpener.On("OpenFile", ".tool-versions").Return(strings.NewReader(testassets.ToolVersions), nil)
	fileOpener.On("OpenFile", "pubspec.lock").Return(strings.NewReader(testassets.PubspecLock), nil)
	fileOpener.On("OpenFile", "pubspec.yaml").Return(strings.NewReader(testassets.PubspecYaml), nil)

	proj := New("", fileOpener)
	sdkVersions, err := proj.FlutterAndDartSDKVersions()
	require.NoError(t, err)

	b, err := json.MarshalIndent(sdkVersions, "", "\t")
	require.NoError(t, err)

	require.Equal(t, string(b), `{
	"FlutterSDKVersions": [
		{
			"Version": "3.7.12",
			"Constraint": null,
			"Source": "fvm_config_json"
		},
		{
			"Version": "3.7.12",
			"Constraint": null,
			"Source": "tool_versions"
		},
		{
			"Version": null,
			"Constraint": "^3.7.12",
			"Source": "pubspec_yaml"
		}
	],
	"DartSDKVersions": [
		{
			"Version": null,
			"Constraint": "\u003e=2.19.6 \u003c3.0.0",
			"Source": "pubspec_lock"
		},
		{
			"Version": null,
			"Constraint": "\u003e=2.19.6 \u003c3.0.0",
			"Source": "pubspec_yaml"
		}
	]
}`)
}