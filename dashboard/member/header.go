package dashboard

import (
	
		"net/http"
		"fmt"
		"text/template"
		//config "../config"
		"strconv"
		"encoding/json"
		datatables "../../datatables"
		login "../../login"
		//"strings"
		//helper "../../helper"
)

type Profile struct {
  	Message string  // error status output or validated data
  	Status bool //pass or fail
  	ID string
}




func DashboardMember(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	username, ses_uID := login.Get_account_info(r)
	
	//username := "static"
	//ses_uID := 1
	str_ses_uID :=strconv.Itoa(ses_uID)
	rights :=r.URL.Query().Get("rights")
	fmt.Println(str_ses_uID)
	_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 

		fmt.Println(r.Method)
	if r.Method =="GET" {	
		
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		tconf := make(map[string]string)
		//tconf["parentID"] =r.URL.Query().Get("parentID")
		tconf["username"] =username
		tconf["Org_id"] = str_OrgID
		tconf["member_id"] = session_user_id_str
		tconf["rights"] = rights
		tconf["test_js"] = `alert("from webserver")`
		//fmt.Println(`SIS_ItemClass_Get 2,`+str_OrgID  )
		//fmt.Println(`SIS_Supplier_Get 2,`+str_OrgID  )
		arr_data := datatables.Data_row(`exec Members_get 2,`+str_OrgID+`,`+session_user_id_str  )

		tconf["name"] = 				arr_data[0]
		tconf["membersince"] = 			arr_data[1] 
		tconf["status"] = 				arr_data[2] 
		

		
		type Data struct {
			Rights	string	
			Tconf 	map[string]string
		
			
		}   
		tmpl := template.New("Hadd.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("dashboard/member/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights , tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {

		r.ParseForm()	
	//	img_dir := r.Form["Image_dr"][0]
	//	item_id := r.Form["Item_id"][0]
	//	Org_id :=login.Get_session_org_id(r)
	//	str_OrgID :=strconv.Itoa(Org_id)
	//	helper.Img_resizer(img_dir,item_id,str_OrgID)

		js, err := json.Marshal(``)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)	
		

	}	    
}
