package timekeeping

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
		"time"

)

type Profile struct {
  	Message string  // error status output or validated data
  	Status bool //pass or fail
  	ID string
}

type RetMessage struct {
	RetNo	int
    RetMsg	string
}

func Holiday_listHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)	
	//fmt.Fprint(w,"mars test here")
	//fmt.Fprint(w, r.URL.Query().Get("rights"))
	rights :=r.URL.Query().Get("rights")
	last_start :=r.URL.Query().Get("last_start")
	last_length :=r.URL.Query().Get("last_length")
	last_search :=r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key,r.URL.Query().Get("rights"))		
	tconf := make(map[string]string)
	ln :=  r.URL.Query().Get("length")
	tconf["Panel_name"] ="Holidays" 
	tconf["Add_new_button"] ="false" 
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="false"
	tconf["Add_new_button_details"] ="false" 
	tconf["upload_button"] ="false"   
	tconf["Add_new_button_url"] ="/timekeeping/holidays/addnew?rights="+rights
	tconf["Action_baseURL"] ="/timekeeping/holidays?rights="+rights
	tconf["ActionCol_param"] ="ID"
    tconf["hdr_id"] = "0"
	tconf["Interpret_button"] ="false"
	tconf["Add_post_button"] ="false"
	tconf["Add_post_button_url"] =""
	
	tconf["level"] ="1" 
	tconf["pageLength"] =last_length 
	tconf["displayStart"] =last_start
	tconf["last_search"] = last_search
	if last_length==""{
		tconf["pageLength"] ="10" 
		tconf["displayStart"] ="0"
	}


	var header = []string{"Action", "ID", "Date","Area Name","Holiday Name","Special" ,"Regular","Regular 100","Remarks" }  //column header
	if ln=="" {
			
		// content access here
		if( strings.Contains(dec_rights, "HAdd") ){
			tconf["Add_new_button"] ="true"
		}
		if( strings.Contains(dec_rights, "Dadd ") ){

		}
			
		//end of content accesss
		var ajaxURLdata = "/timekeeping/holidays?rights="+rights
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("timekeeping/holidays/table.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		dr := r.URL.Query().Get("draw")
		sr := r.URL.Query().Get("search[value]")
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
		tconf["ActionCol"] ="true" // config for action column
	 	tconf["ActionCol_param"] ="ID" // config for parameter of action 
	 	tconf["ActionCol_edit"] ="false" // config for edit click
	 	tconf["ActionCol_edit_is_modal"] ="false" // config for edit click 
	 	tconf["ActionCol_edit_url"] ="/timekeeping/holidays/edit?rights="+rights+"&h_id=" // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/timekeeping/holidays/delete?h_id="
	 	tconf["ActionCol_detail"] ="false" // config for details click 
	 	//tconf["ActionCol_add_child"] ="true" // config for addchild click 
	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
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
		//end rights for tables


/*		fmt.Println("exec LBR_LogHdr_List 1, 2, 1, 1, 1, '"+sortColName+"', '', '"+sr+"'	" )
		fmt.Println(sortColName,str,sr)*/
		fmt.Println( `exec holiday_list 0 , `+str+`, `+ln+`, '`+sortColName+`', '', '`+sr+`'`)
		tconf["sql_total"] =  `exec holiday_list 1 , `+str+`, `+ln+`, '`+sortColName+`', '', '`+sr+`'`
		//tconf["sql_data"]  = `exec holiday_list 0, 2, `+session_user_id_str+` , `+str+`, `+ln+`,  '`+sortColName+`', '', '`+sr+`'`
		tconf["sql_data"]  = `exec holiday_list 0 , `+str+`, `+ln+`, '`+sortColName+`', '', '`+sr+`'`
		
		datatables.DatatableretArray(w,tconf)
		  	 
		
		
		}	    
}

func HolidayAddnewHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	username, _ := login.Get_account_info(r)
	//str_ses_uID :=strconv.Itoa(ses_uID)
	rights :=r.URL.Query().Get("rights")
	//fmt.Fprint(w,rights)
	if r.Method =="GET" {	
		// new struct for Org
		
	   
		    
		type Data struct {
			Rights	string	
		}    
     	tmpl := template.New("holidayAddnew.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("timekeeping/holidays/holidayAddnew.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {
		r.ParseForm()
			
			add_status, lastinsertedID :=  DoHolidayAddnewHandler(username,r.Form["holidaydate"][0],r.Form["name"][0],r.Form["special"][0],r.Form["regular"][0],r.Form["regular_100"][0],r.Form["remarks"][0])

			if( add_status){
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

			
		//}
	}		
}

func DoHolidayAddnewHandler(username,holidaydate string, name string, special string , regular string , regular_100 string ,remarks string ) (bool ,int) {
	
	//rowret ,err, _,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Add', '`+username+`', 0, '`+holidaydaye+`', `+lbr_assign+`,  '`+remarks+`'`,1)
	
	rowret ,err, _,_ := config.Ap_sql(`exec Holiday_Save 'Add' , '`+username+`' , 0 , '`+holidaydate+`' ,2 , '`+name+`' , `+special+`,`+regular+` , '`+regular_100+`' , '`+remarks+`' `,1)

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
type Holiday struct{
	ID 			int
	Holidaydate interface{}
	Area 		interface{} 
	Area_name 	interface{} 
	Name 		interface{} 
	Special 	interface{} 
	Regular 	interface{} 
	Regular_100 interface{} 
	Remarks 	interface{} 

}

func Holiday_Get( Hdr_id string ) Holiday {
	LBR_LogHdr_row ,err, _,_ := config.Ap_sql("exec Holiday_Get 1,"+ Hdr_id,1)
    if err != nil {
    	panic(err.Error())
    }
    var r Holiday
    //var date time.Time
    for LBR_LogHdr_row.Next() {
        err = LBR_LogHdr_row.Scan(&r.ID, &r.Holidaydate,&r.Area,&r.Area_name,&r.Name,&r.Special,&r.Regular,&r.Regular_100,&r.Remarks)
        if err != nil {
                panic(err.Error())
        }
	
    }

	return r
   
}

func HolidayEditHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	if r.Method =="GET" {	
		tconf := make(map[string]string)
		tconf["h_id"] =r.URL.Query().Get("h_id")
		tconf["rights"]=rights
		//tconf["Assign_ID"] = strconv.Itoa(  Holiday_Get(tconf["h_id"]).Lbr_assign )
		fmt.Println("holiday ID " , tconf["h_id"])
		// new struct for Org
		
		//_,session_user_id 	:= login.Get_account_info(r)
		//session_user_id_str := strconv.Itoa(session_user_id) 

		type Data struct {
			Conf map[string]string
			Holiday Holiday
		}


		funcMap := template.FuncMap{
		    "date_formater": func(a time.Time) string {
		        return a.Format(" Jan _2 2006")
		    },
			"func_intertaceToString": func(pval interface{}) string{
					var a string
					switch v := (pval).(type) {
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
						a = fmt.Sprint(v.Format("Jan _2 2006"))
						//a = fmt.Sprint(v.Format("Jan _2 2006 03:04:05PM"))
					default:
				//		fmt.Print(v)
						a= fmt.Sprint(v)
					}
					return a
				},
			}    
		tmpl := template.New("holidayEdit.html").Funcs(funcMap)
		var err error
		if tmpl, err = tmpl.ParseFiles("timekeeping/holidays/holidayEdit.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{tconf,Holiday_Get(tconf["h_id"])} )
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		session_username,_ := login.Get_account_info(r)	
		r.ParseForm()
		//Holiday_Save 'Edit' , 'Admin' , 15 , '2016-1-15' ,2 , 'Mars Day112121' , 0,0 , '' , 'dev test'
		//_ ,err, ex_stat,_ := config.Ap_sql(`exec Holiday_Save 'Edit', '`+session_username+`', `+r.Form["h_id"][0] +`, '`+r.Form["trandate"][0]+`', `+r.Form["lbr_assign"][0] +`, '`+r.Form["remarks"][0]+`'`,3)
		
		fmt.Println(`Holiday_Save 'Edit' , '`+session_username+`' , `+r.Form["h_id"][0]+` , '`+r.Form["holidaydate"][0]+`' ,2 , '`+r.Form["name"][0]+`' , `+r.Form["special"][0]+`,`+r.Form["regular"][0]+` , '`+r.Form["regular_100"][0]+`' , '`+r.Form["remarks"][0]+`'`)
		_ ,err, ex_stat,_ := config.Ap_sql(`Holiday_Save 'Edit' , '`+session_username+`' , `+r.Form["h_id"][0]+` , '`+r.Form["holidaydate"][0]+`' ,2 , '`+r.Form["name"][0]+`' , `+r.Form["special"][0]+`,`+r.Form["regular"][0]+` , '`+r.Form["regular_100"][0]+`' , '`+r.Form["remarks"][0]+`'`,3)

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

func HolidayDeleteHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["h_id"] = r.URL.Query().Get("h_id")
	   	tconf["delete_url"] = "/timekeeping/holidays/delete"
	   	tmpl := template.New("modal_delete.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("timekeeping/holidays/modal_delete.html"); err != nil {
			fmt.Println(err)
		}
		err = tmpl.Execute(w,tconf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} 
	}else{
		session_username,_ := login.Get_account_info(r)	
		r.ParseForm()
		//_ ,err, _,_ := config.Ap_sql(`exec LBR_OTHdr_Save 'Delete', '`+session_username+`', `+r.Form["h_id"][0],3)
		fmt.Println(`Holiday_Save 'Delete' , '`+session_username+`' ,` + r.Form["h_id"][0] + ` , '' ,0 , '' , 0,0 , '' , ''`)
		_ ,err, _,_ := config.Ap_sql(`Holiday_Save 'Delete' , '`+session_username+`' ,` + r.Form["h_id"][0] + ` , '' ,0 , '' , 0,0 , '' , ''`,3)
		
		if err != nil {
	        fmt.Println(err.Error)
	        http.Error(w, err.Error(), http.StatusInternalServerError)
	        panic(err.Error())
		}
	}	   	
}