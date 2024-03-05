package main

import "C"
import (
	"bhd/app"
	"encoding/json"
	_ "net/http/pprof"
	"path/filepath"
	"strings"

	"bhd/dailyrotate"
	"bhd/log"
)

type BackendParams struct {
	Param1 string `json:"param1"`
	Param2 string `json:"param2"`
	Param3 string `json:"param3"`
	Param4 string `json:"param4"`
}

//export Call
func Call(methodName string, data string) *C.char {
	if methodName == "" {
		result := &app.ApiReturnStruct{ErrorID: 1, ErrorDescription: "Error: No method name provided"}
		content, _ := json.Marshal(result)
		return C.CString(string(content))
	}
	switch methodName {
	case "M1":
		var rq BackendParams
		json.Unmarshal([]byte(data), &rq)
		methodResult := app.InitializeWallet(rq.Param1)
		return C.CString(methodResult.ToJsonString())
		// mnemonic string
	case "M2":
		methodResult := app.GetWalletMnemonic()
		return C.CString(methodResult.ToJsonString())
	case "M3":
		methodResult := app.GetPublicBchAddress()
		return C.CString(methodResult.ToJsonString())
	case "M4":
		var rq BackendParams
		json.Unmarshal([]byte(data), &rq)
		methodResult := app.SignTransaction(rq.Param1)
		return C.CString(methodResult.ToJsonString())
	case "M5":
		var rq BackendParams
		json.Unmarshal([]byte(data), &rq)
		methodResult := app.GetBchAddressQrCode(rq.Param1, rq.Param2)
		return C.CString(methodResult.ToJsonString())
	default:
		result := &app.ApiReturnStruct{ErrorID: 2, ErrorDescription: "Error: Unknown method name:" + methodName + " (Not implemented))"}
		content, _ := json.Marshal(result)
		return C.CString(string(content))
	}
}

// log close handler function
func onCloseHappened(path string, didRotate bool) {

}

func main() {
	log.SetLogLevel(log.LevelDebug)
	pathFormat := filepath.Join("dir", "2006-01-02.log")
	f, err := dailyrotate.NewFile(pathFormat, onCloseHappened)
	if err != nil {
		log.Error("Failed to setup log file due to", err)
	}
	log.SetLogDest(f)
	app.InitializeWallet("caught before prosper fiscal glimpse verb badge animal dress property kiss analyst wrist bachelor panda view range either develop advice hidden impulse tail volcano")
	app.SignTransaction(strings.Replace("{'Hash':'','Size':0,'Height':0,'Index':0,'Version':1,'LockTime':0,'Inputs':[{'Sequence':0,'Value':100000,'PrevHash':'bc27832feb4f34e174d7c1bea0a4e4490b30ab222356c829f6649497c2970e68','PrevIndex':0}],'Outputs':[{'Value':2000,'Spent':false,'PkScript':'76a914eeed96fd3e0806986e8d19acfee0053b0366601188ac','Address':''},{'Value':97881,'Spent':false,'PkScript':'76a91464ee6ac83f1a70d8a38ab9a8d8230b540e12c29588ac','Address':''}],'InputVal':100000,'OutputVal':2000}", "'", "\"", -1))

}
