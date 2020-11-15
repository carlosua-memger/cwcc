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
	tconf["Panel_name"] = "Member Mapping"
	tconf["Add_new_button"] = "false"
	tconf["Rights"] = rights
	tconf["Add_back_button"] = "false"
	tconf["Add_new_button_details"] = "false"
	tconf["upload_button"] = "false"
	tconf["Add_new_button_url"] = "/master/members_mapping/HaddHandler?rights=" + rights
	tconf["Action_baseURL"] = "/master/members_mapping/details?rights=" + rights
	tconf["ActionCol_param"] = "ID"
	tconf["hdr_id"] = "0"
	tconf["Interpret_button"] = "false"
	tconf["Add_post_button"] = "false"
	tconf["Add_post_button_url"] = ""
	tconf["Archive"] = "true" // is archive
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
	var header = []string{"Action", "ID", "Status", "Detected", "CHRIS ID", "CHRIS Code", "CHRIS Fullname", "CWCC Partner", "CWCC Member", "Remarks"} //column header
	if ln == "" {
		// content access here
		if strings.Contains(dec_rights, "HdrAdd") {
			tconf["Add_new_button"] = "true"
			//fmt.Println("run herer")
		}

		//end of content accesss
		var ajaxURLdata = "/master/members_mapping?rights=" + rights
		type TblConf struct {
			Headercol []string
			AjaxUrl   string
			Tempconf  map[string]string
		}
		tmpl := template.New("table.html").Funcs(helper.FuncMap)
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("master/members_mapping/table.html"); err != nil {
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
		tconf["ActionCol_edit_url"] = "/master/members_mapping/HEditHandler?rights=" + rights + "&h_id=" // config for edit click
		tconf["ActionCol_view"] = "false"                                                                // config for edit click
		tconf["ActionCol_view_is_modal"] = "false"                                                       // config for edit click
		//tconf["ActionCol_edit_url"] ="/timekeeping/overtime_logs/OvertimeLogsHeaderEdit?rights="+rights+"&h_id=" // config for edit click
		tconf["ActionCol_view_url"] = "/master/members_mapping/HViewHandler?rights=" + rights + "&h_id=" // config for edit click

		tconf["ActionCol_delete"] = "false" // config for delete click
		tconf["ActionCol_delete_url"] = "/master/members_mapping/HDeleteHandler?h_id="
		tconf["ActionCol_detail"] = "false"    // config for details click
		tconf["ActionCol_add_child"] = "false" // config for addchild click
		tconf["ActionCol_loanconsolidation_url"] = "/master/members_mapping/detals/DListHandler?h_id="
		tconf["style_size"] = "12"
		tconf["style_font"] = ""
		tconf["format_number_col"] = "Amount"
		tconf["dr"] = dr //draw
		//geting total record etc
		// rights for tables
		if strings.Contains(dec_rights, "HdrEdit") {
			tconf["ActionCol_edit"] = "true"
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
		tconf["sql_total"] = `EOD_CHRISMemberToMap_list 1, ` + str_OrgID + `,  ` + str + `,  ` + ln + `, '` + sr + `',` + archive
		//tconf["sql_data"]  = `SIS_Item_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"] = `EOD_CHRISMemberToMap_list 0, ` + str_OrgID + `,  ` + str + `,  ` + ln + `, '` + sr + `',` + archive

		datatables.DatatableretArray(w, tconf)

	}
}

type RetMessage struct {
	RetNo  int
	RetMsg string
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
		tconf["Add_post_button_url"] = "/master/members_mapping/post?rights=" + rights + "&hdr_id=" + tconf["h_id"]
		tconf["Add_cancel_button_url"] = "/master/members_mapping/cancel?rights=" + rights + "&hdr_id=" + tconf["h_id"]

		//tconf["data"] = datatables.DataList(sql)
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )

		//_,session_user_id 	:= login.Get_account_info(r)

		type Data struct {
			Rights     string
			Tconf      map[string]string
			Ls_partner [][]string

			//Arr_tag [][]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		arr_data := datatables.Data_row(`exec EOD_CHRISMemberToMap_Get 1,'` + tconf["h_id"] + `'`)

		fmt.Println(arr_data)
		tconf["id"] = arr_data[0]
		tconf["status"] = arr_data[1]
		tconf["detected"] = arr_data[2]
		tconf["employeeno"] = arr_data[3]
		tconf["employeecode"] = arr_data[4]
		tconf["fullname"] = arr_data[5]
		tconf["proposedpartner"] = arr_data[6]
		tconf["proposedmember"] = arr_data[7]
		tconf["proposedpartner_name"] = arr_data[8]
		tconf["proposedmember_name"] = arr_data[9]

		tconf["Action_button_save"] = "true"
		if arr_data[1] == `For Posting` {
			if strings.Contains(decrights, "HdrPost") {
				tconf["Action_button_post"] = "true" // config for post button

			} else {
				tconf["Action_button_post"] = "false" // config for post button

			}
		} else {
			tconf["Action_button_post"] = "false"
		}

		if arr_data[1] == `Posted` {
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

		if arr_data[1] == `Cancelled` {

			tconf["Action_button_cancel"] = "false" // config for cancel button
			tconf["Action_button_post"] = "false"
			tconf["Action_button_save"] = "false"

		}
		//tconf["Add_upload_button_url"] = "/transaction/loan_release/HaddHandler?rights=" + rights + "&fromtype=la&la_id=" + tconf["h_id"]
		tconf["CreateLR_url"] = "/transaction/loan_release/HaddHandler?rights=" + rights + "&fromtype=la&la_id=" + tconf["h_id"]
		tconf["LoanConso_url"] = "/transaction/loan_consolidation/details?rights=" + rights + "&hdr_id=" + tconf["h_id"] + "&member=" + tconf["la_member"]

		//arr_data_itemclass := datatables.DataList(`select id,dbo.SIS_Itemclass_Name(1, 3, id) [itemclass_name] from sis_itemclass`  )
		//arr_data_supplier := datatables.DataList(`select id,dbo.SIS_Supplier_Name(1, 3, id) [supplier_name] from SIS_Supplier`  )
		//arr_data_itemclass := datatables.DataList(`SIS_ItemClass_Get 2,`+str_OrgID  )
		//arr_data_supplier := datatables.DataList(`SIS_Supplier_Get 2,`+str_OrgID  )

		arrDataPartner := datatables.DataList(`exec partners_get 2,` + str_OrgID + `,''`)

		tmpl := template.New("Hedit.html")
		//tmpl := template.New("Hedit.html").Funcs(local_FuncMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("master/members_mapping/Hedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arrDataPartner})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	} else {
		//session_username := `static`
		r.ParseForm()
		//_ ,err, ex_stat,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Edit', '`+session_username+`', `+r.Form["h_id"][0] +`, '`+r.Form["trandate"][0]+`', `+r.Form["lbr_assign"][0] +`, '`+r.Form["remarks"][0]+`'`,3)
		//_ ,err, ex_stat,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Edit', '`+session_username+`', `+r.Form["h_id"][0] +`, '`+r.Form["trandate"][0]+`', `+r.Form["lbr_assign"][0] +`, '`+r.Form["remarks"][0]+`'`,3)

		_, err, ex_stat, _ := config.Ap_sql(`exec DailySumHdr_Save 'Add', 'Mars', `+r.Form["h_id"][0]+` , 1, `+r.Form["branch"][0]+`, '`+r.Form["date"][0]+`', '`+r.Form["remarks"][0]+`'`, 3)

		//exec LBR_LogHdr_Save 'Edit', 'Ian', 3, '11 Jul 2016', 1, 'logfile abc', 'device abc', 'remarks abc'
		if err != nil {
			fmt.Println(err.Error)
			panic(err.Error())
		}
		//lastinsertedID, _ := res.LastInsertId()
		if ex_stat == true {
			profile := Profile{"Data Successfully Update ", true, ""}
			js, err := json.Marshal(profile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("Data Successfully Edited")
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		} else {
			profile := Profile{"Oppsss something went wrong", false, ""}
			js, err := json.Marshal(profile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("Oppsss something went wrong")
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}

	}
}

func HPostHandler(w http.ResponseWriter, r *http.Request) {

	hdr_id := r.URL.Query().Get("hdr_id")
	username, _ := login.Get_account_info(r)
	r.ParseForm()

	fmt.Println(`GOOD TO POST`)
	arr_raw_data := datatables.DataList(`EOD_CHRISMemberToMap_Post '` + username + `',` + hdr_id)
	fmt.Println(arr_raw_data)

	js, err := json.Marshal(arr_raw_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
