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

	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)

	last_start := r.URL.Query().Get("last_start")
	last_length := r.URL.Query().Get("last_length")
	last_search := r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key, rights)
	//dec_rights :="static here"
	tconf := make(map[string]string)
	ln := r.URL.Query().Get("length")
	tconf["Panel_name"] = "Members List"
	tconf["Add_new_button"] = "false"
	tconf["Rights"] = rights
	tconf["Add_back_button"] = "false"
	tconf["Add_new_button_details"] = "false"
	tconf["upload_button"] = "false"
	tconf["Add_new_button_url"] = "/master/members/HaddHandler?rights=" + rights
	//tconf["Action_baseURL"] ="/master/branch?rights="+rights
	tconf["Action_baseURL"] = "/master/members/contribution?rights=" + rights //application link
	tconf["ActionCol_param"] = "ID"
	tconf["hdr_id"] = "0"
	tconf["Interpret_button"] = "false"
	tconf["Add_post_button"] = "false"
	tconf["Add_post_button_url"] = ""

	tconf["Add_cancel_button"] = "false"
	tconf["Add_cancel_button_url"] = ""

	//tconf["username"] = username
	tconf["org_id"] = str_OrgID

	tconf["level"] = "1"
	tconf["pageLength"] = last_length
	tconf["displayStart"] = last_start
	tconf["last_search"] = last_search
	if last_length == "" {
		tconf["pageLength"] = "10"
		tconf["displayStart"] = "0"
	}
	var header = []string{"Action", "ID", "Name", "IssuedIDNo", "TIN", "SSS", "GSIS", "Email", "Mobile No.", "Address", "Remarks"} //column header
	if ln == "" {
		// content access here
		if strings.Contains(dec_rights, "HdrAdd") {
			tconf["Add_new_button"] = "true"
		}
		if strings.Contains(dec_rights, "DtlAdd") {
			tconf["Add_new_button"] = "true"
		}

		//end of content accesss
		var ajaxURLdata = "/master/members?rights=" + rights
		type TblConf struct {
			Headercol []string
			AjaxUrl   string
			Tempconf  map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("master/members/table.html"); err != nil {
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
		tconf["ActionCol"] = "true"                // config for action column
		tconf["ActionCol_param"] = "ID"            // config for parameter of action
		tconf["ActionCol_edit"] = "false"          // config for edit click
		tconf["ActionCol_edit_is_modal"] = "false" // config for edit click
		//tconf["ActionCol_edit_url"] ="/timekeeping/overtime_logs/OvertimeLogsHeaderEdit?rights="+rights+"&h_id=" // config for edit click
		tconf["ActionCol_edit_url"] = "/master/members/HEditHandler?rights=" + rights + "&h_id=" // config for edit click
		tconf["ActionCol_delete"] = "false"                                                      // config for delete click
		tconf["ActionCol_delete_url"] = "/master/members/HDeleteHandler?h_id="
		tconf["ActionCol_detail"] = "false"    // config for details click
		tconf["ActionCol_add_child"] = "false" // config for addchild click
		//tconf["ActionCol_customer_url"] ="/inventory/location_items/customers?rights="+rights+"&h_id="
		tconf["ActionCol_contribution"] = "true"
		tconf["style_size"] = "12"
		tconf["style_font"] = ""
		tconf["format_number_col"] = "Amount"

		tconf["org_id"] = str_OrgID

		tconf["dr"] = dr //draw
		//geting total record etc
		// rights for tables
		if strings.Contains(dec_rights, "HdrEdit") {
			tconf["ActionCol_edit"] = "true"
		}
		if strings.Contains(dec_rights, "HdrDelete") {
			tconf["ActionCol_delete"] = "true" // config for delete click

		}
		if strings.Contains(dec_rights, "HdrDetails") {
			tconf["ActionCol_detail"] = "true" // config for delete click

		}
		if strings.Contains(dec_rights, "DtlAdd") {
			tconf["Add_new_button"] = "true"
		}
		//end rights for tables
		_, session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id)

		fmt.Println(session_user_id_str)

		tconf["sql_total"] = `members_list 1, ` + str_OrgID + `,  ` + str + `,  ` + ln + `, '` + sr + `'`
		tconf["sql_data"] = `members_list 0, ` + str_OrgID + `,  ` + str + `,  ` + ln + `, '` + sr + `'`

		datatables.DatatableretArray(w, tconf)

	}
}

type RetMessage struct {
	RetNo  int
	RetMsg string
}

func DoAdd(branch string, date string, remarks string, username string) (bool, int) {
	//rowret ,err, _,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Add', '`+username+`', 0, '`+trandate+`', `+lbr_assign+`,  '`+remarks+`'`,1)
	rowret, err, _, _ := config.Ap_sql(`exec DailySumHdr_Save 'Add', 'Mars', 0, 1, `+branch+`, '`+date+`', '`+remarks+`'`, 1)
	if err != nil {
		panic(err.Error())
	}
	var r RetMessage
	if err != nil {
		fmt.Println(err.Error)
		panic(err.Error())
	}
	for rowret.Next() {
		err = rowret.Scan(&r.RetNo, &r.RetMsg)
		if err != nil {
			panic(err.Error())
		}
		r = RetMessage{r.RetNo, r.RetMsg}
	}

	if strings.Contains(r.RetMsg, "Success") {
		return true, r.RetNo
	} else {
		return false, r.RetNo
	}

}

var local_FuncMap = template.FuncMap{
	"Sql_list": func(s_qry string, org_id string, tag_id string) [][]string {
		fmt.Println(`sis_custags_get 1, ` + org_id + `, 0 ,` + tag_id)
		retdata := datatables.DataList(`sis_custags_get 1, ` + org_id + `, 0 ,` + tag_id)
		return retdata
	},
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

		//arr_data := datatables.DataList(`select id,dbo.SIS_Itemclass_Name(1, 3, id) [itemclass_name] from sis_itemclass`  )
		//arr_data := datatables.DataList(`SIS_ItemClass_Get 2,`+str_OrgID  )
		//arr_data_supplier := datatables.DataList(`select id,dbo.SIS_Supplier_Name(1, 3, id) [supplier_name] from SIS_Supplier`  )
		//arr_data_courier := datatables.DataList(`SIS_Courier_Get 2,`+str_OrgID  )
		//arr_data_supplier := datatables.DataList(`SIS_Supplier_Get 2,`+str_OrgID  )
		//arr_data_itemtag := datatables.DataList(`sis_custags_list 1,`+str_OrgID+` , 0`  )
		//arr_data_orgmem := datatables.DataList(`exec OrgMem_Get 1,`+str_OrgID+` , 0`  )
		type Data struct {
			Rights string
			Tconf  map[string]string
			//Ls_supplier [][]string
			//Arr_itemtag [][]string
			//Ls_courier [][]string
			//Ls_orgmem [][]string
		}
		tmpl := template.New("Hadd.html").Funcs(local_FuncMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("master/members/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{rights, tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}

	} else {

		r.ParseForm()
		add_status, lastinsertedID := DoAdd(r.Form["branch"][0], r.Form["date"][0], r.Form["remarks"][0], username)

		if add_status {
			str_lastinsertedID := strconv.Itoa(lastinsertedID)
			profile := Profile{"Data Successfully added ", true, str_lastinsertedID}
			js, err := json.Marshal(profile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("Data Successfully added")
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

/*func HAddTagHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method =="POST" {
		r.ParseForm()
		item_id := r.Form["item_id"][0]
		username, _ := login.Get_account_info(r)
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		var returnData[] string
		for key ,_ := range r.Form["tag"] {
			tag := r.Form["tag"][key]
			value_input := r.Form["value_input"][key]


			if(value_input!=""){
	        	arr_data_itemtag := datatables.Data_row(`SIS_custags_save 'Update', '`+username+`',`+str_OrgID+`,`+item_id+`, `+tag+`, '`+value_input+`', ''` )
	        	fmt.Println(arr_data_itemtag)
	        	returnData = append(returnData, arr_data_itemtag[0])
	        	returnData = append(returnData, arr_data_itemtag[1])
				if(strings.Contains(arr_data_itemtag[1] , `Error`)){
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
}*/

//edit here

type LBR_OTHdr struct {
	/* ID 			int
	Status 		string
	Trandate 	interface{}
	Lbr_assign 	int
	Remarks 	interface{} */

	ID            int
	org           int
	name          string
	contactperson interface{}
	TIN           interface{}
	address       interface{}
	area          interface{}
	email         interface{}
	contactno     interface{}
	creditlimit   interface{}
	courier       interface{}
	courier_name  interface{}
	remarks       interface{}
	member        int
}

func LBR_OTHdr_Get_id(Hdr_id string) LBR_OTHdr {
	//db_raw ,err, _,_ := config.Ap_sql(`LBR_OTHdr_Get 1 ,`+Hdr_id,1)
	db_raw, err, _, _ := config.Ap_sql(`SIS_Customer_Get 1, 3,`+Hdr_id, 1)

	if err != nil {
		panic(err.Error())
	}
	var r LBR_OTHdr
	for db_raw.Next() {
		err = db_raw.Scan(&r.ID, &r.org, &r.name, &r.contactperson, &r.TIN, &r.address, &r.area, &r.email, &r.contactno, &r.creditlimit, &r.courier, &r.courier_name, &r.remarks, &r.member)
		if err != nil {
			panic(err.Error())
		}
	}
	return r
}

func HEditHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	rights := r.URL.Query().Get("rights")
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	username, ses_uID := login.Get_account_info(r)
	str_ses_uID := strconv.Itoa(ses_uID)
	fmt.Println(str_ses_uID)
	fmt.Println(username)

	//rights :="static"

	if r.Method == "GET" {

		//Org_id :=login.Get_session_org_id(r)
		//str_OrgID :=strconv.Itoa(Org_id)
		//username, _ := login.Get_account_info(r)

		tconf := make(map[string]string)
		tconf["h_id"] = r.URL.Query().Get("h_id")
		tconf["rights"] = rights
		tconf["Org_id"] = str_OrgID
		tconf["username_sess"] = username

		//tconf["data"] = datatables.DataList(sql)
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )

		_, session_user_id := login.Get_account_info(r)
		// session_user_id := 1  //static here
		session_user_id_str := strconv.Itoa(session_user_id)
		fmt.Println(session_user_id_str)

		type Data struct {
			Conf          map[string]string
			Ls_bank       [][]string
			Ls_membertype [][]string
			Ls_status     [][]string
			Ls_branch     [][]string
			Ls_role       [][]string
			Ls_partner    [][]string

			//Arr_itemtag  [][]string
			//Ls_courier [][]string
			//Ls_orgmem [][]string
		}

		//arr_data := datatables.Data_row(`exec Members_Get 1, '`+str_OrgID+`', '`+tconf["h_id"]+`'`  )
		arr_data := datatables.Data_row(`exec Members_Get 1, '` + str_OrgID + `', '` + tconf["h_id"] + `'`)

		fmt.Println(arr_data)
		tconf["id"] = arr_data[0]
		tconf["id_no"] = arr_data[1]
		tconf["status"] = arr_data[2]
		tconf["status_date"] = arr_data[3]
		tconf["branch"] = arr_data[4]
		tconf["fname"] = arr_data[5]
		tconf["mname"] = arr_data[6]
		tconf["lname"] = arr_data[7]
		tconf["sname"] = arr_data[8]
		tconf["sex"] = arr_data[9]
		tconf["nationality"] = arr_data[10]
		tconf["stat"] = arr_data[11]
		tconf["birthday"] = arr_data[12]
		tconf["birthplace"] = arr_data[13]
		tconf["spousename"] = arr_data[14]
		tconf["bank"] = arr_data[15]
		tconf["acct"] = arr_data[16]
		tconf["tin"] = arr_data[17]
		tconf["sss"] = arr_data[18]
		tconf["gsis"] = arr_data[19]
		tconf["maiden_name"] = arr_data[20]
		tconf["present_add"] = arr_data[21]
		tconf["present_stat"] = arr_data[22]
		tconf["other_add"] = arr_data[23]
		tconf["other_stat"] = arr_data[24]
		tconf["mailing_add"] = arr_data[25]
		tconf["phone_no"] = arr_data[26]
		tconf["cell_no"] = arr_data[27]
		tconf["email"] = arr_data[28]
		tconf["educ_attain"] = arr_data[29]
		tconf["emp_stat"] = arr_data[30]
		tconf["occupation"] = arr_data[31]
		tconf["fund_source"] = arr_data[32]
		tconf["emp"] = arr_data[33]
		tconf["gross"] = arr_data[34]
		tconf["emp_add"] = arr_data[35]
		tconf["emp_no"] = arr_data[36]
		tconf["emp_ind"] = arr_data[37]
		tconf["emp_year"] = arr_data[38]
		tconf["prev_emp"] = arr_data[39]
		tconf["prev_year"] = arr_data[40]
		tconf["dependents"] = arr_data[41]
		tconf["spo_emp_stat"] = arr_data[42]
		tconf["spo_emp"] = arr_data[43]
		tconf["spo_gross"] = arr_data[44]
		tconf["card_bank"] = arr_data[45]
		tconf["vehicles"] = arr_data[46]
		tconf["lend1"] = arr_data[47]
		tconf["lend1_tran"] = arr_data[48]
		tconf["lend1_amt"] = arr_data[49]
		tconf["lend2"] = arr_data[50]
		tconf["lend2_tran"] = arr_data[51]
		tconf["lend2_amt"] = arr_data[52]
		tconf["lend3"] = arr_data[53]
		tconf["lend3_tran"] = arr_data[54]
		tconf["lend3_amt"] = arr_data[55]
		tconf["remarks"] = arr_data[56]
		tconf["terminated"] = arr_data[57]
		tconf["lasttouch"] = arr_data[58]
		tconf["lasttouched"] = arr_data[59]
		tconf["searchtext"] = arr_data[60]
		tconf["savings_deduction"] = arr_data[61]
		tconf["fname_"] = arr_data[62]
		tconf["mname_"] = arr_data[63]
		tconf["lname_"] = arr_data[64]
		tconf["sname_"] = arr_data[65]
		tconf["reg_to_textinquiry"] = arr_data[66]
		tconf["username_user"] = arr_data[67]
		tconf["password_user"] = arr_data[68]
		tconf["enabled"] = arr_data[69]
		tconf["disabled"] = arr_data[70]
		tconf["sysrole"] = arr_data[71]
		tconf["membertype"] = arr_data[72]
		tconf["postcode"] = arr_data[73]
		tconf["partner"] = arr_data[74]
		tconf["partner_id"] = arr_data[75]
		tconf["membership_ref"] = arr_data[76]
		tconf["termination_ref"] = arr_data[77]

		arr_ls_bank := datatables.DataList(`exec banks_get 1, '` + str_OrgID + `'`)
		arr_ls_membertype := datatables.DataList(`exec membertype_get 1, '' `)
		arr_ls_status := datatables.DataList(`exec member_status_get 1, '' `)
		arr_ls_branches := datatables.DataList(`exec branches_get 1, '' `)
		arr_ls_role := datatables.DataList(`exec sysrole_get  1, '' `)
		arr_data_partner := datatables.DataList(`exec partners_get 2,` + str_OrgID + `,''`)
		//ls_itemclass:=datatables.DataList(`exec SIS_Supplier_Get 2, `+str_OrgID+``)
		//arr_data_courier := datatables.DataList(`SIS_Courier_Get 2,`+str_OrgID  )
		//arr_data_itemtag := datatables.DataList(`sis_custags_list 1,`+str_OrgID+` ,`+ tconf["h_id"] )
		//arr_data_orgmem := datatables.DataList(`exec OrgMem_Get 1,`+str_OrgID+` , 0`  )
		//fmt.Println(ls_itemclass)

		tmpl := template.New("Hedit.html").Funcs(local_FuncMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("master/members/Hedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{tconf, arr_ls_bank, arr_ls_membertype, arr_ls_status, arr_ls_branches, arr_ls_role, arr_data_partner})
		//err1 := tmpl.Execute( w,&Data{tconf,arr_data_itemtag,arr_data_courier,arr_data_orgmem} )
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	} else {
		session_username, _ := login.Get_account_info(r)
		r.ParseForm()
		//_ ,err, ex_stat,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Edit', '`+session_username+`', `+r.Form["h_id"][0] +`, '`+r.Form["trandate"][0]+`', `+r.Form["lbr_assign"][0] +`, '`+r.Form["remarks"][0]+`'`,3)
		//_ ,err, ex_stat,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Edit', '`+session_username+`', `+r.Form["h_id"][0] +`, '`+r.Form["trandate"][0]+`', `+r.Form["lbr_assign"][0] +`, '`+r.Form["remarks"][0]+`'`,3)

		_, err, ex_stat, _ := config.Ap_sql(`exec DailySumHdr_Save 'Add', '`+session_username+`', `+r.Form["h_id"][0]+` , 1, `+r.Form["branch"][0]+`, '`+r.Form["date"][0]+`', '`+r.Form["remarks"][0]+`'`, 3)

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

func HDeleteHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	if r.Method == "GET" {

		tconf := make(map[string]string)
		tconf["h_id"] = r.URL.Query().Get("h_id")
		tconf["delete_url"] = "/inventory/customer/HDeleteHandler"
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

		username, _ := login.Get_account_info(r)
		r.ParseForm()
		//_ ,err, _,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Delete', '`+session_username+`', `+r.Form["h_id"][0],3)
		//_ ,err, _,_ := config.Ap_sql(`Supplier_Save 'Delete','programmer',1, ` +r.Form["h_id"][0] ,3)
		//_ ,err, _,_ := config.Ap_sql(`SIS_ItemClass_Save 'Delete','programmer',` +r.Form["h_id"][0]  ,3)

		fmt.Println(`SIS_Customer_Save 'Delete','` + username + `' ,  ` + r.Form["h_id"][0] + ` , ` + str_OrgID + ``)

		arr_data := datatables.Data_row(`SIS_Customer_Save 'Delete','` + username + `' ,  ` + r.Form["h_id"][0] + ``)
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
