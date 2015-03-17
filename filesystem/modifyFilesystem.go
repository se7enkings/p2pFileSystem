package filesystem

func AppendFilesystem(originFileSystem Filesystem, receivedFileSystem Filesystem) Filesystem {
	// because the File.IsLocal is ignored by json, IsLocal in received Filesystem is always default bool value(false).
	// It is possible that duplicate filename exists in the returned Filesystem,
	for hash, file := range receivedFileSystem {
		_, ok := originFileSystem[hash]
		if ok {
			continue
		} else {
			originFileSystem[hash] = file
		}
	}
	return originFileSystem
}

func OnReceiveFilesystem() {

}
