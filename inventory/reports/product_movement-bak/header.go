package product_movement

import (
	
		"net/http"
		"fmt"
		"text/template"
		config "../../../config"
		"strconv"
		datatables "../../../datatables"
		login "../../../login"
		"strings"
		"encoding/json"
		"time"
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
	tconf["Panel_name"] ="Product Movement" 
	tconf["Add_new_button"] ="false" 
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="false"
	tconf["Add_new_button_details"] ="false" 
	tconf["upload_button"] ="false"   
	tconf["Add_new_button_url"] ="/inventory/supplier/HaddHandler?rights="+rights 
	tconf["Action_baseURL"] ="report/product_movement?rights="+rights
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
	var header = []string{"item_name","supplier_name","beg_vol","suprec_vol","cussal_vol","adjent_vol","end_vol"}  //column header
	if ln=="" {
		// content access here
		if( strings.Contains(dec_rights, "HAdd") ){
			tconf["Add_new_button"] ="true"
		}
		if( strings.Contains(dec_rights, "Dadd") ){
			tconf["Add_new_button"] ="true"
		}
			
		//end of content accesss
		var ajaxURLdata = "report/product_movement?rights="+rights
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("inventory/reports/product_movement/table.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		dr := r.URL.Query().Get("draw")
		item_id := r.URL.Query().Get("item_id")
		startdate := r.URL.Query().Get("startdate")
		enddate := r.URL.Query().Get("enddate")

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
		tconf["ActionCol"] ="false" // config for action column
	 	tconf["ActionCol_param"] ="ID" // config for parameter of action 
	 	tconf["ActionCol_edit"] ="false" // config for edit click
	 	tconf["ActionCol_edit_is_modal"] ="false" // config for edit click 
	 	//tconf["ActionCol_edit_url"] ="/timekeeping/overtime_logs/OvertimeLogsHeaderEdit?rights="+rights+"&h_id=" // config for edit click 
	 	tconf["ActionCol_edit_url"] ="/inventory/supplier/HEditHandler?rights="+rights+"&h_id=" // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/inventory/supplier/HDeleteHandler?h_id="
	 	tconf["ActionCol_detail"] ="false" // config for details click 
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
	 	tconf["col_img"] ="false"
	 	tconf["name_with_img"] ="true"
	 	tconf["col_img_num"] ="1"
	 	tconf["img_id_col"] ="0" // 
		tconf["dr"]=dr  //draw
		tconf["col_hidden1"] ="1" // what column to hide
		
		tconf["stat_url"]=`/inventory/report/product_movement_per_item?item_id=`  //link to statistical report

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
		if( strings.Contains(dec_rights, "Dadd") ){
			tconf["Add_new_button"] ="true"
		}
		//end rights for tables
		_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 

		fmt.Println(session_user_id_str)
		

		//tconf["sql_total"] = `SIS_Supplier_List 1, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		//tconf["sql_data"]  = `SIS_Supplier_List 0, `+str_OrgID+`,  `+str+`,  `+ln+`, '`+sr+`'`
		if(item_id==``){
			item_id = `0`
		}
		//fmt.Println(`Sis_rep_product_movement 4 ,`+str_OrgID+` ,`+item_id+` ,0 ,'`+startdate+`' ,'`+enddate+`' ,`+str+` ,`+ln)
		//tconf["sql_total"] = `Sis_rep_product_movement 4 ,`+str_OrgID+` ,`+item_id+` ,0 ,'`+startdate+`' ,'`+enddate+`' ,`+str+` ,`+ln
		//tconf["sql_data"]  = `Sis_rep_product_movement 1 ,`+str_OrgID+` ,`+item_id+` ,0 ,'`+startdate+`' ,'`+enddate+`' ,`+str+` ,`+ln
		
		tconf["sql_total"] = `SIS_Rep_Product_salemovement 0, `+str_OrgID+`, `+item_id+`, '',  0, '`+startdate+`', '`+enddate+`', `+str+`, `+ln
		tconf["sql_data"]  = `SIS_Rep_Product_salemovement 1, `+str_OrgID+`, `+item_id+`, '',  0, '`+startdate+`', '`+enddate+`', `+str+`, `+ln		
		DatatableretArray(w,tconf)
		  	 
		
		
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
		tconf["test_js"] = `alert("from webserver")`
		

		//arr_data := datatables.DataList(`select id,dbo.SIS_Itemclass_Name(1, 3, id) [itemclass_name] from sis_itemclass`  )
		//arr_data := datatables.DataList(`SIS_ItemClass_Get 2,`+str_OrgID  )
		//arr_data_supplier := datatables.DataList(`select id,dbo.SIS_Supplier_Name(1, 3, id) [supplier_name] from SIS_Supplier`  )

		//arr_data_supplier := datatables.DataList(`SIS_Supplier_Get 2,`+str_OrgID  )
		arr_data_itemtag := datatables.DataList(`sis_suptags_list 1,`+str_OrgID+` , 0`  )
		type Data struct {
			Rights	string	
			Tconf 	map[string]string
			//Ls_itemClass [][]string
			//Ls_supplier [][]string
			Arr_tag [][]string
		}  
		tmpl := template.New("Hadd.html").Funcs(local_FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/supplier/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights,tconf,arr_data_itemtag})
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
	        	arr_data_itemtag := datatables.Data_row(`SIS_suptags_save 'Update', '`+username+`',`+str_OrgID+`,`+item_id+`, `+tag+`, '`+value_input+`', '`+remarks+`'` )
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
	//rights :="static"

	if r.Method =="GET" {	


		//Org_id :=login.Get_session_org_id(r)
		//str_OrgID :=strconv.Itoa(Org_id)
		username, _ := login.Get_account_info(r)

		tconf := make(map[string]string)
		tconf["h_id"] =r.URL.Query().Get("h_id")
		tconf["rights"]=rights
		tconf["username"] = username
		tconf["org_id"] = str_OrgID
		//tconf["data"] = datatables.DataList(sql)
		//tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )
		
		_,session_user_id 	:= login.Get_account_info(r)
		// session_user_id := 1  //static here
		session_user_id_str := strconv.Itoa(session_user_id) 
		fmt.Println(session_user_id_str)
		

		type Data struct {
			Conf map[string]string
			//Ls_itemClass [][]string
			Arr_itemtag  [][]string
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		arr_data := datatables.Data_row(`exec SIS_Supplier_Get 1, '`+str_OrgID+`', '`+tconf["h_id"]+`'`  )
		
		fmt.Println(arr_data) 
		tconf["suprec_id"] = 		arr_data[0]
		tconf["item_org"] = 		arr_data[1] 
		tconf["item_name"] = 		arr_data[2] 
		tconf["item_remarks"] = 	arr_data[3] 
		
		//ls_itemclass:=datatables.DataList(`exec SIS_Supplier_Get 2, `+str_OrgID+``)
		arr_data_itemtag := datatables.DataList(`sis_suptags_list 1,`+str_OrgID+` ,`+ tconf["h_id"] )
		//fmt.Println(ls_itemclass)

		tmpl := template.New("Hedit.html").Funcs(local_FuncMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/supplier/Hedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute( w,&Data{tconf,arr_data_itemtag} )
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		session_username,_ := login.Get_account_info(r)	
		r.ParseForm()
		//_ ,err, ex_stat,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Edit', '`+session_username+`', `+r.Form["h_id"][0] +`, '`+r.Form["trandate"][0]+`', `+r.Form["lbr_assign"][0] +`, '`+r.Form["remarks"][0]+`'`,3)
		//_ ,err, ex_stat,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Edit', '`+session_username+`', `+r.Form["h_id"][0] +`, '`+r.Form["trandate"][0]+`', `+r.Form["lbr_assign"][0] +`, '`+r.Form["remarks"][0]+`'`,3)
		
		_ ,err, ex_stat,_ := config.Ap_sql(`exec DailySumHdr_Save 'Add', '`+session_username+`', `+r.Form["h_id"][0] +` , 1, `+r.Form["branch"][0]+`, '`+r.Form["date"][0]+`', '`+r.Form["remarks"][0]+`'`,3)
		
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
	   	tconf["delete_url"] = "/inventory/supplier/HDeleteHandler"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/supplier/modal_delete_loghdr.html"); err != nil {
			fmt.Println(err)
		}
		err = tmpl.Execute(w,tconf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} 
	}else{
		
		username, _ := login.Get_account_info(r)
		r.ParseForm()
		//_ ,err, _,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Delete', '`+session_username+`', `+r.Form["h_id"][0],3)
		//_ ,err, _,_ := config.Ap_sql(`Supplier_Save 'Delete','programmer',1, ` +r.Form["h_id"][0] ,3)
		//_ ,err, _,_ := config.Ap_sql(`SIS_ItemClass_Save 'Delete','programmer',` +r.Form["h_id"][0]  ,3)

		fmt.Println(`SIS_Supplier_Save 'Delete','`+username+`',` +r.Form["h_id"][0] )
		
		arr_data := datatables.Data_row(`SIS_Supplier_Save 'Delete','`+username+`',` +r.Form["h_id"][0]  )
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


func ViewitemHandler(w http.ResponseWriter, r *http.Request) {
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	tconf := make(map[string]string)
	

	tconf["img_name"] = r.URL.Query().Get("img_name")
	tconf["item_id"] = r.URL.Query().Get("item_id")
	tconf["graph_title"] = r.URL.Query().Get("item_name")
	tconf["item_image"] = config.Img_to_render_custom1( tconf["img_name"] ,171)
	tconf["item_name"] = r.URL.Query().Get("item_name")
	tconf["startdate"] = r.URL.Query().Get("startdate")
	tconf["enddate"] = r.URL.Query().Get("enddate")
	
	tconf["ddata"] = config.Decrypt( config.Encrypt_key, r.URL.Query().Get("ddata"))
	fmt.Println(tconf["ddata"])
	arr_ddata := strings.Split(	tconf["ddata"] , `][`)
	tconf["item_image"] 			= config.Img_to_render_custom1(arr_ddata[0],171)
	tconf["item_id_hideme"] 		= arr_ddata[1]
	tconf["item_name"]	 			= arr_ddata[2]
	tconf["supplier_name"] 			= arr_ddata[3]
	tconf["beg_vol"] 				= arr_ddata[4]
	tconf["suprec_vol"] 			= arr_ddata[5]
	tconf["cussal_vol"] 			= arr_ddata[6]
	tconf["adjent_vol"] 			= arr_ddata[7]
	tconf["end_vol"] 				= arr_ddata[8]
	
	fmt.Println(arr_ddata)
	fmt.Println(`SIS_Rep_Product_salemovement_detail 1, `+str_OrgID+`, `+tconf["item_id"]+`, '`+tconf["startdate"]+`', '`+tconf["enddate"]+`'`);
	graph1_Data := datatables.DataList(`SIS_Rep_Product_salemovement_detail 1, `+str_OrgID+`, `+tconf["item_id"]+`, '`+tconf["startdate"]+`', '`+tconf["enddate"]+`'` )
	
	fmt.Println(graph1_Data)
	tconf["graph1_data"] =` [`
	tconf["graph1_data2"] =` [`
	for _ , value := range graph1_Data{
			fmt.Println(value[0])
			tconf["graph1_data"] +=`['`+value[1]+`', `+value[2]+`],`
			tconf["graph1_data2"] +=`['`+value[1]+`', `+value[3]+`],`

	}
	tconf["graph1_data"] +=`]`
	tconf["graph1_data2"] +=`]`

	//
	//	
	fmt.Println(`SIS_Rep_Product_salemovement_detail 2, `+str_OrgID+`, `+tconf["item_id"]+`,'`+tconf["startdate"]+`', '`+tconf["enddate"]+`'`);
	rowret ,err, _,_ := config.Ap_sql(`SIS_Rep_Product_salemovement_detail 2, `+str_OrgID+`, `+tconf["item_id"]+`, '`+tconf["startdate"]+`', '`+tconf["enddate"]+`'`,1)
	if err != nil {
		panic(err.Error())
	}
	type Graph2 struct {
		Date_s	time.Time
	    Vals	string
	    Vals2   string
	}

	var RG2 Graph2
	if err != nil {
        fmt.Println(err.Error)
        panic(err.Error())
    }
    tconf["graph2_data"] =` [`
    tconf["graph2_data2"] =` [`

    for rowret.Next() {
	    err = rowret.Scan(&RG2.Date_s,&RG2.Vals,&RG2.Vals2)
	    if err != nil {
	        panic(err.Error())
	    }
	    fmt.Println(RG2.Date_s)
		fmt.Println(RG2.Vals)
		str_dd 	 :=strconv.Itoa(RG2.Date_s.Day())
		var i_mm int = int(RG2.Date_s.Month()) - 1
		str_mm 	 :=strconv.Itoa(i_mm) 
		str_yyyy :=strconv.Itoa(RG2.Date_s.Year())
		fmt.Println(`Day:` , str_dd)
		fmt.Println(`MM:` , str_mm)
		fmt.Println(`YYYY:` , str_yyyy)
		

		tconf["graph2_data"] +=`[Date.UTC(`+str_yyyy+`, `+str_mm+`, `+str_dd+`), `+RG2.Vals+`],`
		tconf["graph2_data2"] +=`[Date.UTC(`+str_yyyy+`, `+str_mm+`, `+str_dd+`), `+RG2.Vals2+`],`
	 	
	 	//r = RetMessage{r.Date_s,r.Date_s}
	}
	tconf["graph2_data"] +=`]`
	tconf["graph2_data2"] +=`]`
	
	tmpl := template.New("view_item_stat.html")
	//var err error
	if tmpl, err = tmpl.ParseFiles("inventory/reports/product_movement/view_item_stat.html"); err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w,tconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} 

}	



// custom datatable

func DatatableretArray(w http.ResponseWriter, conf map[string]string) {
	rows1 ,err1, _,_  := config.Ap_sql(conf["sql_total"],1)
	    var err error
	    if err1!= nil {
            panic(err1.Error())
	    }
	    type Records struct{
	    	Page_current int
	    	Page_more int
	    	Page_total int
	    	Record_more int
	    	Record_total int
	    }
	    var Recs Records
	    for rows1.Next() {
	    	err = rows1.Scan(&Recs.Page_current, &Recs.Page_more, &Recs.Page_total, &Recs.Record_more, &Recs.Record_total)
	            if err != nil {
	                    panic(err.Error())
	            }
	            
	    } 
	    //"data":[]
	 	 drr1,_:=strconv.Atoi(conf["dr"])
	    type Jsondata struct {
	    	Draw int `json:"draw"`
	    	RecordsTotal int `json:"recordsTotal"`
	    	RecordsFiltered int `json:"recordsFiltered"`
	    	Data [][]string `json:"data"`
	    }
	     type JsondataEmpty struct {
	    	total_count int `json:"total_count"`
	    	RecordsTotal int `json:"recordsTotal"`
	    	RecordsFiltered int `json:"recordsFiltered"`
	    	Data []string `json:"data"`
	    }
		
		objdata :=Datatable(conf["sql_data"],conf)
	    if  Recs.Record_total ==0{
	    	objdata1 := []string{} 
 			dd1 := &JsondataEmpty{drr1, Recs.Record_total , Recs.Record_total ,objdata1}
 			js, err := json.Marshal(dd1)
 			if err != nil {
		    	http.Error(w, err.Error(), http.StatusInternalServerError)
		    	return
	  		}
		  	w.Header().Set("Content-Type", "application/json")
		  	w.Write(js)	  
	    }else{
	    	dd := &Jsondata{drr1, Recs.Record_total , Recs.Record_total ,objdata}
		    js, err := json.Marshal(dd)
		    

			if err != nil {
		    	http.Error(w, err.Error(), http.StatusInternalServerError)
		    	return
	  		}
	  		w.Header().Set("Content-Type", "application/json")
		  	w.Write(js)	  
		}
}



func Datatable(sql string, conf map[string]string ) [][]string {
	
	var m  [][]string
	m1 :=make(map[string]int)
	var Actioncol bool
	Actioncol = true
	rows ,err, _,_ := config.Ap_sql(sql,1)
    if err != nil {
    	panic(err.Error())
    }
   	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
	}
	if cols == nil {
		fmt.Println(err)
	}
 	
	//getting the column name
	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
	}
	//geting columnheader key
	var hkey_action_param int
	for key , value := range cols{
			if ( strings.ToUpper(value)==strings.ToUpper(conf["ActionCol_param"])){
				hkey_action_param = key
			}
	}
	//end geting column key	
	draw1,_:=strconv.Atoi(conf["dr"] )
	m1["draw"] = draw1
	col_img_num,_ :=strconv.Atoi(conf["col_img_num"]) 
	img_id_col,_ :=strconv.Atoi(conf["img_id_col"]) 
	//column hide initialazation
	col_hidden1,_ :=strconv.Atoi(conf["col_hidden1"]) 
	if (col_hidden1==0 && conf["col_hidden1"]==`0`){
		col_hidden1=0
	}else{
		col_hidden1 = 1000
	}	
	//end column hide initialazation
//	m1["recordsTotal"] = strconv.Atoi(conf["recordsTotal"]) 
//	m1["recordsFiltered"] = strconv.Atoi(conf["recordsFiltered"]) 

	//type containerdata [][]interface{}
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//mars at work here
		var a []string
		var fontcolor  string
		var posted bool
		var Cancelled bool
		fontcolor =""
		fmt.Println(fontcolor)
		/*if( strings.Contains( PrintValue(vals[3].(*interface{})) , "Posted") ){
			fontcolor ="style='color:red;'"
		}*/
		if(len(cols)>3){
			if( strings.Contains( PrintValue(vals[3].(*interface{})) , "Posted") ){
				posted = true
			}
		}	
		if( strings.Contains( PrintValue(vals[2].(*interface{})) , "Posted") ){
			posted = true
		}
		if(len(cols)>3){
			if( strings.Contains( PrintValue(vals[3].(*interface{})) , "Cancelled") ){
				Cancelled = true
			}
		}	
		if( strings.Contains( PrintValue(vals[2].(*interface{})) , "Cancelled") ){
			Cancelled = true
		}
		
		for i := 1; i < len(cols); i++ {  // fetching data tables to container map array
			if Actioncol && i ==0{
				var EditAction , DeleteAction ,AddChildAction,AddroleAction,AddDetailAction,AssignmentAction ,MembersAction  string
				
				if conf["ActionCol_edit"] =="true"{
					if conf["ActionCol_edit_is_modal"] =="true"{
						EditAction = `<a data-toggle="modal" href="#basic" onclick="Modalfunc('`+conf["ActionCol_edit_url"] + PrintValue(vals[hkey_action_param].(*interface{}))+`',this);">Edit</a>`
					}else{
						//EditAction =`<a  name="`+conf["ActionCol_param"]+`' href ="#`+ conf["ActionCol_edit_url"] +`?id=`+PrintValue(vals[i].(*interface{})) +`" onclick="Nav2(this,`+PrintValue(vals[hkey_action_param].(*interface{})) +`)" >Edit</a>`
						EditAction =`<a href ="#`+ conf["ActionCol_edit_url"]+PrintValue(vals[i].(*interface{})) +`" onclick="func_edit('`+conf["ActionCol_edit_url"]+PrintValue(vals[hkey_action_param].(*interface{})) +`',this,event)" >Edit</a>`
					}
					if posted {
						EditAction = ``	
					}
					if Cancelled{
						EditAction = ``
					}	
				}
				if conf["ActionCol_delete"] =="true"{
					DeleteAction = `<a data-toggle="modal" href="#modalSession" onclick="Modalfuncdelete('`+conf["ActionCol_delete_url"]+ PrintValue(vals[hkey_action_param].(*interface{})) +`',1,this)" >Delete</a>`
					if posted {
						DeleteAction = ``	
					}
					if Cancelled {
						DeleteAction = ``
					}
				}	
				if conf["ActionCol_add_child"] =="true"{	
					//AddChildAction ="<a style='color: "+fontcolor+";' href ='#"+PrintValue(vals[i].(*interface{})) +"'>Add Child</a>"
					AddChildAction =`<a href ="#`+ conf["ActionCol_edit_url"]+PrintValue(vals[i].(*interface{})) +`" onclick="func_add_child('`+conf["ActionCol_add_child_url"]+PrintValue(vals[hkey_action_param].(*interface{})) +`',this)" >Add Child</a>`
				}
				if conf["ActionCol_roles"] =="true"{	
					AddroleAction ="<a  name='"+conf["ActionCol_param"]+"' href ='#"+ conf["ActionCol_roles_url"] +"?id="+PrintValue(vals[i].(*interface{})) +"' onclick='Nav2(this,"+PrintValue(vals[hkey_action_param].(*interface{})) +")' >Roles</a>"
				}
				if conf["ActionCol_detail"] =="true"{	
					AddDetailAction ="<a  name='"+conf["ActionCol_param"]+"' href ='#"+ conf["ActionCol_detail_url"] +"?id="+PrintValue(vals[i].(*interface{})) +"' onclick='Nav2(event,"+PrintValue(vals[hkey_action_param].(*interface{})) +")' >Details</a>"
				}
				if conf["ActionCol_assignment"] =="true"{	
					AssignmentAction =`<a href ="#`+ conf["ActionCol_assignment_url"]+PrintValue(vals[i].(*interface{})) +`" onclick="func_assignment('`+conf["ActionCol_assignment_url"]+PrintValue(vals[hkey_action_param].(*interface{})) +`',this)" >Assignment</a>`
				}
				if conf["ActionCol_members"] =="true"{	//on hris->assignment module button on members on assignment 
					MembersAction =`<a href ="#`+ conf["ActionCol_members_url"]+PrintValue(vals[i].(*interface{})) +`" onclick="func_assignment('`+conf["ActionCol_members_url"]+PrintValue(vals[hkey_action_param].(*interface{})) +`',this)" >Members</a>`
				}

				a  = append(a,EditAction+" "+DeleteAction+" "+AddChildAction +" "+ AddroleAction +" "+ AddDetailAction + " " +AssignmentAction + "  " + MembersAction)
			}
			//a  = append(a,"<div "+fontcolor+">" + PrintValue(vals[i].(*interface{})) + "</div>"  )
			if conf["col_img"]=="true" && i ==col_img_num{


				a  = append( a , `<img src="`+ config.Img_to_render(`img_uploads/size_80/`+  PrintValue(vals[0].(*interface{}))  +`.jpg` , PrintValue(vals[0].(*interface{})) ,80) + `" width="80" >` )			
			}
			if conf["name_with_img"]=="true" && i ==col_img_num {
				str_to_enc :=PrintValue(vals[0].(*interface{})) +`][`+PrintValue(vals[1].(*interface{})) +`][`+PrintValue(vals[2].(*interface{})) +`][`+PrintValue(vals[3].(*interface{})) +`][`+PrintValue(vals[4].(*interface{})) +`][`+PrintValue(vals[5].(*interface{})) +`][`+PrintValue(vals[6].(*interface{})) +`][`+PrintValue(vals[7].(*interface{})) +`][`+PrintValue(vals[8].(*interface{})) 
				enc_data := config.Encrypt(config.Encrypt_key,str_to_enc)
				a  = append( a ,  config.Img_to_render_custom1( PrintValue(vals[img_id_col].(*interface{})) ,80)  + `<a href="#`+conf["stat_url"]+PrintValue(vals[img_id_col].(*interface{}))+`" onclick="javascript:open_url('`+conf["stat_url"]+PrintValue(vals[1].(*interface{}))+`&img_name=`+PrintValue(vals[0].(*interface{}))+`','`+PrintValue(vals[2].(*interface{})) +`' ,'`+enc_data +`'	);" >`+PrintValue(vals[2].(*interface{})) +`</a>` )			
			}else{
				//	a  = append(a, PrintValue(vals[i].(*interface{}))  )
				
				if(i != col_hidden1 && i != 2){ //hide column
					if(i==14){ //is number, and align left 
						a  = append(a, `<div style="float: right;">`+PrintValue(vals[i].(*interface{})) + `</div>`  )
					}else{
						a  = append(a, PrintValue(vals[i].(*interface{}))  )
					}
				}
			}
		
		}
		m = append( m,a)
		//mars at work here
	}
		
	if rows.Err() != nil {
		fmt.Println(rows.Err)
	}	
	return m
}


func PrintValue(pval *interface{}) string{
	var a string
	switch v := (*pval).(type) {
	case nil:
//		fmt.Print("NULL")
		a =fmt.Sprint("")
	case bool:
		if v {
//			fmt.Print("1")
		a =	fmt.Sprint("1")
		} else {
//			fmt.Print("0")
		a =	fmt.Sprint("0")
		}
	case []byte:
//		fmt.Print(string(v))
		a = fmt.Sprint(string(v))
	case time.Time:
//		fmt.Print(v.Format("2006-01-02 15:04:05.999"))
		a = fmt.Sprint(v.Format("2006-01-02 15:04:05.999"))

	default:
//		fmt.Print(v)
		a= fmt.Sprint(v)
	}
	return a
}

//end custom datatable