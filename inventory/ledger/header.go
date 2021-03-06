package ledger

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
		"time"
)

type Profile struct {
  	Message string  // error status output or validated data
  	Status bool //pass or fail
  	ID string
}




func HListHandler(w http.ResponseWriter, r *http.Request) {
	
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
		//arr_data_itemtag := datatables.DataList(`sis_itemtags_list 1,`+str_OrgID+` , 0`  )
		arr_data_location := datatables.DataList(`exec SIS_Location_Get 6,`+str_OrgID  )
		arr_data_item_status := datatables.DataList(`exec SIS_Item_Status_Get 0,`+str_OrgID  )
		arr_data_price := datatables.DataList(`SIS_Item_price_Get 0 ,`+str_OrgID  )
		arr_data_item_class := datatables.DataList(`SIS_ItemClass_Get 2 ,`+str_OrgID  )
		
		type Data struct {
			Rights	string	
			Tconf 	map[string]string
			Ls_itemClass [][]string
			Ls_supplier [][]string
			Arr_itemtag [][]string
			Arr_location [][]string
			Arr_item_status [][]string
			Arr_price [][]string
			Arr_item_class [][]string



		}   
		tmpl := template.New("Hadd.html").Funcs(local_FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/ledger/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights , tconf,nil,nil,nil,arr_data_location,arr_data_item_status,arr_data_price,arr_data_item_class})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {

		r.ParseForm()	
		img_dir := r.Form["Image_dr"][0]
		item_id := r.Form["Item_id"][0]
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		helper.Img_resizer(img_dir,item_id,str_OrgID)

		js, err := json.Marshal(`resizer success...`)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)	
		

	}

}

func ViewerHandler(w http.ResponseWriter, r *http.Request) {
	tconf := make(map[string]string)
	
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)

	tconf["rpt_type"] 	=r.URL.Query().Get("rpt_type")
	tconf["asofdate"] 	=r.URL.Query().Get("asofdate")
	tconf["item"] 		=r.URL.Query().Get("item")
	tconf["location"] 	=r.URL.Query().Get("location")
	tconf["datefrom"] 	=r.URL.Query().Get("datefrom")
	tconf["dateto"] 	=r.URL.Query().Get("dateto")
	tconf["org_id"] 	= str_OrgID

	tconf["item_key"] 		=r.URL.Query().Get("item_key")
	tconf["zeros"] 			=r.URL.Query().Get("zeros")
	tconf["item_status"]	=r.URL.Query().Get("item_status")
	tconf["sales_type"] 	=r.URL.Query().Get("sales_type")
	tconf["price"] 			=r.URL.Query().Get("price")
	tconf["batch"] 			=r.URL.Query().Get("batch")
	tconf["item_class"] 	=r.URL.Query().Get("item_class")
	
	var header []string
		
	if tconf["rpt_type"] ==`ledgerasof`{
		 header = []string{ "Item Name","Location Name", "Unit Price", "Status", "Ledger Running Bal."}  //column header
	}else if tconf["rpt_type"] ==`ledgerperiod`{
		 header = []string{"Location", "Item Name", "trans Date", "Qty", "Unit Cost", "Balance", "Total Cost", "Item"} //column header
	}
	//fmt.Println(`SIS_Ledger_Report_asof 0, 3, '`+tconf["asofdate"]+`', `+tconf["location"]+`, 0, 1, 20,''`)
	//arr_data := datatables.DataList(`SIS_Ledger_Report_asof 0, 3, '`+tconf["asofdate"]+`', `+tconf["location"]+`, 0, 1, 20,''`  )
	 
	type Data struct {
		Tconf 		map[string]string
		Header 		[]string
		Arr_data 	[][]string
	}
	var str_page string
	var str_page_dir string
	
	if tconf["rpt_type"] ==`ledgerasof`{
		str_page=`viewer.html`
		str_page_dir = `inventory/ledger/viewer.html`
	}else if tconf["rpt_type"] ==`ledgerperiod`{
		str_page=`viewer_period.html`
		str_page_dir = `inventory/ledger/viewer_period.html`
	}

	tmpl := template.New(str_page).Funcs(local_FuncMap) 
	var err error
	if tmpl, err = tmpl.ParseFiles(str_page_dir); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w,Data{tconf,header,nil})
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
} 

func AjaxViewerHandler(w http.ResponseWriter, r *http.Request) {
	tconf := make(map[string]string)

	tconf["rpt_type"] 	=r.URL.Query().Get("rpt_type")
	tconf["asofdate"] 	=r.URL.Query().Get("asofdate")
	tconf["item"] 		=r.URL.Query().Get("item")
	tconf["location"] 	=r.URL.Query().Get("location")
	tconf["page"] 		=r.URL.Query().Get("page")
	tconf["datefrom"] 	=r.URL.Query().Get("datefrom")
	tconf["dateto"] 	=r.URL.Query().Get("dateto")
	tconf["org_id"] 	=r.URL.Query().Get("org_id")

	tconf["item_key"] 		=r.URL.Query().Get("item_key")
	tconf["zeros"] 			=r.URL.Query().Get("zeros")
	tconf["item_status"]	=r.URL.Query().Get("item_status")
	tconf["sales_type"] 	=r.URL.Query().Get("sales_type")
	tconf["price"] 			=r.URL.Query().Get("price")
	tconf["batch"] 			=r.URL.Query().Get("batch")
	tconf["item_class"] 	=r.URL.Query().Get("item_class")
	


	if(tconf["item"]==""){
		tconf["item"]="0"
	}
	/*if(tconf["zeros"]==""){
		tconf["zeros"]="1"
	}
	if(tconf["sales_type"]==""){
		tconf["sales_type"]="1"
	}
	if(tconf["price"]==""){
		tconf["price"]="0"
	}
	if(tconf["item_class"]==""){
		tconf["item_class"]="0"
	}*/

	var arr_data [][]string
	if tconf["rpt_type"] ==`ledgerasof`{
		
		fmt.Println(`SIS_Ledger_Report_asof 0, `+tconf["org_id"]+`, '`+tconf["asofdate"]+`', `+tconf["location"]+`, `+tconf["item"]+`, `+tconf["page"]+`, 50,'`+tconf["item_key"]+`',`+tconf["zeros"] +`,'`+tconf["item_status"]+`',`+tconf["sales_type"]+`,`+tconf["price"]+`,'`+tconf["batch"]+`',`+tconf["item_class"] )
		arr_data = DataList(`SIS_Ledger_Report_asof 0, `+tconf["org_id"]+`, '`+tconf["asofdate"]+`', `+tconf["location"]+`, `+tconf["item"]+`, `+tconf["page"]+`, 50,'`+tconf["item_key"]+`',`+tconf["zeros"] +`,'`+tconf["item_status"]+`',`+tconf["sales_type"]+`,`+tconf["price"]+`,'`+tconf["batch"]+`',`+tconf["item_class"], tconf["org_id"] )

	}else if tconf["rpt_type"] ==`ledgerperiod`{
		fmt.Println(`SIS_Ledger_Report_period 0, `+tconf["org_id"]+`, '`+tconf["datefrom"]+`','`+tconf["dateto"]+`', `+tconf["location"]+`, `+tconf["item"]+` `+tconf["page"]+`, 50, ''  `)
		//uses custom datatables at buttom
		arr_data = DataList(`SIS_Ledger_Report_period 0, `+tconf["org_id"]+`, '`+tconf["datefrom"]+`','`+tconf["dateto"]+`', `+tconf["location"]+`, `+tconf["item"]+`, `+tconf["page"]+`, 50, ''  `, tconf["org_id"]  )
		//end
	}

	//fmt.Println(`SIS_Ledger_Report_asof 0, 3, '`+tconf["asofdate"]+`', `+tconf["location"]+`, 0, `+tconf["page"]+`, 25,''`)
	//arr_data := datatables.DataList(`SIS_Ledger_Report_asof 0, 3, '`+tconf["asofdate"]+`', `+tconf["location"]+`, 0, `+tconf["page"]+`, 25,''`  )
	js, err := json.Marshal(arr_data)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
		}
  	w.Header().Set("Content-Type", "application/json")
  	w.Write(js)	
} 



func AjaxViewerDetailsHandler(w http.ResponseWriter, r *http.Request) {
	tconf := make(map[string]string)
	
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)

	tconf["location"] 	=r.URL.Query().Get("location")
	tconf["item"] 	=r.URL.Query().Get("item")
	tconf["item_name"] 	=r.URL.Query().Get("item_name")
	tconf["location_name"] 	=r.URL.Query().Get("location_name")
	tconf["uprice"] 	=r.URL.Query().Get("uprice")
	tconf["status"] 	=r.URL.Query().Get("status")
	tconf["runbal"] 	=r.URL.Query().Get("runbal")
	
	
	fmt.Println(`sis_ledger_report_asof_detail 0 ,`+str_OrgID+`, `+tconf["location"]+`,`+tconf["item"])
	arr_data := datatables.DataList(`sis_ledger_report_asof_detail 0 ,`+str_OrgID+`, `+tconf["location"]+`,`+tconf["item"]   )
	 
	type Data struct {
		Tconf 		map[string]string
		Arr_data 	[][]string
	}
	var str_page string
	var str_page_dir string
	str_page=`viewer_details.html`
	str_page_dir = `inventory/ledger/viewer_details.html`

	tmpl := template.New(str_page).Funcs(local_FuncMap) 
	var err error
	if tmpl, err = tmpl.ParseFiles(str_page_dir); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w,Data{tconf,arr_data})
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
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
	

		//arr_data := datatables.DataList(`select id,dbo.SIS_Itemclass_Name(1, 3, id) [itemclass_name] from sis_itemclass`  )
		arr_data := datatables.DataList(`SIS_ItemClass_Get 2,`+str_OrgID  )
		//arr_data_supplier := datatables.DataList(`select id,dbo.SIS_Supplier_Name(1, 3, id) [supplier_name] from SIS_Supplier`  )

		arr_data_supplier := datatables.DataList(`SIS_Supplier_Get 2,`+str_OrgID  )
		arr_data_itemtag := datatables.DataList(`sis_itemtags_list 1,`+str_OrgID+` , 0`  )
		type Data struct {
			Rights	string	
			Tconf 	map[string]string
			Ls_itemClass [][]string
			Ls_supplier [][]string
			Arr_itemtag [][]string
		}   
		tmpl := template.New("Hadd.html").Funcs(local_FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/item/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights , tconf,arr_data,arr_data_supplier,arr_data_itemtag})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {

		r.ParseForm()	
		img_dir := r.Form["Image_dr"][0]
		item_id := r.Form["Item_id"][0]
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		helper.Img_resizer(img_dir,item_id,str_OrgID)

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
			
		}
		//arr_data := datatables.Data_row(`exec branch_get 1, 1,`+ tconf["h_id"]  )
		fmt.Println(`exec SIS_Item_Get 1, `+str_OrgID+`, '`+tconf["h_id"]+`'` )
		arr_data := datatables.Data_row(`exec SIS_Item_Get 1, `+str_OrgID+`, '`+tconf["h_id"]+`'`  )
		
		//fmt.Println(arr_data) 
		//arr_img := strings.Split( string(arr_data[8]) , `"`)
		//fmt.Println(arr_img[len(arr_img)-2])

		tconf["item_id"] = 				arr_data[0]
		tconf["item_org"] = 			arr_data[1] 
		tconf["item_name"] = 			arr_data[2]
		tconf["item_itemclass"] = 		arr_data[3] 
		tconf["item_itemclass_name"] = 	arr_data[4]  
		tconf["item_supplier"] = 		arr_data[5] 
		tconf["item_suppliername"] = 	arr_data[6] 
		tconf["item_remarks"] = 		arr_data[7] 
		tconf["item_image"] = 			helper.Img_to_render(`img_uploads/size_171/`+arr_data[1]+`/`+arr_data[0]+`.jpg` ,arr_data[0],arr_data[1])
		

		


		
		//arr_data_itemclass := datatables.DataList(`select id,dbo.SIS_Itemclass_Name(1, 3, id) [itemclass_name] from sis_itemclass`  )
		//arr_data_supplier := datatables.DataList(`select id,dbo.SIS_Supplier_Name(1, 3, id) [supplier_name] from SIS_Supplier`  )
		arr_data_itemclass := datatables.DataList(`SIS_ItemClass_Get 2,`+str_OrgID  )
		arr_data_supplier := datatables.DataList(`SIS_Supplier_Get 2,`+str_OrgID  )
		arr_data_itemtag := datatables.DataList(`sis_itemtags_list 1,`+str_OrgID+` ,`+ tconf["h_id"] )
		tmpl := template.New("Hedit.html").Funcs(local_FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/item/Hedit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute( w,&Data{rights,tconf,arr_data_itemclass,arr_data_supplier,arr_data_itemtag} )
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
	   	tconf["delete_url"] = "/inventory/item/HDeleteHandler"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/item/modal_delete_loghdr.html"); err != nil {
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
		arr_data := datatables.Data_row(`SIS_Item_Save 'Delete','`+username+`',` +r.Form["h_id"][0]  )
		 

		js, err := json.Marshal(arr_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)	

		
	}	   	
}



func DataList(sql string, org string ) [][]string {
	
	var m  [][]string
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
	//type containerdata [][]interface{}
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//mars at work here
		var a []string
		
		for i := 0; i < len(cols); i++ {  // fetching data tables to container map array
			if(i==4){
				a  = append(a, config.Img_to_render_custom1(PrintValue(vals[i].(*interface{}) ),80,org )  )
			}else{
				a  = append(a, PrintValue(vals[i].(*interface{}))  )
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


func Ajax_getItem(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	getq := r.URL.Query().Get("q")
	getid := r.URL.Query().Get("id")
	isbaseonassignid :=r.URL.Query().Get("isbaseonassignid")

	fmt.Println(`is base sa assignID`,isbaseonassignid)


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
		 	sqlstr ="exec SIS_Item_Get 2, "+str_OrgID+",'"+getq+"'"
		 }else{
		 	//sqlstr = "exec LBR_Member_Get "+strtype+",  '"+getq+"'"
		 	sqlstr = "exec SIS_Item_Get 2, "+str_OrgID+",'"+getq+"'"


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
	    InItemRow = Inner_item{r.Id,r.Item_image, config.Img_to_render(`img_uploads/size_80/`+str_OrgID+`/` +r.Item_image+`.jpg`,r.Item_image,80,str_OrgID)}
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




