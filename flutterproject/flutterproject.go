package flutterproject

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/v2/fileutil"
	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/godrei/go-flutter/flutterproject/internal/sdk"
	"gopkg.in/yaml.v3"
)

type FlutterAndDartSDKVersions struct {
	FlutterSDKVersions []sdk.Version
	DartSDKVersions    []sdk.Version
}

type Pubspec struct {
	Name string `yaml:"name"`
}

type Project struct {
	rootDir    string
	pubspecPth string

	fileManager fileutil.FileManager
	pathChecker pathutil.PathChecker
	fileOpener  FileOpener

	flutterAndDartSDKVersions *FlutterAndDartSDKVersions
	pubspec                   *Pubspec
	testDirPth                *string
	iosProjectPth             *string
	androidProjectPth         *string
}

func New(rootDir string, fileManager fileutil.FileManager, pathChecker pathutil.PathChecker) (*Project, error) {
	pubspecPth := filepath.Join(rootDir, sdk.PubspecRelPath)
	exists, err := pathChecker.IsPathExists(pubspecPth)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("not a Flutter project: pubspec.yaml not found at: %s", pubspecPth)
	}

	return &Project{
		rootDir:     rootDir,
		pubspecPth:  pubspecPth,
		fileManager: fileManager,
		pathChecker: pathChecker,
		fileOpener:  NewFileOpener(fileManager),
	}, nil
}

func (p *Project) Pubspec() (Pubspec, error) {
	if p.pubspec != nil {
		return *p.pubspec, nil
	}

	var pubspec Pubspec
	pubspecFile, err := p.fileManager.Open(p.pubspecPth)
	if err != nil {
		return Pubspec{}, err
	}
	if err := yaml.NewDecoder(pubspecFile).Decode(&pubspec); err != nil {
		return Pubspec{}, err
	}

	p.pubspec = &pubspec

	return pubspec, nil
}

func (p *Project) TestDirPth() string {
	if p.testDirPth != nil {
		return *p.testDirPth
	}

	const testDirRelPth = "test"

	hasTests := false
	testsDirPath := filepath.Join(p.rootDir, testDirRelPth)

	if exists, err := p.pathChecker.IsDirExists(testsDirPath); err == nil && exists {
		// TODO: make os.ReadDir testable
		if entries, err := os.ReadDir(testsDirPath); err == nil && len(entries) > 0 {
			for _, entry := range entries {
				if strings.HasSuffix(entry.Name(), "_test.dart") {
					hasTests = true
					break
				}
			}
		}
	}

	if !hasTests {
		testsDirPath = ""
	}

	p.testDirPth = &testsDirPath

	return testsDirPath
}

func (p *Project) IOSProjectPth() string {
	if p.iosProjectPth != nil {
		return *p.iosProjectPth
	}

	const iosProjectRelPth = "ios/Runner.xcworkspace"

	hasIOSProject := false
	iosProjectPth := filepath.Join(p.rootDir, iosProjectRelPth)
	if exists, err := p.pathChecker.IsPathExists(iosProjectPth); err == nil && exists {
		hasIOSProject = true
	}

	if !hasIOSProject {
		iosProjectPth = ""
	}

	p.iosProjectPth = &iosProjectPth

	return iosProjectPth

}

func (p *Project) AndroidProjectPth() string {
	const androidProjectRelPth = "android/build.gradle"
	const androidProjectKtsRelPth = "android/build.gradle.kts"

	hasAndroidProject := false
	androidProjectPth := filepath.Join(p.rootDir, androidProjectRelPth)
	if exists, err := p.pathChecker.IsPathExists(androidProjectPth); err == nil && exists {
		hasAndroidProject = true
	}

	if !hasAndroidProject {
		androidProjectPth = filepath.Join(p.rootDir, androidProjectKtsRelPth)
		if exists, err := p.pathChecker.IsPathExists(androidProjectPth); err == nil && exists {
			hasAndroidProject = true
		}
	}

	if !hasAndroidProject {
		androidProjectPth = ""
	}

	p.androidProjectPth = &androidProjectPth

	return androidProjectPth
}

func (p *Project) FlutterAndDartSDKVersions() (FlutterAndDartSDKVersions, error) {
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
