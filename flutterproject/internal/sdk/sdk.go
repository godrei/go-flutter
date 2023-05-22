package sdk

import "io"

type VersionsReader interface {
	ReadSDKVersions(projectRootDir string) (*Version, *Version, error)
}

type FileOpener interface {
	OpenFile(pth string) (io.Reader, error)
}
