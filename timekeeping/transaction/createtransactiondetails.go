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

func  CreateTransactionDetailsHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	dec_rights := config.Decrypt(config.Encrypt_key,r.URL.Query().Get("rights"))
	tconf := make(map[string]string)
	tconf["Assign_ID"] ="0"
	tconf["hdr_id"] = r.URL.Query().Get("ID") 
	LBR_hdr := LBR_LogHdr_Get_id(tconf["hdr_id"])
	LBR_hdr_ID_str:=strconv.Itoa( LBR_hdr.ID )
	if r.URL.Query().Get("ID")!=""{
		Ghid = r.URL.Query().Get("ID")
		
		tconf["Assign_ID"] =strconv.Itoa( LBR_hdr.Lbr_assign )
	}
	
	ln :=  r.URL.Query().Get("length")
	tconf["Panel_name"] ="Transaction Details ID - ["+ LBR_hdr_ID_str+"]" 
	tconf["Add_new_button"] ="false"
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="true"
	tconf["Add_new_button_details"] ="false"  
	tconf["upload_button"] ="false" 
	tconf["Add_new_button_url"] ="/timekeeping/CreateTransactionDetailsAddnew?rights="+rights+"&hdr_id="
	tconf["Action_baseURL"] ="/timekeeping/CreateTransactionDetails?ID="+tconf["hdr_id"]+"&rights="+rights
	tconf["ActionCol_param"] ="ID"
	tconf["Interpret_button"] ="false"
	tconf["Interpret_button_url"] ="/timekeeping/transactionDetails/interpret?rights="+rights+"&hdr_id="+LBR_hdr_ID_str
	tconf["Add_post_button"] ="false"
	tconf["Add_post_button_url"] ="/timekeeping/transactionDetails/post?rights="+rights+"&hdr_id="+LBR_hdr_ID_str
	tconf["level"] ="2" 
	tconf["pageLength"] ="10" 
	tconf["displayStart"] ="0"
	tconf["last_search"]=""
	var header = []string{"Action", "RecNo","ID", "Member_name", "Assigment","DayCode","Device","logtime","log type","logoutput","Remarks"}  //column header

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
	 	}
		//end rights for tables
	 	

		var ajaxURLdata = "/timekeeping/CreateTransactionDetails?ID="+tconf["hdr_id"]+"&rights="+rights
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
		//table ajax level config
		tconf["ActionCol"] ="true" // config for action column
	 	tconf["ActionCol_param"] ="ID" // config for parameter of action 
	 	tconf["ActionCol_edit"] ="false" // config for edit click 
	 	tconf["ActionCol_edit_is_modal"] ="true" // config for edit click 
	 	tconf["ActionCol_edit_url"] ="/timekeeping/CreateTransactionDetailsEdit?rights="+rights+"&dtl_id=" // config for edit click  URL
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/timekeeping/CreateTransactionDetailsDelete?dtl_id="
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
		if( strings.Contains(LBR_hdr.Status , "Posted") ){
			tconf["ActionCol_edit"] ="false"
	 		tconf["ActionCol_delete"] ="false"	
	 	}
		/*if( strings.Contains(  LBR_hdr.Status , "Posted") ){
			tconf["ActionCol_edit"] ="false"
			tconf["ActionCol_delete"] ="false"
		}*/	
		//end rights for tables
		fmt.Println("exec LBR_LogHdr_List 1, 2, 1, 1, 1, '"+sortColName+"', '', '"+sr+"'	" )
		fmt.Println(sortColName,str,sr)
		//tconf["sql_total"] = "exec LBR_LogHdr_List 1, 2, 1, 1, 50, 'id', '','' "	
		tconf["sql_total"] = "exec LBR_LogDtl_List 1, "+ Ghid +", "+str+", "+r.URL.Query().Get("length")+",  'id', '', '"+sr+"'"	
		//tconf["sql_data"] ="exec LBR_LogHdr_List 0, 2, 1, 1, 50, 'id', '','' "
		tconf["sql_data"] ="exec LBR_LogDtl_List 0, "+ Ghid +", "+str+" , "+r.URL.Query().Get("length")+",  'id', '', '"+sr+"' "
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
	LBR_LOGDTL []LBR_LOGDTL	
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
	Member_id string
	Name string
}
func func_members( assigmentID string,datefrom string,dateto string ) []Members {
	var mode string
	if datefrom == ""{
		mode = "1"
	}else{
		mode = "5"
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


func CreateTransactionDetailsAddnewHAndler(w http.ResponseWriter, r *http.Request) {
	   	login.Session_renew(w,r)
	   	rights :=r.URL.Query().Get("rights")
	   	if r.Method =="GET" {
		   	session_username,_ := login.Get_account_info(r)	

	   		tconf := make(map[string]string)
		   	tconf["hdr_id"] = r.URL.Query().Get("hdr_id") 
		   	tconf["rights"]= rights
		   	AssignID := strconv.Itoa( LBR_LogHdr_Get_id( r.URL.Query().Get("hdr_id")).Lbr_assign )

		   	username, userID :=login.Get_account_info(r)
			fmt.Println ("username",username)
			fmt.Println ("userID",userID)

		    tmpl := template.New("CreateTransactionDetailsAddnew.html")
			var err error
			if tmpl, err = tmpl.ParseFiles("timekeeping/transaction/CreateTransactionDetailsAddnew.html"); err != nil {
				fmt.Println(err)
			}
			err1 := tmpl.Execute(w,htmldata{tconf,func_LBR_Assign_Get(AssignID ,session_username),Daycode(),nil})
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
			var output string
			output = r.Form["output"][0]
			if r.Form["output"][0] == "" {
				output = "0"
			}
			var err error 
			var ex_stat bool
			for _ , memberID := range r.Form["memberid"]{
				                       fmt.Println("exec LBR_LogDtl_Save 'Add','"+session_username+"',0,"+r.Form["hdr_id"][0] +","+ memberID +","+r.Form["lbr_assign"][0]+","+r.Form["daycode_id"][0]+",'','"+r.Form["trndatetime"][0]+"','"+r.Form["trntype"][0]+"',"+output+",'"+r.Form["remarks"][0]+"' ")
				_ ,err, ex_stat,_ = config.Ap_sql("exec LBR_LogDtl_Save 'Add','"+session_username+"',0,"+r.Form["hdr_id"][0] +","+ memberID +","+r.Form["lbr_assign"][0]+","+r.Form["daycode_id"][0]+",'','"+r.Form["trndatetime"][0]+"','"+r.Form["trntype"][0]+"',"+output+",'"+r.Form["remarks"][0]+"' ",3)
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

	type LBR_LOGDTL struct {
		Id 			int64
		Hdr			int
		Lbr_assign 	int	
		Memberid 	int	
		Member_name	string
		Lbr_daycode	int
		Logtime		string
		Logtype		string
		Logoutput	float64
		Remarks     string
	}

func  func_get_LBR_LOGDTL_byid( id string ) []LBR_LOGDTL {
	row ,err, _,_ := config.Ap_sql("exec LBR_LogDtl_Get 1, "+id,1)
	res := []LBR_LOGDTL{}
	for row.Next() {
        var r LBR_LOGDTL 
        err = row.Scan(&r.Id, &r.Hdr, &r.Lbr_assign, &r.Memberid, &r.Member_name, &r.Lbr_daycode, &r.Logtime, &r.Logtype, &r.Logoutput, &r.Remarks)
        if err != nil {
                panic(err.Error())
        }
		datarow := LBR_LOGDTL{
        	r.Id ,			
			r.Hdr,			
			r.Lbr_assign ,	
			r.Memberid ,		
			r.Member_name	,
			r.Lbr_daycode	,
			r.Logtime	,	
			r.Logtype	,	
			r.Logoutput	,
			r.Remarks   ,  
	      
	    }
		res = append(res, datarow)
	}
    return res
}
func CreateTransactionDetailsEditHAndler(w http.ResponseWriter, r *http.Request) {
	   	login.Session_renew(w,r)
	   	rights :=r.URL.Query().Get("rights")
	   	if r.Method =="GET" {
		   	tconf := make(map[string]string)
		   	tconf["dtl_id"] = r.URL.Query().Get("dtl_id") 
			tconf["rights"]=rights
			AssignID :=  r.URL.Query().Get("Assign_ID") 
			fmt.Println("assignID oh",AssignID)	
			fmt.Fprint(w,AssignID)
		/*	var Hdr int 
			var Member sql.NullInt64 
			var Lbr_assign int 
			var Lbr_daycode int 
			var Logtime time.Time 
			var Logtype string 
			var Logoutput float64	
			var Remarks string	    
		    
		    rows ,err, _,_ := config.Ap_sql("select Hdr,Member,Lbr_assign,lbr_daycode,logtime,logtype,logoutput,remarks from LBR_LOGDTL where id =?", tconf["dtl_id"] ).Scan(&name)
		    err = rows.Scan(&r.Id,&r.Name,&r.Parent, &r.Url, &r.Param , &r.Tooltip,&r.Rights, &r.Cte_level, &r.Cte_id, &r.Cte_name)
		    switch {
		    case err == sql.ErrNoRows:
		            log.Printf("No user with that ID.")
		    case err != nil:
		            log.Fatal(err)
		    default:
		            fmt.Printf("Username is %s\n", username)
		    }
*/

		    session_username,_ := login.Get_account_info(r)	
		    tmpl := template.New("CreateTransactionDetailsEdit.html")
			var err error
			if tmpl, err = tmpl.ParseFiles("timekeeping/transaction/CreateTransactionDetailsEdit.html"); err != nil {
				fmt.Println(err)
			}
			err1 := tmpl.Execute(w,htmldata{tconf,func_LBR_Assign_Get(AssignID ,session_username ),Daycode(),func_get_LBR_LOGDTL_byid(tconf["dtl_id"])})
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
			var output string
			output = r.Form["output"][0]
			if r.Form["output"][0] == "" {
				output = "0"
			}
			
			_ ,err, ex_stat,_ := config.Ap_sql("exec LBR_LogDtl_Save 'Edit','"+session_username+"',"+r.Form["dtl_id"][0] +",0,"+r.Form["memberid"][0]+","+r.Form["lbr_assign"][0]+","+r.Form["daycode_id"][0]+",'','"+r.Form["trndatetime"][0]+"','"+r.Form["trntype"][0]+"',"+output+",'"+r.Form["remarks"][0]+"' ",3)
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
func CreateTransactionDetailsDeleteHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {
	   	tconf := make(map[string]string)
	   	tconf["dtl_id"] = r.URL.Query().Get("dtl_id")
	   	tconf["delete_url"] = "/timekeeping/CreateTransactionDetailsDelete"
	   	tmpl := template.New("modal_delete.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("timekeeping/transaction/modal_delete.html"); err != nil {
			fmt.Println(err)
		}
		err = tmpl.Execute(w,tconf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} 
	}else{
		session_username,_ := login.Get_account_info(r)	
		r.ParseForm()
		_ ,err, _,_ := config.Ap_sql("exec LBR_LogDtl_Save 'Delete','"+session_username+"',"+r.Form["dtl_id"][0],3)
		if err != nil {
		        fmt.Println(err.Error)
		        http.Error(w, err.Error(), http.StatusInternalServerError)
		        panic(err.Error())
		}

	}	   	
}

func CreateTransactionDetailsPostHAndler(w http.ResponseWriter, r *http.Request) {
	session_username,_ := login.Get_account_info(r)	
	//var ValidateMsg string
	var validateRecCount int
	Error_Message:=""
	profile := Profile_ajax{"",true}
	sqlstr_validate := `exec LBR_LogHdr_Validate 'Post',`+ r.URL.Query().Get("hdr_id")
	fmt.Println(sqlstr_validate)
	validate_rowret ,err, _,_ := config.Ap_sql(sqlstr_validate,1)
	if err != nil {
		panic(err.Error())
	}
	for validate_rowret.Next() {
	    err = validate_rowret.Scan(&Error_Message)
	    if err != nil {
	        panic(err.Error())
	    }
	 	validateRecCount = validateRecCount + 1
	}
	profile = Profile_ajax{Error_Message,true}
	if validateRecCount==0{
		sqlstr:= `exec LBR_LogHdr_Post '` +  session_username +`' , ` +r.URL.Query().Get("hdr_id")
		fmt.Println(sqlstr)
		
		_ ,err, _,_:= config.Ap_sql(sqlstr,3)
		//rowret ,err, _,_ := config.Ap_sql(sqlstr,1)
		if err != nil {
			panic(err.Error())
		}
		/*var rtmsg RetMessage
		for rowret.Next() {
		    err = rowret.Scan(&rtmsg.RetNo,&rtmsg.RetMsg)
		    if err != nil {
		        panic(err.Error())
		    }
		 	rtmsg = RetMessage{rtmsg.RetNo,rtmsg.RetMsg}
		}
		profile = Profile_ajax{rtmsg.RetMsg,true}*/
		
	}
	js, err := json.Marshal(profile)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
		}
  	w.Header().Set("Content-Type", "application/json")
  	w.Write(js)

		
	   	
}

func CreateTransactionDetailsInterpretHAndler(w http.ResponseWriter, r *http.Request) {
	session_username,_ := login.Get_account_info(r)	
  	//var ValidateMsg string
	validateRecCount :=0
	Error_Message:=""
	sqlstr_validate := `exec LBR_LogHdr_Validate 'Interpret',`+ r.URL.Query().Get("hdr_id")
	fmt.Println(sqlstr_validate)
	validate_rowret ,err, _,_ := config.Ap_sql(sqlstr_validate,1)
	if err != nil {
		panic(err.Error())
	}
	for validate_rowret.Next() {
	    err = validate_rowret.Scan(&Error_Message)
	    if err != nil {
	        panic(err.Error())
	    }
	 	validateRecCount = validateRecCount + 1
	}
	profile := Profile_ajax{Error_Message,true}
	if validateRecCount==0{
		sqlstr:= `exec LBR_LogHdr_Interpret '`+session_username+`',`+ r.URL.Query().Get("hdr_id")
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