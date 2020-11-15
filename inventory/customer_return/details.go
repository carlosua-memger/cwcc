package customer_return

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
	tconf["Panel_name"] ="Customer Return Detail" 
	tconf["Add_new_button"] ="false" 
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="true"
	tconf["Add_new_button_details"] ="false" 
	tconf["upload_button"] ="false"   
	//tconf["Add_new_button_url"] ="/inventory/customer_return/details/DaddHandler?rights="+rights+"&hdr_id="+tconf["hdr_id"]
	tconf["Action_baseURL"] ="/inventory/customer_return/details?rights="+rights
	tconf["ActionCol_param"] ="ID"
    //tconf["hdr_id"] = ""
    tconf["Archive"] ="false" // is archive
	tconf["Interpret_button"] ="false"
	tconf["Add_post_button"] ="false"
	tconf["Add_post_button_url"] ="/inventory/customer_return/details/post?rights="+rights+"&hdr_id="+tconf["hdr_id"]
	
	tconf["Add_cancel_button"] ="false"
	tconf["Add_cancel_button_url"] ="/inventory/customer_return/details/cancel?rights="+rights+"&hdr_id="+tconf["hdr_id"]
	tconf["Add_upload_button"] ="false"
	tconf["Add_upload_button_url"] ="/inventory/customer_return/HaddHandler?rights="+rights+"&fromtype=cp&cp_id="+tconf["hdr_id"] 
	tconf["Add_duplicate_button"] ="false"
	tconf["Add_duplicate_button_url"] ="/inventory/customer_return/details/duplicate_sales?rights="+rights+"&fromtype=cp&cp_id="+tconf["hdr_id"] 
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
	var header = []string{"Action", "Item Name", "Qty", "Unit Price", "Total Price","Supplier","Remarks"}  //column header
	if ln=="" {
		
		tconf["ismoreinfo"] = "true" // to description info
		fmt.Println(`exec SIS_CusRetHdr_Get 1, `+str_OrgID+`, '`+tconf["hdr_id"]+`'`   )
		arr_data := datatables.Data_row( `exec SIS_CusRetHdr_Get 1, `+str_OrgID+`, '`+tconf["hdr_id"]+`'`  )
		fmt.Println(arr_data) 
		
		tconf["cust_ret_id"] = 				arr_data[0]
		tconf["cust_ret_status"] = 			arr_data[1] 
		tconf["cust_ret_location"] = 		arr_data[2] 
		tconf["cust_ret_location_name"] = 	arr_data[3]
		tconf["cust_ret_trndate"] = 		arr_data[4] 
		tconf["cust_ret_customer_name"] = 	arr_data[5] 
		tconf["cust_ret_remarks"] = 		arr_data[6]
		tconf["cust_ret_customer"] = 		arr_data[7] 

		tconf["Add_new_button_url"] ="/inventory/customer_return/details/DaddHandler?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&customer="+tconf["cust_ret_customer"]+"&location="+tconf["cust_ret_location"]
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
				tconf["Add_upload_button"] ="true"	
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
	 		tconf["Add_duplicate_button"] ="false"
		
		}else if(arr_data[1]==`Posted`){
				tconf["Add_new_button_details"] ="false"
				tconf["Add_post_button"] ="false" 	
	 			//tconf["Add_cancel_button"] ="true"
 				tconf["Add_upload_button"] ="true"
 				tconf["Add_duplicate_button"] ="false"


		
		}else if(arr_data[1]==`Empty`){
				tconf["Add_post_button"] ="false" 	
	 			tconf["Add_cancel_button"] ="false"
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
		
		
		
		var ajaxURLdata = "/inventory/customer_return/details?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&status="+ tconf["cust_ret_status"]+"&location="+tconf["cust_ret_location"]
		fmt.Println("/inventory/customer_return/details?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&status="+ tconf["cust_ret_status"]+"&location="+tconf["cust_ret_location"])
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
		if tmpl, err = tmpl.ParseFiles("inventory/customer_return/table.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf,arr_data_tag,nil})
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
	 	tconf["ActionCol_edit_url"] =  "/inventory/customer_return/details/DEditHandler?rights="+rights+"&location="+location+"&h_id=" // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/inventory/customer_return/details/DDeleteHandler?h_id="
	 	tconf["ActionCol_detail"] ="false" // config for details click 
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
		tconf["col_img"] ="false"
	 	tconf["name_with_img"] ="true"
	 	tconf["col_img_num"] ="1"
	 	tconf["img_id_col"] ="6"
	 	tconf["col_hidden1"] ="0" // what column to hide
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
	 	if(status==`Cancelled`){
			tconf["ActionCol_edit"] ="false"
			tconf["ActionCol_delete"] ="false" 
		}else if(status==`Posted`){
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
		tconf["sql_total"] = `SIS_CusRetDtl_List 1, `+str_OrgID+`,`+hdr_id+`,  `+str+`,  `+ln+`, '`+sr+`'`
		//tconf["sql_data"]  = `SIS_Item_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"]  = `SIS_CusRetDtl_List 0, `+str_OrgID+`,`+hdr_id+`,  `+str+`,  `+ln+`, '`+sr+`'`
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
	rights :=r.URL.Query().Get("rights")
	hdr_id :=r.URL.Query().Get("hdr_id")
	isbarcode :=r.URL.Query().Get("isbarcode")



	assignmentID :=r.URL.Query().Get("location")

	fmt.Println ("location ",assignmentID )
	

	
	
	//fmt.Fprint(w,rights)
	
	if r.Method =="GET" {	
		fmt.Println(`this is hdr id=`,hdr_id)
		customer :=r.URL.Query().Get("customer")
		
		tconf := make(map[string]string)
		//tconf["parentID"] =r.URL.Query().Get("parentID")
		tconf["username"] =username
		tconf["Org_id"] = str_OrgID
		tconf["hdr_id"] = hdr_id
		tconf["rights"] = rights
		tconf["test_js"] = `alert("from webserver")`
		tconf["isbarcode"] = isbarcode
		tconf["location"] = assignmentID
		tconf["customer"] = customer
		
		assignmentID :=  r.URL.Query().Get("location") 
		fmt.Println ("location ",assignmentID )
			arr_data_supplier := datatables.DataList(`SIS_Supplier_Get 2,`+str_OrgID  )

		type Data struct {
			Rights	string	
			Tconf 	map[string]string
			Ls_supplier [][]string
		
		}   
		
		
		var err error
		var html_file string = ``
		var html_filename string = ``

		if (isbarcode == `1`){
			html_filename = "Dadd_barcode.html"
			html_file = "inventory/customer_return/Dadd_barcode.html"
		}else{
			html_filename = "Dadd.html"
			html_file = "inventory/customer_return/Dadd.html"
		}

		tmpl := template.New(html_filename)
		if tmpl, err = tmpl.ParseFiles(html_file); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights , tconf,arr_data_supplier})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {

		r.ParseForm()
		//fmt.Println(`SIS_CusSalDtl_Save 'Add', '`+username+`',0,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+`,`+r.Form["qty"][0]+`,`+r.Form["unit_cost"][0]+`,'`+r.Form["remarks"][0]+`'` )
		//arr_raw_data := datatables.Data_row(`SIS_CusSalDtl_Save 'Add', '`+username+`',0,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+`,`+r.Form["qty"][0]+`,`+r.Form["unit_cost"][0]+`,'`+r.Form["remarks"][0]+`'` )
		fmt.Println(`SIS_CusRetDtl_Save 'Add', '`+username+`',0,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+` ,`+r.Form["qty"][0]+` ,`+r.Form["unit_price"][0]+`,`+r.Form["supplier"][0]+`,'`+r.Form["remarks"][0]+`'`)
		arr_raw_data := datatables.Data_row(`SIS_CusRetDtl_Save 'Add', '`+username+`',0,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+` ,`+r.Form["qty"][0]+` ,`+r.Form["unit_price"][0]+`,`+r.Form["supplier"][0]+`,'`+r.Form["remarks"][0]+`'`)
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

	fmt.Println ("location ",assignmentID )
	

	if r.Method =="GET" {	

		

		tconf := make(map[string]string)
		tconf["h_id"] =r.URL.Query().Get("h_id")
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
			Ls_supplier [][]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		fmt.Println(`exec SIS_CusRetDtl_Get 1, `+str_OrgID+`, '`+tconf["h_id"]+`'`  )
		arr_data := datatables.Data_row(`exec SIS_CusRetDtl_Get 1, `+str_OrgID+`, '`+tconf["h_id"]+`'`  )
		
		tconf["custret_detail_id"] = 			arr_data[0]
		tconf["custret_detail_id_item"] = 		arr_data[1] 
		tconf["custret_detail_id_item_name"] = 	arr_data[2] 
		tconf["custret_detail_id_qty"] = 		arr_data[3]
		tconf["custret_detail_id_uprice"] = 	arr_data[4] 
		tconf["custret_detail_id_tprice"] = 	arr_data[5]
		tconf["custret_detail_id_supplier"] = 	arr_data[6] 
		tconf["custret_detail_id_supplier_name"] = 	arr_data[7] 
		tconf["custret_detail_id_remarks"] = 	arr_data[8]
		tconf["custret_detail_id_hdr_id"] = 	arr_data[9]
		
		arr_data_supplier := datatables.DataList(`SIS_SupplierItem_Get 3,`+tconf["custret_detail_id_item"]  )

		

		tmpl := template.New("Dedit.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/customer_return/Dedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute( w,&Data{rights,tconf,arr_data_supplier} )
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		r.ParseForm()
		//fmt.Println(`SIS_CusSalDtl_Save 'Edit', '`+username+`',`+r.Form["id"][0]+`,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+`,`+r.Form["qty"][0]+`,`+r.Form["unit_cost"][0]+`,'`+r.Form["remarks"][0]+`'` )
		//arr_raw_data := datatables.Data_row(`SIS_CusSalDtl_Save 'Edit', '`+username+`',`+r.Form["id"][0]+`,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+`,`+r.Form["qty"][0]+`,`+r.Form["unit_cost"][0]+`,'`+r.Form["remarks"][0]+`'`)
		fmt.Println(`SIS_CusRetDtl_Save 'Edit', '`+username+`',`+r.Form["id"][0]+`,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+` ,`+r.Form["qty"][0]+` ,`+r.Form["unit_price"][0]+`,`+r.Form["supplier"][0]+`,'`+r.Form["remarks"][0]+`'`)
		arr_raw_data := datatables.Data_row(`SIS_CusRetDtl_Save 'Edit', '`+username+`',`+r.Form["id"][0]+`,`+r.Form["hdr_id"][0]+` ,`+r.Form["item"][0]+` ,`+r.Form["qty"][0]+` ,`+r.Form["unit_price"][0]+`,`+r.Form["supplier"][0]+`,'`+r.Form["remarks"][0]+`'`)
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
	   	tconf["delete_url"] = "/inventory/customer_return/details/DDeleteHandler"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/customer_return/modal_delete_loghdr.html"); err != nil {
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
		arr_data := datatables.Data_row(`SIS_CusRetDtl_Save 'Delete','`+username+`',` +r.Form["h_id"][0]  )
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
		 	sqlstr ="exec SIS_CusRet_Misc 2, "+str_OrgID+",0, "+hdr_id+",'"+getq+"'"
		 }else{
		 	//sqlstr = "exec LBR_Member_Get "+strtype+",  '"+getq+"'"
		 	sqlstr ="exec SIS_CusRet_Misc 2, "+str_OrgID+",0, "+hdr_id+",'"+getq+"'"


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
		hdr_id :=r.URL.Query().Get("hdr_id")
		username, _ := login.Get_account_info(r)
		r.ParseForm()
		fmt.Println(`SIS_CusSal_Post_Check `+str_OrgID+`,`+hdr_id )
		arr_raw_data := datatables.DataList(`SIS_CusRet_Post_Check `+str_OrgID+`,`+hdr_id)
		fmt.Println(arr_raw_data)
		if(arr_raw_data==nil){
			fmt.Println(`GOOD TO POST`)
			arr_raw_data = datatables.DataList(`SIS_CusRet_Post '`+username+`',`+str_OrgID+`,`+hdr_id)
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
		fmt.Println(`SIS_CusSal_Cancel_Check `+str_OrgID+`,`+hdr_id )
		//arr_data_tag := datatables.DataList(`SIS_SupRecTags_List 1,`+str_OrgID+` ,`+ tconf["hdr_id"] )

		arr_raw_data := datatables.DataList(`SIS_CusRet_Cancel_Check `+str_OrgID+`,`+hdr_id)
		fmt.Println(arr_raw_data)
		if(arr_raw_data==nil){
			fmt.Println(`GOOD TO POST`)
			arr_raw_data = datatables.DataList(`SIS_CusRet_Cancel '`+username+`',`+str_OrgID+`,`+hdr_id)
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

func DAdd_reference(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["h_id"] = r.URL.Query().Get("h_id")
	   	tconf["rights"] = r.URL.Query().Get("rights")

	   	tconf["delete_url"] = "/inventory/customer_return/details/DDeleteHandler"
	   	tmpl := template.New("Dadd_reference.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/customer_return/Dadd_reference.html"); err != nil {
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
	   	tconf["delete_url"] = "/inventory/customer_return/details/DDeleteHandler"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/customer_return/modal_delete_loghdr.html"); err != nil {
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
	


}

func Ajax_get_reference(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	getq := r.URL.Query().Get("q")
	getid := r.URL.Query().Get("id")
	h_id := r.URL.Query().Get("h_id")

	isbaseonassignid :=r.URL.Query().Get("isbaseonassignid")
	assignmentID := r.URL.Query().Get("assignmentID")
	fmt.Println(`is base sa assignID`,isbaseonassignid)
	fmt.Println(` assignID`,assignmentID)

	getq = getq + getid
		type Inner_item struct {
			Id string `json:"id"`
			Trndate interface{} `json:"trndate"`
			Name string `json:"full_name"`

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
		 	sqlstr = "sis_tranref_list 2, "+str_OrgID+", 4,"+h_id
		 }else{
		 	//sqlstr = "exec LBR_Member_Get "+strtype+",  '"+getq+"'"
		 	
		 	//sqlstr = "exec SIS_Item_Get 2, "+str_OrgID+",'"+getq+"'"
		 	
		 	sqlstr = "sis_tranref_list 2, "+str_OrgID+", 4, "+h_id


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
	        err = rows.Scan(&r.Id,&r.Trndate,&r.Name)
	        if err != nil {
	            panic(err.Error())
	        }
			count = count +1
			post := Inner_item {
	        	r.Id,
	        	config.Interfacetosting(r.Trndate), 
	        	r.Name,  
		    }
		    InItem = append(InItem, post)
	    	InItemRow = Inner_item{r.Id, config.Interfacetosting(r.Trndate) , r.Name}
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
		

func Show_prev_price(w http.ResponseWriter, r *http.Request) {
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	customer := r.URL.Query().Get("customer")
	location := r.URL.Query().Get("location")
	item := r.URL.Query().Get("item")
	if (item==``){
		item = `0`
	}
	sql := `SIS_CusSal_ShowLastRecords 2, `+str_OrgID+`, `+customer+`, `+item+`, `+location+``
	fmt.Println(sql)
    arr_data := datatables.DataList(sql)
	tmpl := template.New("show_prev.html").Funcs(config.FuncMap)
	var err error
	if tmpl, err = tmpl.ParseFiles("inventory/customer_return/show_prev.html"); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w,arr_data)
	if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}