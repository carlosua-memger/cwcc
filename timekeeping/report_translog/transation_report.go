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
		"strings"
	//	"math"

)



func  Lbr_rep_logpair_rpt(w http.ResponseWriter, r *http.Request) {
	tconf := make(map[string]string)
	tconf["datefrom"] = r.URL.Query().Get("datefrom")
	tconf["dateto"] = r.URL.Query().Get("dateto")
	tconf["lbr_assign"] = r.URL.Query().Get("lbr_assign")
	tconf["member"] = r.URL.Query().Get("member")
	fmt.Println(`datefrom` ,tconf["datefrom"])
	fmt.Println(`dateto` ,tconf["dateto"])
	var datevar time.Time
	var datein time.Time
	var dateout time.Time

	type Lbr_rep_logpair struct{
		Assignment_name string
		Member_name string
		Logdate string
		Daycode string
		Login string
		Logout string
	}
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)

	fmt.Println(`lbr_rep_logpair 3, `+str_OrgID+`, `+tconf["member"]+` , `+tconf["lbr_assign"]+`, '`+tconf["datefrom"]+`', '`+tconf["dateto"] +`'`)
	LBR_LogHdr_row ,err, _,_ := config.Ap_sql(`lbr_rep_logpair 3, `+str_OrgID+`, `+tconf["member"]+` , `+tconf["lbr_assign"]+`, '`+tconf["datefrom"]+`', '`+tconf["dateto"] +`'`,1)
    if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	panic(err.Error())
    }
    mfr_list_data := []Lbr_rep_logpair{}
    for LBR_LogHdr_row.Next() {
        var r Lbr_rep_logpair
        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
        err = LBR_LogHdr_row.Scan(&r.Assignment_name, &r.Member_name,&datevar,&r.Daycode,&datein,&dateout)
        if err != nil {
                panic(err.Error())
        }
		post2 := Lbr_rep_logpair{
        	r.Assignment_name,
	        r.Member_name,
	        datevar.Format("2006-01-02"),
	      	r.Daycode,
	      	datein.Format("2006-01-02 15:04:05.999"),
	      	dateout.Format("2006-01-02 15:04:05.999"),
	      }
		mfr_list_data = append(mfr_list_data, post2)
    }

	type Htdata struct {
		Mfr_list_data []Lbr_rep_logpair
		
	}
	var_htdata :=Htdata{mfr_list_data}
	tmpl := template.New("report_Lbr_rep_logpair.html")
	//var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
	if tmpl, err = tmpl.ParseFiles("timekeeping/report_translog/report_Lbr_rep_logpair.html"); err != nil {
		fmt.Println(err)
	}
	//err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf})
	err1 := tmpl.Execute(w,var_htdata)
	if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
	}

}

func  Nav_Report_member_log(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	_,session_user_id := login.Get_account_info(r)
	session_user_id_str := strconv.Itoa(session_user_id) 

	type Htdata struct {
		LBR_Assign_Get []lbr_assign_name
	}
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	
	var_htdata :=Htdata{func_LBR_Assign_Get(session_user_id_str, str_OrgID)}
	tmpl := template.New("report_nav_logreport.html")
	var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
	if tmpl, err = tmpl.ParseFiles("timekeeping/report_translog/report_nav_logreport.html"); err != nil {
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



func Func_lbr_rep_lograw(w http.ResponseWriter, r *http.Request) {
 	
 	tconf := make(map[string]string)
	tconf["datefrom"] = r.URL.Query().Get("datefrom")
	tconf["dateto"] = r.URL.Query().Get("dateto")
	tconf["lbr_assign"] = r.URL.Query().Get("lbr_assign")
	tconf["member"] = r.URL.Query().Get("member")
	
 	type Lbr_rep_lograw struct{
 		Member_name 	interface{}
 		Assignment_name interface{}
 		Logassign		interface{}
 		Login           interface{}
 		Logtime         interface{} 
 		Lbr_logpair     interface{}
 		Lbr_logpaired   interface{}
 	}


 	//fmt.Println(`lbr_rep_logpair 3, 2, `+tconf["member"]+` , `+tconf["lbr_assign"]+`, '`+tconf["datefrom"]+`', '`+tconf["dateto"] +`'`)
 	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	
 	fmt.Println(`lbr_rep_lograw  2, `+str_OrgID+`, `+tconf["member"]+` , `+tconf["lbr_assign"]+` , '', '', '`+tconf["datefrom"]+`', '`+tconf["dateto"]+`'`)
	
	//LBR_LogHdr_row ,err, _,_ := config.Ap_sql(`lbr_rep_logpair 3, 2, `+tconf["member"]+` , `+tconf["lbr_assign"]+`, '`+tconf["datefrom"]+`', '`+tconf["dateto"] +`'`,1)
	lbr_rep_lograw_row ,err, _,_ := config.Ap_sql(`lbr_rep_lograw  2, `+str_OrgID+`, `+tconf["member"]+` , `+tconf["lbr_assign"]+` , '', '', '`+tconf["datefrom"]+`', '`+tconf["dateto"]+`'`,1)
    if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	panic(err.Error())
    }
    Lbr_rep_lograw_data := []Lbr_rep_lograw{}
    for lbr_rep_lograw_row.Next() {
        var r Lbr_rep_lograw
        //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
        err = lbr_rep_lograw_row.Scan(&r.Member_name, &r.Assignment_name,&r.Logassign,&r.Login,&r.Logtime,&r.Lbr_logpair,&r.Lbr_logpaired)
        if err != nil {
                panic(err.Error())
        }
		post2 := Lbr_rep_lograw{
        	r.Member_name,
	        r.Assignment_name,
	        r.Logassign,
	      	r.Login,
	      	r.Logtime,
	      	r.Lbr_logpair,
	      	r.Lbr_logpaired,
	    }

		Lbr_rep_lograw_data = append(Lbr_rep_lograw_data, post2)
    }

	type Htdata struct {
		Lbr_rep_lograw_data []Lbr_rep_lograw
		
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
						a =fmt.Sprint("NULL")
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
						a = fmt.Sprint(v.Format("Jan _2 2006 03:04:05PM"))

					default:
				//		fmt.Print(v)
						a= fmt.Sprint(v)
					}
					return a
				},

			}

	var_htdata :=Htdata{Lbr_rep_lograw_data}
	tmpl := template.New("LBR_Rep_LogRaw_report.html").Funcs(funcMap)

	//var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
	if tmpl, err = tmpl.ParseFiles("timekeeping/report_translog/LBR_Rep_LogRaw_report.html"); err != nil {
		fmt.Println(err)
	}
	//err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf})
	err1 := tmpl.Execute(w,var_htdata)
	if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
	}




}


func  Report_member_log3(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("here..!")
	login.Session_renew(w,r)
	
	tconf := make(map[string]string)
	tconf["datefrom"] = r.URL.Query().Get("datefrom")
	tconf["dateto"] = r.URL.Query().Get("dateto")
	tconf["lbr_assign"] = r.URL.Query().Get("lbr_assign")
	tconf["member"] = r.URL.Query().Get("member")
	tconf["group"] = r.URL.Query().Get("member")

	//conve ,_ := strconv.ParseInt( tconf["member"] , 10, 32)
	conve,c_err := strconv.Atoi(tconf["member"])
	if(c_err != nil){
		tconf["member"] = "0"
		tconf["group"] = r.URL.Query().Get("member")
	}else{
		tconf["group"] = ""
	}
	fmt.Println(c_err)
	conve_page := int(conve)
	fmt.Println(`conve inte here`)
	fmt.Println(conve_page)


	tconf["no_page"] = r.URL.Query().Get("no_page")
	no_page64 ,_ := strconv.ParseInt( tconf["no_page"] , 10, 32)
	no_page := int(no_page64)
	/*type Lbr_Rep_LogRegister struct{
	   		Member_name string
	   		Assignment string
	   		LogDate string 
	   		DayCode string
	   		In_Sched string
	   		Out_Sched string
	   		In_Time string
	   		Out_Time string
	   		Hr_total string
	   		Hr_reg string
	   		Hr_ot string
	   		Hr_nd string
	   		Hr_ndot string
	   		colormered int
		}*/
	/*type Lbr_Rep_LogRegister struct{
	   		Member_name string
	   		Logdate string
	   		Lbr_daycode_name string 
	   		Assignment string
	   		Login string
	   		Logout string
	   		Final_hrreg string
	   		Final_hrot string
	   		Final_hrnd string
	   		Min_late string
	   		Raw_log 		string
	   		
	}*/


	/*var Logdate interface{}
	var Lbr_daycode_name interface{}
	var Login interface{}
	var Logout interface{}
	var Final_hrreg interface{}
	var Final_hrot interface{}
	var Final_hrnd interface{}
	var Min_late  interface{}
	var Raw_log  interface{}*/

	type Lbr_Rep_LogRegister_summary struct{
   		Member_name 	   string
   		Lbr_daycode_name 		interface{}
   		Dayspresent 		interface{} 
   		Total_hrreg 			interface{}
   		Total_hrot 			interface{}
   		Total_hrnd 			interface{}
   		Total_minlate 		interface{}
   	}

   	type Lbr_Rep_LogRegister struct{
			logdate 			interface{}	
			leader 				interface{}	
			leader_name 		string	
			member 				interface{}	
			member_name 		string	
			lbr_assign 			interface{}	
			lbr_assign_name 	interface{}	
			login 				interface{}	
			logout 				interface{}	
			lbr_daycode 		interface{}	
			ot_approved 		interface{}	
			lbr_logtotal 		interface{}	
			reg_reg 			interface{}	
			reg_ot 				interface{}	
			reg_nd 				interface{}	
			reg_ndot 			interface{}	
			reg_late 			interface{}	
			reg_undertime 		interface{}	
			rd_reg 				interface{}	
			rd_ot 				interface{}	
			rd_nd 				interface{}	
			rd_ndot 			interface{}	
			rd_late 			interface{}	
			rd_undertime 		interface{}	
			sh_reg 				interface{}	
			sh_ot 				interface{}	
			sh_nd 				interface{}	
			sh_ndot 			interface{}	
			sh_late 			interface{}	
			sh_undertime 		interface{}	
			shrd_reg 			interface{}	
			shrd_ot 			interface{}	
			shrd_nd 			interface{}	
			shrd_ndot 			interface{}	
			shrd_late 			interface{}	
			shrd_undertime 		interface{}	
			rh_reg 				interface{}	
			rh_ot 				interface{}	
			rh_nd 				interface{}	
			rh_ndot 			interface{}	
			rh_late 			interface{}	
			rh_undertime 		interface{}	
			rh_100 				interface{}	
			rhrd_reg 			interface{}	
			rhrd_ot 			interface{}	
			rhrd_nd 			interface{}	
			rhrd_ndot 			interface{}	
			rhrd_late 			interface{}	
			rhrd_undertime 		interface{}	
   	}
			
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
	
	   	fmt.Println(`lbr_rep_logpair 6,`+str_OrgID+`,`+tconf["member"]+`,`+tconf["lbr_assign"]+`,'`+tconf["datefrom"]+ `','`+tconf["dateto"] + `','`+tconf["group"]+`'`)

	   	LBR_LogHdr_row ,err, _,_ := config.Ap_sql(`lbr_rep_logpair 6,`+str_OrgID+`,`+tconf["member"]+`,`+tconf["lbr_assign"]+`,'`+tconf["datefrom"]+ `','`+tconf["dateto"] + `','`+tconf["group"]+`'`,1)
	   	//LBR_LogHdr_row ,err, _,_ := config.Ap_sql(`lbr_rep_logpair 4, 2, 1, 0, '1 Nov 2015', '15 Nov 2015'`,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    Lbr_Rep_LogRegister_data := []Lbr_Rep_LogRegister{}
	    

	    for LBR_LogHdr_row.Next() {
            var r Lbr_Rep_LogRegister
            err = LBR_LogHdr_row.Scan(
				&r.logdate, 
				&r.leader, 
				&r.leader_name, 
				&r.member, 
				&r.member_name, 
				&r.lbr_assign, 
				&r.lbr_assign_name, 
				&r.login, 
				&r.logout, 
				&r.lbr_daycode, 
				&r.ot_approved, 
				&r.lbr_logtotal, 
				&r.reg_reg, 
				&r.reg_ot, 
				&r.reg_nd, 
				&r.reg_ndot, 
				&r.reg_late, 
				&r.reg_undertime, 
				&r.rd_reg, 
				&r.rd_ot, 
				&r.rd_nd, 
				&r.rd_ndot, 
				&r.rd_late, 
				&r.rd_undertime, 
				&r.sh_reg, 
				&r.sh_ot, 
				&r.sh_nd, 
				&r.sh_ndot, 
				&r.sh_late, 
				&r.sh_undertime, 
				&r.shrd_reg, 
				&r.shrd_ot, 
				&r.shrd_nd, 
				&r.shrd_ndot, 
				&r.shrd_late, 
				&r.shrd_undertime, 
				&r.rh_reg, 
				&r.rh_ot, 
				&r.rh_nd, 
				&r.rh_ndot, 
				&r.rh_late, 
				&r.rh_undertime, 
				&r.rh_100, 
				&r.rhrd_reg, 
				&r.rhrd_ot, 
				&r.rhrd_nd, 
				&r.rhrd_ndot, 
				&r.rhrd_late, 
				&r.rhrd_undertime,
			)
            
            if err != nil {
                    panic(err.Error())
            }

  		    post2 := Lbr_Rep_LogRegister{
					r.logdate, 
					r.leader, 
					r.leader_name, 
					r.member, 
					r.member_name, 
					r.lbr_assign, 
					r.lbr_assign_name, 
					r.login, 
					r.logout, 
					r.lbr_daycode, 
					r.ot_approved, 
					r.lbr_logtotal, 
					r.reg_reg, 
					r.reg_ot, 
					r.reg_nd, 
					r.reg_ndot, 
					r.reg_late, 
					r.reg_undertime, 
					r.rd_reg, 
					r.rd_ot, 
					r.rd_nd, 
					r.rd_ndot, 
					r.rd_late, 
					r.rd_undertime, 
					r.sh_reg, 
					r.sh_ot, 
					r.sh_nd, 
					r.sh_ndot, 
					r.sh_late, 
					r.sh_undertime, 
					r.shrd_reg, 
					r.shrd_ot, 
					r.shrd_nd, 
					r.shrd_ndot, 
					r.shrd_late, 
					r.shrd_undertime, 
					r.rh_reg, 
					r.rh_ot, 
					r.rh_nd, 
					r.rh_ndot, 
					r.rh_late, 
					r.rh_undertime, 
					r.rh_100, 
					r.rhrd_reg, 
					r.rhrd_ot, 
					r.rhrd_nd, 
					r.rhrd_ndot, 
					r.rhrd_late, 
					r.rhrd_undertime,
     	
 		    }

			Lbr_Rep_LogRegister_data = append(Lbr_Rep_LogRegister_data, post2)
	    }


		t := time.Now().Local()
		//newFormat := t.Format("2006-01-02 15:00:00 +0800")
		newFormat := t.Format("Mon Jan _2 2006 03:04:05PM")
		        

		var id int
		var name string
		var location string
		Org_Get_row ,err, _,_ := config.Ap_sql(" Org_Get 1, 2",1)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    for Org_Get_row.Next() {
	    	err = Org_Get_row.Scan(&id, &name,&location)
            if err != nil {
                    panic(err.Error())
            }
        }
	    
        
	    ///version1
	    strt:= `
	    <html>
			<head>
				<title style="font-size: 16px;">Time keeping report</title>
			</head>
			<style type="text/css">
				.rpt_date {
					float: right;
					
				}
				.page_header{
					margin-bottom: 5px;
					    float: right;
					/*display: inline; */
				}
				.page_header_left{
					margin-bottom: 5px;
					    float: left;
					/*display: inline; */
				}
				.name{
					/*float: left;*/
					margin-bottom: 3px;
					font-family: Arial; 
				}
				.covered_data{
					/*float: right;*/
					font-family: Arial; 
				}
				/*.pic_logo {
					width: 200px;

				}*/
				table{
					width: 100%;
				}
				.tb_r{
					text-align: right;
				}
				
				.dtatbl {
				
				border: 0.99px solid #B3B3B3;
				border-collapse: collapse;
				font-family: 'Arial Narrow', Arial, sans-serif;
				font-size: 10px;
				}

				.dtatbl th{
					border: 0.99px solid #B3B3B3;
				border-collapse: collapse;
				font-family: 'Arial Narrow', Arial, sans-serif;
				font-size: 10px;
				}

				.dtatbl td {
					border: 0.99px solid #B3B3B3;
				border-collapse: collapse;
				font-family: 'Arial Narrow', Arial, sans-serif;
				font-size: 10px;
				}
				th, td {
				    /*padding: 2px;*/
				}
				th{
					
					background-color: whiteSmoke;
					/*background-image: -webkit-linear-gradient*/
					
				}
				body {
					font-size: 10px;
				}
				.logo{
					width: 200px;
					margin-bottom: 5px;
				}
				.page_footer{
					margin-top: 20px;
					display: inline-block;
					width: 100%;
				}
				.signaturemember{
					border-top: 0.99px solid #B3B3B3;
					width: 25%;
					float: left;
					font-family: inherit;
				}
				.signatureleft{
					margin-left: 10%;
					border-top: 0.99px solid #B3B3B3;
					width: 25%;
					float: left;
					font-family: inherit; 
				}
				.signatureright{
					border-top: 0.99px solid #B3B3B3;
					width: 25%;
					float: right;
					font-family: inherit; 
				}
				.tbl_foor{
					margin-top: 10px;
				}
				.tbl_foor td div{
					border-top: 0.99px solid #B3B3B3;	
					width: 50%;
					font-size: 10px;
				}
				.tbl_floor{
					margin-top: 10px;
				}
				.tbl_floor td div{
					width: 50%;
					font-size: 10px;
				}
				.isred {
					color: red;
				}
				.log_sum{
					margin-top: 10px;
				}

			</style>

			<body>
		`
		fmt.Fprint(w,strt)
		//counter :=0
		//counter1:=0
		Container :=" "
		counterRemarks :=0
		//TotalHours :=0
		//Contain := ""
		var breakcont int = 0
		for _ , val := range Lbr_Rep_LogRegister_data {
		/*	if counter1==0{
				Contain = val.Member_name
			}
			if val.Member_name==Contain{
				fmt.Fprint(w,`---******` + val.Member_name)
			}

			*/

			if counterRemarks!=0{
				if Container != val.leader_name{
					
							strt=`
							</table>
							



										
									
					</div>
					

					
				
				<table  class="tbl_foor">
					<tbody>
						<tr>
							<td><div>COOP COORDINATOR</div></td>
							<td><div>FIELD ENG'R/FOREMAN/QUANTITY SURVEYOR</div></td>
							<td><div>MDC PHR</div></td>
						</tr>
					</tbody>
				</table>	


							
							`
							fmt.Fprint(w,strt)
					if(breakcont< (no_page-1) ){
						breakcont = breakcont +1
						 strt =`<p style="margin-top: 20px; border: solid 1px;" ></p>`	
						fmt.Fprint(w,strt)
					}else{
						breakcont = 0 
						strt =`<p style="page-break-before:always;" ></p>`	
						fmt.Fprint(w,strt)
					}

				}
			}
			if Container != val.leader_name{
			//if(counter==0){

				//TotalHours := Interfacetosting(val.reg_reg) + Interfacetosting(val.reg_ot) + Interfacetosting(val.rd_ot) + Interfacetosting(val.sh_ot) + Interfacetosting(val.rh_ot) + Interfacetosting(val.rhrd_ot)

				counterRemarks = + 1
				Container = val.leader_name

				/* tot_reg := Interfacetosting(val.reg_reg)
				tot_reg_ot := Interfacetosting(val.reg_ot)
				total_reg, _ := strconv.ParseInt(tot_reg, 10, 0)
				total_reg_ot, _ := strconv.ParseInt(tot_reg_ot, 10, 0) */

				//fmt.Println(total_reg+total_reg_ot)

				strt=	`	<div class="page">
								<div class="rpt_date">
								
						`
						fmt.Fprint(w,strt)		
						/*fmt.Fprint(w,newFormat)*/
									
				strt=		`	</div>
								<div class="logo">
									
									<b class="logo_name"><center> `
								fmt.Fprint(w,strt)	
								/*fmt.Fprint(w,name)	*/

				strt=			`</center></b>
									<div class="logo_desc">
								`
								fmt.Fprint(w,strt)
								/*fmt.Fprint(w,location)	*/

								
				strt=			`	</div>
								</div>
								

								<div class="print_date"> 
									Productivity Monitoring

								</div>				

								<div class="print_date"> 
									Report Date:`
									fmt.Fprint(w,strt)	
									fmt.Fprint(w,newFormat)	
									strt=`
								</div>			


								<div class="page_header_left"> 

									<table class="dtatbl">
										<tr><td class="tb_r">Date:</td><td style='width:50%; text-align: right; '>`+PrintValue(val.logdate,0)+`</td></tr>
										<tr><td class="tb_r">Project:</td><td style='width:50%; text-align: right; '></td></tr>
										<tr><td class="tb_r">Activity:</td><td style='width:50%; text-align: right; '></td></tr>
										<tr><td class="tb_r">Area/Floor/Block:</td><td style='width:50%; text-align: right; '></td></tr>
										<tr><td class="tb_r">WBS/NWA:</td><td style='width:50%; text-align: right; '></td></tr>
									</table>
								</div>

								<div class="page_header">	
									<div class="name" >Leader Name: <b>  `
										fmt.Fprint(w,strt)
										fmt.Fprint(w,val.leader_name)
										strt=`</b> 

										<table class="dtatbl">
										<tr><td class="tb_r">Accomplishment:</td><td  style='width:50%; text-align: right; '></td></tr>
										<tr><td class="tb_r">No. Of Personnel:</td><td  style='width:50%; text-align: right; '></td></tr>
										<tr><td class="tb_r">Total Manhour:</td><td  style='width:50%; text-align: right; '></td></tr>
										<tr><td class="tb_r">Productivity: (Manhour/Accomplishment):</td><td  style='width:50%; text-align: right; '></td></tr>
										<tr><td  style='width:50%; text-align: right; ' colspan = 2>Productivity = Total Manhour/Accomplishment</td></tr>
									</table>
									</div>
								</div>

							<table class="dtatbl" >
								<tr>
									<th colspan="19">	To be accomplished by COOP COORDINATOR </th>
									<th colspan="2">	Joint Assessment </th>
								</tr>
								<tr>
									
									<th rowspan="3">MEMBER NAME 	</th>
									<th rowspan="3">POSITION 		</th>
									<th colspan="2">ACTUAL TIME 	</th>
									<th colspan="14">MANHOURS		</th>
									<th rowspan="3">	TOTAL HOURS </th>
									<th rowspan="2" colspan="2">	INDIVIDUAL </th>
								</tr>
								<tr>
									<th rowspan="2">TIME IN 		</th>
									<th rowspan="2">	TIME OUT	</th>
									<th colspan="4">	REGULAR 	</th>
									<th colspan="2">	RESTDAY 	</th>
									<th colspan="2">	SPECIAL HOLIDAY	</th>
									<th colspan="2">	SPECIAL HOLIDAY RESTDAY	</th>
									<th colspan="2">	REGULAR HOLIDAY	</th>
									<th colspan="2">	REGULAR HOLIDAY	RESTDAY </th>
								</tr>
									<th>	Reg 	</th>
									<th>	OT 	</th>
									<th>	ND	</th>
									<th>	ND OT	</th>
									<th>	Reg	</th>
									<th>	OT 	</th>
									<th>	Reg 	</th>
									<th>	OT 	</th>
									<th>	Reg	</th>
									<th>	OT	</th>
									
									<th>	Reg	</th>
									<th>	OT	</th>
									<th>	Reg	</th>
									<th>	OT	</th>
									<th>	Accomplishment	</th>
									<th>	Productivity	</th>


								</tr>
								`
								fmt.Fprint(w,strt)

					

			//}
			}  	

			//if(val.Member_name==Container){
			//	counter = counter+1
						if( strings.Contains( PrintValue(val.logdate,0) , "  ") ){
							strt =`
										<tr>


											<td colspan="4" >`+PrintValue(val.logdate,0)+`</td>
											<td colspan="4" >`+PrintValue(val.leader_name,1)+`</td>`
											
						}else{
							strt =`
										<tr>


											<td>`+PrintValue(val.member_name,0)+`</td>
											<td>`+PrintValue(val.lbr_assign_name,0)+`</td>
											<td>`+PrintValue(val.login,1)+`</td>
											<td>`+PrintValue(val.logout,1)+`</td> `

						}

							Total_hours := String_to_float(PrintValue(val.reg_reg ,0)) + String_to_float(PrintValue(val.reg_ot ,0)) + String_to_float(PrintValue(val.rd_ot ,0)) + String_to_float(PrintValue(val.sh_reg ,0)) + String_to_float(PrintValue(val.sh_ot ,0)) + String_to_float(PrintValue(val.rh_reg ,0)) + String_to_float(PrintValue(val.rh_ot ,0))  
							str_Total_hours := Float_to_string(Total_hours)
							strt = strt + `
											<td class="tb_r">`+Round(PrintValue(val.reg_reg ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.reg_ot ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.reg_nd,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.reg_ndot,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.rd_reg,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.rd_ot ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.sh_reg ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.sh_ot ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.shrd_reg,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.shrd_ot,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.rh_reg,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.rh_ot,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.rhrd_reg,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.rhrd_ot,0))+`</td>
											
											<td class="tb_r">`+ str_Total_hours +`</td>
											<td class="tb_r">`+" "+`</td>
											<td class="tb_r">`+" "+`</td>
												
											
											

										</tr>
									`
									fmt.Fprint(w,strt)

		}
		strt = `
				</table>



				



				

				<table  class="tbl_floor" >
					<tbody>
						<tr>
							<td><div>Prepared By</div></td>
							<td><div>Joint Assessment By</div></td>
							<td><div>Verified By</div></td>
						</tr>
					</tbody>
				</table>					
				<table  class="tbl_foor">
					<tbody>
						<tr>
							<td><div>COOP COORDINATOR</div></td>
							<td><div>FIELD ENG'R/FOREMAN/QUANTITY SURVEYOR</div></td>
							<td><div>MDC PHR</div></td>
						</tr>
					</tbody>
				</table>	

			</div>
			<p style="page-break-before:always;"></p>
								
		</body>
		
   
		</html>					
		`
		fmt.Fprint(w,strt)

	

}

func  Report_member_log2(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("here..!")
	login.Session_renew(w,r)
	
	tconf := make(map[string]string)
	tconf["datefrom"] = r.URL.Query().Get("datefrom")
	tconf["dateto"] = r.URL.Query().Get("dateto")
	tconf["lbr_assign"] = r.URL.Query().Get("lbr_assign")
	tconf["member"] = r.URL.Query().Get("member")
	tconf["group"] = r.URL.Query().Get("member")

	//conve ,_ := strconv.ParseInt( tconf["member"] , 10, 32)
	conve,c_err := strconv.Atoi(tconf["member"])
	if(c_err != nil){
		tconf["member"] = "0"
		tconf["group"] = r.URL.Query().Get("member")
	}else{
		tconf["group"] = ""
	}
	fmt.Println(c_err)
	conve_page := int(conve)
	fmt.Println(`conve inte here`)
	fmt.Println(conve_page)


	tconf["no_page"] = r.URL.Query().Get("no_page")
	no_page64 ,_ := strconv.ParseInt( tconf["no_page"] , 10, 32)
	no_page := int(no_page64)
	
	/*type Lbr_Rep_LogRegister struct{
	   		Member_name string
	   		Assignment string
	   		LogDate string 
	   		DayCode string
	   		In_Sched string
	   		Out_Sched string
	   		In_Time string
	   		Out_Time string
	   		Hr_total string
	   		Hr_reg string
	   		Hr_ot string
	   		Hr_nd string
	   		Hr_ndot string
	   		colormered int
		}*/
	/*type Lbr_Rep_LogRegister struct{
	   		Member_name string
	   		Logdate string
	   		Lbr_daycode_name string 
	   		Assignment string
	   		Login string
	   		Logout string
	   		Final_hrreg string
	   		Final_hrot string
	   		Final_hrnd string
	   		Min_late string
	   		Raw_log 		string
	   		
	}*/


	/*var Logdate interface{}
	var Lbr_daycode_name interface{}
	var Login interface{}
	var Logout interface{}
	var Final_hrreg interface{}
	var Final_hrot interface{}
	var Final_hrnd interface{}
	var Min_late  interface{}
	var Raw_log  interface{}*/

	type Lbr_Rep_LogRegister_summary struct{
   		Member_name 	   string
   		Lbr_daycode_name 		interface{}
   		Dayspresent 		interface{} 
   		Total_hrreg 			interface{}
   		Total_hrot 			interface{}
   		Total_hrnd 			interface{}
   		Total_minlate 		interface{}
   	}

   	type Lbr_Rep_LogRegister struct{
   		Member            	string
 		Member_name		  	string	
		Assignment_Name  	interface{}
		Dates				interface{}
		Day					interface{}	
		timein				interface{}
		Timeout				interface{}
		Reg_reg 			interface{}
		Reg_ot 				interface{}
		Reg_nd				interface{}
		Reg_ndot			interface{}
		Reg_late			interface{}
		Reg_undertime		interface{}
		Rd_reg				interface{}
		Rd_ot 				interface{}
		Rd_nd 				interface{}
		Rd_ndot 			interface{}
		Rd_late 			interface{}
		Rd_undertime 		interface{}
		Sh_reg				interface{}
		Sh_ot 				interface{}
		Sh_nd				interface{}
		Sh_ndot				interface{}
		Sh_late				interface{}
		Sh_undertime 		interface{}	
		Shrd_reg	 		interface{}
		Shrd_ot				interface{}
		Shrd_nd				interface{}
		Shrd_ndot	 		interface{}
		Shrd_late			interface{}
		Shrd_undertime		interface{}
		Rh_reg				interface{}
		Rh_ot				interface{}
		Rh_nd				interface{}
		Rh_ndot				interface{}
		Rh_late				interface{}
		Rh_undertime		interface{}
		Rh_100				interface{}
		Rhrd_reg			interface{}
		Rhrd_ot				interface{}
		Rhrd_nd				interface{}
		Rhrd_ndot			interface{}
		Rhrd_late			interface{}
		Rhrd_undertime		interface{}
		RawLogs 			interface{}
   	}

			

	/*type Lbr_Rep_LogRegister struct{
	   		Member_name string
	   		Device_id string
	   		LogDate string 
	   		Day_name string
	   		DayCode string
	   		Assignment string
	   		In_Sched string
	   		Out_Sched string
	   		In_Time string
	   		Out_Time string
	   		Hr_total string
	   		Hr_break string
	   		Hr_reg string
	   		Hr_ot string
	   		Hr_nd string
	   		Hr_ndot string
	   		Min_graceperiod  string
	   		Hr_late string
	   		Hr_undertime string
	   		colormered int
		}
			*/
		//fmt.Println("exec lbr_rep_lograw 3, '"+tconf["datefrom"]+"',  '"+tconf["dateto"]+"', "+tconf["lbr_assign"]+", "+tconf["member"])
	   	//fmt.Println(`lbr_rep_logpair 4,0,`+tconf["member"]+`,`+tconf["lbr_assign"]+`,'`+tconf["datefrom"]+ `','`+tconf["dateto"] + `'`)
	   	//LBR_LogHdr_row ,err, _,_ := config.Ap_sql(`lbr_rep_logpair 4,2,`+tconf["member"]+`,`+tconf["lbr_assign"]+`,'`+tconf["datefrom"]+ `','`+tconf["dateto"] + `'`,1)
	   	Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
	
	   	fmt.Println(`lbr_rep_logpair 4,`+str_OrgID+`,`+tconf["member"]+`,`+tconf["lbr_assign"]+`,'`+tconf["datefrom"]+ `','`+tconf["dateto"] + `','`+tconf["group"]+`'`)
	   	
/*

	   	tconf := make(map[string]string)
		tconf["datefrom"] = r.URL.Query().Get("datefrom")
		tconf["dateto"] = r.URL.Query().Get("dateto")
		tconf["lbr_assign"] = r.URL.Query().Get("lbr_assign")
		tconf["member"] = r.URL.Query().Get("member")
		tconf["no_page"] = r.URL.Query().Get("no_page")
*/


	   	LBR_LogHdr_row ,err, _,_ := config.Ap_sql(`lbr_rep_logpair 4,`+str_OrgID+`,`+tconf["member"]+`,`+tconf["lbr_assign"]+`,'`+tconf["datefrom"]+ `','`+tconf["dateto"] + `','`+tconf["group"]+`'`,1)
	   	//LBR_LogHdr_row ,err, _,_ := config.Ap_sql(`lbr_rep_logpair 4, 2, 1, 0, '1 Nov 2015', '15 Nov 2015'`,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    Lbr_Rep_LogRegister_data := []Lbr_Rep_LogRegister{}
	    

	    for LBR_LogHdr_row.Next() {
            var r Lbr_Rep_LogRegister
            //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
            //err = LBR_LogHdr_row.Scan(&r.Member_name, &Logdate,&Lbr_daycode_name,&r.Assignment,&Login,&Logout,&Final_hrreg,&Final_hrot,&Final_hrnd,&Min_late,&Raw_log)
            err = LBR_LogHdr_row.Scan(
									&r.Member            	,
									&r.Member_name	,
									&r.Assignment_Name  	,
									&r.Dates	,
									&r.Day	,
									&r.timein	,
									&r.Timeout	,
									&r.Reg_reg 	,
									&r.Reg_ot 	,
									&r.Reg_nd	,
									&r.Reg_ndot	,
									&r.Reg_late	,
									&r.Reg_undertime	,
									&r.Rd_reg	,
									&r.Rd_ot 	,
									&r.Rd_nd 	,
									&r.Rd_ndot 	,
									&r.Rd_late 	,
									&r.Rd_undertime 	,
									&r.Sh_reg 	,
									&r.Sh_ot 	,
									&r.Sh_nd	,
									&r.Sh_ndot	,
									&r.Sh_late	,
									&r.Sh_undertime 	,
									&r.Shrd_reg	,
									&r.Shrd_ot	,
									&r.Shrd_nd	,
									&r.Shrd_ndot	,
									&r.Shrd_late	,
									&r.Shrd_undertime	,
									&r.Rh_reg	,
									&r.Rh_ot	,
									&r.Rh_nd	,
									&r.Rh_ndot	,
									&r.Rh_late	,
									&r.Rh_undertime	,
									&r.Rh_100	,
									&r.Rhrd_reg	,
									&r.Rhrd_ot	,
									&r.Rhrd_nd	,
									&r.Rhrd_ndot	,
									&r.Rhrd_late	,
									&r.Rhrd_undertime,
									&r.RawLogs	)
            
            if err != nil {
                    panic(err.Error())
            }
			/*post2 := Lbr_Rep_LogRegister{
            	r.Member_name,
		        PrintValue(Logdate,0),
		        PrintValue(Lbr_daycode_name,0),
		      	r.Assignment,
		      	PrintValue(Login,1),
		      	PrintValue(Logout,1),
		      	PrintValue(Final_hrreg,0),
		      	PrintValue(Final_hrot,0),
		      	PrintValue(Final_hrnd,0),
		      	PrintValue(Min_late,0),
		      	PrintValue(Raw_log,0),
		      	
 		    }*/

 		    post2 := Lbr_Rep_LogRegister{
            	r.Member            	,
				r.Member_name			,
				r.Assignment_Name  	,
				r.Dates				,
				r.Day					,
				r.timein				,
				r.Timeout				,
				r.Reg_reg 	,
				r.Reg_ot 	,
				r.Reg_nd	,
				r.Reg_ndot	,
				r.Reg_late	,
				r.Reg_undertime	,
				r.Rd_reg	,
				r.Rd_ot 	,
				r.Rd_nd 	,
				r.Rd_ndot 	,
				r.Rd_late 	,
				r.Rd_undertime 	,
				r.Sh_reg 	,
				r.Sh_ot 	,
				r.Sh_nd	,
				r.Sh_ndot	,
				r.Sh_late	,
				r.Sh_undertime 	,
				r.Shrd_reg	,
				r.Shrd_ot	,
				r.Shrd_nd	,
				r.Shrd_ndot	,
				r.Shrd_late	,
				r.Shrd_undertime	,
				r.Rh_reg	,
				r.Rh_ot	,
				r.Rh_nd	,
				r.Rh_ndot	,
				r.Rh_late	,
				r.Rh_undertime	,
				r.Rh_100	,
				r.Rhrd_reg	,
				r.Rhrd_ot	,
				r.Rhrd_nd	,
				r.Rhrd_ndot	,
				r.Rhrd_late	,
				r.Rhrd_undertime	,
				r.RawLogs ,

		      	
 		    }

			Lbr_Rep_LogRegister_data = append(Lbr_Rep_LogRegister_data, post2)
	    }


		
		t := time.Now().Local()
		//newFormat := t.Format("2006-01-02 15:00:00 +0800")
		newFormat := t.Format("Mon Jan _2 2006 03:04:05PM")
		        

		var id int
		var name string
		var location string
		Org_Get_row ,err, _,_ := config.Ap_sql(" Org_Get 1, 2",1)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    for Org_Get_row.Next() {
	    	err = Org_Get_row.Scan(&id, &name,&location)
            if err != nil {
                    panic(err.Error())
            }
        }
	    
        
	    ///version1
	    strt:= `
	    <html>
			<head>
				<title style="font-size: 16px;">Time keeping report</title>
			</head>
			<style type="text/css">
				.rpt_date {
					float: right;
					
				}
				.page_header{
					margin-bottom: 5px;
					    float: right;
					/*display: inline; */
				}
				.name{
					/*float: left;*/
					margin-bottom: 3px;
					font-family: Arial; 
				}
				.covered_data{
					/*float: right;*/
					font-family: Arial; 
				}
				/*.pic_logo {
					width: 200px;

				}*/
				table{
					width: 100%;
				}
				.tb_r{
					text-align: right;
				}
				
				.dtatbl {
				
				border: 0.99px solid #B3B3B3;
				border-collapse: collapse;
				font-family: 'Arial Narrow', Arial, sans-serif;
				font-size: 10px;
				}

				.dtatbl th{
					border: 0.99px solid #B3B3B3;
				border-collapse: collapse;
				font-family: 'Arial Narrow', Arial, sans-serif;
				font-size: 10px;
				}

				.dtatbl td {
					border: 0.99px solid #B3B3B3;
				border-collapse: collapse;
				font-family: 'Arial Narrow', Arial, sans-serif;
				font-size: 10px;
				}
				th, td {
				    /*padding: 2px;*/
				}
				th{
					
					background-color: whiteSmoke;
					/*background-image: -webkit-linear-gradient*/
					
				}
				body {
					font-size: 10px;
				}
				.logo{
					width: 200px;
					margin-bottom: 5px;
				}
				.page_footer{
					margin-top: 20px;
					display: inline-block;
					width: 100%;
				}
				.signaturemember{
					border-top: 0.99px solid #B3B3B3;
					width: 25%;
					float: left;
					font-family: inherit;
				}
				.signatureleft{
					margin-left: 10%;
					border-top: 0.99px solid #B3B3B3;
					width: 25%;
					float: left;
					font-family: inherit; 
				}
				.signatureright{
					border-top: 0.99px solid #B3B3B3;
					width: 25%;
					float: right;
					font-family: inherit; 
				}
				.tbl_foor{
					margin-top: 10px;
				}
				.tbl_foor td div{
					border-top: 0.99px solid #B3B3B3;	
					width: 50%;
					font-size: 10px;
				}
				.isred {
					color: red;
				}
				.log_sum{
					margin-top: 10px;
				}

			</style>

			<body>
		`
		fmt.Fprint(w,strt)
		//counter :=0
		//counter1:=0
		Container :=""
		counterRemarks :=0
		//Contain := ""
		var breakcont int = 0
		for _ , val := range Lbr_Rep_LogRegister_data {
		/*	if counter1==0{
				Contain = val.Member_name
			}
			if val.Member_name==Contain{
				fmt.Fprint(w,`---******` + val.Member_name)
			}

			*/

			if counterRemarks!=0{
				if Container != val.Member_name{
					
							strt=`
							</table>
							



										
									
					</div>
					

					

					<table  class="tbl_foor">
						<tbody>
							<tr>
								<td><div>Member Signature</div></td>
								<td><div>Verified By</div></td>
								<td><div>Noted By</div></td>
							</tr>
						</tbody>
					</table>	


							
							`
							fmt.Fprint(w,strt)
					if(breakcont< (no_page-1) ){
						breakcont = breakcont +1
						 strt =`<p style="margin-top: 20px; border: solid 1px;" ></p>`	
						fmt.Fprint(w,strt)
					}else{
						breakcont = 0 
						strt =`<p style="page-break-before:always;" ></p>`	
						fmt.Fprint(w,strt)
					}

				}
			}
			if Container != val.Member_name{
			//if(counter==0){
				counterRemarks = + 1
				Container = val.Member_name 
				strt=	`	<div class="page">
								<div class="rpt_date">
								
						`
						fmt.Fprint(w,strt)		
						/*fmt.Fprint(w,newFormat)*/
									
				strt=		`	</div>
								<div class="logo">
									
									<b class="logo_name"><center> `
								fmt.Fprint(w,strt)	
								/*fmt.Fprint(w,name)	*/

				strt=			`</center></b>
									<div class="logo_desc">
								`
								fmt.Fprint(w,strt)
								/*fmt.Fprint(w,location)	*/

								
				strt=			`	</div>
								</div>
								

								<div class="page_header">	
									<div class="name" >Name: <b>  `
										fmt.Fprint(w,strt)
										fmt.Fprint(w,val.Member_name)
										strt=`</b> 
									</div>
									<div class="name" >Position:   `
										fmt.Fprint(w,strt)
										fmt.Fprint(w,val.Assignment_Name)
										strt=`
									</div>
								</div>

								<div class="print_date"> 
									Time keeping register
									
								</div>		

								<div class="print_date"> 
									Report Date:`
									fmt.Fprint(w,strt)	
									fmt.Fprint(w,newFormat)	
									strt=`
								</div>			


							<table class="dtatbl" >
								<tr>
									
									<th rowspan="2">DATE</th>
									<th rowspan="2">DAY</th>
									<th rowspan="2">TIME IN</th>
									<th rowspan="2">	TIME OUT	</th>
									<th colspan="6">	REGULAR 	</th>
									<th colspan="6">	RESTDAY 	</th>
									<th colspan="6">	SPECIAL HOLIDAY	</th>
									<th colspan="6">	SPECIAL HOLIDAY RESTDAY	</th>
									
									
									
									
									<th colspan="7">	REGULAR HOLIDAY	</th>
									<th colspan="6">	REGULAR HOLIDAY	RESTDAY </th>
									<th rowspan ="2"> RAWLOGS </th>
								</tr>
								<tr>
									<th>	Reg 	</th>
									<th>	OT 	</th>
									<th>	ND	</th>
									<th>	ND OT	</th>
									<th>	Late	</th>
									<th>	U/Time	</th>
									<th>	Reg	</th>
									<th>	OT 	</th>
									<th>	ND 	</th>
									<th>	NDOT 	</th>
									<th>	Late 	</th>
									<th>	U/Time 	</th>
									<th>	Reg 	</th>
									<th>	OT 	</th>
									<th>	ND	</th>
									<th>	NDOT	</th>
									<th>	Late	</th>
									<th>	U/Time 	</th>
									<th>	Reg	</th>
									<th>	OT	</th>
									<th>	ND	</th>
									<th>	NDOT	</th>
									<th>	Late	</th>
									<th>	U/Time	</th>
									
									<th>	Reg	</th>
									<th>	OT	</th>
									<th>	ND	</th>
									<th>	NDOT	</th>
									<th>	Late	</th>
									<th>	U/Time	</th>
									<th>	100%	</th>
									<th>	Reg	</th>
									<th>	OT	</th>
									<th>	ND	</th>
									<th>	NDOT	</th>
									<th>	Late	</th>
									<th>	U/Time	</th>


								</tr>
								`
								fmt.Fprint(w,strt)

					

			//}
			}  	

			//if(val.Member_name==Container){
			//	counter = counter+1
						if( strings.Contains( PrintValue(val.Dates,0) , "**") ){
							strt =`
										<tr>


											<td colspan="4" >`+PrintValue(val.Dates,0)+`</td>`
											
						}else{
							strt =`
										<tr>


											<td>`+PrintValue(val.Dates,0)+`</td>
											<td>`+PrintValue(val.Day,0)+`</td>
											<td>`+PrintValue(val.timein,1)+`</td>
											<td>`+PrintValue(val.Timeout,1)+`</td> `

						}

							strt = strt + `
											<td class="tb_r">`+Round(PrintValue(val.Reg_reg ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Reg_ot ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Reg_nd,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Reg_ndot,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Reg_late,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Reg_undertime,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rd_reg,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rd_ot ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rd_nd ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rd_ndot ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rd_late ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rd_undertime ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Sh_reg ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Sh_ot ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Sh_nd,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Sh_ndot,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Sh_late,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Sh_undertime ,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Shrd_reg,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Shrd_ot,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Shrd_nd,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Shrd_ndot,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Shrd_late,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Shrd_undertime,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rh_reg,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rh_ot,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rh_nd,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rh_ndot,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rh_late,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rh_undertime,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rh_100,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rhrd_reg,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rhrd_ot,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rhrd_ot,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rhrd_ndot,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rhrd_late,0))+`</td>
											<td class="tb_r">`+Round(PrintValue(val.Rhrd_undertime,0))+`</td>
											<td>`+PrintValue(val.RawLogs,0)+`</td>
												
											
											

										</tr>
									`
									fmt.Fprint(w,strt)
		}
		strt = `
				</table>



				



				

				<table  class="tbl_foor">
					<tbody>
						<tr>
							<td><div>Member Signature</div></td>
							<td><div>Verified By</div></td>
							<td><div>Noted By</div></td>
						</tr>
					</tbody>
				</table>	

			</div>
			<p style="page-break-before:always;"></p>
								
		</body>
		
   
		</html>					
		`
		fmt.Fprint(w,strt)
		

	

}



func  Report_member_log(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("here..!")
	login.Session_renew(w,r)
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
		

	tconf := make(map[string]string)
	tconf["datefrom"] = r.URL.Query().Get("datefrom")
	tconf["dateto"] = r.URL.Query().Get("dateto")
	tconf["lbr_assign"] = r.URL.Query().Get("lbr_assign")
	tconf["member"] = r.URL.Query().Get("member")
	tconf["no_page"] = r.URL.Query().Get("no_page")
	no_page64 ,_ := strconv.ParseInt( tconf["no_page"] , 10, 32)
	no_page := int(no_page64)
	/*type Lbr_Rep_LogRegister struct{
	   		Member_name string
	   		Assignment string
	   		LogDate string 
	   		DayCode string
	   		In_Sched string
	   		Out_Sched string
	   		In_Time string
	   		Out_Time string
	   		Hr_total string
	   		Hr_reg string
	   		Hr_ot string
	   		Hr_nd string
	   		Hr_ndot string
	   		colormered int
		}*/
	type Lbr_Rep_LogRegister struct{
	   		Member_name string
	   		Logdate string
	   		Lbr_daycode_name string 
	   		Assignment string
	   		Login string
	   		Logout string
	   		Final_hrreg string
	   		Final_hrot string
	   		Final_hrnd string
	   		Min_late string
	   		Raw_log 		string
	   		
		}
	var Logdate interface{}
	var Lbr_daycode_name interface{}
	var Login interface{}
	var Logout interface{}
	var Final_hrreg interface{}
	var Final_hrot interface{}
	var Final_hrnd interface{}
	var Min_late  interface{}
	var Raw_log  interface{}

	type Lbr_Rep_LogRegister_summary struct{
   		Member_name 	   string
   		Lbr_daycode_name 		interface{}
   		Dayspresent 		interface{} 
   		Total_hrreg 			interface{}
   		Total_hrot 			interface{}
   		Total_hrnd 			interface{}
   		Total_minlate 		interface{}
   	}
			

	/*type Lbr_Rep_LogRegister struct{
	   		Member_name string
	   		Device_id string
	   		LogDate string 
	   		Day_name string
	   		DayCode string
	   		Assignment string
	   		In_Sched string
	   		Out_Sched string
	   		In_Time string
	   		Out_Time string
	   		Hr_total string
	   		Hr_break string
	   		Hr_reg string
	   		Hr_ot string
	   		Hr_nd string
	   		Hr_ndot string
	   		Min_graceperiod  string
	   		Hr_late string
	   		Hr_undertime string
	   		colormered int
		}
			*/
		//fmt.Println("exec lbr_rep_lograw 3, '"+tconf["datefrom"]+"',  '"+tconf["dateto"]+"', "+tconf["lbr_assign"]+", "+tconf["member"])
	   	fmt.Println(`lbr_rep_lograw 5,`+str_OrgID+`,`+tconf["member"]+`,`+tconf["lbr_assign"]+`,'',`+`'','`+tconf["datefrom"]+ `','`+tconf["dateto"] + `'`)
	   	LBR_LogHdr_row ,err, _,_ := config.Ap_sql(`lbr_rep_lograw 5,`+str_OrgID+`,`+tconf["member"]+`,`+tconf["lbr_assign"]+`,'',`+`'','`+tconf["datefrom"]+ `','`+tconf["dateto"] + `'`,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    Lbr_Rep_LogRegister_data := []Lbr_Rep_LogRegister{}
	    

	    for LBR_LogHdr_row.Next() {
            var r Lbr_Rep_LogRegister
            //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
            err = LBR_LogHdr_row.Scan(&r.Member_name, &Logdate,&Lbr_daycode_name,&r.Assignment,&Login,&Logout,&Final_hrreg,&Final_hrot,&Final_hrnd,&Min_late,&Raw_log)
            if err != nil {
                    panic(err.Error())
            }
			post2 := Lbr_Rep_LogRegister{
            	r.Member_name,
		        PrintValue(Logdate,0),
		        PrintValue(Lbr_daycode_name,0),
		      	r.Assignment,
		      	PrintValue(Login,1),
		      	PrintValue(Logout,1),
		      	PrintValue(Final_hrreg,0),
		      	PrintValue(Final_hrot,0),
		      	PrintValue(Final_hrnd,0),
		      	PrintValue(Min_late,0),
		      	PrintValue(Raw_log,0),
		      	
 		    }
			Lbr_Rep_LogRegister_data = append(Lbr_Rep_LogRegister_data, post2)
	    }
		
		t := time.Now().Local()
		//newFormat := t.Format("2006-01-02 15:00:00 +0800")
		newFormat := t.Format("Mon Jan _2 2006")
		        

		var id int
		var name string
		var location string
		Org_Get_row ,err, _,_ := config.Ap_sql(" Org_Get 1, 2",1)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    for Org_Get_row.Next() {
	    	err = Org_Get_row.Scan(&id, &name,&location)
            if err != nil {
                    panic(err.Error())
            }
        }
	    
        //start of log summary

        fmt.Println(`lbr_rep_lograw 6,`+str_OrgID+`,`+tconf["member"]+`,`+tconf["lbr_assign"]+`,'',`+`'','`+tconf["datefrom"]+ `','`+tconf["dateto"] + `'`)

        LBR_LogHdr_sum_row ,err, _,_ := config.Ap_sql(`lbr_rep_lograw 6,`+str_OrgID+`,`+tconf["member"]+`,`+tconf["lbr_assign"]+`,'',`+`'','`+tconf["datefrom"]+ `','`+tconf["dateto"] + `'`,1)
	    if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	panic(err.Error())
	    }
	    Lbr_Rep_LogRegister_summary_data := []Lbr_Rep_LogRegister_summary{}
	    for LBR_LogHdr_sum_row.Next() {
            var r Lbr_Rep_LogRegister_summary
            //err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
            err = LBR_LogHdr_sum_row.Scan(&r.Member_name, &r.Lbr_daycode_name,&r.Dayspresent ,&r.Total_hrreg ,&r.Total_hrot, &r.Total_hrnd ,&r.Total_minlate)
            if err != nil {
                    panic(err.Error())
            }
			post2_sum := Lbr_Rep_LogRegister_summary{
            	r.Member_name 	,
		   		r.Lbr_daycode_name 		,
		   		r.Dayspresent 	,	 
		   		r.Total_hrreg 		,	
		   		r.Total_hrot 		,	
		   		r.Total_hrnd   ,	
		   		r.Total_minlate   ,	
		   		 	
 		    }
			Lbr_Rep_LogRegister_summary_data = append(Lbr_Rep_LogRegister_summary_data, post2_sum)
	    }
		
	    //end of log summary
	    ///version1
	    strt:= `
	    <html>
			<head>
				<title style="font-size: 16px;">Time keeping report</title>
			</head>
			<style type="text/css">
				.rpt_date {
					float: right;
					
				}
				.page_header{
					margin-bottom: 5px;
					/*display: inline; */
				}
				.name{
					/*float: left;*/
					margin-bottom: 3px;
					font-family: Arial; 
				}
				.covered_data{
					/*float: right;*/
					font-family: Arial; 
				}
				/*.pic_logo {
					width: 200px;

				}*/
				table{
					width: 100%;
				}
				
				.dtatbl {
				
				border: 0.99px solid #B3B3B3;
				border-collapse: collapse;
				font-family: 'Arial Narrow', Arial, sans-serif;
				font-size: 9px;
				}

				.dtatbl th{
					border: 0.99px solid #B3B3B3;
				border-collapse: collapse;
				font-family: 'Arial Narrow', Arial, sans-serif;
				font-size: 9px;
				}

				.dtatbl td {
					border: 0.99px solid #B3B3B3;
				border-collapse: collapse;
				font-family: 'Arial Narrow', Arial, sans-serif;
				font-size: 9px;
				}
				th, td {
				    /*padding: 2px;*/
				}
				th{
					
					background-color: whiteSmoke;
					/*background-image: -webkit-linear-gradient*/
					
				}
				body {
					font-size: 9px;
				}
				.logo{
					width: 200px;
					margin-bottom: 5px;
				}
				.page_footer{
					margin-top: 20px;
					display: inline-block;
					width: 100%;
				}
				.signaturemember{
					border-top: 0.99px solid #B3B3B3;
					width: 25%;
					float: left;
					font-family: inherit;
				}
				.signatureleft{
					margin-left: 10%;
					border-top: 0.99px solid #B3B3B3;
					width: 25%;
					float: left;
					font-family: inherit; 
				}
				.signatureright{
					border-top: 0.99px solid #B3B3B3;
					width: 25%;
					float: right;
					font-family: inherit; 
				}
				.tbl_foor{
					margin-top: 9px;
				}
				.tbl_foor td div{
					border-top: 0.99px solid #B3B3B3;	
					width: 50%;
					font-size: 9px;
				}
				.isred {
					color: red;
				}
				.log_sum{
					margin-top: 9px;
				}

			</style>

			<body>
		`
		fmt.Fprint(w,strt)
		//counter :=0
		//counter1:=0
		Container :=""
		counterRemarks :=0
		//Contain := ""
		var breakcont int = 0
		for _ , val := range Lbr_Rep_LogRegister_data {
		/*	if counter1==0{
				Contain = val.Member_name
			}
			if val.Member_name==Contain{
				fmt.Fprint(w,`---******` + val.Member_name)
			}

			*/

			if counterRemarks!=0{
				if Container != val.Member_name{
					
							strt=`
							</table>
							<div class="log_sum">
					<b>Log Summary:</b>
					<table style="width:50%;" class="dtatbl">
						<tbody>
							<tr>
								<th>DayCode</th>
								<th>Day Present</th>
								<th>Hr Reg</th>
								<th>Hr OT</th>
								<th>Hr ND</th>
								<th>Min Late</th>


							</tr>`
							fmt.Fprint(w,strt)

							for  _, val_sum := range Lbr_Rep_LogRegister_summary_data {
								if val_sum.Member_name == Container{
									strt =`
										<tr>
											<td>`+PrintValue(val_sum.Lbr_daycode_name,0)+`</td>
											<td>`+PrintValue(val_sum.Dayspresent,0)+`</td>
											<td>`+PrintValue(val_sum.Total_hrreg,0)+`</td>
											<td>`+PrintValue(val_sum.Total_hrot,0)+`</td>
											<td>`+PrintValue(val_sum.Total_hrnd,0)+`</td>
											<td>`+PrintValue(val_sum.Total_minlate,0)+`</td>
											

										</tr>
									`
									fmt.Fprint(w,strt)
								}
							}
							strt =`

						</tbody>
					</table>
				</div> <!-- end logsum -->





										
									
					</div>
					

					

					<table  class="tbl_foor">
						<tbody>
							<tr>
								<td><div>Member Signature</div></td>
								<td><div>Verified By</div></td>
								<td><div>Noted By</div></td>
							</tr>
						</tbody>
					</table>	


							
							`
							fmt.Fprint(w,strt)
					if(breakcont< (no_page-1) ){
						breakcont = breakcont +1
						 strt =`<p style="margin-top: 20px; border: solid 1px;" ></p>`	
						fmt.Fprint(w,strt)
					}else{
						breakcont = 0 
						strt =`<p style="page-break-before:always;" ></p>`	
						fmt.Fprint(w,strt)
					}

				}
			}
			if Container != val.Member_name{
			//if(counter==0){
				counterRemarks = + 1
				Container = val.Member_name 
				strt=	`	<div class="page">
								<div class="rpt_date">
								Report Date:
						`
						fmt.Fprint(w,strt)		
						fmt.Fprint(w,newFormat)
									
				strt=		`	</div>
								<div class="logo">
									
									<b class="logo_name"><center> `
								fmt.Fprint(w,strt)	
								//fmt.Fprint(w,name)	///orgname here

				strt=			`</center></b>
									<div class="logo_desc">
								`
								fmt.Fprint(w,strt)
								fmt.Fprint(w,location)	

								
				strt=			`	</div>
								</div>
								

								<div class="page_header">	
									<div class="name" >Name: <b>  `
										fmt.Fprint(w,strt)
										fmt.Fprint(w,val.Member_name)
										strt=`</b> 
									</div>
							
							
									<div class="covered_data"> 
										Coverd Date:`
										fmt.Fprint(w,strt)	
										fmt.Fprint(w,tconf["datefrom"])
											strt=`
												  to `
											fmt.Fprint(w,strt)
										fmt.Fprint(w,tconf["dateto"])	
										strt=`
									</div>
								</div>			

							<table class="dtatbl" >
								<tr>
									
								    <th>Date</th> 
								    <th>Day Code</th>
								    <th>Asgmnt</th> 
								    <th>Login</th>
								    <th>Logout</th>
								    <th>Reg hrs</th> 
								    <th>OT</th>
								    <th>ND</th>
								    <th>Late mins</th>
								    <th>Raw Log</th>
								    
								</tr>
								`
								fmt.Fprint(w,strt)

					

			//}
			}  	

			//if(val.Member_name==Container){
			//	counter = counter+1
					
					strt = `<tr> 
						<td>`
				

					//fmt.Fprint(w,strt)
					//fmt.Fprint(w,val.Device_id)
					//strt = `</td>`			
					//fmt.Fprint(w,strt)

					strt = `<td>`
					fmt.Fprint(w,strt)
					fmt.Fprint(w,val.Logdate)
					strt = `</td>`			
					fmt.Fprint(w,strt)

					strt = `<td>`
					fmt.Fprint(w,strt)
					fmt.Fprint(w,val.Lbr_daycode_name)
					strt = `</td>`			
					fmt.Fprint(w,strt)

					strt = `<td style="font-size: 6px;">`
					fmt.Fprint(w,strt)
					fmt.Fprint(w,val.Assignment)
					strt = `</td>`			
					fmt.Fprint(w,strt)

					strt = `<td>`
					fmt.Fprint(w,strt)
					fmt.Fprint(w,val.Login)
					strt = `</td>`			
					fmt.Fprint(w,strt)

					strt = `<td>`
					fmt.Fprint(w,strt)
					fmt.Fprint(w,val.Logout)
					strt = `</td>`			
					fmt.Fprint(w,strt)

		      		strt = `<td>`
					fmt.Fprint(w,strt)
					fmt.Fprint(w,val.Final_hrreg)
					strt = `</td>`			
					fmt.Fprint(w,strt)

					strt = `	<td>`
					fmt.Fprint(w,strt)
					fmt.Fprint(w,val.Final_hrot)
					strt = `</td>`			
					fmt.Fprint(w,strt)

					strt = `<td>`
					fmt.Fprint(w,strt)
					fmt.Fprint(w,val.Final_hrnd)
					strt = `</td>`			
					fmt.Fprint(w,strt)

					strt = `<td>`
					fmt.Fprint(w,strt)
					fmt.Fprint(w,val.Min_late)
					strt = `</td>`			
					fmt.Fprint(w,strt)

					strt = `<td style="font-size: 6px;">`
					fmt.Fprint(w,strt)
					fmt.Fprint(w,val.Raw_log)  //Break hours applied
					strt = `</td>`			
					fmt.Fprint(w,strt)


					

					strt = `
					</tr>`		
					fmt.Fprint(w,strt)

		
			
	    		
		}
		strt = `
				</table>



				<div class="log_sum">
					<b>Log Summary:</b>
					<table style="width:50%;" class="dtatbl">
						<tbody>
							<tr>
								<th>DayCode</th>
								<th>Day Present</th>
								<th>Hr Reg</th>
								<th>Hr OT</th>
								<th>Hr ND</th>
								<th>Min late</th>
								


							</tr>
							`
							fmt.Fprint(w,strt)

							for  _, val_sum := range Lbr_Rep_LogRegister_summary_data {
								if val_sum.Member_name == Container{
									strt =`
										<tr>
											<td>`+PrintValue(val_sum.Lbr_daycode_name,0)+`</td>
											<td>`+PrintValue(val_sum.Dayspresent,0)+`</td>
											<td>`+PrintValue(val_sum.Total_hrreg,0)+`</td>
											<td>`+PrintValue(val_sum.Total_hrot,0)+`</td>
											<td>`+PrintValue(val_sum.Total_hrnd,0)+`</td>
											<td>`+PrintValue(val_sum.Total_minlate,0)+`</td>
											

										</tr>
									`
									fmt.Fprint(w,strt)
								}
							}

						strt =	`

						</tbody>
					</table>
				</div> <!-- end logsum -->



				

				<table  class="tbl_foor">
					<tbody>
						<tr>
							<td><div>Member Signature</div></td>
							<td><div>Verified By</div></td>
							<td><div>Noted By</div></td>
						</tr>
					</tbody>
				</table>	

			</div>
			<p style="page-break-before:always;"></p>
								
		</body>
		
   
		</html>					
		`
		fmt.Fprint(w,strt)
		

	

}

func Round(input_str string) string {
	//fmt.Println(" round string" , input_str) 	
 	input ,_ := strconv.ParseFloat(input_str,64)
 	//fmt.Println(" plain float64" , input) 
 	if (input ==0){
 		a := fmt.Sprint("")
 		return a
 	}
 	a := fmt.Sprint(input)
 	return a
 	/*//input float64

 	if input < 0 {
 		input_float64:= math.Ceil(input - 0.5)
 		a := fmt.Sprint(input_float64)
 		fmt.Println(" round float64" , a) 	
 		return a
 	}
 	input_float64:= math.Floor(input + 0.5)

 	a := fmt.Sprint(input_float64)
 	fmt.Println(" round float64" , a) 	
 	return a*/
 }

func String_to_float(input_str string) float64 {
	input ,_ := strconv.ParseFloat(input_str,64)
 	return input
}
func Float_to_string(input_str float64) string {
	a := fmt.Sprint(input_str)
 	return a
}

func PrintValue(pval interface{},istime int) string{
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
		if (istime==1){
			a = fmt.Sprint(v.Format("15:04:05"))
		}else if(istime==0){
			a = fmt.Sprint(v.Format("2006-01-02"))
		}else{
			//fmt.Println(v.Month())
			a = fmt.Sprint(v.Format("2006-01-02 Mon"))

		}

	default:
//		fmt.Print(v)
		a= fmt.Sprint(v)
	}
	return a
}