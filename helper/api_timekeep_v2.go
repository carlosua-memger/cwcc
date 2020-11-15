package helper

import(
	"encoding/json"
	//"fmt"
	//"log"
	"net/http"
	mux "../github.com/gorilla/mux"
	"fmt"
	//datatables "../datatables"
	config "../config"

)
// function route control
func Timekeeping_sub_module_timekeep_v2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if( vars["key1"]==`assignment` ){
		Timekeep_v2_assignment(w,r)
	}
	if( vars["key1"]==`LBR_LogTotal_List` ){
		Timekeep_v2_LBR_LogTotal_List(w,r)
	}
	if( vars["key1"]==`Lbr_assign_get` ){
		Lbr_assign_get(w,r)
	}
	if( vars["key1"]==`LBR_LogRaw_get` ){
		LBR_LogRaw_get(w,r)
	}
	if( vars["key1"]==`LBR_LogPair_get` ){
		LBR_LogPair_get(w,r)
	}
	if(vars["key1"]==`LBR_TK_Pairing_Save`){
		LBR_TK_Pairing_Save(w,r)
	}
	if(vars["key1"]==`lbr_logtotal_update`){
		Lbr_logtotal_update(w,r)
	}
	if(vars["key1"]==`LBR_TK_Raw_Save`){
		LBR_TK_Raw_Save(w,r)
	}
	if(vars["key1"]==`LBR_LogPair_Delete`){
		LBR_LogPair_Delete(w,r)
	}

	


	
	/*if(vars["key1"]==`Lbr_assign_get_raw`){
		Lbr_assign_get(w,r)
	}*/

}

func SQL_to_JSON(sqlstr string ,w http.ResponseWriter){
	fmt.Println(sqlstr)
	rows ,err, _,_  := config.Ap_sql(sqlstr,1)    
    

    if err != nil {
	      http.Error(w, err.Error(), http.StatusInternalServerError)
    		return
	  }
	  defer rows.Close()
	  columns, err := rows.Columns()
	  if err != nil {
	     http.Error(w, err.Error(), http.StatusInternalServerError)
    		return
	  }
	  count := len(columns)
	  tableData := make([]map[string]interface{}, 0)
	  values := make([]interface{}, count)
	  valuePtrs := make([]interface{}, count)
	  for rows.Next() {
	      for i := 0; i < count; i++ {
	          valuePtrs[i] = &values[i]
	      }
	      rows.Scan(valuePtrs...)
	      entry := make(map[string]interface{})
	      for i, col := range columns {
	          var v interface{}
	          val := values[i]
	          b, ok := val.([]byte)
	          if ok {
	              v = string(b)
	          } else {
	              v = val
	          }
	          entry[col] = v
	      }
	      tableData = append(tableData, entry)
	  }
	  //jsonData, err := json.Marshal(tableData)
	  if err != nil {
	      http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
	  }
	 // fmt.Println(string(jsonData))
	  //return string(jsonData), nil 

	js, err := json.Marshal(tableData)
    
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
		}
  	w.Header().Set("Content-Type", "application/json")
  	//w.Header().Set("Content-Type", "application/json")
  	w.Write(js)

}


func LBR_LogPair_Delete(w http.ResponseWriter, r *http.Request) {
	//API endpoint : /api/ajax/timekeeping/timekeep_v2/{org}/Lbr_logtotal_update
	r.ParseForm()


   	LogPair_ID 		:= 	r.Form["LogPair_ID"][0]
	user 			:= 	r.Form["user"][0]
	


	sqlstr := `exec LBR_LogPair_Delete '`+user+ `',`+LogPair_ID
	SQL_to_JSON(sqlstr,w);
}

func LBR_TK_Raw_Save(w http.ResponseWriter, r *http.Request) {
	//API endpoint : /api/ajax/timekeeping/timekeep_v2/{org}/Lbr_logtotal_update
	r.ParseForm()


    logdatetime 	:= 	r.Form["logdatetime"][0]
	memberID 		:= 	r.Form["memberID"][0]
	org 			:= 	r.Form["org"][0]
	user 			:= 	r.Form["user"][0]

	sqlstr := `exec LBR_TK_Raw_Save '`+user+ `',`+org+`,`+memberID+`,'`+logdatetime+`'`
	SQL_to_JSON(sqlstr,w);
}

func Lbr_logtotal_update(w http.ResponseWriter, r *http.Request) {
	//API endpoint : /api/ajax/timekeeping/timekeep_v2/{org}/Lbr_logtotal_update
	
	r.ParseForm()


    lbr_logtotalID 	:= 	r.Form["lbr_logtotalID"][0]
	user 			:= 	r.Form["user"][0]
	hr_reg 			:= 	r.Form["hr_reg"][0]
	lbr_daycodeID 	:= 	r.Form["lbr_daycodeID"][0]
	hr_ot 			:= 	r.Form["hr_ot"][0]
	hr_nd 			:= 	r.Form["hr_nd"][0]
	hr_ndot 		:= 	r.Form["hr_ndot"][0]
	hr_break 		:= 	r.Form["hr_break"][0]
	min_late 		:= 	r.Form["min_late"][0]
	min_undertime 	:= 	r.Form["min_undertime"][0]
	remarks 		:= 	r.Form["remarks"][0]
	


	sqlstr := `exec lbr_logtotal_update   `+lbr_logtotalID+`, '`+user+ `',`+lbr_daycodeID+`,`+hr_reg+`,`+hr_ot+`,`+hr_nd+`,`+ hr_ndot + `,`+hr_break+`,`+min_late+`,`+ min_undertime +`,'`+remarks +`'`
	SQL_to_JSON(sqlstr,w);
}


func LBR_TK_Pairing_Save(w http.ResponseWriter, r *http.Request) {
	//API endpoint : /api/ajax/timekeeping/timekeep_v2/{org}/LBR_TK_Pairing_Save
	
	r.ParseForm()

	mode := 	r.Form["mode"][0]
	FormUser := r.Form["FormUser"][0]
	org := 		r.Form["org"][0]
	member := 	r.Form["member"][0]
	assign := 	r.Form["assign"][0]
	logdate := 	r.Form["logdate"][0]
	login := 	r.Form["login"][0]
	logout := 	r.Form["logout"][0]
	produced := r.Form["produced"][0]
	remarks := 	r.Form["remarks"][0]


	sqlstr := `exec LBR_TK_Pairing_Save   `+mode+`, '`+FormUser+ `',0,`+org+`,`+member+`,`+assign+`,'`+logdate+`','`+login + `','`+logout+`',`+produced+`,'`+remarks +`'`
	SQL_to_JSON(sqlstr,w);
}

func Timekeep_v2_LBR_LogTotal_List(w http.ResponseWriter, r *http.Request) {
	//API endpoint : /api/ajax/timekeeping/timekeep_v2/{org}/LBR_LogTotal_List
	mode := r.URL.Query().Get("mode")
	org := r.URL.Query().Get("org")
	TopAssignment := r.URL.Query().Get("TopAssignment")
	Logdate := r.URL.Query().Get("Logdate")
	PageNo := r.URL.Query().Get("PageNo")
	RowPerPage := r.URL.Query().Get("RowPerPage")
	SearchKey := r.URL.Query().Get("SearchKey")
	sqlstr := `exec LBR_LogTotal_List ` +mode+`, `+org+`, `+ TopAssignment +`, '`+Logdate +`', `+PageNo+`, `+RowPerPage+`, '`+SearchKey+`' `
	SQL_to_JSON(sqlstr,w)
}

func Lbr_assign_get(w http.ResponseWriter, r *http.Request) {
	//API endpoint : /api/ajax/timekeeping/timekeep_v2/{org}/Lbr_assign_get
	mode := r.URL.Query().Get("mode")
	org := r.URL.Query().Get("org")
	TopAssignment := r.URL.Query().Get("TopAssignment")
	/*Logdate := r.URL.Query().Get("Logdate")
	PageNo := r.URL.Query().Get("PageNo")
	RowPerPage := r.URL.Query().Get("RowPerPage")
	SearchKey := r.URL.Query().Get("SearchKey")*/
	sqlstr := `exec lbr_assign_get  `+mode+`,`+org+`,`+TopAssignment
	SQL_to_JSON(sqlstr,w);
}

func LBR_LogRaw_get(w http.ResponseWriter, r *http.Request) {
	//API endpoint : /api/ajax/timekeeping/timekeep_v2/{org}/LBR_LogRaw_get
	mode := r.URL.Query().Get("mode")
	org := r.URL.Query().Get("org")
	MemberID := r.URL.Query().Get("MemberID")
	TopAssignment := r.URL.Query().Get("TopAssignment")
	LogDate := r.URL.Query().Get("LogDate")
	sqlstr := `exec LBR_LogRaw_get   `+mode+` ,`+org+`, `+MemberID+ `,`+TopAssignment+`,'`+LogDate+`','`+LogDate+`'` 
	SQL_to_JSON(sqlstr,w);

}

func LBR_LogPair_get(w http.ResponseWriter, r *http.Request) {
	//API endpoint : /api/ajax/timekeeping/timekeep_v2/{org}/LBR_LogPair_get
	mode := r.URL.Query().Get("mode")
	org := r.URL.Query().Get("org")
	MemberID := r.URL.Query().Get("MemberID")
	TopAssignment := r.URL.Query().Get("TopAssignment")
	LogDate := r.URL.Query().Get("LogDate")
	sqlstr := `exec LBR_LogPair_get  `+mode+` ,0 ,`+org+`, `+MemberID+ `,`+TopAssignment+`,'`+LogDate+`','`+LogDate+`'` 
	SQL_to_JSON(sqlstr,w);

}

/*func Lbr_assign_get(w http.ResponseWriter, r *http.Request) {
	//API endpoint : /api/ajax/timekeeping/timekeep_v2/{org}/Lbr_assign_get
	mode := r.URL.Query().Get("mode")
	org := r.URL.Query().Get("org")
	TopAssignment := r.URL.Query().Get("TopAssignment")
	sqlstr := `exec lbr_assign_get `+mode+`,`+org+`, '`+TopAssignment + `'`
	SQL_to_JSON(sqlstr,w);

}*/

func Timekeep_v2_assignment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	getq := r.URL.Query().Get("q")
	mode := r.URL.Query().Get("mode")
	//API endpoint : /api/ajax/timekeeping/timekeep_v2/{org}/assignment
	orgid := vars["org"]
	type Inner_item struct {
		Id int `json:"id"`
		Fullname string `json:"full_name"`
	}
	type JsondataEmpty struct {
    	Total_count int `json:"total_count"`
    	Incomplete_results bool `json:"incomplete_results"`
    	Items []Inner_item `json:"items"`
    	
    }

	var count int
	var sqlstr string
	sqlstr = `exec lbr_assign_get `+mode+`,`+orgid+`, '`+getq + `'`
	fmt.Println(sqlstr)
	rows ,err, _,_  := config.Ap_sql(sqlstr,1)    
    if err != nil {
		fmt.Println("db error:",err)
	}
	InItem := []Inner_item{}
    //var InItemRow Inner_item
    for rows.Next() {
        var r Inner_item
        err = rows.Scan(&r.Id,&r.Fullname)
        if err != nil {
            panic(err.Error())
        }
		count = count +1
		post := Inner_item {
        	r.Id,
        	r.Fullname, 
	    }
	    InItem = append(InItem, post)
    	//InItemRow = Inner_item{r.Id,r.Fullname}
    }
    

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