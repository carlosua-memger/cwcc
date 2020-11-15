package members

import (
	
		"net/http"
		"fmt"
		"text/template"
		config "../../config"
		login "../../login"
		"encoding/json"
		"strconv"	
		datatables "../../datatables"
		"time"
		"strings"
)	

type members_list struct {
		id int
		username string
		firstname string
		lastname string	
		middlename string
		suffixname string
		birthdate string
}
type Profile struct {
	  	Message string  // error status output or validated data
	  	Status bool //pass or fail
}
func AllAccountHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	last_start :=r.URL.Query().Get("last_start")
	last_length :=r.URL.Query().Get("last_length")
	last_search :=r.URL.Query().Get("last_search")
	dec_rights := config.Decrypt(config.Encrypt_key,r.URL.Query().Get("rights"))
	tconf := make(map[string]string)
	ln :=  r.URL.Query().Get("length")
	tconf["Panel_name"] ="Members"   //panel name
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="false"
	tconf["Add_new_button"] ="false" 
	tconf["level"] ="1" 
	tconf["pageLength"] =last_length 
	tconf["displayStart"] =last_start
	tconf["last_search"] = last_search
	if last_length==""{
		tconf["pageLength"] ="10" 
		tconf["displayStart"] ="0"
	}
	/*if last_search==""{
		tconf["last_search"] = "0" 
	}*/
	tconf["Add_new_button_url"] ="/Member/CreateMember?rights="+rights
	var header = []string{"Action", "recno", "ID","Employee #", "username", "lastname", "middlename", "firstname", "maidenname", "suffixname", "status", "gender", "tin","sss","mobileno","email","enable","disable","lasttouched","remarks"}  //column header
	// rights for tables 	
	if( strings.Contains(dec_rights, "L1add") ){
		tconf["Add_new_button"] ="true"
	}
	//end rights for tables
 
	if ln=="" {
		var ajaxURLdata = "/member?rights="+rights
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("hris/members/table.html"); err != nil {
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
	 	tconf["ActionCol_edit_is_modal"] ="true" // config for edit click 
	 	tconf["ActionCol_edit_url"] ="/Member/EditMember?rights="+rights+"&ID=" // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/Member/DeleteMember?ID="
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	tconf["ActionCol_assignment"] ="false"
	 	tconf["ActionCol_assignment_url"]="/Member/MemberAssignment?rights="+rights+"&ID="
	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
		tconf["dr"]=dr  //draw
		//geting total record etc 
		
		// rights for tables 	
		if( strings.Contains(dec_rights, "L1edit") ){
			tconf["ActionCol_edit"] ="true" // config for edit click 
		}
		if( strings.Contains(dec_rights, "L1delete") ){
			tconf["ActionCol_delete"] ="true" // config for delete click 
		}
		if( strings.Contains(dec_rights, "L1assignment") ){
			tconf["ActionCol_assignment"] ="true" // config for delete click 
		}
		//end rights for tables
		fmt.Println("exec member_list 1, "+str+", "+ln+", 2, '"+sortColName+"', '', '"+sr+"'")
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		tconf["sql_total"] = "exec member_list 1, "+str+", "+ln+", "+str_OrgID+", '"+sortColName+"', '', '"+sr+"'"	
		tconf["sql_data"] ="exec member_list 0, "+str+", "+ln+", "+str_OrgID+", '"+sortColName+"', '', '"+sr+"'"
		datatables.DatatableretArray(w,tconf)
		  	  
		
	}	    

}

type RetMessage struct{
	RetNo int
	RetMsg string
}
func DoCreateMember(org string, gender string, firstname string, middlename string, lastname string, suffixname string ,maidenname string, birthday string, address string,postcode string,TIN string,SSS string,GSIS string, Pagibig string, PHealth string,mobileno string,telno string,email string,civil_stat string,remarks string, hris string )RetMessage{
	rowret ,err, _,_ := config.Ap_sql("Member_Save 'Add','1','0',"+org+",'','"+lastname+"','"+middlename+"','"+firstname+"','"+suffixname+"','"+maidenname+"','"+gender+"','"+birthday+"','"+address+"','"+postcode+"','"+TIN+"','"+SSS+"','"+GSIS+"','"+Pagibig+"','"+PHealth+"','"+mobileno+"','"+telno+"','"+email+"','"+civil_stat+"','"+remarks+"','"+hris+"' ",1)
	if err != nil {
		panic(err.Error())
	}
	var r RetMessage
	for rowret.Next() {
	    err = rowret.Scan(&r.RetNo,&r.RetMsg)
	    if err != nil {
	        panic(err.Error())
	    }
	 	r = RetMessage{r.RetNo,r.RetMsg}
	}
	return r
}
func CreateMemberHAndler(w http.ResponseWriter, r *http.Request) {
	rights :=r.URL.Query().Get("rights")
	login.Session_renew(w,r)
	_,session_userID := login.Get_account_info(r)
	str_userID := strconv.Itoa(session_userID)
	Org_id :=login.Get_session_org_id(r)
	str_OrgID :=strconv.Itoa(Org_id)

	if r.Method =="GET" {
		
		rows ,err, _,_ := config.Ap_sql("SysDropDown_Get 'MemberGender'",1)
	    if err != nil {
	    	panic(err.Error())
	    }

	    type row struct {
	           
	            Name  string
	            
	    }     
		got := []row{}
		    for rows.Next() {
		            var r row
		            err = rows.Scan(&r.Name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post := row{
		            	r.Name,
				        
				    }
					got = append(got, post)
				}
		rows3 ,err, _,_ := config.Ap_sql("SysDropDown_Get 'MemberCivilStatus'",1)
	    if err != nil {
	    	panic(err.Error())
	    }

	    type row3 struct {
	           
	            Name  string
	            
	    }     
		got3 := []row3{}
		    for rows3.Next() {
		            var r row3
		            err = rows3.Scan(&r.Name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post3 := row3{
		            	r.Name,
				        
				    }
					got3 = append(got3, post3)
				}

		// new struct for Org
		rows2 ,err, _,_ := config.Ap_sql(`exec orgmem_get 1, `+str_OrgID+`,`+str_userID+`,''` ,1)
	    if err != nil {
	    	panic(err.Error())
	    }

	    type row2 struct {
	           Org_id int
	            Org_name  string
	       	            
	    }      
		got2 := []row2{}
		    for rows2.Next() {
		            var r row2
		            err = rows2.Scan(&r.Org_id, &r.Org_name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post2 := row2{
		            	r.Org_id,
				        r.Org_name,
				      
				    }
					got2 = append(got2, post2)
		    }
		
	   
		    
		type Data struct {
			MemberGender []row
			MemberCivilStatus []row3
			Org []row2 
			Rights string
			
		}    
     	tmpl := template.New("CreateMember.html")
		//var err error
		if tmpl, err = tmpl.ParseFiles("hris/members/CreateMember.html"); err != nil {
			fmt.Println(err)
		}

		err1 := tmpl.Execute(w,&Data{got,got3, got2,rights})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else {
		r.ParseForm()
			bday :=r.Form["bday_year"][0] + "/" + r.Form["bday_month"][0] + "/" + r.Form["bday_day"][0]
			if(r.Form["postcode"][0]==""){
				r.Form["postcode"][0] = "1"
			}

			retmsg := DoCreateMember(r.Form["org"][0],r.Form["gender"][0], r.Form["firstname"][0],r.Form["middlename"][0],r.Form["lastname"][0],r.Form["maidenname"][0],r.Form["suffixname"][0],bday,r.Form["TIN"][0],r.Form["SSS"][0],r.Form["GSIS"][0],r.Form["Pagibig"][0],r.Form["PHealth"][0],r.Form["mobileno"][0],r.Form["telno"][0],r.Form["email"][0],r.Form["address"][0],r.Form["postcode"][0],r.Form["civil_stat"][0],r.Form["remarks"][0],r.Form["hris"][0]) 
			profile := Profile{retmsg.RetMsg,true}
				js, err := json.Marshal(profile)
				if err != nil {
			    	http.Error(w, err.Error(), http.StatusInternalServerError)
			    	return
		  		}
			  	w.Header().Set("Content-Type", "application/json")
			  	w.Write(js)

			/*if(DoCreateMember(r.Form["org"][0],r.Form["gender"][0], r.Form["firstname"][0],r.Form["middlename"][0],r.Form["lastname"][0],r.Form["maidenname"][0],r.Form["suffixname"][0],bday,r.Form["TIN"][0],r.Form["SSS"][0],r.Form["GSIS"][0],r.Form["Pagibig"][0],r.Form["PHealth"][0],r.Form["mobileno"][0],r.Form["telno"][0],r.Form["email"][0],r.Form["address"][0],r.Form["postcode"][0],r.Form["civil_stat"][0],r.Form["remarks"][0]) ){
				profile := Profile{"Data Successfully added",true}
				js, err := json.Marshal(profile)
				if err != nil {
			    	http.Error(w, err.Error(), http.StatusInternalServerError)
			    	return
		  		}
			  	fmt.Println("Data Successfully added")
			  	w.Header().Set("Content-Type", "application/json")
			  	w.Write(js)
			  		
			}else{
				profile := Profile{"Oppsss something went wrong",false}
				js, err := json.Marshal(profile)
				if err != nil {
			    	http.Error(w, err.Error(), http.StatusInternalServerError)
			    	return
		  		}
			  	fmt.Println("Oppsss something went wrong")
			  	w.Header().Set("Content-Type", "application/json")
			  	w.Write(js)	
			}*/

			
		//}
	}

}

func DoEditMemberHAndler(w http.ResponseWriter, r *http.Request){
	login.Session_renew(w,r)
	session_username,_ := login.Get_account_info(r)
	r.ParseForm()
	
	bday := r.Form["bday_year"][0]+"/"+r.Form["bday_month"][0]+"/"+r.Form["bday_day"][0]
	//bday := "2014-1-1"
	//                           fmt.Println("Member_Save 'Edit','"+login.Uname+"',"+r.Form["rec_no"][0]+","+r.Form["org"][0]+",'username Dare','"+r.Form["lastname"][0]+"','"+r.Form["middlename"][0]+"','"+r.Form["firstname"][0]+"','"+r.Form["maidenname"][0]+"','"+r.Form["suffixname"][0]+"','"+r.Form["gender"][0]+"','"+bday+"','"+r.Form["TIN"][0]+"','"+r.Form["SSS"][0]+"','"+r.Form["GSIS"][0]+"','"+r.Form["Pagibig"][0]+"','"+r.Form["PHealth"][0]+"','"+r.Form["mobileno"][0]+"','"+r.Form["telno"][0]+"','"+r.Form["email"][0]+"','"+r.Form["address"][0]+"',"+r.Form["postcode"][0]+",'"+r.Form["civil_stat"][0]+"','"+r.Form["remarks"][0]+"' ")
	//_ ,err, ex_stat,_ :=config.Ap_sql("Member_Save 'Edit','1','"+r.Form.Get("rec_no")+"','','"+r.Form.Get("username")+"','"+r.Form.Get("lastname")+"','"+r.Form.Get("middlename")+"','"+r.Form.Get("firstname")+"','"+r.Form.Get("maidenname")+"','"+r.Form.Get("suffixname")+"','"+r.Form.Get("gender")+"','"+r.Form.Get("birthday")+"','"+r.Form.Get("TIN")+"','"+r.Form.Get("SSS")+"','"+r.Form.Get("GSIS")+"','"+r.Form.Get("Pagibig")+"','"+r.Form.Get("PHealth")+"','"+r.Form.Get("mobileno")+"','"+r.Form.Get("telno")+"','"+r.Form.Get("email")+"','"+r.Form.Get("address")+"','"+r.Form.Get("postcode")+"','"+r.Form.Get("civil_stat")+"','"+r.Form.Get("remarks")+"',0 ",3)
	//	_ ,err, ex_stat,_ := config.Ap_sql("Member_Save 'Edit','"+login.Uname+"',"+r.Form["rec_no"][0]+","+r.Form["org"][0]+",'"+r.Form["username"][0]+"','"+r.Form["lastname"][0]+"','"+r.Form["middlename"][0]+"','"+r.Form["firstname"][0]+"','"+r.Form["maidenname"][0]+"','"+r.Form["suffixname"][0]+"','"+r.Form["gender"][0]+"','"+bday+"','"+r.Form["TIN"][0]+"','"+r.Form["SSS"][0]+"','"+r.Form["GSIS"][0]+"','"+r.Form["Pagibig"][0]+"','"+r.Form["PHealth"][0]+"','"+r.Form["mobileno"][0]+"','"+r.Form["telno"][0]+"','"+r.Form["email"][0]+"','"+r.Form["address"][0]+"',"+r.Form["postcode"][0]+",'"+r.Form["civil_stat"][0]+"','"+r.Form["remarks"][0]+"' ",1)
	
	

	rowret ,err, _,_ := config.Ap_sql("Member_Save 'Edit','"+session_username+"',"+r.Form["rec_no"][0]+","+r.Form["org"][0]+",'"+r.Form["username"][0]+"','"+r.Form["lastname"][0]+"','"+r.Form["middlename"][0]+"','"+r.Form["firstname"][0]+"','"+r.Form["maidenname"][0]+"','"+r.Form["suffixname"][0]+"','"+r.Form["gender"][0]+"','"+bday+"','"+r.Form["TIN"][0]+"','"+r.Form["SSS"][0]+"','"+r.Form["GSIS"][0]+"','"+r.Form["Pagibig"][0]+"','"+r.Form["PHealth"][0]+"','"+r.Form["mobileno"][0]+"','"+r.Form["telno"][0]+"','"+r.Form["email"][0]+"','"+r.Form["address"][0]+"',"+r.Form["postcode"][0]+",'"+r.Form["civil_stat"][0]+"','"+r.Form["remarks"][0]+"','"+r.Form["hris"][0]+"' ",1)
	if err != nil {
		panic(err.Error())
	}
	var rtmsg RetMessage
	for rowret.Next() {
	    err = rowret.Scan(&rtmsg.RetNo,&rtmsg.RetMsg)
	    if err != nil {
	        panic(err.Error())
	    }
	 	rtmsg = RetMessage{rtmsg.RetNo,rtmsg.RetMsg}
	}
	
	profile := Profile{rtmsg.RetMsg,true}
	js, err := json.Marshal(profile)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
		}
  	w.Header().Set("Content-Type", "application/json")
  	w.Write(js)
}
func EditMemberHAndler(w http.ResponseWriter, r *http.Request){
	//r.ParseForm()
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	tconf := make(map[string]string)
	tconf["rec_no"] =r.URL.Query().Get("ID")
	fmt.Println(r.Form.Get("rec_no"))
	rows ,err, _,_ := config.Ap_sql("Member_Get 1,'"+tconf["rec_no"]+"'",1)
    if err != nil {
    	panic(err.Error())
    }
	 type row struct {
	            ID int
	            Username 	string
	            Password	string
	            authmethod	string
	            Lastname  	string
	            Middlename 	string
	            Firstname 	string
	            Maidenname 	string
	            Suffixname 	string
	            Gender 		string
	            Birthday 	time.Time
	            Tin 		string
	            SSS 		string
	            GSIS 		string
	            Pagibig 	string
	            Phealth 	string
	            Mobileno 	string
	            Telno 		string
	            Email 		string
	            Address 	string
	            Postcode 	string
	            Civilstatus string
	            Remarks 	string
	            enable 		time.Time
	            disable 	time.Time
	            lasttouch 	string
	            lasttouched time.Time
	            Bmonth 		string
	            Byear 		string
	            Bday 		string
	            oldID 		string
	            Hris_id		string
	          /*  Enable string
	          	Disable string
	            Lasttouched string*/
	            //Remarks string

	    }
	var a,b,c,d string
	got := []row{}
	    for rows.Next() {
	            var r row
	            err = rows.Scan(&r.ID,&r.Username, &r.Password, &r.authmethod, &r.Lastname, &r.Middlename, &r.Firstname, &r.Maidenname, &r.Suffixname,&r.Gender, &r.Birthday,  &r.Tin, &r.SSS, &r.GSIS,&r.Pagibig,&r.Phealth, &r.Mobileno,&r.Telno, &r.Email,&r.Address,&r.Postcode,&r.Civilstatus,&r.Remarks,&r.enable,&r.disable,&r.lasttouch,&r.lasttouched,&r.oldID,&r.Hris_id)
	            if err != nil {
	                    panic(err.Error())
	            }
	            a  = strconv.Itoa( int(r.Birthday.Month())  ) 
	            b  = strconv.Itoa(r.Birthday.Year()) 
	            c  = strconv.Itoa(r.Birthday.Day()) 
	           /* d = r.Birthday*/
				fmt.Println("bday",a)
				fmt.Println("bmon",b)
				fmt.Println("byear",c)
				post := row {
	               	r.ID,
	            	r.Username, 
	            	r.Password,
	            	r.authmethod,
		            r.Lastname , 
		            r.Middlename ,
		            r.Firstname ,
		            r.Maidenname ,
		            r.Suffixname ,
		            r.Gender ,
		            r.Birthday, 
		            r.Tin ,
		            r.SSS ,
		            r.GSIS ,
		            r.Pagibig ,
		            r.Phealth ,
		            r.Mobileno ,
		            r.Telno ,
		            r.Email ,
		            r.Address, 
		            r.Postcode, 
		            r.Civilstatus, 
		            r.Remarks ,
		            r.enable ,
		            r.disable ,
		            r.lasttouch, 
		            r.lasttouched ,
		            a,
		            b,
		            c, 
					r.oldID,
		            r.Hris_id, 
		            

			    }
				got = append(got, post)

		}
	    //months
	    type months struct{  
	    	Id string
	    	Month string
	    }
	    //var d1 months
	    month := []months{}
	    d1 := months{"1","JAN"}  
	   	month = append(month,d1)
	   	d1 = months{"2","FEB"}  
	   	month = append(month,d1)
	   	d1 = months{"3","MAR"}  
	   	month = append(month,d1)
	   	d1 = months{"4","APR"}  
	   	month = append(month,d1)
	   	d1 = months{"5","MAY"}  
	   	month = append(month,d1)
	   	d1 = months{"6","JUN"}  
	   	month = append(month,d1)
	   	d1 = months{"7","JUL"}  
	   	month = append(month,d1)
	   	d1 = months{"8","AUG"}  
	   	month = append(month,d1)
	   	d1 = months{"9","SEP"}  
	   	month = append(month,d1)
	   	d1 = months{"10","OCT"}  
	   	month = append(month,d1)
	   	d1 = months{"11","NOV"}  
	   	month = append(month,d1)
	   	d1 = months{"12","DEC"}  
	   	month = append(month,d1)
	   	fmt.Println(month)
	   	// end months
		fmt.Println("boom data",a)
	    fmt.Println("boom data",b)
	    fmt.Println("boom data",c)
	    fmt.Println("boom data",d)
	    var monthday int
	    type Monthdays struct {
	    	Day string
	    }
	    a = `11`
	    switch a {
	    	case "1" :
	    		monthday = 31	
	    	case "2" : 
	    		monthday = 28
	    	case "3" :
	    		monthday = 31
	    	case "4" :
	    		monthday = 30
	    	case "5" :
	    		monthday = 31
	    	case "6" :
	    		monthday = 30
	    	case "7" :
	    		monthday = 31
	    	case "8" :
	    		monthday = 31
	    	case "9" :
	    		monthday = 30
	    	case "10" :
	    		monthday = 31
	    	case "11" :
	    		monthday = 30
	    	case "12" :
	    		monthday = 31
	    	default :	
	    } 
	    mondays := []Monthdays{}
	    for i := 1; i <= monthday ; i++ {
	    	dy  := Monthdays{ strconv.Itoa(i)}
	    	mondays = append(mondays, dy)
	    }

	    rows3 ,err, _,_ := config.Ap_sql("SysDropDown_Get 'MemberGender'",1)
	    if err != nil {
	    	panic(err.Error())
	    }

	    type row3 struct {
	            Name  string
	          
	            
	    }     
		got3 := []row3{}
		    for rows3.Next() {
		            var r row3
		            err = rows3.Scan(&r.Name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post3 := row3{
		            	r.Name,
				        
				    }
					got3 = append(got3, post3)
				}

		rows4 ,err, _,_ := config.Ap_sql("SysDropDown_Get 'MemberCivilStatus'",1)
	    if err != nil {
	    	panic(err.Error())
	    }

	    type row4 struct {
	           
	            Name  string
	            
	    }     
		got4 := []row4{}
		    for rows4.Next() {
		            var r row4
		            err = rows4.Scan(&r.Name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post4 := row4{
		            	r.Name,
				        
				    }
					got4 = append(got4, post4)
				}		// new struct for Org
		

				// new struct for Org
		_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 

		rows2 ,err, _,_ := config.Ap_sql("exec orgmem_get 1, 0,"+session_user_id_str,1)
	    if err != nil {
	    	panic(err.Error())
	    }

	    type row2 struct {
	           Org_id int
	            Org_name  string
	       	            
	    }     
		got2 := []row2{}
		    for rows2.Next() {
		            var r row2
		            err = rows2.Scan(&r.Org_id, &r.Org_name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post2 := row2{
		            	r.Org_id,
				        r.Org_name,
				      
				    }
					got2 = append(got2, post2)
		    }

		type Data2 struct {
			MemberList []row
			MemberGender []row3
			Org []row2
			Months []months
			Monthdays []Monthdays
			MemberCivilStatus []row4
			Rights string
		}    
	var templates101 = template.Must(template.ParseFiles("hris/members/EditMember.html"))
		err1 := templates101.Execute(w,&Data2{got , got3 ,got2,month,mondays,got4,rights})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}    
}
/*type Error_mess struct {
	RetNo	int
    RetMsg	string
	       	            
}
*/
func DeleteMemberHAndler(w http.ResponseWriter, r *http.Request){
	login.Session_renew(w,r)
	if r.Method =="GET" {
		tconf := make(map[string]string)
		tconf["RecID"] = r.URL.Query().Get("ID")
	   	tconf["delete_url"] = "/Member/DeleteMember"
	   	
	   	tmpl := template.New("modal_delete.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("hris/members/modal_delete.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,tconf)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		session_username,_ := login.Get_account_info(r)
		
		r.ParseForm()	
		fmt.Println("Member_Save 'Delete','"+session_username+"' ,"+r.Form["rec_no"][0]+"")
		//_ ,err1, ex_stat,_ := config.Ap_sql("Member_Save 'Delete','"+login.Uname+"' ,"+r.Form["rec_no"][0]+"",3)
  		rows2 ,err, _,_ := config.Ap_sql("Member_Save 'Delete','"+session_username+"' ,"+r.Form["rec_no"][0]+"",1)
  		if err != nil {
	    	panic(err.Error())
	    }
		var errorMsg  Error_mess
	    for rows2.Next() {
            var r Error_mess
            err = rows2.Scan(&r.RetNo, &r.RetMsg)
            if err != nil {
                    panic(err.Error())
            }
			errorMsg = Error_mess{r.RetNo,r.RetMsg}
		      
		}
				
		
	    
		js, err := json.Marshal(errorMsg)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	fmt.Println("Data Successfully added")
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)

	}		

		
}


type Error_mess struct {
	RetNo	int
    RetMsg	string
	       	            
}

type LBR_Memass_list struct {
	Recno 			interface{}
	Id 				interface{}
	Lbr_assign 		interface{}
	Lbr_assign_name	interface{}
	orgmem			interface{}
	Member_name		interface{}
	Processaccess	interface{}
	Txdate			interface{}
	Device_id		interface{}
	Remarks			interface{}
	Enable			interface{}
	Disable 		interface{}
}

func MemberAssigmentHAndler(w http.ResponseWriter, r *http.Request){

	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	dec_rights := config.Decrypt(config.Encrypt_key,r.URL.Query().Get("rights"))
	memberID :=r.URL.Query().Get("ID") 

	tconf := make(map[string]string)
	ln :=  r.URL.Query().Get("length")
	tconf["Panel_name"] =`Member Assigment [ `+ memberID +` ]`   //panel name
	tconf["Rights"] = rights
	tconf["Add_back_button"] ="true"
	tconf["Add_new_button"] ="false"  
	tconf["level"] ="2" 
	tconf["pageLength"] ="10" 
	tconf["displayStart"] ="0"
	tconf["last_search"]=""
	tconf["Add_new_button_url"] ="/Member/MemberAssignmentCreate?rights="+rights+"&memberID="+memberID
	var header = []string{"Action", "Recno", "ID", "Employee #", "Leader", "Lbr_assign", "Lbr_assign_name", "orgmem", "Member_name","Processaccess", "Txdate", "Device_id", "Remarks", "Enable", "Disable"}  //column header
	// rights for tables 	
	if( strings.Contains(dec_rights, "MAadd") ){
		tconf["Add_new_button"] ="true"
	}
	//end rights for tables

	if ln=="" {
		var ajaxURLdata = "/Member/MemberAssignment?rights="+rights +"&ID="+memberID
		type TblConf struct{
			Headercol []string
			AjaxUrl string
			Tempconf map[string]string
		}
		tmpl := template.New("table.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("hris/members/table.html"); err != nil {
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
		tconf["ActionCol_assignment"] ="false" 
		tconf["ActionCol"] ="true" // config for action column
	 	tconf["ActionCol_param"] ="ID" // config for parameter of action 
	 	tconf["ActionCol_edit"] ="false" // config for edit click 
	 	tconf["ActionCol_edit_is_modal"] ="false" // config for edit click 
	 	tconf["ActionCol_edit_url"] ="/Member/MemberAssignmentEdit?rights="+rights+"&ID=" // config for edit click 
	 	tconf["ActionCol_delete"] ="false" // config for delete click 
	 	tconf["ActionCol_delete_url"] ="/Member/MemberAssignmentDelete?ID="
	 	tconf["ActionCol_add_child"] ="false" // config for addchild click 
	 	//tconf["ActionCol_assignment"] ="false"
	 	tconf["style_size"] ="12"
	 	tconf["style_font"] =""
	 	tconf["format_number_col"] ="Amount"
		tconf["dr"]=dr  //draw
		//geting total record etc 
		
		// rights for tables 	
		if( strings.Contains(dec_rights, "MAedit") ){
			tconf["ActionCol_edit"] ="true" // config for edit click 
		}
		if( strings.Contains(dec_rights, "MAdelete") ){
			tconf["ActionCol_delete"] ="true" // config for delete click 
		}
		//end rights for tables
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		fmt.Println("exec LBR_Memass_list 0, "+str_OrgID+", "+memberID+", 0, "+str+", "+ln+", '"+sortColName+"', '', '"+sr+"'")
		tconf["sql_total"] = "exec LBR_Memass_list 1	, "+str_OrgID+", "+memberID+", 0, "+str+", "+ln+", '"+sortColName+"', '', '"+sr+"'"	
		tconf["sql_data"] = "exec LBR_Memass_list 0, "+str_OrgID+", "+memberID+", 0, "+str+", "+ln+", '"+sortColName+"', '', '"+sr+"'"

		datatables.DatatableretArray(w,tconf)
		  	  
		
	}	    
}

func MemberAssigmentCreateHAndler(w http.ResponseWriter, r *http.Request){
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	tconf := make(map[string]string)
	tconf["rights"]=rights
	tconf["memberID"]=r.URL.Query().Get("memberID")
	if r.Method =="GET" {
		_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		LBR_Assign_Get_row ,err, _,_ := config.Ap_sql("exec LBR_Assign_Get 2, "+str_OrgID+","+ session_user_id_str,1)
	    if err != nil {
	    	panic(err.Error())
	    }

	    type LBR_Assign_Get struct {
	           
	            Id  string
	            Name string
	            
	    }     
		LBR_Assign_Get_data := []LBR_Assign_Get{}
		    for LBR_Assign_Get_row.Next() {
		            var r LBR_Assign_Get
		            err = LBR_Assign_Get_row.Scan(&r.Id,&r.Name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post := LBR_Assign_Get{
		            	r.Id,
		            	r.Name,
				        
				    }
					LBR_Assign_Get_data = append(LBR_Assign_Get_data, post)
				}
		
		    
		type Data struct {
			LBR_Assign_Get []LBR_Assign_Get
			Htdata map[string]string
			
		}   
		htdata := Data{LBR_Assign_Get_data,tconf} 
     	tmpl := template.New("MemberAssignmentCreate.html")
		//var err error
		if tmpl, err = tmpl.ParseFiles("hris/members/MemberAssignmentCreate.html"); err != nil {
			fmt.Println(err)
		}

		err1 := tmpl.Execute(w,htdata)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		//do create assigment Hire.
		r.ParseForm()

		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		session_username,_ := login.Get_account_info(r)
		sqlstr:= `exec LBR_Memass_Save 
						'Add',
						'`+ session_username + `' ,
						0,
						`+str_OrgID+`,
						` + r.Form["lbr_assign"][0] + ` ,
						` + r.Form["member"][0]+ `	 ,
						` + r.Form["processaccess"][0]+` ,
						'`+ r.Form["txdate"][0]+`' ,
						'`+ r.Form["device_id"][0]+`' ,
						'`+ r.Form["remarks"][0]+`' ,
						'`+ r.Form["enable"][0]+`' ,
						'`+ r.Form["disable"][0] + `' ,
						'`+ r.Form["leader"][0] + `' ,
						'`+ r.Form["emp_id"][0] + `' `

				
		rowret ,err, _,_ := config.Ap_sql(sqlstr,1)
		if err != nil {
			panic(err.Error())
		}
		var r RetMessage
		for rowret.Next() {
		    err = rowret.Scan(&r.RetNo,&r.RetMsg)
		    if err != nil {
		        panic(err.Error())
		    }
		 	r = RetMessage{r.RetNo,r.RetMsg}
		}

		profile := Profile{r.RetMsg,true}
		js, err := json.Marshal(profile)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)


	}	

}	




func MemberAssigmentEditHAndler(w http.ResponseWriter, r *http.Request){
	login.Session_renew(w,r)
	rights :=r.URL.Query().Get("rights")
	tconf := make(map[string]string)
	tconf["rights"]=rights
	if r.Method =="GET" {
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		tconf["MemassID"]=r.URL.Query().Get("ID")
		_,session_user_id := login.Get_account_info(r)
		session_user_id_str := strconv.Itoa(session_user_id) 

		LBR_Assign_Get_row ,err, _,_ := config.Ap_sql("exec LBR_Assign_Get 2, "+str_OrgID+","+ session_user_id_str,1)
	    if err != nil {
	    	panic(err.Error())
	    }

	    type LBR_Assign_Get struct {
	           
	            Id  string
	            Name string
	            
	    }     
		LBR_Assign_Get_data := []LBR_Assign_Get{}
		    for LBR_Assign_Get_row.Next() {
		            var r LBR_Assign_Get
		            err = LBR_Assign_Get_row.Scan(&r.Id,&r.Name)
		            if err != nil {
		                    panic(err.Error())
		            }
					post := LBR_Assign_Get{
		            	r.Id,
		            	r.Name,
				        
				    }
					LBR_Assign_Get_data = append(LBR_Assign_Get_data, post)
				}
		
		
		type LBR_MemAss struct {
			Id 				interface{}	
			Lbr_assign		interface{} 
			Lbr_assign_name interface{}
			Orgmem 			interface{}
			Member_name		interface{}
			Leader			interface{}
			HRIS			interface{}
			Processaccess	interface{}
			Txdate			interface{}
			Device_id		interface{}
			Remarks			interface{}
			Enable			interface{}
			Disable 		interface{}

		}
		
		LBR_MemAss_raw ,err, _,_ := config.Ap_sql(`exec LBR_MemAss_Get 0, `+tconf["MemassID"]+`,`+str_OrgID+`` ,1)
	 	if err != nil {
	    	panic(err.Error())
		}
		var LBR_MemAss_data  LBR_MemAss
		for LBR_MemAss_raw.Next() {
			err = LBR_MemAss_raw.Scan(&LBR_MemAss_data.Id, &LBR_MemAss_data.Lbr_assign,&LBR_MemAss_data.Lbr_assign_name,&LBR_MemAss_data.Orgmem ,&LBR_MemAss_data.Member_name,&LBR_MemAss_data.Leader,&LBR_MemAss_data.HRIS,&LBR_MemAss_data.Processaccess,&LBR_MemAss_data.Txdate ,&LBR_MemAss_data.Device_id ,&LBR_MemAss_data.Remarks,&LBR_MemAss_data.Enable,&LBR_MemAss_data.Disable)
	            if err != nil {
	                    panic(err.Error())
	            }
	        //LBR_MemAss_data =   LBR_MemAss{}   
	    }

    
		type Data struct {
			LBR_Assign_Get []LBR_Assign_Get
			Htdata map[string]string
			LBR_MemAss LBR_MemAss
			
		}   
		htdata := Data{LBR_Assign_Get_data,tconf ,LBR_MemAss_data} 
     	
     	funcMap := template.FuncMap{
		    "date_formater": func(a time.Time) string {
		        return a.Format(" Jan _2 2006")
		    },
		} 
		tmpl := template.New("MemberAssignmentEdit.html").Funcs(funcMap)
		//var err error
		if tmpl, err = tmpl.ParseFiles("hris/members/MemberAssignmentEdit.html"); err != nil {
			fmt.Println(err)
		}

		err1 := tmpl.Execute(w,htdata)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
		
	}else{
		//do create assigment Hire.
		r.ParseForm()
		Org_id :=login.Get_session_org_id(r)
		str_OrgID :=strconv.Itoa(Org_id)

		session_username,_ := login.Get_account_info(r)
		sqlstr:= `exec LBR_Memass_Save 
						'Edit',
						'`+session_username + `' ,
						`+ r.Form["memassID"][0] + `  ,
						`+str_OrgID+`,
						` + r.Form["lbr_assign"][0] + ` ,
						` + r.Form["member"][0]+ `	 ,
						` + r.Form["processaccess"][0]+` ,
						'`+ r.Form["txdate"][0]+`' ,
						'`+ r.Form["device_id"][0]+`' ,
						'`+ r.Form["remarks"][0]+`' ,
						'`+ r.Form["enable"][0]+`' ,
						'`+ r.Form["disable"][0] + `' ,
						'`+ r.Form["leader"][0] + `' ,
						'`+ r.Form["emp_id"][0] + `' `
				
		rowret ,err, _,_ := config.Ap_sql(sqlstr,1)
		if err != nil {
			panic(err.Error())
		}
		var r RetMessage
		for rowret.Next() {
		    err = rowret.Scan(&r.RetNo,&r.RetMsg)
		    if err != nil {
		        panic(err.Error())
		    }
		 	r = RetMessage{r.RetNo,r.RetMsg}
		}

		profile := Profile{r.RetMsg,true}
		js, err := json.Marshal(profile)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
  		}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)


	}	

}	



func MemberAssigmentDeleteHAndler(w http.ResponseWriter, r *http.Request){
	login.Session_renew(w,r)
	if r.Method =="GET" {
		tconf := make(map[string]string)
		tconf["RecID"] = r.URL.Query().Get("ID")
	   	tconf["delete_url"] = "/Member/MemberAssignmentDelete"
	   	
	   	tmpl := template.New("modal_delete.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("hris/members/modal_delete.html"); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w,tconf)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}else{
		r.ParseForm()
		session_username,_ := login.Get_account_info(r)	
		sqlstr:=`	exec LBR_Memass_Save
		                    'Delete',
		                    '`+session_username+`' ,
		                    `+r.Form["rec_no"][0]+`
				`

  		rows2 ,err, _,_ := config.Ap_sql(sqlstr,1)
  		if err != nil {
	    	panic(err.Error())
	    }
		var errorMsg  Error_mess
	    for rows2.Next() {
            var r Error_mess
            err = rows2.Scan(&r.RetNo, &r.RetMsg)
            if err != nil {
                    panic(err.Error())
            }
			errorMsg = Error_mess{r.RetNo,r.RetMsg}
		      
		}
				
		
	    
		js, err := json.Marshal(errorMsg)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	fmt.Println("Data Successfully added")
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)

	}		



}
