package other_disbursement

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	config "../../config"
	datatables "../../datatables"
	email "../../github.com/email-master"
	helper "../../helper"
	login "../../login"
)

func DListHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	//fmt.Fprint(w,"mars test here")
	//fmt.Fprint(w, r.URL.Query().Get("rights"))
	rights := r.URL.Query().Get("rights")

	//rights :="static here"
	fmt.Println("rights here")
	fmt.Println(rights)

	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)

	last_start := r.URL.Query().Get("last_start")
	last_length := r.URL.Query().Get("last_length")
	last_search := r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key, rights)

	tconf := make(map[string]string)
	tconf["hdr_id"] = r.URL.Query().Get("ID")
	hdr_id := r.URL.Query().Get("ID")
	fmt.Println("hdrID" + hdr_id)
	fmt.Println(tconf["hdr_id"])
	ln := r.URL.Query().Get("length")
	tconf["Panel_name"] = "Other Disbursement Detail"
	tconf["Add_new_button"] = "false"
	tconf["Rights"] = rights
	tconf["Add_back_button"] = "true"
	tconf["Add_new_button_details"] = "false"
	tconf["upload_button"] = "false"

	tconf["Action_baseURL"] = "/transaction/other_disbursement/details?rights=" + rights
	tconf["ActionCol_param"] = "ID"
	//tconf["hdr_id"] = ""
	tconf["Interpret_button"] = "false"
	tconf["Add_post_button"] = "false"
	tconf["Add_post_button_url"] = "/transaction/other_disbursement/details/post?rights=" + rights + "&hdr_id=" + tconf["hdr_id"]

	tconf["Add_deleteall_button_url"] = "/transaction/other_disbursement/details/deleteall?rights=" + rights + "&hdr_id=" + tconf["hdr_id"]
	tconf["ActionCol_void"] = "false" // config for void click
	tconf["Add_duplicate_button"] = "false"
	tconf["Add_duplicate_button_url"] = "/transaction/other_disbursement/details/duplicate?rights=" + rights + "&id=" + tconf["hdr_id"]

	tconf["Add_cancel_button"] = "false"
	tconf["Add_cancel_button_url"] = "/transaction/other_disbursement/details/cancel?rights=" + rights + "&hdr_id=" + tconf["hdr_id"]
	tconf["withtotal"] = "true"
	tconf["Archive"] = "false"
	tconf["org"] = str_OrgID

	tconf["level"] = "2"
	tconf["pageLength"] = last_length
	tconf["displayStart"] = last_start
	tconf["last_search"] = last_search
	if last_length == "" {
		tconf["pageLength"] = "10"
		tconf["displayStart"] = "0"
	}
	var header = []string{"Action", "Member", "Amount", "Account #", "Bank", "Transaction Type", "Remarks"} //column header
	if ln == "" {
		tconf["ismoreinfo"] = "true" // to description info
		arr_data := datatables.Data_row(`exec OtherDisHdr_Get ` + tconf["hdr_id"])
		fmt.Println(arr_data)
		tconf["id"]= arr_data[0]
		tconf["trandate"]= arr_data[1]
		tconf["otherdis_type"]= arr_data[2]
		tconf["remarks"]= arr_data[3]
		tconf["posted"]= arr_data[4]
		tconf["posted_by"]= arr_data[5]
		tconf["cancelled"]= arr_data[6]
		tconf["cancelled_by"]= arr_data[7]
		tconf["downloaded"]= arr_data[8]
		tconf["last_update"]= arr_data[9]
		tconf["partner"]= arr_data[10]
		tconf["searchtext"]= arr_data[11]
		tconf["org"]= arr_data[12]
		tconf["lastTouched"]= arr_data[13]
		tconf["lastTouchedBy"]= arr_data[14]
		tconf["updated_by"]= arr_data[15]
		tconf["updated"]= arr_data[16]
		tconf["Status"]= arr_data[17]
		tconf["otherdis_types_Name"]= arr_data[18]
		tconf["partner_name"]= arr_data[19]
		tconf["total"]= arr_data[20]



		tran_status := arr_data[15]
		tconf["Add_new_button_url"] = "/transaction/other_disbursement/details/DaddHandler?rights=" + rights + "&hdr_id=" + tconf["hdr_id"] + "&member=" + tconf["member"] + "&member_name=" + tconf["member_name"]

		//tconf["Add_new_button_url"] = "/transcation/supplier_order/details/DaddHandler?rights=" + rights + "&hdr_id=" + tconf["hdr_id"] + "&customer=" + tconf["suprec_hdr_customer"] + "&location=" + tconf["suprec_hdr_location"] + "&supplier=" + tconf["suprec_hdr_supplier"]
		//tconf["Email_url"] = "/inventory/supplier_order/details/email?rights=" + rights + "&hdr_id=" + tconf["hdr_id"] + "&customer=" + tconf["suprec_hdr_customer"] + "&location=" + tconf["suprec_hdr_location"]

		//tconf["Add_new_button_url"] ="/inventory/supplier_order/details/DaddHandler?rights="+rights+"&hdr_id="+tconf["hdr_id"]

		// content access here
		// rights for tables
		if strings.Contains(dec_rights, "DtlAdd") {
			tconf["Add_new_button_details"] = "true"
		}
		if strings.Contains(dec_rights, "DUpload") {
			tconf["upload_button"] = "true"

		}
		if strings.Contains(dec_rights, "DtlPost") {
			if tran_status == `Empty` && tran_status == `Posted` {
				tconf["Add_post_button"] = "false"
			} else {
				tconf["Add_post_button"] = "true"
			}
		}

		if strings.Contains(dec_rights, "DtlCancel") {
			if tran_status == `Empty` || tran_status == `Cancelled` || tran_status == `For Posting` {
				tconf["Add_cancel_button"] = "false"
			} else {
				tconf["Add_cancel_button"] = "true"
			}
		}
		if tran_status == `Cancelled` {
			tconf["Add_post_button"] = "false"
			tconf["Add_new_button_details"] = "false"
			tconf["upload_button"] = "false"
			tconf["Add_cancel_button"] = "false"
			tconf["Add_duplicate_button"] = "true"

		} else if tran_status == `Posted` {
			tconf["Add_new_button_details"] = "false"
			tconf["Add_post_button"] = "false"

			//	tconf["Add_cancel_button"] ="true"

		} else if tran_status == `Empty` {
			tconf["Add_post_button"] = "false"
			tconf["Add_cancel_button"] = "false"

		}

		if strings.Contains(dec_rights, "DInterpret") {
			tconf["Interpret_button"] = "true"

		}

		/*if( strings.Contains(LBR_hdr.Status , "Posted") ){
		 		tconf["Add_new_button_details"] ="false"
		 		tconf["upload_button"] ="false"
		 		tconf["Add_post_button"] ="false"
		 		tconf["Interpret_button"] ="false"
		 		if( strings.Contains(dec_rights, "DCancel") ){
		 			tconf["Add_cancel_button"] ="true"
				}
		 	}
		 	if( strings.Contains(LBR_hdr.Status , "Cancelled") ){
		 		tconf["Add_new_button_details"] ="false"
		 		tconf["upload_button"] ="false"
		 		tconf["Add_post_button"] ="false"
		 		tconf["Interpret_button"] ="false"
		 		tconf["Add_cancel_button"] ="false"
		 	}*/
		//end rights for tables

		//end of content accesss

		var ajaxURLdata = "/transaction/other_disbursement/details?rights=" + rights + "&ID=" + tconf["hdr_id"]

		type TblConf struct {
			Headercol []string
			AjaxUrl   string
			Tempconf  map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("transaction/other_disbursement/table.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &TblConf{header, ajaxURLdata, tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	} else {
		dr := r.URL.Query().Get("draw")

		sr := r.URL.Query().Get("search[value]")
		fmt.Println(sr)
		srt1, _ := strconv.Atoi(r.URL.Query().Get("start"))
		ln1, _ := strconv.Atoi(r.URL.Query().Get("length"))
		srt1 = (srt1 / ln1) + 1
		str := strconv.Itoa(srt1)
		sortcol := r.URL.Query().Get("order[0][column]")
		if sortcol == "0" {
			sortcol = "1"
		}
		srtcolindexInt, _ := strconv.Atoi(sortcol)
		var sortColName string
		for key, value := range header {
			if srtcolindexInt == key {
				sortColName = value
			}
		}
		fmt.Println(sortColName)
		tconf["ActionCol"] = "true"               // config for action column
		tconf["ActionCol_param"] = "ID"           // config for parameter of action
		tconf["ActionCol_edit"] = "false"         // config for edit click
		tconf["ActionCol_edit_is_modal"] = "true" // config for edit click  modal design
		//tconf["ActionCol_edit_url"] ="/timekeeping/overtime_logs/OvertimeLogsHeaderEdit?rights="+rights+"&h_id=" // config for edit click
		tconf["ActionCol_edit_url"] = "/transaction/other_disbursement/details/DEditHandler?rights=" + rights + "&h_id=" // config for edit click
		tconf["ActionCol_delete"] = "false"                                                                                   // config for delete click
		tconf["ActionCol_delete_url"] = "/transaction/other_disbursement/details/DDeleteHandler?h_id="
		tconf["ActionCol_detail"] = "false"    // config for details click
		tconf["ActionCol_add_child"] = "false" // config for addchild click
		tconf["style_size"] = "12"
		tconf["style_font"] = ""
		tconf["format_number_col"] = "Amount"
		//tconf["col_img"] = "false"
		//tconf["name_with_img"] = "true"
		//tconf["col_img_num"] = "1"
		//tconf["img_id_col"] = "5"
		tconf["col_hidden1"] = "0"
		tconf["dr"] = dr //draw
		//tconf["void"] = "7"
		//geting total record etc
		// rights for tables
		if strings.Contains(dec_rights, "DtlEdit") {
			tconf["ActionCol_edit"] = "true"
		}
		if strings.Contains(dec_rights, "DtlDelete") {
			tconf["ActionCol_delete"] = "true" // config for delete click

		}

		status := r.URL.Query().Get("status")
		if status == `Cancelled` {
			tconf["ActionCol_edit"] = "false"
			tconf["ActionCol_delete"] = "false"
		} else if status == `Posted` {
			tconf["ActionCol_edit"] = "false"
			tconf["ActionCol_delete"] = "false"
		}

		tconf["status"] = status
		//end rights for tables
		//_,session_user_id := login.Get_account_info(r)

		/*		fmt.Println("exec LBR_LogHdr_List 1, 2, 1, 1, 1, '"+sortColName+"', '', '"+sr+"'	" )
				fmt.Println(sortColName,str,sr)*/
		//fmt.Println(`exec LBR_OTHdr_List 1, 2, `+session_user_id_str+` , `+str+`, `+ln+`,  '`+sortColName+`', '', '`+sr+`'`)
		//tconf["sql_total"] = `exec LBR_OTHdr_List 1, 2, `+session_user_id_str+` , `+str+`, `+ln+`,  '`+sortColName+`', '', '`+sr+`'`
		//tconf["sql_data"]  = `exec LBR_OTHdr_List 0, 2, `+session_user_id_str+` , `+str+`, `+ln+`,  '`+sortColName+`', '', '`+sr+`'`

		//tconf["sql_total"] = `exec dailysumhdr_list 1, 0, 0 , `+str+`, `+ln+`, '`+sortcol+`','','`+sr+`'`
		//tconf["sql_data"]  = `exec dailysumhdr_list 0, 0, 0 , `+str+`, `+ln+`, '`+sortcol+`','','`+sr+`'`

		//tconf["sql_total"] = `SIS_Item_List 1, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_total"] = `OtherDisDtl_List 1, ` + tconf["hdr_id"] + `,  ` + str + `,  ` + ln + `, '` + sr + `'`
		//tconf["sql_data"]  = `SIS_Item_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"] = `OtherDisDtl_List 0, ` + tconf["hdr_id"] + `,  ` + str + `,  ` + ln + `, '` + sr + `'`
		fmt.Println(tconf["sql_data"])
		tconf["org"] = str_OrgID
		datatables.DatatableretArray(w, tconf)

	}
}

func DAddHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	username, _ := login.Get_account_info(r)
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	//username := "static"
	//ses_uID := 1
	//str_ses_uID :=strconv.Itoa(ses_uID)
	rights := r.URL.Query().Get("rights")
	hdr_id := r.URL.Query().Get("hdr_id")

	//fmt.Fprint(w,rights)

	if r.Method == "GET" {
		fmt.Println(`this is hdr id=`, hdr_id)

		tconf := make(map[string]string)
		//tconf["parentID"] =r.URL.Query().Get("parentID")
		tconf["username"] = username
		tconf["Org_id"] = str_OrgID
		tconf["hdr_id"] = hdr_id
		tconf["rights"] = rights

		type Data struct {
			Rights      string
			Tconf       map[string]string
			Ls_trantype [][]string
			Ls_bank 	[][]string
		}

		arr_data_trantype := datatables.DataList(`exec OthDisDtl_trantype_List 0`)
		arr_data_bank := datatables.DataList(`exec Banks_Get 1, ` + str_OrgID + `,''`)

		/*
			tmpl := template.New("Dadd.html")
			var err error
			if tmpl, err = tmpl.ParseFiles("inventory/supplier_order/Dadd.html"); err != nil {
				fmt.Println(err)
			}
			err1 := tmpl.Execute(w,&Data{rights , tconf})
			if err1 != nil {
				http.Error(w, err1.Error(), http.StatusInternalServerError)
			}
		*/

		var err error
		tmpl := template.New("Dadd.html")
		if tmpl, err = tmpl.ParseFiles("transaction/other_disbursement/Dadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arr_data_trantype, arr_data_bank})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	} else {
		r.ParseForm()

		fmt.Println(`OtherDisDtl_SaveAdd '` + username + `',` + r.Form["hdr_id"][0] + `,` + r.Form["member"][0] + `,` + r.Form["amount"][0] + `,` + r.Form["bank"][0] + `, '` + r.Form["acct"][0] + `' ,'` + r.Form["remarks"][0] + `',` + r.Form["trantype"][0])
		arr_raw_data := datatables.Data_row(`OtherDisDtl_SaveAdd '` + username + `',` + r.Form["hdr_id"][0] + `,` + r.Form["member"][0] + `,` + r.Form["amount"][0] + `,` + r.Form["bank"][0] + `,'` + r.Form["acct"][0] +  `','` + r.Form["remarks"][0] + `',` + r.Form["trantype"][0])
		fmt.Println(arr_raw_data)
		js, err := json.Marshal(arr_raw_data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

func DEditHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	rights := r.URL.Query().Get("rights")
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	username, _ := login.Get_account_info(r)

	//rights :="rights"

	if r.Method == "GET" {

		tconf := make(map[string]string)
		
		tconf["h_id"] = r.URL.Query().Get("h_id")
		tconf["rights"] = rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID

		//tconf["data"] = datatables.DataList(sql)
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lb r_assign )

		//_,session_user_id 	:= login.Get_account_info(r)

		type Data struct {
			Rights      string
			Tconf       map[string]string
			Ls_trantype [][]string
			Ls_bank 	[][]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )

		arr_data := datatables.Data_row(`exec OtherDisDtl_Get 1, ` + tconf["h_id"] + `,0`)
		arr_data_bank := datatables.DataList(`exec Banks_Get 1, ` + str_OrgID + `,''`)
		
		fmt.Println(arr_data)
		tconf["id"]= arr_data[0]
		tconf["hdr_id"]= arr_data[1]
		tconf["member"]= arr_data[2]
		tconf["amt"]= arr_data[3]
		tconf["bank"]= arr_data[4]
		tconf["acct"]= arr_data[5]
		tconf["remarks"]= arr_data[6]
		tconf["last_update"]= arr_data[7]
		tconf["savings_ave"]= arr_data[8]
		tconf["savings_int"]= arr_data[9]
		tconf["loan_tot"]= arr_data[10]
		tconf["loan_int"]= arr_data[11]
		tconf["cv"]= arr_data[12]
		tconf["past_due"]= arr_data[13]
		tconf["trantype_dtl"]= arr_data[14]
		tconf["added_by"]= arr_data[15]
		tconf["updated_by"]= arr_data[16]
		tconf["updated"]= arr_data[17]
		tconf["bank_name"]= arr_data[18]
		tconf["MemberName"]= arr_data[19]

		arr_data_trantype := datatables.DataList(`exec OthDisDtl_trantype_List 0`)
		
		tmpl := template.New("Dedit.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/other_disbursement/Dedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arr_data_trantype, arr_data_bank})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	} else {
		r.ParseForm()
		/* USER, ID, HDR, MEMBER, AMOUNT, BANK, ACCT, REMARKS, TRANTYPE*/
		arr_raw_data := datatables.Data_row(`OtherDisDtl_SaveEdit '` + username + `',` + r.Form["id"][0] + `,` + r.Form["hdr_id"][0] + ` ,` + r.Form["member"][0] + `,` + r.Form["amount"][0] + `,` + r.Form["bank"][0] + `,'` + r.Form["acct"][0] + `','` + r.Form["remarks"][0] + `',` + r.Form["trantype"][0])
		fmt.Println(arr_raw_data)
		js, err := json.Marshal(arr_raw_data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

func DDeleteHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	if r.Method == "GET" {
		tconf := make(map[string]string)
		tconf["h_id"] = r.URL.Query().Get("h_id")
		tconf["delete_url"] = "/transaction/other_disbursement/details/DDeleteHandler"
		tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/other_disbursement/modal_delete_loghdr.html"); err != nil {
			fmt.Println(err)
		}
		err = tmpl.Execute(w, tconf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		//session_username,_ := login.Get_account_info(r)
		r.ParseForm()
		username, _ := login.Get_account_info(r)
		arr_data := datatables.Data_row(`OtherDisDtl_SaveDelete '` + username + `',` + r.Form["h_id"][0])
		js, err := json.Marshal(arr_data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func DPostHandler(w http.ResponseWriter, r *http.Request) {
	hdr_id := r.URL.Query().Get("hdr_id")
	username, _ := login.Get_account_info(r)
	r.ParseForm()

	fmt.Println(`GOOD TO POST`)
	arr_raw_data := datatables.DataList(`OtherDisHdr_Post '` + username + `',` + hdr_id)
	fmt.Println(arr_raw_data)

	js, err := json.Marshal(arr_raw_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func DDeleteAllHandler(w http.ResponseWriter, r *http.Request) {
	hdr_id := r.URL.Query().Get("hdr_id")
	username, _ := login.Get_account_info(r)
	r.ParseForm()

	fmt.Println(`GOOD TO DELETE`)
	arr_raw_data := datatables.DataList(`ARHDR_SavePurgeDetail '` + username + `',` + hdr_id)
	fmt.Println(arr_raw_data)

	js, err := json.Marshal(arr_raw_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func DCancelHandler(w http.ResponseWriter, r *http.Request) {
	hdr_id := r.URL.Query().Get("hdr_id")
	username, _ := login.Get_account_info(r)
	r.ParseForm()

	fmt.Println(`GOOD TO POST`)
	arr_raw_data := datatables.DataList(`OtherDisHdr_Cancel '` + username + `',` + hdr_id)
	fmt.Println(arr_raw_data)

	js, err := json.Marshal(arr_raw_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func EmailHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	rights := r.URL.Query().Get("rights")
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	username, _ := login.Get_account_info(r)
	//rights :="rights"
	hdr_id := r.URL.Query().Get("hdr_id")
	//str_hdrID :=strconv.Itoa(int(hdr_id))
	//strconv.Itoa(int(time.Nanoseconds()))

	if r.Method == "GET" {

		tconf := make(map[string]string)
		tconf["h_id"] = hdr_id
		tconf["rights"] = rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID

		type Data struct {
			Rights string
			Tconf  map[string]string
		}

		arr_data := datatables.Data_row(`exec LRHDR_Get 3,` + str_OrgID + `,` + tconf["hdr_id"])
		fmt.Println(arr_data)
		tconf["id"] = arr_data[0]
		tconf["pn"] = arr_data[1]
		tconf["la"] = arr_data[2]
		tconf["trandate"] = arr_data[3]
		tconf["member"] = arr_data[4]
		tconf["partner"] = arr_data[5]
		tconf["loantype"] = arr_data[6]
		tconf["amt"] = arr_data[7]
		tconf["term"] = arr_data[8]
		tconf["paymode"] = arr_data[9]
		tconf["first_due"] = arr_data[10]
		tconf["released"] = arr_data[11]
		tconf["int_rate"] = arr_data[12]
		tconf["service_fee"] = arr_data[13]
		tconf["bank"] = arr_data[14]
		tconf["cheque"] = arr_data[15]
		tconf["acct"] = arr_data[16]
		tconf["cv"] = arr_data[17]
		tconf["cv_amt"] = arr_data[18]
		tconf["remarks"] = arr_data[19]
		tconf["posted"] = arr_data[20]
		tconf["posted_by"] = arr_data[21]
		tconf["cancelled"] = arr_data[22]
		tconf["cancelled_by"] = arr_data[23]
		tconf["last_update"] = arr_data[24]
		tconf["textfiledownloaded"] = arr_data[25]
		tconf["gl_mapping_account_cv"] = arr_data[26]
		tconf["gl_mapping_account_jv"] = arr_data[27]
		tconf["textfiledownloaded_by"] = arr_data[28]
		tconf["cv_date"] = arr_data[29]
		tconf["searchtext"] = arr_data[30]
		tconf["org"] = arr_data[31]
		tconf["added"] = arr_data[32]
		tconf["added_by"] = arr_data[33]
		tconf["updated"] = arr_data[34]
		tconf["updated_by"] = arr_data[35]
		tconf["status"] = arr_data[36]
		tconf["member_name"] = arr_data[37]
		tconf["partner_name"] = arr_data[38]
		tconf["loantype_name"] = arr_data[39]
		tconf["paymode_name"] = arr_data[40]
		tconf["bank_name"] = arr_data[41]
		tconf["LAR_amt"] = arr_data[42]
		tconf["member_email"] = arr_data[43]
		tconf["company_name"] = arr_data[44]
		tconf["cc_email"] = config.Email_id

		tmpl := template.New("email.html").Funcs(helper.FuncMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/loan_release/email.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	} else {
		r.ParseForm()
		hrd_id := r.URL.Query().Get("hdr_id")

		Org_id := login.Get_session_org_id(r)
		str_OrgID := strconv.Itoa(Org_id)
		supplier := "hello"
		fmt.Println(`SIS_Supplier_Get 1 ,  ` + str_OrgID + ` ,` + supplier)
		arr_supplier := datatables.Data_row(`SIS_Supplier_Get 4 ,  ` + str_OrgID + ` ,` + supplier)

		supplier_name := arr_supplier[2]
		org_name := arr_supplier[10]

		fmt.Println(`SIS_SupOrd_Custom 1,` + str_OrgID + ` ,` + hrd_id)
		SIS_SupOrdTags_List := datatables.DataList(`SIS_SupOrd_Custom 2,` + str_OrgID + ` ,` + hrd_id)
		SIS_SupOrdHdr_Report := datatables.DataList(`SIS_SupOrd_Custom 1,` + str_OrgID + ` ,` + hrd_id)
		ht_tbl := `<table width="100%" cellpadding="15">
	            <tr>
	                <th align="left" bgcolor="#E5E7E9"><b>#</th>
	                <th align="left" bgcolor="#E5E7E9"><b>Item</th>
	                <th align="right" bgcolor="#E5E7E9"><b>Qty</th>
	                <th align="left" bgcolor="#E5E7E9"><b>Remarks</th>
	            </tr>        
			`
		for _, value := range SIS_SupOrdTags_List {
			ht_tbl += `<tr>`
			ht_tbl += `<td>` + value[0] + `</td>`
			ht_tbl += `<td>` + value[1] + `</td>`
			ht_tbl += `<td align="right">` + value[2] + `</td>`
			ht_tbl += `<td>` + value[3] + `</td>`
			ht_tbl += `</tr>`
		}
		for _, value2 := range SIS_SupOrdHdr_Report {
			ht_tbl += `<tr>`
			ht_tbl += `<td bgcolor="#E5E7E9"><u><b>` + value2[5] + `</u></b></td>`
			ht_tbl += `<td bgcolor="#E5E7E9"></td>`
			ht_tbl += `<td align="right" bgcolor="#E5E7E9"><b><u>` + value2[4] + `</u></b></td>`
			ht_tbl += `<td bgcolor="#E5E7E9"></td>`
			ht_tbl += `</tr>`

		}
		ht_tbl += `</table>`

		e := email.NewEmail()
		e.From = config.Email_id
		e.To = []string{r.Form["emailto"][0]}
		e.Bcc = []string{config.Email_bcc}
		fmt.Println(`cc email here : `, r.Form["ccto"][0])
		arr_line := strings.Split(r.Form["ccto"][0], ",")
		var arr_CC []string
		for key, val := range arr_line {
			fmt.Println(` key : `, key)
			fmt.Println(` val : `, val)
			arr_CC = append(arr_CC, val)
		}
		e.Cc = []string{}

		if len(arr_line) > 0 {
			e.Cc = arr_CC
		}

		e.Subject = org_name + ` Supplier Order (ID# ` + hrd_id + `)`
		e.HTML = []byte(`Hi ` + supplier_name + `,<br><br>Sending you list of orders, as follows;<br><br>` + ht_tbl + `<br>Additional Message: ` + r.Form["message"][0] + `<br>Thank you<br>` + org_name + `<br><br>This is a system generated email. Please send your reply to ` + r.Form["ccto"][0] + `.<br> Powered by <a href="www.memger.com"> Memger Tech Inc.</a> <i>Inventory like a <b>boss</b>... <a href="www.memger.com">try us for <b>free</b>.</a></i>`)

		retErr := e.Send("smtp.zoho.com:587", smtp.PlainAuth("", e.From, config.Email_password, "smtp.zoho.com"))
		fmt.Println(`Email sending Error : `, retErr)
		var json_ret string

		if retErr == nil {
			json_ret = `Success`
		} else {
			json_ret = `Error`
		}
		fmt.Println("done Sending email")

		js, err := json.Marshal(json_ret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

func Download_LoanBalHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	username, _ := login.Get_account_info(r)
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	//username := "static"
	//ses_uID := 1
	//str_ses_uID :=strconv.Itoa(ses_uID)
	rights := r.URL.Query().Get("rights")
	hdr_id := r.URL.Query().Get("hdr_id")
	member := r.URL.Query().Get("member")
	//isbarcode :=r.URL.Query().Get("isbarcode")

	if r.Method == "GET" {
		fmt.Println(`this is hdr id=`, hdr_id)

		tconf := make(map[string]string)
		//tconf["parentID"] =r.URL.Query().Get("parentID")
		tconf["username"] = username
		tconf["Org_id"] = str_OrgID
		tconf["hdr_id"] = hdr_id
		tconf["rights"] = rights
		tconf["test_js"] = `alert("from webserver")`
		tconf["member"] = member

		type Data struct {
			Rights  string
			Tconf   map[string]string
			LoanBal [][]string
		}
		var err error
		var html_file string = ``
		var html_filename string = ``

		arrdata := datatables.DataList(`exec LRDTL_Get 4,` + member + `,` + hdr_id)
		fmt.Println(arrdata)
		html_filename = "download_LoanBal.html"
		html_file = "transaction/other_disbursement/download_LoanBal.html"

		tmpl := template.New(html_filename)
		if tmpl, err = tmpl.ParseFiles(html_file); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arrdata})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}

	} else {

		r.ParseForm()
		hrd_id := r.Form["hrd_id"][0]
		checkbox_id := r.Form["checkbox_id"][0]
		amt := r.Form["amt"][0]
		member := r.Form["member"][0]

		//ret_add := datatables.Data_row(`SIS_CusSalDtl_DownLoad  1, '` + username + `', ` + hrd_id + `, ` + checkbox_id + `,` + qty)
		ret_add := datatables.Data_row(`LPDTL_SaveAdd '` + username + `',` + hrd_id + ` ,` + member + `,` + checkbox_id + `,` + amt + `,''`)
		js, err := json.Marshal(ret_add)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

type ArDtlUploadReq struct {
	ArHrd            string
	MemberId         string
	PartnerTransType string
	Amount           string
	Remarks          string
	TxnLine          int
	TotalLines       string
	SessionName      string
	W                http.ResponseWriter
	R                *http.Request
}

type Translog struct {
	Txtline    int
	MsgID      int
	Msg        interface{}
	LineDetail interface{}
}

var translog_arr []Translog

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	//get the multipart reader for the request.
	login.Session_renew(w, r)
	session_username, _ := login.Get_account_info(r)
	getid := r.URL.Query().Get("hid")
	//device := r.URL.Query().Get("device")
	partner := r.URL.Query().Get("partner")
	trandate := r.URL.Query().Get("trandate")
	translog_arr = nil

	go config.LogHandler(fmt.Sprint("hid :", getid))
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Println(translog_arr)
	fmt.Println("Upload successful.")
	go config.LogHandler("Upload successfully")
	t0 := time.Now().Local()
	t0_newFormat := `Start Time : ` + t0.Format("15:04:05")

	translog_arr = append(translog_arr, Translog{0, 0, t0_newFormat, ``})
	var txline int = 1

	var ArDltUploadReqList []*ArDtlUploadReq
	//copy each part to destination.
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() == "" {
			continue
		}
		fmt.Println("print file name :", part.FileName())
		if strings.Contains(part.FileName(), ".txt") == false {
			translog_arr = append(translog_arr, Translog{0, 0, `Error: invalid filename.`, part.FileName()})
			fmt.Println(`Error: invalid filename.`)
			continue
		}

		fmt.Println("print file name :", part.FileName())
		config.LogHandler(fmt.Sprint("print file name :", part.FileName()))
		t := time.Now() //get current date time to make name of folder
		datetimefoldername := t.Format("20060102150405")
		//os.Mkdir("C:/Go/muse/uploads/"+datetimefoldername,0777)
		//os.Mkdir("textfiles/upload_lp/"+datetimefoldername, 0777)
		dst, err := os.Create("textfiles/upload_ar/" + datetimefoldername + "_" + partner + "_" + part.FileName())
		defer dst.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, part); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//on parsing data from file
		contents, err := ioutil.ReadFile("textfiles/upload_ar/" + datetimefoldername + "_" + partner + "_" + part.FileName())
		if err != nil {
			fmt.Println("error on reading file :", err)
			go config.LogHandler(fmt.Sprint("error on reading file :", err))
		}
		//fmt.Println(string(contents));
		partnerCheckHeader := strings.TrimLeft(string(contents[1:5]), "0")
		monthCheckHeader := string(contents[5:7])
		dayCheckHeader := string(contents[7:9])
		yearCheckHeader := string(contents[9:13])

		monthHeader := string(trandate[0:2])
		dayHeader := string(trandate[3:5])
		yearHeader := string(trandate[6:10])

		if strings.Contains(string(contents), "^") == true && (partner != partnerCheckHeader || monthCheckHeader != monthHeader || dayCheckHeader != dayHeader || yearCheckHeader != yearHeader) {
			translog_arr = append(translog_arr, Translog{0, 0, `Error: header not match.`, string(contents[0:13])})
			fmt.Println(`Error: header not match`)
			continue
		}

		lines := strings.Split(string(contents), "\n")
		len_lines := len(lines)
		str_len_lines := strconv.Itoa(len_lines)

		fmt.Println(`lines`, lines)
		fmt.Println(`len_lines`, len_lines)
		fmt.Println(`str_len_lines`, str_len_lines)
		//translog_arr = nil
		//lines,_:=strings.ReadString('\n')
		//fmt.Println("data sa arr 0", lines[0][10:len(lines[0])-9])
		for key, _ := range lines {

			if key > 0 && key < (len(lines)) {
				//if string(len(lines[key])) != "40" {
				//	translog_arr = append(translog_arr, Translog{0, 0, `Error: line too short.`})
				//	fmt.Println(`Error: line too short.`)
				//}strconv.ParseFloat("3.1415", 64)
				if len(string(lines[key])) > 29 {
					amount_check, _ := strconv.ParseFloat(string(lines[key][30:39]), 64)
					if strings.Contains(string(lines[key]), "^") == false && amount_check > 0 {
						lrid := strings.TrimLeft(string(lines[key][0:7]), "0")
						member := strings.TrimLeft(string(lines[key][18:26]), "0")
						trantype := string(lines[key][26:30])
						amount := strings.TrimLeft(string(lines[key][30:39]), "0")
						//lrid_len := len(lrid)

						fmt.Println(len(lines))

						fmt.Println(`lrid`, lrid)
						fmt.Println(`member`, member)
						fmt.Println(`trantype`, trantype)
						fmt.Println(`amount`, amount)
						//logid:=lines[key][8:len(lines[key])-29]

						if len(lrid) == 0 {

							ArDltUploadReqList = append(ArDltUploadReqList, &ArDtlUploadReq{
								ArHrd:            getid,
								MemberId:         member,
								PartnerTransType: trantype,
								Amount:           amount,
								Remarks:          string(lines[key]),
								TxnLine:          txline,
								TotalLines:       str_len_lines,
								SessionName:      session_username,
								W:                w,
								R:                r,
							})

						}
					}

				} else {
					translog_arr = append(translog_arr, Translog{txline, 0, `Error: invalid line`, nil})
				}
				fmt.Println("loop line no: %v", txline)
				txline = txline + 1
			}
		}

	}

	workCh := make(chan *ArDtlUploadReq, 10)
	initWorkers(workCh)

	progressWg.Add(len(ArDltUploadReqList))
	for _, arDltUpload := range ArDltUploadReqList {
		workCh <- arDltUpload
	}
	progressWg.Wait()

	//display success message.
	//display(w, "upload", "Upload successful.")

	//fmt.Println(translog_arr)
	fmt.Println("Upload successful.")
	//go config.LogHandler("Upload successful")
	t1 := time.Now().Local()
	t1_newFormat := `End Time : ` + t1.Format("15:04:05")

	translog_arr = append(translog_arr, Translog{0, 0, t1_newFormat, ``})
	js, err := json.Marshal(translog_arr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	//wg.Wait()
}

type User_upload_list struct {
	Username   string
	Desciption string
}

var User_ups_status_list []User_upload_list //save to memory,
func upload_info(w http.ResponseWriter, r *http.Request, lineno int, line_tot string, Username string) {
	str_line := strconv.Itoa(lineno)
	pollstat := str_line + ` / ` + line_tot
	var temp_var_User_ups_status_list []User_upload_list
	for _, val := range User_ups_status_list {
		if val.Username != Username {
			temp_var_User_ups_status_list = append(User_ups_status_list, val)
		}
	}
	User_ups_status_list = temp_var_User_ups_status_list
	User_ups_status_list = append(User_ups_status_list, User_upload_list{Username, pollstat})
	return
}

func AjaxPollingeHandler(w http.ResponseWriter, r *http.Request) {
	go login.Session_renew(w, r)
	session_username, _ := login.Get_account_info(r)
	var uploading_status string = ``
	for _, val := range User_ups_status_list {
		if val.Username == session_username {
			uploading_status = val.Desciption
		}
		//	fmt.Println(`key:`, key , `val:`,val)
	}

	js, err := json.Marshal(uploading_status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func Duplicate(w http.ResponseWriter, r *http.Request) {
	username, _ := login.Get_account_info(r)
	tconf := make(map[string]string)
	tconf["id"] = r.URL.Query().Get("id")

	fmt.Println(`ARHDR_Duplicate  '` + username + `' ,  ` + tconf["id"])
	arr_data := datatables.Data_row(`ARHDR_Duplicate  '` + username + `' ,  ` + tconf["id"])
	js, err := json.Marshal(arr_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
