package dashboard

import (
	
		"net/http"
		"fmt"
		"text/template"
		config "../config"
		"strconv"
		"encoding/json"
		datatables "../datatables"
		login "../login"
		"strings"
		helper "../helper"
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
		//arr_data_location := datatables.DataList(`exec SIS_Location_Get 2,`+str_OrgID  )
		
		type Data struct {
			Rights	string	
			Tconf 	map[string]string
			
		}   
		tmpl := template.New("Hadd.html").Funcs(local_FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("dashboard/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights , tconf})
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
type RetMessage struct {
	RetNo	int
    RetMsg	string
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

		helper.Img_resizer(img_dir,item_id,`3`)

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
		tconf["item_image"] = 			helper.Img_to_render(`img_uploads/size_171/`+str_OrgID+`/`+arr_data[0]+`.jpg` ,arr_data[0], str_OrgID)
		

		


		
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