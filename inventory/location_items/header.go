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
		helper "../../helper"
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
	
	last_start :=r.URL.Query().Get("last_start")
	last_length :=r.URL.Query().Get("last_length")
	last_search :=r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key,rights)		
	//dec_rights :="static here"
	tconf := make(map[string]string)
	ln :=  r.URL.Query().Get("length")
	tconf["Panel_name"] ="Item Locations" 
	tconf["Add_new_button"] ="false" 
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="false"
	tconf["Add_new_button_details"] ="false" 
	tconf["upload_button"] ="false"   
	tconf["Add_new_button_url"] ="/inventory/location_items/HaddHandler?rights="+rights 
	tconf["Action_baseURL"] ="/inventory/location_items?rights="+rights
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
	var header = []string{"Action", "ID", "Name", "Location", "Unit Price","SKU","Remarks"}  //column header
	if ln=="" {
		// content access here
		if( strings.Contains(dec_rights, "HAdd") ){
			tconf["Add_new_button"] ="true"
			//fmt.Println("run herer")
		}
		if( strings.Contains(dec_rights, "Dadd ") ){
			
		}
			
		//end of content accesss
		var ajaxURLdata = "/inventory/location_items?rights="+rights
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("inventory/location_items/table.html"); err != nil {
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
	 	tconf["ActionCol_edit_url"] ="/inventory/location_items/HEditHandler?rights="+rights+"&h_id=" // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/inventory/location_items/HDeleteHandler?h_id="
	 	tconf["ActionCol_detail"] ="false" // config for details click 
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	tconf["ActionCol_duplicate"] ="false"
	 	tconf["ActionCol_duplicate_is_modal"] ="false"
	 	tconf["ActionCol_duplicate_url"] ="/inventory/location_items/HDuplicateHandler?rights="+rights+"&h_id="
	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
	 	tconf["col_img"] ="false"
	 	tconf["name_with_img"] ="true"
	 	tconf["col_img_num"] ="1"
	 	tconf["img_id_col"] ="0"
		tconf["dr"]=dr 
		 //draw
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

		if( strings.Contains(dec_rights, "HReplicate") ){
			tconf["ActionCol_duplicate"] ="true" // config for delete click 
	 	
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
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)



		tconf["sql_total"] = `SIS_Location_Item_List 1, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"]  = `SIS_Location_Item_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["org"] = str_OrgID
		fmt.Println(tconf["sql_data"])
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
			fmt.Println(`sis_itemtags_get 1, `+org_id+`, 0 ,`+tag_id)
			retdata := datatables.DataList(`sis_itemtags_get 1, `+org_id+`, 0 ,`+tag_id)
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
		tconf["test_js"] = `alert("from webserver")`
		fmt.Println(`SIS_ItemClass_Get 2,`+str_OrgID  )
		fmt.Println(`SIS_Supplier_Get 2,`+str_OrgID  )
		fmt.Println(`SIS_Location_Get 6,`+str_OrgID  )
	

		//arr_data := datatables.DataList(`select id,dbo.SIS_Itemclass_Name(1, 3, id) [itemclass_name] from sis_itemclass`  )
		arr_data := datatables.DataList(`SIS_ItemClass_Get 2,`+str_OrgID  )
		//arr_data_supplier := datatables.DataList(`select id,dbo.SIS_Supplier_Name(1, 3, id) [supplier_name] from SIS_Supplier`  )

		arr_data_supplier := datatables.DataList(`SIS_Supplier_Get 2,`+str_OrgID  )
		arr_data_itemtag := datatables.DataList(`sis_itemtags_list 1,`+str_OrgID+` , 0`  )
		arr_data_location := datatables.DataList(`exec SIS_Location_Get 6,`+str_OrgID  )
		type Data struct {
			Rights		string	
			Tconf 		map[string]string
			Ls_itemClass [][]string
			Ls_supplier [][]string
			Arr_itemtag [][]string
			Ls_location [][]string
		}   
		tmpl := template.New("Hadd.html").Funcs(local_FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/location_items/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights , tconf,arr_data,arr_data_supplier,arr_data_itemtag,arr_data_location})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		r.ParseForm()	
		img_dir := r.Form["Image_dr"][0]
		//item_id := r.Form["Item_id"][0]
		item_name := r.Form["Item_name"][0]

		//helper.Img_resizer(img_dir,item_id)
		helper.Img_resizer(img_dir,item_name,str_OrgID)

		js, err := json.Marshal(`resizer success...`)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)	
		

	}		
}

func HAddTagHandler(w http.ResponseWriter, r *http.Request) {
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
	        	arr_data_itemtag := datatables.Data_row(`SIS_itemtags_save 'Update', '`+username+`',`+str_OrgID+`,`+item_id+`, `+tag+`, '`+value_input+`', ''` )

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
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	//rights :="rights"

	if r.Method =="GET" {	

		username, _ := login.Get_account_info(r)

		tconf := make(map[string]string)
		tconf["h_id"] =r.URL.Query().Get("h_id")
		tconf["rights"]=rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID
		
		//tconf["data"] = datatables.DataList(sql)
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )
		
		//_,session_user_id 	:= login.Get_account_info(r)
		session_user_id := 1  //static here
		session_user_id_str := strconv.Itoa(session_user_id) 
		fmt.Println(session_user_id_str)

		type Data struct {
			Rights string
			Conf map[string]string
			Ls_itemClass [][]string
			Ls_supplier  [][]string
			Arr_itemtag  [][]string
			Ls_location [][]string
			
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		fmt.Println("running sql")
		fmt.Println(`exec SIS_Location_Item_Get 1, `+str_OrgID+`, '`+tconf["h_id"]+`'`   )
		arr_data := datatables.Data_row(`exec SIS_Location_Item_Get 1, `+str_OrgID+`, '`+tconf["h_id"]+`'`  )
		
		//fmt.Println(arr_data) 
		//arr_img := strings.Split( string(arr_data[8]) , `"`)
		//fmt.Println(arr_img[len(arr_img)-2])

		/*tconf["item_id"] = 				arr_data[0]
		tconf["item_org"] = 			arr_data[1] 
		tconf["item_name"] = 			arr_data[2]
		tconf["item_itemclass"] = 		arr_data[3] 
		tconf["item_itemclass_name"] = 	arr_data[4]  
		tconf["item_supplier"] = 		arr_data[5] 
		tconf["item_suppliername"] = 	arr_data[6] 
		tconf["item_remarks"] = 		arr_data[7] 
		tconf["item_image"] = 			helper.Img_to_render(`img_uploads/size_171/`+arr_data[1]+`/`+arr_data[0]+`.jpg` ,arr_data[0],arr_data[1])
		tconf["unit_price"] = 			arr_data[9] 
		tconf["unit_cost"] = 			arr_data[10] 
		tconf["sku"] = 					arr_data[11] 
		tconf["barcode"] = 				arr_data[12]
		tconf["uom"] = 					arr_data[13] 
		tconf["supplier_code"] = 		arr_data[14] 
		tconf["specs"] = 				arr_data[15] 
		tconf["status"] = 				arr_data[16] */

		tconf["item_id"] = 				arr_data[0]
		tconf["item_org"] = 			arr_data[1] 
		tconf["item_item"] = 			arr_data[2]
		tconf["item_name"] = 			arr_data[3]
		tconf["item_location"] = 		arr_data[4]
		tconf["item_remarks"] = 		arr_data[5] 
		tconf["item_image"] = 			arr_data[6]
		tconf["unit_price"] = 			arr_data[7] 
		tconf["unit_cost"] = 			arr_data[8] 
		tconf["sku"] = 					arr_data[9]  
		tconf["comm_per"] =				arr_data[10] 
		tconf["comm_amt"] = 			arr_data[11]
		tconf["uprice_date"] = 			arr_data[12]
		tconf["tax_per"] = 				arr_data[13] 
		tconf["tax_amt"] = 				arr_data[14] 
		tconf["ads_per"] = 				arr_data[15] 
		tconf["ads_amt"] = 				arr_data[16] 
		tconf["freight_per"] = 			arr_data[17] 
		tconf["freight_amt"] = 			arr_data[18]  
		tconf["ucost_date"] = 			arr_data[19]  
		tconf["item_group"] = 			arr_data[20]  

		
		//arr_data_itemclass := datatables.DataList(`select id,dbo.SIS_Itemclass_Name(1, 3, id) [itemclass_name] from sis_itemclass`  )
		//arr_data_supplier := datatables.DataList(`select id,dbo.SIS_Supplier_Name(1, 3, id) [supplier_name] from SIS_Supplier`  )
		arr_data_itemclass := datatables.DataList(`SIS_ItemClass_Get 2,`+str_OrgID  )
		arr_data_supplier := datatables.DataList(`SIS_Supplier_Get 2,`+str_OrgID  )
		arr_data_itemtag := datatables.DataList(`sis_itemtags_list 1,`+str_OrgID+` ,`+ tconf["h_id"] )
		arr_data_location := datatables.DataList(`exec SIS_Location_Get 6,`+str_OrgID  )
		tmpl := template.New("Hedit.html").Funcs(local_FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/location_items/Hedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute( w,&Data{rights,tconf,arr_data_itemclass,arr_data_supplier,arr_data_itemtag,arr_data_location} )
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
	login.Session_renew(w,r)
	
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["h_id"] = r.URL.Query().Get("h_id")
	   	tconf["delete_url"] = "/inventory/location_items/HDeleteHandler"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/location_items/modal_delete_loghdr.html"); err != nil {
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

		arr_data := datatables.Data_row(`SIS_Location_Item_Save 'Delete','`+username+`' ,  `+r.Form["h_id"][0] +`` )
		 

		js, err := json.Marshal(arr_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)	

		
	}	   	
}

func HDuplicateHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	//rights :="rights"

	if r.Method =="GET" {	

		username, _ := login.Get_account_info(r)

		tconf := make(map[string]string)
		tconf["h_id"] =r.URL.Query().Get("h_id")
		tconf["rights"]=rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID
		
		//tconf["data"] = datatables.DataList(sql)
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )
		
		//_,session_user_id 	:= login.Get_account_info(r)
		session_user_id := 1  //static here
		session_user_id_str := strconv.Itoa(session_user_id) 
		fmt.Println(session_user_id_str)

		type Data struct {
			Rights string
			Conf map[string]string
			Ls_itemClass [][]string
			Ls_supplier  [][]string
			Arr_itemtag  [][]string
			Ls_location [][]string
			
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		fmt.Println(`exec SIS_Item_Get 1, `+str_OrgID+`, '`+tconf["h_id"]+`'` )
		arr_data := datatables.Data_row(`exec SIS_Location_Item_Get 1, `+str_OrgID+`, '`+tconf["h_id"]+`'`  )
		
		//fmt.Println(arr_data) 
		//arr_img := strings.Split( string(arr_data[8]) , `"`)
		//fmt.Println(arr_img[len(arr_img)-2])

		/*tconf["item_id"] = 				arr_data[0]
		tconf["item_org"] = 			arr_data[1] 
		tconf["item_name"] = 			arr_data[2]
		tconf["item_itemclass"] = 		arr_data[3] 
		tconf["item_itemclass_name"] = 	arr_data[4]  
		tconf["item_supplier"] = 		arr_data[5] 
		tconf["item_suppliername"] = 	arr_data[6] 
		tconf["item_remarks"] = 		arr_data[7] 
		tconf["item_image"] = 			helper.Img_to_render(`img_uploads/size_171/`+arr_data[1]+`/`+arr_data[0]+`.jpg` ,arr_data[0],arr_data[1])
		tconf["unit_price"] = 			arr_data[9] 
		tconf["unit_cost"] = 			arr_data[10] 
		tconf["sku"] = 					arr_data[11] 
		tconf["barcode"] = 				arr_data[12]
		tconf["uom"] = 					arr_data[13] 
		tconf["supplier_code"] = 		arr_data[14] 
		tconf["specs"] = 				arr_data[15] 
		tconf["status"] = 				arr_data[16] */

		tconf["item_id"] = 				arr_data[0]
		tconf["item_org"] = 			arr_data[1] 
		tconf["item_item"] = 			arr_data[2]
		tconf["item_name"] = 			arr_data[3]
		tconf["item_location"] = 		arr_data[4]
		tconf["item_remarks"] = 		arr_data[5] 
		tconf["item_image"] = 			arr_data[6]
		tconf["unit_price"] = 			arr_data[7] 
		tconf["unit_cost"] = 			arr_data[8] 
		tconf["sku"] = 					arr_data[9]  
		tconf["comm_per"] =				arr_data[10] 
		tconf["comm_amt"] = 			arr_data[11]
		tconf["uprice_date"] = 			arr_data[12]
		tconf["tax_per"] = 				arr_data[13] 
		tconf["tax_amt"] = 				arr_data[14] 
		tconf["ads_per"] = 				arr_data[15] 
		tconf["ads_amt"] = 				arr_data[16] 
		tconf["freight_per"] = 			arr_data[17] 
		tconf["freight_amt"] = 			arr_data[18]  
		tconf["ucost_date"] = 			arr_data[19]   


		

		


		
		//arr_data_itemclass := datatables.DataList(`select id,dbo.SIS_Itemclass_Name(1, 3, id) [itemclass_name] from sis_itemclass`  )
		//arr_data_supplier := datatables.DataList(`select id,dbo.SIS_Supplier_Name(1, 3, id) [supplier_name] from SIS_Supplier`  )
		arr_data_itemclass := datatables.DataList(`SIS_ItemClass_Get 2,`+str_OrgID  )
		arr_data_supplier := datatables.DataList(`SIS_Supplier_Get 2,`+str_OrgID  )
		arr_data_itemtag := datatables.DataList(`sis_itemtags_list 1,`+str_OrgID+` ,`+ tconf["h_id"] )
		arr_data_location := datatables.DataList(`exec SIS_Location_Get 6,`+str_OrgID  )

		tmpl := template.New("Hduplicate.html").Funcs(local_FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/location_items/Hduplicate.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute( w,&Data{rights,tconf,arr_data_itemclass,arr_data_supplier,arr_data_itemtag,arr_data_location} )
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

func Ajax_getItem(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	getq := r.URL.Query().Get("q")
	getid := r.URL.Query().Get("id")
	isbaseonassignid :=r.URL.Query().Get("isbaseonassignid")
	assignmentID := r.URL.Query().Get("assignmentID")
	location := r.URL.Query().Get("location")

	fmt.Println(`is base sa assignID`,isbaseonassignid)
	fmt.Println(` assignID`,assignmentID)

	getq = getq + getid
		/*type Inner_item struct {
			Id string `json:"id"`
			Fullname string `json:"full_name"`
			Item_image string `json:"avatar"`
			Price string `json:"price"`

		}*/
		type Inner_item struct {
			Id string `json:"id"`
			Fullname string `json:"full_name"`
			Item_image string `json:"avatar"`
			Price string `json:"price"`
			Ucost string `json:"ucost"`
			Sku string `json:"sku"`
			Comm_per string `json:"comm_per"`
			Comm_amt string `json:"comm_amt"`
			Remarks string `json:"remarks"`
			Tax_per string `json:"tax_per"`
			Tax_amt string `json:"tax_amt"`
			Ads_per string `json:"ads_per"`
			Ads_amt string `json:"ads_amt"`
			Freight_per string `json:"freight_per"`
			Freight_amt string `json:"freight_amt"`
			Uprice_date string `json:"uprice_date"`
			Ucost_date string `json:"ucost_date"`




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
		 	sqlstr ="exec SIS_Item_Get 7, "+str_OrgID+",'"+getq+"',"+location
		 }else{
		 	//sqlstr = "exec LBR_Member_Get "+strtype+",  '"+getq+"'"
		 	sqlstr = "exec SIS_Item_Get 7, "+str_OrgID+",'"+getq+"',"+location


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
	        err = rows.Scan(&r.Id,&r.Fullname,&r.Item_image,&r.Price ,&r.Ucost,&r.Sku,&r.Comm_per,&r.Comm_amt,&r.Remarks,&r.Tax_per,&r.Tax_amt,&r.Ads_per,&r.Ads_amt,&r.Freight_per,&r.Freight_amt,&r.Uprice_date,&r.Ucost_date)
	        if err != nil {
	            panic(err.Error())
	        }
			count = count +1
			post := Inner_item {
	        	r.Id,
	        	r.Fullname, 
	        	config.Img_to_render(`img_uploads/size_80/`+str_OrgID+`/`+r.Item_image+`.jpg`,r.Item_image,80,str_OrgID), 
		    	r.Price,
		    	r.Ucost,
		    	r.Sku,
		    	r.Comm_per,
		    	r.Comm_amt,
		    	r.Remarks,
		    	r.Tax_per,
		    	r.Tax_amt,
		    	r.Ads_per,
		    	r.Ads_amt,
		    	r.Freight_per,
		    	r.Freight_amt,
		    	r.Uprice_date,
		    	r.Ucost_date,
		    }
		    InItem = append(InItem, post)
	    InItemRow = Inner_item{r.Id,r.Fullname, config.Img_to_render(`img_uploads/size_80/`+str_OrgID+`/`+r.Item_image+`.jpg`,r.Item_image,80,str_OrgID),r.Price ,r.Ucost,r.Sku ,r.Comm_per , r.Comm_amt , r.Remarks , r.Tax_per , r.Tax_amt , r.Ads_per , r.Ads_amt , r.Freight_per , r.Freight_amt, r.Uprice_date, r.Ucost_date}
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