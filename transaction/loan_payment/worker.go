package loanpayment

import (
	"fmt"
	"sync"

	config "../../config"
)

var (
	progressWg sync.WaitGroup
)

func initWorkers(workCh <-chan *ArDtlUploadReq) {
	for i := uint64(0); i < 10; i++ {
		go func() {
			for work := range workCh {
				insertTxn(work)
			}
		}()
	}
}

func insertTxn(txn *ArDtlUploadReq) {

	defer progressWg.Done()
	fmt.Println("actual parseLine : %s", txn.TxnLine)

	fmt.Println("LPDTL_Upload " + txn.ArHrd + "," + txn.MemberId + "," + txn.PartnerTransType + "," + txn.Amount + ",'" + txn.Remarks + "'," + txn.Lrid)
	config.LogHandler("LPDTL_Upload " + txn.ArHrd + "," + txn.MemberId + "," + txn.PartnerTransType + "," + txn.Amount + ",'" + txn.Remarks + "'," + txn.Lrid)
	var rr Translog
	sqlstr := "LPDTL_Upload " + txn.ArHrd + "," + txn.MemberId + "," + txn.PartnerTransType + "," + txn.Amount + ",'" + txn.Remarks + "'," + txn.Lrid
	rowret, err, _, _ := config.Ap_sql(sqlstr, 1)
	//_ ,err, _,_ := config.Ap_sql("LBR_LogDtl_Save 'UploadBioText','"+login.Uname+"', 0, "+getid+", 0, 0, 0, '"+logid+"', '"+logtime+"'",3)
	if err != nil {
		fmt.Println(err.Error)
		//	panic(err.Error())
	}
	for rowret.Next() {
		err = rowret.Scan(&rr.MsgID, &rr.Msg)
		if err != nil {
			fmt.Println(err.Error)
			//		panic(err.Error())
		}

		translog_arr = append(translog_arr, Translog{txn.TxnLine, rr.MsgID, rr.Msg, nil})
		upload_info(txn.W, txn.R, txn.TxnLine, txn.TotalLines, txn.SessionName)

	}

}
