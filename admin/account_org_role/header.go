package account_org_role

import (
	
		"net/http"
		"fmt"
		"text/template"
		config "../../config"
		"strconv"
		"encoding/json"
		datatables "../../datatables"
		login "../../login"
		"strings"
)

type Profile struct {
  	Message string  // error status output or validated data
  	Status bool //pass or fail
  	ID string
}

func HListHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)	
	//fmt.Fprint(w,"mars test here")
	//fmt.Fprint(w, r.URL.Query().Get("rights"))
	rights :=r.URL.Query().Get("rights")
	//rights :="static here"

	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	
	last_start :=r.URL.Query().Get("last_start")
	last_length :=r.URL.Query().Get("last_length")
	last_search :=r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key,rights)		
	//dec_rights :="static here"
	tconf := make(map[string]string)
	ln :=  r.URL.Query().Get("length")
	tconf["Panel_name"] ="Acount Organization And Role" 
	tconf["Add_new_button"] ="false" 
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="false"
	tconf["Add_new_button_details"] ="false" 
	tconf["upload_button"] ="false"   
	tconf["Add_new_button_url"] ="/administrator/account_org_role/HaddHandler?rights="+rights 
	tconf["Action_baseURL"] ="/administrator/accountrole?rights="+rights
	tconf["ActionCol_param"] ="ID"
    tconf["hdr_id"] = "0"
	tconf["Interpret_button"] ="false"
	tconf["Add_post_button"] ="false"
	tconf["Add_post_button_url"] =""
	
	tconf["Add_cancel_button"] ="false"
	tconf["Add_cancel_button_url"] =""
		
	//tconf["username"] = username
	tconf["org_id"] = str_OrgID

	tconf["level"] ="1" 
	tconf["pageLength"] =last_length 
	tconf["displayStart"] =last_start
	tconf["last_search"] = last_search
	if last_length==""{
		tconf["pageLength"] ="10" 
		tconf["displayStart"] ="0"
	}
	var header = []string{"Action", "ID", "Status","Join Date","Member Name","Role"}  //column header
	if ln=="" {
		// content access here
		if( strings.Contains(dec_rights, "HAdd") ){
			tconf["Add_new_button"] ="true"
		}
		if( strings.Contains(dec_rights, "Dadd") ){
			tconf["Add_new_button"] ="true"
		}
			
		//end of content accesss
		var ajaxURLdata = "/administrator/account_org_role?rights="+rights
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("admin/account_org_role/table.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		dr := r.URL.Query().Get("draw")

		sr := r.URL.Query().Get("search[value]")
		fmt.Println(sr)
		srt1,_ :=strconv.Atoi(r.URL.Query().Get("start"))
	    ln1,_ :=strconv.Atoi(r.URL.Query().Get("length")) 
	    srt1 = (srt1 / ln1) + 1
	    str := strconv.Itoa(srt1)
	    sortcol :=  r.URL.Query().Get("order[0][column]")	
		if sortcol=="0"{
			sortcol="1"
		}
		srtcolindexInt,_ :=strconv.Atoi(sortcol)
	    var sortColName string
		for key , value := range header{
			if ( srtcolindexInt==key){
				sortColName = value
			}
		}
		fmt.Println(sortColName)
		tconf["ActionCol"] ="true" // config for action column
	 	tconf["ActionCol_param"] ="ID" // config for parameter of action 
	 	tconf["ActionCol_edit"] ="false" // config for edit click
	 	tconf["ActionCol_edit_is_modal"] ="false" // config for edit click 
	 	//tconf["ActionCol_edit_url"] ="/timekeeping/overtime_logs/OvertimeLogsHeaderEdit?rights="+rights+"&h_id=" // config for edit click 
	 	tconf["ActionCol_edit_url"] ="/administrator/account_org_role/HEditHandler?rights="+rights+"&h_id=" // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/administrator/account_org_role/HDeleteHandler?h_id="
	 	tconf["ActionCol_detail"] ="false" // config for details click 
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	tconf["ActionCol_roles"] ="false" // config for account role click 
	 	tconf["ActionCol_roles_url"] ="/administrator/accountrole?rights="+rights // config for roles
		tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"

	 	tconf["org_id"] = str_OrgID

		tconf["dr"]=dr  //draw
		//geting total record etc 
		// rights for tables 	
		if( strings.Contains(dec_rights, "HEdit") ){
			tconf["ActionCol_edit"] ="true"
		}
		if( strings.Contains(dec_rights, "HDelete") ){
			tconf["ActionCol_delete"] ="true" // config for delete click 
	 	
		}
		if( strings.Contains(dec_rights, "HDetails") ){
			tconf["ActionCol_detail"] ="true" // config for delete click 
	 	
		}
		if( strings.Contains(dec_rights, "Hroles") ){
			tconf["ActionCol_roles"] ="true" // config for roles click 
	 	
		}
		if( strings.Contains(dec_rights, "Dadd") ){
			tconf["Add_new_button"] ="true"
		}
		//end rights for tables
		_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 

		fmt.Println(session_user_id_str)
		
		tconf["sql_total"] = `OrgMem_List 1, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"]  = `OrgMem_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		
		datatables.DatatableretArray(w,tconf)
		  	 
		
		
		}	    
}
type RetMessage struct {
	RetNo	int
    RetMsg	string
}
func DoAdd(branch string, date string, remarks string , username string) (bool ,int) {
	//rowret ,err, _,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Add', '`+username+`', 0, '`+trandate+`', `+lbr_assign+`,  '`+remarks+`'`,1)
	rowret ,err, _,_ := config.Ap_sql(`exec DailySumHdr_Save 'Add', 'Mars', 0, 1, `+branch+`, '`+date+`', '`+remarks+`'`,1)
	if err != nil {
		panic(err.Error())
	}
	var r RetMessage
	if err != nil {
        fmt.Println(err.Error)
        panic(err.Error())
    }
    for rowret.Next() {
	    err = rowret.Scan(&r.RetNo,&r.RetMsg)
	    if err != nil {
	        panic(err.Error())
	    }
	 	r = RetMessage{r.RetNo,r.RetMsg}
	}
	
	if( strings.Contains(r.RetMsg, "Success") ){
		return true ,r.RetNo
	}else{
		return false ,r.RetNo
	}
	
}

var local_FuncMap = template.FuncMap{
		"Sql_list": func(s_qry string , org_id string, tag_id string) [][]string{
			fmt.Println(`sis_suptags_get 1, `+org_id+`, 0 ,`+tag_id)
			retdata := datatables.DataList(`sis_suptags_get 1, `+org_id+`, 0 ,`+tag_id)
			return retdata
		},
	}

func HAddHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	username, ses_uID := login.Get_account_info(r)
	
	//username := "static"
	//ses_uID := 1
	str_ses_uID :=strconv.Itoa(ses_uID)
	rights :=r.URL.Query().Get("rights")
	fmt.Println(str_ses_uID)
	//fmt.Fprint(w,rights)
	
	if r.Method =="GET" {	
		
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		tconf := make(map[string]string)
		//tconf["parentID"] =r.URL.Query().Get("parentID")
		tconf["username"] =username
		tconf["Org_id"] = str_OrgID

		//arr_data := datatables.DataList(`select id,dbo.SIS_Itemclass_Name(1, 3, id) [itemclass_name] from sis_itemclass`  )
		//arr_data := datatables.DataList(`SIS_ItemClass_Get 2,`+str_OrgID  )
		//arr_data_supplier := datatables.DataList(`select id,dbo.SIS_Supplier_Name(1, 3, id) [supplier_name] from SIS_Supplier`  )

		//arr_data_supplier := datatables.DataList(`SIS_Supplier_Get 2,`+str_OrgID  )
		arr_data_roleperaccount := datatables.DataList(`exec account_roles_peraccount 0,` +str_OrgID  )
		fmt.Println(arr_data_roleperaccount)
		type Data struct {
			Rights	string	
			Tconf 	map[string]string
			Arr_RoleAccount [][]string
		}  
		tmpl := template.New("Hadd.html").Funcs(local_FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("admin/account_org_role/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights,tconf,arr_data_roleperaccount})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		r.ParseForm()
		memberid :=r.Form["membername"][0]
		//transdate :=r.Form["transdate"][0]
		OrgmemroleID := `0`
		role := r.Form["role"][0]
		
		/*fmt.Println(`OrgMem_Save 'Add','`+username+`',0,`+ str_OrgID + `,`+memberid)
		orgmem_save := datatables.DataList(`OrgMem_Save 'Add','`+username+`',0,`+ str_OrgID + `,`+memberid)
		fmt.Println(orgmem_save)*/


		fmt.Println(`OrgMemRol_Save 'Add','`+username+`',`+OrgmemroleID+`,`+ str_OrgID + `,`+memberid+ `,`+role)
		orgmemrole_save := datatables.DataList(`OrgMemRol_Save 'Add','`+username+`',`+OrgmemroleID+`,`+ str_OrgID + `,`+memberid+ `,`+role)
		fmt.Println(orgmemrole_save)

		js, err := json.Marshal(orgmemrole_save)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)	

	}		
}



//edit here




type LBR_OTHdr struct{
	ID 			int
	Status 		string
	Trandate 	interface{} 
	Lbr_assign 	int
	Remarks 	interface{}
}

type Dailysumhdr_get struct{
	ID 			int
	Branch 		interface{}
	Docdate 	interface{} 
	Remarks 	interface{}
}

func HEditHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	//rights :="static"
	username, _ := login.Get_account_info(r)
	if r.Method =="GET" {	


		//Org_id :=login.Get_session_org_id(r)
		//str_OrgID :=strconv.Itoa(Org_id)
		

		tconf := make(map[string]string)
		tconf["h_id"] =r.URL.Query().Get("h_id")
		tconf["Rights"]=rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID
		_,session_user_id 	:= login.Get_account_info(r)
		// session_user_id := 1  //static here
		session_user_id_str := strconv.Itoa(session_user_id) 
		fmt.Println(session_user_id_str)
		

		type Data struct {
			Tconf map[string]string
			//Ls_itemClass [][]string
			Arr_RoleAccount  [][]string
		}
		
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		fmt.Println(`orgmemrol_get `+tconf["h_id"]  )
		arr_data := datatables.Data_row(`orgmemrol_get `+tconf["h_id"]  )
		
		fmt.Println(arr_data) 

		tconf["orgmemroleID"] = 	arr_data[0]
		tconf["SysRole_id"] = 		arr_data[1] 
		tconf["SysRole_name"] = 	arr_data[2] 
		tconf["OrgMem_txdate"] = 	arr_data[3] 
		tconf["OrgMem_status"] = 	arr_data[4] 
		tconf["OrgMem_remarks"] = 	arr_data[5] 
		tconf["Member_id"] = 		arr_data[6] 
		tconf["Member_name"] = 		arr_data[7]

		
		//ls_itemclass:=datatables.DataList(`exec SIS_Supplier_Get 2, `+str_OrgID+``)
		arr_data_account_role := datatables.DataList(`exec account_roles_peraccount `  + tconf["h_id"]+ `,` +str_OrgID )
		fmt.Println(arr_data_account_role)

		tmpl := template.New("Hedit.html").Funcs(local_FuncMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("admin/account_org_role/Hedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute( w,&Data{tconf,arr_data_account_role} )
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{

		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		r.ParseForm()
		memberid :=r.Form["membername"][0]
		//transdate :=r.Form["transdate"][0]
		OrgmemroleID :=r.Form["OrgmemroleID"][0]

		role := r.Form["role"][0]
		
		fmt.Println(`OrgMemRol_Save 'Edit','`+username+`',`+OrgmemroleID+`,`+ str_OrgID + `,`+memberid+ `,`+role)
		orgmem_save := datatables.DataList(`OrgMemRol_Save 'Edit','`+username+`',`+OrgmemroleID+`,`+ str_OrgID + `,`+memberid+ `,`+role)
		fmt.Println(orgmem_save)

		js, err := json.Marshal(orgmem_save)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)	
	}	
}

func HDeleteHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {


	   	tconf := make(map[string]string)
	   	tconf["h_id"] = r.URL.Query().Get("h_id")
	   	tconf["delete_url"] = "/administrator/account_org_role/HDeleteHandler"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("admin/account_org_role/modal_delete_loghdr.html"); err != nil {
			fmt.Println(err)
		}
		err = tmpl.Execute(w,tconf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} 
	}else{
		
		username, _ := login.Get_account_info(r)
		r.ParseForm()
		fmt.Println(`OrgMemRol_Save 'Delete','`+username+`',` +r.Form["h_id"][0]  )
		arr_data := datatables.Data_row(`OrgMemRol_Save 'Delete','`+username+`',` +r.Form["h_id"][0]  )
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