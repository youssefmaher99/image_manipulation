package util

import (
	"context"
	"fmt"
	"os"
	"path"
	"server/data"
	"server/presist"

	"github.com/mholt/archiver/v4"
)

func Archive(imageNames []string, uid string, filterName string) error {
	images := make(map[string]string)
	for i := 0; i < len(imageNames); i++ {
		key := "filtered/" + uid + "_" + filterName + "_" + imageNames[i]
		fmt.Println(key)
		images[key] = ""
	}
	// fmt.Println(images)
	archive, err := archiver.FilesFromDisk(nil, images)
	if err != nil {
		return err
	}

	archiverExt := ExtBasedOnPlatform()

	out, err := os.Create(path.Join("archives", uid+archiverExt))
	if err != nil {
		return err
	}
	defer out.Close()

	format := archiver.CompressedArchive{
		Archival: archiver.Tar{},
	}
	err = format.Archive(context.Background(), out, archive)
	if err != nil {
		return err
	}

	// HACK
	data.InMemoryArchives.Add(uid, struct{}{})
	presist.AddArchive(uid)
	return nil
}
