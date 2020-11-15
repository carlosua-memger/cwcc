package item_class

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
	tconf["Panel_name"] = "Members Admin List"
	tconf["Add_new_button"] = "false"
	tconf["Rights"] = rights
	tconf["Add_back_button"] = "false"
	tconf["Add_new_button_details"] = "false"
	tconf["upload_button"] = "false"
	tconf["Add_new_button_url"] = "/master/members_admin/HaddHandler?rights=" + rights
	//tconf["Action_baseURL"] ="/master/branch?rights="+rights
	tconf["Action_baseURL"] = "/master/members_admin/location?rights=" + rights //application link
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
		var ajaxURLdata = "/master/members_admin?rights=" + rights
		type TblConf struct {
			Headercol []string
			AjaxUrl   string
			Tempconf  map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("master/members_admin/table.html"); err != nil {
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
		tconf["ActionCol_edit_url"] = "/master/members_admin/HEditHandler?rights=" + rights + "&h_id=" // config for edit click
		tconf["ActionCol_delete"] = "false"                                                            // config for delete click
		tconf["ActionCol_delete_url"] = "/master/members/HDeleteHandler?h_id="
		tconf["ActionCol_detail"] = "false"    // config for details click
		tconf["ActionCol_add_child"] = "false" // config for addchild click
		//tconf["ActionCol_customer_url"] ="/inventory/location_items/customers?rights="+rights+"&h_id="

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
			Conf map[string]string
			//Ls_bank [][]string
			//Ls_membertype [][]string
			//Ls_status [][]string
			//Ls_branch [][]string
			Ls_role [][]string

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

		//arr_ls_bank := datatables.DataList(`exec banks_get 1, '`+str_OrgID+`'`   )
		//arr_ls_membertype := datatables.DataList(`exec membertype_get 1, '' ` )
		//arr_ls_status := datatables.DataList(`exec member_status_get 1, '' ` )
		//arr_ls_branches := datatables.DataList(`exec branches_get 1, '' ` )
		arr_ls_role := datatables.DataList(`exec sysrole_get  1, '' `)
		//ls_itemclass:=datatables.DataList(`exec SIS_Supplier_Get 2, `+str_OrgID+``)
		//arr_data_courier := datatables.DataList(`SIS_Courier_Get 2,`+str_OrgID  )
		//arr_data_itemtag := datatables.DataList(`sis_custags_list 1,`+str_OrgID+` ,`+ tconf["h_id"] )
		//arr_data_orgmem := datatables.DataList(`exec OrgMem_Get 1,`+str_OrgID+` , 0`  )
		//fmt.Println(ls_itemclass)

		tmpl := template.New("Hedit.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("master/members_admin/Hedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{tconf, arr_ls_role})
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

func EmailHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	//rights :=r.URL.Query().Get("rights")
	//Org_id :=login.Get_session_org_id(r)
	//str_OrgID :=strconv.Itoa(Org_id)
	//username, _ := login.Get_account_info(r)
	//rights :="rights"

	//str_hdrID :=strconv.Itoa(int(hdr_id))
	//strconv.Itoa(int(time.Nanoseconds()))
	r.ParseForm()

	//Org_id :=login.Get_session_org_id(r)
	//str_OrgID :=strconv.Itoa(Org_id)

	//fmt.Println(`cc email here : `,r.Form["fname"][0]);

	e := email.NewEmail()
	e.From = config.Email_id
	e.To = []string{r.Form["username_user"][0]}
	e.Bcc = []string{config.Email_bcc}
	fmt.Println(`cc email here : `, r.Form["username_user"][0])

	member := r.Form["fname"][0]
	password := r.Form["password"][0]
	username_user := r.Form["username_user"][0]

	e.Subject = `Welcome to MyCommonWealth`
	e.HTML = []byte(`Hi ` + member + `,<br><br>Welcome to MyCommonWealth, to get started here is your username and temporary password, kindly change your password immediately.<br><br><b>Username: </b>` + username_user + `<br><b>Password: </b>` + password + `<br><br>Click <a href="https://member.cwcc.ph/">here</a> to login your account.<br>Thank you.<br><b><i>This is an automated email, do not reply to this email.</b></i><br><i><b>Powered by Memgertech Inc.</b></i>`)

	retErr := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", e.From, config.Email_password, "smtp.gmail.com"))
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
