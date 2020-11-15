package loanpayment

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
	tconf["Panel_name"] = "Resign Process"
	tconf["Add_new_button"] = "false"
	tconf["Rights"] = rights
	tconf["Add_back_button"] = "false"
	tconf["Add_new_button_details"] = "false"
	tconf["upload_button"] = "false"
	tconf["Add_new_button_url"] = "/transaction/resign_process/HaddHandler?rights=" + rights
	tconf["Action_baseURL"] = "/transaction/resign_process/details?rights=" + rights
	tconf["ActionCol_param"] = "ID"
	tconf["hdr_id"] = "0"
	tconf["Interpret_button"] = "false"
	tconf["Add_post_button"] = "false"
	tconf["Add_post_button_url"] = ""
	tconf["Archive"] = "true" // is archive
	tconf["member"] = ""
	tconf["MemberName"] = ""

	tconf["Add_cancel_button"] = "false"
	tconf["Add_clearance_button"] = "false"
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
	var header = []string{"Action", "ID", "Status", "Trandate", "Member", "Contribution Adjustment", "Loan Adjustment", "Process Type", "Remarks"} //column header
	if ln == "" {
		// content access here
		if strings.Contains(dec_rights, "HdrAdd") {
			tconf["Add_new_button"] = "true"
			//fmt.Println("run herer")
		}

		//end of content accesss
		var ajaxURLdata = "/transaction/resign_process?rights=" + rights
		type TblConf struct {
			Headercol []string
			AjaxUrl   string
			Tempconf  map[string]string
		}
		tmpl := template.New("table.html").Funcs(helper.FuncMap)
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("transaction/resign_process/table.html"); err != nil {
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
		tconf["ActionCol_edit_url"] = "/transaction/resign_process/HEditHandler?rights=" + rights + "&h_id=" // config for edit click
		tconf["ActionCol_view"] = "false"                                                                    // config for edit click
		tconf["ActionCol_view_is_modal"] = "false"                                                           // config for edit click
		//tconf["ActionCol_edit_url"] ="/timekeeping/overtime_logs/OvertimeLogsHeaderEdit?rights="+rights+"&h_id=" // config for edit click
		tconf["ActionCol_view_url"] = "/transaction/resign_process/HViewHandler?rights=" + rights + "&h_id=" // config for edit click

		tconf["ActionCol_delete"] = "false" // config for delete click
		tconf["ActionCol_delete_url"] = "/transaction/resign_process/HDeleteHandler?h_id="
		tconf["ActionCol_detail"] = "true"     // config for details click
		tconf["ActionCol_add_child"] = "false" // config for addchild click
		tconf["ActionCol_loanconsolidation_url"] = "/transaction/resign_process/details/DListHandler?h_id="
		tconf["style_size"] = "12"
		tconf["style_font"] = ""
		tconf["format_number_col"] = "Amount"
		tconf["dr"] = dr //draw
		//geting total record etc
		// rights for tables
		if strings.Contains(dec_rights, "HdrEdit") {
			tconf["ActionCol_edit"] = "true"
		}

		if strings.Contains(dec_rights, "HdrView") && archive == "1" {
			tconf["ActionCol_view"] = "true"
		}
		if strings.Contains(dec_rights, "HdrDelete") {
			tconf["ActionCol_delete"] = "true" // config for delete click

		}
		if strings.Contains(dec_rights, "HdrDetails") {
			tconf["ActionCol_detail"] = "true" // config for details click

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
		tconf["sql_total"] = `ARLR_List 1, ` + str_OrgID + `,  ` + str + `,  ` + ln + `, '` + sr + `',` + archive
		//tconf["sql_data"]  = `SIS_Item_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"] = `ARLR_List 0, ` + str_OrgID + `,  ` + str + `,  ` + ln + `, '` + sr + `',` + archive

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

		arr_data_proctype := datatables.DataList(`exec proctype_get 1,0`)

		//arr_data_tag := datatables.DataList(`exec SIS_SupOrdTags_List 1,`+str_OrgID+` , 0`  )
		type Data struct {
			Rights      string
			Tconf       map[string]string
			Ls_proctype [][]string
		}
		tmpl := template.New("Hadd.html")
		//tmpl := template.New("Hadd.html").Funcs(local_FuncMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/resign_process/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arr_data_proctype})
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

	//rights :="rights"

	if r.Method == "GET" {

		username, _ := login.Get_account_info(r)

		tconf := make(map[string]string)
		tconf["Add_post_button"] = "false"
		tconf["Add_cancel_button"] = "false"
		tconf["h_id"] = r.URL.Query().Get("h_id")
		tconf["rights"] = rights
		tconf["username"] = username
		tconf["Org_id"] = str_OrgID

		type Data struct {
			Rights      string
			Tconf       map[string]string
			Ls_proctype [][]string

			//Arr_tag [][]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		arr_data := datatables.Data_row(`exec ARLR_Get 1,` + tconf["h_id"])

		fmt.Println(arr_data)
		tconf["id"] = arr_data[0]
		tconf["trandate"] = arr_data[1]
		tconf["member"] = arr_data[2]
		tconf["AEHDR"] = arr_data[3]
		tconf["LAEHDR"] = arr_data[4]
		tconf["posted"] = arr_data[5]
		tconf["posted_by"] = arr_data[6]
		tconf["cancelled"] = arr_data[7]
		tconf["cancelled_by"] = arr_data[8]
		tconf["remarks"] = arr_data[9]
		tconf["last_update"] = arr_data[10]
		tconf["proctype"] = arr_data[11]
		tconf["textfiledownloaded"] = arr_data[12]
		tconf["textfiledownloaded_by"] = arr_data[13]
		tconf["searchtext"] = arr_data[14]
		tconf["org"] = arr_data[15]
		tconf["added"] = arr_data[16]
		tconf["added_by"] = arr_data[17]
		tconf["updated"] = arr_data[18]
		tconf["updated_by"] = arr_data[19]
		tconf["status"] = arr_data[20]
		tconf["MemberName"] = arr_data[21]
		tconf["ProcessType"] = arr_data[22]

		tconf["Action_button_save"] = "true"

		arr_data_proctype := datatables.DataList(`exec ProcType_Get 1, 0`)

		tmpl := template.New("Hedit.html")
		//tmpl := template.New("Hedit.html").Funcs(local_FuncMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/resign_process/Hedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arr_data_proctype})
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

	//rights :="rights"

	if r.Method == "GET" {

		username, _ := login.Get_account_info(r)

		tconf := make(map[string]string)
		tconf["Add_post_button"] = "false"
		tconf["Add_cancel_button"] = "false"
		tconf["h_id"] = r.URL.Query().Get("h_id")
		tconf["rights"] = rights
		tconf["username"] = username
		tconf["Org_id"] = str_OrgID

		type Data struct {
			Rights      string
			Tconf       map[string]string
			Ls_proctype [][]string

			//Arr_tag [][]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		arr_data := datatables.Data_row(`exec ARLR_Get 1,` + tconf["h_id"])

		fmt.Println(arr_data)
		tconf["id"] = arr_data[0]
		tconf["trandate"] = arr_data[1]
		tconf["member"] = arr_data[2]
		tconf["AEHDR"] = arr_data[3]
		tconf["LAEHDR"] = arr_data[4]
		tconf["posted"] = arr_data[5]
		tconf["posted_by"] = arr_data[6]
		tconf["cancelled"] = arr_data[7]
		tconf["cancelled_by"] = arr_data[8]
		tconf["remarks"] = arr_data[9]
		tconf["last_update"] = arr_data[10]
		tconf["proctype"] = arr_data[11]
		tconf["textfiledownloaded"] = arr_data[12]
		tconf["textfiledownloaded_by"] = arr_data[13]
		tconf["searchtext"] = arr_data[14]
		tconf["org"] = arr_data[15]
		tconf["added"] = arr_data[16]
		tconf["added_by"] = arr_data[17]
		tconf["updated"] = arr_data[18]
		tconf["updated_by"] = arr_data[19]
		tconf["status"] = arr_data[20]
		tconf["MemberName"] = arr_data[21]
		tconf["ProcessType"] = arr_data[22]

		tconf["Action_button_save"] = "true"

		arr_data_proctype := datatables.DataList(`exec ProcType_Get 1, 0`)

		tmpl := template.New("Hview.html")
		//tmpl := template.New("Hedit.html").Funcs(local_FuncMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/resign_process/Hview.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arr_data_proctype})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}
}

func HDeleteHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)

	if r.Method == "GET" {
		tconf := make(map[string]string)
		tconf["h_id"] = r.URL.Query().Get("h_id")
		tconf["delete_url"] = "/transaction/resign_process/HDeleteHandler"
		tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("transaction/resign_process/modal_delete_loghdr.html"); err != nil {
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
		arr_data := datatables.Data_row(`ARLR_SaveDelete '` + username + `' , ` + r.Form["h_id"][0] + ``)
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
