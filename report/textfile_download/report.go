package report

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	//"strconv"
	//"encoding/json"
	config "../../config"
	datatables "../../datatables"
	login "../../login"
	//"strings"
)

type Profile struct {
	Message string // error status output or validated data
	Status  bool   //pass or fail
	ID      string
}

func TextFileDownload(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	//_, session_user_id := login.Get_account_info(r)
	//session_user_id_str := strconv.Itoa(session_user_id)
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)

	type Htdata struct {
		Ls_partner [][]string
	}

	partner := datatables.DataList(`exec partners_get 4,` + str_OrgID + `,''`)

	var_htdata := Htdata{partner}
	tmpl := template.New("rawreports_nav.html")
	var err error
	if tmpl, err = tmpl.ParseFiles("report/textfile_download/rawreports_nav.html"); err != nil {
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
)

func Download(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
	login.Session_renew(w, r)
	username, _ := login.Get_account_info(r)
	transdate := r.URL.Query().Get("transdate")
	partner := r.URL.Query().Get("partner")
	partner_id = r.URL.Query().Get("partner")
	partner_name := r.URL.Query().Get("partner_name")

	trandate = r.URL.Query().Get("transdate")

	t := time.Now()
	DatetimeFileAddname := t.Format("20060102150405")
	//a := fmt.Sprint(lines[1])
	downloadfilename := partner_name + "_" + DatetimeFileAddname

	//var i string
	//fmt.Scan(&i)
	//os.Exit(0)
	tconf := make(map[string]string)
	tconf["filename"] = downloadfilename
	sql_sp = `Textfile_download 1,'` + transdate + `',` + partner + `,'` + username + `'`

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

	newFile, err := os.Create("txtfile/" + downloadfilename + ".txt")
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
	ForceDownload(w, r, "txtfile/"+downloadfilename+".txt")
	fmt.Println("eto" + downloadfilename)

}

/* func Datagrid(fpath string) {
	//var m  [][]string

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

	newFile, err := os.Create("txtfile/" + fpath + ".txt")
	if err != nil {
		fmt.Println(err)
		newFile.Close()
		return
	}
	var header string
	header = lpad(partner_id, "0", 4) + "," + trandate + "\r\n"
	header = strings.TrimSpace(header)
	_, err = fmt.Fprintln(newFile, header+"\r")
	if err != nil {
		fmt.Println("the  error was here %v", err)
		return
	}

	var totalamount float64
	count = 0
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
			lineString += strings.TrimSpace(PrintValue(vals[i].(*interface{})))
			if i != (len(cols) - 1) {
				//_, err = file.WriteString(strings.TrimSpace("") + "," + "\t")
				lineString += ","
			}
			if i == 3 {
				num, err := strconv.ParseFloat(PrintValue(vals[i].(*interface{})), 64)
				if err == nil {
					totalamount += num
				}

			}

		}
		fmt.Println("line string :%v", lineString)
		_, err = fmt.Fprintln(newFile, lineString+"\r")
		if err != nil {
			fmt.Println("the  error was here %v", err)
			return
		}

		//total = total + vals[i][3]
		count = count + 1
		fmt.Println("counter %s ", count)
		//m = append( m,a)
	}
	var footer string
	footer = strconv.Itoa(count) + "," + fmt.Sprintf("%.2f", totalamount)
	footer = strings.TrimSpace(footer)

	//	totalLine := count ++
	_, err = fmt.Fprintln(newFile, footer)
	if err != nil {
		fmt.Println("the  error was here %v", err)
		return
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
	fmt.Println(totalamount)

} */

func Ap_sql(sq string, sq_type int) (*sql.Rows, error, bool, sql.Result) {
	//dsn := "server=marlon-it\\sqlexpress;user id=sa;password=masterkey;database=ap_web_db"
	//dsn := "server=marlon-it\\sqlexpress2008;user id=sa;password=masterkey;database=tap_11aug15;encrypt=disable"
	dsn := strdns
	//dsn := "server=marlon-it\\sqlexpress2008;user id=sa;password=masterkey;database=tap_11aug15;encrypt=disable"
	//dsn := "server=gleceper-lenovo\\sqlexpress;user id=sa;password=masterkey;database=tap;encrypt=disable"

	//dsn := "server=SQL5009.Smarterasp.net;user id=DB_9D5B76_muse_admin;password=nolram321;database=DB_9D5B76_muse;encrypt=disable"
	//fmt.Println(dsn)
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
		a = fmt.Sprint(v.Format("2006-01-02 15:04:05.999"))

	default:
		//		fmt.Print(v)
		a = fmt.Sprint(v)
	}
	return a
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

func ForceDownload(w http.ResponseWriter, r *http.Request, file string) {

	downloadBytes, err := ioutil.ReadFile(file)

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
