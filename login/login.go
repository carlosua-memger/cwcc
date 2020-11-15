package login

import (
	"fmt"
	"net/http"
	"text/template"

	sessions "../github.com/gorilla/sessions"

	//email "../github.com/email"
	"database/sql"
	"encoding/json"
	"net/smtp"
	"strconv"
	"time"

	config "../config"
	datatables "../datatables"
	email "../github.com/email-master"
	//"time"
)

type access_list struct {
	Id          sql.NullInt64
	Name        string
	Parent      sql.NullInt64
	Url         sql.NullString
	Param       sql.NullString
	Tooltip     sql.NullString
	Rights      string
	Cte_level   int64
	Cte_id      sql.NullInt64
	Cte_name    sql.NullString
	Member_id   int
	Menurank    string
	Menudefault string
}
type SideBar struct {
	HostURL     string
	SideBarCont []access_list
}
type Content struct {
	HostURL string
}
type MainPageData struct {
	LogInUsername string
	HostURL       string
	SideBar       *SideBar
	Content       *Content
	Org_list      []Orgmem_get
	Tconf         map[string]string
}
type Profile struct {
	Message string // error status output or validated data
	Stats   string //pass or fail
	//Hobbies []string
}

var Uname string
var UserID int64
var key = []byte("mu5e5ecret_m3mG3r") // 32 bytes
var Store = sessions.NewCookieStore([]byte("something-very-secret-cwcc"))
var Store_uname = sessions.NewCookieStore([]byte("something-very-secret-cwcc"))

type Login struct {
	Username  string    // error status output or validated data
	Logtime   time.Time //pass or fail
	UserAgent string
	//Hobbies []string
}

var Arr_login []Login

//var Arr_login []Login
//var Arr_login_map make(map[string][]Login)
var Arr_login_map = make(map[string]Login)

/*type Orgmem_get struct{
	Org_id int
	Org_name string
}
*/
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tconf := make(map[string]string)
	tconf["login_status"] = r.URL.Query().Get("status")
	session, _ := Store.Get(r, "muse")
	session_username, _ := Store_uname.Get(r, "uname")
	a := session.Values["username"]
	fmt.Println("method:", r.Host) //get request method
	if a != nil {
		http.Redirect(w, r, "/mainframe", 302)
	}
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"

	//	w.Header().Set("Access-Control-Allow-Origin", "https://preview.c9user.io")
	w.Header().Set("Access-Control-Allow-Origin", "localhost:5656")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")

	if r.Method == "GET" {
		tmpl := template.New("landpage.html")
		var err error
		if tmpl, err = tmpl.ParseFiles("login/landpage.html"); err != nil {
			fmt.Println(err)
		}
		err = tmpl.Execute(w, tconf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	} else {
		r.ParseForm()
		fmt.Println(`user Agent here:`, r.UserAgent())
		fmt.Println(`__________________________________`)
		fmt.Println(`user Referer here:`, r.Referer())
		fmt.Println(`__________________________________`)

		sess_check := Server_ses_storage(template.HTMLEscapeString(r.Form.Get("username")))
		if sess_check != `` {
			fmt.Println(sess_check)
			http.Redirect(w, r, `/login?status=This user account is already logged in. Need additional user? <a id="hf-chat" col>click here.</a>`, 302)
		} else {
			client_info := `client ip: ` + r.RemoteAddr + `user agent:` + r.UserAgent()
			auth_type := r.Form["auth_type"][0]
			//auth_type =1  via google
			//auth_type = 0 normal login
			var param_authtype string = ``
			if auth_type == `1` {
				param_authtype = `,'Google'`
				r.Form["password"][0] = r.Form["username"][0]
			}
			fmt.Println("members_login '" + template.HTMLEscapeString(r.Form.Get("username")) + "','" + r.Form["password"][0] + "','" + client_info + "'" + param_authtype)
			//rows ,err, _,_ := config.Ap_sql("member_login '"+template.HTMLEscapeString(r.Form.Get("username")) +"','"+template.HTMLEscapeString(r.Form.Get("password")) +"'",1)
			rows, err, _, _ := config.Ap_sql("members_login '"+template.HTMLEscapeString(r.Form.Get("username"))+"','"+r.Form["password"][0]+"','"+client_info+"'"+param_authtype, 1)
			var rowConter int64
			var userid int
			//var Member_id
			for rows.Next() {
				var r access_list
				err = rows.Scan(&r.Id, &r.Name, &r.Parent, &r.Url, &r.Param, &r.Tooltip, &r.Rights, &r.Cte_level, &r.Cte_id, &r.Cte_name, &r.Member_id, &r.Menurank, &r.Menudefault)
				if err != nil {
					panic(err.Error())
				}
				rowConter = rowConter + 1
				//fmt.Println("Agi sya dare")
				userid = r.Member_id
				//Member_id = r.Member_id
				//fmt.Println("Agi sya dare 22")
			}
			fmt.Println("naihap na raws:", rowConter)

			if rowConter > 0 {
				Uname = template.HTMLEscapeString(r.Form.Get("username"))

				Store.Options = &sessions.Options{
					Path:     "/",
					MaxAge:   config.SessionMAxTime,
					HttpOnly: true,
				}
				// Set some session values.
				fmt.Println("session.Values['UserID'] = ", userid)
				session.Values["UserID"] = userid
				//session.Values["Member_id"] = Member_id
				session.Values["username"] = r.Form["username"][0]
				session.Values["password"] = r.Form["password"][0]
				session.Values["isLogin"] = "1"
				// session.Values["logdatetime"] = time.

				// Save it.
				session.Save(r, w)

				Store_uname.Options = &sessions.Options{
					Path:     "/",
					MaxAge:   0,
					HttpOnly: true,
				}
				session_username.Values["username"] = r.Form["username"][0]
				session_username.Save(r, w)
				// session server saving
				t := time.Now()
				//Arr_login = append(Arr_login,Login{Uname,t,} )
				// end session server saving
				Arr_login_map[Uname] = Login{Uname, t, r.UserAgent()}
				http.Redirect(w, r, "/mainframe", 302)
			} else {
				//renderTemplate(w,&HtmlPage{htdata:INDEX_HTML})
				//http.Redirect(w, r, "/login_attempt", 302)

				tmpl := template.New("landpage.html")
				tconf := make(map[string]string)
				tconf["login_status"] = "Invalid username or password."
				var err error
				if tmpl, err = tmpl.ParseFiles("login/landpage.html"); err != nil {
					fmt.Println(err)
				}
				err = tmpl.Execute(w, tconf)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)

				}
			}
		}

	}
}

func LoginAttemptHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("landpage.html")
	tconf := make(map[string]string)
	tconf["login_status"] = r.URL.Query().Get("status")
	var err error
	if tmpl, err = tmpl.ParseFiles("login/landpage.html"); err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, tconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logout here")
	ses_username, _ := Get_account_info(r)
	fmt.Println(`Member_Logout '` + ses_username + `'`)
	_, err, _, _ := config.Ap_sql(`Member_Logout '`+ses_username+`'`, 3)
	if err != nil {
		panic(err.Error())
	}
	session, _ := Store.Get(r, "muse")
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   1,
		HttpOnly: true,
	}
	session.Values["username"] = nil
	session.Values["password"] = nil
	session.Values["isLogin"] = nil
	session.Save(r, w)
	Delete_Session := Server_ses_storage_delete(ses_username)
	fmt.Println(Delete_Session)
	for key, seval := range Arr_login_map {
		fmt.Println(`key sa login`, key)
		fmt.Println(`val sa login usrname`, seval.Username)
		fmt.Println(`val sa login logtime`, seval.Logtime)
	}

	http.Redirect(w, r, "/", 302)
}

func ReloginHandler(w http.ResponseWriter, r *http.Request) {
	sess_check := Server_ses_storage(template.HTMLEscapeString(r.Form.Get("username")))
	if sess_check != `` {
		fmt.Println(sess_check)
		profile := Profile{`This user account is already logged in. Need additional user? <a id="hf-chat"> click here. </a>`, "0"}
		js, err := json.Marshal(profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		session, _ := Store.Get(r, "muse")
		fmt.Println("session here  --- ::   ", session.Values["username"])
		r.ParseForm()
		client_info := `client ip: ` + r.RemoteAddr + `user agent:` + r.UserAgent()
		//rows ,err, _,_ := config.Ap_sql("member_login '"+template.HTMLEscapeString(r.Form.Get("username")) +"','"+template.HTMLEscapeString(r.Form.Get("password")) +"'",1)
		rows, err, _, _ := config.Ap_sql("members_login '"+template.HTMLEscapeString(r.Form.Get("username"))+"','"+template.HTMLEscapeString(r.Form.Get("password"))+"','"+client_info+"'", 1)
		//rows ,err, _,_ := config.Ap_sql("member_login '"+template.HTMLEscapeString(r.Form.Get("username")) +"','"+template.HTMLEscapeString(r.Form.Get("password")) +"'",1)
		var rowConter int64
		var userid int
		for rows.Next() {
			var r access_list
			err = rows.Scan(&r.Id, &r.Name, &r.Parent, &r.Url, &r.Param, &r.Tooltip, &r.Rights, &r.Cte_level, &r.Cte_id, &r.Cte_name, &r.Member_id, &r.Menurank, &r.Menudefault)
			if err != nil {
				panic(err.Error())
			}
			userid = r.Member_id
			rowConter = rowConter + 1
		}
		if rowConter > 0 {

			session_uname, _ := Store_uname.Get(r, "uname")
			if session_uname.Values["username"] != r.Form["username"][0] {

				profile := Profile{"this user is new, refresh browser is required.", "0"}
				js, err := json.Marshal(profile)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)

			} else {

				Store.Options = &sessions.Options{
					Path:     "/",
					MaxAge:   config.SessionMAxTime,
					HttpOnly: true,
				}

				// Set some session values.
				session.Values["UserID"] = userid
				session.Values["username"] = r.Form["username"][0]
				session.Values["password"] = r.Form["password"][0]

				session.Values["isLogin"] = "1"
				// Save it.
				t := time.Now()
				Arr_login_map[Uname] = Login{r.Form["username"][0], t, r.UserAgent()}

				session.Save(r, w)

				profile := Profile{"", "1"}
				js, err := json.Marshal(profile)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)

			}
		} else {
			profile := Profile{"invalid username or password", "0"}
			js, err := json.Marshal(profile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)

		}
	}
}

func FrontIndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("index_ap.html")
	var err error
	if tmpl, err = tmpl.ParseFiles("tpl/index/index_ap.html"); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w, nil)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func SessionExpireHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Host)
	isprocess := r.URL.Query().Get("process")
	if isprocess == "true" {
		Session_renew_max(w, r)
	}

	var isExpire int
	isExpire = 0
	session, _ := Store.Get(r, "muse")
	// Set some session values.
	a := session.Values["isLogin"]
	b, _ := a.(string)

	//Uname :=session.Values["username"].(string)
	fmt.Println("sessin data1" + b)

	if a == nil {
		isExpire = 1
		fmt.Println("session expire")
	}
	var err error
	tmpl := template.New("session_expire.html")
	var err0 error
	if tmpl, err0 = tmpl.ParseFiles("tpl/session_expire.html"); err != nil {
		fmt.Println(err0)
	}
	err1 := tmpl.Execute(w, struct{ IsEx int }{isExpire})
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

type Orgmem_get struct {
	Org_id   int
	Org_name string
}

func MainFrameHandler(w http.ResponseWriter, r *http.Request) {

	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   config.SessionMAxTime,
		HttpOnly: true,
	}
	session, _ := Store.Get(r, "muse")
	// Set some session values.
	a := session.Values["username"]
	sesuserid, _ := session.Values["UserID"].(int)
	b, _ := a.(string)
	Uname = b
	fmt.Println("user ID : ", UserID)
	//pas :=session.Values["password"]
	//ses_pas ,_ := pas.(string)
	//fmt.Println("sessin data1"+ b)
	if a == nil {
		http.Redirect(w, r, "/login", 302)
	} else {

		var ses_Org_id int
		loopCount := 0

		str_userID := strconv.Itoa(sesuserid)

		orgmem_get_row, err, _, _ := config.Ap_sql(`sysaccess_get 1 ,`+str_userID+`,0`, 1)
		var r_org Orgmem_get
		orgmem_get_data := []Orgmem_get{}
		for orgmem_get_row.Next() {
			err = orgmem_get_row.Scan(&r_org.Org_id, &r_org.Org_name)
			if err != nil {
				panic(err.Error())
			}
			if loopCount == 0 {
				ses_Org_id = r_org.Org_id
			}
			loopCount = loopCount + 1
			orgmem_get_data = append(orgmem_get_data, Orgmem_get{r_org.Org_id, r_org.Org_name})
		}

		// get slide bar data
		//rows ,err, _,_ := ap_sql("member_login '"+decrypt(key,cookie_username.Value)+"','',1 ",1)
		//rows ,err, _,_ := config.Ap_sql("member_login '"+b+"','"+ses_pas+"' ",1)

		str_OrgID := strconv.Itoa(ses_Org_id)
		//session.Values["org_id"] = str_OrgID
		tconf := make(map[string]string)
		tconf["org_id"] = str_OrgID
		tconf["userID"] = str_userID
		url_org_id := r.URL.Query().Get("org_id")
		fmt.Println(`mao ni sa org id url sidde :`, url_org_id)
		int_org_id, conv_errr := strconv.Atoi(url_org_id)
		fmt.Println("convi error:")
		fmt.Println(conv_errr)
		if conv_errr != nil {
			int_org_id = 0
		}
		fmt.Println("url org id here")
		if int_org_id > 0 {
			fmt.Println(url_org_id)
			session.Values["Org_id"] = int_org_id
			tconf["org_id"] = url_org_id
		} else {
			session.Values["Org_id"] = ses_Org_id

		}
		fmt.Println(`mao ni sa org id na sa session ma save :`, session.Values["Org_id"])
		session.Save(r, w)
		//sysaccess_get 2, @MemberID, @Org

		fmt.Println(`exec orgmemrol_get 0, ` + tconf["org_id"] + `, 0, ` + str_userID + `,1`)
		rows, err, _, _ := config.Ap_sql(`exec sysaccess_get 2, `+str_userID+`,`+tconf["org_id"], 1)
		//rows ,err, _,_ := config.Ap_sql(`exec orgmemrol_get 0, `+tconf["org_id"]+`, 0, `+str_userID+`,1`,1)
		if err != nil {
			panic(err.Error())
		}
		got := []access_list{}
		for rows.Next() {
			var r access_list
			//err = rows.Scan(&r.Account_id,&r.Username,&r.Id, &r.Name, &r.Parent , &r.Url, &r.Param, &r.Tooltip, &r.Leveling)
			err = rows.Scan(&r.Id, &r.Name, &r.Parent, &r.Url, &r.Param, &r.Tooltip, &r.Rights, &r.Cte_level, &r.Cte_id, &r.Cte_name, &r.Member_id, &r.Menurank, &r.Menudefault)
			if err != nil {
				panic(err.Error())
			}
			if !r.Parent.Valid {
				r.Parent = sql.NullInt64{Int64: 0, Valid: true}
			}
			if !r.Tooltip.Valid {
				r.Tooltip = sql.NullString{String: "Null", Valid: true}
			}

			post := access_list{
				r.Id,
				r.Name,
				r.Parent,
				r.Url,
				r.Param,
				r.Tooltip,
				config.Encrypt(config.Encrypt_key, r.Rights),
				r.Cte_level,
				r.Cte_id,
				r.Cte_name,
				r.Member_id,
				r.Menurank,
				r.Menudefault,
			}
			got = append(got, post)

		}

		//end slidebar data
		//costum template on go
		funcMap := template.FuncMap{
			"equals": func(a sql.NullInt64, b sql.NullInt64) bool {
				a1 := a.Int64
				b1 := b.Int64
				//fmt.Println(a1,"----",b1)

				return a1 == b1
			},
			"str_to_int": func(D1 string) int {
				int_val, _ := strconv.Atoi(D1)
				return int_val
			},
		}

		pagedata := &MainPageData{LogInUsername: Uname, HostURL: r.Host, SideBar: &SideBar{HostURL: r.Host, SideBarCont: got}, Org_list: orgmem_get_data, Tconf: tconf}
		tmpl := template.New("MainFrame.html").Funcs(funcMap)
		var err0 error
		if tmpl, err0 = tmpl.ParseFiles("tpl/MainFrame.html"); err != nil {
			fmt.Println(err0)
		}
		if tmpl, err0 = tmpl.ParseFiles("tpl/sidebar.html"); err != nil {
			fmt.Println(err0)
		}
		/*if tmpl, err0 = tmpl.ParseFiles("tpl/DashBoardContent.html"); err != nil {
			fmt.Println(err0)
		}*/
		err1 := tmpl.Execute(w, pagedata)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}
	}
}

type Error_mess struct {
	RetNo  int
	RetMsg string
}

func SignupverifyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("decrypt token")
	token := r.URL.Query().Get("token")
	fmt.Println(token)

	hkey := config.Decrypt(config.Encrypt_key, token)
	fmt.Println("decrypt token")
	fmt.Println(hkey)
	arr_ret := datatables.Data_row(`signup_get 1, '` + hkey + `'`)
	fmt.Println(len(arr_ret))
	fmt.Println(arr_ret)
	tconf := make(map[string]string)
	tconf["org_name"] = ``
	if len(arr_ret) > 0 {
		tconf["org_name"] = arr_ret[0]
	}

	tconf["email"] = hkey

	tmpl := template.New("signup_notify.html")
	var err error
	if tmpl, err = tmpl.ParseFiles("login/signup_notify.html"); err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, tconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	tconf := make(map[string]string)

	ses_username, int_memgerid := Get_account_info(r)
	tconf["username"] = ses_username
	tconf["memberid"] = strconv.Itoa(int_memgerid)
	tconf["flag"] = ``

	tmpl := template.New("change_password.html")
	var err error
	if tmpl, err = tmpl.ParseFiles("login/change_password.html"); err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, tconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func DoChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	tconf := make(map[string]string)

	r.ParseForm()
	username := r.Form["username"][0]
	memberid := r.Form["memberid"][0]
	oldpassword := r.Form["oldpassword"][0]
	password := r.Form["password"][0]

	fmt.Println(`Member_ChangePassword  1, '` + username + `', ` + memberid + `, '` + oldpassword + `', '` + password + `'`)
	ret_changepasword := datatables.Data_row(`Member_ChangePassword  1, '` + username + `', ` + memberid + `, '` + oldpassword + `', '` + password + `'`)
	fmt.Println(ret_changepasword[1])
	tconf["username"] = username
	tconf["memberid"] = memberid
	tconf["flag"] = ret_changepasword[1]

	tmpl := template.New("change_password.html")
	var err error
	if tmpl, err = tmpl.ParseFiles("login/change_password.html"); err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, tconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signning up here....")
	r.ParseForm()

	orgName := r.Form["orgName"][0]
	post_email := r.Form["email"][0]
	password := r.Form["password"][0]
	fmt.Println(post_email)

	arr_ret_check := datatables.Data_row(`Signup_viaEmail_Check '','` + post_email + `'`)
	fmt.Println(`check email`, arr_ret_check)

	type Tag struct {
		Status   bool
		HtmlData string `json:"HtmlData"`
	}
	var ret_data Tag

	if len(arr_ret_check) == 0 {

		arr_ret := datatables.Data_row(`Signup_Save 'Add',0,'','','','` + post_email + `','` + orgName + `','` + password + `' `)
		fmt.Println(`signUP reply`)
		fmt.Println(arr_ret)

		hkey := config.Encrypt(config.Encrypt_key, post_email)

		//tconf := make(map[string]string)
		fmt.Println(config.Email_id)
		fmt.Println(config.Email_password)

		e := email.NewEmail()
		e.From = config.Email_id
		e.To = []string{post_email}
		e.Bcc = []string{}
		e.Cc = []string{config.Email_cc}

		e.Subject = `Complete your sign up and get more done with Memger`

		e.HTML = []byte(`Thanks for getting started with Memger! We need a little more information to complete your registration, including confirmation of your email address. Click below to confirm your email address: http://memger.com/signup-verify-email?token=` + hkey + `
			<br>If you have problems, please paste the above URL into your web browser.
			<br>
			<br>
			P.S.: This is an auto-generated message. <br>Please do not reply to this email.
			<br>
			<br>
			Thanks,
			<br>Memger Support`)

		retErr := e.Send("smtp.zoho.com:587", smtp.PlainAuth("", e.From, config.Email_password, "smtp.zoho.com"))
		fmt.Println(`Email sending Error : `, retErr)

		/*
				email.SendEmail(
			        "smtp.zoho.com",
			        587,
			        config.Email_id,
			        config.Email_password,
			        []string{post_email},
			        "Complete your sign up and get more done with Memger",
			        `

						Thanks for getting started with Memger! We need a little more information to complete your registration, including confirmation of your email address. Click below to confirm your email address: http://memger.com//signup-verify-email?token=`+hkey+`
						<br>If you have problems, please paste the above URL into your web browser.
						<br>
						<br>
						P.S.: This is an auto-generated message. <br>Please do not reply to this email.
						<br>
						<br>
						Thanks,
						<br>Memger Support

			        `,
			        tconf)*/

		var ht string = `<div class="note note-success">`
		ht += `<h4 class="block">Signup Sucess</h4>`
		ht += ` <p>`
		ht += `Email verification sent to ` + post_email + ` , to continue please check your email and click the given link. Thank you.`
		ht += `</p>`
		ht += `</div>`
		ret_data = Tag{true, ht}
	} else {
		ret_data = Tag{false, `This email is already Exist`}
	}

	js, err := json.Marshal(ret_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func SignupverifyValidateHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	org_name := ``
	if len(r.Form["orgName"]) > 0 {
		org_name = r.Form["orgName"][0]
	}
	email := ``
	if len(r.Form["email"]) > 0 {
		email = r.Form["email"][0]
	}
	fmt.Println(`org typed`, org_name)
	sqlstr := `Signup_viaEmail_Check '` + org_name + `' ,'` + email + `'`
	fmt.Println(sqlstr)
	arr_ret := datatables.Data_row(sqlstr)
	fmt.Println(arr_ret)

	type Tag struct {
		Valid bool `json:"valid"`
	}
	var ret_data Tag
	if len(arr_ret) > 0 {
		ret_data = Tag{false}
	} else {
		ret_data = Tag{true}
	}
	js, err := json.Marshal(ret_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func SignupViaEmailHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	org_name := r.Form["orgName"][0]
	email := r.Form["email"][0]
	lastname := r.Form["lastname"][0]
	middlename := r.Form["middlename"][0]
	firstname := r.Form["firstname"][0]
	password := r.Form["password"][0]

	//sqlstr:=`Signup_Save 'Add',0,'`+r.Form["lastname"][0]+`','`+r.Form["middlename"][0]+`','`+r.Form["firstname"][0]+`','`+r.Form["email"][0]+`','`+r.Form["contactno"][0]+`','Assignment Name on latest timesheet?','`+r.Form["a1"][0]+`','Biometrics ID on latest Timesheet?','`+r.Form["a2"][0]+`','',''`
	sqlstr := `signup_viaEmail '` + org_name + `', '` + email + `', '` + lastname + `', '` + middlename + `', '` + firstname + `','` + password + `'`
	fmt.Println(sqlstr)
	arr_ret := datatables.Data_row(sqlstr)
	fmt.Println(arr_ret)
	/*js, err := json.Marshal(arr_ret)
		if err != nil {
	    	http.Error(w, err.Error(), http.StatusInternalServerError)
	    	return
			}
	  	w.Header().Set("Content-Type", "application/json")
	  	w.Write(js)	*/
	if len(arr_ret) > 0 {
		if arr_ret[1] == `Success` {
			//http.Redirect(w, r, "/login", 302)
			session, _ := Store.Get(r, "muse")
			session_username, _ := Store_uname.Get(r, "uname")
			Uname = email

			Store.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   config.SessionMAxTime,
				HttpOnly: true,
			}
			// string to int conversion
			v := arr_ret[0]
			member_id_int, _ := strconv.Atoi(v)

			//end of conversion

			// Set some session values.
			fmt.Println("session.Values['UserID'] from signup= ", arr_ret[0])
			session.Values["UserID"] = member_id_int
			//session.Values["Member_id"] = Member_id
			session.Values["username"] = email
			session.Values["password"] = password
			session.Values["isLogin"] = "1"
			// session.Values["logdatetime"] = time.

			// Save it.
			session.Save(r, w)

			Store_uname.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   0,
				HttpOnly: true,
			}
			session_username.Values["username"] = email
			session_username.Save(r, w)
			// session server saving
			//t := time.Now()
			//Arr_login = append(Arr_login,Login{Uname,t,} )
			// end session server saving
			//Arr_login_map[Uname] = Login{Uname,t,r.UserAgent(),}
			http.Redirect(w, r, "/mainframe", 302)

		} else {
			fmt.Fprint(w, `Opss the the system dedect multiple posting of data. or may this email-address already exist.`)
		}

	}

}

func Get_account_info(r *http.Request) (string, int) {

	tconf := make(map[string]string)
	tconf["login_status"] = r.URL.Query().Get("status")
	session, _ := Store.Get(r, "muse")
	a := session.Values["username"].(string)
	b := session.Values["UserID"].(int)
	return a, b
}

func Get_account_info_id(r *http.Request) (int, string) {

	tconf := make(map[string]string)
	tconf["login_status"] = r.URL.Query().Get("status")
	session, _ := Store.Get(r, "muse")
	a := session.Values["UserID"].(int)
	b := session.Values["username"].(string)
	return a, b
}

func Get_session_org_id(r *http.Request) int {

	tconf := make(map[string]string)
	tconf["login_status"] = r.URL.Query().Get("status")
	session, _ := Store.Get(r, "muse")
	a := session.Values["Org_id"].(int)
	return a
}

func Session_renew(w http.ResponseWriter, r *http.Request) {
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   config.SessionMAxTime,
		HttpOnly: true,
	}
	session, _ := Store.Get(r, "muse")
	// Set some session values.
	a := session.Values["username"]
	sesuserid := session.Values["UserID"].(int)
	fmt.Println("session.Values[UserID].(int64) --- ", sesuserid)
	b, _ := a.(string)
	Uname = b
	//UserID = sesuserid
	fmt.Println("user ID : ", UserID)
	//pas :=session.Values["password"]
	//ses_pas ,_ := pas.(string)
	fmt.Println("sessin data1" + b)
	if a == nil {
		fmt.Fprint(w, "session error")
		//http.Redirect(w, r, "/login", 302)
	} else {
		t := time.Now()
		//Arr_login = append(Arr_login,Login{Uname,t,} )
		// end session server saving
		Arr_login_map[Uname] = Login{Uname, t, r.UserAgent()}
		session.Save(r, w)

	}
}

func Session_renew_max(w http.ResponseWriter, r *http.Request) {
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	}
	session, _ := Store.Get(r, "muse")
	// Set some session values.
	a := session.Values["username"]
	sesuserid := session.Values["UserID"].(int)
	fmt.Println("session.Values[UserID].(int64) --- ", sesuserid)
	b, _ := a.(string)
	Uname = b
	//UserID = sesuserid
	fmt.Println("user ID : ", UserID)
	//pas :=session.Values["password"]
	//ses_pas ,_ := pas.(string)
	fmt.Println("sessin data1" + b)
	if a == nil {
		//fmt.Fprint(w,"session error")
		//http.Redirect(w, r, "/login", 302)
	} else {
		session.Save(r, w)
	}
}
func Server_ses_storage(uname string) string {
	var ret_str = ``
	/*for key , seval := range Arr_login {
		fmt.Println(`key sa login`,key)
		fmt.Println(`val sa login usrname`,seval.Username)
		fmt.Println(`val sa login logtime`,seval.Logtime)
		if(seval.Username==uname){
			ret_str = `Account already Login`
		}
	}*/
	var diff float64
	for key, seval := range Arr_login_map {
		fmt.Println(`key sa login`, key)
		fmt.Println(`val sa login usrname`, seval.Username)
		fmt.Println(`val sa login logtime`, seval.Logtime)
		t1 := time.Now()
		diff = t1.Sub(seval.Logtime).Minutes()
		fmt.Println("time diff:", diff)
	}
	if (Arr_login_map[uname].Username == uname) && (diff <= 1) {
		ret_str = `This user account is already logged in. Need additional user? <a id="hf-chat">click here.</a>`
	}

	fmt.Println(`Arr_login_map[uname] = `, Arr_login_map[uname])
	return ret_str
}

func Server_ses_storage_delete(uname string) string {
	delete(Arr_login_map, uname)
	var ret_str = `deleted`
	return ret_str
}
