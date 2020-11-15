package supplier_receipt

import (
	
		"net/http"
		"fmt"
		"text/template"
		config "../../config"
		"strconv"
		"encoding/json"
		datatables "../../datatables"
		login "../../login"
		helper "../../helper"
		"strings"
)




func DListHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)	
	//fmt.Fprint(w,"mars test here")
	//fmt.Fprint(w, r.URL.Query().Get("rights"))
	rights :=r.URL.Query().Get("rights")
	location :=r.URL.Query().Get("location")
	customer :=r.URL.Query().Get("customer")
	cusapp :=r.URL.Query().Get("cusapp")
	//rights :="static here"
	//rights :="static here"
	fmt.Println("rights here")
	fmt.Println(rights)

	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	hdr_id := r.URL.Query().Get("hdr_id")

		
	fmt.Println(cusapp)
	last_start :=r.URL.Query().Get("last_start")
	last_length :=r.URL.Query().Get("last_length")
	last_search :=r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key,rights)		
	fmt.Println("dec rights here")
	fmt.Println(dec_rights)
	tconf := make(map[string]string)
	tconf["hdr_id"] = r.URL.Query().Get("ID") 
	ln :=  r.URL.Query().Get("length")
	tconf["Panel_name"] ="Account Application Details" 
	tconf["Add_new_button"] ="false" 
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="true"
	tconf["Add_new_button_details"] ="false" 
	tconf["upload_button"] ="false"   
	//tconf["Add_new_button_url"] ="/inventory/customer_adjustment/details/DaddHandler?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&customer="+tconf["suprec_hdr_customer"]
		// content access here
	tconf["Action_baseURL"] ="/inventory/customer_adjustment/details?rights="+rights
	tconf["ActionCol_param"] ="ID"
    //tconf["hdr_id"] = ""
    tconf["Archive"] ="false" // is archive
	tconf["Interpret_button"] ="false"
	tconf["Add_post_button"] ="false"
	//tconf["Add_post_button_url"] ="/inventory/customer_adjustment/details/post?rights="+rights+"&cusapp="+tconf["cusapp"]+"&hdr_id="+tconf["hdr_id"]
	
	tconf["Add_cancel_button"] ="false"
	tconf["Add_cancel_button_url"] ="/inventory/customer_adjustment/details/cancel?rights="+rights+"&hdr_id="+tconf["hdr_id"]
	tconf["Add_upload_button"] ="false"
	tconf["Add_upload_button_url"] ="/inventory/customer_adjustment/HaddHandler?rights="+rights+"&fromtype=cp&cp_id="+tconf["hdr_id"] 
	tconf["Add_duplicate_button"] ="false"
	tconf["Add_duplicate_button_url"] ="/inventory/customer_adjustment/details/duplicate_sales?rights="+rights+"&fromtype=cp&cp_id="+tconf["hdr_id"] 
	tconf["withtotal"] = "true"  // table total at bottom
	tconf["org"] =str_OrgID
	

	tconf["level"] ="2" 
	tconf["pageLength"] =last_length 
	tconf["displayStart"] =last_start
	tconf["last_search"] = last_search
	if last_length==""{
		tconf["pageLength"] ="10" 
		tconf["displayStart"] ="0"
	}
	var header = []string{"Action", "ID", "Transaction Type", "Transaction ID", "Transaction Date", "Transaction Amount", "Applied Amount", "Other Adj", "Remarks", "Transaction Balance"}  //column header
	if ln=="" {
		
		tconf["ismoreinfo"] = "true" // to description info

		fmt.Println(`SIS_CusAdjHdr_Get 1, `+str_OrgID+`, '`+tconf["hdr_id"]+`'` )
		arr_data_cusadj := datatables.Data_row(`SIS_CusAdjHdr_Get 1, `+str_OrgID+`, '`+tconf["hdr_id"]+`'` )
		fmt.Println(arr_data_cusadj) 
		
		tconf["cusadj_id"] = 					arr_data_cusadj[0]
		tconf["cusadj_status"] = 				arr_data_cusadj[1] 
		tconf["cusadj_trandate"] = 				arr_data_cusadj[2] 
		tconf["cusadj_customer"] = 				arr_data_cusadj[3] 
		tconf["cusadj_customer_name"] = 		arr_data_cusadj[4] 
		tconf["cusadj_cusadjtype"] = 			arr_data_cusadj[5]
		tconf["cusadj_adjtype"] = 				arr_data_cusadj[6] 
		tconf["cusadj_amt"] = 					arr_data_cusadj[7]
		tconf["cusadj_remarks"] = 				arr_data_cusadj[8]




		//fmt.Println(`SIS_CusAcc_Get 1, `+tconf["cusadj_cusadjtype"]+`, '`+tconf["hdr_id"]+`', `+str_OrgID+`, `+tconf["cusadj_customer"]+`` )
		//arr_data := datatables.Data_row(`SIS_CusAcc_Get 1,  `+tconf["cusadj_cusadjtype"]+`, '`+tconf["hdr_id"]+`', `+str_OrgID+`, `+tconf["cusadj_customer"]+`` )
		fmt.Println(`SIS_CusAcc_Get 1, 'CusAdj', '`+tconf["hdr_id"]+`'` )
		arr_data := datatables.Data_row(`SIS_CusAcc_Get 1,  'CusAdj', '`+tconf["hdr_id"]+`'` )
		fmt.Println(arr_data) 
		
		tconf["suprec_id"] = 					arr_data[0]
		tconf["suprec_hdr_customer_name"] = 	arr_data[1] 
		tconf["suprec_hdr_Adjtype"] = 			arr_data[2] 
		tconf["suprec_hdr_id"] = 				arr_data[3] 
		tconf["suprec_hdr_transdate"] = 		arr_data[4] 
		tconf["suprec_hdr_amount"] = 			arr_data[5]
		tconf["suprec_hdr_balance"] = 			arr_data[6] 
		tconf["suprec_hdr_status"] = 			arr_data[7]

		tconf["Add_new_button_url"] ="/inventory/customer_adjustment/details/DaddHandler?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&customer="+tconf["cusadj_customer"]+"&cusapp="+tconf["suprec_id"]
		tconf["Add_post_button_url"] ="/inventory/customer_adjustment/details/post?rights="+rights+"&cusapp="+tconf["suprec_id"]+"&hdr_id="+tconf["hdr_id"]
	
		// content access here
		// rights for tables 	
		if( strings.Contains(dec_rights, "DAdd") ){
			tconf["Add_new_button_details"] ="true"
		}
		if( strings.Contains(dec_rights, "DUpload") ){
			tconf["upload_button"] ="true"  
	 	
		}
		if( strings.Contains(dec_rights, "DPost") ){
			if(arr_data[7]==`App Posted`){
				tconf["Add_post_button"] ="false" 
				//tconf["Add_upload_button"] ="true"	
			}else{
				tconf["Add_post_button"] ="true" 
			}
		}

		if( strings.Contains(dec_rights, "DCancel") ){
			if(arr_data[1]==`Empty` || arr_data[1]==`Cancelled` || arr_data[1]==`For Posting`){
				tconf["Add_cancel_button"] ="false" 	
			}else{
				tconf["Add_cancel_button"] ="true" 
			}
		}
		if(arr_data[7]==`App Cancelled`){
			tconf["Add_post_button"] ="false" 	
			tconf["Add_new_button_details"] ="false"
			tconf["upload_button"] ="false"  
	 		tconf["Add_cancel_button"] ="false"
	 		tconf["Add_duplicate_button"] ="false"
		
		}else if(arr_data[7]==`App Posted`){
				tconf["Add_new_button_details"] ="false"
				tconf["Add_post_button"] ="false" 	
	 			//tconf["Add_cancel_button"] ="true"
 				//tconf["Add_upload_button"] ="true"
 				tconf["Add_duplicate_button"] ="false"


		
		}else if(arr_data[7]==`Empty`){
				tconf["Add_post_button"] ="false" 	
	 			//tconf["Add_cancel_button"] ="false"
	 			tconf["Add_duplicate_button"] ="false"
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
		
		arr_data_tag := datatables.DataList(`SIS_CusSalTags_List 1,`+str_OrgID+` ,`+ tconf["hdr_id"] )
		arr_data_refID := datatables.DataList(`SIS_TranRef_List 0, `+str_OrgID+`,4,`+ tconf["hdr_id"] )
		
		
		var ajaxURLdata = "/inventory/customer_adjustment/details?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&status="+ tconf["suprec_hdr_status"]+"&location="+tconf["suprec_hdr_location"]+"&customer="+tconf["cusadj_customer"]+"&cusapp="+tconf["suprec_id"]
		fmt.Println("/inventory/customer_adjustment/details?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&status="+ tconf["suprec_hdr_status"]+"&location="+tconf["suprec_hdr_location"]+"&customer="+tconf["cusadj_customer"]+"&cusapp="+tconf["suprec_id"])
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
			Arr_tag [][]string
			Arr_refID [][]string
		}
		tmpl := template.New("table.html").Funcs(helper.FuncMap) 
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("inventory/customer_adjustment/table.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf,arr_data_tag,arr_data_refID})
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
	 	tconf["ActionCol_edit_url"] =  "/inventory/customer_adjustment/details/DEditHandler?rights="+rights+"&hdr_id="+hdr_id+"&cusapp="+cusapp+"&customer="+customer+"&location="+location+"&h_id=" // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/inventory/customer_adjustment/details/DDeleteHandler?h_id="
	 	tconf["ActionCol_detail"] ="false" // config for details click 
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
		/* tconf["col_img"] ="false"
	 	tconf["name_with_img"] ="true"
	 	tconf["col_img_num"] ="1"
	 	tconf["img_id_col"] ="6"*/
	 	//tconf["col_hidden1"] ="0" // what column to hide
		tconf["dr"]=dr  //draw

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
	 	if(status==`App Cancelled`){
			tconf["ActionCol_edit"] ="false"
			tconf["ActionCol_delete"] ="false" 
		}else if(status==`App Posted`){
			tconf["ActionCol_edit"] ="false"
			tconf["ActionCol_delete"] ="false" 
		}
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
		tconf["sql_total"] = `SIS_CusAccApp_List  1, 'CusAdj', `+hdr_id+`,  `+str+`,  `+ln+`, '`+sr+`'`
		//tconf["sql_data"]  = `SIS_Item_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"]  = `SIS_CusAccApp_List  0, 'CusAdj', `+hdr_id+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["org"] = str_OrgID
		fmt.Println(tconf["sql_data"])
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
	customer :=r.URL.Query().Get("customer")
	cusapp :=r.URL.Query().Get("cusapp")

	fmt.Println ("location ",assignmentID )
	fmt.Println ("cusapp ",cusapp )
	fmt.Println ("customer ",customer )
	

	
	
	//fmt.Fprint(w,rights)
	
	if r.Method =="GET" {	
		fmt.Println(`this is hdr id=`,hdr_id)
		customer :=r.URL.Query().Get("customer")
		cusapp :=r.URL.Query().Get("cusapp")
		
		tconf := make(map[string]string)
		//tconf["parentID"] =r.URL.Query().Get("parentID")
		tconf["username"] =username
		tconf["Org_id"] = str_OrgID
		tconf["hdr_id"] = hdr_id
		tconf["rights"] = rights
		tconf["test_js"] = `alert("from webserver")`
		tconf["isbarcode"] = isbarcode
		tconf["customer"] = customer
		tconf["location"] = assignmentID
		tconf["cusapp"] = cusapp
		

		arr_data_customer := datatables.DataList(`SIS_Customer_Get 2,`+str_OrgID  )
		arr_data_cusacc := datatables.DataList(`SIS_CusAcc_Get 2,  'CusAdj', '`+tconf["hdr_id"]+`'` )

		arr_data_adjtype := datatables.DataList(`exec SIS_CusAdjtype_Get 1,`+str_OrgID  )

		fmt.Println(arr_data_cusacc)
		
		type Data struct {
			Rights	string	
			Tconf 	map[string]string
			Ls_customer [][]string
			Ls_cusacc [][]string
			Ls_adjtype  [][]string
			
		}   
		
		
		var err error
		var html_file string = ``
		var html_filename string = ``

		if (isbarcode == `1`){
			html_filename = "Dadd_barcode.html"
			html_file = "inventory/customer_adjustment/Dadd_barcode.html"
		}else{
			html_filename = "Dadd.html"
			html_file = "inventory/customer_adjustment/Dadd.html"
		}

		tmpl := template.New(html_filename)
		if tmpl, err = tmpl.ParseFiles(html_file); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights , tconf,arr_data_customer,arr_data_cusacc,arr_data_adjtype})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {

		r.ParseForm()
		//fmt.Println(`SIS_CusAdjDtl_Save 'Add', '`+username+`',0,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+`,`+r.Form["qty"][0]+`,`+r.Form["unit_cost"][0]+`,'`+r.Form["remarks"][0]+`'` )
		//arr_raw_data := datatables.Data_row(`SIS_CusAdjDtl_Save 'Add', '`+username+`',0,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+`,`+r.Form["qty"][0]+`,`+r.Form["unit_cost"][0]+`,'`+r.Form["remarks"][0]+`'` )
		fmt.Println(`SIS_CusAccApp_Save 'Add','`+username+`',0, `+r.Form["cusapp"][0]+`,`+r.Form["cussal"][0]+` ,`+r.Form["amount"][0]+` ,`+r.Form["otheradj"][0]+`,'`+r.Form["remarks"][0]+`'`)
		arr_raw_data := datatables.Data_row(`SIS_CusAccApp_Save 'Add','`+username+`',0, `+r.Form["cusapp"][0]+`,`+r.Form["cussal"][0]+` ,`+r.Form["amount"][0]+` ,`+r.Form["otheradj"][0]+`,'`+r.Form["remarks"][0]+`'`)
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



	hdr_id :=r.URL.Query().Get("hdr_id")
	assignmentID :=r.URL.Query().Get("location")
	customer :=r.URL.Query().Get("customer")

	fmt.Println ("location ",assignmentID )
	fmt.Println ("customer ",customer )
	

	if r.Method =="GET" {	

		

		tconf := make(map[string]string)
		tconf["h_id"] =r.URL.Query().Get("h_id")
		tconf["hdr_id"] = hdr_id
		tconf["rights"]=rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID
		tconf["location"] = assignmentID
		tconf["customer"] = customer
		
		assignmentID :=  r.URL.Query().Get("location") 
		customer :=r.URL.Query().Get("customer")
		fmt.Println ("assignmentID ",assignmentID )
		fmt.Println ("customer ",customer )

		fmt.Println ("ID", tconf["h_id"])
		fmt.Println ("hdr_ID", tconf["hdr_id"])
		
		
		//tconf["data"] = datatables.DataList(sql)
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )
		
		//_,session_user_id 	:= login.Get_account_info(r)

		
		type Data struct {
			Rights string
			Tconf map[string]string
			Ls_cusacc [][]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		fmt.Println(`SIS_CusAcc_Get 3, 'CusAdj', 0, '`+tconf["h_id"]+`'` )
		arr_data := datatables.Data_row(`SIS_CusAcc_Get 3,  'CusAdj', 0, '`+tconf["h_id"]+`'` )
		fmt.Println(arr_data) 


		arr_data_cusacc := datatables.DataList(`SIS_CusAcc_Get 2,  'CusAdj', '`+tconf["hdr_id"]+`'` )
		fmt.Println(arr_data_cusacc) 
		
		tconf["suprec_id"] = 					arr_data[0]
		tconf["suprec_hdr_cusacc"] = 			arr_data[1] 
		tconf["suprec_hdr_cusacc_desc"] = 		arr_data[2] 
		tconf["suprec_hdr_applied_to"] = 		arr_data[3] 
		tconf["suprec_hdr_applied_desc"] = 		arr_data[4] 
		tconf["suprec_hdr_amount"] = 			arr_data[5]
		tconf["AppOtheradj"] = 					arr_data[6]
		tconf["AppRemarks"] = 					arr_data[7]
		

		tmpl := template.New("Dedit.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/customer_adjustment/Dedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute( w,&Data{rights,tconf,arr_data_cusacc} )
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		r.ParseForm()
		//fmt.Println(`SIS_CusAdjDtl_Save 'Edit', '`+username+`',`+r.Form["id"][0]+`,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+`,`+r.Form["qty"][0]+`,`+r.Form["unit_cost"][0]+`,'`+r.Form["remarks"][0]+`'` )
		//arr_raw_data := datatables.Data_row(`SIS_CusAdjDtl_Save 'Edit', '`+username+`',`+r.Form["id"][0]+`,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+`,`+r.Form["qty"][0]+`,`+r.Form["unit_cost"][0]+`,'`+r.Form["remarks"][0]+`'`)


		fmt.Println(`SIS_CusAccApp_Save 'Edit','`+username+`',`+r.Form["id"][0]+`, `+r.Form["cusapp"][0]+`,`+r.Form["cussal"][0]+` ,`+r.Form["amount"][0]+` ,`+r.Form["otheradj"][0]+`,'`+r.Form["remarks"][0]+`'`)
		arr_raw_data := datatables.Data_row(`SIS_CusAccApp_Save 'Edit','`+username+`',`+r.Form["id"][0]+`, `+r.Form["cusapp"][0]+`,`+r.Form["cussal"][0]+` ,`+r.Form["amount"][0]+` ,`+r.Form["otheradj"][0]+`,'`+r.Form["remarks"][0]+`'`)


		//fmt.Println(`SIS_CusAdjDtl_Save 'Edit', '`+username+`',`+r.Form["id"][0]+`,`+r.Form["hdr_id"][0]+` ,`+r.Form["amount"][0]+`,'`+r.Form["cussal"][0]+`'' ,`+r.Form["remarks"][0]+` ` )
		//arr_raw_data := datatables.Data_row(`SIS_CusAdjDtl_Save 'Edit', '`+username+`',`+r.Form["id"][0]+`,`+r.Form["amount"][0]+`,'`+r.Form["cussal"][0]+`'' ,`+r.Form["remarks"][0]+` ` )
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
	   	tconf["delete_url"] = "/inventory/customer_adjustment/details/DDeleteHandler"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/customer_adjustment/modal_delete_loghdr.html"); err != nil {
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
		arr_data := datatables.Data_row(`SIS_CusAccApp_Save 'Delete','`+username+`',` +r.Form["h_id"][0]  )
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
	fmt.Println(`is base sa assignID`,isbaseonassignid)
	fmt.Println(` assignID`,assignmentID)

	getq = getq + getid
		type Inner_item struct {
			Id string `json:"id"`
			Fullname string `json:"full_name"`
			Item_image string `json:"avatar"`
			Price string `json:"price"`

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
		 	sqlstr ="exec SIS_Item_Get 5, "+str_OrgID+",'"+getq+"', "+assignmentID+""
		 }else{
		 	//sqlstr = "exec LBR_Member_Get "+strtype+",  '"+getq+"'"
		 	sqlstr = "exec SIS_Item_Get 5, "+str_OrgID+",'"+getq+"'"


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
	        err = rows.Scan(&r.Id,&r.Fullname,&r.Item_image,&r.Price)
	        if err != nil {
	            panic(err.Error())
	        }
			count = count +1
			post := Inner_item {
	        	r.Id,
	        	r.Fullname, 
	        	config.Img_to_render(`img_uploads/size_80/`+str_OrgID+`/`+r.Item_image+`.jpg`,r.Item_image,80,str_OrgID), 
		    	r.Price,
		    }
		    InItem = append(InItem, post)
	    InItemRow = Inner_item{r.Id,r.Item_image, config.Img_to_render(`img_uploads/size_80/`+str_OrgID+`/`+r.Item_image+`.jpg`,r.Item_image,80,str_OrgID),r.Price}
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

func Ajax_getItem_isbarcode(w http.ResponseWriter, r *http.Request) {
			Org_id :=login.Get_session_org_id(r)
			str_OrgID :=strconv.Itoa(Org_id)
			barcode := r.URL.Query().Get("barcode")
			mode := r.URL.Query().Get("mode")


			arr_data := datatables.Data_row(`exec SIS_Barcode_Search `+mode+`,`+str_OrgID+`,'`+barcode+`'`  )
			js, err := json.Marshal(arr_data)
		    
			if err != nil {
		    	http.Error(w, err.Error(), http.StatusInternalServerError)
		    	return
				}
		  	w.Header().Set("Content-Type", "application/json")
		  	//w.Header().Set("Content-Type", "application/json")
		  	w.Write(js)
}



func DPostHandler(w http.ResponseWriter, r *http.Request) {
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)


		//hdr_id :=r.URL.Query().Get("hdr_id")
			cusapp :=r.URL.Query().Get("cusapp")
			fmt.Println ("cusapp ",cusapp )
		


		username, _ := login.Get_account_info(r)
		r.ParseForm()
		fmt.Println(`SIS_CusAcc_appPost '`+username+`',`+str_OrgID+`,`+cusapp )
		arr_raw_data := datatables.DataList(`SIS_CusAcc_appPost '`+username+`',`+str_OrgID+`,`+cusapp)
		fmt.Println(arr_raw_data)
		/*if(arr_raw_data==nil){
			fmt.Println(`GOOD TO POST`)
			arr_raw_data = datatables.DataList(`SIS_CusAdj_Post '`+username+`',`+str_OrgID+`,`+cusapp)
			fmt.Println(arr_raw_data)
		} */
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
		fmt.Println(`sis_cusadj_cancel '`+username+`',`+str_OrgID+`,`+hdr_id)
		//arr_data_tag := datatables.DataList(`SIS_SupRecTags_List 1,`+str_OrgID+` ,`+ tconf["hdr_id"] )

		// sis_cusadj_cancel @formuser, @org, @id


		arr_raw_data := datatables.DataList(`sis_cusadj_cancel '`+username+`',`+str_OrgID+`,`+hdr_id)
		fmt.Println(arr_raw_data)
		/* if(arr_raw_data==nil){
			fmt.Println(`GOOD TO POST`)
			arr_raw_data = datatables.DataList(`SIS_CusAdj_Cancel '`+username+`',`+str_OrgID+`,`+hdr_id)
			fmt.Println(arr_raw_data)
		} */
		js, err := json.Marshal(arr_raw_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)
}

func DAdd_reference(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["h_id"] = r.URL.Query().Get("h_id")
	   	tconf["rights"] = r.URL.Query().Get("rights")

	   	tconf["delete_url"] = "/inventory/customer_adjustment/details/DDeleteHandler"
	   	tmpl := template.New("Dadd_reference.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/customer_adjustment/Dadd_reference.html"); err != nil {
			fmt.Println(err)
		}
		err = tmpl.Execute(w,tconf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} 
	}else{
		//session_username,_ := login.Get_account_info(r)	
		r.ParseForm()
		h_id:=r.Form["h_id"][0]
		ref_ID :=r.Form["item"][0]
		ary_ref_ID := strings.Split(ref_ID,",")
		

		username, _ := login.Get_account_info(r)
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		var arr_retdata [][]string
		var db_data []string
		for _,val := range ary_ref_ID {

			//Customer Sales Detail > Add Reference calls sis_tranref_list 2, 3 (hardcoded?), 6 (incorrect tranID), 4 
			//Must call SP sis_tranref_list 2, @org, 4, @transaction ID
				fmt.Println(`SIS_TranRef_Save 'Add' ,'`+username+`' ,0 ,`+str_OrgID+` ,4 ,`+h_id+` , 9 ,`+val )
				db_data = datatables.Data_row(`SIS_TranRef_Save 'Add' ,'`+username+`' ,0 ,`+str_OrgID+` ,4 ,`+h_id+` , 9 ,`+val  )
				db_data = append(db_data,val)
				//arr_retdata = append(arr_retdata,db_data)
				if(len(db_data)>0){
					arr_retdata = append(arr_retdata,db_data)
				}
				fmt.Println(db_data)		
		}	
		
		js, err := json.Marshal(arr_retdata)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)	
	}
}

func Ddelete_reference(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["trn_id"] = r.URL.Query().Get("trn_id")
	   	tconf["delete_url"] = "/inventory/customer_adjustment/details/DDeleteHandler"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/customer_adjustment/modal_delete_loghdr.html"); err != nil {
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
		arr_data := datatables.Data_row(`SIS_TranRef_Save 'Delete','`+username+`',` +r.Form["trn_id"][0]  )
		js, err := json.Marshal(arr_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)	
	}
}


func Dduplicate_sales(w http.ResponseWriter, r *http.Request) {
		username, _ := login.Get_account_info(r)
		tconf := make(map[string]string)
		tconf["id"] =r.URL.Query().Get("cp_id")
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		fmt.Println(`SIS_CusAdjHdr_Save 'Duplicate','`+username+`' ,  `+tconf["id"] +` , `+str_OrgID+`` )
		arr_data := datatables.Data_row(`SIS_CusAdjHdr_Save 'Duplicate','`+username+`' ,  `+tconf["id"] +` , `+str_OrgID+`` )
		js, err := json.Marshal(arr_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)		


}

func Ajax_get_reference(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	getq := r.URL.Query().Get("q")
	getid := r.URL.Query().Get("id")
	//h_id := r.URL.Query().Get("h_id")

	isbaseonassignid :=r.URL.Query().Get("isbaseonassignid")
	assignmentID := r.URL.Query().Get("assignmentID")
	customer := r.URL.Query().Get("customer")
	fmt.Println(`is base sa assignID`,isbaseonassignid)
	fmt.Println(` assignID`,assignmentID)
	fmt.Println(` customer`,customer)

	getq = getq + getid
		type Inner_item struct {
			Id string `json:"id"`
			Trndate interface{} `json:"trndate"`
			trantype string `json:"full_name"`
			tranID string `json:"tranID"`
			amount string `json:"amount"`
			Balance string `json:"Balance"`
			status string `json:"status"`

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
		 	//sqlstr ="exec SIS_Item_Get 3, "+str_OrgID+",'"+getq+"'"
		 	//fmt.Println(sqlstr)
		 	sqlstr = "SIS_CusAcc_Get 2, 'CusAdj', 0, "+str_OrgID+", "+customer
		 }else{
		 	//sqlstr = "exec LBR_Member_Get "+strtype+",  '"+getq+"'"
		 	
		 	//sqlstr = "exec SIS_Item_Get 2, "+str_OrgID+",'"+getq+"'"
		 	
		 	sqlstr = "SIS_CusAcc_Get 2,  'CusAdj', 0, "+str_OrgID+", "+customer


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
	        err = rows.Scan(&r.Id,&r.Trndate,&r.trantype,&r.tranID,&r.amount,&r.Balance,&r.status)
	        if err != nil {
	            panic(err.Error())
	        }
			count = count +1
			post := Inner_item {
	        	r.Id,
	        	config.Interfacetosting(r.Trndate), 
	        	r.trantype,
	        	r.tranID,
	        	r.amount,
	        	r.Balance,
	        	r.status,  
		    }
		    InItem = append(InItem, post)
	    	InItemRow = Inner_item{r.Id, config.Interfacetosting(r.Trndate) , r.trantype, r.tranID, r.amount, r.Balance,r.status}
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