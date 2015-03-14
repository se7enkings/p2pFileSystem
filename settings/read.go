package settings

// TODO

const BlockSize = 4194304 // 4 M

func GetSharePath() string {
    return "D:/Audio/Music/Animation/星之声"
}

func IsIgnoredDir(name string) bool {
    return name == ".dropbox.cache"
}

func IsIgnoredFile(name string) bool{
    return name == "Thumbs.db"
}
