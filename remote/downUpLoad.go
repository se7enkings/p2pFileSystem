package remote

import (
	"fmt"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"github.com/CRVV/p2pFileSystem/transfer"
	"os"
)

func DownloadFile(hash string) error {
	name, path, size, owners, err := filesystem.GetRemoteFile(hash)
	if err != nil {
		return err
	}
	//	owners := make(map[string]int)
	//	for _, ownerName := range ownerList {
	//		owners[ownerName] = 0
	//	}
	blockCount, lastBlockSize := getBlockCountAndLastBlockSize(size)

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

	completeBlockNumChan := make(chan [2]int32)
	requestMessageChan := make(chan *FBRMessage, settings.MaxDownloadThreads*2)

	ownerNum := 0
	for i := 0; i < settings.MaxDownloadThreads; i++ {
		go downloadFileBlock(tempFile, owners[ownerNum], requestMessageChan, completeBlockNumChan)
		ownerNum++
		if ownerNum == len(owners) {
			ownerNum = 0
		}
	}
	requestMessage := FBRMessage{
		DestinationName: "", // assigned in function downloadFileBlock()
		Username:        settings.GetSettings().GetUsername(),
		FileHash:        hash,
		BlockNum:        -1,                     // assigned in function getMessageForQueue()
		BlockSize:       settings.FileBlockSize, // may be changed in function getMessageForQueue()
	}
	go func() {
		for blockNum, _ := range toBeCompletedBlocks {
			requestMessageChan <- getMessageForQueue(requestMessage, blockNum, blockCount, lastBlockSize)
		}
	}()

	for len(toBeCompletedBlocks) != 0 {
		blockComplete := <-completeBlockNumChan
		if blockComplete[1] < 0 {
			logger.Info(fmt.Sprintf("block %d downloading failed. add it to queue", blockComplete[0]))
			requestMessageChan <- getMessageForQueue(requestMessage, blockComplete[0], blockCount, lastBlockSize)
			continue
		}
		logger.Info(fmt.Sprintf("block %d complete", blockComplete[0]))
		delete(toBeCompletedBlocks, blockComplete[0])
	}
	for i := 0; i < settings.MaxDownloadThreads; i++ {
		requestMessageChan <- nil
	}

	//	tempFile.Sync()
	tempFile.Close()
	os.MkdirAll(settings.GetSettings().GetSharePath()+path, 0774)
	err = os.Rename(settings.GetSettings().GetSharePath()+"/.temp/"+hash, settings.GetSettings().GetSharePath()+path+"/"+name)
	logger.Warning(err)

	filesystem.RefreshLocalFile()
	return nil
}
func getMessageForQueue(message FBRMessage, blockNum int32, blockCount int32, lastBlockSize int32) *FBRMessage {
	message.BlockNum = blockNum
	if blockNum == blockCount-1 {
		message.BlockSize = lastBlockSize
	}
	return &message
}
func getBlockCountAndLastBlockSize(size int64) (blockCount int32, lastBlockSize int32) {
	lastBlockSize = int32(size % settings.FileBlockSize)
	if lastBlockSize == 0 {
		blockCount = int32(size / settings.FileBlockSize)
		lastBlockSize = settings.FileBlockSize
	} else {
		blockCount = int32(size/settings.FileBlockSize + 1)
	}
	logger.Info(fmt.Sprintf("this file have %d blocks", blockCount))
	return
}
func downloadFileBlock(tempFile *os.File, destinationName string, requestMessageChan chan *FBRMessage, completeBlockNumChan chan [2]int32) {
	for {
		requestMessage := <-requestMessageChan
		if requestMessage == nil {
			logger.Info(fmt.Sprintf("thread for download from %s down", destinationName))
			break
		}
		logger.Info(fmt.Sprintf("start to download block %d", requestMessage.BlockNum))
		requestMessage.DestinationName = destinationName
		fileData, err := transfer.TcpConnectionForReceiveFile(requestMessage)
		if err != nil {
			logger.Warning(err)
			completeBlockNumChan <- [2]int32{requestMessage.BlockNum, -1}
			return
		}
		offset := int64(requestMessage.BlockNum) * settings.FileBlockSize
		_, err = tempFile.WriteAt(fileData, offset)
		logger.Info(fmt.Sprintf("write to file at %d, size %d", offset, len(fileData)))
		if err != nil {
			logger.Warning(err)
			completeBlockNumChan <- [2]int32{requestMessage.BlockNum, -1}
			return
		}
		completeBlockNumChan <- [2]int32{requestMessage.BlockNum, 0}
		logger.Info(fmt.Sprintf("download block %d complete", requestMessage.BlockNum))
	}
}
func onRequestedFileBlock(requestMessage *FBRMessage) []byte {
	name, path, _, _, err := filesystem.GetLocalFile(requestMessage.FileHash)
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
		logger.Warning(fmt.Sprintf("I am requested a too big block and its size is %d bytes", requestMessage.BlockSize))
		return nil
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
