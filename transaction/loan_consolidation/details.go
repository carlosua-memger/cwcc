package loanconsolidation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	config "../../config"
	datatables "../../datatables"
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
	fmt.Println(w)

	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	hdr_id := r.URL.Query().Get("hdr_id")
	fromlr := r.URL.Query().Get("fromlr")
	lrid := r.URL.Query().Get("lrid")

	last_start := r.URL.Query().Get("last_start")
	last_length := r.URL.Query().Get("last_length")
	last_search := r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key, rights)
	fmt.Println("dec rights here")
	fmt.Println(dec_rights)
	tconf := make(map[string]string)
	tconf["hdr_id"] = hdr_id
	member := r.URL.Query().Get("member")
	ln := r.URL.Query().Get("length")
	isposted := r.URL.Query().Get("isposted")
	tconf["rights"] = rights
	tconf["Panel_name"] = "Loan Consolidation"
	tconf["Add_new_button"] = "false"
	tconf["Rights"] = rights
	tconf["Add_back_button"] = "true"
	tconf["Add_new_button_details"] = "true"
	tconf["upload_button"] = "false"
	//tconf["Add_new_button_url"] ="/inventory/supplier_order/details/DaddHandler?rights="+rights+"&hdr_id="+tconf["hdr_id"]
	tconf["Action_baseURL"] = "/transaction/loan_consolidation/details?rights=" + rights
	tconf["ActionCol_param"] = "ID"
	//tconf["hdr_id"] = ""

	tconf["withtotal"] = "true"
	tconf["Archive"] = "false"
	tconf["org"] = str_OrgID

	tconf["level"] = "2"
	tconf["pageLength"] = last_length
	tconf["displayStart"] = last_start
	tconf["last_search"] = last_search
	tconf["fromlr"] = fromlr
	tconf["lrid"] = lrid
	tconf["member"] = member
	tconf["isposted"] = isposted
	if fromlr == "1" {
		tconf["ActionCol_back_url"] = "/transaction/loan_release/details?rights=" + rights + "&ID=" + lrid
	} else if fromlr == "2" {
		tconf["ActionCol_back_url"] = "/transaction/loan_applications/HViewHandler?rights=" + rights + "&h_id=" + hdr_id

	} else {
		tconf["ActionCol_back_url"] = "/transaction/loan_applications/HEditHandler?rights=" + rights + "&hdr_id=" + hdr_id
	}

	if last_length == "" {
		tconf["pageLength"] = "10"
		tconf["displayStart"] = "0"
	}
	var header = []string{"Action", "Description", "Amount", "Remarks"} //column header
	if ln == "" {

		tconf["ismoreinfo"] = "true" // to description info
		arr_data := datatables.Data_row(`LA_Get 1 , ` + hdr_id)
		fmt.Println(arr_data)

		tconf["la_id"] = arr_data[0]
		tconf["la_trandate"] = arr_data[1]
		tconf["la_member"] = arr_data[2]
		tconf["la_partner"] = arr_data[3]
		tconf["la_loantype"] = arr_data[4]
		tconf["la_purpose"] = arr_data[5]
		tconf["la_amt"] = arr_data[6]
		tconf["la_amt_generated"] = arr_data[7]
		tconf["la_term"] = arr_data[8]
		tconf["la_paymode"] = arr_data[9]
		tconf["la_paysched"] = arr_data[10]
		tconf["la_coor"] = arr_data[11]
		tconf["la_coor_endorsed"] = arr_data[12]
		tconf["la_received"] = arr_data[13]
		tconf["la_remarks"] = arr_data[14]
		tconf["la_posted"] = arr_data[15]
		tconf["la_posted_by"] = arr_data[16]
		tconf["la_cancelled"] = arr_data[17]
		tconf["la_cancelled_by"] = arr_data[18]
		tconf["la_last_update"] = arr_data[19]
		tconf["la_org"] = arr_data[20]
		tconf["la_last_updateby"] = arr_data[21]
		tconf["la_searchkey"] = arr_data[22]
		tconf["la_status"] = arr_data[23]
		tconf["la_member_name"] = arr_data[24]
		tconf["la_partner_name"] = arr_data[25]
		tconf["la_loantype_name"] = arr_data[26]
		tconf["la_int_rate"] = arr_data[27]
		tconf["la_service_fee"] = arr_data[28]
		tconf["la_bank"] = arr_data[29]
		tconf["la_acct"] = arr_data[30]

		tconf["Add_new_button_url"] = "/transaction/loan_consolidation/details/DaddHandler?rights=" + rights + "&hdr_id=" + hdr_id + "&member=" + tconf["la_member"]

		tconf["Download_LoanBal_url"] = "/transaction/loan_consolidation/details/DownloadLoanBal?rights=" + rights + "&hdr_id=" + hdr_id + "&member=" + tconf["la_member"] + "&fromlr=" + fromlr + "&lrid=" + lrid

		//tconf["Add_new_button_url"] ="/inventory/supplier_order/details/DaddHandler?rights="+rights+"&hdr_id="+tconf["hdr_id"]

		// content access here
		// rights for tables
		if strings.Contains(dec_rights, "DAdd") && isposted == "0" {
			tconf["Add_new_button_details"] = "true"
		}
		if strings.Contains(dec_rights, "DUpload") {
			tconf["upload_button"] = "true"

		}
		if strings.Contains(dec_rights, "DPost") {
			if arr_data[1] == `Empty` && arr_data[1] == `Posted` {
				tconf["Add_post_button"] = "false"
			} else {
				tconf["Add_post_button"] = "true"
			}
		}

		if strings.Contains(dec_rights, "DCancel") {
			if arr_data[1] == `Empty` || arr_data[1] == `Cancelled` || arr_data[1] == `For Posting` {
				tconf["Add_cancel_button"] = "false"
			} else {
				tconf["Add_cancel_button"] = "true"
			}
		}
		if arr_data[1] == `Cancelled` {
			tconf["Add_post_button"] = "false"
			tconf["Add_new_button_details"] = "false"
			tconf["upload_button"] = "false"
			tconf["Add_cancel_button"] = "false"

		} else if arr_data[1] == `Posted` {
			tconf["Add_new_button_details"] = "false"
			tconf["Add_post_button"] = "false"
			//	tconf["Add_cancel_button"] ="true"

		} else if arr_data[1] == `Empty` {
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

		var ajaxURLdata = "/transaction/loan_consolidation/details?rights=" + rights + "&hdr_id=" + hdr_id + "&member=" + tconf["la_member"] + "&isposted=" + isposted
		fmt.Println("/transaction/loan_consolidation/details?rights=" + rights + "&hdr_id=" + hdr_id)
		type TblConf struct {
			Headercol []string
			AjaxURL   string
			Tconf     map[string]string
		}
		tmpl := template.New("table.html").Funcs(helper.FuncMap)
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("transaction/loan_consolidation/table.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &TblConf{header, ajaxURLdata, tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	} else {
		dr := r.URL.Query().Get("draw")
		isposted := r.URL.Query().Get("isposted")

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
		tconf["ActionCol"] = "true"     // config for action column
		tconf["ActionCol_param"] = "ID" // config for parameter of action

		tconf["ActionCol_edit"] = "true"   // config for edit click
		tconf["ActionCol_delete"] = "true" // config for delete click

		tconf["ActionCol_edit_is_modal"] = "true" // config for edit click  modal design
		//tconf["ActionCol_edit_url"] ="/timekeeping/overtime_logs/OvertimeLogsHeaderEdit?rights="+rights+"&h_id=" // config for edit click
		tconf["ActionCol_edit_url"] = "/transaction/loan_consolidation/details/DEditHandler?rights=" + rights + "&member=" + member + "&hdr_id=" + hdr_id + "&h_id=" // config for edit click

		tconf["ActionCol_delete_url"] = "/transaction/loan_consolidation/details/DDeleteHandler?h_id="
		tconf["ActionCol_detail"] = "false"    // config for details click
		tconf["ActionCol_add_child"] = "false" // config for addchild click
		tconf["style_size"] = "12"
		tconf["style_font"] = ""
		tconf["format_number_col"] = "Amount"
		tconf["col_img"] = "false"
		tconf["name_with_img"] = "false"
		tconf["col_img_num"] = "1"
		tconf["img_id_col"] = "5"
		tconf["col_hidden1"] = "0"
		tconf["dr"] = dr //draw
		tconf["void"] = "7"
		tconf["withtotal"] = "true"
		//geting total record etc
		// rights for tables
		if strings.Contains(dec_rights, "DEdit") {

			tconf["ActionCol_edit"] = "true"

		}
		if isposted == "1" {
			tconf["ActionCol_edit"] = "false"
			tconf["ActionCol_delete"] = "false"
			tconf["ActionCol"] = "false"
		}

		if strings.Contains(dec_rights, "DDelete") {

			tconf["ActionCol_delete"] = "true" // config for delete click
		}

		if strings.Contains(dec_rights, "HDetails") {
			tconf["ActionCol_detail"] = "true" // config for delete click

		}

		status := r.URL.Query().Get("status")
		if status == `Cancelled` {
			tconf["ActionCol_edit"] = "false"
			tconf["ActionCol_delete"] = "false"
		} else if status == `Posted` {
			tconf["ActionCol_edit"] = "false"
			tconf["ActionCol_delete"] = "false"
		}

		tconf["la_status"] = status
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
		tconf["sql_total"] = `LAR_List 1, ` + str_OrgID + `, ` + str + `,  ` + ln + `, '` + sr + `', ` + hdr_id
		//tconf["sql_data"]  = `SIS_Item_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"] = `LAR_List 0, ` + str_OrgID + `, ` + str + `,  ` + ln + `, '` + sr + `', ` + hdr_id
		fmt.Println(tconf["sql_data"])
		fmt.Println("dito" + isposted)
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
	member := r.URL.Query().Get("member")

	//fmt.Fprint(w,rights)

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
			Rights   string
			Tconf    map[string]string
			Ls_loans [][]string
		}
		arrDataLoans := datatables.DataList(`exec LRDTL_Get 2,` + member + `,` + hdr_id)

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
		if tmpl, err = tmpl.ParseFiles("transaction/loan_consolidation/Dadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arrDataLoans})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}

	} else {

		r.ParseForm()
		fmt.Println(`LAR_SaveAdd '` + username + `,` + r.Form["hdr_id"][0] + ` ,` + r.Form["lrdtl"][0] + `,` + r.Form["amount"][0] + `,'` + r.Form["remarks"][0] + `'`)
		arr_raw_data := datatables.Data_row(`LAR_SaveAdd '` + username + `',` + str_OrgID + `,` + r.Form["hdr_id"][0] + ` ,` + r.Form["lrdtl"][0] + `,` + r.Form["amount"][0] + `,'` + r.Form["remarks"][0] + `'`)
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
	member := r.URL.Query().Get("member")
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	username, _ := login.Get_account_info(r)
	hdr_id := r.URL.Query().Get("hdr_id")

	//rights :="rights"

	if r.Method == "GET" {

		tconf := make(map[string]string)
		tconf["h_id"] = r.URL.Query().Get("h_id")
		tconf["hdr_id"] = hdr_id
		tconf["rights"] = rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID
		tconf["member"] = member

		//tconf["data"] = datatables.DataList(sql)
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )

		//_,session_user_id 	:= login.Get_account_info(r)

		type Data struct {
			Rights   string
			Tconf    map[string]string
			Ls_loans [][]string
		}
		arrDataLoans := datatables.DataList(`exec LRDTL_Get 2,` + member + `,` + hdr_id)
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		fmt.Println(`exec LAR_Get 1, ` + tconf["h_id"] + `'`)
		arr_data := datatables.Data_row(`exec LAR_Get 1, ` + tconf["h_id"])

		fmt.Println(arr_data)
		tconf["lar_id"] = arr_data[0]
		tconf["lar_la"] = arr_data[1]
		tconf["lar_lrdtl"] = arr_data[2]
		tconf["lar_amt"] = arr_data[3]
		tconf["lar_remarks"] = arr_data[4]
		tconf["lar_lrdtl_name"] = arr_data[5]

		tmpl := template.New("Dedit.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/loan_consolidation/Dedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arrDataLoans})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	} else {
		r.ParseForm()

		fmt.Println(`LAR_SaveEdit '` + username + `,` + r.Form["hdr_id"][0] + ` ,` + r.Form["lrdtl"][0] + `,` + r.Form["amount"][0] + `,'` + r.Form["remarks"][0] + `'`)
		arr_raw_data := datatables.Data_row(`LAR_SaveEdit '` + username + `',` + str_OrgID + `,` + r.Form["id"][0] + `,` + r.Form["hdr_id"][0] + ` ,` + r.Form["lrdtl"][0] + `,` + r.Form["amount"][0] + `,'` + r.Form["remarks"][0] + `'`)
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
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	if r.Method == "GET" {
		tconf := make(map[string]string)
		tconf["h_id"] = r.URL.Query().Get("h_id")
		tconf["delete_url"] = "/transaction/loan_consolidation/details/DDeleteHandler"
		tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/loan_consolidation/modal_delete_loghdr.html"); err != nil {
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
		arr_data := datatables.Data_row(`LAR_SaveDelete '` + username + `',` + str_OrgID + `,` + r.Form["h_id"][0])

		js, err := json.Marshal(arr_data)
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
	fromlr := r.URL.Query().Get("fromlr")
	lrid := r.URL.Query().Get("lrid")

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
		tconf["la_member"] = member
		tconf["fromlr"] = fromlr
		tconf["lrid"] = lrid

		type Data struct {
			Rights  string
			Tconf   map[string]string
			LoanBal [][]string
		}
		var err error
		var html_file string = ``
		var html_filename string = ``

		arrdata := datatables.DataList(`exec LRDTL_Get 2,` + member + `,` + hdr_id)
		fmt.Println(arrdata)
		html_filename = "download_LoanBal.html"
		html_file = "transaction/loan_consolidation/download_LoanBal.html"

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
		amt := r.Form["qty"][0]

		//ret_add := datatables.Data_row(`SIS_CusSalDtl_DownLoad  1, '` + username + `', ` + hrd_id + `, ` + checkbox_id + `,` + qty)
		ret_add := datatables.Data_row(`LAR_SaveAdd '` + username + `',` + str_OrgID + `,` + hrd_id + ` ,` + checkbox_id + `,` + amt + `,''`)

		js, err := json.Marshal(ret_add)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}
