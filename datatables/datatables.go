package datatables

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	config "../config"
)

type Jsondata struct {
	Draw            int        `json:"draw"`
	RecordsTotal    int        `json:"recordsTotal"`
	RecordsFiltered int        `json:"recordsFiltered"`
	Data            [][]string `json:"data"`
}

type Records struct {
	Page_current int
	Page_more    int
	Page_total   int
	Record_more  int
	Record_total int
}

func DatatableretArray(w http.ResponseWriter, conf map[string]string) {
	getTotal, err1, _, _ := config.Ap_sql(conf["sql_total"], 1)
	var err error
	if err1 != nil {
		panic(err1.Error())
	}

	var Recs Records
	for getTotal.Next() {
		err = getTotal.Scan(&Recs.Page_current, &Recs.Page_more, &Recs.Page_total, &Recs.Record_more, &Recs.Record_total)
		if err != nil {
			panic(err.Error())
		}

	}

	draw, _ := strconv.Atoi(conf["dr"])

	var objdata = [][]string{}

	if Recs.Record_total > 0 {
		objdata = Datatable(conf["sql_data"], conf)
	}

	rawdatatable := &Jsondata{draw, Recs.Record_total, Recs.Record_total, objdata}
	js, err := json.Marshal(rawdatatable)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	objdata = nil
	rawdatatable = nil
	js = nil

}

func Datatable(sql string, conf map[string]string) [][]string {

	var m [][]string
	m1 := make(map[string]int)
	var Actioncol bool
	Actioncol = true
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
	//geting columnheader key
	var hkey_action_param int
	for key, value := range cols {
		if strings.ToUpper(value) == strings.ToUpper(conf["ActionCol_param"]) {
			hkey_action_param = key
		}
	}
	//end geting column key
	draw1, _ := strconv.Atoi(conf["dr"])
	m1["draw"] = draw1
	col_img_num, _ := strconv.Atoi(conf["col_img_num"])
	//column hide initialazation
	col_hidden1, _ := strconv.Atoi(conf["col_hidden1"])
	if col_hidden1 == 0 && conf["col_hidden1"] == `0` {
		col_hidden1 = 0
	} else {
		col_hidden1 = 1000
	}
	//end column hide initialazation
	img_id_col, _ := strconv.Atoi(conf["img_id_col"])
	void, _ := strconv.Atoi(conf["void"])
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
		var fontcolor string
		var posted bool
		var Cancelled bool
		var Empty bool
		var ForPosting bool
		fontcolor = ""
		fmt.Println(fontcolor)
		fmt.Println("void value here")
		//fmt.Println(PrintValue(vals[void].(*interface{})))
		/*if( strings.Contains( PrintValue(vals[3].(*interface{})) , "Posted") ){
			fontcolor ="style='color:red;'"
		}*/
		if len(cols) > 3 {
			if strings.Contains(PrintValue(vals[3].(*interface{})), "Posted") {
				posted = true
			}
		}
		if strings.Contains(PrintValue(vals[2].(*interface{})), "Posted") {
			posted = true
		}

		if strings.Contains(PrintValue(vals[1].(*interface{})), "Posted") {
			posted = true
		}
		if len(cols) > 3 {
			if strings.Contains(PrintValue(vals[3].(*interface{})), "Cancelled") {
				Cancelled = true
			}
		}
		if strings.Contains(PrintValue(vals[2].(*interface{})), "Cancelled") {
			Cancelled = true
		}

		if strings.Contains(PrintValue(vals[1].(*interface{})), "Cancelled") {
			Cancelled = true
		}
		if len(cols) > 3 {
			if strings.Contains(PrintValue(vals[3].(*interface{})), "Empty") {
				Empty = true
			}
		}
		if strings.Contains(PrintValue(vals[2].(*interface{})), "Empty") {
			Empty = true
		}

		if strings.Contains(PrintValue(vals[1].(*interface{})), "Empty") {
			Empty = true
		}

		if conf["status"] == "For Posting" {
			ForPosting = true
		}

		if conf["status"] == "Cancelled" {
			Cancelled = true
		}

		for i := 0; i < len(cols); i++ { // fetching data tables to container map array
			if Actioncol && i == 0 {
				var EditAction, AllocationAction, DeleteAction, AddChildAction, AddroleAction, AddDetailAction, AssignmentAction, MembersAction, DuplicateAction, PostAction, ApplicationAction, CustomerAction, LocationAction, ItemAction, SupplierAction, AddChildInheritAction, VoidAction, ViewAction, ContributionAction, QuickEditAction string

				if conf["ActionCol_view"] == "true" {
					if conf["ActionCol_view_is_modal"] == "true" {
						ViewAction = `<a data-toggle="modal" href="#basic" onclick="Modalfunc('` + conf["ActionCol_view_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this);">View</a>`
					} else {
						//EditAction =`<a  name="`+conf["ActionCol_param"]+`' href ="#`+ conf["ActionCol_edit_url"] +`?id=`+PrintValue(vals[i].(*interface{})) +`" onclick="Nav2(this,`+PrintValue(vals[hkey_action_param].(*interface{})) +`)" >Edit</a>`
						ViewAction = `<a href ="#` + conf["ActionCol_view_url"] + PrintValue(vals[i].(*interface{})) + `" onclick="func_view('` + conf["ActionCol_view_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this,event)" >View</a>`
					}
					if posted {
						ViewAction = `<a href ="#` + conf["ActionCol_view_url"] + PrintValue(vals[i].(*interface{})) + `" onclick="func_view('` + conf["ActionCol_view_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this,event)" >View</a>`
					}
					if Cancelled {
						ViewAction = `<a href ="#` + conf["ActionCol_view_url"] + PrintValue(vals[i].(*interface{})) + `" onclick="func_view('` + conf["ActionCol_view_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this,event)" >View</a>`
					}
				}

				if conf["ActionCol_edit"] == "true" {
					if conf["ActionCol_edit_is_modal"] == "true" {
						EditAction = `<a data-toggle="modal" href="#basic" onclick="Modalfunc('` + conf["ActionCol_edit_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this);">Edit</a>`
					} else {
						//EditAction =`<a  name="`+conf["ActionCol_param"]+`' href ="#`+ conf["ActionCol_edit_url"] +`?id=`+PrintValue(vals[i].(*interface{})) +`" onclick="Nav2(this,`+PrintValue(vals[hkey_action_param].(*interface{})) +`)" >Edit</a>`
						EditAction = `<a href ="#` + conf["ActionCol_edit_url"] + PrintValue(vals[i].(*interface{})) + `" onclick="func_edit('` + conf["ActionCol_edit_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this,event)" >Edit</a>`
					}
					if posted {
						EditAction = ``
					}
					if Cancelled {
						EditAction = ``
					}
				}

				if conf["ActionCol_duplicate"] == "true" {
					if conf["ActionCol_duplicate_is_modal"] == "true" {
						DuplicateAction = `<a data-toggle="modal" href="#basic" onclick="Modalfunc('` + conf["ActionCol_duplicate_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this);">Replicate</a>`
					} else {
						//DuplicateAction =`<a  name="`+conf["ActionCol_param"]+`' href ="#`+ conf["ActionCol_Duplicate_url"] +`?id=`+PrintValue(vals[i].(*interface{})) +`" onclick="Nav2(this,`+PrintValue(vals[hkey_action_param].(*interface{})) +`)" >Duplicate</a>`
						DuplicateAction = `<a href ="#` + conf["ActionCol_duplicate_url"] + PrintValue(vals[i].(*interface{})) + `" onclick="func_duplicate('` + conf["ActionCol_duplicate_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this,event)" >Replicate</a>`
					}
					if posted {
						DuplicateAction = ``
					}
					if Cancelled {
						DuplicateAction = ``
					}
				}
				if conf["ActionCol_delete"] == "true" {
					DeleteAction = `<a data-toggle="modal" href="#modalSession" onclick="Modalfuncdelete('` + conf["ActionCol_delete_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',1,this)" >Delete</a>`
					if posted {
						DeleteAction = ``
					}
					if Cancelled {
						DeleteAction = ``
					}
				}

				if conf["ActionCol_void"] == "true" {
					VoidAction = `<a data-toggle="modal" href="#modalSession" onclick="Modalfuncvoid('` + conf["ActionCol_void_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',1,this)" >Void</a>`
					if posted && PrintValue(vals[void].(*interface{})) == "0" {
						VoidAction = `<a data-toggle="modal" href="#modalSession" onclick="Modalfuncvoid('` + conf["ActionCol_void_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',1,this)" >Void</a>`
					}
					if Cancelled {
						VoidAction = ``
					}

					if ForPosting {
						VoidAction = ``
					}

					if PrintValue(vals[void].(*interface{})) == "1" {

						VoidAction = ``
					}

				}
				if conf["ActionCol_post"] == "true" {
					PostAction = `<a data-toggle="modal" href="#modalSession" onclick="Modalfuncpost('` + conf["ActionCol_post_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',1,this)" >Post</a>`
					if posted {
						PostAction = ``
					}
					if Cancelled {
						PostAction = ``
					}
				}

				if conf["ActionCol_add_child_inherit"] == "true" {
					if conf["ActionCol_add_child_inherit_is_modal"] == "true" {
						AddChildInheritAction = `<a data-toggle="modal" href="#basic" onclick="Modalfunc('` + conf["ActionCol_add_child_inherit_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this);">Add Child</a>`
					} else {
						//DuplicateAction =`<a  name="`+conf["ActionCol_param"]+`' href ="#`+ conf["ActionCol_Duplicate_url"] +`?id=`+PrintValue(vals[i].(*interface{})) +`" onclick="Nav2(this,`+PrintValue(vals[hkey_action_param].(*interface{})) +`)" >Duplicate</a>`
						AddChildInheritAction = `<a href ="#` + conf["ActionCol_add_child_inherit_url"] + PrintValue(vals[i].(*interface{})) + `" onclick="func_add_child_inherit('` + conf["ActionCol_add_child_inherit_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this,event)" >Add Child</a>`
					}
					if posted {
						AddChildInheritAction = ``
					}
					if Cancelled {
						AddChildInheritAction = ``
					}
				}

				if conf["ActionCol_add_child"] == "true" {
					//AddChildAction ="<a style='color: "+fontcolor+";' href ='#"+PrintValue(vals[i].(*interface{})) +"'>Add Child</a>"
					AddChildAction = `<a href ="#` + conf["ActionCol_edit_url"] + PrintValue(vals[i].(*interface{})) + `" onclick="func_add_child('` + conf["ActionCol_add_child_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this)" >Add Child</a>`
				}
				if conf["ActionCol_roles"] == "true" {
					AddroleAction = "<a  name='" + conf["ActionCol_param"] + "' href ='#" + conf["ActionCol_roles_url"] + "?id=" + PrintValue(vals[i].(*interface{})) + "' onclick='Nav2(this," + PrintValue(vals[hkey_action_param].(*interface{})) + ")' >Roles</a>"
				}
				if conf["ActionCol_detail"] == "true" {
					AddDetailAction = "<a  name='" + conf["ActionCol_param"] + "' href ='#" + conf["ActionCol_detail_url"] + "?id=" + PrintValue(vals[i].(*interface{})) + "' onclick='Nav2(event," + PrintValue(vals[hkey_action_param].(*interface{})) + ")' >Details</a>"
				}
				if conf["ActionCol_customer"] == "true" {
					CustomerAction = "<a  name='" + conf["ActionCol_param"] + "' href ='#" + conf["ActionCol_customer_url"] + "?id=" + PrintValue(vals[i].(*interface{})) + "' onclick='Nav2(event," + PrintValue(vals[hkey_action_param].(*interface{})) + ")' >Customers</a>"
				}
				if conf["ActionCol_location"] == "true" {
					LocationAction = "<a  name='" + conf["ActionCol_param"] + "' href ='#" + conf["ActionCol_location_url"] + "?id=" + PrintValue(vals[i].(*interface{})) + "' onclick='Nav2(event," + PrintValue(vals[hkey_action_param].(*interface{})) + ")' >Locations</a>"
				}
				if conf["ActionCol_item"] == "true" {
					LocationAction = "<a  name='" + conf["ActionCol_param"] + "' href ='#" + conf["ActionCol_item_url"] + "?id=" + PrintValue(vals[i].(*interface{})) + "' onclick='Nav2(event," + PrintValue(vals[hkey_action_param].(*interface{})) + "," + `"items"` + ")' >Items</a>"
				}
				if conf["ActionCol_supplier"] == "true" {
					LocationAction = "<a  name='" + conf["ActionCol_param"] + "' href ='#" + conf["ActionCol_supplier_url"] + "?id=" + PrintValue(vals[i].(*interface{})) + "' onclick='Nav2(event," + PrintValue(vals[hkey_action_param].(*interface{})) + "," + `"suppliers"` + ")' >Suppliers</a>"
				}
				if conf["ActionCol_contribution"] == "true" {
					ContributionAction = "<a  name='" + conf["ActionCol_param"] + "' href ='#" + conf["ActionCol_contribution_url"] + "?id=" + PrintValue(vals[i].(*interface{})) + "' onclick='Nav2(event," + PrintValue(vals[hkey_action_param].(*interface{})) + "," + `"contribution"` + ")' >Contribution</a>"
				}
				if conf["ActionCol_allocation"] == "true" {
					if conf["ActionCol_allocation_is_modal"] == "true" {
						AllocationAction = `<a data-toggle="modal" href="#basic" onclick="Modalfunc('` + conf["ActionCol_allocation_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this);">Allocations</a>`
					} else {
						//DuplicateAction =`<a  name="`+conf["ActionCol_param"]+`' href ="#`+ conf["ActionCol_Duplicate_url"] +`?id=`+PrintValue(vals[i].(*interface{})) +`" onclick="Nav2(this,`+PrintValue(vals[hkey_action_param].(*interface{})) +`)" >Duplicate</a>`
						AllocationAction = `<a href ="#` + conf["ActionCol_allocation_url"] + PrintValue(vals[i].(*interface{})) + `" onclick="func_allocation('` + conf["ActionCol_allocation_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this,event)" >Allocations</a>`
					}

				}
				if conf["ActionCol_application"] == "true" {
					ApplicationAction = ``
					if posted {
						ApplicationAction = "<a  name='" + conf["ActionCol_param"] + "' href ='#" + conf["ActionCol_application_url"] + "?id=" + PrintValue(vals[i].(*interface{})) + "' onclick='Nav2(event," + PrintValue(vals[hkey_action_param].(*interface{})) + "," + `"application"` + ")' >Application</a>"
					}
					if Empty {
						ApplicationAction = ``
					}
				}
				if conf["ActionCol_assignment"] == "true" {
					AssignmentAction = `<a href ="#` + conf["ActionCol_assignment_url"] + PrintValue(vals[i].(*interface{})) + `" onclick="func_assignment('` + conf["ActionCol_assignment_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this)" >Assignment</a>`
				}
				if conf["ActionCol_members"] == "true" { //on hris->assignment module button on members on assignment
					MembersAction = `<a href ="#` + conf["ActionCol_members_url"] + PrintValue(vals[i].(*interface{})) + `" onclick="func_assignment('` + conf["ActionCol_members_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this)" >Members</a>`
				}

				if conf["ActionCol_quick_edit"] == "true" {
					QuickEditAction = `<a href="javascript:showQuickEditModal('` + PrintValue(vals[hkey_action_param].(*interface{})) + `', '` + PrintValue(vals[1].(*interface{})) + `');"><i class="fas fa-user-edit"></i>QuickEdit</a>`
				}

				a = append(a, EditAction+"  "+DeleteAction+" "+AllocationAction+" "+AddChildAction+" "+AddroleAction+" "+AddDetailAction+" "+AssignmentAction+"  "+MembersAction+" "+DuplicateAction+" "+PostAction+" "+AddChildInheritAction+" "+ApplicationAction+" "+CustomerAction+" "+LocationAction+" "+ItemAction+" "+SupplierAction+" "+VoidAction+" "+ViewAction+" "+ContributionAction+" "+QuickEditAction)

				if conf["ActionCol_allocation"] == "true" {
					if conf["ActionCol_allocation_is_modal"] == "true" {
						AllocationAction = `<a data-toggle="modal" href="#basic" onclick="Modalfunc('` + conf["ActionCol_allocation_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this);">Allocations</a>`
					} else {
						//DuplicateAction =`<a  name="`+conf["ActionCol_param"]+`' href ="#`+ conf["ActionCol_Duplicate_url"] +`?id=`+PrintValue(vals[i].(*interface{})) +`" onclick="Nav2(this,`+PrintValue(vals[hkey_action_param].(*interface{})) +`)" >Duplicate</a>`
						AllocationAction = `<a href ="#` + conf["ActionCol_allocation_url"] + PrintValue(vals[i].(*interface{})) + `" onclick="func_allocation('` + conf["ActionCol_allocation_url"] + PrintValue(vals[hkey_action_param].(*interface{})) + `',this,event)" >Allocations</a>`
					}

				}
			}
			//a  = append(a,"<div "+fontcolor+">" + PrintValue(vals[i].(*interface{})) + "</div>"  )
			if conf["col_img"] == "true" && i == col_img_num {
				a = append(a, `<img src="`+config.Img_to_render(`img_uploads/size_80/`+conf["org"]+`/`+PrintValue(vals[0].(*interface{}))+`.jpg`, PrintValue(vals[0].(*interface{})), 80, conf["org"])+`" width="80" >`)
			}
			if conf["name_with_img"] == "true" && i == col_img_num {
				a = append(a, config.Img_to_render_custom1(PrintValue(vals[img_id_col].(*interface{})), 80, conf["org"])+PrintValue(vals[i].(*interface{})))
			} else {
				if i != col_hidden1 { //hide column
					a = append(a, PrintValue(vals[i].(*interface{})))
				}
			}
		}
		m = append(m, a)
		//mars at work here
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
		a = fmt.Sprint("NULL")
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
		//a = fmt.Sprint(v.Format("2006-01-02 15:04:05.999"))
		a = fmt.Sprint(v.Format("01/02/2006"))

	default:
		//		fmt.Print(v)
		a = fmt.Sprint(v)
	}
	return a
}

func DataList(sql string) [][]string {

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
			a = append(a, PrintValue(vals[i].(*interface{})))
		}
		m = append(m, a)
		//mars at work here
	}

	if rows.Err() != nil {
		fmt.Println(rows.Err)
	}
	return m

}

func Data_row(sql string) []string {

	var m []string
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
			a = append(a, PrintValue(vals[i].(*interface{})))

		}
		m = a
		//mars at work here
	}

	if rows.Err() != nil {
		fmt.Println(rows.Err)
	}
	return m

}
