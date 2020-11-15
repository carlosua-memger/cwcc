package adjustmententry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"text/template"

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
	tconf["Panel_name"] = "Consumer Entry Detail"
	tconf["Add_new_button"] = "false"
	tconf["Rights"] = rights
	tconf["Add_back_button"] = "true"
	tconf["Add_new_button_details"] = "false"
	tconf["upload_button"] = "false"

	tconf["Action_baseURL"] = "/transaction/bulk_entry/details?rights=" + rights
	tconf["ActionCol_param"] = "ID"
	//tconf["hdr_id"] = ""
	tconf["Interpret_button"] = "false"
	tconf["Add_post_button"] = "false"
	tconf["Add_post_button_url"] = "/transaction/bulk_entry/details/post?rights=" + rights + "&hdr_id=" + tconf["hdr_id"]
	tconf["ActionCol_void"] = "false" // config for void click
	tconf["Add_deleteall_button_url"] = "/transaction/bulk_entry/details/deleteall?rights=" + rights + "&hdr_id=" + tconf["hdr_id"]

	tconf["Add_cancel_button"] = "false"
	tconf["Add_cancel_button_url"] = "/transaction/bulk_entry/details/cancel?rights=" + rights + "&hdr_id=" + tconf["hdr_id"]
	tconf["withtotal"] = "true"
	tconf["Archive"] = "false"
	tconf["org"] = str_OrgID
	tconf["Add_duplicate_button"] = "false"
	tconf["Add_duplicate_button_url"] = "/transaction/bulk_entry/details/duplicate?rights=" + rights + "&id=" + tconf["hdr_id"]

	tconf["level"] = "2"
	tconf["pageLength"] = last_length
	tconf["displayStart"] = last_start
	tconf["last_search"] = last_search
	if last_length == "" {
		tconf["pageLength"] = "10"
		tconf["displayStart"] = "0"
	}
	var header = []string{"Action", "Member", "Amount", "PN", "LAID", "LRID", "Remarks"} //column header
	if ln == "" {

		tconf["ismoreinfo"] = "true" // to description info
		arr_data := datatables.Data_row(`exec BULK_LOANHDR_Get 2,` + str_OrgID + `,` + tconf["hdr_id"] + `,0`)
		fmt.Println(arr_data)
		tconf["id"] = arr_data[0]
		tconf["partner"] = arr_data[1]
		tconf["trandate"] = arr_data[2]
		tconf["loantype"] = arr_data[3]
		tconf["first_due"] = arr_data[4]
		tconf["released"] = arr_data[5]
		tconf["noofpayments"] = arr_data[6]
		tconf["payment_mode"] = arr_data[7]
		tconf["remarks"] = arr_data[8]
		tconf["posted"] = arr_data[9]
		tconf["cancelled"] = arr_data[10]
		tconf["posted_by"] = arr_data[11]
		tconf["cancelled_by"] = arr_data[12]
		tconf["searchtext"] = arr_data[13]
		tconf["org"] = arr_data[14]
		tconf["added"] = arr_data[15]
		tconf["added_by"] = arr_data[16]
		tconf["updated"] = arr_data[17]
		tconf["updated_by"] = arr_data[18]
		tconf["status"] = arr_data[19]
		tconf["PartnerName"] = arr_data[20]
		tconf["LoantypeName"] = arr_data[21]
		tconf["PaymodeName"] = arr_data[22]
		tconf["dtl_total"] = arr_data[23]

		tran_status := arr_data[19]
		tconf["Add_new_button_url"] = "/transaction/bulk_entry/details/DaddHandler?rights=" + rights + "&hdr_id=" + tconf["hdr_id"] + "&member=" + tconf["member"] + "&member_name=" + tconf["member_name"]

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

		var ajaxURLdata = "/transaction/bulk_entry/details?rights=" + rights + "&ID=" + tconf["hdr_id"]

		type TblConf struct {
			Headercol []string
			AjaxUrl   string
			Tempconf  map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("transaction/bulk_entry/table.html"); err != nil {
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
		tconf["ActionCol_edit_url"] = "/transaction/bulk_entry/details/DEditHandler?rights=" + rights + "&h_id=" // config for edit click
		tconf["ActionCol_delete"] = "false"                                                                      // config for delete click
		tconf["ActionCol_delete_url"] = "/transaction/bulk_entry/details/DDeleteHandler?h_id="
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
		tconf["sql_total"] = `BULK_LOANDTL_List 1, ` + tconf["hdr_id"] + `,  ` + str + `,  ` + ln + `, '` + sr + `'`
		//tconf["sql_data"]  = `SIS_Item_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"] = `BULK_LOANDTL_List 0, ` + tconf["hdr_id"] + `,  ` + str + `,  ` + ln + `, '` + sr + `'`
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
			Rights string
			Tconf  map[string]string
		}

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
		if tmpl, err = tmpl.ParseFiles("transaction/bulk_entry/Dadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	} else {
		r.ParseForm()

		amount := strings.Replace(r.Form["amount"][0], ",", "", -1)

		fmt.Println(`BULK_LOANDTL_SaveAdd '` + username + `',` + r.Form["hdr_id"][0] + ` ,` + r.Form["member"][0] + `,'` + r.Form["pn"][0] + `',` + amount + `,'` + r.Form["remarks"][0] + `'`)
		arr_raw_data := datatables.Data_row(`BULK_LOANDTL_SaveAdd '` + username + `',` + r.Form["hdr_id"][0] + ` ,` + r.Form["member"][0] + `,'` + r.Form["pn"][0] + `',` + amount + `,'` + r.Form["remarks"][0] + `'`)
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
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )

		//_,session_user_id 	:= login.Get_account_info(r)

		type Data struct {
			Rights      string
			Tconf       map[string]string
			Ls_trantype [][]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		fmt.Println(`exec AEDTL_Get 3, 0,` + tconf["h_id"])
		arr_data := datatables.Data_row(`exec BULK_LOANDTL_Get 1, ` + tconf["h_id"] + `,0`)

		fmt.Println(arr_data)
		tconf["id"] = arr_data[0]
		tconf["hdr_id"] = arr_data[1]
		tconf["member"] = arr_data[2]
		tconf["amount"] = arr_data[3]
		tconf["pn"] = arr_data[4]
		tconf["laid"] = arr_data[5]
		tconf["lrid"] = arr_data[6]
		tconf["remarks"] = arr_data[7]
		tconf["added"] = arr_data[8]
		tconf["added_by"] = arr_data[9]
		tconf["updated"] = arr_data[10]
		tconf["updated_by"] = arr_data[11]
		tconf["MemberName"] = arr_data[12]

		arr_data_trantype := datatables.DataList(`exec TranType_Get 1, 0`)
		tmpl := template.New("Dedit.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/bulk_entry/Dedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arr_data_trantype})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	} else {
		r.ParseForm()
		amount := strings.Replace(r.Form["amount"][0], ",", "", -1)
		fmt.Println(`BULK_LOANDTL_SaveEdit '` + username + `',` + r.Form["id"][0] + `,` + r.Form["hdr_id"][0] + ` ,` + r.Form["member"][0] + `,'` + r.Form["pn"][0] + `',` + amount + `,'` + r.Form["remarks"][0] + `'`)
		arr_raw_data := datatables.Data_row(`BULK_LOANDTL_SaveEdit '` + username + `',` + r.Form["id"][0] + `,` + r.Form["hdr_id"][0] + ` ,` + r.Form["member"][0] + `,'` + r.Form["pn"][0] + `',` + amount + `,'` + r.Form["remarks"][0] + `'`)
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
		tconf["delete_url"] = "/transaction/bulk_entry/details/DDeleteHandler"
		tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/bulk_entry/modal_delete_loghdr.html"); err != nil {
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
		arr_data := datatables.Data_row(`BULK_LOANDTL_SaveDelete '` + username + `',` + r.Form["h_id"][0])
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
	arr_raw_data := datatables.DataList(`BULK_LOANHDR_Post '` + username + `',` + hdr_id)
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
	arr_raw_data := datatables.DataList(`BULK_LOANHDR_Cancel '` + username + `',` + hdr_id)
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

func DDeleteAllHandler(w http.ResponseWriter, r *http.Request) {
	hdr_id := r.URL.Query().Get("hdr_id")
	username, _ := login.Get_account_info(r)
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	r.ParseForm()

	fmt.Println(`GOOD TO DELETE`)
	arr_raw_data := datatables.DataList(`BULK_LOANHDR_SavePurgeDetail '` + username + `',` + str_OrgID + `,` + hdr_id)
	fmt.Println(arr_raw_data)

	js, err := json.Marshal(arr_raw_data)
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

	fmt.Println(`BULK_LOANHDR_Duplicate  '` + username + `' ,  ` + tconf["id"])
	arr_data := datatables.Data_row(`BULK_LOANHDR_Duplicate  '` + username + `' ,  ` + tconf["id"])
	js, err := json.Marshal(arr_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
