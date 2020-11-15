package supplier_receipt

import (
	"fmt"
	"net/http"

	//"text/template"
	config "../../config"
	//"strconv"
	"encoding/json"
	//datatables "../../datatables"
	login "../../login"
	//helper "../../helper"
	//"strings"
)

func Ajax_getItem(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	getq := r.URL.Query().Get("q")
	getid := r.URL.Query().Get("id")
	//isbaseonassignid :=r.URL.Query().Get("isbaseonassignid")
	//assignmentID := r.URL.Query().Get("assignmentID")
	//fmt.Println(`is base sa assignID`,isbaseonassignid)
	//fmt.Println(` assignID`,assignmentID)

	getq = getq + getid
	type Inner_item struct {
		Id       string `json:"id"`
		Fullname string `json:"full_name"`
		//Item_image string `json:"avatar"`
		//Price string `json:"price"`
	}
	type JsondataEmpty struct {
		Total_count        int          `json:"total_count"`
		Incomplete_results bool         `json:"incomplete_results"`
		Items              []Inner_item `json:"items"`
	}
	var count int
	/*var strtype string
	if getid !=""{
		strtype = "1"
	}else{
		strtype = "2"
	}	*/
	//`exec lbr_memass_search 1, 395 , '' ,'',''`
	var sqlstr string
	//Org_id :=login.Get_session_org_id(r)
	//str_OrgID :=strconv.Itoa(Org_id)
	sqlstr = "exec Members_search 1,'" + getq + "'"
	fmt.Println(sqlstr)
	rows, err, _, _ := config.Ap_sql(sqlstr, 1)
	if err != nil {
		fmt.Println("db error:", err)
	}
	InItem := []Inner_item{}
	var InItemRow Inner_item
	for rows.Next() {
		var r Inner_item
		err = rows.Scan(&r.Id, &r.Fullname)
		//err = rows.Scan(&r.Id,&r.Fullname,&r.Item_image,&r.Price)
		if err != nil {
			panic(err.Error())
		}
		count = count + 1
		post := Inner_item{
			r.Id,
			r.Fullname,
			//config.Img_to_render(`img_uploads/size_80/`+str_OrgID+`/`+r.Item_image+`.jpg`,r.Item_image,80,str_OrgID),
			//r.Price,
		}
		InItem = append(InItem, post)
		//InItemRow = Inner_item{r.Id,r.Item_image, config.Img_to_render(`img_uploads/size_80/`+str_OrgID+`/` +r.Item_image+`.jpg`,r.Item_image,80,str_OrgID),r.Price}
		InItemRow = Inner_item{r.Id, r.Fullname}
	}

	if getid != "" {

		//jsondata := JsondataEmpty{count,false,InItem}
		js, err := json.Marshal(InItemRow)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	} else {

		jsondata := JsondataEmpty{count, false, InItem}
		js, err := json.Marshal(jsondata)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}

}

func SearchNoStatus(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	getq := r.URL.Query().Get("q")
	getid := r.URL.Query().Get("id")
	//isbaseonassignid :=r.URL.Query().Get("isbaseonassignid")
	//assignmentID := r.URL.Query().Get("assignmentID")
	//fmt.Println(`is base sa assignID`,isbaseonassignid)
	//fmt.Println(` assignID`,assignmentID)

	getq = getq + getid
	type Inner_item struct {
		Id       string `json:"id"`
		Fullname string `json:"full_name"`
		//Item_image string `json:"avatar"`
		//Price string `json:"price"`
	}
	type JsondataEmpty struct {
		Total_count        int          `json:"total_count"`
		Incomplete_results bool         `json:"incomplete_results"`
		Items              []Inner_item `json:"items"`
	}
	var count int
	/*var strtype string
	if getid !=""{
		strtype = "1"
	}else{
		strtype = "2"
	}	*/
	//`exec lbr_memass_search 1, 395 , '' ,'',''`
	var sqlstr string
	//Org_id :=login.Get_session_org_id(r)
	//str_OrgID :=strconv.Itoa(Org_id)
	sqlstr = "exec Members_search 2,'" + getq + "'"
	fmt.Println(sqlstr)
	rows, err, _, _ := config.Ap_sql(sqlstr, 1)
	if err != nil {
		fmt.Println("db error:", err)
	}
	InItem := []Inner_item{}
	var InItemRow Inner_item
	for rows.Next() {
		var r Inner_item
		err = rows.Scan(&r.Id, &r.Fullname)
		//err = rows.Scan(&r.Id,&r.Fullname,&r.Item_image,&r.Price)
		if err != nil {
			panic(err.Error())
		}
		count = count + 1
		post := Inner_item{
			r.Id,
			r.Fullname,
			//config.Img_to_render(`img_uploads/size_80/`+str_OrgID+`/`+r.Item_image+`.jpg`,r.Item_image,80,str_OrgID),
			//r.Price,
		}
		InItem = append(InItem, post)
		//InItemRow = Inner_item{r.Id,r.Item_image, config.Img_to_render(`img_uploads/size_80/`+str_OrgID+`/` +r.Item_image+`.jpg`,r.Item_image,80,str_OrgID),r.Price}
		InItemRow = Inner_item{r.Id, r.Fullname}
	}

	if getid != "" {

		//jsondata := JsondataEmpty{count,false,InItem}
		js, err := json.Marshal(InItemRow)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	} else {

		jsondata := JsondataEmpty{count, false, InItem}
		js, err := json.Marshal(jsondata)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}

}
