package app

import "github.com/spf13/afero"

func initFileSystem(dir string) afero.Fs {
	var afs afero.Fs
	if dir == "" || dir == "memory" {
		afs = afero.NewMemMapFs()
	} else {
		afs = afero.NewBasePathFs(afs, dir)
	}

	return afs
}
