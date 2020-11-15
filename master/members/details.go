package item_class

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	config "../../config"
	datatables "../../datatables"
	login "../../login"
)

func DListHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	//fmt.Fprint(w,"mars test here")
	//fmt.Fprint(w, r.URL.Query().Get("rights"))
	rights := r.URL.Query().Get("rights")
	//customer :=r.URL.Query().Get("customer")
	//rights :="static here"
	fmt.Println("rights here")
	fmt.Println(rights)

	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	hdr_id := r.URL.Query().Get("hdr_id")
	//customer_name := r.URL.Query().Get("customer_name")
	//location_name := r.URL.Query().Get("location_name")

	last_start := r.URL.Query().Get("last_start")
	last_length := r.URL.Query().Get("last_length")
	last_search := r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key, rights)
	fmt.Println("dec rights here")
	fmt.Println(dec_rights)
	tconf := make(map[string]string)

	tconf["hdr_id"] = r.URL.Query().Get("ID")
	//LBR_hdr := LBR_OTHdr_Get_id(tconf["hdr_id"])
	//LBR_hdr_ID_str:=strconv.Itoa( LBR_hdr.ID )

	ln := r.URL.Query().Get("length")
	tconf["Panel_name"] = " Contribution "
	//["+ LBR_hdr_ID_str+"]"
	tconf["Add_new_button"] = "false"
	tconf["Rights"] = rights
	tconf["Add_back_button"] = "true"
	tconf["Add_new_button_details"] = "false"
	tconf["upload_button"] = "false"
	//tconf["Add_new_button_url"] ="/inventory/customer/location/DaddHandler?rights="+rights+"&hdr_id="+tconf["hdr_id"]
	tconf["Action_baseURL"] = "/master/members/contribution?rights=" + rights
	tconf["ActionCol_param"] = "ID"
	//tconf["hdr_id"] = ""
	tconf["Interpret_button"] = "false"
	tconf["Add_post_button"] = "false"
	tconf["Add_post_button_url"] = "/inventory/customer/location/post?rights=" + rights + "&hdr_id=" + tconf["hdr_id"]

	tconf["Add_cancel_button"] = "false"
	tconf["Add_cancel_button_url"] = "/inventory/customer/location/cancel?rights=" + rights + "&hdr_id=" + tconf["hdr_id"]
	tconf["withtotal"] = "true"
	tconf["Archive"] = "false"

	tconf["level"] = "2"
	tconf["pageLength"] = last_length
	tconf["displayStart"] = last_start
	tconf["last_search"] = last_search
	if last_length == "" {
		tconf["pageLength"] = "10"
		tconf["displayStart"] = "0"
	}
	var header = []string{"Action", "ID", "Trantype", "Limit", "Deductible", "Suspended", "Remarks"} //column header
	if ln == "" {

		tconf["ismoreinfo"] = "false" // to description info

		/* arr_data_loccust := datatables.Data_row(`sis_locationcustomer_get 2 ,'`+tconf["hdr_id"]+`'`  )
		fmt.Println(arr_data)

		tconf["customer_id"] = 						arr_data_loccust[0]
		tconf["customer_name"] = 					arr_data_loccust[1] */

		tconf["Add_new_button_url"] = "/master/members/contribution/DaddHandler?rights=" + rights + "&hdr_id=" + tconf["hdr_id"]

		// content access here
		// rights for tables
		if strings.Contains(dec_rights, "DtlAdd") {
			tconf["Add_new_button_details"] = "true"
		}
		if strings.Contains(dec_rights, "DUpload") {
			tconf["upload_button"] = "true"

		}

		/*if( strings.Contains(LBR_hdr.Status , "Posted") ){
		 		tconf["Add_new_button_customers"] ="false"
		 		tconf["upload_button"] ="false"
		 		tconf["Add_post_button"] ="false"
		 		tconf["Interpret_button"] ="false"
		 		if( strings.Contains(dec_rights, "DCancel") ){
		 			tconf["Add_cancel_button"] ="true"
				}
		 	}
		 	if( strings.Contains(LBR_hdr.Status , "Cancelled") ){
		 		tconf["Add_new_button_customers"] ="false"
		 		tconf["upload_button"] ="false"
		 		tconf["Add_post_button"] ="false"
		 		tconf["Interpret_button"] ="false"
		 		tconf["Add_cancel_button"] ="false"
		 	}*/
		//end rights for tables

		//end of content accesss

		//arr_data_tag := datatables.DataList(`SIS_LocOrdTags_List 1,`+str_OrgID+` ,`+ tconf["hdr_id"] )
		//arr_data_refID := datatables.DataList(`SIS_TranRef_List 0, `+str_OrgID+`, 6,`+ tconf["hdr_id"] )

		var ajaxURLdata = "/master/members/contribution?rights=" + rights + "&hdr_id=" + tconf["hdr_id"]
		fmt.Println("/master/members/contribution?rights=" + rights + "&hdr_id=" + tconf["hdr_id"])
		type TblConf struct {
			Headercol []string
			AjaxUrl   string
			Tempconf  map[string]string
			//Arr_tag [][]string
			//Arr_refID [][]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("master/members/table.html"); err != nil {
			fmt.Println(err)
		}
		//err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf,arr_data_tag,arr_data_refID})
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
		tconf["ActionCol_edit_url"] = "/master/members/contribution/DEditHandler?rights=" + rights + "&hdr_id=" + hdr_id + "&h_id=" // config for edit click
		tconf["ActionCol_delete"] = "false"                                                                                         // config for delete click
		tconf["ActionCol_delete_url"] = "/inventory/customer/location/DDeleteHandler?h_id="
		tconf["ActionCol_detail"] = "false"    // config for customers click
		tconf["ActionCol_add_child"] = "false" // config for addchild click
		tconf["style_size"] = "12"
		tconf["style_font"] = ""
		tconf["format_number_col"] = "Amount"
		/*tconf["col_img"] ="false"
		 	tconf["name_with_img"] ="true"
		 	tconf["col_img_num"] ="1"
		 	tconf["img_id_col"] ="6"
			tconf["col_hidden1"] ="0"*/
		tconf["dr"] = dr //draw
		//geting total record etc
		// rights for tables
		if strings.Contains(dec_rights, "HEdit") {
			tconf["ActionCol_edit"] = "true"
		}
		if strings.Contains(dec_rights, "HDelete") {
			tconf["ActionCol_delete"] = "true" // config for delete click

		}
		if strings.Contains(dec_rights, "Hcontribution") {
			tconf["ActionCol_detail"] = "true" // config for delete click

		}

		if strings.Contains(dec_rights, "DtlEdit") {
			tconf["ActionCol_edit"] = "true"

		}

		status := r.URL.Query().Get("status")
		if status == `Cancelled` {
			tconf["ActionCol_edit"] = "false"
			tconf["ActionCol_delete"] = "false"
		} else if status == `Posted` {
			tconf["ActionCol_edit"] = "false"
			tconf["ActionCol_delete"] = "false"
		}
		//end rights for tables
		//_,session_user_id := login.Get_account_info(r)

		tconf["sql_total"] = `Member_Trantypes_List 1, ` + str_OrgID + `,` + hdr_id + `,  ` + str + `,  ` + ln + `, '` + sr + `'`
		tconf["sql_data"] = `Member_Trantypes_List 0, ` + str_OrgID + `,` + hdr_id + `,  ` + str + `,  ` + ln + `, '` + sr + `'`
		fmt.Println(tconf["sql_data"])
		tconf["org"] = str_OrgID
		datatables.DatatableretArray(w, tconf)

	}
}

func DaddHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	username, _ := login.Get_account_info(r)
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	//username := "static"
	//ses_uID := 1
	//str_ses_uID :=strconv.Itoa(ses_uID)
	rights := r.URL.Query().Get("rights")
	hdr_id := r.URL.Query().Get("hdr_id")
	isbarcode := r.URL.Query().Get("isbarcode")

	assignmentID := r.URL.Query().Get("customer")

	fmt.Println("customer ", assignmentID)

	if r.Method == "GET" {
		fmt.Println(`this is hdr id=`, hdr_id)

		tconf := make(map[string]string)
		tconf["username"] = username
		tconf["Org_id"] = str_OrgID
		tconf["hdr_id"] = hdr_id
		tconf["rights"] = rights
		tconf["customer"] = assignmentID

		type Data struct {
			Rights      string
			Tconf       map[string]string
			Ls_trantype [][]string
		}

		//arr_data_customer := datatables.DataList(`SIS_Location_Get 5,`+str_OrgID+` ,'`+tconf["hdr_id"]+`'`   )

		arr_data_trantype := datatables.DataList(`exec TranType_Get 3, 0`)

		var err error
		var html_file string = ``
		var html_filename string = ``

		if isbarcode == `1` {
			html_filename = "Dadd_barcode.html"
			html_file = "master/members/Dadd_barcode.html"
		} else {
			html_filename = "Dadd.html"
			html_file = "master/members/Dadd.html"
		}

		tmpl := template.New(html_filename)
		if tmpl, err = tmpl.ParseFiles(html_file); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arr_data_trantype})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}

	} else {

		r.ParseForm()

		limit := r.Form["limit"][0]
		limit = strings.Replace(limit, ",", "", -1)

		fmt.Println(`Member_Trantypes_SaveAdd '` + username + `',` + str_OrgID + `,` + r.Form["hdr_id"][0] + `,` + r.Form["trantype"][0] + `,` + limit + `,` + r.Form["deductible"][0] + `,` + r.Form["isSuspended"][0] + `,'` + r.Form["remarks"][0] + `'`)
		arr_raw_data := datatables.Data_row(`Member_Trantypes_SaveAdd '` + username + `',` + str_OrgID + `,` + r.Form["hdr_id"][0] + ` ,` + r.Form["trantype"][0] + `,` + limit + `,` + r.Form["deductible"][0] + `,` + r.Form["isSuspended"][0] + `,'` + r.Form["remarks"][0] + `'`)
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

func DAddTagHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		item_id := r.Form["item_id"][0]
		username, _ := login.Get_account_info(r)
		Org_id := login.Get_session_org_id(r)
		str_OrgID := strconv.Itoa(Org_id)
		var returnData []string
		for key, _ := range r.Form["tag"] {
			tag := r.Form["tag"][key]
			value_input := r.Form["value_input"][key]
			remarks := r.Form["remarks"][key]

			if value_input != "" {
				arr_data_itemtag := datatables.Data_row(`SIS_itemtags_save 'Update', '` + username + `',` + str_OrgID + `,` + item_id + `, ` + tag + `, '` + value_input + `', '` + remarks + `'`)
				fmt.Println(arr_data_itemtag)
				returnData = append(returnData, arr_data_itemtag[0])
				returnData = append(returnData, arr_data_itemtag[1])
				if strings.Contains(arr_data_itemtag[1], `Error`) {
					returnData = append(returnData, `HasError`)
				}
			}
		}
		js, err := json.Marshal(returnData)
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

	hdr_id := r.URL.Query().Get("hdr_id")
	assignmentID := r.URL.Query().Get("customer")

	fmt.Println("customer ", assignmentID)

	if r.Method == "GET" {

		tconf := make(map[string]string)
		tconf["h_id"] = r.URL.Query().Get("h_id")
		tconf["hdr_id"] = hdr_id
		tconf["rights"] = rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID
		tconf["customer"] = assignmentID

		assignmentID := r.URL.Query().Get("customer")
		fmt.Println("assignmentID ", assignmentID)

		//tconf["data"] = datatables.DataList(sql)
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )

		//_,session_user_id 	:= login.Get_account_info(r)

		type Data struct {
			Rights      string
			Tconf       map[string]string
			Ls_trantype [][]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		arr_data := datatables.Data_row(`Member_Trantypes_Get 1, '` + tconf["h_id"] + `'`)
		//arr_data_customer := datatables.DataList(`SIS_Customer_Get 3,`+str_OrgID+` ,'`+tconf["hdr_id"]+`'` )

		arr_data_trantype := datatables.DataList(`exec TranType_Get 3, 0`)

		fmt.Println(arr_data)
		tconf["id"] = arr_data[0]
		tconf["org"] = arr_data[1]
		tconf["member"] = arr_data[2]
		tconf["trantype"] = arr_data[3]
		tconf["limit"] = arr_data[4]
		tconf["deductible"] = arr_data[5]
		tconf["suspended"] = arr_data[6]
		tconf["remarks"] = arr_data[7]
		tconf["added"] = arr_data[8]
		tconf["added_by"] = arr_data[9]
		tconf["updated"] = arr_data[10]
		tconf["updated_by"] = arr_data[11]
		tconf["MemberName"] = arr_data[12]
		tconf["TranTypeName"] = arr_data[13]
		tconf["suspended_trigger"] = arr_data[14]

		tmpl := template.New("Dedit.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("master/members/Dedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf, arr_data_trantype})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	} else {
		r.ParseForm()

		fmt.Println(`Member_Trantypes_SaveEdit '` + username + `',` + r.Form["id"][0] + `, ` + str_OrgID + ` ,` + r.Form["hdr_id"][0] + `,` + r.Form["trantype"][0] + `,` + r.Form["limit"][0] + `,` + r.Form["deductible"][0] + `,` + r.Form["isSuspended"][0] + `,'` + r.Form["remarks"][0] + `'`)
		arr_raw_data := datatables.Data_row(`Member_Trantypes_SaveEdit '` + username + `',` + r.Form["id"][0] + `, ` + str_OrgID + ` ,` + r.Form["hdr_id"][0] + `,` + r.Form["trantype"][0] + `,` + r.Form["limit"][0] + `,` + r.Form["deductible"][0] + `,` + r.Form["isSuspended"][0] + `,'` + r.Form["remarks"][0] + `'`)
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
		tconf["delete_url"] = "/inventory/customer/location/DDeleteHandler"
		tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/customer/modal_delete_loghdr.html"); err != nil {
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
		arr_data := datatables.Data_row(`SIS_LocationCustomer_Save 'Delete','` + username + `',` + r.Form["h_id"][0])
		js, err := json.Marshal(arr_data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
