package item_class

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
	//login.Session_renew(w,r)	
	//fmt.Fprint(w,"mars test here")
	//fmt.Fprint(w, r.URL.Query().Get("rights"))
	//rights :=r.URL.Query().Get("rights")
	rights :=r.URL.Query().Get("rights")
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)

	last_start :=r.URL.Query().Get("last_start")
	last_length :=r.URL.Query().Get("last_length")
	last_search :=r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key,rights)		
	tconf := make(map[string]string)
	ln :=  r.URL.Query().Get("length")
	tconf["Panel_name"] ="Item Class Record" 
	tconf["Add_new_button"] ="false" 
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="false"
	tconf["Add_new_button_details"] ="false" 
	tconf["upload_button"] ="false"   
	tconf["Add_new_button_url"] ="/inventory/item_class/HaddHandler?rights="+rights 
	tconf["Action_baseURL"] ="/master/branch?rights="+rights
	tconf["ActionCol_param"] ="ID"
    tconf["hdr_id"] = "0"
	tconf["Interpret_button"] ="false"
	tconf["Add_post_button"] ="false"
	tconf["Add_post_button_url"] =""
	
	tconf["Add_cancel_button"] ="false"
	tconf["Add_cancel_button_url"] =""
	

	tconf["level"] ="1" 
	tconf["pageLength"] =last_length 
	tconf["displayStart"] =last_start
	tconf["last_search"] = last_search
	if last_length==""{
		tconf["pageLength"] ="10" 
		tconf["displayStart"] ="0"
	}
	var header = []string{"Action", "ID", "Name", "Remarks"}  //column header
	if ln=="" {
		// content access here
		if( strings.Contains(dec_rights, "HAdd") ){
			tconf["Add_new_button"] ="true"
			tconf["ActionCol_add_child"] ="true"
		}
		if( strings.Contains(dec_rights, "Dadd ") ){

		}
			
		//end of content accesss
		var ajaxURLdata = "/inventory/item_class?rights="+rights
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("inventory/item_class/table.html"); err != nil {
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
	 	tconf["ActionCol_edit_is_modal"] ="true" // config for edit click 
	 	//tconf["ActionCol_edit_url"] ="/timekeeping/overtime_logs/OvertimeLogsHeaderEdit?rights="+rights+"&h_id=" // config for edit click 
	 	tconf["ActionCol_edit_url"] ="/inventory/item_class/HEditHandler?rights="+rights+"&h_id=" // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/inventory/item_class/HDeleteHandler?h_id="
	 	tconf["ActionCol_detail"] ="false" // config for details click 
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
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
		//end rights for tables
		//_,session_user_id := login.Get_account_info(r)
		session_user_id := 1  //static here
		session_user_id_str := strconv.Itoa(session_user_id) 
		fmt.Println(session_user_id_str)
/*		fmt.Println("exec LBR_LogHdr_List 1, 2, 1, 1, 1, '"+sortColName+"', '', '"+sr+"'	" )
		fmt.Println(sortColName,str,sr)*/
		//fmt.Println(`exec LBR_OTHdr_List 1, 2, `+session_user_id_str+` , `+str+`, `+ln+`,  '`+sortColName+`', '', '`+sr+`'`)
		//tconf["sql_total"] = `exec LBR_OTHdr_List 1, 2, `+session_user_id_str+` , `+str+`, `+ln+`,  '`+sortColName+`', '', '`+sr+`'`
		//tconf["sql_data"]  = `exec LBR_OTHdr_List 0, 2, `+session_user_id_str+` , `+str+`, `+ln+`,  '`+sortColName+`', '', '`+sr+`'`
		
		//tconf["sql_total"] = `exec dailysumhdr_list 1, 0, 0 , `+str+`, `+ln+`, '`+sortcol+`','','`+sr+`'`
		//tconf["sql_data"]  = `exec dailysumhdr_list 0, 0, 0 , `+str+`, `+ln+`, '`+sortcol+`','','`+sr+`'`
		
		

		tconf["sql_total"] = `SIS_ItemClass_List 1, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"]  = `SIS_ItemClass_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		
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
		tconf := make(map[string]string)
		tconf["parentID"] =r.URL.Query().Get("parentID")
		tconf["username"] = username
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		tconf["Org_id"] = str_OrgID

		type Data struct {
			Rights	string	
			Tconf 	map[string]string
		}    
     	tmpl := template.New("Hadd.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/item_class/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights , tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {

		r.ParseForm()
		add_status, lastinsertedID := DoAdd(r.Form["branch"][0],r.Form["date"][0],r.Form["remarks"][0],username)

		if(add_status){
			str_lastinsertedID :=strconv.Itoa(lastinsertedID) 
			profile := Profile{"Data Successfully added ",true ,str_lastinsertedID}
			js, err := json.Marshal(profile)
			if err != nil {
		    	http.Error(w, err.Error(), http.StatusInternalServerError)
		    	return
	  		}
		  	fmt.Println("Data Successfully added")
		  	w.Header().Set("Content-Type", "application/json")
		  	w.Write(js)
		  		
		}else{
			profile := Profile{"Oppsss something went wrong",false,""}
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
func LBR_OTHdr_Get_id( Hdr_id string ) Dailysumhdr_get {
	//db_raw ,err, _,_ := config.Ap_sql(`LBR_OTHdr_Get 1 ,`+Hdr_id,1)
	db_raw ,err, _,_ := config.Ap_sql(`dailysumhdr_get 1,`+Hdr_id,1)

    if err != nil {
    	panic(err.Error())
    }
    var r Dailysumhdr_get
    for db_raw.Next() {
        err = db_raw.Scan(&r.ID, &r.Branch,&r.Docdate,&r.Remarks)
        if err != nil {
            panic(err.Error())
        }
	}
	return r
} 
  
func HEditHandler(w http.ResponseWriter, r *http.Request) {
	//login.Session_renew(w,r)
	//rights :=r.URL.Query().Get("rights")
	rights :=r.URL.Query().Get("rights")
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	if r.Method =="GET" {	


		tconf := make(map[string]string)
		tconf["h_id"] =r.URL.Query().Get("h_id")
		tconf["rights"]=rights
		username, _ := login.Get_account_info(r)
		tconf["username"] = username
		tconf["rights"]=username
		

		type Data struct {
			Conf map[string]string
			Ls_itemClass [][]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		arr_data := datatables.Data_row(`exec SIS_ItemClass_Get 1, `+str_OrgID+`, '`+tconf["h_id"]+`'`  )
		
		fmt.Println(arr_data) 
		tconf["item_id"] = 			arr_data[0]
		tconf["item_org"] = 		arr_data[1] 
		tconf["item_name"] = 		arr_data[2] 
		tconf["item_parent"] = 		arr_data[3] 
		tconf["item_parent_name"] = arr_data[4] 
		tconf["item_remarks"] = 	arr_data[5] 
		
		ls_itemclass:=datatables.DataList(`exec SIS_ItemClass_Get 2, `+str_OrgID+``)
		fmt.Println(ls_itemclass)

		tmpl := template.New("Hedit.html").Funcs(config.FuncMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/item_class/Hedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute( w,&Data{tconf,ls_itemclass} )
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		//session_username := `static`
		r.ParseForm()
		//_ ,err, ex_stat,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Edit', '`+session_username+`', `+r.Form["h_id"][0] +`, '`+r.Form["trandate"][0]+`', `+r.Form["lbr_assign"][0] +`, '`+r.Form["remarks"][0]+`'`,3)
		//_ ,err, ex_stat,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Edit', '`+session_username+`', `+r.Form["h_id"][0] +`, '`+r.Form["trandate"][0]+`', `+r.Form["lbr_assign"][0] +`, '`+r.Form["remarks"][0]+`'`,3)
		
		_ ,err, ex_stat,_ := config.Ap_sql(`exec DailySumHdr_Save 'Add', 'Mars', `+r.Form["h_id"][0] +` , 1, `+r.Form["branch"][0]+`, '`+r.Form["date"][0]+`', '`+r.Form["remarks"][0]+`'`,3)
		
		//exec LBR_LogHdr_Save 'Edit', 'Ian', 3, '11 Jul 2016', 1, 'logfile abc', 'device abc', 'remarks abc'
		if err != nil {
			        fmt.Println(err.Error)
			        panic(err.Error())
			    }
			//lastinsertedID, _ := res.LastInsertId()
			if ex_stat == true{
				profile := Profile{"Data Successfully Update ",true ,""}
				js, err := json.Marshal(profile)
				if err != nil {
			    	http.Error(w, err.Error(), http.StatusInternalServerError)
			    	return
		  		}
			  	fmt.Println("Data Successfully Edited")
			  	w.Header().Set("Content-Type", "application/json")
			  	w.Write(js)	
			}else{
				profile := Profile{"Oppsss something went wrong",false,""}
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
	//login.Session_renew(w,r)
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["h_id"] = r.URL.Query().Get("h_id")
	   	tconf["delete_url"] = "/inventory/item_class/HDeleteHandler"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/item_class/modal_delete_loghdr.html"); err != nil {
			fmt.Println(err)
		}
		err = tmpl.Execute(w,tconf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} 
	}else{
		//session_username,_ := login.Get_account_info(r)	
		r.ParseForm()
		username, _ := login.Get_account_info(r)
		//_ ,err, _,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Delete', '`+session_username+`', `+r.Form["h_id"][0],3)
		//_ ,err, _,_ := config.Ap_sql(`Supplier_Save 'Delete','programmer',1, ` +r.Form["h_id"][0] ,3)
		//_ ,err, _,_ := config.Ap_sql(`SIS_ItemClass_Save 'Delete','programmer',` +r.Form["h_id"][0]  ,3)
		arr_data := datatables.Data_row(`SIS_ItemClass_Save 'Delete','`+username+`',` +r.Form["h_id"][0]  )
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