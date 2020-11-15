package account_role

import (
	
		"net/http"
		 "fmt"
		"text/template"
		datatables "../../datatables"
		 "strconv"
		config "../../config"
		login "../../login"

)

func  AccountRoleHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	

	if r.Method =="GET" {	
		
		rights :=r.URL.Query().Get("rights")

		for key , value := range r.URL.Query().Get("ID"){
					fmt.Println("sa key",key)
					fmt.Print("    sa value",value)
				}
		getid,_:=strconv.Atoi(r.URL.Query().Get("ID"))
		if getid > 0 { // start roles button handler
			contentconf := make(map[string]string)
			contentconf["MemberID"] =r.URL.Query().Get("ID")
		    contentconf["postURL"] = "/administrator/accountrole" // on submit url
			type Sysrole struct{
				Id int
				Name string
				Isdata string
			}
			var Lastname,Firstname,Middlename,Username string
			rows ,err, _,_ :=config.Ap_sql("exec account_roles_peraccount "+r.URL.Query().Get("ID") ,1)
			exportedSysrole := []Sysrole{}
		    for rows.Next() {
		            var r Sysrole
		            err = rows.Scan(&r.Id, &r.Name,&r.Isdata)
		            if err != nil {
		                    panic(err.Error())
		            }
					data := Sysrole{
		            	r.Id, 
		            	r.Name,
		            	r.Isdata,
				    }
					exportedSysrole = append(exportedSysrole, data)
		    }
		    memberrows ,err, _,_ :=config.Ap_sql("select lastname, firstname,middlename,Username from member inner join orgMem on orgMem.member=member.id where orgMem.id= "+r.URL.Query().Get("ID"),1)
			for memberrows.Next() {
		            err = memberrows.Scan(&Lastname, &Firstname,&Middlename,&Username)
		            if err != nil {
		                    panic(err.Error())
		            }
		         	contentconf["lastname"] = Lastname
		         	contentconf["firstname"] =Firstname
		         	contentconf["middlename"] =Middlename
		         	contentconf["username"] =Username
		         	contentconf["Rights"] =rights
		         	
					
		    }
		    type Contentdata struct {
		    	SysroleData []Sysrole
		    	Contentcon map[string]string
		    }
		    tmpl := template.New("account_role.html")
			//var err error
			if tmpl, err = tmpl.ParseFiles("admin/account_role/account_role.html"); err != nil {
				fmt.Println(err)
			}
			err1 := tmpl.Execute(w,&Contentdata{exportedSysrole,contentconf})
			if err1 != nil {
				http.Error(w, err1.Error(), http.StatusInternalServerError)
			}		
		
		//end roles button handler
		}else{
				

			tconf := make(map[string]string)
			ln :=  r.URL.Query().Get("length")
			tconf["Panel_name"] ="Account Role" 
			tconf["Add_new_button"] ="false"  
			tconf["Add_new_button_url"] ="administrator/accountrole"
			tconf["Action_baseURL"] ="/administrator/accountrole" // on clicking action column 
			var header = []string{"Action", "ID", "username", "lastname", "middlename", "firstname", "maidenname", "suffixname"}  //column header
			tconf["ActionCol_param"] ="ID"
			if ln=="" {
				var ajaxURLdata = "/administrator/accountrole"
				type TblConf struct{
					Headercol []string
					AjaxUrl string
					Tempconf map[string]string
				}
				tmpl := template.New("table.html")
				var err error
				//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
				if tmpl, err = tmpl.ParseFiles("admin/account_role/table.html"); err != nil {
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
					sortcol="2"
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
			 	tconf["ActionCol_edit"] ="false" // config for edit click 
			 	tconf["ActionCol_delete"] ="false" // config for delete click 
			 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
			 	tconf["ActionCol_roles"] ="true" // config for roles
			 	tconf["ActionCol_roles_url"] ="/administrator/accountrole" // config for roles
			 	tconf["ActionCol_roles_account_edit"] ="false"
			 	tconf["style_size"] ="12"
			 	tconf["style_font"] =""
			 	tconf["format_number_col"] ="Amount"
				tconf["dr"]=dr  //draw
				//geting total record etc 
				fmt.Println("exec member_list 1, "+str+", "+ln+", 2, '"+sortColName+"', '', '"+sr+"'")
				tconf["sql_total"] ="select 1 as Page_current,1 as Page_more, 1 as Page_total ,1 as Record_more , 1 as Record_total"	
				tconf["sql_data"] ="exec account_roles"
				datatables.DatatableretArray(w,tconf)
			}	  	  
		}		
	}else{
		session_username,_ := login.Get_account_info(r)	
		r.ParseForm()
		/*r.Form["firstname"][0]
		r.Form["lastname"][0]
		r.Form["member_id"][0]
		r.Form["middlename"][0]
		r.Form["username"][0]*/
		fmt.Println("delete orgmemrol where  orgmem ="+r.Form["member_id"][0])
		_ ,err, _,_ := config.Ap_sql("delete orgmemrol where  orgmem ="+r.Form["member_id"][0],3)
			if err!= nil {
			panic(err.Error())
			}
		for _ , value := range r.Form["sysroleid"] {
			fmt.Println("insert into orgmemrol(sysrole,orgmem,lasttouch) Values("+value+","+r.Form["member_id"][0]+",'"+session_username+"')")
			_ ,err, _,_ := config.Ap_sql("insert into orgmemrol(sysrole,orgmem,lasttouch) Values("+value+","+r.Form["member_id"][0]+",'"+session_username+"')",3)
			if err!= nil {
			panic(err.Error())
			}
		}



	}
}


