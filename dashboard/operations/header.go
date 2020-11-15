package product_movement

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	config "../../config"
	datatables "../../datatables"
	helper "../../helper"
	login "../../login"
)

type Profile struct {
	Message string // error status output or validated data
	Status  bool   //pass or fail
	ID      string
}

func ViewitemHandler(w http.ResponseWriter, r *http.Request) {
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	tconf := make(map[string]string)

	//tconf["img_name"] = r.URL.Query().Get("img_name")
	//tconf["item_id"] = r.URL.Query().Get("item_id")
	//tconf["graph_title"] = r.URL.Query().Get("item_name")
	//tconf["item_image"] = config.Img_to_render_custom1( tconf["img_name"] ,171)
	//tconf["item_image"] = ""

	//tconf["item_name"] = r.URL.Query().Get("item_name")
	//tconf["startdate"] = r.URL.Query().Get("startdate")
	//tconf["enddate"] = r.URL.Query().Get("enddate")

	tconf["ddata"] = ""
	fmt.Println(tconf["ddata"])
	arr_ddata := strings.Split(tconf["ddata"], `][`)
	//tconf["item_image"] 			= "config.Img_to_render_custom1(arr_ddata[0],171)"

	tconf["item_image"] = strings.Replace(r.URL.Query().Get("item_img"), "size_80", "size_171", 2)
	fmt.Println(`img src : `, tconf["item_image"])
	tconf["item_id_hideme"] = " "
	tconf["supplier_name"] = " "
	tconf["beg_vol"] = " "
	tconf["suprec_vol"] = " "
	tconf["cussal_vol"] = " "
	tconf["adjent_vol"] = " "
	tconf["end_vol"] = " "

	fmt.Println(arr_ddata)
	fmt.Println(`EOD_LoanPastDueAgingReport_Show 1, ` + str_OrgID)

	//EOD_LoanPastDueAgingReport_Show 1, 1, 0
	//
	//
	/* 	fmt.Println(`EOD_LoanPastDueAgingReport_Show 1, ` + str_OrgID)
	   rowret, err, _, _ := config.Ap_sql(`EOD_LoanPastDueAgingReport_Show 1, `+str_OrgID, 1)
	   if err != nil {
		   panic(err.Error())
	   }
	   type Graph2 struct {
		   Date_s time.Time
		   Vals   string
		   Vals2  string
	   }

	   var RG2 Graph2
	   if err != nil {
		   fmt.Println(err.Error)
		   panic(err.Error())
	   }
	   tconf["graph2_data"] = ` [`
	   tconf["graph2_data2"] = ` [`

	   for rowret.Next() {
		   err = rowret.Scan(&RG2.Date_s, &RG2.Vals, &RG2.Vals2)
		   if err != nil {
			   panic(err.Error())
		   }
		   fmt.Println(RG2.Date_s)
		   fmt.Println(RG2.Vals)
		   str_dd := strconv.Itoa(RG2.Date_s.Day())
		   var i_mm int = int(RG2.Date_s.Month()) - 1
		   str_mm := strconv.Itoa(i_mm)
		   str_yyyy := strconv.Itoa(RG2.Date_s.Year())
		   fmt.Println(`Day:`, str_dd)
		   fmt.Println(`MM:`, str_mm)
		   fmt.Println(`YYYY:`, str_yyyy)

		   tconf["graph2_data"] += `[Date.UTC(` + str_yyyy + `, ` + str_mm + `, ` + str_dd + `), ` + RG2.Vals + `],`
		   tconf["graph2_data2"] += `[Date.UTC(` + str_yyyy + `, ` + str_mm + `, ` + str_dd + `), ` + RG2.Vals2 + `],`

		   //r = RetMessage{r.Date_s,r.Date_s}
	   }
	   tconf["graph2_data"] += `]`
	   tconf["graph2_data2"] += `]` */
	var err error
	reportUrl := "dashboard/operations/view_item_stat.html"
	arr_reportname := strings.Split(string(reportUrl), `/`)
	var default_tpl_name = arr_reportname[len(arr_reportname)-1]
	//var default_tpl_dir = string(reportUrl)

	tmpl := template.New(default_tpl_name).Funcs(helper.FuncMap)

	//var err error
	if tmpl, err = tmpl.ParseFiles("dashboard/operations/view_item_stat.html"); err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, tconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DashboardLoanPorfolio(w http.ResponseWriter, r *http.Request) {

	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	tconf := make(map[string]string)

	//tconf["img_name"] = r.URL.Query().Get("img_name")
	//tconf["item_id"] = r.URL.Query().Get("item_id")
	//tconf["graph_title"] = r.URL.Query().Get("item_name")
	//tconf["item_image"] = config.Img_to_render_custom1( tconf["img_name"] ,171)
	//tconf["item_image"] = ""

	//tconf["item_name"] = r.URL.Query().Get("item_name")
	//tconf["startdate"] = r.URL.Query().Get("startdate")
	//tconf["enddate"] = r.URL.Query().Get("enddate")

	tconf["ddata"] = ""
	fmt.Println(tconf["ddata"])
	arr_ddata := strings.Split(tconf["ddata"], `][`)
	//tconf["item_image"] 			= "config.Img_to_render_custom1(arr_ddata[0],171)"

	tconf["item_image"] = strings.Replace(r.URL.Query().Get("item_img"), "size_80", "size_171", 2)
	fmt.Println(`img src : `, tconf["item_image"])
	tconf["item_id_hideme"] = " "
	tconf["supplier_name"] = " "
	tconf["beg_vol"] = " "
	tconf["suprec_vol"] = " "
	tconf["cussal_vol"] = " "
	tconf["adjent_vol"] = " "
	tconf["end_vol"] = " "

	fmt.Println(arr_ddata)
	fmt.Println(`EOD_LoanPastDueAgingReport_Show 1, ` + str_OrgID)
	graph3_Data := datatables.DataList(`EOD_LoanRiskPortfolioReport_Show 1, ` + str_OrgID)

	fmt.Println(graph3_Data)
	tconf["partner"] = ` [`
	tconf["par_percentage"] = ` [`
	tconf["partner_id"] = ` [`

	for _, value := range graph3_Data {
		fmt.Println(value[0])
		tconf["par_percentage"] += `{type: 'column', y:` + value[1] + `, partner_id:` + value[4] + `,olr:` + value[2] + `,par:` + value[3] + `,partner_name:'` + value[0] + `'},`

		tconf["partner"] += `'` + value[0] + `',`

		tconf["partner_id"] += `[` + value[4] + `],`

	}

	tconf["partner"] += `]`
	tconf["par_percentage"] += ` ]`
	tconf["partner_id"] += ` ]`
	tconf["par_partner_url"] = "/dashboard/operations/par_partner?partner_id="

	var err error
	tmpl := template.New("dashboard_loanportfolio.html")
	//var err error
	if tmpl, err = tmpl.ParseFiles("dashboard/operations/dashboard_loanportfolio.html"); err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, tconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DashboardLoanRelease(w http.ResponseWriter, r *http.Request) {
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	tconf := make(map[string]string)

	//tconf["img_name"] = r.URL.Query().Get("img_name")
	//tconf["item_id"] = r.URL.Query().Get("item_id")
	//tconf["graph_title"] = r.URL.Query().Get("item_name")
	//tconf["item_image"] = config.Img_to_render_custom1( tconf["img_name"] ,171)
	//tconf["item_image"] = ""

	//tconf["item_name"] = r.URL.Query().Get("item_name")
	//tconf["startdate"] = r.URL.Query().Get("startdate")
	//tconf["enddate"] = r.URL.Query().Get("enddate")

	tconf["ddata"] = ""
	fmt.Println(tconf["ddata"])
	arr_ddata := strings.Split(tconf["ddata"], `][`)
	//tconf["item_image"] 			= "config.Img_to_render_custom1(arr_ddata[0],171)"

	tconf["item_image"] = strings.Replace(r.URL.Query().Get("item_img"), "size_80", "size_171", 2)
	fmt.Println(`img src : `, tconf["item_image"])
	tconf["item_id_hideme"] = " "
	tconf["supplier_name"] = " "
	tconf["beg_vol"] = " "
	tconf["suprec_vol"] = " "
	tconf["cussal_vol"] = " "
	tconf["adjent_vol"] = " "
	tconf["end_vol"] = " "

	fmt.Println(arr_ddata)
	fmt.Println(`EOD_LoanPastDueAgingReport_Show 1, ` + str_OrgID)
	graph2_Data := datatables.DataList(`EOD_LoanApplicationReleaseReport_Show 1, ` + str_OrgID)

	fmt.Println(graph2_Data)
	tconf["date"] = ` [`
	tconf["la_value"] = ` [`
	tconf["la_volume"] = ` [`
	tconf["lr_value"] = ` [`
	tconf["lr_volume"] = ` [`

	for _, value := range graph2_Data {
		fmt.Println(value[0])
		tconf["la_value"] += `{y:` + value[1] + `, url:'datefrom=` + value[5] + `&dateto=` + value[5] + `&reportname=LR_Register%203,@@org,%20%27@@datefrom%27,%20%27@@dateto%27,%20@@partner&export_type=print&group=0&group_name=All&partner=0&partner_name=All&loantype=0&loantype_name=All&paymode=0&paymode_name=All&bank=0&bank_name=All&filter_name=&status=0&status_name=All&coor=0&coor_name=All&trantype=0&trantype_name=All&template=default_template&report_name=Register:%20Loan%20Release%20Detail&member=&member_name=undefined&laetype=9laetype_name=Adjustment'},`
		tconf["la_volume"] += `{y:` + value[2] + `, url:'datefrom=` + value[5] + `&dateto=` + value[5] + `&reportname=LR_Register%203,@@org,%20%27@@datefrom%27,%20%27@@dateto%27,%20@@partner&export_type=print&group=0&group_name=All&partner=0&partner_name=All&loantype=0&loantype_name=All&paymode=0&paymode_name=All&bank=0&bank_name=All&filter_name=&status=0&status_name=All&coor=0&coor_name=All&trantype=0&trantype_name=All&template=default_template&report_name=Register:%20Loan%20Release%20Detail&member=&member_name=undefined&laetype=9laetype_name=Adjustment'},`
		tconf["lr_value"] += `{y:` + value[3] + `, url:'datefrom=` + value[5] + `&dateto=` + value[5] + `&reportname=LR_Register%203,@@org,%20%27@@datefrom%27,%20%27@@dateto%27,%20@@partner&export_type=print&group=0&group_name=All&partner=0&partner_name=All&loantype=0&loantype_name=All&paymode=0&paymode_name=All&bank=0&bank_name=All&filter_name=&status=0&status_name=All&coor=0&coor_name=All&trantype=0&trantype_name=All&template=default_template&report_name=Register:%20Loan%20Release%20Detail&member=&member_name=undefined&laetype=9laetype_name=Adjustment'},`
		tconf["lr_volume"] += `{y:` + value[4] + `, url:'datefrom=` + value[5] + `&dateto=` + value[5] + `&reportname=LR_Register%203,@@org,%20%27@@datefrom%27,%20%27@@dateto%27,%20@@partner&export_type=print&group=0&group_name=All&partner=0&partner_name=All&loantype=0&loantype_name=All&paymode=0&paymode_name=All&bank=0&bank_name=All&filter_name=&status=0&status_name=All&coor=0&coor_name=All&trantype=0&trantype_name=All&template=default_template&report_name=Register:%20Loan%20Release%20Detail&member=&member_name=undefined&laetype=9laetype_name=Adjustment'},`

		tconf["date"] += `'` + value[0] + `',`

	}

	tconf["date"] += `]`
	tconf["la_value"] += ` ]`
	tconf["la_volume"] += ` ]`
	tconf["lr_value"] += ` ]`
	tconf["lr_volume"] += ` ]`

	var err error
	tmpl := template.New("dashboard_loanrelease.html")
	//var err error
	if tmpl, err = tmpl.ParseFiles("dashboard/operations/dashboard_loanrelease.html"); err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, tconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func DashboardLoanAging(w http.ResponseWriter, r *http.Request) {
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	tconf := make(map[string]string)

	//tconf["img_name"] = r.URL.Query().Get("img_name")
	//tconf["item_id"] = r.URL.Query().Get("item_id")
	//tconf["graph_title"] = r.URL.Query().Get("item_name")
	//tconf["item_image"] = config.Img_to_render_custom1( tconf["img_name"] ,171)
	//tconf["item_image"] = ""

	//tconf["item_name"] = r.URL.Query().Get("item_name")
	//tconf["startdate"] = r.URL.Query().Get("startdate")
	//tconf["enddate"] = r.URL.Query().Get("enddate")

	tconf["ddata"] = ""
	fmt.Println(tconf["ddata"])
	arr_ddata := strings.Split(tconf["ddata"], `][`)
	//tconf["item_image"] 			= "config.Img_to_render_custom1(arr_ddata[0],171)"

	tconf["item_image"] = strings.Replace(r.URL.Query().Get("item_img"), "size_80", "size_171", 2)
	fmt.Println(`img src : `, tconf["item_image"])
	tconf["item_id_hideme"] = " "
	tconf["supplier_name"] = " "
	tconf["beg_vol"] = " "
	tconf["suprec_vol"] = " "
	tconf["cussal_vol"] = " "
	tconf["adjent_vol"] = " "
	tconf["end_vol"] = " "

	fmt.Println(arr_ddata)
	fmt.Println(`EOD_LoanPastDueAgingReport_Show 1, ` + str_OrgID)
	graph1_Data := datatables.DataList(`EOD_LoanPastDueAgingReport_Show 1, ` + str_OrgID)

	fmt.Println(graph1_Data)
	tconf["age130"] = ` [`
	tconf["age3160"] = ` [`
	tconf["age6190"] = ` [`
	tconf["age91120"] = ` [`
	tconf["age121150"] = ` [`
	tconf["age151up"] = ` [`
	tconf["graph1_data"] = ` [`
	tconf["graph1_data2"] = ` [`
	tconf["categories"] = ` [`
	for _, value := range graph1_Data {
		fmt.Println(value[0])
		tconf["age130"] += `[` + value[1] + `],`
		tconf["age3160"] += `[` + value[2] + `],`
		tconf["age6190"] += `[` + value[3] + `],`
		tconf["age91120"] += `[` + value[4] + `],`
		tconf["age121150"] += `[` + value[5] + `],`
		tconf["age151up"] += `[` + value[6] + `],`
		tconf["categories"] += `'` + value[0] + `',`

		tconf["graph1_data"] += `[` + value[7] + `,` + value[6] + `],`

		//tconf["graph1_data"] += `['` + value[0] + `', ` + value[7] + `,` + value[6] + `],`
		tconf["graph1_data2"] += `['` + value[1] + `', ` + value[3] + `],`

	}
	tconf["graph1_data"] += `]`
	tconf["graph1_data2"] += `]`
	tconf["categories"] += `]`
	tconf["age130"] += ` ]`
	tconf["age3160"] += ` ]`
	tconf["age6190"] += ` ]`
	tconf["age91120"] += ` ]`
	tconf["age121150"] += ` ]`
	tconf["age151up"] += ` ]`

	var err error
	tmpl := template.New("dashboard_loanaging.html")
	//var err error
	if tmpl, err = tmpl.ParseFiles("dashboard/operations/dashboard_loanaging.html"); err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, tconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Dashboard_parpartner(w http.ResponseWriter, r *http.Request) {
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	tconf := make(map[string]string)

	tconf["partner_id"] = r.URL.Query().Get("partner_id")
	//tconf["item_id"] = r.URL.Query().Get("item_id")
	//tconf["graph_title"] = r.URL.Query().Get("item_name")
	//tconf["item_image"] = config.Img_to_render_custom1( tconf["img_name"] ,171)
	//tconf["item_image"] = ""

	//tconf["item_name"] = r.URL.Query().Get("item_name")
	//tconf["startdate"] = r.URL.Query().Get("startdate")
	//tconf["enddate"] = r.URL.Query().Get("enddate")

	tconf["ddata"] = ""
	fmt.Println(tconf["ddata"])
	arr_ddata := strings.Split(tconf["ddata"], `][`)
	//tconf["item_image"] 			= "config.Img_to_render_custom1(arr_ddata[0],171)"

	/* tconf["item_image"] = strings.Replace(r.URL.Query().Get("item_img"), "size_80", "size_171", 2)
	fmt.Println(`img src : `, tconf["item_image"])
	tconf["item_id_hideme"] = " "
	tconf["supplier_name"] = " "
	tconf["beg_vol"] = " "
	tconf["suprec_vol"] = " "
	tconf["cussal_vol"] = " "
	tconf["adjent_vol"] = " "
	tconf["end_vol"] = " " */

	fmt.Println(arr_ddata)
	fmt.Println(`EOD_LoanRiskPortfolioReport_Show  2, ` + str_OrgID)
	graph4_Data := datatables.DataList(`EOD_LoanRiskPortfolioReport_Show  2, ` + str_OrgID + `,` + tconf["partner_id"])

	fmt.Println(graph4_Data)
	tconf["partner2"] = ` [`
	tconf["par_percentage2"] = ` [`
	tconf["partner_id2"] = ` [`

	for _, value := range graph4_Data {
		fmt.Println(value[0])
		tconf["par_percentage2"] += `{y:` + value[1] + `,olr:` + value[2] + `,par:` + value[3] + `,partner_name:'` + value[0] + `'},`

		tconf["partner2"] += `'` + value[0] + `',`

		tconf["partner_id2"] += `[` + value[3] + `],`

	}

	tconf["partner2"] += `]`
	tconf["par_percentage2"] += ` ]`
	tconf["partner_id2"] += ` ]`
	//EOD_LoanPastDueAgingReport_Show 1, 1, 0
	//
	//
	/* 	fmt.Println(`EOD_LoanPastDueAgingReport_Show 1, ` + str_OrgID)
	   rowret, err, _, _ := config.Ap_sql(`EOD_LoanPastDueAgingReport_Show 1, `+str_OrgID, 1)
	   if err != nil {
		   panic(err.Error())
	   }
	   type Graph2 struct {
		   Date_s time.Time
		   Vals   string
		   Vals2  string
	   }

	   var RG2 Graph2
	   if err != nil {
		   fmt.Println(err.Error)
		   panic(err.Error())
	   }
	   tconf["graph2_data"] = ` [`
	   tconf["graph2_data2"] = ` [`

	   for rowret.Next() {
		   err = rowret.Scan(&RG2.Date_s, &RG2.Vals, &RG2.Vals2)
		   if err != nil {
			   panic(err.Error())
		   }
		   fmt.Println(RG2.Date_s)
		   fmt.Println(RG2.Vals)
		   str_dd := strconv.Itoa(RG2.Date_s.Day())
		   var i_mm int = int(RG2.Date_s.Month()) - 1
		   str_mm := strconv.Itoa(i_mm)
		   str_yyyy := strconv.Itoa(RG2.Date_s.Year())
		   fmt.Println(`Day:`, str_dd)
		   fmt.Println(`MM:`, str_mm)
		   fmt.Println(`YYYY:`, str_yyyy)

		   tconf["graph2_data"] += `[Date.UTC(` + str_yyyy + `, ` + str_mm + `, ` + str_dd + `), ` + RG2.Vals + `],`
		   tconf["graph2_data2"] += `[Date.UTC(` + str_yyyy + `, ` + str_mm + `, ` + str_dd + `), ` + RG2.Vals2 + `],`

		   //r = RetMessage{r.Date_s,r.Date_s}
	   }
	   tconf["graph2_data"] += `]`
	   tconf["graph2_data2"] += `]` */
	var err error
	tmpl := template.New("par_partner.html")
	//var err error
	if tmpl, err = tmpl.ParseFiles("dashboard/operations/par_partner.html"); err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, tconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func DataList(sql string, org string) [][]string {

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
			if i == 0 {
				a = append(a, config.Img_to_render_custom1(PrintValue(vals[i].(*interface{})), 80, org))
			} else {
				a = append(a, PrintValue(vals[i].(*interface{})))
			}
		}
		m = append(m, a)
	}

	if rows.Err() != nil {
		fmt.Println(rows.Err)
	}
	return m

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
