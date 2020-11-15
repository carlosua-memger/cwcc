package assignments

import (
	
		"net/http"
		/*_"go-mssqldb-master"*/
		"fmt"
		"strings"
		"text/template"
		_"time"
		config "../../config"
		login "../../login"
		_"encoding/json"
		/*_"github.com/gorilla/sessions"*/
		"time"
		"strconv"
		_"encoding/json"	
		datatables "../../datatables"
		"encoding/json"
)		


func  AssignmentHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	last_start :=r.URL.Query().Get("last_start")
	last_length :=r.URL.Query().Get("last_length")
	last_search :=r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key,r.URL.Query().Get("rights"))
	tconf := make(map[string]string)
	ln :=  r.URL.Query().Get("length")
	tconf["Panel_name"] ="Assignment" 
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="false"  
	tconf["Add_new_button"] ="false"
	tconf["Add_new_button_url"] ="/assignment/CreateAssignment?parent=null&rights="+rights
	tconf["level"] ="1"
	tconf["pageLength"] =last_length 
	tconf["displayStart"] =last_start
	tconf["last_search"] = last_search
	if last_length==""{
		tconf["pageLength"] ="10" 
		tconf["displayStart"] ="0"
	}
	// rights for tables 	
	if( strings.Contains(dec_rights, "L1add") ){
		tconf["Add_new_button"] ="true"
	}
	//end rights for tables

	var header = []string{"Action", "Recno", "ID", "Name"}  //column header
	if ln=="" {
		var ajaxURLdata = "/assignment?rights="+rights
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("hris/assignments/table.html"); err != nil {
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
			sortcol="3"
		}
		srtcolindexInt,_ :=strconv.Atoi(sortcol)
	    var sortColName string
		for key , value := range header{
			if ( srtcolindexInt==key){
				sortColName = value
			}
		}
		_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 
		tconf["ActionCol"] ="true" // config for action column
	 	tconf["ActionCol_param"] ="ID" // config for parameter of action 
	 	tconf["ActionCol_edit"] ="false" // config for edit click 
	 	tconf["ActionCol_edit_url"] ="/assignment/EditAssignment?rights="+rights+"&ID=" // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/assignment/DeleteAssignment?ID="
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	tconf["ActionCol_add_child_url"] ="/assignment/CreateAssignment?rights="+rights+"&parent=" // config for addchild click 
	 	tconf["ActionCol_duplicate"] ="true"
	 	tconf["ActionCol_duplicate_is_modal"] ="false"
	 	tconf["ActionCol_duplicate_url"] ="/assignment/DuplicateAssignment?rights="+rights+"&ID=" // config for edit click 
	 	tconf["ActionCol_rate"] ="true"
	 	tconf["ActionCol_rate_is_modal"] ="false"
	 	tconf["ActionCol_rate_url"] ="/assignment/CreateRate?rights="+rights+"&ID=" // config for edit click 
	 	tconf["ActionCol_members"] ="false"
	 	tconf["ActionCol_members_url"] ="/assignment/AssignmentMembers?rights="+rights+"&ID="
	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
		tconf["dr"]=dr  //draw
		// rights for tables 	
		if( strings.Contains(dec_rights, "L1edit") ){
			tconf["ActionCol_edit"] ="true" // config for edit click 
		}
		if( strings.Contains(dec_rights, "L1delete") ){
			tconf["ActionCol_delete"] ="true" // config for delete click 
		}
		if( strings.Contains(dec_rights, "L1addChild") ){
			tconf["ActionCol_add_child"] ="true" // config for delete click 
		}
		if( strings.Contains(dec_rights, "L1members") ){
			tconf["ActionCol_members"] ="true" // config for member assign click 
		}
		//end rights for tables
		//geting total record etc 
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		fmt.Println("exec LBR_Assign_List 1, "+str+", "+ln+", "+str_OrgID+", "+session_user_id_str+", '"+sortColName+"', '', '"+sr+"'	")
		tconf["sql_total"] = "exec LBR_Assign_List 1, "+str+", "+ln+", "+str_OrgID+", "+session_user_id_str+", '"+sortColName+"', '', '"+sr+"'	"	
		tconf["sql_data"] ="exec LBR_Assign_List 0, "+str+", "+ln+", "+str_OrgID+", "+session_user_id_str+", '"+sortColName+"', '', '"+sr+"'	"
		datatables.DatatableretArray(w,tconf)  //datatables return array data
		  	  
		
	}	    



}

type Profile struct {
	  	Message string  // error status output or validated data
	  	Status bool //pass or fail
}
type RetMessage struct{
	RetNo int
	RetMsg string
}
func CreateAssignmentHAndler(w http.ResponseWriter, r *http.Request) {

	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	
	if r.Method =="GET" {	
		_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		
		parent :=r.URL.Query().Get("parent")
		// new struct for Org
		rows2 ,err, _,_ := config.Ap_sql("orgmem_get 1, 0,"+session_user_id_str,1)
	    if err != nil {
	    	panic(err.Error())
	    }
		type row2 struct {
	            Org_id int64
	            Org_name  string
	    }     
		got2 := []row2{}
		    for rows2.Next() {
		            var r row2
		            err = rows2.Scan(&r.Org_id, &r.Org_name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post2 := row2{
		            	r.Org_id,
				        r.Org_name,
				      
				    }
					got2 = append(got2, post2)
		    }

		 //district dropdown
		/*type District struct {
			   DisID int
	           District string
	    }*/
	    rows2 ,err, _,_ = config.Ap_sql("SysDropDown_Get 'AssignmentType'",1)
	    var name string
	    var Arrname []string
    	for rows2.Next() {
            err = rows2.Scan(&name)
            if err != nil {
                    panic(err.Error())
            }
			Arrname = append(Arrname, name)
		      
		   
		}

		type Data struct {
			Org []row2 
			Type []string
			Parent string 
			Arr_itemtag  [][]string 
			Rights string
		}    


		arr_data_itemtag := datatables.DataList(`LBR_Assign_Get 2, `+str_OrgID+`,`+session_user_id_str)


     	tmpl := template.New("CreateAssignment.html")
		//var err error
		if tmpl, err = tmpl.ParseFiles("hris/assignments/CreateAssignment.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{got2,Arrname,parent ,arr_data_itemtag,rights })
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}	
	}else{	
		r.ParseForm()
		session_username,_ := login.Get_account_info(r)
		if r.Form["postcode"][0]==""{
			r.Form["postcode"][0] = "1"
		}
		if strings.TrimSpace(r.Form["parent"][0])==""{
			r.Form["parent"][0] = "null"
		}
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		sqlstr :=`
			exec LBR_Assign_Save 
				 'Add' 
				, '`+session_username+`'
				, 0 
				, `+str_OrgID+` 
				, '`+r.Form["name"][0] +`'
				, '`+r.Form["type"][0] +`'
				, `+r.Form["parent"][0] +`
				, '`+r.Form["location"][0] + `'
				, `+r.Form["postcode"][0] + `
				, '`+r.Form["sched_start"][0] + `'
				, `+r.Form["req_hrs_day"][0] + `
				, `+r.Form["req_hrs_wk"][0] + `
				, `+r.Form["graceperiod"][0] + `
				, `+r.Form["reg_late_interval"][0] + `
				, `+r.Form["reg_break_interval"][0] + `
				, `+r.Form["break_hr"][0] + `
				, `+r.Form["OT_hrs_min"][0] + `
				, `+r.Form["OT_hrs_interval"][0] + `
				, `+r.Form["ND_hrs_min"][0] + `
				, `+r.Form["ND_hrs_interval"][0] + `
				, '`+r.Form["default_in"][0] + `'
				, '`+r.Form["default_out"][0] + `'
				,  `+r.Form["sun_rest"][0] + `
				, '`+r.Form["sun_in"][0] + `'
				, '`+r.Form["sun_out"][0] + `'
				,  `+r.Form["mon_rest"][0] + `
				, '`+r.Form["mon_in"][0] + `' 
				, '`+r.Form["mon_out"][0] + `'
				,  `+r.Form["tue_rest"][0] + `
				, '`+r.Form["tue_in"][0] + `'
				, '`+r.Form["tue_out"][0] + `'
				,  `+r.Form["wed_rest"][0] + `
				, '`+r.Form["wed_in"][0] + `'
				, '`+r.Form["wed_out"][0] + `'
				,  `+r.Form["thu_rest"][0] + `
				, '`+r.Form["thu_in"][0] + `'
				, '`+r.Form["thu_out"][0] + `'
				,  `+r.Form["fri_rest"][0] + `
				, '`+r.Form["fri_in"][0] + `'
				, '`+r.Form["fri_out"][0] + `'
				,  `+r.Form["sat_rest"][0] + `
				, '`+r.Form["sat_in"][0] + `'
				, '`+r.Form["sat_out"][0] + `'
				, `+r.Form["target_all"][0] + `
				, `+r.Form["target_per"][0] + `
				, '`+r.Form["remarks"][0] + `'
				, '`+r.Form["mrate"][0] + `'
				, '`+r.Form["drate"][0] + `'
				, '`+r.Form["hrate"][0] + `'
		`
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
	
	profile := Profile{rtmsg.RetMsg,true}
	js, err := json.Marshal(profile)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
		}
  	w.Header().Set("Content-Type", "application/json")
  	w.Write(js)

	}	
}
type Error_mess struct {
	RetNo	int
    RetMsg	string
	       	            
}
func DeleteAssignmentHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {
		tconf := make(map[string]string)
		tconf["RecID"] = r.URL.Query().Get("ID")
	   	tconf["delete_url"] = "/assignment/DeleteAssignment"
	   	
	   	tmpl := template.New("modal_delete.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("hris/assignments/modal_delete.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,tconf)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		r.ParseForm()
		session_username,_ := login.Get_account_info(r)	
		fmt.Println("Member_Save 'Delete','"+session_username+"' ,"+r.Form["rec_no"][0]+"")
		sqlstr:= `LBR_Assign_Save 'Delete','`+session_username+`',`+r.Form["rec_no"][0]
  		rows2 ,err, _,_ := config.Ap_sql(sqlstr,1)
  		if err != nil {
	    	panic(err.Error())
	    }
		var errorMsg  Error_mess
	    for rows2.Next() {
            var r Error_mess
            err = rows2.Scan(&r.RetNo, &r.RetMsg)
            if err != nil {
                    panic(err.Error())
            }
			errorMsg = Error_mess{r.RetNo,r.RetMsg}
		      
		}
				
		
	    
		js, err := json.Marshal(errorMsg)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	fmt.Println("Data Successfully added")
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)

	}		


}
type LBR_Assign_Get struct {
	Id 			int
	Org 		int
	Org_name 	string
	Name 		interface{}
	Type 		string
	Parent 		interface{}
	Parent_name string
	Location string
	Postcode int
	Postcode_name 		string
	Sched_start 		time.Time
	Req_hrs_day string
	Req_hrs_wk 			string
	Graceperiod 		string
	Reg_late_interval 	string
	Reg_break_interval 	string 
	Break_hr 			string
	OT_hrs_min			string
	OT_hrs_interval 	string
	ND_hrs_min 			string 
	ND_hrs_interval 	string
	Default_in 			interface{}
	Default_out 		interface{}
	Sun_in 				interface{}
	Sun_out 			interface{}
	Mon_in 				interface{}
	Mon_out 			interface{}
	Tue_in 				interface{}
	Tue_out 			interface{}
	Wed_in 				interface{}
	Wed_out 			interface{}
	Thu_in 				interface{}
	Thu_out 			interface{}
	Fri_in 				interface{}
	Fri_out 			interface{}
	Sat_in 				interface{}
	Sat_out 			interface{}
	
	Sun_rest			interface{}
	Mon_rest			interface{}
	Tue_rest			interface{}
	Wed_rest			interface{}
	Thu_rest			interface{}
	Fri_rest			interface{}
	Sat_rest			interface{}


	Target_all 			interface{}
	Target_per 			interface{}
	Remarks 			interface{}
	RateperMonth 		interface{}
	RateperDay 			interface{}
	RateperHour 		interface{}
}

func nil_to_string(interfacedata interface{} ) interface{} {
	var v2 interface{}
	if interfacedata==nil{
		v2 = ``
	}else{
		//v2=interfacedata.Format("15:04:05")
		 switch v := (interfacedata).(type) {
			 case nil:
			v2 =fmt.Sprint("NULL")
		case bool:
			if v {
			v2 =	fmt.Sprint("1")
			} else {
			v2 =	fmt.Sprint("0")
			}
		case []byte:
			v2 = fmt.Sprint(string(v))
		case time.Time:
			v2 = fmt.Sprint(v.Format("15:04:05"))

		default:
			v2= fmt.Sprint(v)
		}	
	}
	return	v2
}
type LBR_Assign struct {
	Id interface{}
	Name string
}
func LBR_Assign_Get_drowdown_raw(usernameID string) []LBR_Assign {
	

	LBR_Assign_Get_drowdown_raw ,err, _,_ := config.Ap_sql(`exec LBR_Assign_Get 2, 2, `+usernameID,1)
	
 	if err != nil {
    	panic(err.Error())
	}
	var LBR_Assign_data []LBR_Assign
	for LBR_Assign_Get_drowdown_raw.Next() {
		var	r LBR_Assign
		err = LBR_Assign_Get_drowdown_raw.Scan(&r.Id, &r.Name)
            if err != nil {
                    panic(err.Error())
            }
        data := LBR_Assign{r.Id,r.Name }
        LBR_Assign_data =  append(LBR_Assign_data,data)    
    }        
    return LBR_Assign_data

}

 

func EditAssignmentHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {	
		rights :=r.URL.Query().Get("rights")
		tconf := make(map[string]string)
		parent :=r.URL.Query().Get("parent")
		tconf["RecID"] = r.URL.Query().Get("ID")
		_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 
		
		// LBR_Assign_Get --
		LBR_Assign_Get_raw ,err, _,_ := config.Ap_sql(`LBR_Assign_Get 1, ` + tconf["RecID"] ,1)
	    if err != nil {
	    	panic(err.Error())
	    }
		    
		//var LBR_Assign_Get_data LBR_Assign_Get
	    var LBR_Assign_Get_data LBR_Assign_Get
	    for LBR_Assign_Get_raw.Next() {
	            var r LBR_Assign_Get
	            err = LBR_Assign_Get_raw.Scan(&r.Id, &r.Org, &r.Org_name, &r.Name, &r.Type, &r.Parent, &r.Parent_name, &r.Location, &r.Postcode, &r.Postcode_name, &r.Sched_start, &r.Req_hrs_day, &r.Req_hrs_wk, &r.Graceperiod, &r.Reg_late_interval, &r.Reg_break_interval, &r.Break_hr, &r.OT_hrs_min, &r.OT_hrs_interval, &r.ND_hrs_min, &r.ND_hrs_interval    , &r.Default_in, &r.Default_out, &r.Sun_in, &r.Sun_out, &r.Mon_in, &r.Mon_out, &r.Tue_in ,&r.Tue_out     , &r.Wed_in, &r.Wed_out, &r.Thu_in, &r.Thu_out, &r.Fri_in, &r.Fri_out, &r.Sat_in, &r.Sat_out , &r.Sun_rest,&r.Mon_rest,&r.Tue_rest,&r.Wed_rest,&r.Thu_rest,&r.Fri_rest,&r.Sat_rest,  &r.Target_all,&r.Target_per, &r.Remarks, &r.RateperMonth, &r.RateperDay, &r.RateperHour)
	            if err != nil {
	                    panic(err.Error())
	            }
	             //fmt.Println("name", r.Parent.IsNil())
	             fmt.Println(r.Sched_start.Format(" Jan _2 2006") )
	             LBR_Assign_Get_data = LBR_Assign_Get{	r.Id,
				 										r.Org, 
				 										r.Org_name,
				 										r.Name,
				 										r.Type, 
				 										nil_to_string(r.Parent), 
				 										r.Parent_name, 
				 										r.Location, 
				 										r.Postcode, 
				 										r.Postcode_name, 
				 										r.Sched_start, 
				 										r.Req_hrs_day, 
				 										r.Req_hrs_wk, 
				 										r.Graceperiod, 
				 										r.Reg_late_interval, 
				 										r.Reg_break_interval, 
				 										r.Break_hr, 
				 										r.OT_hrs_min, 
				 										r.OT_hrs_interval, 
				 										r.ND_hrs_min, 
				 										r.ND_hrs_interval, 
				 										nil_to_string(r.Default_in), 
				 										nil_to_string(r.Default_out), 
				 										nil_to_string(r.Sun_in), 
				 										nil_to_string(r.Sun_out), 
				 										nil_to_string(r.Mon_in), 
				 										nil_to_string(r.Mon_out), 
				 										nil_to_string(r.Tue_in) ,
				 										nil_to_string(r.Tue_out)     , 
				 										nil_to_string(r.Wed_in), 
				 										nil_to_string(r.Wed_out), 
				 										nil_to_string(r.Thu_in), 
				 										nil_to_string(r.Thu_out), 
				 										nil_to_string(r.Fri_in), 
				 										nil_to_string(r.Fri_out), 
				 										nil_to_string(r.Sat_in), 
				 										nil_to_string(r.Sat_out) ,

				 										nil_to_string(r.Sun_rest) ,
	 													nil_to_string(r.Mon_rest) ,
	 													nil_to_string(r.Tue_rest) ,
	 													nil_to_string(r.Wed_rest) ,
	 													nil_to_string(r.Thu_rest) ,
	 													nil_to_string(r.Fri_rest) ,
	 													nil_to_string(r.Sat_rest) ,

				 										nil_to_string(r.Target_all),
				 										nil_to_string(r.Target_per), 
				 										nil_to_string(r.Remarks), 
				 										nil_to_string(r.RateperMonth), 
				 										nil_to_string(r.RateperDay), 
				 										nil_to_string(r.RateperHour) }
				/*post2 := row2{
	            	r.Org_id,
			        r.Org_name,
			      
			    }
				got2 = append(got2, post2)*/
	    }

		// new struct for Org
		rows2 ,err, _,_ := config.Ap_sql("orgmem_get 1, 0,"+session_user_id_str,1)
	    if err != nil {
	    	panic(err.Error())
	    }
		type row2 struct {
	            Org_id int64
	            Org_name  string
	    }     
		got2 := []row2{}
		    for rows2.Next() {
		            var r row2
		            err = rows2.Scan(&r.Org_id, &r.Org_name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post2 := row2{
		            	r.Org_id,
				        r.Org_name,
				      
				    }
					got2 = append(got2, post2)
		    }

		 //district dropdown
		/*type District struct {
			   DisID int
	           District string
	    }*/
	    rows2 ,err, _,_ = config.Ap_sql("SysDropDown_Get 'AssignmentType'",1)
	    var name string
	    var Arrname []string
    	for rows2.Next() {
            err = rows2.Scan(&name)
            if err != nil {
                    panic(err.Error())
            }
			Arrname = append(Arrname, name)
		      
		   
		}

	    type Data struct {
			Org []row2 
			Type []string
			Parent string 
			LBR_Assign_Get LBR_Assign_Get 
			LBR_Assign_Get_drowdown_raw []LBR_Assign
			Rights string
		}
		htdata:= &Data{got2,Arrname,parent,LBR_Assign_Get_data,LBR_Assign_Get_drowdown_raw(session_user_id_str),rights}

		funcMap := template.FuncMap{
		    "date_formater": func(a time.Time) string {
		        return a.Format(" Jan _2 2006")
		    },
		    "time_formater": func(a time.Time) string {
		        return a.Format("15:04:05")
		    },

		}   
     	tmpl := template.New("assignementEdit.html").Funcs(funcMap)
		//var err error
		if tmpl, err = tmpl.ParseFiles("hris/assignments/assignementEdit.html"); err != nil {
			fmt.Println(err)
		}

		err1 := tmpl.Execute(w,htdata)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}	
	}else{	 //on edit
		r.ParseForm()
		session_username,_ := login.Get_account_info(r)	
		if r.Form["postcode"][0]==""{
			r.Form["postcode"][0] = "1"
		}
		if strings.TrimSpace(r.Form["parent"][0])==""{
			r.Form["parent"][0] = "null"
		} 
		sqlstr :=`
			exec LBR_Assign_Save 
				 'Edit' 
				, '`+session_username+`'
				, `+r.Form["Id"][0] +`
				, 2 
				, '`+r.Form["name"][0] +`'
				, '`+r.Form["type"][0] +`'
				, `+r.Form["parent"][0] +`
				, '`+r.Form["location"][0] + `'
				, `+r.Form["postcode"][0] + `
				, '`+r.Form["sched_start"][0] + `'
				, `+r.Form["req_hrs_day"][0] + `
				, `+r.Form["req_hrs_wk"][0] + `
				, `+r.Form["graceperiod"][0] + `
				, `+r.Form["reg_late_interval"][0] + `
				, `+r.Form["reg_break_interval"][0] + `
				, `+r.Form["break_hr"][0] + `
				, `+r.Form["OT_hrs_min"][0] + `
				, `+r.Form["OT_hrs_interval"][0] + `
				, `+r.Form["ND_hrs_min"][0] + `
				, `+r.Form["ND_hrs_interval"][0] + `
				, '`+r.Form["default_in"][0] + `'
				, '`+r.Form["default_out"][0] + `'
				, `+r.Form["sun_rest"][0] + `
				, '`+r.Form["sun_in"][0] + `'
				, '`+r.Form["sun_out"][0] + `'
				, `+r.Form["mon_rest"][0] + `
				, '`+r.Form["mon_in"][0] + `'
				, '`+r.Form["mon_out"][0] + `'
				,  `+r.Form["tue_rest"][0] + `
				, '`+r.Form["tue_in"][0] + `'
				, '`+r.Form["tue_out"][0] + `'
				,  `+r.Form["wed_rest"][0] + `
				, '`+r.Form["wed_in"][0] + `'
				, '`+r.Form["wed_out"][0] + `'
				,  `+r.Form["thu_rest"][0] + `
				, '`+r.Form["thu_in"][0] + `'
				, '`+r.Form["thu_out"][0] + `'
				,  `+r.Form["fri_rest"][0] + `
				, '`+r.Form["fri_in"][0] + `'
				, '`+r.Form["fri_out"][0] + `'
				,  `+r.Form["sat_rest"][0] + `
				, '`+r.Form["sat_in"][0] + `'
				, '`+r.Form["sat_out"][0] + `'
				, `+r.Form["target_all"][0] + `
				, `+r.Form["target_per"][0] + `
				, '`+r.Form["remarks"][0] + `'
				, '`+r.Form["mrate"][0] + `'
				, '`+r.Form["drate"][0] + `'
				, '`+r.Form["hrate"][0] + `'
		`
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
	
	profile := Profile{rtmsg.RetMsg,true}
	js, err := json.Marshal(profile)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
		}
  	w.Header().Set("Content-Type", "application/json")
  	w.Write(js)

	}		

	
}


func AssignmentMembersHAndler(w http.ResponseWriter, r *http.Request){

	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	dec_rights := config.Decrypt(config.Encrypt_key,r.URL.Query().Get("rights"))
	AssignmentID :=r.URL.Query().Get("ID") 

	tconf := make(map[string]string)
	ln :=  r.URL.Query().Get("length")
	tconf["Rights"] = rights
	tconf["Panel_name"] =`Member/s on this Assigment [ `+ AssignmentID +` ]`   //panel name
	tconf["Add_back_button"] ="true"  
	tconf["Add_new_button"] ="true"  
	tconf["level"] ="2" 
	tconf["pageLength"] ="10" 
	tconf["displayStart"] ="0"
	tconf["last_search"]=""
	tconf["Add_new_button_url"] ="/assignment/AssignmentMembersCreate?rights="+rights+"&assignmentid="+AssignmentID
	var header = []string{"Action", "Recno", "ID","Employee #", "Leader", "Lbr_assign", "Lbr_assign_name","orgmem", "Member_name", "Processaccess", "Txdate", "Device_id", "Remarks", "Enable", "Disable"}  //column header
	// rights for tables 	
	if( strings.Contains(dec_rights, "AMadd") ){
		tconf["Add_new_button"] ="true"
	}
	//end rights for tables

	if ln=="" {
		var ajaxURLdata = "/assignment/AssignmentMembers?rights="+rights +"&ID="+AssignmentID
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("hris/assignments/table.html"); err != nil {
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
		tconf["ActionCol_assignment"] ="false" 
		tconf["ActionCol"] ="true" // config for action column
	 	tconf["ActionCol_param"] ="ID" // config for parameter of action 
	 	tconf["ActionCol_edit"] ="false" // config for edit click 
	 	tconf["ActionCol_edit_is_modal"] ="false" // config for edit click 
	 	tconf["ActionCol_edit_url"] ="/assignment/AssignmentMembersEdit?rights="+rights+"&ID=" // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/assignment/AssignmentMembersDelete?ID="
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	//tconf["ActionCol_assignment"] ="false"
	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
		tconf["dr"]=dr  //draw
		//geting total record etc 
		
		// rights for tables 	
		if( strings.Contains(dec_rights, "AMedit") ){
			tconf["ActionCol_edit"] ="true" // config for edit click 
		}
		if( strings.Contains(dec_rights, "AMdelete") ){
			tconf["ActionCol_delete"] ="true" // config for delete click 
		}
		//end rights for tables
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		fmt.Println("exec LBR_Memass_list 0, "+str_OrgID+", 0, "+AssignmentID+", "+str+", "+ln+", '"+sortColName+"', '', '"+sr+"'")
		tconf["sql_total"] = "exec LBR_Memass_list 1, "+str_OrgID+", 0, "+AssignmentID+", "+str+", "+ln+", '"+sortColName+"', '', '"+sr+"'"	
		tconf["sql_data"] =  "exec LBR_Memass_list 0, "+str_OrgID+", 0, "+AssignmentID+", "+str+", "+ln+", '"+sortColName+"', '', '"+sr+"'"

		datatables.DatatableretArray(w,tconf)
		  	  
		
	}	    
}

func AssignmentMembersCreateHAndler(w http.ResponseWriter, r *http.Request){
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	tconf := make(map[string]string)
	tconf["rights"]=rights
	tconf["assignmentid"]=r.URL.Query().Get("assignmentid")
	if r.Method =="GET" {
		    
		type Data struct {
			Htdata map[string]string
			
		}   
		htdata := Data{tconf} 
     	tmpl := template.New("assignmentmemberCreate.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("hris/assignments/assignmentmemberCreate.html"); err != nil {
			fmt.Println(err)
		}

		err1 := tmpl.Execute(w,htdata)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else{
		//do create assigment Hire.
		r.ParseForm()

		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		session_username,_ := login.Get_account_info(r)
		
		arr_members :=strings.Split(r.Form["member"][0] ,",")
		
		var retr RetMessage
			
		for _,members := range arr_members {
			fmt.Println(members)
			sqlstr:= `exec LBR_Memass_Save 
						'Add',
						'`+session_username + `' ,
						0,
						`+str_OrgID+`,
						` + r.Form["lbr_assign"][0] + ` ,
						` + members +`	 ,
						` + r.Form["processaccess"][0]+` ,
						'`+ r.Form["txdate"][0]+`' ,
						'`+ r.Form["device_id"][0]+`' ,
						'`+ r.Form["remarks"][0]+`' ,
						'`+ r.Form["enable"][0]+`' ,
						'`+ r.Form["disable"][0] + `',
						'`+ r.Form["leader"][0] + `' ,
						'`+ r.Form["emp_id"][0] + `' `

			rowret ,err, _,_ := config.Ap_sql(sqlstr,1)
			if err != nil {
			panic(err.Error())
			}
			for rowret.Next() {
			    err = rowret.Scan(&retr.RetNo,&retr.RetMsg)
			    if err != nil {
			        panic(err.Error())
			    }
			 	retr = RetMessage{retr.RetNo,retr.RetMsg}
			}
		}
		
		

		profile := Profile{retr.RetMsg,true}
		js, err := json.Marshal(profile)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)


	}	

}

func AssignmentMembersEditHAndler(w http.ResponseWriter, r *http.Request){
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	tconf := make(map[string]string)
	tconf["rights"]=rights
	tconf["ID"]=r.URL.Query().Get("ID")
	if r.Method =="GET" {
		    
		
		type LBR_MemAss struct {
			Id 				interface{}	
			Lbr_assign		interface{} 
			Lbr_assign_name interface{}
			Orgmem 			interface{}
			Member_name		interface{}
			Leader	 		interface{}
			HRIS			interface{}
			Processaccess	interface{}
			Txdate			interface{}
			Device_id		interface{}
			Remarks			interface{}
			Enable			interface{}
			Disable 		interface{}

		}
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		LBR_MemAss_raw ,err, _,_ := config.Ap_sql(`exec LBR_MemAss_Get 0, `+tconf["ID"]+`,`+str_OrgID+`` ,1)
	 	if err != nil {
	    	panic(err.Error())
		}

		var LBR_MemAss_data  LBR_MemAss
		for LBR_MemAss_raw.Next() {
			err = LBR_MemAss_raw.Scan(&LBR_MemAss_data.Id, &LBR_MemAss_data.Lbr_assign,&LBR_MemAss_data.Lbr_assign_name,&LBR_MemAss_data.Orgmem ,&LBR_MemAss_data.Member_name,&LBR_MemAss_data.Leader,&LBR_MemAss_data.HRIS,&LBR_MemAss_data.Processaccess,&LBR_MemAss_data.Txdate ,&LBR_MemAss_data.Device_id ,&LBR_MemAss_data.Remarks,&LBR_MemAss_data.Enable,&LBR_MemAss_data.Disable)
	            if err != nil {
	                    panic(err.Error())
	            }
	        //LBR_MemAss_data =   LBR_MemAss{}   
	    }

	    type Data struct {
			Htdata map[string]string
			LBR_MemAss LBR_MemAss
		}        
		htdata := Data{tconf,LBR_MemAss_data} 
     	
     	funcMap := template.FuncMap{
		    "date_formater": func(a time.Time) string {
		        return a.Format(" Jan _2 2006")
		    },
		} 
		tmpl := template.New("assignmentmemberEdit.html").Funcs(funcMap)
		//var err error
		if tmpl, err = tmpl.ParseFiles("hris/assignments/assignmentmemberEdit.html"); err != nil {
			fmt.Println(err)
		}

		err1 := tmpl.Execute(w,htdata)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else{
		//do create assigment Hire.
		session_username,_ := login.Get_account_info(r)
		r.ParseForm()

		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		arr_members :=strings.Split(r.Form["member"][0] ,",")
		
		var retr RetMessage
			
		for _,members := range arr_members {
			fmt.Println(members)
			sqlstr:= `exec LBR_Memass_Save 
						'Edit',
						'`+session_username + `' ,
						`+ r.Form["assinmentID"][0] +  ` ,
						`+str_OrgID+`,
						` + r.Form["lbr_assign"][0] + ` ,
						` + members +`	 ,
						` + r.Form["processaccess"][0]+` ,
						'`+ r.Form["txdate"][0]+`' ,
						'`+ r.Form["device_id"][0]+`' ,
						'`+ r.Form["remarks"][0]+`' ,
						'`+ r.Form["enable"][0]+`' ,
						'`+ r.Form["disable"][0] + `' ,
						'`+ r.Form["leader"][0] + `',
						'`+ r.Form["emp_id"][0] + `' `

			fmt.Println(sqlstr)
			rowret ,err, _,_ := config.Ap_sql(sqlstr,1)
			if err != nil {
			panic(err.Error())
			}
			for rowret.Next() {
			    err = rowret.Scan(&retr.RetNo,&retr.RetMsg)
			    if err != nil {
			        panic(err.Error())
			    }
			 	retr = RetMessage{retr.RetNo,retr.RetMsg}
			}
		}
		
		

		profile := Profile{retr.RetMsg,true}
		js, err := json.Marshal(profile)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)


	}	

}


func AssignmentMembersDeleteHAndler(w http.ResponseWriter, r *http.Request){
	login.Session_renew(w,r)
	if r.Method =="GET" {
		tconf := make(map[string]string)
		tconf["RecID"] = r.URL.Query().Get("ID")
	   	tconf["delete_url"] = "/assignment/AssignmentMembersDelete"
	   	

	   

	   	tmpl := template.New("modal_delete.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("hris/assignments/modal_delete.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,tconf)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		session_username,_ := login.Get_account_info(r)
		r.ParseForm()	
		sqlstr:=`	exec LBR_Memass_Save
		                    'Delete',
		                    '`+session_username+`' ,
		                    `+r.Form["rec_no"][0]+`
				`

  		rows2 ,err, _,_ := config.Ap_sql(sqlstr,1)
  		if err != nil {
	    	panic(err.Error())
	    }
		var errorMsg  Error_mess
	    for rows2.Next() {
            var r Error_mess
            err = rows2.Scan(&r.RetNo, &r.RetMsg)
            if err != nil {
                    panic(err.Error())
            }
			errorMsg = Error_mess{r.RetNo,r.RetMsg}
		      
		}
				
		
	    
		js, err := json.Marshal(errorMsg)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	fmt.Println("Data Successfully added")
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)

	}		



}

func DuplicateAssignmentHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {	
		rights :=r.URL.Query().Get("rights")
		tconf := make(map[string]string)
		parent :=r.URL.Query().Get("parent")
		tconf["RecID"] = r.URL.Query().Get("ID")
		_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 
		
		// LBR_Assign_Get --
		LBR_Assign_Get_raw ,err, _,_ := config.Ap_sql(`LBR_Assign_Get 1, ` + tconf["RecID"] ,1)
	    if err != nil {
	    	panic(err.Error())
	    }
		    
		//var LBR_Assign_Get_data LBR_Assign_Get
	    var LBR_Assign_Get_data LBR_Assign_Get
	    for LBR_Assign_Get_raw.Next() {
	            var r LBR_Assign_Get
	            err = LBR_Assign_Get_raw.Scan(&r.Id, &r.Org, &r.Org_name, &r.Name, &r.Type, &r.Parent, &r.Parent_name, &r.Location, &r.Postcode, &r.Postcode_name, &r.Sched_start, &r.Req_hrs_day, &r.Req_hrs_wk, &r.Graceperiod, &r.Reg_late_interval, &r.Reg_break_interval, &r.Break_hr, &r.OT_hrs_min, &r.OT_hrs_interval, &r.ND_hrs_min, &r.ND_hrs_interval    , &r.Default_in, &r.Default_out, &r.Sun_in, &r.Sun_out, &r.Mon_in, &r.Mon_out, &r.Tue_in ,&r.Tue_out     , &r.Wed_in, &r.Wed_out, &r.Thu_in, &r.Thu_out, &r.Fri_in, &r.Fri_out, &r.Sat_in, &r.Sat_out , &r.Sun_rest,&r.Mon_rest,&r.Tue_rest,&r.Wed_rest,&r.Thu_rest,&r.Fri_rest,&r.Sat_rest,  &r.Target_all,&r.Target_per, &r.Remarks, &r.RateperMonth, &r.RateperDay, &r.RateperHour)
	            if err != nil {
	                    panic(err.Error())
	            }
	             //fmt.Println("name", r.Parent.IsNil())
	             fmt.Println(r.Sched_start.Format(" Jan _2 2006") )
	             LBR_Assign_Get_data = LBR_Assign_Get{	r.Id,
				 										r.Org, 
				 										r.Org_name,
				 										r.Name,
				 										r.Type, 
				 										nil_to_string(r.Parent), 
				 										r.Parent_name, 
				 										r.Location, 
				 										r.Postcode, 
				 										r.Postcode_name, 
				 										r.Sched_start, 
				 										r.Req_hrs_day, 
				 										r.Req_hrs_wk, 
				 										r.Graceperiod, 
				 										r.Reg_late_interval, 
				 										r.Reg_break_interval, 
				 										r.Break_hr, 
				 										r.OT_hrs_min, 
				 										r.OT_hrs_interval, 
				 										r.ND_hrs_min, 
				 										r.ND_hrs_interval, 
				 										nil_to_string(r.Default_in), 
				 										nil_to_string(r.Default_out), 
				 										nil_to_string(r.Sun_in), 
				 										nil_to_string(r.Sun_out), 
				 										nil_to_string(r.Mon_in), 
				 										nil_to_string(r.Mon_out), 
				 										nil_to_string(r.Tue_in) ,
				 										nil_to_string(r.Tue_out)     , 
				 										nil_to_string(r.Wed_in), 
				 										nil_to_string(r.Wed_out), 
				 										nil_to_string(r.Thu_in), 
				 										nil_to_string(r.Thu_out), 
				 										nil_to_string(r.Fri_in), 
				 										nil_to_string(r.Fri_out), 
				 										nil_to_string(r.Sat_in), 
				 										nil_to_string(r.Sat_out) ,

				 										nil_to_string(r.Sun_rest) ,
	 													nil_to_string(r.Mon_rest) ,
	 													nil_to_string(r.Tue_rest) ,
	 													nil_to_string(r.Wed_rest) ,
	 													nil_to_string(r.Thu_rest) ,
	 													nil_to_string(r.Fri_rest) ,
	 													nil_to_string(r.Sat_rest) ,

				 										nil_to_string(r.Target_all),
				 										nil_to_string(r.Target_per), 
				 										nil_to_string(r.Remarks), 
				 										nil_to_string(r.RateperMonth), 
				 										nil_to_string(r.RateperDay), 
				 										nil_to_string(r.RateperHour) }
				/*post2 := row2{
	            	r.Org_id,
			        r.Org_name,
			      
			    }
				got2 = append(got2, post2)*/
	    }

		// new struct for Org
		rows2 ,err, _,_ := config.Ap_sql("orgmem_get 1, 0,"+session_user_id_str,1)
	    if err != nil {
	    	panic(err.Error())
	    }
		type row2 struct {
	            Org_id int64
	            Org_name  string
	    }     
		got2 := []row2{}
		    for rows2.Next() {
		            var r row2
		            err = rows2.Scan(&r.Org_id, &r.Org_name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post2 := row2{
		            	r.Org_id,
				        r.Org_name,
				      
				    }
					got2 = append(got2, post2)
		    }

		 //district dropdown
		/*type District struct {
			   DisID int
	           District string
	    }*/
	    rows2 ,err, _,_ = config.Ap_sql("SysDropDown_Get 'AssignmentType'",1)
	    var name string
	    var Arrname []string
    	for rows2.Next() {
            err = rows2.Scan(&name)
            if err != nil {
                    panic(err.Error())
            }
			Arrname = append(Arrname, name)
		      
		   
		}

	    type Data struct {
			Org []row2 
			Type []string
			Parent string 
			LBR_Assign_Get LBR_Assign_Get 
			LBR_Assign_Get_drowdown_raw []LBR_Assign
			Rights string
		}
		htdata:= &Data{got2,Arrname,parent,LBR_Assign_Get_data,LBR_Assign_Get_drowdown_raw(session_user_id_str),rights}

		funcMap := template.FuncMap{
		    "date_formater": func(a time.Time) string {
		        return a.Format(" Jan _2 2006")
		    },
		    "time_formater": func(a time.Time) string {
		        return a.Format("15:04:05")
		    },

		}   
     	tmpl := template.New("assignementDuplicate.html").Funcs(funcMap)
		//var err error
		if tmpl, err = tmpl.ParseFiles("hris/assignments/assignementDuplicate.html"); err != nil {
			fmt.Println(err)
		}

		err1 := tmpl.Execute(w,htdata)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}	
	}else{	 //on edit
		r.ParseForm()
		session_username,_ := login.Get_account_info(r)	
		if r.Form["postcode"][0]==""{
			r.Form["postcode"][0] = "1"
		}
		if strings.TrimSpace(r.Form["parent"][0])==""{
			r.Form["parent"][0] = "null"
		}
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		sqlstr :=`
			exec LBR_Assign_Save 
				 'Add' 
				, '`+session_username+`'
				, 0 
				, `+str_OrgID+` 

				, '`+r.Form["name"][0] +`'
				, '`+r.Form["type"][0] +`'
				, `+r.Form["parent"][0] +`
				, '`+r.Form["location"][0] + `'
				, `+r.Form["postcode"][0] + `
				, '`+r.Form["sched_start"][0] + `'
				, `+r.Form["req_hrs_day"][0] + `
				, `+r.Form["req_hrs_wk"][0] + `
				, `+r.Form["graceperiod"][0] + `
				, `+r.Form["reg_late_interval"][0] + `
				, `+r.Form["reg_break_interval"][0] + `
				, `+r.Form["break_hr"][0] + `
				, `+r.Form["OT_hrs_min"][0] + `
				, `+r.Form["OT_hrs_interval"][0] + `
				, `+r.Form["ND_hrs_min"][0] + `
				, `+r.Form["ND_hrs_interval"][0] + `
				, '`+r.Form["default_in"][0] + `'
				, '`+r.Form["default_out"][0] + `'
				,  `+r.Form["sun_rest"][0] + `
				, '`+r.Form["sun_in"][0] + `'
				, '`+r.Form["sun_out"][0] + `'
				,  `+r.Form["mon_rest"][0] + `
				, '`+r.Form["mon_in"][0] + `' 
				, '`+r.Form["mon_out"][0] + `'
				,  `+r.Form["tue_rest"][0] + `
				, '`+r.Form["tue_in"][0] + `'
				, '`+r.Form["tue_out"][0] + `'
				,  `+r.Form["wed_rest"][0] + `
				, '`+r.Form["wed_in"][0] + `'
				, '`+r.Form["wed_out"][0] + `'
				,  `+r.Form["thu_rest"][0] + `
				, '`+r.Form["thu_in"][0] + `'
				, '`+r.Form["thu_out"][0] + `'
				,  `+r.Form["fri_rest"][0] + `
				, '`+r.Form["fri_in"][0] + `'
				, '`+r.Form["fri_out"][0] + `'
				,  `+r.Form["sat_rest"][0] + `
				, '`+r.Form["sat_in"][0] + `'
				, '`+r.Form["sat_out"][0] + `'
				, `+r.Form["target_all"][0] + `
				, `+r.Form["target_per"][0] + `
				, '`+r.Form["remarks"][0] + `'
				, '`+r.Form["mrate"][0] + `'
				, '`+r.Form["drate"][0] + `'
				, '`+r.Form["hrate"][0] + `'
		`
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
	
	profile := Profile{rtmsg.RetMsg,true}
	js, err := json.Marshal(profile)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
		}
  	w.Header().Set("Content-Type", "application/json")
  	w.Write(js)

	}		

	
}

func RateAssignmentHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	if r.Method =="GET" {	
		rights :=r.URL.Query().Get("rights")
		tconf := make(map[string]string)
		parent :=r.URL.Query().Get("parent")
		tconf["RecID"] = r.URL.Query().Get("ID")
		_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 
		
		// LBR_Assign_Get --
		LBR_Assign_Get_raw ,err, _,_ := config.Ap_sql(`LBR_Assign_Get 1, ` + tconf["RecID"] ,1)
	    if err != nil {
	    	panic(err.Error())
	    }
		    
		//var LBR_Assign_Get_data LBR_Assign_Get
	    var LBR_Assign_Get_data LBR_Assign_Get
	    for LBR_Assign_Get_raw.Next() {
	            var r LBR_Assign_Get
	            err = LBR_Assign_Get_raw.Scan(&r.Id, &r.Org, &r.Org_name, &r.Name, &r.Type, &r.Parent, &r.Parent_name, &r.Location, &r.Postcode, &r.Postcode_name, &r.Sched_start, &r.Req_hrs_day, &r.Req_hrs_wk, &r.Graceperiod, &r.Reg_late_interval, &r.Reg_break_interval, &r.Break_hr, &r.OT_hrs_min, &r.OT_hrs_interval, &r.ND_hrs_min, &r.ND_hrs_interval    , &r.Default_in, &r.Default_out, &r.Sun_in, &r.Sun_out, &r.Mon_in, &r.Mon_out, &r.Tue_in ,&r.Tue_out     , &r.Wed_in, &r.Wed_out, &r.Thu_in, &r.Thu_out, &r.Fri_in, &r.Fri_out, &r.Sat_in, &r.Sat_out , &r.Sun_rest,&r.Mon_rest,&r.Tue_rest,&r.Wed_rest,&r.Thu_rest,&r.Fri_rest,&r.Sat_rest,  &r.Target_all,&r.Target_per, &r.Remarks, &r.RateperMonth, &r.RateperDay, &r.RateperHour)
	            if err != nil {
	                    panic(err.Error())
	            }
	             //fmt.Println("name", r.Parent.IsNil())
	             fmt.Println(r.Sched_start.Format(" Jan _2 2006") )
	             LBR_Assign_Get_data = LBR_Assign_Get{	r.Id,
				 										r.Org, 
				 										r.Org_name,
				 										r.Name,
				 										r.Type, 
				 										nil_to_string(r.Parent), 
				 										r.Parent_name, 
				 										r.Location, 
				 										r.Postcode, 
				 										r.Postcode_name, 
				 										r.Sched_start, 
				 										r.Req_hrs_day, 
				 										r.Req_hrs_wk, 
				 										r.Graceperiod, 
				 										r.Reg_late_interval, 
				 										r.Reg_break_interval, 
				 										r.Break_hr, 
				 										r.OT_hrs_min, 
				 										r.OT_hrs_interval, 
				 										r.ND_hrs_min, 
				 										r.ND_hrs_interval, 
				 										nil_to_string(r.Default_in), 
				 										nil_to_string(r.Default_out), 
				 										nil_to_string(r.Sun_in), 
				 										nil_to_string(r.Sun_out), 
				 										nil_to_string(r.Mon_in), 
				 										nil_to_string(r.Mon_out), 
				 										nil_to_string(r.Tue_in) ,
				 										nil_to_string(r.Tue_out)     , 
				 										nil_to_string(r.Wed_in), 
				 										nil_to_string(r.Wed_out), 
				 										nil_to_string(r.Thu_in), 
				 										nil_to_string(r.Thu_out), 
				 										nil_to_string(r.Fri_in), 
				 										nil_to_string(r.Fri_out), 
				 										nil_to_string(r.Sat_in), 
				 										nil_to_string(r.Sat_out) ,

				 										nil_to_string(r.Sun_rest) ,
	 													nil_to_string(r.Mon_rest) ,
	 													nil_to_string(r.Tue_rest) ,
	 													nil_to_string(r.Wed_rest) ,
	 													nil_to_string(r.Thu_rest) ,
	 													nil_to_string(r.Fri_rest) ,
	 													nil_to_string(r.Sat_rest) ,

				 										nil_to_string(r.Target_all),
				 										nil_to_string(r.Target_per), 
				 										nil_to_string(r.Remarks), 
				 										nil_to_string(r.RateperMonth), 
				 										nil_to_string(r.RateperDay), 
				 										nil_to_string(r.RateperHour) }
				/*post2 := row2{
	            	r.Org_id,
			        r.Org_name,
			      
			    }
				got2 = append(got2, post2)*/
	    }

		// new struct for Org
		rows2 ,err, _,_ := config.Ap_sql("orgmem_get 1, 0,"+session_user_id_str,1)
	    if err != nil {
	    	panic(err.Error())
	    }
		type row2 struct {
	            Org_id int64
	            Org_name  string
	    }     
		got2 := []row2{}
		    for rows2.Next() {
		            var r row2
		            err = rows2.Scan(&r.Org_id, &r.Org_name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post2 := row2{
		            	r.Org_id,
				        r.Org_name,
				      
				    }
					got2 = append(got2, post2)
		    }

		 //district dropdown
		/*type District struct {
			   DisID int
	           District string
	    }*/
	    rows2 ,err, _,_ = config.Ap_sql("SysDropDown_Get 'AssignmentType'",1)
	    var name string
	    var Arrname []string
    	for rows2.Next() {
            err = rows2.Scan(&name)
            if err != nil {
                    panic(err.Error())
            }
			Arrname = append(Arrname, name)
		      
		   
		}

	    type Data struct {
			Org []row2 
			Type []string
			Parent string 
			LBR_Assign_Get LBR_Assign_Get 
			LBR_Assign_Get_drowdown_raw []LBR_Assign
			Rights string
		}
		htdata:= &Data{got2,Arrname,parent,LBR_Assign_Get_data,LBR_Assign_Get_drowdown_raw(session_user_id_str),rights}

		funcMap := template.FuncMap{
		    "date_formater": func(a time.Time) string {
		        return a.Format(" Jan _2 2006")
		    },
		    "time_formater": func(a time.Time) string {
		        return a.Format("15:04:05")
		    },

		}   
     	tmpl := template.New("CreateRate.html").Funcs(funcMap)
		//var err error
		if tmpl, err = tmpl.ParseFiles("hris/assignments/CreateRate.html"); err != nil {
			fmt.Println(err)
		}

		err1 := tmpl.Execute(w,htdata)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}	
	}else{	 //on edit
		r.ParseForm()
		session_username,_ := login.Get_account_info(r)	
		if r.Form["postcode"][0]==""{
			r.Form["postcode"][0] = "1"
		}
		if strings.TrimSpace(r.Form["parent"][0])==""{
			r.Form["parent"][0] = "null"
		} 
		sqlstr :=`
			exec LBR_Assign_Save 
				 'Add' 
				, '`+session_username+`'
				, 0 
				, 2 
				, '`+r.Form["name"][0] +`'
				, '`+r.Form["type"][0] +`'
				, `+r.Form["parent"][0] +`
				, '`+r.Form["location"][0] + `'
				, `+r.Form["postcode"][0] + `
				, '`+r.Form["sched_start"][0] + `'
				, `+r.Form["req_hrs_day"][0] + `
				, `+r.Form["req_hrs_wk"][0] + `
				, `+r.Form["graceperiod"][0] + `
				, `+r.Form["reg_late_interval"][0] + `
				, `+r.Form["reg_break_interval"][0] + `
				, `+r.Form["break_hr"][0] + `
				, `+r.Form["OT_hrs_min"][0] + `
				, `+r.Form["OT_hrs_interval"][0] + `
				, `+r.Form["ND_hrs_min"][0] + `
				, `+r.Form["ND_hrs_interval"][0] + `
				, '`+r.Form["default_in"][0] + `'
				, '`+r.Form["default_out"][0] + `'
				,  `+r.Form["sun_rest"][0] + `
				, '`+r.Form["sun_in"][0] + `'
				, '`+r.Form["sun_out"][0] + `'
				,  `+r.Form["mon_rest"][0] + `
				, '`+r.Form["mon_in"][0] + `' 
				, '`+r.Form["mon_out"][0] + `'
				,  `+r.Form["tue_rest"][0] + `
				, '`+r.Form["tue_in"][0] + `'
				, '`+r.Form["tue_out"][0] + `'
				,  `+r.Form["wed_rest"][0] + `
				, '`+r.Form["wed_in"][0] + `'
				, '`+r.Form["wed_out"][0] + `'
				,  `+r.Form["thu_rest"][0] + `
				, '`+r.Form["thu_in"][0] + `'
				, '`+r.Form["thu_out"][0] + `'
				,  `+r.Form["fri_rest"][0] + `
				, '`+r.Form["fri_in"][0] + `'
				, '`+r.Form["fri_out"][0] + `'
				,  `+r.Form["sat_rest"][0] + `
				, '`+r.Form["sat_in"][0] + `'
				, '`+r.Form["sat_out"][0] + `'
				, `+r.Form["target_all"][0] + `
				, `+r.Form["target_per"][0] + `
				, '`+r.Form["remarks"][0] + `'
				, '`+r.Form["mrate"][0] + `'
				, '`+r.Form["drate"][0] + `'
				, '`+r.Form["hrate"][0] + `'
		`
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
	
	profile := Profile{rtmsg.RetMsg,true}
	js, err := json.Marshal(profile)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
		}
  	w.Header().Set("Content-Type", "application/json")
  	w.Write(js)

	}		

	
}