package role

import (
	
		"net/http"
		 "fmt"
		"text/template"
		 "strconv"
		datatables "../../datatables"
		config "../../config"
		login "../../login"
)	


func  RoleHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	tconf := make(map[string]string)
	ln :=  r.URL.Query().Get("length")
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)
	tconf["Panel_name"] ="Access Role" 
	tconf["Add_new_button"] ="true"  
	tconf["Add_new_button_url"] ="/administrator/createrole"
	tconf["ActionCol_param"] ="ID"
	tconf["ActionCol_edit"] ="true"
	tconf["Action_baseURL"] ="/administrator/createrole" // on clicking action column 
	
	var header = []string{"Action", "ID", "Name", "Remarks","Lasttouch","Lasttouched"}  //column header
	getid,_:=strconv.Atoi(r.URL.Query().Get("ID"))
	if getid > 0 { 
		tconf["MemberID"] =r.URL.Query().Get("ID")
	    type Sysrole struct{
			Id int
			Name string
			Isdata string
		}


	}else if ln=="" {
		var ajaxURLdata = "/administrator/role"
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
		if tmpl, err = tmpl.ParseFiles("admin/role/table.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		dr := r.URL.Query().Get("draw")
		sr := r.URL.Query().Get("search[value]")
		srt1,_ :=strconv.Atoi(r.URL.Query().Get("start"))
	    ln1,_ :=strconv.Atoi(r.URL.Query().Get("length")) 
	    srt1 = (srt1 / ln1) + 1
	    str := strconv.Itoa(srt1)
	    sortcol :=  r.URL.Query().Get("order[0][column]")	
		if sortcol=="0"{
			sortcol="3"
		}
		srtcolindexInt,_ :=strconv.Atoi(sortcol)
	    var sortColName string
		for key , value := range header{
			if ( srtcolindexInt==key){
				sortColName = value
			}
		}
		tconf["ActionCol"] ="true" // config for action column
	 	tconf["ActionCol_param"] ="ID" // config for parameter of action 
	 	tconf["ActionCol_edit"] ="true" // config for edit click 
	 	tconf["ActionCol_edit_url"] ="/administrator/createrole?ID="
		tconf["ActionCol_edit_is_modal"] ="false" // config for edit click 
		tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
		tconf["dr"]=dr  //draw
		//geting total record etc 
		fmt.Println("exec LBR_Assign_List 1, "+str+", "+ln+", 2, 1, '"+sortColName+"', '', '"+sr+"'	")
		fmt.Println("exec LBR_Assign_List 1,	")
		tconf["sql_total"] = "select 1 as Page_current,1 as Page_more, 1 as Page_total ,1 as Record_more , 1 as Record_total "	
		tconf["sql_data"] = `SysRole_Get 3, `+str_OrgID 
		datatables.DatatableretArray(w,tconf)  //datatables return array data
	 
		
	}	    

}


func CreateRoleHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)

	if r.Method =="GET" {	
		//getid,_:=strconv.Atoi(r.URL.Query().Get("ID"))
	/*	if getid > 0 {

		}else{*/

			contentconf := make(map[string]string)
			contentconf["roleid"]= r.URL.Query().Get("ID")
			idinInt,_:=strconv.Atoi(r.URL.Query().Get("ID")) 
			contentconf["Panel_name"] ="Update Role" 
			if idinInt == 0{
				contentconf["roleid"]= "0"
				contentconf["Panel_name"] ="Create Role" 
				contentconf["rolename"] = ""
			}
			fmt.Println("naka string pa",contentconf["roleid"]);
			fmt.Println("naka int na",idinInt);
			contentconf["postURL"] = "/administrator/createrole"
			
			type Sysaccess struct {
				ID int64
				Name string
				Url string
				Rights string
				Ischeck string
			}
			var rolename string
			rows ,err, _,_ := config.Ap_sql("select name from sysrole where id="+ contentconf["roleid"],1)	
			//rows ,err, _,_ := config.Ap_sql("select id, name, url from sysaccess" ,1)	
			 if err != nil {
		    	panic(err.Error())
		    }
		    for rows.Next() {
		            err = rows.Scan(&rolename)
		            if err != nil {
		                    panic(err.Error())
		            }
					contentconf["rolename"] = rolename	
			}
			Sysaccessrows ,err, _,_ := config.Ap_sql("exec SysAccess_Get 3,"+contentconf["roleid"],1)
			if err != nil {
		    	panic(err.Error())
		    }
		    SysaccArr := []Sysaccess{}
		    var s Sysaccess
		    for Sysaccessrows.Next() {
		            err = Sysaccessrows.Scan(&s.ID,&s.Name,&s.Url,&s.Rights,&s.Ischeck)
		            if err != nil {
		                panic(err.Error())
		            }
					data := Sysaccess{
		            	s.ID,
		            	s.Name,
		            	s.Url,
		            	s.Rights,
		            	s.Ischeck,
				    }
					SysaccArr = append(SysaccArr, data)	
			}
			type Contentdata struct {
		    	Sysaccessdata []Sysaccess
		    	Contentcon map[string]string
		    }
			tmpl := template.New("createrole.html")
				//var err error
				if tmpl, err = tmpl.ParseFiles("admin/role/createrole.html"); err != nil {
					fmt.Println(err)
				}
				err1 := tmpl.Execute(w,&Contentdata{SysaccArr,contentconf})
				if err1 != nil {
					http.Error(w, err1.Error(), http.StatusInternalServerError)
				}
		//}		
	}else{
		//add new
		session_username,_ := login.Get_account_info(r)	
		r.ParseForm()
		roleidhtml,_:=strconv.Atoi(r.Form["roleID"][0]) 
		var lastinsertedIDString string
		if roleidhtml > 0{
			_ ,err, _,_:= config.Ap_sql("update sysrole set name = '"+r.Form["rolename"][0]+"',lasttouch='"+session_username+"' where id ="+r.Form["roleID"][0],3)
			if err!= nil {
				panic(err.Error())
			}
			_ ,err, _, _= config.Ap_sql("delete from sysrolacc where sysrole="+r.Form["roleID"][0],3)
			if err!= nil {
				panic(err.Error())
			}
			lastinsertedIDString = r.Form["roleID"][0]
		}else{
			fmt.Println("insert into sysrole(name,remarks,lasttouch) values('"+r.Form["rolename"][0]+"','','"+session_username+"')")
			_ ,err, _,res:= config.Ap_sql("insert into sysrole(name,remarks,lasttouch) values('"+r.Form["rolename"][0]+"','','"+session_username+"')" ,3)
			if err!= nil {
				panic(err.Error())
			}
			lastinsertedID, _ := res.LastInsertId()
			lastinsertedIDString  = strconv.FormatInt(lastinsertedID,10)
		}
		for key ,val := range r.Form["sysaccessid"] {
			fmt.Println("insert into sysrolacc(sysrole,sysaccess,rights,lasttouch) values("+lastinsertedIDString+","+val+",'"+r.Form["rights["+val+"]"][0]+"','"+session_username+"')")
			var	err error
			_ ,err, _,_ = config.Ap_sql("insert into sysrolacc(sysrole,sysaccess,rights,lasttouch) values("+lastinsertedIDString+","+val+",'"+r.Form["rights["+val+"]"][0]+"','"+session_username+"')",3)	
        	if err!= nil {
				panic(err.Error())
			}
        	fmt.Println("key-mao:", key)
        	fmt.Println("value mao ni" , val)
        //end add new 


        }

	}			
	
}