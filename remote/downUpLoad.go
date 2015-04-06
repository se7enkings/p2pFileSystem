package remote

import (
	"fmt"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"github.com/CRVV/p2pFileSystem/transfer"
	"os"
	"strconv"
)

func DownloadFile(hash string) error {
	name, path, size, owners, err := filesystem.GetRemoteFile(hash)
	if err != nil {
		return err
	}

	var blockCount int32
	lastBlockSize := int32(size % settings.FileBlockSize)
	if lastBlockSize == 0 {
		blockCount = int32(size / settings.FileBlockSize)
	} else {
		blockCount = int32(size/settings.FileBlockSize + 1)
	}
	logger.Info("this file have " + strconv.Itoa(int(blockCount)) + " blocks")
	os.Mkdir(settings.GetSettings().GetSharePath()+"/.temp", 0774)
	tempFile, err := os.Create(settings.GetSettings().GetSharePath() + "/.temp/" + hash)
	if err != nil {
		logger.Warning(err)
		return err
	}
	defer tempFile.Close()

	toBeCompletedBlocks := make(map[int32]bool)
	for i := int32(0); i < blockCount; i++ {
		toBeCompletedBlocks[i] = true
	}

	completeBlockNumChan := make(chan int32)
	ownerCount := len(owners)
	for {
		ownerNum := 0
		runningRoutines := 0
		for blockNum, _ := range toBeCompletedBlocks {
			requestMessage := FBRMessage{
				DestinationName: owners[ownerNum],
				Username:        settings.GetSettings().GetUsername(),
				FileHash:        hash,
				BlockNum:        blockNum}
			requestMessage.BlockSize = int32(settings.FileBlockSize)
			if blockNum == blockCount-1 && lastBlockSize != 0 {
				requestMessage.BlockSize = lastBlockSize
			}

			//TODO: need thread pool
			go downloadFileBlock(tempFile, &requestMessage, completeBlockNumChan)
			runningRoutines++

			ownerNum++
			if ownerNum == ownerCount {
				ownerNum = 0
			}
		}
		if runningRoutines == 0 {
			break
		}
		for i := 0; i < runningRoutines; i++ {
			blockNum := <-completeBlockNumChan
			if blockNum < 0 {
				logger.Info("block number " + strconv.Itoa(int(blockNum)) + " downloaded but fail")
				continue
			}
			logger.Info("block number " + strconv.Itoa(int(blockNum)) + " complete")
			delete(toBeCompletedBlocks, blockNum)
		}
	}
	tempFile.Sync()
	tempFile.Close()
	os.MkdirAll(settings.GetSettings().GetSharePath()+path, 0774)
	err = os.Rename(settings.GetSettings().GetSharePath()+"/.temp/"+hash, settings.GetSettings().GetSharePath()+path+"/"+name)
	logger.Warning(err)

	filesystem.RefreshLocalFile()
	return nil
}
func downloadFileBlock(tempFile *os.File, requestMessage *FBRMessage, completeBlockNumChan chan int32) {
	logger.Info("start to download block " + strconv.Itoa(int(requestMessage.BlockNum)))
	fileData, err := transfer.TcpConnectionForReceiveFile(requestMessage)
	if err != nil {
		completeBlockNumChan <- -1
		return
	}
	offset := int64(requestMessage.BlockNum) * settings.FileBlockSize
	_, err = tempFile.WriteAt(fileData, offset)
	logger.Info(fmt.Sprintf("write to file at %d, size %d", offset, len(fileData)))
	if err != nil {
		completeBlockNumChan <- -1
		return
	}
	completeBlockNumChan <- requestMessage.BlockNum
	logger.Info("download block " + strconv.Itoa(int(requestMessage.BlockNum)) + " complete")
}
func onRequestedFileBlock(requestMessage *FBRMessage) []byte {
	name, path, _, _, err := filesystem.GetRemoteFile(requestMessage.FileHash)
	if err != nil {
		logger.Warning("I am requested a file but I do not have it")
		return nil
	}
	f, err := os.Open(settings.GetSettings().GetSharePath() + path + "/" + name)
	if err != nil {
		logger.Warning(err)
		return nil
	}
	if requestMessage.BlockSize > settings.FileBlockSize {
		//TODO: should not panic
		panic(requestMessage.BlockSize)
	}
	buff := make([]byte, requestMessage.BlockSize)
	_, err = f.ReadAt(buff, int64(requestMessage.BlockNum)*settings.FileBlockSize)
	if err != nil {
		logger.Warning(err)
		return nil
	}
	return buff
}
func UploadFile(path string) {

}
