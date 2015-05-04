package logic

import (
	"catchendb/src/node"
	"catchendb/src/util"
	"net/url"
	"sync"
)

var (
	transactionID    int
	transactionMutex *sync.Mutex
)

func getTransactionId() (id int) {
	transactionMutex.Lock()
	transactionID += 1
	id = transactionID
	transactionMutex.Unlock()
	return
}

type transaction struct {
	ID        int
	ChangeLog []transactionLog
}

type transactionLog struct {
	typ  int
	data node.Data
}

func handleBegin(keyword url.Values) []byte {
	rsp := Rsp{}
	code := keyword.Get(URL_CMD)
	_ = code
	return util.JsonOut(rsp)
}

func handleCommit(keyword url.Values) []byte {
	rsp := Rsp{}
	code := keyword.Get(URL_CMD)
	_ = code
	return util.JsonOut(rsp)
}

func handleRollBack(keyword url.Values) []byte {
	rsp := Rsp{}
	code := keyword.Get(URL_CMD)
	_ = code
	return util.JsonOut(rsp)
}

func initTransaction() {
	registerCMD(CMD_BEGIN, 1, handleBegin, TYPE_W)
	registerCMD(CMD_COMMIT, 1, handleCommit, TYPE_W)
	registerCMD(CMD_ROLLBACK, 1, handleRollBack, TYPE_W)
}

func init() {
	transactionID = 0
	transactionMutex = new(sync.Mutex)
}
