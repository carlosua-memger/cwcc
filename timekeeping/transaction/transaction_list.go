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

)

type Profile struct {
  	Message string  // error status output or validated data
  	Status bool //pass or fail
  	ID string
}

func  TransactionHandler(w http.ResponseWriter, r *http.Request) {
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
	tconf["Panel_name"] ="Transaction Logs" 
	tconf["Add_new_button"] ="false" 
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="false"
	tconf["Add_new_button_details"] ="false" 
	tconf["upload_button"] ="false"   
	tconf["Add_new_button_url"] ="/timekeeping/CreateTransaction?rights="+rights
	tconf["Action_baseURL"] ="/timekeeping/CreateTransactionDetails?rights="+rights
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

	var header = []string{"Action", "Recno", "ID", "Trandate","Status","Assignment Name","logfile Name","Device","Remarks"}  //column header
	if ln=="" {
			
		// content access here
		if( strings.Contains(dec_rights, "HAdd") ){
			tconf["Add_new_button"] ="true"
		}
		if( strings.Contains(dec_rights, "Dadd ") ){

		}
			
		//end of content accesss
		var ajaxURLdata = "/timekeeping?rights="+rights
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("timekeeping/transaction/table.html"); err != nil {
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
			sortcol="2"
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
	 	tconf["ActionCol_edit_url"] ="/timekeeping/EditTransactionDetails?rights="+rights+"&h_id=" // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/timekeeping/CreateTransactionDelete?h_id="
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
		_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 

/*		fmt.Println("exec LBR_LogHdr_List 1, 2, 1, 1, 1, '"+sortColName+"', '', '"+sr+"'	" )
		fmt.Println(sortColName,str,sr)*/
		//tconf["sql_total"] = "exec LBR_LogHdr_List 1, 2, 1, 1, 50, 'id', '','' "	
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
	
		tconf["sql_total"] = "exec LBR_LogHdr_List 1, "+str_OrgID+", "+session_user_id_str+", "+str+", "+ln+", '"+sortColName+"', '', '"+sr+"'	"	
		//tconf["sql_data"] ="exec LBR_LogHdr_List 0, 2, 1, 1, 50, 'id', '','' "
		tconf["sql_data"] =   "exec LBR_LogHdr_List 0, "+str_OrgID+", "+session_user_id_str+", "+str+", "+ln+", '"+sortColName+"', '', '"+sr+"'	"
		datatables.DatatableretArray(w,tconf)
		  	 
		
		
		}	    
}

func DoCreateTransaction(trandate string, lbr_assign string, logfilename string, device string, remarks string,date_from string , date_to string) (bool ,int) {
	rowret ,err, _,_ := config.Ap_sql("LBR_LogHdr_Save 'Add','1','0','"+trandate+"','"+lbr_assign+"','"+logfilename+"','"+device+"','"+remarks+"','"+date_from+"','"+date_to+"'",1)
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

func CreateTransactionHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	_, ses_uID := login.Get_account_info(r)
	str_ses_uID :=strconv.Itoa(ses_uID)
	rights :=r.URL.Query().Get("rights")
	//fmt.Fprint(w,rights)
	if r.Method =="GET" {	
		// new struct for Org
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
	
		rows2 ,err, _,_ := config.Ap_sql("LBR_Assign_Get 2, "+str_OrgID+", "+str_ses_uID,1)
	    if err != nil {
	    	panic(err.Error())
	    }
		type row2 struct {
	            Id int64
	            Name  string
	    }     
		got2 := []row2{}
		    for rows2.Next() {
		            var r row2
		            err = rows2.Scan(&r.Id, &r.Name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post2 := row2{
		            	r.Id,
				        r.Name,
				      
				    }
					got2 = append(got2, post2)
		    }
		 //district dropdown
		/*type District struct {
			   DisID int
	           District string
	    }*/
	   
		    
		type Data struct {
			Assign []row2 
			Rights	string	
		}    
     	tmpl := template.New("CreateTransaction.html")
		//var err error
		if tmpl, err = tmpl.ParseFiles("timekeeping/transaction/CreateTransaction.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{got2,rights})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {
		r.ParseForm()
			
			add_status, lastinsertedID :=  DoCreateTransaction(r.Form["trandate"][0],r.Form["lbr_assign"][0], r.Form["logfilename"][0],r.Form["device"][0],r.Form["remarks"][0],r.Form["date_from"][0],r.Form["date_to"][0])

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

func EditTransactionDetailsHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	if r.Method =="GET" {	
		tconf := make(map[string]string)
		tconf["h_id"] =r.URL.Query().Get("h_id")
		tconf["rights"]=rights
		tconf["Assign_ID"] = strconv.Itoa(  LBR_LogHdr_Get_id(tconf["h_id"]).Lbr_assign )
		fmt.Println("assign id on head",tconf["Assign_ID"])
		// new struct for Org
		
		_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 

		rows2 ,err, _,_ := config.Ap_sql("LBR_Assign_Get 2,2,"+session_user_id_str,1)
	    if err != nil {
	    	panic(err.Error())
	    }
		type row2 struct {
	            Id int64
	            Name  string
	    }     
		got2 := []row2{}
		    for rows2.Next() {
		            var r row2
		            err = rows2.Scan(&r.Id, &r.Name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post2 := row2{
		            	r.Id,
				        r.Name,
				      
				    }
					got2 = append(got2, post2)
		    }
		type Data struct {
			Assign []row2 
			Conf map[string]string
			LBR_LogHdr LBR_LogHdr
		}    
	 	tmpl := template.New("edit_transaction_head.html")
		//var err error
		if tmpl, err = tmpl.ParseFiles("timekeeping/transaction/edit_transaction_head.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{got2,tconf,LBR_LogHdr_Get_id(tconf["h_id"])} )
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		session_username,_ := login.Get_account_info(r)	
		r.ParseForm()
		fmt.Println("exec LBR_LogHdr_Save 'Edit','"+session_username+"',"+r.Form["h_id"][0] +",'"+r.Form["trandate"][0] +"',"+r.Form["lbr_assign"][0] +",'"+r.Form["logfilename"][0]+"','"+r.Form["device"][0]+"','"+r.Form["remarks"][0]+"'")
		_ ,err, ex_stat,_ := config.Ap_sql("exec LBR_LogHdr_Save 'Edit','"+session_username+"',"+r.Form["h_id"][0] +",'"+r.Form["trandate"][0] +"',"+r.Form["lbr_assign"][0] +",'"+r.Form["logfilename"][0]+"','"+r.Form["device"][0]+"','"+r.Form["remarks"][0]+"'",3)
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

func CreateTransactionDeleteHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["h_id"] = r.URL.Query().Get("h_id")
	   	tconf["delete_url"] = "/timekeeping/CreateTransactionDelete"
	   	tmpl := template.New("modal_delete_loghdr.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("timekeeping/transaction/modal_delete_loghdr.html"); err != nil {
			fmt.Println(err)
		}
		err = tmpl.Execute(w,tconf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} 
	}else{
		session_username,_ := login.Get_account_info(r)	
		r.ParseForm()
		_ ,err, _,_ := config.Ap_sql("exec LBR_LogHdr_Save 'Delete','"+session_username+"',"+r.Form["h_id"][0],3)
		if err != nil {
	        fmt.Println(err.Error)
	        http.Error(w, err.Error(), http.StatusInternalServerError)
	        panic(err.Error())
		}
	}	   	
}
