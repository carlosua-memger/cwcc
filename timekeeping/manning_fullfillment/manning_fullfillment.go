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
		//"time"

)

func  Nav_Report_mfr(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	_,session_user_id := login.Get_account_info(r)
	session_user_id_str := strconv.Itoa(session_user_id) 

	type Htdata struct {
		LBR_Assign_Get []lbr_assign_name
	}
	var_htdata :=Htdata{func_LBR_Assign_Get(session_user_id_str)}
	tmpl := template.New("report_nav_mfr.html")
	var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
	if tmpl, err = tmpl.ParseFiles("timekeeping/manning_fullfillment/report_nav_mfr.html"); err != nil {
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

func func_LBR_Assign_Get(  session_user_id_str string) []lbr_assign_name {
	row ,err, _,_ := config.Ap_sql("exec LBR_Assign_Get 2, 2, "+session_user_id_str,1)
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


func  Report(w http.ResponseWriter, r *http.Request) {
	assignment  := r.URL.Query().Get("lbr_assign")
	transdate  := r.URL.Query().Get("date")

	type mfr_list struct{
		Assignment string
		Deployed int
		Present int
		Absent 	int
		Rating string
	}
	fmt.Println(`lbr_rep_logpair 2, 2, 0, `+assignment+`, '', '`+transdate+`'`)

	   	LBR_LogHdr_row ,err, _,_ := config.Ap_sql(`lbr_rep_logpair 2, 2, 0, `+assignment+`, '', '`+transdate+`'`,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    mfr_list_data := []mfr_list{}
	    for LBR_LogHdr_row.Next() {
            var r mfr_list
            //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
            err = LBR_LogHdr_row.Scan(&r.Assignment, &r.Deployed,&r.Present,&r.Absent,&r.Rating)
            if err != nil {
                    panic(err.Error())
            }
			post2 := mfr_list{
            	r.Assignment,
		        r.Deployed,
		        r.Present,
		      	r.Absent,
		      	r.Rating,
 		    }
			mfr_list_data = append(mfr_list_data, post2)
	    }


	type Htdata struct {
		Mfr_list_data []mfr_list
		Assignment string
		Date string
	}
	var_htdata :=Htdata{mfr_list_data,func_LBR_Assign_Get_byID(assignment),transdate}
	tmpl := template.New("report.html")
	//var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
	if tmpl, err = tmpl.ParseFiles("timekeeping/manning_fullfillment/report.html"); err != nil {
		fmt.Println(err)
	}
	//err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf})
	err1 := tmpl.Execute(w,var_htdata)
	if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
	}

}



func  Report_main(w http.ResponseWriter, r *http.Request) {
	assignment  := r.URL.Query().Get("lbr_assign")
	transdate  := r.URL.Query().Get("date")
	fmt.Println(transdate)
	type mfr_list struct{
		Assignment string
		Deployed int
		Present int
		Absent 	int
		Rating string
	}
	fmt.Println(`lbr_rep_logpair 2, 2, 0, `+assignment+`, '', '`+transdate+`'`)

	   	LBR_LogHdr_row ,err, _,_ := config.Ap_sql(`lbr_rep_logpair 2, 2, 0, `+assignment+`, '', '`+transdate+`'`,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    mfr_list_data := []mfr_list{}
	    for LBR_LogHdr_row.Next() {
            var r mfr_list
            //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
            err = LBR_LogHdr_row.Scan(&r.Assignment, &r.Deployed,&r.Present,&r.Absent,&r.Rating)
            if err != nil {
                    panic(err.Error())
            }
			post2 := mfr_list{
            	r.Assignment,
		        r.Deployed,
		        r.Present,
		      	r.Absent,
		      	r.Rating,
 		    }
			mfr_list_data = append(mfr_list_data, post2)
	    }


	type Htdata struct {
		Mfr_list_data []mfr_list
		Assignment string
		Date string
	}
	var_htdata :=Htdata{mfr_list_data,func_LBR_Assign_Get_byID(assignment),transdate}
	tmpl := template.New("report_main.html")
	//var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
	if tmpl, err = tmpl.ParseFiles("timekeeping/manning_fullfillment/report_main.html"); err != nil {
		fmt.Println(err)
	}
	//err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf})
	err1 := tmpl.Execute(w,var_htdata)
	if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
	}

}