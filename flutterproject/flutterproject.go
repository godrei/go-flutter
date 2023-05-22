package flutterproject

import "github.com/godrei/go-flutter/flutterproject/internal/sdk"

type FlutterAndDartSDKVersions struct {
	FlutterSDKVersions []sdk.Version
	DartSDKVersions    []sdk.Version
}

type Project struct {
	rootDir    string
	fileOpener FileOpener

	flutterAndDartSDKVersions *FlutterAndDartSDKVersions
}

func New(rootDir string, fileOpener FileOpener) Project {
	return Project{
		rootDir:    rootDir,
		fileOpener: fileOpener,
	}
}

func (p Project) FlutterAndDartSDKVersions() (FlutterAndDartSDKVersions, error) {
	if p.flutterAndDartSDKVersions != nil {
		return *p.flutterAndDartSDKVersions, nil
	}

	versionReaders := []sdk.VersionsReader{
		sdk.NewFVMVersionReader(p.fileOpener),
		sdk.NewASDFVersionReader(p.fileOpener),
		sdk.NewPubspecLockVersionReader(p.fileOpener),
		sdk.NewPubspecVersionReader(p.fileOpener),
	}

	var flutterSDKVersions []sdk.Version
	var dartSDKVersions []sdk.Version
	for _, versionReader := range versionReaders {
		flutterSDKVersion, dartSDKVersion, err := versionReader.ReadSDKVersions(p.rootDir)
		if err != nil {
			return FlutterAndDartSDKVersions{}, err
		}
		if flutterSDKVersion != nil {
			flutterSDKVersions = append(flutterSDKVersions, *flutterSDKVersion)
		}
		if dartSDKVersion != nil {
			dartSDKVersions = append(dartSDKVersions, *dartSDKVersion)
		}
	}

	flutterAndDartSDKVersions := FlutterAndDartSDKVersions{
		FlutterSDKVersions: flutterSDKVersions,
		DartSDKVersions:    dartSDKVersions,
	}

	p.flutterAndDartSDKVersions = &flutterAndDartSDKVersions

	return flutterAndDartSDKVersions, nil
}
