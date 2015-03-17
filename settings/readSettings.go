package settings

// TODO

const BlockSize = 4194304 // 4 M

func GetSharePath() string {
    return "test/testLocalFolder"
}

func IsIgnoredDir(name string) bool {
    return name == ".dropbox.cache"
}

func IsIgnoredFile(name string) bool{
    return name == "Thumbs.db"
}

func GetUserName() string {
    return "crvv.pku@gmail.com"
}

