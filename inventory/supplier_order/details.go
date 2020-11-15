package supplier_order

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
		helper "../../helper"
		email "../../github.com/email-master"
		"net/smtp"
)




func DListHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)	
	//fmt.Fprint(w,"mars test here")
	//fmt.Fprint(w, r.URL.Query().Get("rights"))
	rights :=r.URL.Query().Get("rights")
	location :=r.URL.Query().Get("location")
		supplier :=r.URL.Query().Get("supplier")
	//rights :="static here"
	fmt.Println("rights here")
	fmt.Println(rights)

	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	hdr_id := r.URL.Query().Get("hdr_id")

		

	last_start :=r.URL.Query().Get("last_start")
	last_length :=r.URL.Query().Get("last_length")
	last_search :=r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key,rights)		
	fmt.Println("dec rights here")
	fmt.Println(dec_rights)
	tconf := make(map[string]string)
	tconf["hdr_id"] = r.URL.Query().Get("ID") 
	ln :=  r.URL.Query().Get("length")
	tconf["Panel_name"] ="Supplier Order Detail" 
	tconf["Add_new_button"] ="false" 
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="true"
	tconf["Add_new_button_details"] ="true" 
	tconf["upload_button"] ="false"   
	//tconf["Add_new_button_url"] ="/inventory/supplier_order/details/DaddHandler?rights="+rights+"&hdr_id="+tconf["hdr_id"]
	tconf["Action_baseURL"] ="/inventory/supplier_order/details?rights="+rights
	tconf["ActionCol_param"] ="ID"
    //tconf["hdr_id"] = ""
	tconf["Interpret_button"] ="false"
	tconf["Add_post_button"] ="false"
	tconf["Add_post_button_url"] ="/inventory/supplier_order/details/post?rights="+rights+"&hdr_id="+tconf["hdr_id"]
	tconf["ActionCol_void"] ="true" // config for void click 
	tconf["ActionCol_void_url"] ="/inventory/supplier_order/details/DVoidHandler?h_id="
	tconf["Add_cancel_button"] ="false"
	tconf["Add_cancel_button_url"] ="/inventory/supplier_order/details/cancel?rights="+rights+"&hdr_id="+tconf["hdr_id"]
	tconf["withtotal"] = "true" 
	tconf["Archive"] = "false"
	tconf["org"] =str_OrgID

	tconf["level"] ="2" 
	tconf["pageLength"] =last_length 
	tconf["displayStart"] =last_start
	tconf["last_search"] = last_search
	if last_length==""{
		tconf["pageLength"] ="10" 
		tconf["displayStart"] ="0"
	}
	var header = []string{"Action", "Item Name", "Qty", "Unit cost", "Total cost", "Remarks"}  //column header
	if ln=="" {
		
		tconf["ismoreinfo"] = "true" // to description info
		arr_data := datatables.Data_row(`SIS_SupOrdHdr_Get 1 ,  `+str_OrgID+` ,'`+tconf["hdr_id"]+`'`  )
		fmt.Println(arr_data) 
		
		tconf["suprec_hdr_id"] = 			arr_data[0]
		tconf["suprec_hdr_status"] = 		arr_data[1] 
		tconf["suprec_hdr_location"] = 		arr_data[2]
		tconf["suprec_hdr_location_name"] = arr_data[3] 
		tconf["suprec_hdr_trnsdate"] = 		arr_data[4] 
		tconf["suprec_hdr_supplier"] = 		arr_data[5]  
		tconf["suprec_hdr_supplier_name"] = arr_data[6]  
		tconf["suprec_hdr_term"] = 			arr_data[7]  
		tconf["suprec_hdr_remarks"] = 		arr_data[8]  
		tconf["Add_new_button_url"] ="/inventory/supplier_order/details/DaddHandler?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&customer="+tconf["suprec_hdr_customer"]+"&location="+tconf["suprec_hdr_location"]+"&supplier="+tconf["suprec_hdr_supplier"]
		tconf["Email_url"] ="/inventory/supplier_order/details/email?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&customer="+tconf["suprec_hdr_customer"]+"&location="+tconf["suprec_hdr_location"]

//tconf["Add_new_button_url"] ="/inventory/supplier_order/details/DaddHandler?rights="+rights+"&hdr_id="+tconf["hdr_id"]
	
		// content access here
		// rights for tables 	
		if( strings.Contains(dec_rights, "DAdd") ){
			tconf["Add_new_button_details"] ="true"
		}
		if( strings.Contains(dec_rights, "DUpload") ){
			tconf["upload_button"] ="true"  
	 	
		}
		if( strings.Contains(dec_rights, "DPost") ){
			if(arr_data[1]==`Empty` && arr_data[1]==`Posted`){
				tconf["Add_post_button"] ="false" 	
			}else{
				tconf["Add_post_button"] ="true" 
			}
		}

		if( strings.Contains(dec_rights, "DCancel") ){
			if(arr_data[1]==`Empty` || arr_data[1]==`Cancelled` || arr_data[1]==`For Posting` ){
				tconf["Add_cancel_button"] ="false" 	
			}else{
				tconf["Add_cancel_button"] ="true" 
			}
		}
		if(arr_data[1]==`Cancelled`){
			tconf["Add_post_button"] ="false" 	
			tconf["Add_new_button_details"] ="false"
			tconf["upload_button"] ="false"  
	 		tconf["Add_cancel_button"] ="false"
		
		}else if(arr_data[1]==`Posted`){
				tconf["Add_new_button_details"] ="false"
				tconf["Add_post_button"] ="false" 	
	 		//	tconf["Add_cancel_button"] ="true"

		
		}else if(arr_data[1]==`Empty`){
				tconf["Add_post_button"] ="false" 	
	 			tconf["Add_cancel_button"] ="false"

		}

		if( strings.Contains(dec_rights, "DInterpret") ){
	 		tconf["Interpret_button"] ="true"
	
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
		
		arr_data_tag := datatables.DataList(`SIS_SupOrdTags_List 1,`+str_OrgID+` ,`+ tconf["hdr_id"] )
		
		
		var ajaxURLdata = "/inventory/supplier_order/details?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&status="+tconf["suprec_hdr_status"]+"&location="+tconf["suprec_hdr_location"]+"&supplier="+tconf["suprec_hdr_supplier"]
		fmt.Println("/inventory/supplier_order/details?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&status="+ tconf["suprec_hdr_status"]+"&location="+tconf["suprec_hdr_location"])
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
			Arr_tag [][]string
		}
			tmpl := template.New("table.html").Funcs(helper.FuncMap)
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("inventory/supplier_order/table.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf,arr_data_tag})
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
	 	tconf["ActionCol_edit_is_modal"] ="true" // config for edit click  modal design
	 	//tconf["ActionCol_edit_url"] ="/timekeeping/overtime_logs/OvertimeLogsHeaderEdit?rights="+rights+"&h_id=" // config for edit click 
	 	tconf["ActionCol_edit_url"] =  "/inventory/supplier_order/details/DEditHandler?rights="+rights+"&location="+location+"&supplier="+supplier+"&h_id="+tconf["hdr_id"] // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/inventory/supplier_order/details/DDeleteHandler?h_id="
	 	tconf["ActionCol_detail"] ="false" // config for details click 
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
	 	tconf["col_img"] ="false"
	 	tconf["name_with_img"] ="true"
	 	tconf["col_img_num"] ="1"
	 	tconf["img_id_col"] ="5"
	 	tconf["col_hidden1"] ="0"
		tconf["dr"]=dr  //draw
		tconf["void"] ="7"
		//geting total record etc 
		// rights for tables 	
		if( strings.Contains(dec_rights, "DEdit") ){
			tconf["ActionCol_edit"] ="true"
		}
		if( strings.Contains(dec_rights, "DDelete") ){
			tconf["ActionCol_delete"] ="true" // config for delete click 
	 	
		}
		if( strings.Contains(dec_rights, "HDetails") ){
			tconf["ActionCol_detail"] ="true" // config for delete click 
	 	
		}
		
		status := r.URL.Query().Get("status")
	 	if(status==`Cancelled`){
			tconf["ActionCol_edit"] ="false"
			tconf["ActionCol_delete"] ="false" 
		}else if(status==`Posted`){
			tconf["ActionCol_edit"] ="false"
			tconf["ActionCol_delete"] ="false" 
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
		tconf["sql_total"] = `SIS_SupOrdDtl_List 1, `+str_OrgID+`,`+hdr_id+`,  `+str+`,  `+ln+`, '`+sr+`'`
		//tconf["sql_data"]  = `SIS_Item_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"]  = `SIS_SupOrdDtl_List 0, `+str_OrgID+`,`+hdr_id+`,  `+str+`,  `+ln+`, '`+sr+`'`
		fmt.Println(tconf["sql_data"])
		tconf["org"] = str_OrgID
		datatables.DatatableretArray(w,tconf)
		  	 
		
		
		}	    
}

func DAddHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	username, _ := login.Get_account_info(r)
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	//username := "static"
	//ses_uID := 1
	//str_ses_uID :=strconv.Itoa(ses_uID)
	rights :=r.URL.Query().Get("rights")
	hdr_id :=r.URL.Query().Get("hdr_id")
	isbarcode :=r.URL.Query().Get("isbarcode")



	assignmentID :=r.URL.Query().Get("location")

	supplier :=r.URL.Query().Get("supplier")

	fmt.Println ("location ",assignmentID )
	

	
	
	//fmt.Fprint(w,rights)
	
	if r.Method =="GET" {	
		fmt.Println(`this is hdr id=`,hdr_id)
		
		tconf := make(map[string]string)
		//tconf["parentID"] =r.URL.Query().Get("parentID")
		tconf["username"] =username
		tconf["Org_id"] = str_OrgID
		tconf["hdr_id"] = hdr_id
		tconf["rights"] = rights
		tconf["test_js"] = `alert("from webserver")`
		tconf["location"] = assignmentID
		tconf["supplier"] = supplier


		type Data struct {
			Rights	string	
			Tconf 	map[string]string
			
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
		var html_file string = ``
		var html_filename string = ``

		if (isbarcode == `1`){
			html_filename = "Dadd_barcode.html"
			html_file = "inventory/supplier_order/Dadd_barcode.html"
		}else{
			html_filename = "Dadd.html"
			html_file = "inventory/supplier_order/Dadd.html"
		}

		tmpl := template.New(html_filename)
		if tmpl, err = tmpl.ParseFiles(html_file); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights , tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {

		r.ParseForm()
		fmt.Println(`SIS_SupOrdDtl_Save 'Add', '`+username+`',0,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+`,`+r.Form["qty"][0]+`,`+r.Form["unit_cost"][0]+`,'`+r.Form["remarks"][0]+`'` )
		arr_raw_data := datatables.Data_row(`SIS_SupOrdDtl_Save 'Add', '`+username+`',0,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+`,`+r.Form["qty"][0]+`,`+r.Form["unit_cost"][0]+`,'`+r.Form["remarks"][0]+`'`)
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
			remarks := r.Form["remarks"][key]

			if(value_input!=""){
	        	arr_data_itemtag := datatables.Data_row(`SIS_itemtags_save 'Update', '`+username+`',`+str_OrgID+`,`+item_id+`, `+tag+`, '`+value_input+`', '`+remarks+`'` )
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
}

func DEditHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	username, _ := login.Get_account_info(r)
	//rights :="rights"



	assignmentID :=r.URL.Query().Get("location")

	supplier :=r.URL.Query().Get("supplier")

	fmt.Println ("location ",assignmentID )
	

	if r.Method =="GET" {	

		

		tconf := make(map[string]string)
		tconf["h_id"] =r.URL.Query().Get("h_id")
		tconf["rights"]=rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID
		tconf["location"] = assignmentID
		tconf["supplier"] = supplier
		
		
		assignmentID :=  r.URL.Query().Get("location") 
		fmt.Println ("assignmentID ",assignmentID )
		
		//tconf["data"] = datatables.DataList(sql)
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )
		
		//_,session_user_id 	:= login.Get_account_info(r)

		
		type Data struct {
			Rights string
			Tconf map[string]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		fmt.Println(`exec SIS_SupOrdDtl_Get 1, `+str_OrgID+`, '`+tconf["h_id"]+`'` )
		arr_data := datatables.Data_row(`exec SIS_SupOrdDtl_Get 1, `+str_OrgID+`, '`+tconf["h_id"]+`'`  )
		
		fmt.Println(arr_data) 
		tconf["suprec_hdr_id"] = 			arr_data[0]
		tconf["suprec_hdr_item"] = 			arr_data[1] 
		tconf["suprec_hdr_itemname"] = 		arr_data[2] 
		tconf["suprec_hdr_qty"] = 			arr_data[3]
		tconf["suprec_hdr_ucost"] = 		arr_data[4]
		tconf["suprec_hdr_tcost"] = 		arr_data[5] 
		tconf["suprec_hdr_remarks"] = 		arr_data[6]  
		tconf["suprec_hdr_hdr_id"] = 		arr_data[7]  


		

		tmpl := template.New("Dedit.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/supplier_order/Dedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute( w,&Data{rights,tconf} )
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		r.ParseForm()

		fmt.Println(`SIS_SupOrdDtl_Save 'Edit', '`+username+`',`+r.Form["id"][0]+`,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+`,`+r.Form["qty"][0]+`,`+r.Form["unit_cost"][0]+`,'`+r.Form["remarks"][0]+`'` )
		arr_raw_data := datatables.Data_row(`SIS_SupOrdDtl_Save 'Edit', '`+username+`',`+r.Form["id"][0]+`,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+`,`+r.Form["qty"][0]+`,`+r.Form["unit_cost"][0]+`,'`+r.Form["remarks"][0]+`'`)
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
	login.Session_renew(w,r)
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["h_id"] = r.URL.Query().Get("h_id")
	   	tconf["delete_url"] = "/inventory/supplier_order/details/DDeleteHandler"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/supplier_order/modal_delete_loghdr.html"); err != nil {
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
		arr_data := datatables.Data_row(`SIS_SupOrdDtl_Save 'Delete','`+username+`',` +r.Form["h_id"][0]  )
		js, err := json.Marshal(arr_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)	
	}	   	
}

func Ajax_getItem(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	getq := r.URL.Query().Get("q")
	getid := r.URL.Query().Get("id")
	isbaseonassignid :=r.URL.Query().Get("isbaseonassignid")
	assignmentID := r.URL.Query().Get("assignmentID")
	//supplier := r.URL.Query().Get("supplier")
	hdr_id := r.URL.Query().Get("hdr_id")
	fmt.Println(`is base sa assignID`,isbaseonassignid)
	fmt.Println(` assignID`,assignmentID)

	getq = getq + getid
		type Inner_item struct {
			Id string `json:"id"`
			Fullname string `json:"full_name"`
			Item_image string `json:"avatar"`
			

		}
		type JsondataEmpty struct {
	    	Total_count int `json:"total_count"`
	    	Incomplete_results bool `json:"incomplete_results"`
	    	Items []Inner_item `json:"items"`
	    	
	    }
		var count int
		/*var strtype string
		if getid !=""{
			strtype = "1"
		}else{
			strtype = "2"
		}	*/
		//`exec lbr_memass_search 1, 395 , '' ,'',''`	
		var sqlstr string
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		 if (isbaseonassignid ==`1`){
		 	//sqlstr = `exec lbr_memass_search 1, `+assignmentID+` , '' ,'',''`
		 	sqlstr ="exec SIS_SupOrd_Misc 3, "+str_OrgID+",0, "+hdr_id+",'"+getq+"'"
		 }else{
		 	//sqlstr = "exec LBR_Member_Get "+strtype+",  '"+getq+"'"
		 	sqlstr ="exec SIS_SupOrd_Misc 3, "+str_OrgID+",0, "+hdr_id+",'"+getq+"'"


		 }
		 fmt.Println(sqlstr)
		 rows ,err, _,_  := config.Ap_sql(sqlstr,1)    
	    if err != nil {
			fmt.Println("db error:",err)
		}
		InItem := []Inner_item{}
	    var InItemRow Inner_item
	    for rows.Next() {
	        var r Inner_item
	        err = rows.Scan(&r.Id,&r.Fullname,&r.Item_image)
	        if err != nil {
	            panic(err.Error())
	        }
			count = count +1
			post := Inner_item {
	        	r.Id,
	        	r.Fullname, 
	        	config.Img_to_render(`img_uploads/size_80/`+str_OrgID+`/`+r.Item_image+`.jpg`,r.Item_image,80,str_OrgID), 
		    	
		    }
		    InItem = append(InItem, post)
	    InItemRow = Inner_item{r.Id,r.Item_image, config.Img_to_render(`img_uploads/size_80/`+str_OrgID+`/`+r.Item_image+`.jpg`,r.Item_image,80,str_OrgID)}
	    }
	    

	  	if getid !=""{
		
			//jsondata := JsondataEmpty{count,false,InItem}
		    js, err := json.Marshal(InItemRow)
		    
			if err != nil {
		    	http.Error(w, err.Error(), http.StatusInternalServerError)
		    	return
				}
		  	w.Header().Set("Content-Type", "application/json")
		  	//w.Header().Set("Content-Type", "application/json")
		  	w.Write(js)
		
		}else{
		
			jsondata := JsondataEmpty{count,false,InItem}
		     js, err := json.Marshal(jsondata)
		    
			if err != nil {
		    	http.Error(w, err.Error(), http.StatusInternalServerError)
		    	return
				}
		  	w.Header().Set("Content-Type", "application/json")
		  	//w.Header().Set("Content-Type", "application/json")
		  	w.Write(js)

		}	

}


func DPostHandler(w http.ResponseWriter, r *http.Request) {
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		hdr_id :=r.URL.Query().Get("hdr_id")
		username, _ := login.Get_account_info(r)
		r.ParseForm()
		fmt.Println(`SIS_SupOrd_Post_Check `+str_OrgID+`,`+hdr_id )
		arr_raw_data := datatables.DataList(`SIS_SupOrd_Post_Check `+str_OrgID+`,`+hdr_id)
		fmt.Println(arr_raw_data)
		if(arr_raw_data==nil){
			fmt.Println(`GOOD TO POST`)
			arr_raw_data = datatables.DataList(`SIS_SupOrd_Post '`+username+`',`+str_OrgID+`,`+hdr_id)
			fmt.Println(arr_raw_data)
		}
		js, err := json.Marshal(arr_raw_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)
}

func DCancelHandler(w http.ResponseWriter, r *http.Request) {
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		hdr_id :=r.URL.Query().Get("hdr_id")
		username, _ := login.Get_account_info(r)
		r.ParseForm()
		fmt.Println(`SIS_SupOrd_Cancel_Check `+str_OrgID+`,`+hdr_id )
		//arr_data_tag := datatables.DataList(`SIS_SupOrdTags_List 1,`+str_OrgID+` ,`+ tconf["hdr_id"] )

		arr_raw_data := datatables.DataList(`SIS_SupOrd_Cancel_Check `+str_OrgID+`,`+hdr_id)
		fmt.Println(arr_raw_data)
		if(arr_raw_data==nil){
			fmt.Println(`GOOD TO POST`)
			arr_raw_data = datatables.DataList(`SIS_SupOrd_Cancel '`+username+`',`+str_OrgID+`,`+hdr_id)
			fmt.Println(arr_raw_data)
		}
		js, err := json.Marshal(arr_raw_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)
}

func EmailHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	username, _ := login.Get_account_info(r)
	//rights :="rights"
	hdr_id :=r.URL.Query().Get("hdr_id")
	//str_hdrID :=strconv.Itoa(int(hdr_id))
	 //strconv.Itoa(int(time.Nanoseconds()))

	assignmentID :=r.URL.Query().Get("location")


	fmt.Println ("location ",assignmentID )
	

	if r.Method =="GET" {	

		

		tconf := make(map[string]string)
		tconf["h_id"] =hdr_id
		tconf["rights"]=rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID
		tconf["location"] = assignmentID
		
		assignmentID :=  r.URL.Query().Get("location") 
		fmt.Println ("assignmentID ",assignmentID )
		
		//tconf["data"] = datatables.DataList(sql)
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )
		
		//_,session_user_id 	:= login.Get_account_info(r)

		
		type Data struct {
			Rights string
			Tconf map[string]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		arr_data := datatables.Data_row(`exec SIS_SupOrdHdr_Get 2, `+str_OrgID+`, `+hdr_id)
		
		fmt.Println(arr_data) 
		tconf["suprec_hdr_id"] = 			arr_data[0]
		tconf["suprec_hdr_status"] = 		arr_data[1]
		tconf["suprec_hdr_locationid"] = 	arr_data[2] 
		tconf["suprec_hdr_locationname"] = 	arr_data[3]
		tconf["suprec_hdr_trandate"] = 		arr_data[4]  
		tconf["suprec_hdr_supplerid"] = 	arr_data[5] 
		tconf["suprec_hdr_suppliername"] = 	arr_data[6] 
		tconf["suprec_detail_remarks"] = 	arr_data[7] 
		tconf["suprec_hdr_email"] = 		arr_data[8]
		tconf["suprec_hdr_org"] = 			arr_data[9] 
		tconf["suprec_hdr_locationemail"] = arr_data[10] 
	

		

		tmpl := template.New("email.html").Funcs(helper.FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/supplier_order/email.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute( w,&Data{rights,tconf} )
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
	r.ParseForm()
	hrd_id := r.URL.Query().Get("hdr_id")
	supplier := r.URL.Query().Get("supplier")
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)

		fmt.Println(`SIS_Supplier_Get 1 ,  `+str_OrgID+` ,`+supplier)
	arr_supplier := datatables.Data_row(`SIS_Supplier_Get 4 ,  `+str_OrgID+` ,`+supplier  )
	
	supplier_name :=arr_supplier[2];
	org_name :=arr_supplier[10];


	
	fmt.Println(`SIS_SupOrd_Custom 1,`+str_OrgID+` ,`+ hrd_id)
	SIS_SupOrdTags_List := datatables.DataList(`SIS_SupOrd_Custom 2,`+str_OrgID+` ,`+ hrd_id )
	SIS_SupOrdHdr_Report := datatables.DataList(`SIS_SupOrd_Custom 1,`+str_OrgID+` ,`+ hrd_id )
    ht_tbl:=`<table width="100%" cellpadding="15">
	            <tr>
	                <th align="left" bgcolor="#E5E7E9"><b>#</th>
	                <th align="left" bgcolor="#E5E7E9"><b>Item</th>
	                <th align="right" bgcolor="#E5E7E9"><b>Qty</th>
	                <th align="left" bgcolor="#E5E7E9"><b>Remarks</th>
	            </tr>        
			`;
    for _ , value := range SIS_SupOrdTags_List{
    	ht_tbl+=`<tr>`;
	    	ht_tbl+=`<td>`+value[0]+`</td>`;
	    	ht_tbl+=`<td>`+value[1]+`</td>`;
	    	ht_tbl+=`<td align="right">`+value[2]+`</td>`;
	    	ht_tbl+=`<td>`+value[3]+`</td>`;
    	ht_tbl+=`</tr>`;
    }
    for _ , value2 := range SIS_SupOrdHdr_Report{
    	ht_tbl+=`<tr>`;
	    	ht_tbl+=`<td bgcolor="#E5E7E9"><u><b>`+value2[5]+`</u></b></td>`;
	    	ht_tbl+=`<td bgcolor="#E5E7E9"></td>`;
	    	ht_tbl+=`<td align="right" bgcolor="#E5E7E9"><b><u>`+value2[4]+`</u></b></td>`;
	    	ht_tbl+=`<td bgcolor="#E5E7E9"></td>`;
    	ht_tbl+=`</tr>`;


    }
	ht_tbl+=`</table>`




	e := email.NewEmail()
	e.From =   config.Email_id
	e.To =  []string{r.Form["emailto"][0]}
	e.Bcc = []string{config.Email_bcc}
	fmt.Println(`cc email here : `, r.Form["ccto"][0]);
	arr_line := strings.Split( r.Form["ccto"][0] , ",")
	var	arr_CC []string
	for key , val := range arr_line {
		fmt.Println(` key : `, key );
		fmt.Println(` val : `, val );
		arr_CC = append(arr_CC,val);
	}
    e.Cc = 	[]string{}

	if (len(arr_line) > 0){
		e.Cc = 	arr_CC
	}

	e.Subject = org_name + ` Supplier Order (ID# ` +hrd_id+ `)`
	e.HTML = []byte(`Hi `+supplier_name+`,<br><br>Sending you list of orders, as follows;<br><br>`+ht_tbl+`<br>Additional Message: `+r.Form["message"][0]+`<br>Thank you<br>`+org_name+`<br><br>This is a system generated email. Please send your reply to `+ r.Form["ccto"][0] + `.<br> Powered by <a href="www.memger.com"> Memger Tech Inc.</a> <i>Inventory like a <b>boss</b>... <a href="www.memger.com">try us for <b>free</b>.</a></i>`)

	retErr := e.Send("smtp.zoho.com:587", smtp.PlainAuth("", e.From, config.Email_password, "smtp.zoho.com"))
	fmt.Println(`Email sending Error : ` , retErr);
    var json_ret string 
    
    if(retErr == nil){
    	json_ret = `Success`
    }else{
    	json_ret = `Error`
    }
    fmt.Println("done Sending email"); 

    	js, err := json.Marshal(json_ret)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)   

	}	
}


func Show_prev_price(w http.ResponseWriter, r *http.Request) {
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	supplier := r.URL.Query().Get("supplier")
	location := r.URL.Query().Get("location")
	item := r.URL.Query().Get("item")
	if (item==``){
		item = `0`
	}
	sql := `SIS_SupOrd_ShowLastRecords 2, `+str_OrgID+`, `+supplier+`, `+item+`, `+location+``
	fmt.Println(sql)
    arr_data := datatables.DataList(sql)
	tmpl := template.New("show_prev.html").Funcs(config.FuncMap)
	var err error
	if tmpl, err = tmpl.ParseFiles("inventory/supplier_order/show_prev.html"); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w,arr_data)
	if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func DVoidHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["h_id"] = r.URL.Query().Get("h_id")
	   	tconf["delete_url"] = "/inventory/supplier_order/details/DVoidHandler"
	   	tmpl := template.New("modal_void_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/supplier_order/modal_void_loghdr.html"); err != nil {
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
		arr_data := datatables.Data_row(`SIS_SupOrd_Void '`+username+`',` +r.Form["h_id"][0]  )
		js, err := json.Marshal(arr_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)	
	}	   	
}
