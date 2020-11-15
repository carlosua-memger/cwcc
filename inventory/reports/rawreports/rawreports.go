package reports

import (
	
		"net/http"
		"fmt"
		"text/template"
		config "../../config"
	//	"strconv"
		_"encoding/json"
		//datatables "../../datatables"
		//login "../../login"
		//"time"
		"time"
		"strings"
		helper "../../helper"

)
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

type Stuct_rawrports struct{
	Id int
	Reportname string
	Sql string
}

func func_get_rawreportslist(userID string) []Stuct_rawrports{
	row ,err, _,_ := config.Ap_sql(`get_rawreport 1,`+userID,1)
	res := []Stuct_rawrports{}
	for row.Next() {
        var r Stuct_rawrports
        err = row.Scan(&r.Id, &r.Reportname,&r.Sql)
        if err != nil {
                panic(err.Error())
        }
		datarow := Stuct_rawrports{
        	r.Id,
	        r.Reportname,
	        r.Sql,

	    }
		res = append(res, datarow)

    }
    return res
}

func  Rawreports_nav(w http.ResponseWriter, r *http.Request) {
	/*login.Session_renew(w,r)
	_,session_user_id := login.Get_account_info(r)*/
// 	session_user_id_str := "1"

	type Htdata struct {
		LBR_Assign_Get []lbr_assign_name
		RawreportsList []Stuct_rawrports
	}
	
	var_htdata :=Htdata{nil,func_get_rawreportslist("1")}
	tmpl := template.New("rawreports_nav.html")
	var err error
	if tmpl, err = tmpl.ParseFiles("reports/rawreports/rawreports_nav.html"); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w,var_htdata)
	if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
	}

}

func Rawreports_doprint(w http.ResponseWriter, r *http.Request) {
	tconf := make(map[string]string)
	ary_reportname := strings.Split(r.URL.Query().Get("reportname"),"][")
	ary_branch := strings.Split(r.URL.Query().Get("branch"),"][")
	ary_trntype := strings.Split(r.URL.Query().Get("trntype"),"][")
	var ary_product []string
	//fmt.Println(`len here`, len(r.URL.Query().Get("product")));
	if (len(r.URL.Query().Get("product"))==0 ){
		ary_product = strings.Split("0][all","][")
	}else{
		ary_product = strings.Split(r.URL.Query().Get("product"),"][")
	}
	//ary_product := strings.Split("0][all","][")
	//ary_product[1] ="All"
	//InItem = append(InItem, post)
	ary_costcode := strings.Split(r.URL.Query().Get("costcode"),"][")
	ary_personnel := strings.Split(r.URL.Query().Get("personnel"),"][")


	fmt.Println(ary_reportname[0])

	tconf["reportname"] = ary_reportname[0]
	tconf["branch"] = ary_branch[0]
	tconf["trntype"] = ary_trntype[0]
	tconf["product"] = ary_product[0]
	tconf["costcode"] = ary_costcode[0]
	tconf["personnel"] = ary_personnel[0]

	tconf["reportname_info"] = ary_reportname[1]
	tconf["branch_info"] = "N/A"
	tconf["trntype_info"] = "N/A"
	tconf["product_info"] = "N/A"
	tconf["costcode_info"] = "N/A"
	tconf["personnel_info"] = "N/A"
	/*if(r.URL.Query().Get("refno")==""){
		tconf["refno"]="0"
	}else{
		tconf["refno"] = r.URL.Query().Get("refno")
	}*/
	tconf["refno"] = r.URL.Query().Get("refno")
	tconf["datefrom"] = r.URL.Query().Get("datefrom")
	tconf["dateto"] = r.URL.Query().Get("dateto")

	tconf["refno_info"] = "N/A"
	tconf["datefrom_info"] = "N/A"
	tconf["dateto_info"] = "N/A"


	tconf["export_type"] = r.URL.Query().Get("export_type")

	tconf["orgid"]="1"
	
	//strt := tconf["reportname"]

	//lbr_rep_logpair 3, 2, 0 , 46, '21 Oct 2015', '3 Nov 2015'
	arr_reportname := strings.Split(string(tconf["reportname"]), "@@")
	fsqlstr:=arr_reportname[0]
	for i := 1; i < len(arr_reportname); i++ {  
		str_added :=""
		if( strings.Contains(arr_reportname[i], "orgid") ){
			str_added =  strings.Replace(arr_reportname[i],"orgid", tconf["orgid"],-1)	
			fsqlstr = fsqlstr + str_added
		}
		if( strings.Contains(arr_reportname[i], "branch") ){
			str_added =  strings.Replace(arr_reportname[i],"branch", tconf["branch"],-1)	
			fsqlstr = fsqlstr + str_added
			tconf["branch_info"] = ary_branch[1]
		}
		if( strings.Contains(arr_reportname[i], "trntype") ){
			str_added =  strings.Replace(arr_reportname[i],"trntype", tconf["trntype"],-1)	
			fsqlstr = fsqlstr + str_added
			tconf["trntype_info"] = ary_trntype[1]
		}
		if( strings.Contains(arr_reportname[i], "product") ){
			str_added =  strings.Replace(arr_reportname[i],"product", tconf["product"],-1)	
			fsqlstr = fsqlstr + str_added
			tconf["product_info"] = ary_product[1]
		}
		if( strings.Contains(arr_reportname[i], "costcode") ){
			str_added =  strings.Replace(arr_reportname[i],"costcode", tconf["costcode"],-1)	
			fsqlstr = fsqlstr + str_added
			tconf["costcode_info"] = ary_costcode[1]
		}
		if( strings.Contains(arr_reportname[i], "personnel") ){
			str_added =  strings.Replace(arr_reportname[i],"personnel", tconf["personnel"],-1)	
			fsqlstr = fsqlstr + str_added
			tconf["personnel_info"] = ary_personnel[1]
		}
		if( strings.Contains(arr_reportname[i], "refno") ){
			str_added =  strings.Replace(arr_reportname[i],"refno", tconf["refno"],-1)	
			fsqlstr = fsqlstr + str_added
			tconf["refno_info"] = r.URL.Query().Get("refno")
		}
		if( strings.Contains(arr_reportname[i], "datefrom") ){
			str_added =  strings.Replace(arr_reportname[i],"datefrom", tconf["datefrom"],-1)	
			fsqlstr = fsqlstr + str_added
			tconf["datefrom_info"] = r.URL.Query().Get("datefrom")
		}
		if( strings.Contains(arr_reportname[i], "dateto") ){
			str_added =  strings.Replace(arr_reportname[i],"dateto", tconf["dateto"],-1)	
			fsqlstr = fsqlstr + str_added
			tconf["dateto_info"] = r.URL.Query().Get("dateto")
		}
	
	}
	fmt.Println(fsqlstr)

	data, hdata := Datatable(fsqlstr,nil)
	var arrname []string
	type Htdata struct {
		Data 	[][]string
		Hdata 	[]string
		Arrname []string
		Tconf map[string]string
	}
	var_htdata :=Htdata{data,hdata,nil,tconf}
	///rawreports_daily_lograw.html
	var default_tpl_name = "rawreports_tpl.html"
	var default_tpl_dir  = "reports/rawreports/rawreports_tpl.html"
	
	if (arr_reportname[0]==`lbr_rep_lograw 4, `){
		default_tpl_name = "rawreports_daily_lograw.html"
		default_tpl_dir  = "reports/rawreports/rawreports_daily_lograw.html"
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
		var_htdata =Htdata{data,hdata,arrname,tconf}
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

