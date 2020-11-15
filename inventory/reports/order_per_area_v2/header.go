package order_per_area

import (
	
		"net/http"
		"fmt"
		"text/template"
		config "../../../config"
		"strconv"
		"encoding/json"
		datatables "../../../datatables"
		login "../../../login"
		"strings"
		//helper "../../../helper"
		"time"
)

type Profile struct {
  	Message string  // error status output or validated data
  	Status bool //pass or fail
  	ID string
}




func HListHandler(w http.ResponseWriter, r *http.Request) {
	
	//login.Session_renew(w,r)
	//username, ses_uID := login.Get_account_info(r)
		tconf := make(map[string]string)
	   	Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		rights :=r.URL.Query().Get("rights")
		tconf["Org_id"] 	= str_OrgID
		tconf["Rights"] = rights
		tconf["Action_baseURL"] ="/inventory/customer_order/details?rights="+rights
		//tconf["ActionCol_param"] ="ID"


		tmpl := template.New("Hadd.html").Funcs(local_FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/reports/order_per_area_v2/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,tconf)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	

}

func ViewerHandler(w http.ResponseWriter, r *http.Request) {
	tconf := make(map[string]string)
	
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)

	tconf["item_id"] 	=r.URL.Query().Get("item_id")
	tconf["datefrom"] 	=r.URL.Query().Get("datefrom")
	tconf["dateto"] 	=r.URL.Query().Get("dateto")
	tconf["org_id"] 	= str_OrgID
	
	
	header := []string{ "Item Name","Supplier Name", "Beginning Balance", "Supplier Reciept", "Customer Sale", "Adjustment Entry", "Ending Balance"}  //column header	
	 
	type Data struct {
		Tconf 		map[string]string
		Header 		[]string
		Arr_data 	[][]string
	}
	str_page :=`viewer.html`
	str_page_dir := `inventory/reports/product_movement/viewer.html`
	

	tmpl := template.New(str_page).Funcs(local_FuncMap) 
	var err error
	if tmpl, err = tmpl.ParseFiles(str_page_dir); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w,Data{tconf,header,nil})
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
} 

func AjaxViewerHandler(w http.ResponseWriter, r *http.Request) {
	tconf := make(map[string]string)

	tconf["item_id"] 	=r.URL.Query().Get("item_id")
	tconf["datefrom"] 	=r.URL.Query().Get("datefrom")
	tconf["dateto"] 	=r.URL.Query().Get("dateto")
	tconf["org_id"] 	=r.URL.Query().Get("org_id")
	tconf["page"] 		=r.URL.Query().Get("page")
	tconf["item_key"] 		=r.URL.Query().Get("item_key")

	


	if(tconf["item_id"]==""){
		tconf["item_id"]="0"
	}
	

	fmt.Println("wetwew")
	fmt.Println(`SIS_Rep_Product_salemovement 1,`+tconf["org_id"]+`,`+tconf["item_id"]+`,'`+tconf["item_key"]+`',0, '`+tconf["datefrom"]+`','`+tconf["dateto"]+`',`+tconf["page"]+`,100`)
	arr_data := DataList(`SIS_Rep_Product_salemovement 1,`+tconf["org_id"]+`,`+tconf["item_id"]+`,'`+tconf["item_key"]+`',0,'`+tconf["datefrom"]+`','`+tconf["dateto"]+`',`+tconf["page"]+`,100`, tconf["org_id"])


	//fmt.Println(`SIS_Ledger_Report_asof 0, 3, '`+tconf["asofdate"]+`', `+tconf["location"]+`, 0, `+tconf["page"]+`, 25,''`)
	//arr_data := datatables.DataList(`SIS_Ledger_Report_asof 0, 3, '`+tconf["asofdate"]+`', `+tconf["location"]+`, 0, `+tconf["page"]+`, 25,''`  )
	js, err := json.Marshal(arr_data)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
		}
  	w.Header().Set("Content-Type", "application/json")
  	w.Write(js)	
} 



func AjaxViewerDetailsHandler(w http.ResponseWriter, r *http.Request) {
	tconf := make(map[string]string)
	
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)

	tconf["location"] 	=r.URL.Query().Get("location")
	tconf["item"] 	=r.URL.Query().Get("item")
	tconf["item_name"] 	=r.URL.Query().Get("item_name")
	tconf["location_name"] 	=r.URL.Query().Get("location_name")
	tconf["uprice"] 	=r.URL.Query().Get("uprice")
	tconf["status"] 	=r.URL.Query().Get("status")
	tconf["runbal"] 	=r.URL.Query().Get("runbal")
	
	
	fmt.Println(`sis_ledger_report_asof_detail 0 ,`+str_OrgID+`, `+tconf["location"]+`,`+tconf["item"])
	arr_data := datatables.DataList(`sis_ledger_report_asof_detail 0 ,`+str_OrgID+`, `+tconf["location"]+`,`+tconf["item"]   )
	 
	type Data struct {
		Tconf 		map[string]string
		Arr_data 	[][]string
	}
	var str_page string
	var str_page_dir string
	str_page=`viewer_details.html`
	str_page_dir = `inventory/ledger/viewer_details.html`

	tmpl := template.New(str_page).Funcs(local_FuncMap) 
	var err error
	if tmpl, err = tmpl.ParseFiles(str_page_dir); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w,Data{tconf,arr_data})
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}


type RetMessage struct {
	RetNo	int
    RetMsg	string
}

var local_FuncMap = template.FuncMap{
    	"Sql_list": func(s_qry string , org_id string, tag_id string) [][]string{
			fmt.Println(`sis_itemtags_get 1, `+org_id+`, 0 ,`+tag_id)
			retdata := datatables.DataList(`sis_itemtags_get 1, `+org_id+`, 0 ,`+tag_id)
			return retdata
		},
	}



func ViewitemHandler(w http.ResponseWriter, r *http.Request) {
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	tconf := make(map[string]string)
	

	tconf["img_name"] = r.URL.Query().Get("img_name")
	tconf["item_id"] = r.URL.Query().Get("item_id")
	tconf["graph_title"] = r.URL.Query().Get("item_name")
	//tconf["item_image"] = config.Img_to_render_custom1( tconf["img_name"] ,171)
	tconf["item_image"] = ""

	tconf["item_name"] = r.URL.Query().Get("item_name")
	tconf["startdate"] = r.URL.Query().Get("startdate")
	tconf["enddate"] = r.URL.Query().Get("enddate")
	
	tconf["ddata"] =""
	fmt.Println(tconf["ddata"])
	arr_ddata := strings.Split(	tconf["ddata"] , `][`)
	//tconf["item_image"] 			= "config.Img_to_render_custom1(arr_ddata[0],171)"
	
	tconf["item_image"] 			= strings.Replace(r.URL.Query().Get("item_img"), "size_80", "size_171", 2)
	fmt.Println(`img src : `, tconf["item_image"])
	tconf["item_id_hideme"] 		= " "
	tconf["supplier_name"] 			= " "
	tconf["beg_vol"] 				= " "
	tconf["suprec_vol"] 			= " "
	tconf["cussal_vol"] 			= " "
	tconf["adjent_vol"] 			= " "
	tconf["end_vol"] 				= " "
	
	fmt.Println(arr_ddata)
	fmt.Println(`SIS_Rep_Product_salemovement_detail 1, `+str_OrgID+`, `+tconf["item_id"]+`, '`+tconf["startdate"]+`', '`+tconf["enddate"]+`'`);
	graph1_Data := datatables.DataList(`SIS_Rep_Product_salemovement_detail 1, `+str_OrgID+`, `+tconf["item_id"]+`, '`+tconf["startdate"]+`', '`+tconf["enddate"]+`'` )
	
	fmt.Println(graph1_Data)
	tconf["graph1_data"] =` [`
	tconf["graph1_data2"] =` [`
	for _ , value := range graph1_Data{
			fmt.Println(value[0])
			tconf["graph1_data"] +=`['`+value[1]+`', `+value[2]+`],`
			tconf["graph1_data2"] +=`['`+value[1]+`', `+value[3]+`],`

	}
	tconf["graph1_data"] +=`]`
	tconf["graph1_data2"] +=`]`

	//
	//	
	fmt.Println(`SIS_Rep_Product_salemovement_detail 2, `+str_OrgID+`, `+tconf["item_id"]+`,'`+tconf["startdate"]+`', '`+tconf["enddate"]+`'`);
	rowret ,err, _,_ := config.Ap_sql(`SIS_Rep_Product_salemovement_detail 2, `+str_OrgID+`, `+tconf["item_id"]+`, '`+tconf["startdate"]+`', '`+tconf["enddate"]+`'`,1)
	if err != nil {
		panic(err.Error())
	}
	type Graph2 struct {
		Date_s	time.Time
	    Vals	string
	    Vals2   string
	}

	var RG2 Graph2
	if err != nil {
        fmt.Println(err.Error)
        panic(err.Error())
    }
    tconf["graph2_data"] =` [`
    tconf["graph2_data2"] =` [`

    for rowret.Next() {
	    err = rowret.Scan(&RG2.Date_s,&RG2.Vals,&RG2.Vals2)
	    if err != nil {
	        panic(err.Error())
	    }
	    fmt.Println(RG2.Date_s)
		fmt.Println(RG2.Vals)
		str_dd 	 :=strconv.Itoa(RG2.Date_s.Day())
		var i_mm int = int(RG2.Date_s.Month()) - 1
		str_mm 	 :=strconv.Itoa(i_mm) 
		str_yyyy :=strconv.Itoa(RG2.Date_s.Year())
		fmt.Println(`Day:` , str_dd)
		fmt.Println(`MM:` , str_mm)
		fmt.Println(`YYYY:` , str_yyyy)
		

		tconf["graph2_data"] +=`[Date.UTC(`+str_yyyy+`, `+str_mm+`, `+str_dd+`), `+RG2.Vals+`],`
		tconf["graph2_data2"] +=`[Date.UTC(`+str_yyyy+`, `+str_mm+`, `+str_dd+`), `+RG2.Vals2+`],`
	 	
	 	//r = RetMessage{r.Date_s,r.Date_s}
	}
	tconf["graph2_data"] +=`]`
	tconf["graph2_data2"] +=`]`
	
	tmpl := template.New("view_item_stat.html")
	//var err error
	if tmpl, err = tmpl.ParseFiles("inventory/reports/product_movement/view_item_stat.html"); err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w,tconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} 

}	


func DataList(sql string, org string ) [][]string {
	
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
			if(i==0){
				a  = append(a, config.Img_to_render_custom1(PrintValue(vals[i].(*interface{}) ),80,org )  )
			}else{
				a  = append(a, PrintValue(vals[i].(*interface{}))  )
			}
		}
		m = append( m,a)
	}
		
	if rows.Err() != nil {
		fmt.Println(rows.Err)
	}	
	return m

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




