package helper

import(
	"encoding/json"
	//"fmt"
	//"log"
	"net/http"
	mux "../github.com/gorilla/mux"
	"fmt"
	datatables "../datatables"
	//config "../config"

)

type ReturnData struct {
	RecordCount int
	LocationID  string

}

func GetLatestLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	module := vars["module"]
	orgid := vars["org"]
	//key1 := vars["key1"]
	fmt.Println(module)
	
	//mux.HandleFunc("/api/ajax/inventory/getlatestlocation/{module}/{orgid}/{key1}",helper.GetLatestLocation).Methods("GET");	
	//todos:= `wewrwew`
	
	data_ret := datatables.Data_row(`SIS_location_get 8,`+orgid+`,''`)
	var ss ReturnData
	if (len(data_ret)>0) {
		 ss= ReturnData{1, data_ret[0]}
	}else{
		 ss= ReturnData{0,`0`}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(  ss ); err != nil {
        panic(err)
    }

}

func DefaultApiFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	module := vars["module"]
	
	/// API module control
	if(vars["module"]==`timekeeping`){
		Timekeeping_module(w,r)
	
	}else if(vars["module"]==`inventory`){
		
		if( vars["sub_module"]=="get-to-json" ){

			sqlstr := r.URL.Query().Get("sql_str")
			fmt.Printf("%v", sqlstr)
			SQL_to_JSON(sqlstr,w)

		}

	}else{
		//sub_module := vars["sub_module"]
		orgid := vars["org"]
		//key1 := vars["key1"]
		fmt.Println(module)
		fmt.Println(orgid)

		
		//mux.HandleFunc("/api/ajax/{module}/{sub_module}/{org}/{key1}",helper.DefaultApiFunc)

		type Data struct {
			Data [][]string	`json:"data"`
		}

		data_ret := datatables.DataList(`SIS_SupRetHdr_Misc 1, 23`)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	    w.WriteHeader(http.StatusOK)
	    if err := json.NewEncoder(w).Encode(  Data{data_ret} ); err != nil {
	        panic(err)
	    }
	}    
		

}


// timekeeping API function module
func Timekeeping_module(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sub_module := vars["sub_module"]
	//API Sub-module control
	if (sub_module==`timekeep_v2`){
		Timekeeping_sub_module_timekeep_v2(w,r)
	}

}

	


func DefaultApiFunc_all(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	transaction := vars["transaction"]
	version := vars["version"]
	fmt.Println(transaction)
	fmt.Println(version)

	str_OrgID :=r.URL.Query().Get("org_id")
	hdr_id	:= r.URL.Query().Get("hdr_id")
	h_id	:= r.URL.Query().Get("h_id")
	search_key :=r.URL.Query().Get("search_key") 
	//mux.HandleFunc("/api/ajax/api/{transaction}/{version}",helper.DefaultApiFunc_all)	
	type Data struct {
		Data [][]string	`json:"items"`
	}
	//SIS_Item_Get 13,69,'336', 0, 0 , 'mikata'
	
	if(transaction==`get_int_account`){
		fmt.Println(`exec SIS_IntAct_Get 2,`+str_OrgID+`, 'IntOrd', '`+h_id+`' , '`+hdr_id+`','`+search_key+`'`    );
		data_ret := datatables.DataList(`exec SIS_IntAct_Get 2,`+str_OrgID+`, 'IntOrd', '`+h_id+`' , '`+hdr_id+`' ,'`+search_key+`'`   )
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	    w.WriteHeader(http.StatusOK)
	    if err := json.NewEncoder(w).Encode(  Data{data_ret} ); err != nil {
	        panic(err)
	    }
	} else {
		fmt.Println(`SIS_Item_Get 13,`+str_OrgID+`,'`+hdr_id+`', 0, 0,'`+search_key+`'` );
		data_ret := datatables.DataList(`SIS_Item_Get 13,`+str_OrgID+`,'`+hdr_id+`', 0, 0,'`+search_key+`'`   )
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	    w.WriteHeader(http.StatusOK)
	    if err := json.NewEncoder(w).Encode(  Data{data_ret} ); err != nil {
	        panic(err)
	    }
	}

	

}







func Get(w http.ResponseWriter, r *http.Request) {
	//API endpoint : /api/ajax/timekeeping/timekeep_v2/{org}/Lbr_logtotal_update
	r.ParseForm()


   	LogPair_ID 		:= 	r.Form["LogPair_ID"][0]
	user 			:= 	r.Form["user"][0]
	


	sqlstr := `exec LBR_LogPair_Delete '`+user+ `',`+LogPair_ID
	SQL_to_JSON(sqlstr,w);
}
