package loanapplication

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

type Profile struct {
	Message string // error status output or validated data
	Status  bool   //pass or fail
	ID      string
}

func HListHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	//fmt.Fprint(w,"mars test here")
	//fmt.Fprint(w, r.URL.Query().Get("rights"))
	rights := r.URL.Query().Get("rights")
	//rights :="static here"

	last_start := r.URL.Query().Get("last_start")
	last_length := r.URL.Query().Get("last_length")
	last_search := r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key, rights)
	//dec_rights :="static here"
	tconf := make(map[string]string)
	ln := r.URL.Query().Get("length")
	tconf["Panel_name"] = "Savings Withdrawal"
	tconf["Add_new_button"] = "false"
	tconf["Rights"] = rights
	tconf["Add_back_button"] = "false"
	tconf["Add_new_button_details"] = "false"
	tconf["upload_button"] = "false"
	tconf["Add_new_button_url"] = "/transaction/savings_withdrawal/HaddHandler?rights=" + rights
	tconf["Action_baseURL"] = "/transaction/savings_withdrawal/details?rights=" + rights
	tconf["ActionCol_param"] = "ID"
	tconf["hdr_id"] = "0"
	tconf["Interpret_button"] = "false"
	tconf["Add_post_button"] = "false"
	tconf["Add_post_button_url"] = ""
	tconf["Archive"] = "true" // is archive

	if strings.Contains(dec_rights, "DtlPost") {
		tconf["ActionCol_post"] = "true"
	}

	tconf["ActionCol_post_url"] = "/transaction/savings_withdrawal/LPostHandler?hdr_id="

	tconf["Add_upload_button"] = "false"

	tconf["Add_cancel_button"] = "false"
	tconf["Add_cancel_button_url"] = ""
	tconf["ismoreinfo"] = "false" // to description info
	tconf["withtotal"] = "false"

	tconf["level"] = "1"
	tconf["pageLength"] = last_length
	tconf["displayStart"] = last_start
	tconf["last_search"] = last_search
	if last_length == "" {
		tconf["pageLength"] = "10"
		tconf["displayStart"] = "0"
	}
	var header = []string{"Action", "ID", "Status", "Trandate", "Member", "Partner", "Transaction Type", "Amount Requested", "Bank", "Cheque", "Acct", "CV", "Amount Generated", "Amount Issued", "Date Issued", "Remarks"} //column header
	if ln == "" {
		// content access here
		if strings.Contains(dec_rights, "HdrAdd") {
			tconf["Add_new_button"] = "true"
			//fmt.Println("run herer")
		}

		//end of content accesss
		var ajaxURLdata = "/transaction/savings_withdrawal?rights=" + rights
		type TblConf struct {
			Headercol []string
			AjaxUrl   string
			Tempconf  map[string]string
		}
		tmpl := template.New("table.html").Funcs(helper.FuncMap)
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("transaction/savings_withdrawal/table.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &TblConf{header, ajaxURLdata, tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	} else {
		dr := r.URL.Query().Get("draw")
		archive := r.URL.Query().Get("arc")
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
		tconf["ActionCol"] = "true"                // config for action column
		tconf["ActionCol_param"] = "ID"            // config for parameter of action
		tconf["ActionCol_edit"] = "false"          // config for edit click
		tconf["ActionCol_edit_is_modal"] = "false" // config for edit click
		//tconf["ActionCol_edit_url"] ="/timekeeping/overtime_logs/OvertimeLogsHeaderEdit?rights="+rights+"&h_id=" // config for edit click
		tconf["ActionCol_edit_url"] = "/transaction/savings_withdrawal/HEditHandler?rights=" + rights + "&h_id=" // config for edit click
		tconf["ActionCol_view"] = "false"                                                                        // config for edit click
		tconf["ActionCol_view_is_modal"] = "false"                                                               // config for edit click
		//tconf["ActionCol_edit_url"] ="/timekeeping/overtime_logs/OvertimeLogsHeaderEdit?rights="+rights+"&h_id=" // config for edit click
		tconf["ActionCol_view_url"] = "/transaction/savings_withdrawal/HViewHandler?rights=" + rights + "&h_id=" // config for edit click

		tconf["ActionCol_delete"] = "false" // config for delete click
		tconf["ActionCol_delete_url"] = "/transaction/savings_withdrawal/HDeleteHandler?h_id="
		tconf["ActionCol_detail"] = "false"    // config for details click
		tconf["ActionCol_add_child"] = "false" // config for addchild click
		tconf["ActionCol_loanconsolidation_url"] = "/transaction/loan_consolidation/detals/DListHandler?h_id="
		tconf["style_size"] = "12"
		tconf["style_font"] = ""
		tconf["format_number_col"] = "Amount"
		tconf["dr"] = dr //draw
		//geting total record etc
		// rights for tables
		if strings.Contains(dec_rights, "HdrEdit") {
			tconf["ActionCol_edit"] = "true"
		}

		if strings.Contains(dec_rights, "HdrView") {
			tconf["ActionCol_view"] = "true"
		}
		if strings.Contains(dec_rights, "HdrDelete") {
			tconf["ActionCol_delete"] = "true" // config for delete click

		}
		if strings.Contains(dec_rights, "HDetails") {
			tconf["ActionCol_detail"] = "true" // config for details click

		}

		if strings.Contains(dec_rights, "HdrView") {
			tconf["ActionCol_view"] = "true" // config for view click

		}

		//end rights for tables
		//_,session_user_id := login.Get_account_info(r)
		session_user_id := 1 //static here
		session_user_id_str := strconv.Itoa(session_user_id)
		fmt.Println(session_user_id_str)
		/*		fmt.Println("exec LBR_LogHdr_List 1, 2, 1, 1, 1, '"+sortColName+"', '', '"+sr+"'	" )
				fmt.Println(sortColName,str,sr)*/
		//fmt.Println(`exec LBR_OTHdr_List 1, 2, `+session_user_id_str+` , `+str+`, `+ln+`,  '`+sortColName+`', '', '`+sr+`'`)
		//tconf["sql_total"] = `exec LBR_OTHdr_List 1, 2, `+session_user_id_str+` , `+str+`, `+ln+`,  '`+sortColName+`', '', '`+sr+`'`
		//tconf["sql_data"]  = `exec LBR_OTHdr_List 0, 2, `+session_user_id_str+` , `+str+`, `+ln+`,  '`+sortColName+`', '', '`+sr+`'`

		//tconf["sql_total"] = `exec dailysumhdr_list 1, 0, 0 , `+str+`, `+ln+`, '`+sortcol+`','','`+sr+`'`
		//tconf["sql_data"]  = `exec dailysumhdr_list 0, 0, 0 , `+str+`, `+ln+`, '`+sortcol+`','','`+sr+`'`
		Org_id := login.Get_session_org_id(r)
		str_OrgID := strconv.Itoa(Org_id)

		//tconf["sql_total"] = `SIS_Item_List 1, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_total"] = `SW_List 1, ` + str_OrgID + `,  ` + str + `,  ` + ln + `, '` + sr + `',` + archive
		//tconf["sql_data"]  = `SIS_Item_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"] = `SW_List 0, ` + str_OrgID + `,  ` + str + `,  ` + ln + `, '` + sr + `',` + archive

		datatables.DatatableretArray(w, tconf)

	}
}

type RetMessage struct {
	RetNo  int
	RetMsg string
}

func HAddHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	username, ses_uID := login.Get_account_info(r)

	//username := "static"
	//ses_uID := 1
	str_ses_uID := strconv.Itoa(ses_uID)
	rights := r.URL.Query().Get("rights")
	fmt.Println(str_ses_uID)
	//fmt.Fprint(w,rights)

	if r.Method == "GET" {

		Org_id := login.Get_session_org_id(r)
		str_OrgID := strconv.Itoa(Org_id)
		tconf := make(map[string]string)
		//tconf["parentID"] =r.URL.Query().Get("parentID")
		tconf["username"] = username
		tconf["Org_id"] = str_OrgID
		tconf["test_js"] = `alert("from webserver")`
		tconf["Add_upload_button"] = "false"

		arr_data_partner := datatables.DataList(`exec partners_get 2,` + str_OrgID + `,''`)
		arr_data_trantype := datatables.DataList(`exec TranType_Get 3 ,0`)
		arr_data_bank := datatables.DataList(`exec Banks_Get 1, ` + str_OrgID + `,''`)
		arr_data_cvto := datatables.DataList(`exec Mas_GLCV_Get_droplist 'add',0`)
		arr_data_jvto := datatables.DataList(`exec Mas_GLJV_Get_droplist 'add', 0`)

		//arr_data_tag := datatables.DataList(`exec SIS_SupOrdTags_List 1,`+str_OrgID+` , 0`  )
		type Data struct {
			Rights      string
			Tconf       map[string]string
			Ls_partner  [][]string
			Ls_trantype [][]string
			Ls_bank     [][]string
			Ls_cvto     [][]string
			Ls_jvto     [][]string

			//Arr_tag [][]string
		}
		tmpl := template.New("Hadd.html")
		//tmpl := template.New("Hadd.html").Funcs(local_FuncMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/savings_withdrawal/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arr_data_partner, arr_data_trantype, arr_data_bank, arr_data_cvto, arr_data_jvto})
		//err1 := tmpl.Execute(w,&Data{rights , tconf,arr_data_location,arr_data_supplier,arr_data_tag}
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}

	}
}

func HEditHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	rights := r.URL.Query().Get("rights")
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	decrights := config.Decrypt(config.Encrypt_key, rights)

	//rights :="rights"

	if r.Method == "GET" {

		username, _ := login.Get_account_info(r)

		tconf := make(map[string]string)
		tconf["Add_upload_button"] = "false"
		tconf["Add_post_button"] = "false"
		tconf["Add_cancel_button"] = "false"
		tconf["h_id"] = r.URL.Query().Get("h_id")
		if tconf["h_id"] > "0" {
		} else {
			tconf["h_id"] = r.URL.Query().Get("hdr_id")
		}
		tconf["rights"] = rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID
		//tconf["LoanConso_url"] = "/transaction/loan_consolidation/details/DListHandler?rights=" + rights + "&hdr_id=" + tconf["h_id"]
		tconf["Add_post_button_url"] = "/transaction/savings_withdrawal/post?rights=" + rights + "&hdr_id=" + tconf["h_id"]
		tconf["Add_cancel_button_url"] = "/transaction/savings_withdrawal/cancel?rights=" + rights + "&hdr_id=" + tconf["h_id"]
		tconf["Add_new_button_url"] = "/transaction/savings_withdrawal/HaddHandler?rights=" + rights
		//tconf["data"] = datatables.DataList(sql)
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )

		//_,session_user_id 	:= login.Get_account_info(r)

		type Data struct {
			Rights      string
			Tconf       map[string]string
			Ls_partner  [][]string
			Ls_trantype [][]string
			Ls_bank     [][]string
			Ls_cvto     [][]string
			Ls_jvto     [][]string
			//Arr_tag [][]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		arr_data := datatables.Data_row(`exec SW_Get 1,'` + tconf["h_id"] + `'`)

		fmt.Println(arr_data)
		tconf["id"] = arr_data[0]
		tconf["trandate"] = arr_data[1]
		tconf["member"] = arr_data[2]
		tconf["partner"] = arr_data[3]
		tconf["trantype"] = arr_data[4]
		tconf["amt_requested"] = arr_data[5]
		tconf["bank"] = arr_data[6]
		tconf["cheque"] = arr_data[7]
		tconf["acct"] = arr_data[8]
		tconf["amt_issued"] = arr_data[9]
		tconf["date_issued"] = arr_data[10]
		tconf["cv"] = arr_data[11]
		tconf["remarks"] = arr_data[12]
		tconf["posted"] = arr_data[13]
		tconf["posted_by"] = arr_data[14]
		tconf["cancelled"] = arr_data[15]
		tconf["cancelled_by"] = arr_data[16]
		tconf["last_update"] = arr_data[17]
		tconf["amt_generated"] = arr_data[18]
		tconf["textfiledownloaded"] = arr_data[19]
		tconf["textfiledownloaded_by"] = arr_data[20]
		tconf["cv_date"] = arr_data[21]
		tconf["cvto"] = arr_data[22]
		tconf["jvto"] = arr_data[23]
		tconf["searchtext"] = arr_data[24]
		tconf["org"] = arr_data[25]
		tconf["added"] = arr_data[26]
		tconf["added_by"] = arr_data[27]
		tconf["updated"] = arr_data[28]
		tconf["updated_by"] = arr_data[29]
		tconf["Status"] = arr_data[30]
		tconf["MemberName"] = arr_data[31]
		tconf["PartnerName"] = arr_data[32]
		tconf["TrantypeName"] = arr_data[33]
		tconf["bank_name"] = arr_data[34]
		tconf["amt_generated_noformat"] = arr_data[35]

		tconf["Action_button_save"] = "true"
		if arr_data[30] == `For Posting` {
			if strings.Contains(decrights, "HdrPost") {
				tconf["Action_button_post"] = "true" // config for post button

			} else {
				tconf["Action_button_post"] = "false" // config for post button

			}
		} else {
			tconf["Action_button_post"] = "false"
		}

		if arr_data[30] == `Posted` {
			if strings.Contains(decrights, "HdrCancel") {
				tconf["Action_button_cancel"] = "true" // config for cancel button
				tconf["Action_button_save"] = "false"
				tconf["Add_upload_button"] = "true"

			} else {
				tconf["Action_button_cancel"] = "false" // config for cancel button
				tconf["Add_upload_button"] = "true"
			}
		} else {
			tconf["Action_button_cancel"] = "false"
		}

		if arr_data[30] == `Cancelled` {

			tconf["Action_button_cancel"] = "false" // config for cancel button
			tconf["Action_button_post"] = "false"
			tconf["Action_button_save"] = "false"

		}
		//tconf["Add_upload_button_url"] = "/transaction/loan_release/HaddHandler?rights=" + rights + "&fromtype=la&la_id=" + tconf["h_id"]

		//arr_data_itemclass := datatables.DataList(`select id,dbo.SIS_Itemclass_Name(1, 3, id) [itemclass_name] from sis_itemclass`  )
		//arr_data_supplier := datatables.DataList(`select id,dbo.SIS_Supplier_Name(1, 3, id) [supplier_name] from SIS_Supplier`  )
		//arr_data_itemclass := datatables.DataList(`SIS_ItemClass_Get 2,`+str_OrgID  )
		//arr_data_supplier := datatables.DataList(`SIS_Supplier_Get 2,`+str_OrgID  )

		arr_data_partner := datatables.DataList(`exec partners_get 2,` + str_OrgID + `,''`)
		arr_data_trantype := datatables.DataList(`exec TranType_Get 3 ,0`)
		arr_data_bank := datatables.DataList(`exec Banks_Get 1, ` + str_OrgID + `,''`)
		arr_data_cvto := datatables.DataList(`exec Mas_GLCV_Get_droplist 'add',0`)
		arr_data_jvto := datatables.DataList(`exec Mas_GLJV_Get_droplist 'add', 0`)

		tmpl := template.New("Hedit.html")
		//tmpl := template.New("Hedit.html").Funcs(local_FuncMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/savings_withdrawal/Hedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arr_data_partner, arr_data_trantype, arr_data_bank, arr_data_cvto, arr_data_jvto})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}
}

func HViewHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	rights := r.URL.Query().Get("rights")
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	decrights := config.Decrypt(config.Encrypt_key, rights)

	//rights :="rights"

	if r.Method == "GET" {

		username, _ := login.Get_account_info(r)

		tconf := make(map[string]string)
		tconf["Add_upload_button"] = "false"
		tconf["Add_post_button"] = "false"
		tconf["Add_cancel_button"] = "false"
		tconf["h_id"] = r.URL.Query().Get("h_id")
		if tconf["h_id"] > "0" {
		} else {
			tconf["h_id"] = r.URL.Query().Get("hdr_id")
		}
		tconf["rights"] = rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID
		//tconf["LoanConso_url"] = "/transaction/loan_consolidation/details/DListHandler?rights=" + rights + "&hdr_id=" + tconf["h_id"]
		tconf["Add_post_button_url"] = "/transaction/savings_withdrawal/post?rights=" + rights + "&hdr_id=" + tconf["h_id"]
		tconf["Add_cancel_button_url"] = "/transaction/savings_withdrawal/cancel?rights=" + rights + "&hdr_id=" + tconf["h_id"]
		tconf["Add_new_button_url"] = "/transaction/savings_withdrawal/HaddHandler?rights=" + rights
		//tconf["data"] = datatables.DataList(sql)
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )

		//_,session_user_id 	:= login.Get_account_info(r)

		type Data struct {
			Rights      string
			Tconf       map[string]string
			Ls_partner  [][]string
			Ls_trantype [][]string
			Ls_bank     [][]string
			Ls_cvto     [][]string
			Ls_jvto     [][]string
			//Arr_tag [][]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		arr_data := datatables.Data_row(`exec SW_Get 1,'` + tconf["h_id"] + `'`)

		fmt.Println(arr_data)
		tconf["id"] = arr_data[0]
		tconf["trandate"] = arr_data[1]
		tconf["member"] = arr_data[2]
		tconf["partner"] = arr_data[3]
		tconf["trantype"] = arr_data[4]
		tconf["amt_requested"] = arr_data[5]
		tconf["bank"] = arr_data[6]
		tconf["cheque"] = arr_data[7]
		tconf["acct"] = arr_data[8]
		tconf["amt_issued"] = arr_data[9]
		tconf["date_issued"] = arr_data[10]
		tconf["cv"] = arr_data[11]
		tconf["remarks"] = arr_data[12]
		tconf["posted"] = arr_data[13]
		tconf["posted_by"] = arr_data[14]
		tconf["cancelled"] = arr_data[15]
		tconf["cancelled_by"] = arr_data[16]
		tconf["last_update"] = arr_data[17]
		tconf["amt_generated"] = arr_data[18]
		tconf["textfiledownloaded"] = arr_data[19]
		tconf["textfiledownloaded_by"] = arr_data[20]
		tconf["cv_date"] = arr_data[21]
		tconf["cvto"] = arr_data[22]
		tconf["jvto"] = arr_data[23]
		tconf["searchtext"] = arr_data[24]
		tconf["org"] = arr_data[25]
		tconf["added"] = arr_data[26]
		tconf["added_by"] = arr_data[27]
		tconf["updated"] = arr_data[28]
		tconf["updated_by"] = arr_data[29]
		tconf["Status"] = arr_data[30]
		tconf["MemberName"] = arr_data[31]
		tconf["PartnerName"] = arr_data[32]
		tconf["TrantypeName"] = arr_data[33]
		tconf["bank_name"] = arr_data[34]

		tconf["Action_button_save"] = "true"
		if arr_data[30] == `For Posting` {
			if strings.Contains(decrights, "HdrPost") {
				tconf["Action_button_post"] = "true" // config for post button

			} else {
				tconf["Action_button_post"] = "false" // config for post button

			}
		} else {
			tconf["Action_button_post"] = "false"
		}

		if arr_data[30] == `Posted` {
			if strings.Contains(decrights, "HdrCancel") {
				tconf["Action_button_cancel"] = "true" // config for cancel button
				tconf["Action_button_save"] = "false"
				tconf["Add_upload_button"] = "true"

			} else {
				tconf["Action_button_cancel"] = "false" // config for cancel button
				tconf["Add_upload_button"] = "true"
			}
		} else {
			tconf["Action_button_cancel"] = "false"
		}

		if arr_data[30] == `Cancelled` {

			tconf["Action_button_cancel"] = "false" // config for cancel button
			tconf["Action_button_post"] = "false"
			tconf["Action_button_save"] = "false"

		}
		//tconf["Add_upload_button_url"] = "/transaction/loan_release/HaddHandler?rights=" + rights + "&fromtype=la&la_id=" + tconf["h_id"]

		//arr_data_itemclass := datatables.DataList(`select id,dbo.SIS_Itemclass_Name(1, 3, id) [itemclass_name] from sis_itemclass`  )
		//arr_data_supplier := datatables.DataList(`select id,dbo.SIS_Supplier_Name(1, 3, id) [supplier_name] from SIS_Supplier`  )
		//arr_data_itemclass := datatables.DataList(`SIS_ItemClass_Get 2,`+str_OrgID  )
		//arr_data_supplier := datatables.DataList(`SIS_Supplier_Get 2,`+str_OrgID  )

		arr_data_partner := datatables.DataList(`exec partners_get 2,` + str_OrgID + `,''`)
		arr_data_trantype := datatables.DataList(`exec TranType_Get 3 ,0`)
		arr_data_bank := datatables.DataList(`exec Banks_Get 1, ` + str_OrgID + `,''`)
		arr_data_cvto := datatables.DataList(`exec Mas_GLCV_Get_droplist 'add',0`)
		arr_data_jvto := datatables.DataList(`exec Mas_GLJV_Get_droplist 'add', 0`)

		tmpl := template.New("Hview.html")
		//tmpl := template.New("Hedit.html").Funcs(local_FuncMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/savings_withdrawal/Hview.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arr_data_partner, arr_data_trantype, arr_data_bank, arr_data_cvto, arr_data_jvto})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}
}

func HDeleteHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	//Org_id := login.Get_session_org_id(r)
	//str_OrgID := strconv.Itoa(Org_id)
	if r.Method == "GET" {
		tconf := make(map[string]string)
		tconf["h_id"] = r.URL.Query().Get("h_id")
		tconf["delete_url"] = "/transaction/savings_withdrawal/HDeleteHandler"
		tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/savings_withdrawal/modal_delete_loghdr.html"); err != nil {
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
		//_ ,err, _,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Delete', '`+session_username+`', `+r.Form["h_id"][0],3)
		//_ ,err, _,_ := config.Ap_sql(`Supplier_Save 'Delete','programmer',1, ` +r.Form["h_id"][0] ,3)
		//_ ,err, _,_ := config.Ap_sql(`SIS_ItemClass_Save 'Delete','programmer',` +r.Form["h_id"][0]  ,3)
		arr_data := datatables.Data_row(`SW_SaveDelete '` + username + `' , ` + r.Form["h_id"][0] + ``)
		/*
			fmt.Println(arr_data)
			tconf["item_id"] = 			arr_data[0]
			tconf["item_org"] = 		arr_data[1]
			tconf["item_name"] = 		arr_data[2]
			tconf["item_parent"] = 		arr_data[3]
			tconf["item_parent_name"] = arr_data[4]
			tconf["item_remarks"] = 	arr_data[5] */
		js, err := json.Marshal(arr_data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

		/*if err != nil {
		        fmt.Println(err.Error)
		        http.Error(w, err.Error(), http.StatusInternalServerError)
		        panic(err.Error())
			}*/
	}
}

func Show_savings_ledger(w http.ResponseWriter, r *http.Request) {
	member := r.URL.Query().Get("member")

	fmt.Println(`LedgerS_Display 1,` + member + `,0`)
	arr_data := datatables.DataList(`LedgerS_Display 1,` + member + `,0`)

	type Data struct {
		//Tconf 		map[string]string
		Arr_data [][]string
	}
	var str_page string
	var str_page_dir string
	str_page = `member_ledgerS.html`
	str_page_dir = `report/member_ledgerS.html`

	tmpl := template.New(str_page).Funcs(helper.FuncMap)
	var err error
	if tmpl, err = tmpl.ParseFiles(str_page_dir); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w, Data{arr_data})
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func HPostHandler(w http.ResponseWriter, r *http.Request) {

	hdr_id := r.URL.Query().Get("hdr_id")
	username, _ := login.Get_account_info(r)
	r.ParseForm()

	fmt.Println(`GOOD TO POST`)
	arr_raw_data := datatables.DataList(`SW_Post '` + username + `',` + hdr_id)
	fmt.Println(arr_raw_data)

	js, err := json.Marshal(arr_raw_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func HCancelHandler(w http.ResponseWriter, r *http.Request) {

	hdr_id := r.URL.Query().Get("hdr_id")
	username, _ := login.Get_account_info(r)
	r.ParseForm()

	fmt.Println(`GOOD TO POST`)
	arr_raw_data := datatables.DataList(`SW_Cancel '` + username + `',` + hdr_id)
	fmt.Println(arr_raw_data)

	js, err := json.Marshal(arr_raw_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func LPostHandler(w http.ResponseWriter, r *http.Request) {

	login.Session_renew(w, r)

	if r.Method == "GET" {
		tconf := make(map[string]string)
		tconf["h_id"] = r.URL.Query().Get("hdr_id")
		tconf["post_url"] = "/transaction/savings_withdrawal/LPostHandler?hdr_id=" + tconf["h_id"]
		tmpl := template.New("modal_post_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/savings_withdrawal/modal_post_loghdr.html"); err != nil {
			fmt.Println(err)
		}
		err = tmpl.Execute(w, tconf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {

		r.ParseForm()

		hdr_id := r.URL.Query().Get("hdr_id")
		username, _ := login.Get_account_info(r)

		fmt.Println(`GOOD TO POST`)
		arr_raw_data := datatables.DataList(`SW_Post '` + username + `',` + hdr_id)
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
