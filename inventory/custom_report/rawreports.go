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
		helper "../../helper"
		//"time"
		"time"
		"strings"
		datatables "../../datatables"

)
type location_name struct{
	Id int
	Name string
}

type Supplier struct{
	Id int
	Name string
}
type Customer struct{
	Id int
	Name string
	Courier int
	Soldby int
}

func func_Location_Get(  session_user_id_str string, w http.ResponseWriter, r *http.Request) []location_name {
	
	
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
	row ,err, _,_ := config.Ap_sql("exec SIS_Location_Get 6,"+str_OrgID+",''",1)
	res := []location_name{}
	for row.Next() {
        var r location_name
        err = row.Scan(&r.Id, &r.Name)
        if err != nil {
                panic(err.Error())
        }
		datarow := location_name{
        	r.Id,
	        r.Name,
	      
	    }
		res = append(res, datarow)

    }
    return res		   		
}

func func_customer_Get(  session_user_id_str string, w http.ResponseWriter, r *http.Request) []Customer {
	
	
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
	row ,err, _,_ := config.Ap_sql("exec SIS_Customer_Get 2,"+str_OrgID+",''",1)
	res := []Customer{}
	for row.Next() {
        var r Customer
        err = row.Scan(&r.Id, &r.Name, &r.Courier, &r.Soldby)
        if err != nil {
                panic(err.Error())
        }
		datarow := Customer{
        	r.Id,
	        r.Name,
	         r.Courier,
	      	r.Soldby,
	      
	    }
		res = append(res, datarow)

    }
    return res		   		
}


func func_Supplier_Get(  session_user_id_str string, w http.ResponseWriter, r *http.Request) []Supplier {
	
	
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	row ,err, _,_ := config.Ap_sql("exec  SIS_Supplier_Get 2,"+str_OrgID+",''",1)
	res := []Supplier{}
	for row.Next() {
        var r Supplier
        err = row.Scan(&r.Id, &r.Name )
        if err != nil {
                panic(err.Error())
        }
		datarow := Supplier{
        	r.Id,
	        r.Name,
	       
	    }
		res = append(res, datarow)

    }
    return res		   		
}



func  Rawreports_nav(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	_,session_user_id := login.Get_account_info(r)
	session_user_id_str := strconv.Itoa(session_user_id) 
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)

	type Htdata struct {
		Location_Get []location_name
		Supplier []Supplier
		CustomerList []Customer
		Report_list [][]string
	}
	report_list := datatables.DataList(`SIS_CustomReport_List 1,`+str_OrgID )
	var_htdata :=Htdata{func_Location_Get(session_user_id_str,w,r),func_Supplier_Get(session_user_id_str,w,r),func_customer_Get(session_user_id_str,w,r),report_list}
	tmpl := template.New("rawreports_nav.html")
	var err error
	if tmpl, err = tmpl.ParseFiles("inventory/custom_report/templates/rawreports_nav.html"); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w,var_htdata)
	if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
	}

}



func Rawreports_doprint(w http.ResponseWriter, r *http.Request) {
	tconf := make(map[string]string)
	tconf["name"] = r.URL.Query().Get("reportname")
	tconf["datefrom"] = r.URL.Query().Get("datefrom")
	tconf["dateto"] = r.URL.Query().Get("dateto")
	tconf["location"] = r.URL.Query().Get("location")
	fmt.Println(tconf["location"])
	tconf["export_type"] = r.URL.Query().Get("export_type")
	tconf["supplier"] = r.URL.Query().Get("supplier")
	tconf["trantype"] = r.URL.Query().Get("trantype")
	tconf["trn_id"] = strings.Trim(r.URL.Query().Get("trn_id"), ` `)

	tconf["soldby"] = r.URL.Query().Get("soldby")
	tconf["customer"] = r.URL.Query().Get("customer")
	tconf["transferto"] = r.URL.Query().Get("transferto")
	tconf["transferfrom"] = r.URL.Query().Get("transferfrom")



	if (len(tconf["trn_id"]) ==0){
		tconf["trn_id"] = `0`
	}
	if (len(tconf["soldby"]) ==0){
		tconf["soldby"] = `0`
	}
	//tconf["orgid"]="2"
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	fmt.Println("sp run here.............")
	fmt.Println(tconf["name"])
	
	//strt := tconf["reportname"]

	//lbr_rep_logpair 3, 2, 0 , 46, '21 Oct 2015', '3 Nov 2015'
	arr_reportname := strings.Split(string(tconf["name"]), "@@")
	fsqlstr:=arr_reportname[0]
	for i := 1; i < len(arr_reportname); i++ {  
		str_added :=""
		
		if( strings.Contains(arr_reportname[i], "id") ){
			str_added =  strings.Replace(arr_reportname[i],"id", tconf["trn_id"] ,-1)	
			fsqlstr = fsqlstr + str_added
		}

		if( strings.Contains(arr_reportname[i], "org") ){
			str_added =  strings.Replace(arr_reportname[i],"org", str_OrgID,-1)	
			fsqlstr = fsqlstr + str_added
		}
	
		if( strings.Contains(arr_reportname[i], "location") ){
			str_added =  strings.Replace(arr_reportname[i],"location", tconf["location"],-1)	
			fsqlstr = fsqlstr + str_added
		}
		if( strings.Contains(arr_reportname[i], "datefrom") ){
			str_added =  strings.Replace(arr_reportname[i],"datefrom", tconf["datefrom"],-1)	
			fsqlstr = fsqlstr + str_added
		}
		if( strings.Contains(arr_reportname[i], "dateto") ){
			str_added =  strings.Replace(arr_reportname[i],"dateto", tconf["dateto"],-1)	
			fsqlstr = fsqlstr + str_added
		}
		if( strings.Contains(arr_reportname[i], "supplier") ){
			str_added =  strings.Replace(arr_reportname[i],"supplier", tconf["supplier"],-1)	
			fsqlstr = fsqlstr + str_added
		}
		if( strings.Contains(arr_reportname[i], "sold_by") ){
			str_added =  strings.Replace(arr_reportname[i],"sold_by", tconf["soldby"],-1)	
			fsqlstr = fsqlstr + str_added
		}
		if( strings.Contains(arr_reportname[i], "customer") ){
			str_added =  strings.Replace(arr_reportname[i],"customer", tconf["customer"],-1)	
			fsqlstr = fsqlstr + str_added
		}
		if( strings.Contains(arr_reportname[i], "transferto") ){
			str_added =  strings.Replace(arr_reportname[i],"transferto", tconf["transferto"],-1)	
			fsqlstr = fsqlstr + str_added
		}
		if( strings.Contains(arr_reportname[i], "transferfr") ){
			str_added =  strings.Replace(arr_reportname[i],"transferfr", tconf["transferfrom"],-1)	
			fsqlstr = fsqlstr + str_added
		}
		if( strings.Contains(arr_reportname[i], "trantype") ){
			str_added =  strings.Replace(arr_reportname[i],"trantype", tconf["trantype"],-1)	
			fsqlstr = fsqlstr + str_added
		}
	
	}
	fmt.Println(fsqlstr)

	data, hdata := Datatable(fsqlstr,nil)
	var arrname []string
	type Htdata struct {
		Data 	[][]string
		Hdata 	[]string
		Arrname []string
	}
	var_htdata :=Htdata{data,hdata,nil}
	///rawreports_daily_lograw.html
	var default_tpl_name = "rawreports_tpl.html"
	var default_tpl_dir  = "inventory/custom_report/templates/rawreports_tpl.html"
	
	if (arr_reportname[0]==`lbr_rep_lograw 4, `){
		default_tpl_name = "rawreports_daily_lograw.html"
		default_tpl_dir  = "inventory/custom_report/templates/rawreports_daily_lograw.html"
		var tempname ="*****"
		for _,val := range data {
			if(val[0]==tempname){
				tempname = val[0]
			}
			if(val[0]!=tempname){
				//fmt.Println(val)
				arrname = append(arrname,val[0])
				tempname = val[0]	
			}

		}	
		var_htdata =Htdata{data,hdata,arrname}
	}
	tmpl := template.New(default_tpl_name).Funcs(config.FuncMap)
	var err error
	if tmpl, err = tmpl.ParseFiles(default_tpl_dir); err != nil {
		fmt.Println(err)
	}
	if(tconf["export_type"]=="excel"){
		w.Header().Set("Content-Disposition", "attachment; filename=RAW_REPORT.xls")
		w.Header().Set("Content-Type", "application/excel")
		w.Header().Set("Content-Type", "application/vnd.ms-excel")

	}else if(tconf["export_type"]=="word"){
		w.Header().Set("Content-Disposition", "attachment; filename=RAW_REPORT.doc")
		w.Header().Set("Content-Type", "application/msword")
	
	}else if(tconf["export_type"]=="pdf"){
		w.Header().Set("Content-Disposition", "attachment; filename=RAW_REPORT.pdf")
		w.Header().Set("Content-Type", "application/pdf")
	}
	err1 := tmpl.Execute(w,var_htdata)
	if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
	}

}	


func Special_rpt(w http.ResponseWriter, r *http.Request	){
	tconf := make(map[string]string)
	tconf["sql_code"] = r.URL.Query().Get("sql_code")
	tconf["rpt_url"] = r.URL.Query().Get("rpt_url")
	tconf["org_id"] = r.URL.Query().Get("org_id")
	tconf["trnID"] = r.URL.Query().Get("trnID")
	
	username, _ := login.Get_account_info(r)

	tconf["username"] = username
	fmt.Println("on printing")
	fmt.Println(tconf["sql_code"])
	fmt.Println(tconf["rpt_url"])
	//fmt.Println(path.Base(tconf["rpt_url"]))
	
	arr_reportname := strings.Split(string(tconf["rpt_url"]), `/`)
	fmt.Println(arr_reportname)
	fmt.Println(len(arr_reportname))
	fmt.Println(arr_reportname[len(arr_reportname)-1])

	
	
	data, hdata := Datatable(tconf["sql_code"],nil)
	//var arrname []string
	type Htdata struct {
		Data 	[][]string
		Hdata 	[]string
		Arrname []string
		Conf 	map[string]string
	}
	var_htdata :=Htdata{data,hdata,nil,tconf}
	///rawreports_dailys_lograw.html
	var default_tpl_name = arr_reportname[len(arr_reportname)-1]
	var default_tpl_dir  = string(tconf["rpt_url"])
	
	tmpl := template.New(default_tpl_name).Funcs(helper.FuncMap)
	var err error
	if tmpl, err = tmpl.ParseFiles(default_tpl_dir); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w,var_htdata)
	if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
	}


}

func Datatable(sql string, conf map[string]string ) ([][]string, []string) {
	
	var m  [][]string
	rows ,err, _,_ := config.Ap_sql(sql,1)
    if err != nil {
    	panic(err.Error())
    }
   	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
	}
	if cols == nil {
		fmt.Println(err)
	}
 	
	//getting the column name
	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
	}
	//end geting column key	

	//type containerdata [][]interface{}
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//mars at work here
		var a []string
		
		
		for i := 0; i < len(cols); i++ {  // fetching data tables to container map array
			a  = append(a, PrintValue(vals[i].(*interface{}))   )
		}
		m = append( m,a)
		//mars at work here
	}
		
	if rows.Err() != nil {
		fmt.Println(rows.Err)
	}	
	return m , cols

}

func PrintValue(pval *interface{}) string{
	var a string
	switch v := (*pval).(type) {
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
		a = fmt.Sprint(v.Format("2006-01-02 15:04:05.999"))

	default:
//		fmt.Print(v)
		a= fmt.Sprint(v)
	}
	return a
}

