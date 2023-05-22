package sdk

import (
	"io"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const pubspecLockRelPath = "pubspec.lock"

type PubspecLockVersionReader struct {
	fileOpener FileOpener
}

func NewPubspecLockVersionReader(fileOpener FileOpener) PubspecLockVersionReader {
	return PubspecLockVersionReader{
		fileOpener: fileOpener,
	}
}

func (r PubspecLockVersionReader) ReadSDKVersions(projectRootDir string) (*Version, *Version, error) {
	pubspecLockPth := filepath.Join(projectRootDir, pubspecLockRelPath)
	f, err := r.fileOpener.OpenFile(pubspecLockPth)
	if err != nil {
		return nil, nil, err
	}

	if f == nil {
		return nil, nil, nil
	}

	flutterVersionStr, dartVersionStr, err := parsePubspecLockSDKVersions(f)
	if err != nil {
		return nil, nil, err
	}

	var flutterVersion *Version
	if flutterVersionStr != "" {
		flutterVersion, err = NewVersion(flutterVersionStr, PubspecLockVersionSource)
		if err != nil {
			return nil, nil, err
		}
	}

	var dartVersion *Version
	if dartVersionStr != "" {
		dartVersion, err = NewVersion(dartVersionStr, PubspecLockVersionSource)
		if err != nil {
			return nil, nil, err
		}
	}

	return flutterVersion, dartVersion, nil
}

func parsePubspecLockSDKVersions(pubspecLockReader io.Reader) (string, string, error) {
	type pubspecLock struct {
		SDKs struct {
			Dart    string `yaml:"dart"`
			Flutter string `yaml:"flutter"`
		} `yaml:"sdks"`
	}

	var config pubspecLock
	d := yaml.NewDecoder(pubspecLockReader)
	if err := d.Decode(&config); err != nil {
		return "", "", err
	}

	return config.SDKs.Flutter, config.SDKs.Dart, nil
}
