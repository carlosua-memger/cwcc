package timekeeping
import (
		"net/http"
		"fmt"
		datatables "../../datatables"
		"strconv"
		"text/template"
		config "../../config"
		"encoding/json"
		login "../../login"
		"strings"
		)

		
var Ghid string

type LBR_LogHdr struct{
	ID int
	Transdate string
	Lbr_assign int 
	Logfilename string
	Device string
	Remarks string
	Status string

}

func LBR_LogHdr_Get_id( Hdr_id string ) LBR_LogHdr {
	LBR_LogHdr_row ,err, _,_ := config.Ap_sql("exec LBR_LogHdr_Get 1,"+ Hdr_id,1)
    if err != nil {
    	panic(err.Error())
    }
    var r LBR_LogHdr
    for LBR_LogHdr_row.Next() {
        err = LBR_LogHdr_row.Scan(&r.ID, &r.Transdate,&r.Lbr_assign,&r.Logfilename,&r.Device,&r.Remarks,&r.Status)
        if err != nil {
                panic(err.Error())
        }
	
    }

	return r
   
}

func  LBR_OTDtl_ListHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	dec_rights := config.Decrypt(config.Encrypt_key,r.URL.Query().Get("rights"))
	tconf := make(map[string]string)
	tconf["Assign_ID"] ="0"
	tconf["hdr_id"] = r.URL.Query().Get("ID") 
	LBR_hdr := LBR_OTHdr_Get_id(tconf["hdr_id"])
	LBR_hdr_ID_str:=strconv.Itoa( LBR_hdr.ID )
	if r.URL.Query().Get("ID")!=""{
		Ghid = r.URL.Query().Get("ID")
		
		tconf["Assign_ID"] =strconv.Itoa( LBR_hdr.Lbr_assign )
	}
	
	ln :=  r.URL.Query().Get("length")
	tconf["Panel_name"] ="Overtime Logs Details ID - ["+ LBR_hdr_ID_str+"]" 
	tconf["Add_new_button"] ="false"
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="true"
	tconf["Add_new_button_details"] ="false"  
	tconf["upload_button"] ="false" 
	tconf["Add_new_button_url"] ="/timekeeping/overtime_logs/details/OvertimeLogsDetailCreate?rights="+rights+"&hdr_id="
	tconf["Action_baseURL"] ="/timekeeping/overtime_logs/details?ID="+tconf["hdr_id"]+"&rights="+rights
	tconf["ActionCol_param"] ="ID"
	tconf["Interpret_button"] ="false"
	tconf["Interpret_button_url"] ="/timekeeping/overtime_logs/details/OvertimeLogsDetailInterpet?rights="+rights+"&hdr_id="+LBR_hdr_ID_str
	tconf["Add_post_button"] ="false"
	tconf["Add_post_button_url"] ="/timekeeping/overtime_logs/details/OvertimeLogsDetailPost?rights="+rights+"&hdr_id="+LBR_hdr_ID_str
	tconf["Add_cancel_button"] ="false"
	tconf["Add_cancel_button_url"] ="/timekeeping/overtime_logs/details/OvertimeLogsDetailCancel?rights="+rights+"&hdr_id="+LBR_hdr_ID_str
	tconf["level"] ="2" 
	tconf["pageLength"] ="10" 
	tconf["displayStart"] ="0"
	tconf["last_search"]=""
	
	var header = []string{"Action","ID", "Member_Name", "Applied","Hours","Remarks"}  //column header

	if ln=="" {
		// rights for tables 	
		if( strings.Contains(dec_rights, "DAdd") ){
			tconf["Add_new_button_details"] ="true"
		}
		if( strings.Contains(dec_rights, "DUpload") ){
			tconf["upload_button"] ="true"  
	 	
		}
		if( strings.Contains(dec_rights, "DPost") ){
			tconf["Add_post_button"] ="true" 
	 	
		}
		if( strings.Contains(dec_rights, "DInterpret") ){
	 		tconf["Interpret_button"] ="true"
	
	 	}
		if( strings.Contains(LBR_hdr.Status , "Posted") ){
	 		tconf["Add_new_button_details"] ="false"
	 		tconf["upload_button"] ="false" 
	 		tconf["Add_post_button"] ="false" 
	 		tconf["Interpret_button"] ="false"
	 		if( strings.Contains(dec_rights, "DCancel") ){
	 			tconf["Add_cancel_button"] ="true"
			}
	 	}
	 	if( strings.Contains(LBR_hdr.Status , "Cancelled") ){
	 		tconf["Add_new_button_details"] ="false"
	 		tconf["upload_button"] ="false" 
	 		tconf["Add_post_button"] ="false" 
	 		tconf["Interpret_button"] ="false"
	 		tconf["Add_cancel_button"] ="false"
	 	}
		//end rights for tables
	 	
		var ajaxURLdata = "/timekeeping/overtime_logs/details?ID="+tconf["hdr_id"]+"&rights="+rights
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("timekeeping/overtime_logs/table.html"); err != nil {
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
		//table ajax level config
		tconf["ActionCol"] ="true" // config for action column
	 	tconf["ActionCol_param"] ="ID" // config for parameter of action 
	 	tconf["ActionCol_edit"] ="false" // config for edit click 
	 	tconf["ActionCol_edit_is_modal"] ="true" // config for edit click 
	 	tconf["ActionCol_edit_url"] ="/timekeeping/overtime_logs/details/OvertimeLogsDetailEdit?rights="+rights+"&dtl_id=" // config for edit click  URL
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/timekeeping/overtime_logs/details/OvertimeLogsDetailDelete?dtl_id="
	 	//tconf["ActionCol_add_child"] ="true" // config for addchild click 
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
			tconf["ActionCol_delete"] ="true"  
	 	
		}
		fmt.Println(LBR_hdr.ID,`status:`,LBR_hdr.Status)
		if( strings.Contains(LBR_hdr.Status , "Posted") ){
	 		tconf["ActionCol_edit"] ="false"
	 		tconf["ActionCol_delete"] ="false"
	 	}

	 	if( strings.Contains(LBR_hdr.Status , "Cancelled") ){
	 		tconf["ActionCol_edit"] ="false"
	 		tconf["ActionCol_delete"] ="false"
	 	}
		/*if( strings.Contains(  LBR_hdr.Status , "Posted") ){
			tconf["ActionCol_edit"] ="false"
			tconf["ActionCol_delete"] ="false"
		}*/	
		//end rights for tables
		
		//tconf["sql_total"] = "exec LBR_LogHdr_List 1, 2, 1, 1, 50, 'id', '','' "	
		tconf["sql_total"] = `exec LBR_OTDtl_List 1, `+Ghid+`, `+str+`, `+r.URL.Query().Get("length")+`,  '`+sortColName+`', '', '`+sr+`'`
		//tconf["sql_data"] ="exec LBR_LogHdr_List 0, 2, 1, 1, 50, 'id', '','' "
		tconf["sql_data"] =`exec LBR_OTDtl_List 0, `+Ghid+`, `+str+`, `+r.URL.Query().Get("length")+`,  '`+sortColName+`', '', '`+sr+`'`
		datatables.DatatableretArray(w,tconf)
	}	    
}
type lbr_assign_name struct{
	Id int
	Name string
}
type htmldata struct {
	Htconf map[string]string
	Lbr_assign_name	[]lbr_assign_name
	//Members []Members
	Lbr_daycode_get []lbr_daycode_get
	LBR_OTDtl LBR_OTDtl	
}
func func_LBR_Assign_Get( AssignID string ,Username string) []lbr_assign_name {

	row ,err, _,_ := config.Ap_sql("exec LBR_Assign_Get 3, "+AssignID+", "+Username,1)
	res := []lbr_assign_name{}
	for row.Next() {
        var r lbr_assign_name
        err = row.Scan(&r.Id, &r.Name)
        if err != nil {
                panic(err.Error())
        }
		datarow := lbr_assign_name{
        	r.Id,
	        r.Name,
	      
	    }
		res = append(res, datarow)

    }
    return res		   		
}
type Members struct {
	Member_id int
	Name string
}
func func_members( assigmentID string,datefrom string,dateto string ) []Members {
	var mode string
	if datefrom == ""{
		mode = "1"
	}else{
		mode = "2"
	}
	fmt.Println("exec lbr_memass_search "+mode+", "+assigmentID+" , '' ,'"+datefrom+"','"+dateto+"'")
	row ,err, _,_ := config.Ap_sql("exec lbr_memass_search "+mode+", "+assigmentID+" , '' ,'"+datefrom+"','"+dateto+"'",1)
	res := []Members{}
	for row.Next() {
        var r Members
        err = row.Scan(&r.Member_id, &r.Name)
        if err != nil {
                panic(err.Error())
        }
		datarow := Members{
        	r.Member_id,
	        r.Name,
	      
	    }
		res = append(res, datarow)

    }
    return res
}
type lbr_daycode_get struct {
	Id int
	Name string	
}
func  Daycode() []lbr_daycode_get {
	row ,err, _,_ := config.Ap_sql("exec lbr_daycode_get 1",1)
	res := []lbr_daycode_get{}
	for row.Next() {
        var r lbr_daycode_get 
        err = row.Scan(&r.Id, &r.Name)
        if err != nil {
                panic(err.Error())
        }
		datarow := lbr_daycode_get{
        	r.Id,
	        r.Name,
	      
	    }
		res = append(res, datarow)
	}
    return res
}


func LBR_OTDtl_CreateHandler(w http.ResponseWriter, r *http.Request) {
	   	login.Session_renew(w,r)
	   	rights :=r.URL.Query().Get("rights")
	   	if r.Method =="GET" {

	   		tconf := make(map[string]string)
		   	tconf["hdr_id"] = r.URL.Query().Get("hdr_id") 
		   	tconf["rights"]= rights

		   	username, userID :=login.Get_account_info(r)
			fmt.Println ("username",username)
			fmt.Println ("userID",userID)

		    tmpl := template.New("OvertimelogsDetailCreate.html")
			var err error
			if tmpl, err = tmpl.ParseFiles("timekeeping/overtime_logs/OvertimelogsDetailCreate.html"); err != nil {
				fmt.Println(err)
			}
			type htmldata struct{
				Htconf map[string]string
			}
			err1 := tmpl.Execute(w,htmldata{tconf})
			if err1 != nil {
				http.Error(w, err1.Error(), http.StatusInternalServerError)
			}
		}else{
			session_username,_ := login.Get_account_info(r)	
			r.ParseForm()
			getid,_:=strconv.Atoi(r.Form["hdr_id"][0])
			type Profile struct {
			  	Message string  // error status output or validated data
			  	Status bool //pass or fail
			  	ID int
			}
			
			var err error 
			var ex_stat bool
			for _ , memberID := range r.Form["memberid"]{
				_ ,err, ex_stat,_ = config.Ap_sql(`LBR_OTDtl_Save 'Add', '`+session_username+`',0,`+r.Form["hdr_id"][0] +`,`+memberID+`,'`+r.Form["applied"][0]+`',`+r.Form["hours"][0]+`,'`+r.Form["remarks"][0]+`'`,3)
			}	
			//lastinsertedID, _ := res.LastInsertId()
			if err != nil {
			        fmt.Println(err.Error)
			        panic(err.Error())
			    }
			if ex_stat == true{
				profile := Profile{"Data Successfully added ",true ,getid}
				js, err := json.Marshal(profile)
				if err != nil {
			    	http.Error(w, err.Error(), http.StatusInternalServerError)
			    	return
		  		}
			  	fmt.Println("Data Successfully added")
			  	w.Header().Set("Content-Type", "application/json")
			  	w.Write(js)	
			}else{
				profile := Profile{"Oppsss something went wrong",false,0}
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

	type LBR_OTDtl struct {
		Id 			interface{}
		Hdr			interface{}
		Memberid 	interface{}	
		Member_name	interface{}
		Applied		interface{}
		Hours 		float64
		Remarks		interface{}
		
	}

func  func_LBR_OTDtl_byid( id string ) LBR_OTDtl {

	db_row ,err, _,_ := config.Ap_sql(`LBR_OTDtl_Get 1 ,` + id,1)
	var r LBR_OTDtl 
    for db_row.Next() {
        err = db_row.Scan(&r.Id, &r.Hdr, &r.Memberid, &r.Member_name, &r.Applied, &r.Hours, &r.Remarks)
        if err != nil {
                panic(err.Error())
        }
	}
    return r
}
func LBR_OTDtl_EditHandler(w http.ResponseWriter, r *http.Request) {
	   	login.Session_renew(w,r)
	   	rights :=r.URL.Query().Get("rights")
	   	if r.Method =="GET" {
		   	tconf := make(map[string]string)
		   	tconf["dtl_id"] = r.URL.Query().Get("dtl_id") 
			tconf["rights"]=rights
			AssignID :=  r.URL.Query().Get("Assign_ID") 
			fmt.Println("assignID oh",AssignID)	
			fmt.Fprint(w,AssignID)
		
		    //session_username,_ := login.Get_account_info(r)	
		    tmpl := template.New("OvertimelogsDetailEdit.html")
			var err error
			if tmpl, err = tmpl.ParseFiles("timekeeping/overtime_logs/OvertimelogsDetailEdit.html"); err != nil {
				fmt.Println(err)
			}
			err1 := tmpl.Execute(w,htmldata{tconf,nil,nil,func_LBR_OTDtl_byid(tconf["dtl_id"])})
			if err1 != nil {
				http.Error(w, err1.Error(), http.StatusInternalServerError)
			}
		}else{
			session_username,_ := login.Get_account_info(r)	
			r.ParseForm()
			//getid,_:=strconv.Atoi(r.Form["dtl_id"][0])
			hdr_id,_:=strconv.Atoi(r.Form["hdr_id"][0])
			type Profile struct {
			  	Message string  // error status output or validated data
			  	Status bool //pass or fail
			  	ID int
			}
			_ ,err, ex_stat,_ := config.Ap_sql(`LBR_OTDtl_Save 'Edit', '`+session_username+`',`+r.Form["dtl_id"][0] +`,`+r.Form["hdr_id"][0] +`,`+r.Form["memberid"][0]+`,'`+r.Form["applied"][0]+`',`+r.Form["hours"][0]+`,'`+r.Form["remarks"][0]+`'`,3)
			if err != nil {
			        fmt.Println(err.Error)
			        panic(err.Error())
			    }
			//lastinsertedID, _ := res.LastInsertId()
			if ex_stat == true{
				profile := Profile{"Data Successfully Update ",true ,hdr_id}
				js, err := json.Marshal(profile)
				if err != nil {
			    	http.Error(w, err.Error(), http.StatusInternalServerError)
			    	return
		  		}
			  	fmt.Println("Data Successfully Edited")
			  	w.Header().Set("Content-Type", "application/json")
			  	w.Write(js)	
			}else{
				profile := Profile{"Oppsss something went wrong",false,0}
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
func LBR_OTDtl_DeleteHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["dtl_id"] = r.URL.Query().Get("dtl_id")
	   	tconf["delete_url"] = "/timekeeping/overtime_logs/details/OvertimeLogsDetailDelete"
	   	tmpl := template.New("modal_delete.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("timekeeping/overtime_logs/modal_delete.html"); err != nil {
			fmt.Println(err)
		}
		err = tmpl.Execute(w,tconf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} 
	}else{
		session_username,_ := login.Get_account_info(r)	
		r.ParseForm()

		_ ,err, _,_ := config.Ap_sql(`LBR_OTDtl_Save 'Delete', '`+session_username+`',`+r.Form["dtl_id"][0],3)
		if err != nil {
	        fmt.Println(err.Error)
	        http.Error(w, err.Error(), http.StatusInternalServerError)
	        panic(err.Error())
		}

	}	   	
}

func LBR_OTDtl_PostHandler(w http.ResponseWriter, r *http.Request) {
	session_username,_ := login.Get_account_info(r)	
	var ValidateMsg string
	var validateRecCount int
	sqlstr_validate := `exec LBR_OTHdr_Validate 'Post',`+ r.URL.Query().Get("hdr_id")
	fmt.Println(sqlstr_validate)
	validate_rowret ,err, _,_ := config.Ap_sql(sqlstr_validate,1)
	if err != nil {
		panic(err.Error())
	}
	for validate_rowret.Next() {
	    err = validate_rowret.Scan(&ValidateMsg)
	    if err != nil {
	        panic(err.Error())
	    }
	 	validateRecCount = validateRecCount + 1
	}
	profile := Profile_ajax{ValidateMsg,true}
	if validateRecCount==0{
		sqlstr:= `exec LBR_OTHdr_Post '` +  session_username +`' , ` +r.URL.Query().Get("hdr_id")
		fmt.Println(sqlstr)
		rowret ,err, _,_ := config.Ap_sql(sqlstr,1)
		if err != nil {
			panic(err.Error())
		}
		var rtmsg RetMessage
		for rowret.Next() {
		    err = rowret.Scan(&rtmsg.RetNo,&rtmsg.RetMsg)
		    if err != nil {
		        panic(err.Error())
		    }
		 	rtmsg = RetMessage{rtmsg.RetNo,rtmsg.RetMsg}
		}
		profile = Profile_ajax{rtmsg.RetMsg,true}
		
	}
	js, err := json.Marshal(profile)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
		}
  	w.Header().Set("Content-Type", "application/json")
  	w.Write(js)

		
	   	
}

type RetMessage struct {
	RetNo	int
    RetMsg	string
}
type Profile_ajax struct {
	  	Message string  // error status output or validated data
	  	Status bool //pass or fail
}




func LBR_OTDtl_CancelHandler(w http.ResponseWriter, r *http.Request) {
	session_username,_ := login.Get_account_info(r)	
	sqlstr:= `exec LBR_OTHdr_Cancel '` +  session_username +`' , ` +r.URL.Query().Get("hdr_id")
	fmt.Println(sqlstr)
	rowret ,err, _,_ := config.Ap_sql(sqlstr,1)
	
	if err != nil {
		panic(err.Error())
	}
	var rtmsg RetMessage
	for rowret.Next() {
	    err = rowret.Scan(&rtmsg.RetNo,&rtmsg.RetMsg)
	    if err != nil {
	        panic(err.Error())
	    }
	 	rtmsg = RetMessage{rtmsg.RetNo,rtmsg.RetMsg}
	}
	profile := Profile_ajax{rtmsg.RetMsg,true}
		
	
	js, err := json.Marshal(profile)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
		}
  	w.Header().Set("Content-Type", "application/json")
  	w.Write(js)
}
