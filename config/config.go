package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"

	//"go-mssqldb-master"
	_ "../github.com/go-mssqldb-master"

	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"
	//"path/filepath"
)

var (
	file           *os.File
	part           []byte
	prefix         bool
	dsn_con        string
	enable_log     string
	Email_id       string
	Email_password string
	Email_to       string
	Email_bcc      string
	Email_cc       string
)

func init() {
	lines, err := readLines("config.dat")
	if err != nil {
		fmt.Println("error: empty config.")
	}
	dsn_con = lines[0]
	enable_log = strings.Replace(lines[1], "enable_log=", "", -1)
	Email_id = strings.Replace(lines[2], "email_id=", "", -1)
	Email_password = strings.Replace(lines[3], "email_password=", "", -1)
	Email_to = strings.Replace(lines[4], "email_to=", "", -1)
	Email_bcc = strings.Replace(lines[5], "email_bcc=", "", -1)
	Email_cc = strings.Replace(lines[6], "email_cc=", "", -1)

	fmt.Println(dsn_con)
	fmt.Println(enable_log)

}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

var SessionMAxTime = 1800
var Encrypt_key = []byte(`Encrypt_key_mu5e5ecreT2015_@!234`)

func Ap_sql(sq string, sq_type int) (*sql.Rows, error, bool, sql.Result) {
	//dsn := "MARLON-IT\\SQLEXPRESS2008;user id=sa;password=masterkey;database=tap_Aug212015;encrypt=disable"
	//dsn := "server=gleceper-lenovo\\sqlexpress;user id=sa;password=masterkey;database=tap;encrypt=disable"
	dsn := dsn_con
	//dsn := "server=SQL5009.Smarterasp.net;user id=DB_9D5B76_muse_admin;password=nolram321;database=DB_9D5B76_muse;encrypt=disable"
	//dsn := "server=SQL5014.Smarterasp.net;user id=DB_9D60DB_tap_admin;password=masterkey;database=DB_9D60DB_tap;encrypt=disable"
	//fmt.Println(dsn)
	db, err := sql.Open("mssql", dsn)
	//db.SetConnMaxLifetime(300)
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		go LogHandler(fmt.Sprint("Cannot connect: ", err.Error()))
		defer db.Close()
		return nil, err, false, nil
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		defer db.Close()
		return nil, err, false, nil
	}

	stmt, err_pr := db.Prepare(sq)
	if err_pr != nil {
		fmt.Println(`error db prepare:`, err_pr)
		//log.Fatal(err_pr)
	}

	if sq_type == 1 {
		rows, err := stmt.Query()
		fmt.Println("option 1 err: ", err)
		defer db.Close()
		return rows, err, false, nil
	} else if sq_type == 2 {
		rows, err := stmt.Query()
		//fmt.Println("option 2 err: ", err)
		defer db.Close()
		return rows, err, false, nil
	} else {
		res, err := stmt.Exec()
		//fmt.Println("option 3 err: ", err)
		defer db.Close()
		return nil, err, true, res
	}

	fmt.Println(err)
	go LogHandler(fmt.Sprint(err))
	defer db.Close()
	//fmt.Println("%v", err)
	return nil, err, true, nil

}

func Encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func Decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

func PrintValue(pval *interface{}) {
	switch v := (*pval).(type) {
	case nil:
		fmt.Print("NULL")
	case bool:
		if v {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
	case []byte:
		fmt.Print(string(v))
	case time.Time:
		fmt.Print(v.Format("2006-01-02 15:04:05.999"))
	default:
		fmt.Print(v)
	}
}

func Interfacetosting(pval interface{}) string {
	var a string
	switch v := (pval).(type) {
	case nil:
		a = fmt.Sprint("")
	case bool:
		if v {
			a = fmt.Sprint("1")
		} else {
			a = fmt.Sprint("0")
		}
	case []byte:
		a = fmt.Sprint(string(v))
	case time.Time:
		a = fmt.Sprint(v.Format("Jan _2 2006 03:04:05PM"))

	default:
		a = fmt.Sprint(v)
	}
	return a

}

func Img_to_render(img_dir string, item_id string, size int, org string) string {
	ret_data := ``
	if size == 80 {
		ret_data = `/src/img_uploads/default/default_80.jpg`
	}
	if size == 171 {
		ret_data = `/src/img_uploads/default/default_171.jpg`
	}

	file_info, err := os.Stat(img_dir)
	//fmt.Println()
	if err == nil {
		if size == 80 {
			ret_data = `/src/img_uploads/size_80/` + org + `/` + item_id + `.jpg?` + fmt.Sprint(file_info.ModTime())
		}
		if size == 171 {
			ret_data = `/src/img_uploads/size_171/` + org + `/` + item_id + `.jpg?` + fmt.Sprint(file_info.ModTime())
		}
	}
	return ret_data
}
func Img_to_render_custom1(item_id string, size int, org string) string {
	ret_data := ``
	fmt.Println(`itemid:`, item_id)
	file_info, err := os.Stat(`img_uploads/size_80/` + org + `/` + item_id + `.jpg`)
	//fmt.Println()
	if err == nil {
		if size == 80 {
			ret_data = `<img src="/src/img_uploads/size_80/` + org + `/` + item_id + `.jpg?` + fmt.Sprint(file_info.ModTime()) + `" width="80" >`

		}
		if size == 171 {
			ret_data = `/src/img_uploads/size_171/` + org + `/` + item_id + `.jpg?` + fmt.Sprint(file_info.ModTime())
		}
	}
	return ret_data

}

func LogHandler(writestring string) {
	if enable_log == `true` {
		t := time.Now()
		DatetimeFileAddname := t.Format("20060102")
		fpath := `logs\` + `muse_logs_` + DatetimeFileAddname + `.dat`
		contents, _ := ioutil.ReadFile(fpath)

		var file *os.File
		var err error
		if file, err = os.Create(fpath); err != nil { ///data file name fpath
			fmt.Println(err)
		}

		//defer file.Close()
		//_,err = file.WriteString(strings.TrimSpace("")+`starting log ` + DatetimeFileAddname )
		_, err = file.Write(contents)
		_, err = file.WriteString(strings.TrimSpace("") + writestring)
		_, err = file.WriteString(strings.TrimSpace("") + "\n")
	}
}

func Trimquoute(str string) string {
	StrFirst := strings.Trim(str, `"`)
	StrLast := strings.Replace(StrFirst, `"`, ``, 3)
	return StrLast
}

var FuncMap = template.FuncMap{
	"date_formater": func(a time.Time) string {
		return a.Format(" Jan _2 2006")
	},

	"footer_ads": func() string {
		ads := `<font size="3"><b>Powered by Memgertech Inc.</b> <i>Inventory like a <b>boss</b>... Try us for free.</i></font>`
		return (ads)

	},
	"get_datetime": func() string {

		t := time.Now().Local()
		return t.Format("_2 Jan 2006 03:04:05PM")
	},
	"func_intertaceToString": func(pval interface{}) string {
		var a string
		switch v := (pval).(type) {
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
			a = fmt.Sprint(v.Format("Jan _2 2006 03:04:05PM"))

		default:
			//		fmt.Print(v)
			a = fmt.Sprint(v)
		}
		return a
	},

	"Return_tr_tble": func(filtername string, arrdata [][]string) string {
		htdata := ``
		for _, val := range arrdata {
			if val[0] == filtername {
				htdata = htdata + `<tr>`
				for _, nesteddata := range val {
					htdata = htdata + `<td>`
					htdata = htdata + nesteddata
					htdata = htdata + `</td>`
				}
				htdata = htdata + `</tr>`
			}

		}

		return htdata

	},

	"col_attrib": func(arr_key int, arr_header []string) string {
		htdata := ``
		if strings.Contains(arr_header[arr_key], "hideme") {
			htdata = `style="display: none;"`
		}

		if strings.Contains(arr_header[arr_key], "isnum") {
			htdata = `align="right"`
		}

		return htdata

	},

	"col_trim_header": func(str string) string {
		htdata := str
		if strings.Contains(str, "_isnum") {
			htdata = strings.Trim(str, "_isnum")
		}
		return htdata
	},

	"combine_string": func(s1 ...string) string {
		ret_str := ``
		for _, val := range s1 {
			ret_str += val
		}
		return ret_str
	},

	"Img_render": func(item_id string, size int) string {
		ret_data := ``
		file_info, err := os.Stat(`img_uploads/size_80/` + item_id + `.jpg`)
		//fmt.Println()
		if err == nil {
			if size == 80 {
				ret_data = `<img src="/src/img_uploads/size_80/` + item_id + `.jpg?` + fmt.Sprint(file_info.ModTime()) + `" width="80" >`
			}
			if size == 171 {
				ret_data = `/src/img_uploads/size_171/` + item_id + `.jpg?` + fmt.Sprint(file_info.ModTime())
			}
		}
		return ret_data
	},

	"Img_to_render": func(img_dir string, item_id string, size int, org string) string {
		ret_data := ``
		if size == 80 {
			ret_data = `/src/img_uploads/default/default_80.jpg`
		}
		if size == 171 {
			ret_data = `/src/img_uploads/default/default_171.jpg`
		}

		file_info, err := os.Stat(img_dir)
		//fmt.Println()
		if err == nil {
			if size == 80 {
				ret_data = `/src/img_uploads/size_80/` + org + `/` + item_id + `.jpg?` + fmt.Sprint(file_info.ModTime())
			}
			if size == 171 {
				ret_data = `/src/img_uploads/size_171/` + org + `/` + item_id + `.jpg?` + fmt.Sprint(file_info.ModTime())
			}
		}
		return ret_data
	},
}
