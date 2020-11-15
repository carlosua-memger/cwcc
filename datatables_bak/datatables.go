package datatables	
import (
	"fmt"
	"time"
	"strconv"
	config "../config"
	"net/http"
	"encoding/json"
	 "strings"
)



func DatatableretArray(w http.ResponseWriter, conf map[string]string) {
	rows1 ,err1, _,_  := config.Ap_sql(conf["sql_total"],1)
	    var err error
	    if err1!= nil {
            panic(err1.Error())
	    }
	    type Records struct{
	    	Page_current int
	    	Page_more int
	    	Page_total int
	    	Record_more int
	    	Record_total int
	    }
	    var Recs Records
	    for rows1.Next() {
	    	err = rows1.Scan(&Recs.Page_current, &Recs.Page_more, &Recs.Page_total, &Recs.Record_more, &Recs.Record_total)
	            if err != nil {
	                    panic(err.Error())
	            }
	            
	    } 
	    //"data":[]
	 	 drr1,_:=strconv.Atoi(conf["dr"])
	    type Jsondata struct {
	    	Draw int `json:"draw"`
	    	RecordsTotal int `json:"recordsTotal"`
	    	RecordsFiltered int `json:"recordsFiltered"`
	    	Data [][]string `json:"data"`
	    }
	     type JsondataEmpty struct {
	    	total_count int `json:"total_count"`
	    	RecordsTotal int `json:"recordsTotal"`
	    	RecordsFiltered int `json:"recordsFiltered"`
	    	Data []string `json:"data"`
	    }
		
		objdata :=Datatable(conf["sql_data"],conf)
	    if  Recs.Record_total ==0{
	    	objdata1 := []string{} 
 			dd1 := &JsondataEmpty{drr1, Recs.Record_total , Recs.Record_total ,objdata1}
 			js, err := json.Marshal(dd1)
 			if err != nil {
		    	http.Error(w, err.Error(), http.StatusInternalServerError)
		    	return
	  		}
		  	w.Header().Set("Content-Type", "application/json")
		  	w.Write(js)	  
	    }else{
	    	dd := &Jsondata{drr1, Recs.Record_total , Recs.Record_total ,objdata}
		    js, err := json.Marshal(dd)
		    

			if err != nil {
		    	http.Error(w, err.Error(), http.StatusInternalServerError)
		    	return
	  		}
	  		w.Header().Set("Content-Type", "application/json")
		  	w.Write(js)	  
		}
}



func Datatable(sql string, conf map[string]string ) [][]string {
	
	var m  [][]string
	m1 :=make(map[string]int)
	var Actioncol bool
	Actioncol = true
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
	//geting columnheader key
	var hkey_action_param int
	for key , value := range cols{
			if ( strings.ToUpper(value)==strings.ToUpper(conf["ActionCol_param"])){
				hkey_action_param = key
			}
	}
	//end geting column key	
	draw1,_:=strconv.Atoi(conf["dr"] )
	m1["draw"] = draw1
//	m1["recordsTotal"] = strconv.Atoi(conf["recordsTotal"]) 
//	m1["recordsFiltered"] = strconv.Atoi(conf["recordsFiltered"]) 

	//type containerdata [][]interface{}
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//mars at work here
		var a []string
		var fontcolor  string
		var posted bool
		var Cancelled bool
		fontcolor =""
		/*if( strings.Contains( PrintValue(vals[3].(*interface{})) , "Posted") ){
			fontcolor ="style='color:red;'"
		}*/
		if( strings.Contains( PrintValue(vals[3].(*interface{})) , "Posted") ){
			posted = true
		}
		if( strings.Contains( PrintValue(vals[2].(*interface{})) , "Posted") ){
			posted = true
		}

		if( strings.Contains( PrintValue(vals[3].(*interface{})) , "Cancelled") ){
			Cancelled = true
		}
		if( strings.Contains( PrintValue(vals[2].(*interface{})) , "Cancelled") ){
			Cancelled = true
		}
		
		for i := 0; i < len(cols); i++ {  // fetching data tables to container map array
			if Actioncol && i ==0{
				var EditAction , DeleteAction ,AddChildAction,AddroleAction,AddDetailAction,AssignmentAction ,MembersAction  string
				
				if conf["ActionCol_edit"] =="true"{
					if conf["ActionCol_edit_is_modal"] =="true"{
						EditAction = `<a data-toggle="modal" href="#basic" onclick="Modalfunc('`+conf["ActionCol_edit_url"] + PrintValue(vals[hkey_action_param].(*interface{}))+`',this);">Edit</a>`
					}else{
						//EditAction =`<a  name="`+conf["ActionCol_param"]+`' href ="#`+ conf["ActionCol_edit_url"] +`?id=`+PrintValue(vals[i].(*interface{})) +`" onclick="Nav2(this,`+PrintValue(vals[hkey_action_param].(*interface{})) +`)" >Edit</a>`
						EditAction =`<a href ="#`+ conf["ActionCol_edit_url"]+PrintValue(vals[i].(*interface{})) +`" onclick="func_edit('`+conf["ActionCol_edit_url"]+PrintValue(vals[hkey_action_param].(*interface{})) +`',this)" >Edit</a>`
					}
					if posted {
						EditAction = ``	
					}
					if Cancelled{
						EditAction = ``
					}	
				}
				if conf["ActionCol_delete"] =="true"{
					DeleteAction = `<a data-toggle="modal" href="#modalSession" onclick="Modalfuncdelete('`+conf["ActionCol_delete_url"]+ PrintValue(vals[hkey_action_param].(*interface{})) +`',1,this)" >Delete</a>`
					if posted {
						DeleteAction = ``	
					}
					if Cancelled {
						DeleteAction = ``
					}
				}	
				if conf["ActionCol_add_child"] =="true"{	
					//AddChildAction ="<a style='color: "+fontcolor+";' href ='#"+PrintValue(vals[i].(*interface{})) +"'>Add Child</a>"
					AddChildAction =`<a href ="#`+ conf["ActionCol_edit_url"]+PrintValue(vals[i].(*interface{})) +`" onclick="func_add_child('`+conf["ActionCol_add_child_url"]+PrintValue(vals[hkey_action_param].(*interface{})) +`',this)" >Add Child</a>`
				}
				if conf["ActionCol_roles"] =="true"{	
					AddroleAction ="<a  name='"+conf["ActionCol_param"]+"' href ='#"+ conf["ActionCol_roles_url"] +"?id="+PrintValue(vals[i].(*interface{})) +"' onclick='Nav2(this,"+PrintValue(vals[hkey_action_param].(*interface{})) +")' >Roles</a>"
				}
				if conf["ActionCol_detail"] =="true"{	
					AddDetailAction ="<a  name='"+conf["ActionCol_param"]+"' href ='#"+ conf["ActionCol_detail_url"] +"?id="+PrintValue(vals[i].(*interface{})) +"' onclick='Nav2(this,"+PrintValue(vals[hkey_action_param].(*interface{})) +")' >Details</a>"
				}
				if conf["ActionCol_assignment"] =="true"{	
					AssignmentAction =`<a href ="#`+ conf["ActionCol_assignment_url"]+PrintValue(vals[i].(*interface{})) +`" onclick="func_assignment('`+conf["ActionCol_assignment_url"]+PrintValue(vals[hkey_action_param].(*interface{})) +`',this)" >Assignment</a>`
				}
				if conf["ActionCol_members"] =="true"{	//on hris->assignment module button on members on assignment 
					MembersAction =`<a href ="#`+ conf["ActionCol_members_url"]+PrintValue(vals[i].(*interface{})) +`" onclick="func_assignment('`+conf["ActionCol_members_url"]+PrintValue(vals[hkey_action_param].(*interface{})) +`',this)" >Members</a>`
				}

				a  = append(a,EditAction+" "+DeleteAction+" "+AddChildAction +" "+ AddroleAction +" "+ AddDetailAction + " " +AssignmentAction + "  " + MembersAction)
			}
			a  = append(a,"<div "+fontcolor+">" + PrintValue(vals[i].(*interface{})) + "</div>"  )
		}
		m = append( m,a)
		//mars at work here
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
		a = fmt.Sprint(v.Format("2006-01-02 15:04:05.999"))

	default:
//		fmt.Print(v)
		a= fmt.Sprint(v)
	}
	return a
}



func DataList(sql string ) [][]string {
	
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
			a  = append(a, PrintValue(vals[i].(*interface{}))  )
			
		}
		m = append( m,a)
		//mars at work here
	}
		
	if rows.Err() != nil {
		fmt.Println(rows.Err)
	}	
	return m

}

func Data_row(sql string ) []string {
	
	var m  []string
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
			a  = append(a, PrintValue(vals[i].(*interface{}))  )
			
		}
		m = a
		//mars at work here
	}
		
	if rows.Err() != nil {
		fmt.Println(rows.Err)
	}	
	return m

}