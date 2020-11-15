package timekeeping
import (
	"net/http"
	"encoding/json"
	)

func  AjaxgetmemberbyassigmentHAndler(w http.ResponseWriter, r *http.Request) {
	//profile := Profile{"Oppsss something went wrong",false,0}
	assigmentID:=r.URL.Query().Get("assigmentID") 
	datefrom:=r.URL.Query().Get("datefrom")
	dateto:=r.URL.Query().Get("dateto") 


	js, err := json.Marshal(func_members(assigmentID,datefrom,dateto))
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
  	w.Header().Set("Content-Type", "application/json")
  	w.Write(js)			

}
/*type Members struct {
	Member_id int
	Name string
}	
func func_members() []Members {
	row ,err, _,_ := config.Ap_sql("exec lbr_memass_search 1, 8 , 'mar'",1)
	res := []Members{}
	for row.Next() {
        var r Members
        err = row.Scan(&r.Member_id, &r.Name)
        if err != nil {
                panic(err.Error())
        }
		datarow := Members{
        	r.Member_id,
	        r.Name,
	      
	    }
		res = append(res, datarow)

    }
    return res
}*/