package helper

import (
	//"database/sql"
	"fmt"
	//"crypto/aes"
	//"crypto/cipher"
	//"crypto/rand"
	//"encoding/base64"
	//"io"
	//"go-mssqldb-master"
	//_ "../github.com/go-mssqldb-master"
	datatables "../datatables"
	// "bufio"
	// "os"
	"time"
	//"io/ioutil"
	"bufio"
	"encoding/base64"
	"image/jpeg"
	"log"
	"os"
	"strings"
	"text/template"

	resize "../github.com/resize_img"
	//"path/filepath"
)

func Img_to_render(img_dir string, item_id string, org string) string {
	ret_data := `/src/img_uploads/default/default_171.png`
	file_info, err := os.Stat(img_dir)
	//fmt.Println()
	if err == nil {
		ret_data = `/src/img_uploads/size_171/` + org + `/` + item_id + `.jpg?` + fmt.Sprint(file_info.ModTime())
	}
	return ret_data
}

func Img_resizer(img_dir string, item_id string, org_id string) {
	// open "test.jpg"
	file, err := os.Open(img_dir)
	if err != nil {
		fmt.Println(err)
	}
	//create folder by organazation
	mkdir_err171 := os.Mkdir("img_uploads/size_171/"+org_id, 0777)
	fmt.Println("mkdir err 171:", mkdir_err171)
	mkdir_err80 := os.Mkdir("img_uploads/size_80/"+org_id, 0777)
	fmt.Println("mkdir err 80:", mkdir_err80)
	//end create folder by organazation
	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(171, 0, img, resize.Lanczos3)

	_, err_info1 := os.Stat("img_uploads/size_171/" + org_id + "/" + item_id + ".jpg")
	if err_info1 == nil {
		remove_err1 := os.Remove("img_uploads/size_171/" + item_id + ".jpg")
		if remove_err1 != nil {
			fmt.Println(remove_err1)
		}
	}

	out, err := os.Create("img_uploads/size_171/" + org_id + "/" + item_id + ".jpg")
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	/// new size 80px

	m2 := resize.Resize(80, 0, img, resize.Lanczos3)

	_, err_info2 := os.Stat("img_uploads/size_80/" + org_id + "/" + item_id + ".jpg")
	if err_info2 == nil {
		remove_err2 := os.Remove("img_uploads/size_80/" + org_id + "/" + item_id + ".jpg")
		if remove_err2 != nil {
			fmt.Println(remove_err2)
		}
	}

	out2, err2 := os.Create("img_uploads/size_80/" + org_id + "/" + item_id + ".jpg")
	if err2 != nil {
		log.Fatal(err2)
	}
	defer out2.Close()

	// write new image to file
	jpeg.Encode(out2, m2, nil)

}

func Base64_image_encoding(img_dir string) string {

	imgFile, err := os.Open(img_dir) // source image

	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
		return `<img alt="100%x180" data-src="holder.js/100%x180" src="data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9InllcyI/PjxzdmcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB3aWR0aD0iMTcxIiBoZWlnaHQ9IjE4MCIgdmlld0JveD0iMCAwIDE3MSAxODAiIHByZXNlcnZlQXNwZWN0UmF0aW89Im5vbmUiPjwhLS0KU291cmNlIFVSTDogaG9sZGVyLmpzLzEwMCV4MTgwCkNyZWF0ZWQgd2l0aCBIb2xkZXIuanMgMi42LjAuCkxlYXJuIG1vcmUgYXQgaHR0cDovL2hvbGRlcmpzLmNvbQooYykgMjAxMi0yMDE1IEl2YW4gTWFsb3BpbnNreSAtIGh0dHA6Ly9pbXNreS5jbwotLT48ZGVmcz48c3R5bGUgdHlwZT0idGV4dC9jc3MiPjwhW0NEQVRBWyNob2xkZXJfMTU2NDhhNmMxYTUgdGV4dCB7IGZpbGw6I0FBQUFBQTtmb250LXdlaWdodDpib2xkO2ZvbnQtZmFtaWx5OkFyaWFsLCBIZWx2ZXRpY2EsIE9wZW4gU2Fucywgc2Fucy1zZXJpZiwgbW9ub3NwYWNlO2ZvbnQtc2l6ZToxMHB0IH0gXV0+PC9zdHlsZT48L2RlZnM+PGcgaWQ9ImhvbGRlcl8xNTY0OGE2YzFhNSI+PHJlY3Qgd2lkdGg9IjE3MSIgaGVpZ2h0PSIxODAiIGZpbGw9IiNFRUVFRUUiLz48Zz48dGV4dCB4PSI2MSIgeT0iOTQuNSI+MTcxeDE4MDwvdGV4dD48L2c+PC9nPjwvc3ZnPg==" data-holder-rendered="true" style="height: 180px; display: block;">`

	} else {
		defer imgFile.Close()

		// create a new buffer base on file size
		fInfo, _ := imgFile.Stat()
		var size int64 = fInfo.Size()
		buf := make([]byte, size)

		// read file content into buffer
		fReader := bufio.NewReader(imgFile)
		fReader.Read(buf)

		// if you create a new image instead of loading from file, encode the image to buffer instead with png.Encode()

		// png.Encode(&buf, image)

		// convert the buffer bytes to base64 string - use buf.Bytes() for new image
		imgBase64Str := base64.StdEncoding.EncodeToString(buf)

		// Embed into an html without PNG file
		/*"<img alt="100%x180" data-src="holder.js/100%x180" src="data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9InllcyI/PjxzdmcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB3aWR0aD0iMTcxIiBoZWlnaHQ9IjE4MCIgdmlld0JveD0iMCAwIDE3MSAxODAiIHByZXNlcnZlQXNwZWN0UmF0aW89Im5vbmUiPjwhLS0KU291cmNlIFVSTDogaG9sZGVyLmpzLzEwMCV4MTgwCkNyZWF0ZWQgd2l0aCBIb2xkZXIuanMgMi42LjAuCkxlYXJuIG1vcmUgYXQgaHR0cDovL2hvbGRlcmpzLmNvbQooYykgMjAxMi0yMDE1IEl2YW4gTWFsb3BpbnNreSAtIGh0dHA6Ly9pbXNreS5jbwotLT48ZGVmcz48c3R5bGUgdHlwZT0idGV4dC9jc3MiPjwhW0NEQVRBWyNob2xkZXJfMTU2NDhhNmMxYTUgdGV4dCB7IGZpbGw6I0FBQUFBQTtmb250LXdlaWdodDpib2xkO2ZvbnQtZmFtaWx5OkFyaWFsLCBIZWx2ZXRpY2EsIE9wZW4gU2Fucywgc2Fucy1zZXJpZiwgbW9ub3NwYWNlO2ZvbnQtc2l6ZToxMHB0IH0gXV0+PC9zdHlsZT48L2RlZnM+PGcgaWQ9ImhvbGRlcl8xNTY0OGE2YzFhNSI+PHJlY3Qgd2lkdGg9IjE3MSIgaGVpZ2h0PSIxODAiIGZpbGw9IiNFRUVFRUUiLz48Zz48dGV4dCB4PSI2MSIgeT0iOTQuNSI+MTcxeDE4MDwvdGV4dD48L2c+PC9nPjwvc3ZnPg==" data-holder-rendered="true" style="height: 180px; display: block;">
		"*/

		img2html := `<img alt="100%x180" src="data:image/png;base64,` + imgBase64Str + `" />`
		return img2html
	}
}

var FuncMap = template.FuncMap{
	"date_formater": func(a time.Time) string {
		return a.Format(" Jan _2 2006")
	},

	"get_datetime": func() string {

		t := time.Now().Local()
		return t.Format("Mon Jan _2 2006 03:04:05PM")
	},

	"footer_ads": func() string {
		ads := `<font size="2"><b>Powered by Memgertech Inc.</b> <i>Inventory like a <b>boss</b>... Try us for free.</i></font>`
		return (ads)

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

	"col_null_header": func(str string) string {
		htdata := str
		if strings.Contains(str, "NULL") {
			htdata = strings.Trim(str, " ")
		}
		return htdata
	},

	"tpl_sql": func(sql string) [][]string {
		fmt.Println(sql)
		arr_data := datatables.DataList(sql)
		return arr_data
	},

	"tpl_sql_row": func(sql string) []string {
		fmt.Println(sql)
		arr_data := datatables.Data_row(sql)
		return arr_data
	},

	"combine_string": func(s1 ...string) string {
		ret_str := ``
		for _, val := range s1 {
			ret_str += val
		}
		return ret_str
	},

	"special_report": func(orgid string, id string, trntype string) string {
		//exec SIS_SpecialReport_Get 1,3,1
		arr_data := datatables.DataList(`exec SIS_SpecialReport_Get 1,` + orgid + `,` + trntype)
		fmt.Println(arr_data)
		htdata := ``
		for _, val := range arr_data {
			fmt.Println(val[1])
			fmt.Println(val[3])
			fmt.Println(val[2])

			arr_reportname := strings.Split(string(val[3]), "@@")
			fsqlstr := arr_reportname[0]
			for i := 1; i < len(arr_reportname); i++ {
				str_added := ""
				if strings.Contains(arr_reportname[i], "mode") {
					str_added = strings.Replace(arr_reportname[i], "mode", `1`, -1)
					fsqlstr = fsqlstr + str_added
				}
				if strings.Contains(arr_reportname[i], "org") {
					str_added = strings.Replace(arr_reportname[i], "org", orgid, -1)
					fsqlstr = fsqlstr + str_added
				}
				if strings.Contains(arr_reportname[i], "id") {
					str_added = strings.Replace(arr_reportname[i], "id", id, -1)
					fsqlstr = fsqlstr + str_added
				}
			}

			fmt.Println(fsqlstr)

			htdata += `  <li> <a href="inventory/special_rpt?sql_code=` + fsqlstr + `&rpt_url=` + val[2] + `&org_id=` + orgid + `&trnID=` + id + `" target="_blank" >` + val[1] + `</a> </li>`
		}

		return htdata
	},

	"Img_to_render": func(img_dir string, item_id string, size int, org string) string {
		ret_data := ``
		if size == 80 {
			ret_data = `/src/img_uploads/default/default_80.jpg`
		}
		if size == 171 {
			ret_data = `/src/img_uploads/default/default_171.png`
		}
		img_loc := img_dir + org + `/` + item_id + `.jpg`
		fmt.Println(`img_loc: `, img_loc)
		file_info, err := os.Stat(img_loc)
		fmt.Println(err)
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
