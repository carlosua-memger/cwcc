package timekeeping

import (
	
		"net/http"
		"fmt"
		"text/template"
		config "../../config"
		"strconv"
		_"encoding/json"
		//datatables "../../datatables"
		login "../../login"
		"time"
		"encoding/json"

)

func  Nav_assignmentandlogs(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	_,session_user_id := login.Get_account_info(r)
	session_user_id_str := strconv.Itoa(session_user_id) 

	type Htdata struct {
		LBR_Assign_Get []lbr_assign_name
	}
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	var_htdata :=Htdata{func_LBR_Assign_Get(session_user_id_str, str_OrgID)}
	tmpl := template.New("assignmentandlogs_nav.html")
	var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
	if tmpl, err = tmpl.ParseFiles("timekeeping/time_keep/assignmentandlogs_nav.html"); err != nil {
		fmt.Println(err)
	}
	//err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf})
	err1 := tmpl.Execute(w,var_htdata)
	if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
	}

}

type lbr_assign_name struct{
	Id int
	Name string
}

func func_LBR_Assign_Get(  session_user_id_str string , org_id string) []lbr_assign_name {
	row ,err, _,_ := config.Ap_sql("exec LBR_Assign_Get 2, "+org_id+", "+session_user_id_str,1)
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


func func_LBR_Assign_Get_byID(  AssignID string) string {
	var name string
	row ,err, _,_ := config.Ap_sql(`select name from LBR_Assign  where id='`+AssignID+`'`,1)
	for row.Next() {
        err = row.Scan(&name)
        if err != nil {
                panic(err.Error())
        }
	}
    return name		   		
}
func  Dragsave(w http.ResponseWriter, r *http.Request) {
	go login.Session_renew(w,r)
	session_username,_ := login.Get_account_info(r)

	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	
	r.ParseForm()
	trn_type := r.Form["trn_type"][0]
	if(trn_type=="dragsave"){
		member	:=	r.Form["member"][0]
		assignid:=	r.Form["assignid"][0]
		logtimein	:=	r.Form["logtimein"][0]
		logtimeout	:=	r.Form["logtimeout"][0]
		logdate	:=	r.Form["logdate"][0]
		fmt.Println(member)
		fmt.Println(assignid)
		

		type LBR_LogPair_Save_struct struct{
			Id 			string
			Logdate 	string
			Assignment 	string
			Login 		string
			Logout 		string
			Daycode 	string
			Hr_reg 		interface{}
			Hr_ot 		interface{} 
			Hr_nd 		interface{}
		}
		var varlogdate time.Time
		var varlogin time.Time
		var varlogout time.Time
		
		//var varlogtime time.Time
		//fmt.Println(`lbr_logpair_save 'Add', 'session_username', 2, `+member+`, `+assignid+`, '`+login+`', '', '`+logdate+`', 'test'`)
		var sqlstr string
			sqlstr =`lbr_logpair_save 'Add', '`+session_username+`', `+str_OrgID+`, `+member+`, `+assignid+`, `+logtimein+`,`+logtimeout+` , '`+logdate+`', ''`
		fmt.Println(sqlstr)
		sql_row ,err, _,_ := config.Ap_sql(sqlstr,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    LBR_LogPair_Save_data := []LBR_LogPair_Save_struct{}
	    for sql_row.Next() {
	        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
	        var r LBR_LogPair_Save_struct
	        err = sql_row.Scan(&r.Id,&varlogdate, &r.Assignment, &varlogin,&varlogout,&r.Daycode,&r.Hr_reg,&r.Hr_ot, &r.Hr_nd)
	        if err != nil {
	                panic(err.Error())
	        }
			post2 := LBR_LogPair_Save_struct{
	        	r.Id,
	        	varlogdate.Format("20060102"),
		        r.Assignment ,
		        varlogin.Format("15:04:05") ,
		        varlogout.Format("15:04:05"), 
		        r.Daycode ,
		        config.Interfacetosting(r.Hr_reg) ,
		        config.Interfacetosting(r.Hr_ot) ,
		        config.Interfacetosting(r.Hr_nd) ,
		    }

			LBR_LogPair_Save_data = append(LBR_LogPair_Save_data, post2)
	    }

	    fmt.Println(LBR_LogPair_Save_data)
	    js, err := json.Marshal(LBR_LogPair_Save_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	w.Header().Set("Content-Type", "application/json")
	  	//w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)
	}else if(trn_type=="pair_delete"){
		
		pairID	:= r.Form["pairID"][0]
		fmt.Println(`LBR_LogPair_Delete '`+ session_username + `',`+pairID)
		sql_row ,err, _,_ := config.Ap_sql(`LBR_LogPair_Delete '`+ session_username + `',`+pairID,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    var r Return_stat
	    for sql_row.Next() {
	        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
	       err = sql_row.Scan(&r.RetNo,&r.RetMsg)
	        if err != nil {
	                panic(err.Error())
	        }
			r = Return_stat{r.RetNo,r.RetMsg}
		}			
		fmt.Println(r)
	    js, err := json.Marshal(r)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	w.Header().Set("Content-Type", "application/json")
	  	//w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)
	
	}else if(trn_type=="addlog"){

		
		member		:= r.Form["member"][0]
		assignid	:= r.Form["assignid"][0]
		logdate		:= r.Form["logdate"][0]
		logtime		:= r.Form["logtime"][0]
		
		fmt.Println(`LBR_LogRaw_Save '`+ session_username + `',`+str_OrgID+`,`+member+`,`+assignid+`,'`+ logdate+` `+ logtime+`'`)
		sql_row ,err, _,_ := config.Ap_sql(`LBR_LogRaw_Save '`+ session_username + `',`+str_OrgID+`,`+member+`,`+assignid+`,'`+ logdate +` `+ logtime+`'`,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    var r Return_stat
	    for sql_row.Next() {
	        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
	       err = sql_row.Scan(&r.RetNo,&r.RetMsg)
	        if err != nil {
	                panic(err.Error())
	        }
			r = Return_stat{r.RetNo,r.RetMsg}
		}			
		fmt.Println(r)
	    js, err := json.Marshal(r)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	w.Header().Set("Content-Type", "application/json")
	  	//w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)
	}else if(trn_type=="mem_nav"){
		
		member		:= r.Form["member"][0]
		assignid	:= r.Form["assignid"][0]
		nav_type	:= r.Form["nav_type"][0]

		type LBR_MemAss_Search_struct struct{
			Id string
			Member string 
		}

		var sql string
		if(nav_type=="0"){
			sql = `LBR_MemAss_Search 4, `+ assignid + `,`+member 
		}else{
			sql = `LBR_MemAss_Search 3, `+ assignid + `,`+member
		}
		fmt.Println(sql)
		sql_row ,err, _,_ := config.Ap_sql(sql ,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    var r LBR_MemAss_Search_struct
	    for sql_row.Next() {
	        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
	       err = sql_row.Scan(&r.Id,&r.Member)
	        if err != nil {
	                panic(err.Error())
	        }
			r = LBR_MemAss_Search_struct{r.Id,r.Member}
		}			
		fmt.Println(r)
	    js, err := json.Marshal(r)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)
	
	}else if(trn_type=="updateassignment"){

		pid		:= r.Form["pairID"][0]
		assignid	:= r.Form["assignid"][0]
		
		fmt.Println( `LBR_LogPair_UpdateAssignment '`+ session_username + `',`+pid+`,`+assignid)
		sql_row ,err, _,_ := config.Ap_sql( `LBR_LogPair_UpdateAssignment '`+ session_username + `',`+pid+`,`+assignid ,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    var r Return_stat
	    for sql_row.Next() {
	        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
	       err = sql_row.Scan(&r.RetNo,&r.RetMsg)
	        if err != nil {
	                panic(err.Error())
	        }
			r = Return_stat{r.RetNo,r.RetMsg}
		}			
		fmt.Println(r)
	    js, err := json.Marshal(r)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	w.Header().Set("Content-Type", "application/json")
	  	//w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)

	}else if(trn_type=="updatelineND"){
		pid		:= r.Form["pid"][0]
		hr_nd	:= r.Form["hr_nd"][0]
		
		fmt.Println( `LBR_LogPair_UpdateHR_ND '`+ session_username + `',`+pid+`,`+hr_nd)
		sql_row ,err, _,_ := config.Ap_sql( `LBR_LogPair_UpdateHR_ND '`+ session_username + `',`+pid+`,`+hr_nd ,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    var r Return_stat
	    for sql_row.Next() {
	        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
	       err = sql_row.Scan(&r.RetNo,&r.RetMsg)
	        if err != nil {
	                panic(err.Error())
	        }
			r = Return_stat{r.RetNo,r.RetMsg}
		}			
		fmt.Println(r)
	    js, err := json.Marshal(r)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	w.Header().Set("Content-Type", "application/json")
	  	//w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)

	}else if(trn_type=="LBR_LogTotal_Update"){
		id 				:= r.Form["Id"][0]
		lasttouch		:= session_username
		lbr_daycode 	:= r.Form["Lbr_daycode"][0]
		hr_reg 			:= r.Form["Hr_reg"][0]
		hr_ot 			:= r.Form["Hr_ot"][0]
		hr_nd 			:= r.Form["Hr_nd"][0]
		hr_ndot 		:= r.Form["Hr_ndot"][0]
		hr_break 		:= r.Form["Hr_break"][0]
		min_late 		:= r.Form["Min_late"][0]
		min_undertime 	:= r.Form["Min_undertime"][0]
		remarks 		:= r.Form["Remarks"][0]


		
		fmt.Println( `LBR_LogTotal_Update `+id+`,'`+ lasttouch + `',`+lbr_daycode + `,`+ hr_reg + `,`+hr_ot+`,`+ hr_nd+`,`+hr_ndot+`,`+hr_break+ `,`+min_late+ `,`+min_undertime+`,'`+remarks+`'`)
		sql_row ,err, _,_ := config.Ap_sql( `LBR_LogTotal_Update `+id+`,'`+ lasttouch + `',`+lbr_daycode + `,`+ hr_reg + `,`+hr_ot+`,`+ hr_nd+`,`+hr_ndot+`,`+hr_break+ `,`+min_late+ `,`+min_undertime+`,'`+remarks+`'` ,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    var r Return_stat
	    for sql_row.Next() {
	        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
	       err = sql_row.Scan(&r.RetNo,&r.RetMsg)
	        if err != nil {
	                panic(err.Error())
	        }
			r = Return_stat{r.RetNo,r.RetMsg}
		}			
		fmt.Println(r)
	    js, err := json.Marshal(r)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	w.Header().Set("Content-Type", "application/json")
	  	//w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)
	}	
}

type Return_stat struct{
	RetNo string
	RetMsg string
}

func Ajax_lbr_lograw_get(w http.ResponseWriter, r *http.Request) {
	member := r.URL.Query().Get("member")
	lbr_assign := r.URL.Query().Get("lbr_assign")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	sp := r.URL.Query().Get("sp")
	fmt.Println("sp no",sp)

	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)

	if(sp=="sp1"){


		type Struct_lbr_lograw_get struct{
			Id string
			LogDate string
			Logtime string
			Lbr_logpair string
		}
		var strID string
		var varlogdate time.Time
		var varlogtime time.Time
		var strLbr_logpair string
		fmt.Println(`lbr_lograw_get 1, `+str_OrgID+`, `+member+`, `+lbr_assign+`, '`+from+`' , '`+to+`'`)
		sql_row ,err, _,_ := config.Ap_sql(`lbr_lograw_get 1, `+str_OrgID+`, `+member+`, `+lbr_assign+`, '`+from+`' , '`+to+`'`,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    lbr_lograw_get_data := []Struct_lbr_lograw_get{}
	    for sql_row.Next() {
	        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
	        err = sql_row.Scan(&strID,&varlogdate, &varlogtime,&strLbr_logpair)
	        if err != nil {
	                panic(err.Error())
	        }
			post2 := Struct_lbr_lograw_get{
	        	strID,
	        	varlogdate.Format("20060102"),
		        varlogtime.Format("15:04:05"),
		        strLbr_logpair,
		    }
			lbr_lograw_get_data = append(lbr_lograw_get_data, post2)
	    }

	    fmt.Println(lbr_lograw_get_data)
	    js, err := json.Marshal(lbr_lograw_get_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	w.Header().Set("Content-Type", "application/json")
	  	//w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)
	}else if(sp=="sp2"){

		type Struct_lbr_logpair_get struct{
			Id string
			Logdate string
			Assignment string
			Login string
			Logout string
			Daycode string
			Hr_reg interface{}
			Hr_ot interface{}
			Hr_nd interface{}
		}
		var varlogdate time.Time
		var vartime time.Time
		var vartimeout time.Time
		//var varlogtime time.Time
		fmt.Println(`lbr_logpair_get 1,0, `+str_OrgID+`, `+member+`, `+lbr_assign+`, '`+from+`' , '`+to+`'`)
		sql_row ,err, _,_ := config.Ap_sql(`lbr_logpair_get 1,0, `+str_OrgID+`, `+member+`, `+lbr_assign+`, '`+from+`' , '`+to+`'`,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    lbr_logpair_get_data := []Struct_lbr_logpair_get{}
	    for sql_row.Next() {
	    	var r Struct_lbr_logpair_get
	        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
	        err = sql_row.Scan(&r.Id,&varlogdate,&r.Assignment,&vartime,&vartimeout,&r.Daycode,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd)
	        if err != nil {
	                panic(err.Error())
	        }
			post2 := Struct_lbr_logpair_get{
				r.Id,
				varlogdate.Format("20060102"),
		        r.Assignment,
		        vartime.Format("15:04:05"),
		        vartimeout.Format("15:04:05"),
		        r.Daycode,
		        config.Interfacetosting(r.Hr_reg),
		        config.Interfacetosting(r.Hr_ot),
		        config.Interfacetosting(r.Hr_nd),
			}

			lbr_logpair_get_data = append(lbr_logpair_get_data, post2)
	    }

	    fmt.Println(lbr_logpair_get_data)
	    js, err := json.Marshal(lbr_logpair_get_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	w.Header().Set("Content-Type", "application/json")
	  	//w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)

	}else if(sp=="sp3"){


		type Struct_lbr_OT_get struct{
			LogDate string
			Hours string
		}
		var varlogdate time.Time
		var hrs string
		//var varlogtime time.Time
		fmt.Println(`lbr_OT_get 1, `+member+`, `+lbr_assign+`, '`+from+`' , '`+to+`'`)
		sql_row ,err, _,_ := config.Ap_sql(`lbr_OT_get 1, `+member+`, `+lbr_assign+`, '`+from+`' , '`+to+`'`,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    lbr_OT_get_data := []Struct_lbr_OT_get{}
	    for sql_row.Next() {
	        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
	        
	        err = sql_row.Scan(&varlogdate, &hrs)
	        if err != nil {
	                panic(err.Error())
	        }
			post2 := Struct_lbr_OT_get{
	        	varlogdate.Format("20060102"),
		        hrs,
		    }
			lbr_OT_get_data = append(lbr_OT_get_data, post2)
	    }

	    fmt.Println(lbr_OT_get_data)
	    js, err := json.Marshal(lbr_OT_get_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	w.Header().Set("Content-Type", "application/json")
	  	//w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)
	}else if(sp=="LBR_Assign_Get"){   //LBR_LogPair_UpdateAssignment	
		
		type LBR_LogPair_UpdateAssignment_struct struct{
			Id string 		`json:"id"`
			Name string 	`json:"text"`  	
		}
		fmt.Println(`LBR_Assign_Get 3, `+ lbr_assign )
		sql_row ,err, _,_ := config.Ap_sql(`LBR_Assign_Get 3,`+lbr_assign,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    LBR_LogPair_UpdateAssignment_data := []LBR_LogPair_UpdateAssignment_struct{}
	    for sql_row.Next() {
	        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
	        var r LBR_LogPair_UpdateAssignment_struct
	        err = sql_row.Scan(&r.Id, &r.Name)
	        if err != nil {
	                panic(err.Error())
	        }
			post2 := LBR_LogPair_UpdateAssignment_struct{
	        	r.Id,
		        r.Name,
		    }
			LBR_LogPair_UpdateAssignment_data = append(LBR_LogPair_UpdateAssignment_data, post2)
	    }

	    fmt.Println(LBR_LogPair_UpdateAssignment_data)
	    js, err := json.Marshal(LBR_LogPair_UpdateAssignment_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	w.Header().Set("Content-Type", "application/json")
	  	//w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)

	}else if(sp=="sp4"){

		type LBR_LogTotal_Get struct{
			Id 					interface{}
			Org 				interface{}
			Member 				interface{}
			Logdate 			string
			Lbr_daycode 		interface{}
			Lbr_daycode_name 	interface{}
			Hr_reg 				interface{}
			Hr_ot 				interface{}
			Hr_nd 				interface{}
			Hr_ndot 			interface{}
			Hr_break 			interface{}
			Min_late 			interface{}
			Min_undertime  		interface{}
			Remarks 			interface{}
		
		}
		var varlogdate time.Time
		/*var vartime time.Time
		var vartimeout time.Time*/
		//var varlogtime time.Time
		fmt.Println(`LBR_LogTotal_Get 1,0, `+str_OrgID+`, `+member+`,'`+from+`' , '`+to+`'`)
		sql_row ,err, _,_ := config.Ap_sql(`LBR_LogTotal_Get 1,0, `+str_OrgID+`, `+member+`,'`+from+`' , '`+to+`'`,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    LBR_LogTotal_Get_data := []LBR_LogTotal_Get{}
	    for sql_row.Next() {
	    	var r LBR_LogTotal_Get
	        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
	        //err = sql_row.Scan(&r.Id,&varlogdate,&r.Assignment,&vartime,&vartimeout,&r.Daycode,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd)
	        err = sql_row.Scan(&r.Id,&r.Org,&r.Member,&varlogdate,&r.Lbr_daycode,&r.Lbr_daycode_name,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot, &r.Hr_break,&r.Min_late,&r.Min_undertime,&r.Remarks)
	        
	        if err != nil {
	                panic(err.Error())
	        }
			post2 := LBR_LogTotal_Get{
				config.Interfacetosting(r.Id),
				config.Interfacetosting(r.Org),
				config.Interfacetosting(r.Member),
				varlogdate.Format("20060102"),
		        config.Interfacetosting(r.Lbr_daycode),
		        config.Interfacetosting(r.Lbr_daycode_name),
		        config.Interfacetosting(r.Hr_reg),
		        config.Interfacetosting(r.Hr_ot),
		        config.Interfacetosting(r.Hr_nd),
		        config.Interfacetosting(r.Hr_ndot), 
		        config.Interfacetosting(r.Hr_break),
		        config.Interfacetosting(r.Min_late),
		        config.Interfacetosting(r.Min_undertime),
		        config.Interfacetosting(r.Remarks),

		        

		       
			}

			LBR_LogTotal_Get_data = append(LBR_LogTotal_Get_data, post2)
	    }

	    fmt.Println(LBR_LogTotal_Get_data)
	    js, err := json.Marshal(LBR_LogTotal_Get_data)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	w.Header().Set("Content-Type", "application/json")
	  	//w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)

	}	
}