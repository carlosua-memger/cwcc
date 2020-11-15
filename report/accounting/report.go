package report

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	config "../../config"
	//"strconv"
	//"encoding/json"
	datatables "../../datatables"
	login "../../login"

	//"strings"
	"../../github.com/360EntSecGroup-Skylar/excelize"
	helper "../../helper"
)

type Profile struct {
	Message string // error status output or validated data
	Status  bool   //pass or fail
	ID      string
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
	tconf := make(map[string]string)
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
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
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
	rights := r.URL.Query().Get("rights")
	dec_rights := config.Decrypt(config.Encrypt_key, rights)
	tconf := make(map[string]string)
	tconf["ActionBut_setcv"] = "false"
	if strings.Contains(dec_rights, "SetCV") {
		tconf["ActionBut_setcv"] = "true"
	}
	type Htdata struct {
		RawreportsList []Stuct_rawrports
		Ls_group       [][]string
		Ls_partner     [][]string
		Ls_loantype    [][]string
		Ls_paymode     [][]string
		Ls_bank        [][]string
		Ls_lpartype    [][]string
		Ls_coop        [][]string
		Tempconf       map[string]string
	}

	group := datatables.DataList(`Partner_Groups_Get 1,0`)
	partner := datatables.DataList(`exec partners_get 2,` + str_OrgID + `,''`)
	loantype := datatables.DataList(`exec LoanTypes_Get 1,` + str_OrgID + `,0`)
	paymode := datatables.DataList(`exec Paymodes_Get 1, 0`)
	bank := datatables.DataList(`exec Banks_Get 1, ` + str_OrgID + `,''`)
	lpartype := datatables.DataList(`exec Mas_LoanSWTypes_Get_droplist 0`)
	coop := datatables.DataList(`exec Mas_Coops_Get_droplist 0`)

	var_htdata := Htdata{func_get_rawreportslist(session_user_id_str), group, partner, loantype, paymode, bank, lpartype, coop, tconf}
	tmpl := template.New("rawreports_nav.html")
	var err error
	if tmpl, err = tmpl.ParseFiles("report/accounting/rawreports_nav.html"); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w, var_htdata)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}

}

const (
	layoutISO = "01-02-2006"
	layoutUS  = "January 2, 2006"
)

func Bank_Advise(w http.ResponseWriter, r *http.Request) {
	dateSelected := r.URL.Query().Get("date")
	partner := r.URL.Query().Get("partner")
	partner_name := r.URL.Query().Get("partner_name")

	lpartype := r.URL.Query().Get("lpartype")
	lpartype_name := r.URL.Query().Get("lpartype_name")
	paymode := r.URL.Query().Get("paymode")
	//paymode_name := r.URL.Query().Get("paymode_name")
	bank := r.URL.Query().Get("bank")
	bank_name := r.URL.Query().Get("bank_name")
	group := r.URL.Query().Get("group")
	//group_name := r.URL.Query().Get("group_name")
	coop := r.URL.Query().Get("coop")
	//coop_name := r.URL.Query().Get("coop_name")
	member := r.URL.Query().Get("member")

	username, _ := login.Get_account_info(r)
	t, _ := time.Parse(layoutISO, dateSelected)

	fmt.Println(`Accounting_BankAdvice 1,'` + dateSelected + `',` + partner + `,` + member + `,` + lpartype + `,` + paymode + `,` + bank + `,` + group + `,` + coop)
	arr_data := datatables.DataList(`Accounting_BankAdvice 1,'` + dateSelected + `',` + partner + `,` + member + `,` + lpartype + `,` + paymode + `,` + bank + `,` + group + `,` + coop)

	type Data struct {
		Tconf      map[string]string
		Arr_data   [][]string
		MemberName string
	}
	var str_page string
	var str_page_dir string
	str_page = `bank_advise.html`
	str_page_dir = `report/accounting/bank_advise.html`
	tconf := make(map[string]string)
	tconf["username"] = username
	tconf["bank_name"] = bank_name
	tconf["lpartype_name"] = lpartype_name
	tconf["date"] = t.Format(layoutUS)
	tmpl := template.New(str_page).Funcs(helper.FuncMap)
	var err error
	if tmpl, err = tmpl.ParseFiles(str_page_dir); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w, Data{tconf, arr_data, partner_name})
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func Debit_Memo(w http.ResponseWriter, r *http.Request) {
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	tconf := make(map[string]string)
	tconf["org_id"] = str_OrgID

	username, _ := login.Get_account_info(r)
	user_id, _ := login.Get_account_info_id(r)

	dateSelected := r.URL.Query().Get("date")
	partner := r.URL.Query().Get("partner")
	lpartype := r.URL.Query().Get("lpartype")
	paymode := r.URL.Query().Get("paymode")
	//paymode_name := r.URL.Query().Get("paymode_name")
	bank := r.URL.Query().Get("bank")
	group := r.URL.Query().Get("group")
	//group_name := r.URL.Query().Get("group_name")
	coop := r.URL.Query().Get("coop")
	t, _ := time.Parse(layoutISO, dateSelected)

	tconf["dateSelected"] = dateSelected
	tconf["partner"] = partner

	tconf["lpartype"] = lpartype

	member := r.URL.Query().Get("member")
	tconf["member"] = member
	tconf["paymode"] = paymode
	tconf["bank"] = bank

	tconf["group"] = group
	tconf["coop"] = coop
	tconf["date"] = t.Format(layoutUS)

	tconf["username"] = username
	tconf["user_id"] = strconv.Itoa(user_id)
	fmt.Println("on printing")
	fmt.Println(tconf["sql_code"])
	fmt.Println(tconf["rpt_url"])
	//fmt.Println(path.Base(tconf["rpt_url"]))
	reportUrl := "report/accounting/debit_memo.html"
	arr_reportname := strings.Split(string(reportUrl), `/`)
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
	var default_tpl_dir = string(reportUrl)

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

func Transmittal(w http.ResponseWriter, r *http.Request) {
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	tconf := make(map[string]string)
	tconf["org_id"] = str_OrgID

	username, _ := login.Get_account_info(r)
	user_id, _ := login.Get_account_info_id(r)

	dateSelected := r.URL.Query().Get("date")
	partner := r.URL.Query().Get("partner")
	lpartype := r.URL.Query().Get("lpartype")
	paymode := r.URL.Query().Get("paymode")
	//paymode_name := r.URL.Query().Get("paymode_name")
	bank := r.URL.Query().Get("bank")
	bank_name := r.URL.Query().Get("bank_name")
	group := r.URL.Query().Get("group")
	//group_name := r.URL.Query().Get("group_name")
	coop := r.URL.Query().Get("coop")
	coop_name := r.URL.Query().Get("coop_name")
	t, _ := time.Parse(layoutISO, dateSelected)
	tconf["dateSelected"] = dateSelected
	tconf["partner"] = partner

	member := r.URL.Query().Get("member")
	tconf["member"] = member
	tconf["lpartype"] = lpartype

	tconf["paymode"] = paymode
	tconf["bank"] = bank
	tconf["bank_name"] = bank_name

	tconf["group"] = group
	tconf["coop"] = coop
	tconf["coop_name"] = coop_name
	tconf["date"] = t.Format(layoutUS)

	tconf["username"] = username
	tconf["user_id"] = strconv.Itoa(user_id)
	fmt.Println("on printing")
	fmt.Println(tconf["sql_code"])
	fmt.Println(tconf["rpt_url"])
	//fmt.Println(path.Base(tconf["rpt_url"]))
	reportUrl := "report/accounting/transmittal.html"
	arr_reportname := strings.Split(string(reportUrl), `/`)
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
	var default_tpl_dir = string(reportUrl)

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

func Transmittal2(w http.ResponseWriter, r *http.Request) {
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	tconf := make(map[string]string)
	tconf["org_id"] = str_OrgID

	username, _ := login.Get_account_info(r)
	user_id, _ := login.Get_account_info_id(r)

	dateSelected := r.URL.Query().Get("date")
	partner := r.URL.Query().Get("partner")
	lpartype := r.URL.Query().Get("lpartype")
	paymode := r.URL.Query().Get("paymode")
	//paymode_name := r.URL.Query().Get("paymode_name")
	bank := r.URL.Query().Get("bank")
	bank_name := r.URL.Query().Get("bank_name")
	group := r.URL.Query().Get("group")
	//group_name := r.URL.Query().Get("group_name")
	coop := r.URL.Query().Get("coop")
	coop_name := r.URL.Query().Get("coop_name")
	member := r.URL.Query().Get("member")
	tconf["member"] = member
	t, _ := time.Parse(layoutISO, dateSelected)
	tconf["dateSelected"] = dateSelected
	tconf["partner"] = partner

	tconf["lpartype"] = lpartype

	tconf["paymode"] = paymode
	tconf["bank"] = bank
	tconf["bank_name"] = bank_name

	tconf["group"] = group
	tconf["coop"] = coop
	tconf["coop_name"] = coop_name
	tconf["date"] = t.Format(layoutUS)

	tconf["username"] = username
	tconf["user_id"] = strconv.Itoa(user_id)
	fmt.Println("on printing")
	fmt.Println(tconf["sql_code"])
	fmt.Println(tconf["rpt_url"])
	//fmt.Println(path.Base(tconf["rpt_url"]))
	reportUrl := "report/accounting/transmittal_2.html"
	arr_reportname := strings.Split(string(reportUrl), `/`)
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
	var default_tpl_dir = string(reportUrl)

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

func DMSummary(w http.ResponseWriter, r *http.Request) {
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	tconf := make(map[string]string)
	tconf["org_id"] = str_OrgID

	username, _ := login.Get_account_info(r)
	user_id, _ := login.Get_account_info_id(r)

	dateSelected := r.URL.Query().Get("date")
	partner := r.URL.Query().Get("partner")
	lpartype := r.URL.Query().Get("lpartype")
	paymode := r.URL.Query().Get("paymode")
	//paymode_name := r.URL.Query().Get("paymode_name")
	bank := r.URL.Query().Get("bank")
	bank_name := r.URL.Query().Get("bank_name")
	group := r.URL.Query().Get("group")
	//group_name := r.URL.Query().Get("group_name")
	coop := r.URL.Query().Get("coop")
	coop_name := r.URL.Query().Get("coop_name")
	t, _ := time.Parse(layoutISO, dateSelected)
	member := r.URL.Query().Get("member")
	tconf["member"] = member
	tconf["dateSelected"] = dateSelected
	tconf["partner"] = partner

	tconf["lpartype"] = lpartype

	tconf["paymode"] = paymode
	tconf["bank"] = bank
	tconf["bank_name"] = bank_name

	tconf["group"] = group
	tconf["coop"] = coop
	tconf["coop_name"] = coop_name
	tconf["date"] = t.Format(layoutUS)

	tconf["username"] = username
	tconf["user_id"] = strconv.Itoa(user_id)
	fmt.Println("on printing")
	fmt.Println(tconf["sql_code"])
	fmt.Println(tconf["rpt_url"])
	//fmt.Println(path.Base(tconf["rpt_url"]))
	reportUrl := "report/accounting/dm_summary.html"
	arr_reportname := strings.Split(string(reportUrl), `/`)
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
	var default_tpl_dir = string(reportUrl)

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

var (
	file                    *os.File
	part                    []byte
	prefix                  bool
	strdns                  string
	sql_sp                  string
	delimit                 string
	branchcode              string
	sp_import               string
	delete_after_restore    bool
	copy_to_temp_folder     bool
	copy_to_temp_folder_dir string
	trandate                string
	partner_id              string
	total                   int
	count                   int
	fpath                   string
)

func Download(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("hello")
	login.Session_renew(w, r)

	format := r.URL.Query().Get("format")
	fileextension := r.URL.Query().Get("fileextension")

	if format == "100" {
		DEFAULT(w, r)
	} else if fileextension == "xls" {
		XLS(w, r)
	}

}

func XLS(w http.ResponseWriter, r *http.Request) {

	login.Session_renew(w, r)

	partner := r.URL.Query().Get("partner")
	lpartype := r.URL.Query().Get("lpartype")
	transdate := r.URL.Query().Get("transdate")
	bank := r.URL.Query().Get("bank")
	group := r.URL.Query().Get("group")
	coop := r.URL.Query().Get("coop")
	member := r.URL.Query().Get("member")
	filename := r.URL.Query().Get("filename")
	fileextension := r.URL.Query().Get("fileextension")
	//sequence := r.URL.Query().Get("sequence")
	//code := r.URL.Query().Get("code")
	//t, err := time.Parse("01-02-2006", transdate)
	t := time.Now()
	DatetimeFileAddname := t.Format("20060102")
	fname := filename + "_" + DatetimeFileAddname + "." + fileextension
	fpath := "./report/xls_archive/" + fname
	sql_sp = `Bank_Textfile 1,'` + transdate + `',` + partner + `,` + bank + `,` + "0" + `,` + group + `,` + member + `,` + lpartype + `,` + coop

	data, hdata := Datatable(sql_sp, nil)

	type Htdata struct {
		Data    [][]string
		Hdata   []string
		Arrname []string
	}

	//var_htdata := Htdata{data, hdata, nil}

	Excelize(data, hdata, fpath)

	f1, err := os.Open(fpath)
	if err != nil {
		// handle error
	}
	defer f1.Close()
	fi, err := f1.Stat()
	if err != nil {
		// handle error
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fname)
	http.ServeContent(w, r, fname, fi.ModTime(), f1)

	//os.Remove(fpath)

	//fmt.Println("HTDATA: ", var_htdata)

	///rawreports_daily_lograw.html
	// template_location := "excel_template.html"
	// var default_tpl_name = template_location
	// var default_tpl_dir = "report/accounting/" + template_location

	// tmpl := template.New(default_tpl_name).Funcs(helper.FuncMap)

	// if tmpl, err = tmpl.ParseFiles(default_tpl_dir); err != nil {
	// 	fmt.Println(err)
	// }

	// w.Header().Set("Content-Disposition", "attachment; filename="+fpath)
	// w.Header().Set("Content-Type", "application/excel")
	// w.Header().Set("Content-Type", "application/vnd.ms-excel")
	// w.Header().Set("Content-Type", "application/x-excel")
	// w.Header().Set("Content-Type", "application/x-msexcel")

	// w.Write([])
	// err1 := tmpl.Execute(w, var_htdata)
	// if err1 != nil {
	// 	http.Error(w, err1.Error(), http.StatusInternalServerError)
	// }

}

func Excelize(data [][]string, hdata []string, fpath string) {
	f := excelize.NewFile()

	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	s := "Sheet1"
	style, _ := f.NewStyle(`{"alignment":{"horizontal": "center"},"font":{"bold": true}}`)

	//SET CELLS FOR HEADERS WITH DESIGN
	for i := 0; i < len(hdata); i++ {
		nextCell := string(alphabet[i]) + "1"
		f.SetCellValue(s, nextCell, hdata[i])
		f.SetCellStyle(s, nextCell, nextCell, style)
	}

	//SET CELLS FOR DATA
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {

			float, err := strconv.ParseFloat(data[i][j], 32)
			if err == nil && strings.Contains(data[i][j], ".") {
				f.SetCellValue(s, string("ABCDEFGHIJKLMNOPQRSTUVWXYZ"[j])+strconv.Itoa(i+2), float)
			} else {
				f.SetCellValue(s, string("ABCDEFGHIJKLMNOPQRSTUVWXYZ"[j])+strconv.Itoa(i+2), data[i][j])
			}

		}
	}

	// Save xlsx file by the given path.
	if err := f.SaveAs(fpath); err != nil {
		fmt.Println(err)
	}
}

func DEFAULT(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)

	partner := r.URL.Query().Get("partner")
	lpartype := r.URL.Query().Get("lpartype")
	transdate := r.URL.Query().Get("transdate")
	bank := r.URL.Query().Get("bank")
	group := r.URL.Query().Get("group")
	coop := r.URL.Query().Get("coop")
	member := r.URL.Query().Get("member")
	filename := r.URL.Query().Get("filename")
	fileextension := r.URL.Query().Get("fileextension")
	sequence := r.URL.Query().Get("sequence")
	code := r.URL.Query().Get("code")
	t, err := time.Parse("01-02-2006", transdate)
	//t := time.Now()
	DatetimeFileAddname := t.Format("20060102")
	Day := t.Format("02")
	//a := fmt.Sprint(lines[1])

	if bank == "40" {
		fpath = "S275" + code + Day + ".TX" + sequence
	} else if fileextension == "01" {
		fpath = filename + "." + lpad(sequence, "0", 2)
	} else {
		fpath = filename + "_" + DatetimeFileAddname + "." + fileextension
	}
	fmt.Println(transdate)
	tconf := make(map[string]string)
	tconf["filename"] = fpath
	sql_sp = `Bank_Textfile 1,'` + transdate + `',` + partner + `,` + bank + `,` + "0" + `,` + group + `,` + member + `,` + lpartype + `,` + coop + `,` + sequence

	//Datagrid(downloadfilename)

	rows, err, _, _ := config.Ap_sql(sql_sp, 1)
	fmt.Println(sql_sp)
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
	
	newFile, err := os.Create("banktextfile/" + fpath)
	
	if err != nil {
		fmt.Println(err)
		newFile.Close()
		return
	}

	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//		var a []string

		var lineString string
		for i := 0; i < len(cols); i++ { // fetching data tables to container map array

			//_, err = file.WriteString(strings.TrimSpace(PrintValue(vals[i].(*interface{}))))
			//	a  = append(a, PrintValue(vals[i].(*interface{})) )
			fmt.Println(strings.TrimSpace(PrintValue(vals[i].(*interface{}))))
			lineString += PrintValue(vals[i].(*interface{}))
			if i != (len(cols) - 1) {
				//_, err = file.WriteString(strings.TrimSpace("") + "," + "\t")
				lineString += ","
			}

		}
		fmt.Println("line string :%v", lineString)
		_, err = fmt.Fprintln(newFile, lineString+"\r")
		if err != nil {
			fmt.Println("the  error was here %v", err)
			return
		}
	}
	err = newFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("file written successfully")

	if rows.Err() != nil {
		fmt.Println(rows.Err)
	}

	ForceDownload(w, r,"banktextfile/"+fpath,fpath)

}

func Ap_sql(sq string, sq_type int) (*sql.Rows, error, bool, sql.Result) {

	dsn := strdns

	db, err := sql.Open("mssql", dsn)
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		defer db.Close()
		return nil, err, false, nil
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		defer db.Close()
		return nil, err, false, nil
	}
	if sq_type == 1 {
		rows, err := db.Query(sq)
		if err != nil {
			fmt.Println("option 1 err: ", err)
		}
		return rows, err, false, nil
	} else if sq_type == 2 {
		rows, err := db.Query(sq)
		if err != nil {
			fmt.Println("option 2 err: ", err)
		}
		return rows, err, false, nil
	} else {
		res, err := db.Exec(sq)
		if err != nil {
			fmt.Println("option 3 err: ", err)
		}
		return nil, err, true, res
	}

	fmt.Println(err)
	defer db.Close()
	//fmt.Println("%v", err)
	return nil, err, true, nil

}

func lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}

func stringToFloat(val string) (*float64, error) {
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil, err
	}
	return &num, nil
}

func ForceDownload(w http.ResponseWriter, r *http.Request, path string, file string) {

	downloadBytes, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
	}

	// set the default MIME type to send
	mime := http.DetectContentType(downloadBytes)

	fileSize := len(string(downloadBytes))

	// Generate the server headers
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", "attachment; filename="+file+"")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(fileSize))
	w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

	//b := bytes.NewBuffer(downloadBytes)
	//if _, err := b.WriteTo(w); err != nil {
	//              fmt.Fprintf(w, "%s", err)
	//      }

	// force it down the client's.....
	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(downloadBytes))

}

func SetCV(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	login.Session_renew(w, r)

	dateSelected := r.URL.Query().Get("date")
	partner := r.URL.Query().Get("partner")
	lpartype := r.URL.Query().Get("lpartype")
	paymode := r.URL.Query().Get("paymode")
	//paymode_name := r.URL.Query().Get("paymode_name")
	bank := r.URL.Query().Get("bank")
	group := r.URL.Query().Get("group")
	//group_name := r.URL.Query().Get("group_name")
	coop := r.URL.Query().Get("coop")
	member := r.URL.Query().Get("member")
	cv := r.URL.Query().Get("cv")

	fmt.Println(`GOOD TO POST`)
	arr_raw_data := datatables.DataList(`Accounting_SetCV 1,'` + dateSelected + `',` + partner + `,` + member + `,` + lpartype + `,` + paymode + `,` + bank + `,` + group + `,` + coop + `,` + cv)
	fmt.Println(arr_raw_data)

	js, err := json.Marshal(arr_raw_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
