package application

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
	rights :=r.URL.Query().Get("rights")
	transtype :=r.URL.Query().Get("transtype")
	transid :=r.URL.Query().Get("ID")
	baseURL := r.URL.Query().Get("baseURL")

	cusapp :=r.URL.Query().Get("cusapp")
	//rights :="static here"
	fmt.Println("rights here")
	fmt.Println(rights)

	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)

		

	last_start :=r.URL.Query().Get("last_start")
	last_length :=r.URL.Query().Get("last_length")
	last_search :=r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key,rights)		
	fmt.Println("dec rights here")
	fmt.Println(dec_rights)
	tconf := make(map[string]string)
	tconf["baseURL"] =baseURL
	tconf["hdr_id"] = r.URL.Query().Get("ID") 
	ln :=  r.URL.Query().Get("length")
	tconf["Panel_name"] ="Account Application Details" 
	tconf["Add_new_button"] ="false" 
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="true"
	tconf["Add_new_button_details"] ="true" 
	tconf["upload_button"] ="false"  

	tconf["transtype"] =transtype
	tconf["transid"] =transid
	//tconf["Add_new_button_url"] ="/inventory/customer_sales/details/DaddHandler?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&customer="+tconf["suprec_hdr_customer"]
		// content access here
	tconf["Action_baseURL"] ="/inventory/helper/application?rights="+rights
	tconf["ActionCol_param"] ="ID"
    //tconf["hdr_id"] = ""
    tconf["Archive"] ="false" // is archive
	tconf["Interpret_button"] ="false"
	tconf["Add_post_button"] ="false"
	tconf["Add_cancelpay_button"] ="false"
	
	tconf["Add_cancel_button"] ="false"
	//tconf["Add_cancel_button_url"] ="/inventory/helper/application/cancel?rights="+rights+"&hdr_id="+cusapp
	tconf["Add_upload_button"] ="false"
	tconf["Add_upload_button_url"] ="/inventory/helper/application/HaddHandler?rights="+rights+"&fromtype=cp&cp_id="+tconf["hdr_id"] 
	tconf["Add_duplicate_button"] ="false"
	tconf["Add_duplicate_button_url"] ="/inventory/helper/application/duplicate_sales?rights="+rights+"&fromtype=cp&cp_id="+tconf["hdr_id"] 
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
	var header = []string{"Action", "Transaction Type", "Transaction ID", "Transaction Date", "Transaction Amount",  "Other Adj","Applied Amount", "Remarks", "Transaction Balance"}  //column header
	if ln=="" {
		
		tconf["ismoreinfo"] = "true" // to description info
		fmt.Println(`SIS_CusAcc_Get  1 ,  '`+transtype+`' ,`+transid )
		arr_data := datatables.Data_row(`SIS_CusAcc_Get  1 ,  '`+transtype+`' ,`+transid)
		fmt.Println(arr_data) 
		
		tconf["CusAccID"] = 			arr_data[0]
		tconf["CusAccCustomerName"] = 	arr_data[1] 
		tconf["CusAccTransType"] = 		arr_data[2]
		tconf["CustAccTransID"] = 		arr_data[3] 
		tconf["CustAccTransDate"] = 	arr_data[4] 
		tconf["CustAccTransAmount"] = 	arr_data[5]  
		tconf["CustAccTransBal"] = 		arr_data[6]  
		tconf["CustAccTransStatus"] = 	arr_data[7]  
		tconf["CustAccTransOrg"] = 		arr_data[8]  
		tconf["Balance"] 		= 		arr_data[9]  
		
		transtype :=r.URL.Query().Get("transtype")
		transid :=r.URL.Query().Get("ID")

		
		tconf["Add_new_button_url"] ="/inventory/helper/application/DaddHandler?rights="+rights+"&ID="+transid+"&transtype="+transtype+"&cusapp="+tconf["CusAccID"]+"&balance="+tconf["Balance"]
		tconf["Add_post_button_url"] ="/inventory/helper/application/post?rights="+rights+"&cusapp="+tconf["CusAccID"]+"&hdr_id="+transid	
		tconf["Add_cancel_button_url"] ="/inventory/helper/application/cancel?rights="+rights+"&cusapp="+tconf["CusAccID"]
		tconf["Add_cancelpay_button_url"] ="/inventory/helper/application/cancelpay?rights="+rights+"&cusapp="+tconf["CusAccID"]+"&ID="+transid
	

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
				tconf["Add_cancelpay_button"] ="false"
				//tconf["Add_upload_button"] ="true"	
			}else{
				tconf["Add_post_button"] ="true" 
				tconf["Add_cancelpay_button"] ="false"
			}
		}
		if(arr_data[7]==`App Cancelled`){
			tconf["Add_post_button"] ="false" 	
			tconf["Add_new_button_details"] ="false"
			tconf["upload_button"] ="false"  
	 		tconf["Add_cancel_button"] ="false"
	 		tconf["Add_duplicate_button"] ="false"
	 		tconf["Add_cancelpay_button"] ="false"
		
		}else if(arr_data[7]==`App Posted`){
			if(transtype==`CusPay`){			
				tconf["Add_new_button_details"] ="false"
				tconf["Add_post_button"] ="false" 	
	 			tconf["Add_cancel_button"] ="true"
	 			tconf["Add_cancelpay_button"] ="false"
 				//tconf["Add_upload_button"] ="true"
 				tconf["Add_duplicate_button"] ="false"
 			}else {
 				tconf["Add_new_button_details"] ="false"
				tconf["Add_post_button"] ="false" 	
	 			tconf["Add_cancel_button"] ="true"
 				//tconf["Add_upload_button"] ="true"
 				tconf["Add_duplicate_button"] ="false"				
 			}

		
		}else if(arr_data[7]==`Empty`){
				tconf["Add_post_button"] ="false" 	
	 			tconf["Add_cancel_button"] ="false"
	 			tconf["Add_duplicate_button"] ="false"
	 			tconf["Add_cancelpay_button"] ="true"
		}

		if( strings.Contains(dec_rights, "DInterpret") ){
	 		tconf["Interpret_button"] ="true"
	
	 	}

	
		arr_data_tag := datatables.DataList(`SIS_CusSalTags_List 1,`+str_OrgID+` ,`+ tconf["hdr_id"] )
		arr_data_refID := datatables.DataList(`SIS_TranRef_List 0, `+str_OrgID+`,4,`+ tconf["hdr_id"] )
		
		
		var ajaxURLdata = "/inventory/helper/application?rights="+rights+"&ID="+transid+"&transtype="+ transtype+"&cusapp="+tconf["suprec_id"]+"&status="+ tconf["CustAccTransStatus"]
		fmt.Println("/inventory/helper/application?rights="+rights+"&ID="+transid+"&transtype="+ transtype+"&cusapp="+tconf["suprec_id"]+"&status="+ tconf["CustAccTransStatus"])
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
			Arr_tag [][]string
			Arr_refID [][]string
		}
		tmpl := template.New("table.html").Funcs(helper.FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("helper/application/table.html"); err != nil {
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
	 	tconf["ActionCol_edit_url"] =  "/inventory/helper/application/DEditHandler?rights="+rights+"&ID="+transid+"&cusapp="+cusapp+"&transtype="+transtype+"&h_id=" // config for edit click 
	 	tconf["ActionCol_delete"] ="true" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/inventory/helper/application/DDeleteHandler?h_id="
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
		if( strings.Contains(dec_rights, "HEdit") ){
			tconf["ActionCol_edit"] ="true"
		}
		if( strings.Contains(dec_rights, "HDelete") ){
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
		
		tconf["sql_total"] = `SIS_CusAccApp_List 1, '`+transtype+`',`+transid+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"]  = `SIS_CusAccApp_List 0, '`+transtype+`',`+transid+`,  `+str+`,  `+ln+`, '`+sr+`'`
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
	balance :=r.URL.Query().Get("balance")



	assignmentID :=r.URL.Query().Get("location")
	cusapp :=r.URL.Query().Get("cusapp")

	fmt.Println ("location ",assignmentID )
	fmt.Println ("cusapp ",cusapp )
	
	
	if r.Method =="GET" {	
		fmt.Println(`this is hdr id=`,hdr_id)
		customer :=r.URL.Query().Get("customer")
		transtype :=r.URL.Query().Get("transtype")
		transid :=r.URL.Query().Get("ID")
		
		tconf := make(map[string]string)
		tconf["username"] =username
		tconf["Org_id"] = str_OrgID
		tconf["hdr_id"] = hdr_id
		tconf["rights"] = rights
		tconf["test_js"] = `alert("from webserver")`
		tconf["isbarcode"] = isbarcode
		tconf["customer"] = customer
		tconf["location"] = assignmentID
		tconf["cusapp"] = cusapp
		tconf["transid"] = transid
		tconf["transtype"] = transtype
		tconf["balance"] = balance
		
		fmt.Println(`SIS_CusAcc_Get 2,'`+transtype+`' ,`+ transid)
		ArrCustAA := datatables.DataList(`SIS_CusAcc_Get 2,'`+transtype+`' ,`+ transid )

		type Data struct {
			Rights	string	
			Tconf 	map[string]string
			ArrCustAA [][]string
		}   
		
		var err error
		var html_file string = ``
		var html_filename string = ``

		if (isbarcode == `1`){
			html_filename = "Dadd_barcode.html"
			html_file = "helper/application/Dadd_barcode.html"
		}else{
			html_filename = "Dadd.html"
			html_file = "helper/application/Dadd.html"
		}

		tmpl := template.New(html_filename)
		if tmpl, err = tmpl.ParseFiles(html_file); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights , tconf, ArrCustAA})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {

		r.ParseForm()
		fmt.Println(`SIS_CusAccApp_Save 'Add', '`+username+`',0,`+r.Form["cusapp"][0]+` ,`+r.Form["applierTo"][0]+` ,`+r.Form["appliedAmt"][0]+` ,`+r.Form["otheradj"][0]+`,'`+r.Form["remarks"][0]+`'`)
		arr_raw_data := datatables.Data_row(`SIS_CusAccApp_Save 'Add', '`+username+`',0,`+r.Form["cusapp"][0]+` ,`+r.Form["applierTo"][0]+` ,`+r.Form["appliedAmt"][0] +` ,`+r.Form["otheradj"][0]+`,'`+r.Form["remarks"][0]+`'`)
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


func DEditHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	username, _ := login.Get_account_info(r)


	assignmentID :=r.URL.Query().Get("location")
	customer :=r.URL.Query().Get("customer")
	balance :=r.URL.Query().Get("balance")


	fmt.Println ("location ",assignmentID )
	fmt.Println ("customer ",customer )
	

	if r.Method =="GET" {	

		transtype :=r.URL.Query().Get("transtype")
		transid :=r.URL.Query().Get("ID")

		tconf := make(map[string]string)
		tconf["h_id"] =r.URL.Query().Get("h_id")
		tconf["rights"]=rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID
		tconf["location"] = assignmentID
		tconf["customer"] = customer
		tconf["balance"] = balance
		
		assignmentID :=  r.URL.Query().Get("location") 
		customer :=r.URL.Query().Get("customer")
		fmt.Println ("assignmentID ",assignmentID )
		fmt.Println ("customer ",customer )
		fmt.Println(`SIS_CusAcc_Get 2,'`+transtype+`' ,`+ transid)
		ArrCustAA := datatables.DataList(`SIS_CusAcc_Get 2,'`+transtype+`' ,`+ transid )
		
		type Data struct {
			Rights string
			Tconf map[string]string
			ArrCustAA [][]string
		}
		sql :=`SIS_CusAcc_Get 3, '',0,`+ tconf["h_id"]
		fmt.Println(sql)
		arr_data := datatables.Data_row(sql)
		
		fmt.Println(arr_data) 
		tconf["AppID"] = 					arr_data[0]
		tconf["AppCusAcc"] = 				arr_data[1]
		tconf["AppCusAcc_Description"] = 	arr_data[2]
		tconf["AppAppliedTo"] = 			arr_data[3]
		tconf["AppAppliedTo_Description"] = arr_data[4]
		tconf["AppAppliedAmt"] = 			arr_data[5]
		tconf["AppOtheradj"] = 				arr_data[6]
		tconf["AppRemarks"] = 				arr_data[7]
		

		tmpl := template.New("Dedit.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("helper/application/Dedit.html"); err != nil {
			fmt.Println(err)
		}		

		err1 := tmpl.Execute( w,&Data{rights,tconf ,ArrCustAA} )
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		r.ParseForm()
		fmt.Println(`SIS_CusAccApp_Save 'Edit', '`+username+`',`+r.Form["id"][0]+`,`+r.Form["cusapp"][0]+` ,`+r.Form["applierTo"][0]+` ,`+r.Form["appliedAmt"][0]+` ,`+r.Form["otheradj"][0]+`,'`+r.Form["remarks"][0]+`'`)
		arr_raw_data := datatables.Data_row(`SIS_CusAccApp_Save 'Edit', '`+username+`',`+r.Form["id"][0]+`,`+r.Form["cusapp"][0]+` ,`+r.Form["applierTo"][0]+` ,`+r.Form["appliedAmt"][0]+` ,`+r.Form["otheradj"][0]+`,'`+r.Form["remarks"][0]+`'`)
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
	   	tconf["delete_url"] = "/inventory/helper/application/DDeleteHandler"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("helper/application/modal_delete_loghdr.html"); err != nil {
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

func DPostHandler(w http.ResponseWriter, r *http.Request) {
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		//hdr_id :=r.URL.Query().Get("hdr_id")
		cusapp :=r.URL.Query().Get("cusapp")
		fmt.Println ("cusapp ",cusapp )
		username, _ := login.Get_account_info(r)
		r.ParseForm()
		arr_raw_data := datatables.DataList(`SIS_CusAcc_appPost  '`+username+`',`+str_OrgID+`,`+cusapp)
		fmt.Println(arr_raw_data)
		js, err := json.Marshal(arr_raw_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)
}

func DCancelHandler(w http.ResponseWriter, r *http.Request) {
		//Org_id :=login.Get_session_org_id(r)
		//str_OrgID :=strconv.Itoa(Org_id)
		cusapp :=r.URL.Query().Get("cusapp")
		//username, _ := login.Get_account_info(r)
		r.ParseForm()
		arr_raw_data := datatables.DataList(`SIS_CusAcc_appCancel `+cusapp)
			fmt.Println(arr_raw_data)
		js, err := json.Marshal(arr_raw_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)
}


func DCancelPayHandler(w http.ResponseWriter, r *http.Request) {
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		transid :=r.URL.Query().Get("ID")
		username, _ := login.Get_account_info(r)
		r.ParseForm()
		arr_raw_data := datatables.DataList(`SIS_CusPay_Cancel '`+username+`',`+str_OrgID+`,`+transid)
			fmt.Println(arr_raw_data)
		js, err := json.Marshal(arr_raw_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)
}

