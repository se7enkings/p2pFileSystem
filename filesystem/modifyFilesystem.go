package filesystem

func AppendFilesystem(originFileSystem Filesystem, receivedFileSystem Filesystem) Filesystem {
	// because the File.IsLocal is ignored by json, IsLocal in received Filesystem is always default bool value(false).
    // TODO: handle duplication of name
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
