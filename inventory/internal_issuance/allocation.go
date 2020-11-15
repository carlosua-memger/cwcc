package internal_issuance

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

func AListHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	transtype :=r.URL.Query().Get("transtype")
	status :=r.URL.Query().Get("status")
	//rights :="static here"
	fmt.Println("rights here")
	fmt.Println(rights)

	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	hdr_id := r.URL.Query().Get("hdr_id")
	d_id := r.URL.Query().Get("d_id")

	last_start :=r.URL.Query().Get("last_start")
	last_length :=r.URL.Query().Get("last_length")
	last_search :=r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key,rights)		
	fmt.Println("dec rights here")
	fmt.Println(dec_rights)
	tconf := make(map[string]string)
	tconf["hdr_id"] = r.URL.Query().Get("ID") 
	ln :=  r.URL.Query().Get("length")
	tconf["Panel_name"] ="Internal Member Allocation" 
	tconf["Add_new_button"] ="false" 
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="true"
	tconf["Add_new_button_details"] ="true" 
	tconf["Add_new_button_allocation"] ="false" 
	tconf["upload_button"] ="false"   
	tconf["Action_baseURL"] ="/inventory/internal_issuance/details?rights="+rights
	tconf["transtype"] =transtype
	tconf["ActionCol_param"] ="ID"
	tconf["org"] =str_OrgID
    //tconf["hdr_id"] = ""
    tconf["Archive"] ="false" // is archive
	tconf["Interpret_button"] ="false"
	tconf["Add_post_button"] ="false"
	tconf["Add_post_button_url"] ="/inventory/internal_issuance/details/allocation/post?rights="+rights+"&hdr_id="+tconf["hdr_id"]
	
	tconf["Add_cancel_button"] ="false"
	tconf["Add_cancel_button_url"] ="/inventory/internal_issuance/details/allocation/cancel?rights="+rights+"&hdr_id="+tconf["hdr_id"]
	tconf["Add_upload_button"] ="false"
	tconf["Add_upload_button_url"] ="/inventory/customer_sales/HaddHandler?rights="+rights+"&fromtype=co&co_id="+tconf["hdr_id"] 
	tconf["withtotal"] = "false"  // table total at bottom

	tconf["level"] ="3" 
	tconf["pageLength"] =last_length 
	tconf["displayStart"] =last_start
	tconf["last_search"] = last_search
	if last_length==""{
		tconf["pageLength"] ="10" 
		tconf["displayStart"] ="0"
	}
	var header = []string{"Action", "ID", "Member", "Qty", "Remarks"}  //column header
	if ln=="" {
		
		tconf["ismoreinfo"] = "false" // to description info

		fmt.Println(`exec SIS_IntIssDtl_Get 1, `+str_OrgID+`, '`+tconf["hdr_id"]+`'`  )
		arr_data := datatables.Data_row(`exec SIS_IntIssDtl_Get 1, `+str_OrgID+`, '`+tconf["hdr_id"]+`'`  )
		
		fmt.Println(arr_data) 
		tconf["suprec_detail_id"] = 		arr_data[0]
		tconf["suprec_detail_item"] = 		arr_data[1] 
		tconf["suprec_detail_itemname"] = 	arr_data[2] 
		tconf["suprec_detail_qty"] = 		arr_data[3]
		tconf["suprec_detail_uprice"] = 	`00.00`
		tconf["suprec_detail_remarks"] = 	arr_data[4] 
		tconf["suprec_detail_hdr_id"] = 	arr_data[5]  

		transtype :=r.URL.Query().Get("transtype")
		fmt.Println(transtype)

		tconf["Add_new_button_url"] ="/inventory/internal_issuance/details/allocation/DaddHandler?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&d_id="+tconf["suprec_detail_hdr_id"]
		// content access here		
		// rights for tables 	
		if( strings.Contains(dec_rights, "DAdd") ){
			tconf["Add_new_button_details"] ="true"
		}
		/*if( strings.Contains(dec_rights, "DUpload") ){
			tconf["upload_button"] ="true"  
	 	
		}*/
		if( strings.Contains(dec_rights, "DPost") ){
			if(arr_data[1]==`Empty` && arr_data[1]==`Posted`){
				tconf["Add_post_button"] ="false" 
				tconf["upload_button"] ="true"	
			}else{
				tconf["Add_post_button"] ="false" 
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
	 			tconf["Add_cancel_button"] ="true"
	 			tconf["Add_upload_button"] ="true"

		
		}else if(arr_data[1]==`Empty`){
				tconf["Add_post_button"] ="false" 	
	 			tconf["Add_cancel_button"] ="false"

		}

		if( strings.Contains(dec_rights, "DInterpret") ){
	 		tconf["Interpret_button"] ="true"
	
	 	}

		//end rights for tables
	 	
		//end of content accesss
		
		
		
		var ajaxURLdata = "/inventory/internal_issuance/details/allocation?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&d_id="+tconf["suprec_detail_hdr_id"]+"&status="+status
		fmt.Println("/inventory/internal_issuance/details/allocation?rights="+rights+"&hdr_id="+tconf["hdr_id"]+"&d_id="+tconf["suprec_detail_hdr_id"]+"&status="+status)
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html").Funcs(helper.FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/internal_issuance/table.html"); err != nil {
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
	 	tconf["ActionCol_edit_is_modal"] ="true" // config for edit click  modal design
	 	tconf["ActionCol_edit_url"] =  "/inventory/internal_issuance/details/allocation/DEditHandler?rights="+rights+"&d_id="+d_id+"&hdr_id="+hdr_id+"&h_id=" // config for edit click 
	 	tconf["ActionCol_edit_is_modal"] ="true" // config for edit click  modal design
	 	tconf["ActionCol_delete"] ="true" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/inventory/internal_issuance/details/allocation/DDeleteHandler?h_id="
	 	tconf["ActionCol_detail"] ="false" // config for details click 
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
		/*tconf["col_img"] ="false"
	 	tconf["name_with_img"] ="true"
	 	tconf["col_img_num"] ="1"
	 	tconf["img_id_col"] ="4" 
	 	tconf["col_hidden1"] ="0" */ // what column to hide
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
		
		fmt.Println(status)
		//status := r.URL.Query().Get("status")
	 	if(status==`Cancelled`){
			tconf["ActionCol_edit"] ="false"
			tconf["ActionCol_delete"] ="false" 
			tconf["Add_new_button_details"] ="false" 
	 		tconf["Add_duplicate_button"] ="false"
		}else if(status==`Posted`){
			tconf["ActionCol_edit"] ="false"
			tconf["ActionCol_delete"] ="false" 
			tconf["Add_new_button_details"] ="false" 
	 		tconf["Add_duplicate_button"] ="false"
		}
		//end rights for tables
		//_,session_user_id := login.Get_account_info(r)

		//SP SIS_IntAcc_List 1|0, ‘IntIss’, @headerID, @detailID, @pageno, @rowperpage, @searchkey

		tconf["sql_total"] = `SIS_IntAcc_List 1, 'IntIss', `+d_id+`,`+hdr_id+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"]  = `SIS_IntAcc_List 0, 'IntIss', `+d_id+`,`+hdr_id+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["org"] = str_OrgID
		fmt.Println(tconf["sql_data"])
		datatables.DatatableretArray(w,tconf)
		  	 
		
		
		}	    
}

func AAddHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	username, _ := login.Get_account_info(r)
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	rights :=r.URL.Query().Get("rights")
	hdr_id :=r.URL.Query().Get("hdr_id")
	d_id :=r.URL.Query().Get("d_id")
	isbarcode :=r.URL.Query().Get("isbarcode")
	assignmentID :=r.URL.Query().Get("location")

	fmt.Println ("location ",assignmentID )
	//fmt.Fprint(w,rights)
	
	if r.Method =="GET" {	
		fmt.Println(`this is hdr id=`,hdr_id)
		fmt.Println(`this is d id=`,d_id)
		
		tconf := make(map[string]string)
		//tconf["parentID"] =r.URL.Query().Get("parentID")
		tconf["username"] =username
		tconf["Org_id"] = str_OrgID
		tconf["hdr_id"] = hdr_id
		tconf["d_id"] = d_id
		tconf["rights"] = rights
		tconf["test_js"] = `alert("from webserver")`
		tconf["isbarcode"] = isbarcode
		tconf["location"] = assignmentID

		type Data struct {
			Rights	string	
			Tconf 	map[string]string
			SupRepHdr [][]string		
		}   
	
		arrdata:= datatables.DataList(`exec SIS_IntAct_Get 3,`+str_OrgID+`, 'IntIss', '`+tconf["d_id"]+`' , '`+tconf["hdr_id"]+`' `   )
		fmt.Println(arrdata);

		var err error
		var html_file string = ``
		var html_filename string = ``

		if (isbarcode == `1`){
			html_filename = "Dadd_barcode.html"
			html_file = "inventory/internal_issuance/Dadd_barcode.html"
		}else{
			html_filename = "DAllocation.html"
			html_file = "inventory/internal_issuance/DAllocation.html"
		}

		tmpl := template.New(html_filename)
		if tmpl, err = tmpl.ParseFiles(html_file); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights , tconf, arrdata})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {

		r.ParseForm()	
		hrd_id := r.Form["hrd_id"][0]
		h_id := r.Form["h_id"][0]
		checkbox_id := r.Form["checkbox_id"][0]
		qty := r.Form["qty"][0]
		remarks := r.Form["remarks"][0]
		
		ret_add := datatables.Data_row(`SIS_IntAcc_Save  'Add',  '`+username+`', 0, `+str_OrgID+`, 'IntIss', `+h_id+`, `+hrd_id+`, `+checkbox_id+`,`+qty+`, '`+remarks+`'` )

		fmt.Println(ret_add)
	    js, err := json.Marshal(ret_add)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)	

	}		
}

func AEditHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	username, _ := login.Get_account_info(r)
	//rights :="rights"

	d_id :=r.URL.Query().Get("d_id")
	h_id :=r.URL.Query().Get("h_id")
	hdr_id :=r.URL.Query().Get("hdr_id")
	assignmentID :=r.URL.Query().Get("location")

	fmt.Println ("location ",assignmentID )
	

	if r.Method =="GET" {	
		tconf := make(map[string]string)
		//tconf["d_id"] =r.URL.Query().Get("d_id")
		tconf["d_id"] = d_id
		tconf["h_id"] = h_id
		tconf["hdr_id"] = hdr_id
		tconf["rights"]=rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID
		tconf["location"] = assignmentID
		
		assignmentID :=  r.URL.Query().Get("location") 
		fmt.Println ("assignmentID ",assignmentID )
		
		type Data struct {
			Rights string
			Tconf map[string]string
			//SupRepHdr [][]string
		}

		//arrdata:= datatables.DataList(`exec SIS_IntAct_Get 1,`+str_OrgID+``  )
		//fmt.Println(arrdata);


		arr_data := datatables.Data_row(`SIS_IntAct_Get 1,'`+str_OrgID+`', '`+tconf["h_id"]+`' `  )

		fmt.Println(arr_data) 
		tconf["intact_detail_id"] = 					arr_data[0]
		tconf["intact_detail_org"] = 					arr_data[1] 
		tconf["intact_detail_trantype"] = 				arr_data[2] 
		tconf["intact_detail_header"] = 				arr_data[3]
		tconf["intact_detail_detail"] = 				arr_data[4] 
		tconf["intact_detail_member"] = 				arr_data[5] 
		tconf["intact_detail_member_name"] = 			arr_data[6] 
		tconf["intact_detail_qty"] = 					arr_data[7]  
		tconf["intact_detail_remarks"] = 				arr_data[8]  

		tmpl := template.New("Aedit.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/internal_issuance/Aedit.html"); err != nil {
			fmt.Println(err)
		}
		//err1 := tmpl.Execute( w,&Data{rights,tconf,arrdata} )
		err1 := tmpl.Execute( w,&Data{rights,tconf} )
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
	  	r.ParseForm()


		fmt.Println(`SIS_IntAcc_Save  'edit',  '`+username+`', `+r.Form["h_id"][0]+`, `+str_OrgID+`, 'IntIss', 0,0, `+r.Form["leader"][0]+`, `+r.Form["qty"][0]+`, '`+r.Form["remarks"][0]+`'` )
		arr_raw_data := datatables.Data_row(`SIS_IntAcc_Save  'edit',  '`+username+`', `+r.Form["h_id"][0]+`, `+str_OrgID+`, 'IntIss', 0,0, `+r.Form["leader"][0]+`, `+r.Form["qty"][0]+`, '`+r.Form["remarks"][0]+`'` )
		
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

func ADeleteHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["h_id"] = r.URL.Query().Get("h_id")
	   	tconf["delete_url"] = "/inventory/internal_issuance/details/allocation/DDeleteHandler"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/internal_issuance/modal_delete_loghdr.html"); err != nil {
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
		arr_data := datatables.Data_row(`SIS_IntAcc_Save 'Delete','`+username+`',` +r.Form["h_id"][0]  )
		js, err := json.Marshal(arr_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)	
	}	   	
}
