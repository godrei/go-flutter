package sdk

import (
	"io"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const pubspecRelPath = "pubspec.yaml"

type PubspecVersionReader struct {
	fileOpener FileOpener
}

func NewPubspecVersionReader(fileOpener FileOpener) PubspecVersionReader {
	return PubspecVersionReader{
		fileOpener: fileOpener,
	}
}

func (r PubspecVersionReader) ReadSDKVersions(projectRootDir string) (*Version, *Version, error) {
	pubspecPth := filepath.Join(projectRootDir, pubspecRelPath)
	f, err := r.fileOpener.OpenFile(pubspecPth)
	if err != nil {
		return nil, nil, err
	}

	if f == nil {
		return nil, nil, nil
	}

	flutterVersionStr, dartVersionStr, err := parsePubspecSDKVersions(f)
	if err != nil {
		return nil, nil, err
	}

	var flutterVersion *Version
	if flutterVersionStr != "" {
		flutterVersion, err = NewVersion(flutterVersionStr, PubspecVersionSource)
		if err != nil {
			return nil, nil, err
		}
	}

	var dartVersion *Version
	if dartVersionStr != "" {
		dartVersion, err = NewVersion(dartVersionStr, PubspecVersionSource)
		if err != nil {
			return nil, nil, err
		}
	}

	return flutterVersion, dartVersion, nil
}

func parsePubspecSDKVersions(pubspecReader io.Reader) (string, string, error) {
	type pubspec struct {
		Environment struct {
			Dart    string `yaml:"sdk"`
			Flutter string `yaml:"flutter"`
		} `yaml:"environment"`
	}

	var config pubspec
	d := yaml.NewDecoder(pubspecReader)
	if err := d.Decode(&config); err != nil {
		return "", "", err
	}

	return config.Environment.Flutter, config.Environment.Dart, nil
}