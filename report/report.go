package report

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	config "../config"
	//"strconv"
	//"encoding/json"
	datatables "../datatables"
	login "../login"

	//"strings"
	helper "../helper"
)

type Profile struct {
	Message string // error status output or validated data
	Status  bool   //pass or fail
	ID      string
}

func Show_savings_ledger(w http.ResponseWriter, r *http.Request) {
	member := r.URL.Query().Get("member")
	trantype := r.URL.Query().Get("trantype")
	member_name := r.URL.Query().Get("member_name")
	username, _ := login.Get_account_info(r)

	fmt.Println(`LedgerS_Display 1,` + member + `,0`)
	arr_data := datatables.DataList(`LedgerS_Display 1,` + member + `,` + trantype + ``)

	type Data struct {
		Tconf      map[string]string
		Arr_data   [][]string
		MemberName string
	}
	var str_page string
	var str_page_dir string
	str_page = `member_ledgerS.html`
	str_page_dir = `report/member_ledgerS.html`
	tconf := make(map[string]string)
	tconf["username"] = username
	tmpl := template.New(str_page).Funcs(helper.FuncMap)
	var err error
	if tmpl, err = tmpl.ParseFiles(str_page_dir); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w, Data{tconf, arr_data, member_name})
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func Show_loans_ledger(w http.ResponseWriter, r *http.Request) {
	member := r.URL.Query().Get("member")
	trantype := r.URL.Query().Get("trantype")
	member_name := r.URL.Query().Get("member_name")
	username, _ := login.Get_account_info(r)

	fmt.Println(`LedgerL_Display 1,` + member + `,0`)
	arr_data := datatables.DataList(`LedgerL_Display 1,` + member + `,` + trantype + ``)

	type Data struct {
		Tconf      map[string]string
		Arr_data   [][]string
		MemberName string
	}
	var str_page string
	var str_page_dir string
	str_page = `member_ledgerL.html`
	str_page_dir = `report/member_ledgerL.html`
	tconf := make(map[string]string)
	tconf["username"] = username
	tmpl := template.New(str_page).Funcs(helper.FuncMap)
	var err error
	if tmpl, err = tmpl.ParseFiles(str_page_dir); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w, Data{tconf, arr_data, member_name})
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func Show_loans_ledger_withbalance(w http.ResponseWriter, r *http.Request) {
	member := r.URL.Query().Get("member")
	trantype := r.URL.Query().Get("trantype")
	member_name := r.URL.Query().Get("member_name")
	username, _ := login.Get_account_info(r)

	fmt.Println(`LedgerL_Display 1,` + member + `,0`)
	arr_data := datatables.DataList(`LedgerL_Display 4,` + member + `,` + trantype + ``)

	type Data struct {
		Tconf      map[string]string
		Arr_data   [][]string
		MemberName string
	}
	var str_page string
	var str_page_dir string
	str_page = `member_ledgerLnobalance.html`
	str_page_dir = `report/member_ledgerLnobalance.html`
	tconf := make(map[string]string)
	tconf["username"] = username
	tmpl := template.New(str_page).Funcs(helper.FuncMap)
	var err error
	if tmpl, err = tmpl.ParseFiles(str_page_dir); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w, Data{tconf, arr_data, member_name})
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func Show_savings_ledger_modal(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)

	//username, _ := login.Get_account_info(r)
	//Org_id :=login.Get_session_org_id(r)
	//str_OrgID :=strconv.Itoa(Org_id)
	//username := "static"
	//ses_uID := 1
	//str_ses_uID :=strconv.Itoa(ses_uID)
	member := r.URL.Query().Get("member")
	trantype := r.URL.Query().Get("trantype")
	username, _ := login.Get_account_info(r)

	//isbarcode :=r.URL.Query().Get("isbarcode")

	//assignmentID :=r.URL.Query().Get("location")

	fmt.Println("member ", member)
	if r.Method == "GET" {
		//.Println(`this is hdr id=`,hdr_id)

		tconf := make(map[string]string)
		//tconf["parentID"] =r.URL.Query().Get("parentID")

		type Data struct {
			//Rights	string
			Tconf     map[string]string
			SupRepHdr [][]string
		}
		var err error
		var html_file string = ``
		var html_filename string = ``

		tconf["username"] = username

		arrdata := datatables.DataList(`ledgers_display 1,` + member + `,` + trantype)
		fmt.Println(arrdata)
		html_filename = "member_ledgerS_modal.html"
		html_file = "report/member_ledgerS_modal.html"

		tmpl := template.New(html_filename)
		if tmpl, err = tmpl.ParseFiles(html_file); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{tconf, arrdata})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}

	}
}

func Show_loans_ledger_modal(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	//username, _ := login.Get_account_info(r)
	//Org_id :=login.Get_session_org_id(r)
	//str_OrgID :=strconv.Itoa(Org_id)
	//username := "static"
	//ses_uID := 1
	//str_ses_uID :=strconv.Itoa(ses_uID)
	member := r.URL.Query().Get("member")
	lrid := r.URL.Query().Get("lrid")

	//isbarcode :=r.URL.Query().Get("isbarcode")

	//assignmentID :=r.URL.Query().Get("location")

	fmt.Println("member ", member)
	if r.Method == "GET" {
		//.Println(`this is hdr id=`,hdr_id)

		tconf := make(map[string]string)
		//tconf["parentID"] =r.URL.Query().Get("parentID")

		type Data struct {
			//Rights	string
			Tconf     map[string]string
			SupRepHdr [][]string
		}
		var err error
		var html_file string = ``
		var html_filename string = ``

		tconf["lrid"] = lrid

		arrdata := datatables.DataList(`ledgerl_display 2,` + member + `,` + lrid)
		fmt.Println(arrdata)
		html_filename = "member_ledgerL_modal.html"
		html_file = "report/member_ledgerL_modal.html"

		tmpl := template.New(html_filename)
		if tmpl, err = tmpl.ParseFiles(html_file); err != nil {
			fmt.Println(err)
		}
		err1 := tmpl.Execute(w, &Data{tconf, arrdata})
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
		}

	}
}

func Special_rpt(w http.ResponseWriter, r *http.Request) {
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	tconf := make(map[string]string)
	tconf["sql_code"] = r.URL.Query().Get("sql_code")
	tconf["rpt_url"] = r.URL.Query().Get("rpt_url")
	tconf["org_id"] = str_OrgID
	tconf["trnID"] = r.URL.Query().Get("trnID")
	tconf["trantype"] = r.URL.Query().Get("trantype")

	username, _ := login.Get_account_info(r)
	user_id, _ := login.Get_account_info_id(r)

	tconf["username"] = username
	tconf["user_id"] = strconv.Itoa(user_id)
	fmt.Println("on printing")
	fmt.Println(tconf["sql_code"])
	fmt.Println(tconf["rpt_url"])
	//fmt.Println(path.Base(tconf["rpt_url"]))

	arr_reportname := strings.Split(string(tconf["rpt_url"]), `/`)
	fmt.Println(arr_reportname)
	fmt.Println(len(arr_reportname))
	fmt.Println(arr_reportname[len(arr_reportname)-1])

	data, hdata := Datatable(tconf["sql_code"], nil)
	//var arrname []string
	type Htdata struct {
		Data    [][]string
		Hdata   []string
		Arrname []string
		Conf    map[string]string
	}
	var_htdata := Htdata{data, hdata, nil, tconf}
	///rawreports_dailys_lograw.html
	var default_tpl_name = arr_reportname[len(arr_reportname)-1]
	var default_tpl_dir = string(tconf["rpt_url"])

	tmpl := template.New(default_tpl_name).Funcs(helper.FuncMap)
	var err error
	if tmpl, err = tmpl.ParseFiles(default_tpl_dir); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w, var_htdata)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}

}

func Datatable(sql string, conf map[string]string) ([][]string, []string) {

	var m [][]string
	rows, err, _, _ := config.Ap_sql(sql, 1)
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

		for i := 0; i < len(cols); i++ { // fetching data tables to container map array
			a = append(a, PrintValue(vals[i].(*interface{})))
		}
		m = append(m, a)
		//mars at work here
	}

	if rows.Err() != nil {
		fmt.Println(rows.Err)
	}
	return m, cols

}

func PrintValue(pval *interface{}) string {
	var a string
	switch v := (*pval).(type) {
	case nil:
		//		fmt.Print("NULL")
		a = fmt.Sprint("")
	case bool:
		if v {
			//			fmt.Print("1")
			a = fmt.Sprint("1")
		} else {
			//			fmt.Print("0")
			a = fmt.Sprint("0")
		}
	case []byte:
		//		fmt.Print(string(v))
		a = fmt.Sprint(string(v))
	case time.Time:
		//		fmt.Print(v.Format("2006-01-02 15:04:05.999"))
		a = fmt.Sprint(v.Format("2006-01-02 15:04:05.999"))

	default:
		//		fmt.Print(v)
		a = fmt.Sprint(v)
	}
	return a
}

func Rawreports_doprint(w http.ResponseWriter, r *http.Request) {
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	username, _ := login.Get_account_info(r)
	user_id, _ := login.Get_account_info_id(r)

	tconf := make(map[string]string)
	tconf["username"] = username
	tconf["user_id"] = strconv.Itoa(user_id)
	tconf["org_id"] = str_OrgID
	tconf["name"] = r.URL.Query().Get("reportname")
	tconf["report_name"] = r.URL.Query().Get("report_name")
	tconf["template"] = r.URL.Query().Get("template")

	tconf["datefrom"] = r.URL.Query().Get("datefrom")
	tconf["dateto"] = r.URL.Query().Get("dateto")

	tconf["export_type"] = r.URL.Query().Get("export_type")

	//tconf["trn_id"] = strings.Trim(r.URL.Query().Get("trn_id"), ` `)
	tconf["member"] = r.URL.Query().Get("member")
	tconf["member_name"] = r.URL.Query().Get("member_name")
	tconf["group"] = r.URL.Query().Get("group")
	tconf["group_name"] = r.URL.Query().Get("group_name")
	tconf["partner"] = r.URL.Query().Get("partner")
	tconf["partner_name"] = r.URL.Query().Get("partner_name")
	tconf["paymode"] = r.URL.Query().Get("paymode")
	tconf["paymode_name"] = r.URL.Query().Get("paymode_name")
	tconf["status_name"] = r.URL.Query().Get("status_name")
	tconf["status"] = r.URL.Query().Get("status")
	tconf["coor"] = r.URL.Query().Get("coor")
	tconf["coor_name"] = r.URL.Query().Get("coor_name")

	tconf["trantype"] = r.URL.Query().Get("trantype")
	tconf["trantype_name"] = r.URL.Query().Get("trantype_name")

	tconf["bank"] = r.URL.Query().Get("bank")
	tconf["bank_name"] = r.URL.Query().Get("bank_name")

	tconf["loantype"] = r.URL.Query().Get("loantype")
	tconf["loantype_name"] = r.URL.Query().Get("loantype_name")

	tconf["strt_date_display"] = r.URL.Query().Get("datefrom")
	tconf["end_date_display"] = r.URL.Query().Get("dateto")
	tconf["filter_name"] = r.URL.Query().Get("filter_name")
	tconf["printtype"] = r.URL.Query().Get("printtype")
	tconf["printtype_name"] = r.URL.Query().Get("printtype_name")
	tconf["laetype"] = r.URL.Query().Get("laetype")
	tconf["laetype_name"] = r.URL.Query().Get("laetype_name")
	tconf["aetype"] = r.URL.Query().Get("aetype")
	tconf["aetype_name"] = r.URL.Query().Get("aetype_name")
	tconf["asof"] = r.URL.Query().Get("asof")
	tconf["lptype"] = r.URL.Query().Get("lptype")
	tconf["lptype_name"] = r.URL.Query().Get("lptype_name")
	tconf["othdtltype"] = r.URL.Query().Get("othdtltype")
	tconf["othdtltype_name"] = r.URL.Query().Get("othdtltype_name")
	tconf["partregioncoop"] = r.URL.Query().Get("partregioncoop")
	tconf["partregioncoop_name"] = r.URL.Query().Get("partregioncoop_name")

	tconf["show_datefrom"] = "0"
	tconf["show_dateto"] = "0"
	tconf["show_group"] = "0"
	tconf["show_partner"] = "0"
	tconf["show_loantype"] = "0"
	tconf["show_paymode"] = "0"
	tconf["show_bank"] = "0"
	tconf["show_filter"] = "0"
	tconf["show_status"] = "0"
	tconf["show_coor"] = "0"
	tconf["show_printtype"] = "0"
	tconf["show_trantype"] = "0"
	tconf["show_member"] = "0"
	tconf["show_laetype"] = "0"
	tconf["show_aetype"] = "0"
	tconf["show_asof"] = "0"
	tconf["show_lptype"] = "0"
	tconf["show_othdtltype"] = "0"
	tconf["show_partregioncoop"] = "0"

	/*
		tconf["show_id"] = "1"
		tconf["show_location"] = "1"
		tconf["show_datefrom"] = "1"
		tconf["show_dateto"] = "1"
		tconf["show_supplier"] = "1"
		tconf["show_soldby"] = "1"
		tconf["show_customer"] = "1"
		tconf["show_transferto"] = "1"
		tconf["show_transferfr"] = "1"
	*/
	/* if len(tconf["trn_id"]) == 0 {
		tconf["trn_id"] = `0`
	}
	if len(tconf["soldby"]) == 0 {
		tconf["soldby"] = `0`
		tconf["soldby_name"] = `All`
	} */
	//tconf["orgid"]="2"

	tconf["org_id"] = str_OrgID
	fmt.Println("sp run here.............")
	fmt.Println(tconf["name"])

	if tconf["member"] == "" {
		tconf["member"] = "0"
	}

	//strt := tconf["reportname"]

	//lbr_rep_logpair 3, 2, 0 , 46, '21 Oct 2015', '3 Nov 2015'
	arr_reportname := strings.Split(string(tconf["name"]), "@@")
	fsqlstr := arr_reportname[0]
	for i := 1; i < len(arr_reportname); i++ {
		str_added := ""

		if strings.Contains(arr_reportname[i], "org") {
			str_added = strings.Replace(arr_reportname[i], "org", str_OrgID, -1)
			fsqlstr = fsqlstr + str_added

		}

		if strings.Contains(arr_reportname[i], "member") {
			str_added = strings.Replace(arr_reportname[i], "member", tconf["member"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_member"] = "1"
		}

		if strings.Contains(arr_reportname[i], "group") {
			str_added = strings.Replace(arr_reportname[i], "group", tconf["group"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_group"] = "1"
		}

		if strings.Contains(arr_reportname[i], "datefrom") {
			str_added = strings.Replace(arr_reportname[i], "datefrom", tconf["datefrom"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_datefrom"] = "1"
		}

		if strings.Contains(arr_reportname[i], "dateto") {
			str_added = strings.Replace(arr_reportname[i], "dateto", tconf["dateto"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_dateto"] = "1"
		}

		if strings.Contains(arr_reportname[i], "partner") {
			str_added = strings.Replace(arr_reportname[i], "partner", tconf["partner"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_partner"] = "1"
		}

		if strings.Contains(arr_reportname[i], "loantype") {
			str_added = strings.Replace(arr_reportname[i], "loantype", tconf["loantype"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_loantype"] = "1"
		}

		if strings.Contains(arr_reportname[i], "customer") {
			str_added = strings.Replace(arr_reportname[i], "customer", tconf["customer"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_customer"] = "1"
		}

		if strings.Contains(arr_reportname[i], "paymode") {
			str_added = strings.Replace(arr_reportname[i], "paymode", tconf["paymode"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_paymode"] = "1"

		}
		if strings.Contains(arr_reportname[i], "bank") {
			str_added = strings.Replace(arr_reportname[i], "bank", tconf["bank"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_bank"] = "1"
		}

		if strings.Contains(arr_reportname[i], "filter") {
			str_added = strings.Replace(arr_reportname[i], "filter", tconf["filter"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_filter"] = "1"
		}

		if strings.Contains(arr_reportname[i], "trantype") {
			str_added = strings.Replace(arr_reportname[i], "trantype", tconf["trantype"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_trantype"] = "1"
		}

		if strings.Contains(arr_reportname[i], "status") {
			str_added = strings.Replace(arr_reportname[i], "status", tconf["status"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_status"] = "1"
		}

		if strings.Contains(arr_reportname[i], "coor") {
			str_added = strings.Replace(arr_reportname[i], "coor", tconf["coor"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_coor"] = "1"
		}

		if strings.Contains(arr_reportname[i], "laetype") {
			str_added = strings.Replace(arr_reportname[i], "laetype", tconf["laetype"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_laetype"] = "1"
		}

		if strings.Contains(arr_reportname[i], "adjtype") {
			str_added = strings.Replace(arr_reportname[i], "adjtype", tconf["aetype"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_aetype"] = "1"
		}

		if strings.Contains(arr_reportname[i], "asof") {
			str_added = strings.Replace(arr_reportname[i], "asof", tconf["asof"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_asof"] = "1"
		}

		if strings.Contains(arr_reportname[i], "lptype") {
			str_added = strings.Replace(arr_reportname[i], "lptype", tconf["lptype"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_lptype"] = "1"
		}

		if strings.Contains(arr_reportname[i], "othdtltype") {
			str_added = strings.Replace(arr_reportname[i], "othdtltype", tconf["othdtltype"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_othdtltype"] = "1"
		}

		if strings.Contains(arr_reportname[i], "partregioncoop") {
			str_added = strings.Replace(arr_reportname[i], "partregioncoop", tconf["partregioncoop"], -1)
			fsqlstr = fsqlstr + str_added
			tconf["show_partregioncoop"] = "1"
		}

	}
	fmt.Println(fsqlstr)
	fmt.Println(tconf["printtype"] + "eto")

	data, hdata := Datatable(fsqlstr, nil)
	var arrname []string
	type Htdata struct {
		Data    [][]string
		Hdata   []string
		Arrname []string
		Tconf   map[string]string
	}
	var_htdata := Htdata{data, hdata, nil, tconf}
	///rawreports_daily_lograw.html
	template_location := string(tconf["template"]) + ".html"
	var default_tpl_name = template_location
	var default_tpl_dir = "report/" + template_location

	if arr_reportname[0] == `lbr_rep_lograw 4, ` {
		default_tpl_name = "rawreports_daily_lograw.html"
		default_tpl_dir = "inventory/reports/rawreports/rawreports_daily_lograw.html"
		var tempname = "*****"
		for _, val := range data {
			if val[0] == tempname {
				tempname = val[0]
			}
			if val[0] != tempname {
				//fmt.Println(val)
				arrname = append(arrname, val[0])
				tempname = val[0]
			}

		}
		var_htdata = Htdata{data, hdata, arrname, tconf}
	}
	tmpl := template.New(default_tpl_name).Funcs(helper.FuncMap)
	var err error
	if tmpl, err = tmpl.ParseFiles(default_tpl_dir); err != nil {
		fmt.Println(err)
	}
	if tconf["export_type"] == "excel" {
		w.Header().Set("Content-Disposition", "attachment; filename=RAW_REPORT.xls")
		w.Header().Set("Content-Type", "application/excel")
		w.Header().Set("Content-Type", "application/vnd.ms-excel")

	} else if tconf["export_type"] == "word" {
		w.Header().Set("Content-Disposition", "attachment; filename=RAW_REPORT.doc")
		w.Header().Set("Content-Type", "application/msword")

	} else if tconf["export_type"] == "pdf" {
		w.Header().Set("Content-Disposition", "attachment; filename=RAW_REPORT.pdf")
		w.Header().Set("Content-Type", "application/pdf")
	}
	err1 := tmpl.Execute(w, var_htdata)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}

}

type Stuct_rawrports struct {
	Id         int
	Reportname string
	Sql        string
	Template   string
}

func func_get_rawreportslist(userID string) []Stuct_rawrports {
	row, err, _, _ := config.Ap_sql(`Reports_List 1,0`, 1)
	res := []Stuct_rawrports{}
	for row.Next() {
		var r Stuct_rawrports
		err = row.Scan(&r.Id, &r.Reportname, &r.Sql, &r.Template)
		if err != nil {
			panic(err.Error())
		}
		datarow := Stuct_rawrports{
			r.Id,
			r.Reportname,
			r.Sql,
			r.Template,
		}
		res = append(res, datarow)

	}
	return res
}

func Rawreports_nav(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	_, session_user_id := login.Get_account_info(r)
	session_user_id_str := strconv.Itoa(session_user_id)
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)

	type Htdata struct {
		RawreportsList    []Stuct_rawrports
		Ls_group          [][]string
		Ls_partner        [][]string
		Ls_loantype       [][]string
		Ls_paymode        [][]string
		Ls_bank           [][]string
		Ls_sloa           [][]string
		Ls_trantype       [][]string
		Ls_laetype        [][]string
		Ls_aetype         [][]string
		Ls_lptype         [][]string
		Ls_othdtltype     [][]string
		Ls_partregioncoop [][]string
	}

	group := datatables.DataList(`Partner_Groups_Get 1,0`)
	partner := datatables.DataList(`exec partners_get 2,` + str_OrgID + `,''`)
	loantype := datatables.DataList(`exec LoanTypes_Get 1,` + str_OrgID + `,0`)
	paymode := datatables.DataList(`exec Paymodes_Get 1, 0`)
	bank := datatables.DataList(`exec Banks_Get 1, ` + str_OrgID + `,''`)
	sloa := datatables.DataList(`exec Coors_Get 1,` + str_OrgID + `,0`)
	trantype := datatables.DataList(`exec TranType_Get 3 ,0`)
	laetype := datatables.DataList(`exec LAEType_Get 1, 0`)
	aetype := datatables.DataList(`exec AEType_Get 1, 0`)
	lptype := datatables.DataList(`exec LPType_Get 1, 0`)
	othdtltype := datatables.DataList(`exec OthDisDtl_trantype_List 0`)
	partregioncoop := datatables.DataList(`exec partners_get 5,` + str_OrgID + `,''`)

	var_htdata := Htdata{func_get_rawreportslist(session_user_id_str), group, partner, loantype, paymode, bank, sloa, trantype, laetype, aetype, lptype, othdtltype, partregioncoop}
	tmpl := template.New("rawreports_nav.html")
	var err error
	if tmpl, err = tmpl.ParseFiles("report/rawreports_nav.html"); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w, var_htdata)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}

}
