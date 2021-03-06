package product_movement

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
		helper "../../../helper"
		"time"
)

type Profile struct {
  	Message string  // error status output or validated data
  	Status bool //pass or fail
  	ID string
}







func HListHandler(w http.ResponseWriter, r *http.Request) {
	
	login.Session_renew(w,r)
	username, ses_uID := login.Get_account_info(r)
	
	//username := "static"
	//ses_uID := 1
	str_ses_uID :=strconv.Itoa(ses_uID)
	rights :=r.URL.Query().Get("rights")
	fmt.Println(str_ses_uID)
	//fmt.Fprint(w,rights)
	
	if r.Method =="GET" {	
		
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		tconf := make(map[string]string)
		//tconf["parentID"] =r.URL.Query().Get("parentID")
		tconf["username"] =username
		tconf["Org_id"] = str_OrgID
		

		//arr_data := datatables.DataList(`select id,dbo.SIS_Itemclass_Name(1, 3, id) [itemclass_name] from sis_itemclass`  )
		//arr_data := datatables.DataList(`SIS_ItemClass_Get 2,`+str_OrgID  )
		//arr_data_supplier := datatables.DataList(`select id,dbo.SIS_Supplier_Name(1, 3, id) [supplier_name] from SIS_Supplier`  )

		//arr_data_supplier := datatables.DataList(`SIS_Supplier_Get 2,`+str_OrgID  )
		//arr_data_itemtag := datatables.DataList(`sis_itemtags_list 1,`+str_OrgID+` , 0`  )
		arr_data_supplier:= datatables.DataList(`exec SIS_Supplier_Get 2,`+str_OrgID  )
		arr_data_location := datatables.DataList(`exec SIS_Location_Get 6,`+str_OrgID  )
		/*
		arr_data_item_status := datatables.DataList(`exec SIS_Item_Status_Get 0,`+str_OrgID  )
		arr_data_price := datatables.DataList(`SIS_Item_price_Get 0 ,`+str_OrgID  )
		arr_data_item_class := datatables.DataList(`SIS_ItemClass_Get 2 ,`+str_OrgID  )
		*/
		

		type Data struct {
			Rights	string	
			Tconf 	map[string]string
			Ls_itemClass [][]string
			Ls_supplier [][]string
			Arr_itemtag [][]string
			Arr_location [][]string
			Arr_item_status [][]string
			Arr_price [][]string
			Arr_item_class [][]string



		}   
		tmpl := template.New("Hadd.html").Funcs(local_FuncMap) 
		var err error
		if tmpl, err = tmpl.ParseFiles("inventory/reports/item_movement/Hadd.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&Data{rights , tconf,nil,arr_data_supplier,nil,arr_data_location,nil,nil,nil})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {

		r.ParseForm()	
		img_dir := r.Form["Image_dr"][0]
		item_id := r.Form["Item_id"][0]
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		helper.Img_resizer(img_dir,item_id,str_OrgID)

		js, err := json.Marshal(`resizer success...`)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)	
		

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
	tconf["location"] 	=r.URL.Query().Get("location")
	tconf["tranchkbox"] =r.URL.Query().Get("tranchkbox")
	
	
	header := []string{"Item","Date", "Transaction Type", "Qty"}  //column header	
	 
	type Data struct {
		Tconf 		map[string]string
		Header 		[]string
		Arr_data 	[][]string
	}
	str_page :=`viewer.html`
	str_page_dir := `inventory/reports/item_movement/viewer.html`
	

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
	tconf["item_key"] 	=r.URL.Query().Get("item_key")
	tconf["location"] 	=r.URL.Query().Get("location")
	tconf["tranchkbox"] =r.URL.Query().Get("tranchkbox")
	//tconf["item_image"] = strings.Replace(r.URL.Query().Get("item_id"), "size_80", "size_171", 2)
	//tconf["name_with_img"] ="false"
	//tconf["col_img_num"] ="1"
	//tconf["img_id_col"] ="0"
	//tconf["org"]  =r.URL.Query().Get("org_id")
	//tconf["col_hidden1"] ="0"
	//tconf["col_img"] ="false"
	


	if(tconf["item_id"]!=""){
		tconf["item_id"]="id:" +tconf["item_id"]
	}else{
			tconf["item_id"]=""

	}
	
	var itemfilter string


	itemfilter = tconf["item_id"] + tconf["item_key"]


	


	fmt.Println("wetwew")
	fmt.Println(`SIS_Rep_Item_Movement 1,`+tconf["org_id"]+`,`+tconf["location"]+`,'`+itemfilter+`','`+tconf["datefrom"]+`','`+tconf["dateto"]+`',`+tconf["tranchkbox"]+`,100`)
	arr_data := DataList(`SIS_Rep_Item_Movement 1,`+tconf["org_id"]+`,`+tconf["location"]+`,'`+itemfilter+`','`+tconf["datefrom"]+`','`+tconf["dateto"]+`','`+tconf["tranchkbox"]+`'`, tconf["org_id"])


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


func Ajax_getItem(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	getq := r.URL.Query().Get("q")
	getid := r.URL.Query().Get("id")
	isbaseonassignid :=r.URL.Query().Get("isbaseonassignid")
	assignmentID := r.URL.Query().Get("assignmentID")
	fmt.Println(`is base sa assignID`,isbaseonassignid)
	fmt.Println(` assignID`,assignmentID)

	getq = getq + getid
		type Inner_item struct {
			Id string `json:"id"`
			Fullname string `json:"full_name"`
			Item_image string `json:"avatar"`
			Price string `json:"price"`
		}
		type JsondataEmpty struct {
	    	Total_count int `json:"total_count"`
	    	Incomplete_results bool `json:"incomplete_results"`
	    	Items []Inner_item `json:"items"`
	    	
	    }
		var count int
		/*var strtype string
		if getid !=""{
			strtype = "1"
		}else{
			strtype = "2"
		}	*/
		//`exec lbr_memass_search 1, 395 , '' ,'',''`	
		var sqlstr string
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)
		 if (isbaseonassignid ==`1`){
		 	//sqlstr = `exec lbr_memass_search 1, `+assignmentID+` , '' ,'',''`
		 	sqlstr ="exec SIS_Item_Get 3, "+str_OrgID+",'"+getq+"'"
		 }else{
		 	//sqlstr = "exec LBR_Member_Get "+strtype+",  '"+getq+"'"
		 	sqlstr = "exec SIS_Item_Get 6, "+str_OrgID+",'"+getq+"'"


		 }
		 fmt.Println(sqlstr)
		 rows ,err, _,_  := config.Ap_sql(sqlstr,1)    
	    if err != nil {
			fmt.Println("db error:",err)
		}
		InItem := []Inner_item{}
	    var InItemRow Inner_item
	    for rows.Next() {
	        var r Inner_item
	        err = rows.Scan(&r.Id,&r.Fullname,&r.Item_image,&r.Price)
	        if err != nil {
	            panic(err.Error())
	        }
			count = count +1
			post := Inner_item {
	        	r.Id,
	        	r.Fullname, 
	        	config.Img_to_render(`img_uploads/size_80/`+str_OrgID+`/`+r.Item_image+`.jpg`,r.Item_image,80,str_OrgID), 
		    	r.Price,
		    }
		    InItem = append(InItem, post)
	    InItemRow = Inner_item{r.Id,r.Item_image, config.Img_to_render(`img_uploads/size_80/`+str_OrgID+`/` +r.Item_image+`.jpg`,r.Item_image,80,str_OrgID),r.Price}
	    }
	    

	  	if getid !=""{
		
			//jsondata := JsondataEmpty{count,false,InItem}
		    js, err := json.Marshal(InItemRow)
		    
			if err != nil {
		    	http.Error(w, err.Error(), http.StatusInternalServerError)
		    	return
				}
		  	w.Header().Set("Content-Type", "application/json")
		  	//w.Header().Set("Content-Type", "application/json")
		  	w.Write(js)
		
		}else{
		
			jsondata := JsondataEmpty{count,false,InItem}
		     js, err := json.Marshal(jsondata)
		    
			if err != nil {
		    	http.Error(w, err.Error(), http.StatusInternalServerError)
		    	return
				}
		  	w.Header().Set("Content-Type", "application/json")
		  	//w.Header().Set("Content-Type", "application/json")
		  	w.Write(js)

		}	

}



