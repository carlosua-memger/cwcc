package customer_payment

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

type Profile struct {
  	Message string  // error status output or validated data
  	Status bool //pass or fail
  	ID string
}




func HListHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)	
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
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
	tconf["Panel_name"] ="Customer Payment" 
	tconf["Add_new_button"] ="false" 
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="false"
	tconf["Add_new_button_details"] ="false" 
	tconf["upload_button"] ="false"   
	tconf["Add_new_button_url"] ="/inventory/customer_payment/HaddHandler?rights="+rights 
	tconf["Action_baseURL"] ="/inventory/customer_payment/details?rights="+rights
	tconf["ActionCol_application_url"] ="/inventory/helper/application?rights="+rights+"&transid=" //application link 	
	tconf["ActionCol_param"] ="ID"
    tconf["hdr_id"] = "0"
	tconf["Interpret_button"] ="false"
	tconf["Add_post_button"] ="false"
	tconf["Add_post_button_url"] =""
	tconf["Archive"] ="true" // is archive
	tconf["org"] =str_OrgID
	
	tconf["Add_cancel_button"] ="false"
	tconf["Add_cancel_button_url"] =""
	tconf["Add_upload_button"] ="false"
	tconf["Add_upload_button_url"] =""
	tconf["Add_duplicate_button"] ="false"
	tconf["Add_duplicate_button_url"] ="" 
	
	tconf["ismoreinfo"] = "false" // to description info	
	tconf["withtotal"] = "false"  // table total at bottom

	
	tconf["level"] ="1" 
	tconf["pageLength"] =last_length 
	tconf["displayStart"] =last_start
	tconf["last_search"] = last_search
	if last_length==""{
		tconf["pageLength"] ="10" 
		tconf["displayStart"] ="0"
	}
	var header = []string{"Action","ID" ,"Customer Name" , "Status" , "Trandate" ,"Pay Type" ,"Pay Ins","Pay Ref" ,"Amount","Application Balance","Application Status" ,"Remarks" }  //column header
	if ln=="" {



		// content access here
		if( strings.Contains(dec_rights, "HAdd") ){
			tconf["Add_new_button"] ="true"
			//fmt.Println("run herer")
		}
		if( strings.Contains(dec_rights, "Dadd ") ){

		}

		//end of content accesss
		var ajaxURLdata = "/inventory/customer_payment?rights="+rights
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html").Funcs(helper.FuncMap) 
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("inventory/customer_payment/table.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		dr := r.URL.Query().Get("draw")

		sr := r.URL.Query().Get("search[value]")
		archive := r.URL.Query().Get("arc")

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
	 	tconf["ActionCol_edit_url"] =  "/inventory/customer_payment/HEditHandler?rights="+rights+"&h_id=" // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/inventory/customer_payment/HDeleteHandler?h_id="
	 	tconf["ActionCol_detail"] ="false" // config for details click 
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	tconf["ActionCol_post"] ="false" // config for details click 
	 	tconf["ActionCol_post_url"] ="/inventory/customer_payment/HPostHandler?h_id="
	 	tconf["ActionCol_application"] ="true" // config for application click 

	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
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
		if( strings.Contains(dec_rights, "HPost") ){
			tconf["ActionCol_post"] ="true" // config for delete click 
	 	
		}
		/* if( strings.Contains(dec_rights, "HApplication") ){
			tconf["ActionCol_application"] ="true" // config for delete click 
	 	
		} */


		//end rights for tables

		_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 

		fmt.Println(session_user_id_str)
		
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		tconf["sql_total"] = `SIS_CusPayHdr_List 1, `+str_OrgID+`,0, `+str+`,  `+ln+`, '`+sr+`',`+archive
		//tconf["sql_data"]  = `SIS_Item_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		tconf["sql_data"]  = `SIS_CusPayHdr_List 0, `+str_OrgID+`,0,  `+str+`,  `+ln+`, '`+sr+`',`+archive
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
			retdata := datatables.DataList(`SIS_SupTags_Get 1, `+org_id+`, 0 ,`+tag_id)
			return retdata
		},
	}

func HAddHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	username, ses_uID := login.Get_account_info(r)
	
	str_ses_uID :=strconv.Itoa(ses_uID)
	rights :=r.URL.Query().Get("rights")
	fmt.Println(str_ses_uID)
	//fmt.Fprint(w,rights)
	
	if r.Method =="GET" {	
		
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		tconf := make(map[string]string)
		tconf["username"] =username
		tconf["Org_id"] = str_OrgID

		fromtype :=r.URL.Query().Get("fromtype")
		fmt.Println(`---------------testing print`)
		tconf["fromtype"] = fromtype
		if(fromtype==`co`){
			fmt.Println(fromtype)
			co_id :=r.URL.Query().Get("co_id")
			fmt.Println(`SIS_CusSalHdr_Get 3,`+str_OrgID+`,`+co_id )
			arr_data_co := datatables.Data_row(`SIS_CusSalHdr_Get 3,`+str_OrgID+`,`+co_id )
			tconf["co_id"] = 				arr_data_co[0]
			tconf["co_trandate"] = 			arr_data_co[1] 
			tconf["co_customer_name"] = 	arr_data_co[2] 
			tconf["co_customer_id"] = 	 	arr_data_co[3]
			tconf["co_amount"] = 	 		arr_data_co[4]
		
			
		}
		
		arr_data_custpaytype := datatables.DataList(`SIS_CusPaytype_Get 1,`+str_OrgID  )
		arr_data_customer := datatables.DataList(`SIS_Customer_Get 2,`+str_OrgID  )
		arr_data_bank := datatables.DataList(`SIS_FinIns_Get 2,`+str_OrgID  )
	
		
		type Data struct {
			Rights	string	
			Tconf 	map[string]string
			Ls_customerpaytype [][]string
			Ls_customer [][]string
			Ls_bank [][]string
			Arr_tag [][]string
		}   
		tmpl := template.New("Hadd.html").Funcs(local_FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/customer_payment/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights , tconf,arr_data_custpaytype,arr_data_customer,arr_data_bank,nil})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else{

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

func HAddTagHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method =="POST" {
		r.ParseForm()	
		item_id := r.Form["suprec_id"][0]
		username, _ := login.Get_account_info(r)
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		var returnData[] string	
		for key ,_ := range r.Form["tag"] {
			tag := r.Form["tag"][key]
			value_input := r.Form["value_input"][key]
			remarks := r.Form["remarks"][key]

			if(value_input!=""){
				fmt.Println(`SIS_CusSalTags_Save 'Update', '`+username+`',`+str_OrgID+`,`+item_id+`, `+tag+`, '`+value_input+`', '`+remarks+`'` )
	        	arr_data_itemtag := datatables.Data_row(`SIS_CusAdjTags_Save 'Update', '`+username+`',`+str_OrgID+`,`+item_id+`, `+tag+`, '`+value_input+`', '`+remarks+`'` )
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
	
		type Data struct {
			Rights string
			Conf map[string]string
			Ls_customer  [][]string
			Ls_customerpaytype  [][]string
			Ls_bank  [][]string
			Arr_tag  [][]string
		}
		fmt.Println(`exec SIS_CusPayHdr_Get 1, `+str_OrgID+`, '`+tconf["h_id"]+`'` )
		arr_data := datatables.Data_row(`exec SIS_CusPayHdr_Get 1, `+str_OrgID+`, '`+tconf["h_id"]+`'`  )
		
		fmt.Println(arr_data) 
		
		tconf["custpay_id"] = 				arr_data[0]
		tconf["custpay_status"] = 			arr_data[1] 
		tconf["custpay_transdate"] = 		arr_data[2] 
		tconf["custpay_custName"] = 		arr_data[3] 
		tconf["custpay_paytypeName"] = 		arr_data[4] 
		tconf["custpay_customer"] = 	 	arr_data[5]
		tconf["custpay_paytype"] = 			arr_data[6] 
		tconf["custpay_payIns"] = 			arr_data[7] 
		tconf["custpay_payRef"] = 			arr_data[8] 
		tconf["custpay_payAmt"] = 			arr_data[9] 
		tconf["custpay_remarks"] = 			arr_data[10] 
		tconf["custpay_checkdate"] = 		arr_data[11] 

		arr_data_custpaytype := datatables.DataList(`SIS_CusPaytype_Get 1,`+str_OrgID  )
		arr_data_customer := datatables.DataList(`SIS_Customer_Get 2,`+str_OrgID  )
		arr_data_bank := datatables.DataList(`SIS_FinIns_Get 2,`+str_OrgID  )

		tmpl := template.New("Hedit.html").Funcs(local_FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/customer_payment/Hedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute( w,&Data{rights,tconf,arr_data_customer,arr_data_custpaytype, arr_data_bank, nil} )
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		//session_username := `static`
		r.ParseForm()
	
		_ ,err, ex_stat,_ := config.Ap_sql(`exec DailySumHdr_Save 'Add', 'Mars', `+r.Form["h_id"][0] +` , 1, `+r.Form["branch"][0]+`, '`+r.Form["date"][0]+`', '`+r.Form["remarks"][0]+`'`,3)
		
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
	//rights :=r.URL.Query().Get("rights")
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["h_id"] = r.URL.Query().Get("h_id")
	   	tconf["delete_url"] = "/inventory/customer_payment/HDeleteHandler"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/customer_payment/modal_delete_loghdr.html"); err != nil {
			fmt.Println(err)
		}
		err = tmpl.Execute(w,tconf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} 
	}else{
		r.ParseForm()
		username, _ := login.Get_account_info(r)
		fmt.Println(`SIS_CusPayHdr_Save 'Delete','`+username+`' ,  `+r.Form["h_id"][0] +` , `+str_OrgID+``)
		arr_data := datatables.Data_row(`SIS_CusPayHdr_Save 'Delete','`+username+`' ,  `+r.Form["h_id"][0] +` , `+str_OrgID+`` )
		js, err := json.Marshal(arr_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)	
	}	   	
}


func HPostHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	//rights :=r.URL.Query().Get("rights")
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["h_id"] = r.URL.Query().Get("h_id")
	   	tconf["post_url"] = "/inventory/customer_payment/HPostHandler"
	   	tmpl := template.New("modal_post_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/customer_payment/modal_post_loghdr.html"); err != nil {
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
		fmt.Println(`SIS_CusPay_Post '`+username+`'  , `+str_OrgID+`,  `+r.Form["h_id"][0] +``)
		arr_data := datatables.Data_row(`SIS_CusPay_Post '`+username+`'  , `+str_OrgID+`,  `+r.Form["h_id"][0] +`` )
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