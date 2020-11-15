package main

import (
	"database/sql"
	_ "flag"
	"fmt"
	"io"
	"io/ioutil"
	_ "log"
	_ "net"
	"net/http"
	_ "net/url"
	"os"
	_ "regexp"
	"text/template"

	account_org_role "./admin/account_org_role"
	account_role "./admin/account_role"
	admin "./admin/role"
	config "./config"
	sessions "./github.com/gorilla/sessions"
	assignment "./hris/assignments"
	member "./hris/members"
	login "./login"

	inventory_banks "./inventory/banks"
	inventory_item "./inventory/item"
	inventory_itemclass "./inventory/item_class"
	inventory_location_item "./inventory/location_items"
	inventory_supplier "./inventory/supplier"
	inventory_supplierreceipt "./inventory/supplier_receipt"
	inventory_supplierreturn "./inventory/supplier_return"
	timekeeping "./timekeeping/transaction"

	inventory_actualcount "./inventory/actual_count"
	inventory_adjustmententry "./inventory/adjustment_entry"
	inventory_brochure "./inventory/brochure"
	inventory_customer "./inventory/customer"
	inventory_customer_adjustment "./inventory/customer_adjustment"
	inventory_customer_bill "./inventory/customer_bill"
	inventory_customer_order "./inventory/customer_order"
	inventory_customer_replacement "./inventory/customer_replacement"
	inventory_customer_return "./inventory/customer_return"
	inventory_customersales "./inventory/customer_sales"
	inventory_internal "./inventory/invent_internal"
	inventory_location "./inventory/location"
	inventory_location_order "./inventory/location_order"
	inventory_location_reciept "./inventory/location_reciept"
	inventory_location_transfer "./inventory/location_transfer"
	inventory_reports "./inventory/reports"
	inventory_supplier_order "./inventory/supplier_order"
	inventory_supplierreplacement "./inventory/supplier_replacement"

	inventory_order_per_area "./inventory/reports/order_per_area"
	inventory_order_per_area_v2 "./inventory/reports/order_per_area_v2"

	inventory_customer_payment "./inventory/customer_payment"

	inventory_internal_issuance "./inventory/internal_issuance"
	inventory_internal_order "./inventory/internal_order"
	inventory_internal_return "./inventory/internal_return"

	inventory_courier "./inventory/courier"
	inventory_item_movement "./inventory/reports/item_movement"
	inventory_product_movement "./inventory/reports/product_movement"

	inventory_custom_report "./inventory/custom_report"

	ledger "./inventory/ledger"
	item_profit "./inventory/reports/item_profit"
	non_selling_aging "./inventory/reports/non_selling_aging"

	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	dashboard "./dashboard"
	datatables "./datatables"
	standard_report "./inventory/reports"
	assignment_and_logs "./timekeeping/assignmentandlogs"
	holidays "./timekeeping/holidays"
	mfr "./timekeeping/manning_fullfillment"
	timekeeping_overtime_logs "./timekeeping/overtime_logs"
	rawreports "./timekeeping/rawreports"
	report_translog "./timekeeping/report_translog"
	time_keep "./timekeeping/time_keep"
	time_keep_uploader "./timekeeping/time_keep_uploader"
	time_keep_v2 "./timekeeping/time_keep_v2"

	//"os/exec"
	context "./github.com/gorilla/context"
	mux "./github.com/gorilla/mux"
	resize "./github.com/resize_img"

	//"log"
	"image/jpeg"

	helper "./helper"
	helper_application "./helper/application"

	///////////////////////////CWCC START////////////////////////////////////////////////
	dashboard_member "./dashboard/member"
	dashboard_operations "./dashboard/operations"
	members "./master/members"
	members_admin "./master/members_admin"
	members_mapping "./master/members_mapping"
	members_search "./master/members_search"
	partners "./master/partners"
	report "./report"
	accounting "./report/accounting"

	accounting_entry "./report/accounting_entry"
	reportTextFileDownload "./report/textfile_download"
	acknowledgement_receipt "./transaction/acknowledgement_receipt"
	adjustment_entry "./transaction/adjustment_entry"
	bulk_entry "./transaction/bulk_entry"
	loan_adjustment "./transaction/loan_adjustment"
	loan_applications "./transaction/loan_applications"
	loan_consolidation "./transaction/loan_consolidation"
	loan_payment "./transaction/loan_payment"
	loan_release "./transaction/loan_release"
	merge_member "./transaction/merge_member"
	other_disbursement "./transaction/other_disbursement"
	resign_process "./transaction/resign_process"
	savings_withdrawal "./transaction/savings_withdrawal"
	// dashboard_apply_loan "./dashboard/apply_loan"
	///////////////////////////CWCC END////////////////////////////////////////////////
)

type Page struct {
	Title string
	Body  string
}
type HtmlPage struct {
	htdata []byte
}
type access_list struct {
	Account_id int64
	Username   string
	Id         sql.NullInt64
	Name       string
	Parent     sql.NullInt64
	Url        sql.NullString
	Param      sql.NullString
	Tooltip    sql.NullString
	Leveling   int64
}

type HeaderData struct {
	HostURL string
}

type sys_account_data struct {
	Lastname   string
	Firstname  string
	Middlename string
	Username   string
	Password   string
}

var uname string
var key = []byte("topsecret2014_of_asiapro") // 32 bytes
var store = sessions.NewCookieStore([]byte("something-very-secret-cwcc"))

func main() {
	go config.LogHandler(`starting logs`)

	//mux.HandleFunc("/" , login.FrontIndexHandler)
	//mux := http.NewServeMux()
	mux := mux.NewRouter()

	//mux.HandleFunc("/" , login.FrontIndexHandler)
	mux.HandleFunc("/", login.IndexHandler)
	mux.HandleFunc("/login", login.IndexHandler)
	mux.HandleFunc("/login_attempt", login.LoginAttemptHandler)
	mux.HandleFunc("/LogOut", login.LogOutHandler)
	mux.HandleFunc("/Signup", login.SignupHandler)
	mux.HandleFunc("/signupeviaemail", login.SignupViaEmailHandler)
	mux.HandleFunc("/signup-verify-email", login.SignupverifyHandler)
	mux.HandleFunc("/signup-verify-email_validate", login.SignupverifyValidateHandler)
	mux.HandleFunc("/changepassword", login.ChangePasswordHandler)
	mux.HandleFunc("/dochangepassword", login.DoChangePasswordHandler)

	mux.HandleFunc("/Relogin", login.ReloginHandler)
	mux.HandleFunc("/SessionExpire", login.SessionExpireHandler)
	mux.HandleFunc("/mainframe", login.MainFrameHandler)
	mux.HandleFunc("/AjaxValidate", AjaxValidateHandler)
	mux.HandleFunc("/AjaxPolling", AjaxPollingeHandler)

	mux.HandleFunc("/Ajax/GetPostalcode", AjaxGetPostalCodeHandler)
	mux.HandleFunc("/Ajax/GetMembers", AjaxGetMemberHandler)

	///////////////////////////CWCC START////////////////////////////////////////////////

	//////////////////////////CWCC_MEMBERS_SEARCH_START//////////////////////////
	mux.HandleFunc("/master/members_search/search/Ajax/Ajax_getItem", members_search.Ajax_getItem)
	mux.HandleFunc("/master/members_search/search/Ajax/SearchNoStatus", members_search.SearchNoStatus)

	//////////////////////////CWCC_MEMBERS_SEARCH_END//////////////////////////

	//////////////////////////CWCC_REPORT_START//////////////////////////
	mux.HandleFunc("/report/textfile_download", reportTextFileDownload.TextFileDownload)
	mux.HandleFunc("/report/do_textfile_download", reportTextFileDownload.Download)
	mux.HandleFunc("/report/Show_savings_ledger", report.Show_savings_ledger)
	mux.HandleFunc("/report/Show_loans_ledger", report.Show_loans_ledger)
	mux.HandleFunc("/report/Show_loans_ledger_withbalance", report.Show_loans_ledger_withbalance)
	mux.HandleFunc("/report/Show_savings_ledger_modal", report.Show_savings_ledger_modal)
	mux.HandleFunc("/report/Show_loans_ledger_modal", report.Show_loans_ledger_modal)
	mux.HandleFunc("/report/Special_rpt", report.Special_rpt)
	mux.HandleFunc("/report/do_print", report.Rawreports_doprint)
	mux.HandleFunc("/report", report.Rawreports_nav)

	//////////////////////////CWCC_REPORT_END//////////////////////////
	//////////////////////////CWCC_ACCOUNTING_START//////////////////////////

	mux.HandleFunc("/report/accounting/Special_rpt", accounting.Special_rpt)
	mux.HandleFunc("/report/accounting/do_print", accounting.Rawreports_doprint)
	mux.HandleFunc("/report/accounting", accounting.Rawreports_nav)
	mux.HandleFunc("/report/accounting/bank_advise", accounting.Bank_Advise)
	mux.HandleFunc("/report/accounting/debit_memo", accounting.Debit_Memo)
	mux.HandleFunc("/report/accounting/transmittal", accounting.Transmittal)
	mux.HandleFunc("/report/accounting/transmittal2", accounting.Transmittal2)
	mux.HandleFunc("/report/accounting/dm_summary", accounting.DMSummary)
	mux.HandleFunc("/report/accounting/do_textfile_download", accounting.Download)
	mux.HandleFunc("/report/accounting/setCV", accounting.SetCV)

	//////////////////////////CWCC_ACCOUNTING_END//////////////////////////

	//////////////////////////CWCC_ACCOUNTING_START//////////////////////////

	mux.HandleFunc("/report/accounting_entry/Special_rpt", accounting_entry.Special_rpt)
	mux.HandleFunc("/report/accounting_entry/do_print", accounting_entry.Rawreports_doprint)
	mux.HandleFunc("/report/accounting_entry", accounting_entry.Rawreports_nav)

	//////////////////////////CWCC_ACCOUNTING_END//////////////////////////
	//////////////////////////CWCC_MEMBERS_START//////////////////////////
	mux.HandleFunc("/master/members", members.HListHandler)
	mux.HandleFunc("/master/members/HaddHandler", members.HAddHandler)
	mux.HandleFunc("/master/members/HEditHandler", members.HEditHandler)
	mux.HandleFunc("/master/members/HDeleteHandler", members.HDeleteHandler)
	mux.HandleFunc("/master/members/contribution", members.DListHandler)
	mux.HandleFunc("/master/members/contribution/DaddHandler", members.DaddHandler)
	mux.HandleFunc("/master/members/contribution/DEditHandler", members.DEditHandler)
	mux.HandleFunc("/master/members/contribution/DDeleteHandler", members.DDeleteHandler)
	//////////////////////////CWCC_MEMBERS_END//////////////////////////

	//////////////////////////CWCC_MEMBERS_ADMIN_START//////////////////////////
	mux.HandleFunc("/master/members_admin", members_admin.HListHandler)
	mux.HandleFunc("/master/members_admin/HEditHandler", members_admin.HEditHandler)

	mux.HandleFunc("/master/members_admin/location", members_admin.DListHandler)
	mux.HandleFunc("/master/members_admin/location/DaddHandler", members_admin.DaddHandler)
	mux.HandleFunc("/master/members_admin/location/DEditHandler", members_admin.DEditHandler)
	mux.HandleFunc("/master/members_admin/location/DDeleteHandler", members_admin.DDeleteHandler)
	mux.HandleFunc("/master/members_admin/email", members_admin.EmailHandler)
	//////////////////////////CWCC_MEMBERS_ADMIN_END//////////////////////////
	//////////////////////////CWCC_MEMBERS_MAPPING_START//////////////////////////
	mux.HandleFunc("/master/members_mapping", members_mapping.HListHandler)

	mux.HandleFunc("/master/members_mapping/HEditHandler", members_mapping.HEditHandler)

	mux.HandleFunc("/master/members_mapping/post", members_mapping.HPostHandler)

	//////////////////////////CWCC_MEMBERS_MAPPING_END//////////////////////////
	//////////////////////////CWCC_LOAN_APPLICATION_START//////////////////////////
	mux.HandleFunc("/transaction/loan_applications", loan_applications.HListHandler)
	mux.HandleFunc("/transaction/loan_applications/HaddHandler", loan_applications.HAddHandler)
	mux.HandleFunc("/transaction/loan_applications/HEditHandler", loan_applications.HEditHandler)
	mux.HandleFunc("/transaction/loan_applications/HDeleteHandler", loan_applications.HDeleteHandler)
	mux.HandleFunc("/transaction/loan_applications/HViewHandler", loan_applications.HViewHandler)
	mux.HandleFunc("/transaction/loan_applications/details", loan_applications.DListHandler)
	mux.HandleFunc("/transaction/loan_applications/details/DaddHandler", loan_applications.DAddHandler)
	mux.HandleFunc("/transaction/loan_applications/details/DEditHandler", loan_applications.DEditHandler)
	mux.HandleFunc("/transaction/loan_applications/details/DDeleteHandler", loan_applications.DDeleteHandler)
	mux.HandleFunc("/transaction/loan_applications/post", loan_applications.HPostHandler)
	mux.HandleFunc("/transaction/loan_applications/cancel", loan_applications.HCancelHandler)
	//mux.HandleFunc("/inventory/supplier_order/details/Ajax/Ajax_getRef",inventory_supplier_order.Ajax_get_reference)
	mux.HandleFunc("/transaction/loan_applications/details/email", loan_applications.EmailHandler)
	//mux.HandleFunc("/transaction/loan_applications/Show_savings_ledger",loan_applications.Show_savings_ledger)
	mux.HandleFunc("/transaction/loan_applications/details/DVoidHandler", loan_applications.DVoidHandler)
	//mux.HandleFunc("/transaction/loan_applications/details/Ajax/Ajax_getRef",escs_loan_applications.Ajax_get_reference)
	//mux.HandleFunc("/inventory/supplier_order/details/Ajax/Ajax_getRef",inventory_supplier_order.Ajax_get_reference)
	mux.HandleFunc("/transaction/loan_applications/details/email", loan_applications.EmailHandler)
	mux.HandleFunc("/transaction/loan_applications/details/Show_prev_price", loan_applications.Show_prev_price)
	//////////////////////////CWCC_LOAN_APPLICATIONS_END//////////////////////////

	//////////////////////////CWCC_LOAN_CONSOLIDATION_START//////////////////////////
	mux.HandleFunc("/transaction/loan_consolidation/details", loan_consolidation.DListHandler)
	mux.HandleFunc("/transaction/loan_consolidation/details/DaddHandler", loan_consolidation.DAddHandler)
	mux.HandleFunc("/transaction/loan_consolidation/details/DEditHandler", loan_consolidation.DEditHandler)
	mux.HandleFunc("/transaction/loan_consolidation/details/DDeleteHandler", loan_consolidation.DDeleteHandler)
	mux.HandleFunc("/transaction/loan_consolidation/details/DownloadLoanBal", loan_consolidation.Download_LoanBalHandler)
	//////////////////////////CWCC_LOAN_CONSOLIDATION_END//////////////////////////

	/////////////////////////CWCC DASHBOARD START/////////////////////////////////

	mux.HandleFunc("/dashboard/member", dashboard_member.DashboardMember)
	mux.HandleFunc("/dashboard/operations", dashboard_operations.ViewitemHandler)
	mux.HandleFunc("/dashboard/operations/par_partner", dashboard_operations.Dashboard_parpartner)
	mux.HandleFunc("/dashboard/operations/DashboardLoanAging", dashboard_operations.DashboardLoanAging)
	mux.HandleFunc("/dashboard/operations/DashboardLoanRelease", dashboard_operations.DashboardLoanRelease)
	mux.HandleFunc("/dashboard/operations/DashboardLoanPorfolio", dashboard_operations.DashboardLoanPorfolio)
	//mux.HandleFunc("/dashboard/apply_loan",dashboard_apply_loan.ApplyLoan)
	//mux.HandleFunc("/dashboard/apply_loan/email" , dashboard_apply_loan.EmailHandler)

	/////////////////////////CWCC DASHBOARD END/////////////////////////////////

	//////////////////////////CWCC_MERGE_MEMBER_START//////////////////////////
	mux.HandleFunc("/transaction/merge_member", merge_member.HListHandler)
	mux.HandleFunc("/transaction/merge_member/HaddHandler", merge_member.HAddHandler)
	mux.HandleFunc("/transaction/merge_member/HEditHandler", merge_member.HEditHandler)
	mux.HandleFunc("/transaction/merge_member/HDeleteHandler", merge_member.HDeleteHandler)
	mux.HandleFunc("/transaction/merge_member/HViewHandler", merge_member.HViewHandler)
	mux.HandleFunc("/transaction/merge_member/details", merge_member.DListHandler)
	mux.HandleFunc("/transaction/merge_member/details/DaddHandler", merge_member.DAddHandler)
	mux.HandleFunc("/transaction/merge_member/details/DEditHandler", merge_member.DEditHandler)
	mux.HandleFunc("/transaction/merge_member/details/DDeleteHandler", merge_member.DDeleteHandler)
	mux.HandleFunc("/transaction/merge_member/post", merge_member.HPostHandler)
	mux.HandleFunc("/transaction/merge_member/details/cancel", merge_member.DCancelHandler)
	mux.HandleFunc("/transaction/merge_member/details/email", merge_member.EmailHandler)
	mux.HandleFunc("/transaction/merge_member/details/DownloadLoanBal", merge_member.Download_LoanBalHandler)
	mux.HandleFunc("/transaction/merge_member/details/Upload", merge_member.UploadHandler)
	mux.HandleFunc("/transaction/merge_member/details/AjaxPolling", merge_member.AjaxPollingeHandler)
	mux.HandleFunc("/transaction/merge_member/details/deleteall", merge_member.DDeleteAllHandler)

	//////////////////////////CWCC_MERGE_MEMBER_ENd//////////////////////////

	//////////////////////////CWCC_PARTNERS_START//////////////////////////
	mux.HandleFunc("/master/partners", partners.HListHandler)
	mux.HandleFunc("/master/partners/HaddHandler", partners.HAddHandler)
	mux.HandleFunc("/master/partners/HEditHandler", partners.HEditHandler)
	mux.HandleFunc("/master/partners/HDeleteHandler", partners.HDeleteHandler)
	mux.HandleFunc("/master/partners/location", partners.DListHandler)
	mux.HandleFunc("/master/partners/location/DaddHandler", partners.DaddHandler)
	mux.HandleFunc("/master/partners/location/DEditHandler", partners.DEditHandler)
	mux.HandleFunc("/master/partners/location/DDeleteHandler", partners.DDeleteHandler)
	//////////////////////////CWCC_PARTNERS_END//////////////////////////

	//////////////////////////CWCC_LOAN_RELEASE_START//////////////////////////
	mux.HandleFunc("/transaction/loan_release", loan_release.HListHandler)
	mux.HandleFunc("/transaction/loan_release/HaddHandler", loan_release.HAddHandler)
	mux.HandleFunc("/transaction/loan_release/HEditHandler", loan_release.HEditHandler)
	mux.HandleFunc("/transaction/loan_release/HDeleteHandler", loan_release.HDeleteHandler)
	mux.HandleFunc("/transaction/loan_release/HViewHandler", loan_release.HViewHandler)
	mux.HandleFunc("/transaction/loan_release/details", loan_release.DListHandler)
	mux.HandleFunc("/transaction/loan_release/details/DaddHandler", loan_release.DAddHandler)
	mux.HandleFunc("/transaction/loan_release/details/DEditHandler", loan_release.DEditHandler)
	mux.HandleFunc("/transaction/loan_release/details/DDeleteHandler", loan_release.DDeleteHandler)
	mux.HandleFunc("/transaction/loan_release/details/post", loan_release.DPostHandler)
	mux.HandleFunc("/transaction/loan_release/details/cancel", loan_release.DCancelHandler)
	mux.HandleFunc("/transaction/loan_release/details/post", loan_release.DPostHandler)
	//mux.HandleFunc("/inventory/supplier_order/details/Ajax/Ajax_getRef",inventory_supplier_order.Ajax_get_reference)
	mux.HandleFunc("/transaction/loan_release/HPostHandler", loan_release.HPostHandler)
	//mux.HandleFunc("/transaction/loan_applications/Show_savings_ledger",loan_applications.Show_savings_ledger)

	//mux.HandleFunc("/transaction/loan_applications/details/Ajax/Ajax_getRef",escs_loan_applications.Ajax_get_reference)
	//mux.HandleFunc("/inventory/supplier_order/details/Ajax/Ajax_getRef",inventory_supplier_order.Ajax_get_reference)
	mux.HandleFunc("/transaction/loan_release/details/email", loan_release.EmailHandler)

	//////////////////////////CWCC_LOAN_RELEASE_ENd//////////////////////////

	//////////////////////////CWCC_LOAN_PAYMENT_START//////////////////////////
	mux.HandleFunc("/transaction/loan_payment", loan_payment.HListHandler)
	mux.HandleFunc("/transaction/loan_payment/HaddHandler", loan_payment.HAddHandler)
	mux.HandleFunc("/transaction/loan_payment/HEditHandler", loan_payment.HEditHandler)
	mux.HandleFunc("/transaction/loan_payment/HDeleteHandler", loan_payment.HDeleteHandler)
	mux.HandleFunc("/transaction/loan_payment/HViewHandler", loan_payment.HViewHandler)
	mux.HandleFunc("/transaction/loan_payment/details", loan_payment.DListHandler)
	mux.HandleFunc("/transaction/loan_payment/details/DaddHandler", loan_payment.DAddHandler)
	mux.HandleFunc("/transaction/loan_payment/details/DEditHandler", loan_payment.DEditHandler)
	mux.HandleFunc("/transaction/loan_payment/details/DDeleteHandler", loan_payment.DDeleteHandler)
	mux.HandleFunc("/transaction/loan_payment/details/post", loan_payment.DPostHandler)
	mux.HandleFunc("/transaction/loan_payment/details/cancel", loan_payment.DCancelHandler)
	mux.HandleFunc("/transaction/loan_payment/details/email", loan_payment.EmailHandler)
	mux.HandleFunc("/transaction/loan_payment/details/DownloadLoanBal", loan_payment.Download_LoanBalHandler)
	mux.HandleFunc("/transaction/loan_payment/details/Upload", loan_payment.UploadHandler)
	mux.HandleFunc("/transaction/loan_payment/details/AjaxPolling", loan_payment.AjaxPollingeHandler)
	mux.HandleFunc("/transaction/loan_payment/details/deleteall", loan_payment.DDeleteAllHandler)
	mux.HandleFunc("/transaction/loan_payment/details/duplicate", loan_payment.Duplicate)

	//////////////////////////CWCC_LOAN_PAYMENT_ENd//////////////////////////

	//////////////////////////CWCC_LOAN_ADJUSTMENT_START//////////////////////////
	mux.HandleFunc("/transaction/loan_adjustment", loan_adjustment.HListHandler)
	mux.HandleFunc("/transaction/loan_adjustment/HaddHandler", loan_adjustment.HAddHandler)
	mux.HandleFunc("/transaction/loan_adjustment/HEditHandler", loan_adjustment.HEditHandler)
	mux.HandleFunc("/transaction/loan_adjustment/HDeleteHandler", loan_adjustment.HDeleteHandler)
	mux.HandleFunc("/transaction/loan_adjustment/HViewHandler", loan_adjustment.HViewHandler)
	mux.HandleFunc("/transaction/loan_adjustment/details", loan_adjustment.DListHandler)
	mux.HandleFunc("/transaction/loan_adjustment/details/DaddHandler", loan_adjustment.DAddHandler)
	mux.HandleFunc("/transaction/loan_adjustment/details/DEditHandler", loan_adjustment.DEditHandler)
	mux.HandleFunc("/transaction/loan_adjustment/details/DDeleteHandler", loan_adjustment.DDeleteHandler)
	mux.HandleFunc("/transaction/loan_adjustment/details/post", loan_adjustment.DPostHandler)
	mux.HandleFunc("/transaction/loan_adjustment/details/cancel", loan_adjustment.DCancelHandler)
	mux.HandleFunc("/transaction/loan_adjustment/details/email", loan_adjustment.EmailHandler)
	mux.HandleFunc("/transaction/loan_adjustment/details/DownloadLoanBal", loan_adjustment.Download_LoanBalHandler)
	mux.HandleFunc("/transaction/loan_adjustment/details/duplicate", loan_adjustment.Duplicate)

	//////////////////////////CWCC_LOAN_ADJUSTMENT_ENd//////////////////////////

	//////////////////////////CWCC_ADJUSTMENT_ENTRY_START//////////////////////////
	mux.HandleFunc("/transaction/adjustment_entry", adjustment_entry.HListHandler)
	mux.HandleFunc("/transaction/adjustment_entry/HaddHandler", adjustment_entry.HAddHandler)
	mux.HandleFunc("/transaction/adjustment_entry/HEditHandler", adjustment_entry.HEditHandler)
	mux.HandleFunc("/transaction/adjustment_entry/HDeleteHandler", adjustment_entry.HDeleteHandler)
	mux.HandleFunc("/transaction/adjustment_entry/HViewHandler", adjustment_entry.HViewHandler)
	mux.HandleFunc("/transaction/adjustment_entry/details", adjustment_entry.DListHandler)
	mux.HandleFunc("/transaction/adjustment_entry/details/DaddHandler", adjustment_entry.DAddHandler)
	mux.HandleFunc("/transaction/adjustment_entry/details/DEditHandler", adjustment_entry.DEditHandler)
	mux.HandleFunc("/transaction/adjustment_entry/details/DDeleteHandler", adjustment_entry.DDeleteHandler)
	mux.HandleFunc("/transaction/adjustment_entry/details/post", adjustment_entry.DPostHandler)
	mux.HandleFunc("/transaction/adjustment_entry/details/cancel", adjustment_entry.DCancelHandler)
	mux.HandleFunc("/transaction/adjustment_entry/details/email", adjustment_entry.EmailHandler)
	mux.HandleFunc("/transaction/adjustment_entry/details/duplicate", adjustment_entry.Duplicate)

	//////////////////////////CWCC_ADJUSMENT_ENTRY_ENd//////////////////////////

	//////////////////////////CWCC_BULK_ENTRY_START//////////////////////////
	mux.HandleFunc("/transaction/bulk_entry", bulk_entry.HListHandler)
	mux.HandleFunc("/transaction/bulk_entry/HaddHandler", bulk_entry.HAddHandler)
	mux.HandleFunc("/transaction/bulk_entry/HEditHandler", bulk_entry.HEditHandler)
	mux.HandleFunc("/transaction/bulk_entry/HDeleteHandler", bulk_entry.HDeleteHandler)
	mux.HandleFunc("/transaction/bulk_entry/HViewHandler", bulk_entry.HViewHandler)
	mux.HandleFunc("/transaction/bulk_entry/details", bulk_entry.DListHandler)
	mux.HandleFunc("/transaction/bulk_entry/details/DaddHandler", bulk_entry.DAddHandler)
	mux.HandleFunc("/transaction/bulk_entry/details/DEditHandler", bulk_entry.DEditHandler)
	mux.HandleFunc("/transaction/bulk_entry/details/DDeleteHandler", bulk_entry.DDeleteHandler)
	mux.HandleFunc("/transaction/bulk_entry/details/post", bulk_entry.DPostHandler)
	mux.HandleFunc("/transaction/bulk_entry/details/cancel", bulk_entry.DCancelHandler)
	mux.HandleFunc("/transaction/bulk_entry/details/email", bulk_entry.EmailHandler)
	mux.HandleFunc("/transaction/bulk_entry/details/deleteall", bulk_entry.DDeleteAllHandler)
	mux.HandleFunc("/transaction/bulk_entry/details/duplicate", bulk_entry.Duplicate)

	//////////////////////////CWCC_BULK_ENTRY_ENd//////////////////////////

	//////////////////////////CWCC_BULK_ENTRY_START//////////////////////////
	mux.HandleFunc("/transaction/savings_withdrawal", savings_withdrawal.HListHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/HaddHandler", savings_withdrawal.HAddHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/HEditHandler", savings_withdrawal.HEditHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/HDeleteHandler", savings_withdrawal.HDeleteHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/HViewHandler", savings_withdrawal.HViewHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/details", savings_withdrawal.DListHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/details/DaddHandler", savings_withdrawal.DAddHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/details/DEditHandler", savings_withdrawal.DEditHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/details/DDeleteHandler", savings_withdrawal.DDeleteHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/details/post", savings_withdrawal.DPostHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/details/cancel", savings_withdrawal.DCancelHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/details/email", savings_withdrawal.EmailHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/post", savings_withdrawal.HPostHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/cancel", savings_withdrawal.HCancelHandler)
	mux.HandleFunc("/transaction/savings_withdrawal/LPostHandler", savings_withdrawal.LPostHandler)

	//////////////////////////CWCC_BULK_ENTRY_ENd//////////////////////////

	//////////////////////////CWCC_LOAN_PAYMENT_START//////////////////////////
	mux.HandleFunc("/transaction/acknowledgement_receipt", acknowledgement_receipt.HListHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/HaddHandler", acknowledgement_receipt.HAddHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/HEditHandler", acknowledgement_receipt.HEditHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/HDeleteHandler", acknowledgement_receipt.HDeleteHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/HViewHandler", acknowledgement_receipt.HViewHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/details", acknowledgement_receipt.DListHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/details/DaddHandler", acknowledgement_receipt.DAddHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/details/DEditHandler", acknowledgement_receipt.DEditHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/details/DDeleteHandler", acknowledgement_receipt.DDeleteHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/details/post", acknowledgement_receipt.DPostHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/details/deleteall", acknowledgement_receipt.DDeleteAllHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/details/cancel", acknowledgement_receipt.DCancelHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/details/email", acknowledgement_receipt.EmailHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/details/DownloadLoanBal", acknowledgement_receipt.Download_LoanBalHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/details/Upload", acknowledgement_receipt.UploadHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/details/AjaxPolling", acknowledgement_receipt.AjaxPollingeHandler)
	mux.HandleFunc("/transaction/acknowledgement_receipt/details/duplicate", acknowledgement_receipt.Duplicate)

	//////////////////////////CWCC_LOAN_PAYMENT_ENd//////////////////////////

	//////////////////////////CWCC_LOAN_PAYMENT_START//////////////////////////
	mux.HandleFunc("/transaction/other_disbursement", other_disbursement.HListHandler)
	mux.HandleFunc("/transaction/other_disbursement/HaddHandler", other_disbursement.HAddHandler)
	mux.HandleFunc("/transaction/other_disbursement/HEditHandler", other_disbursement.HEditHandler)
	mux.HandleFunc("/transaction/other_disbursement/HDeleteHandler", other_disbursement.HDeleteHandler)
	mux.HandleFunc("/transaction/other_disbursement/HViewHandler", other_disbursement.HViewHandler)
	mux.HandleFunc("/transaction/other_disbursement/details", other_disbursement.DListHandler)
	mux.HandleFunc("/transaction/other_disbursement/details/DaddHandler", other_disbursement.DAddHandler)
	mux.HandleFunc("/transaction/other_disbursement/details/DEditHandler", other_disbursement.DEditHandler)
	mux.HandleFunc("/transaction/other_disbursement/details/DDeleteHandler", other_disbursement.DDeleteHandler)
	mux.HandleFunc("/transaction/other_disbursement/details/post", other_disbursement.DPostHandler)
	mux.HandleFunc("/transaction/other_disbursement/details/deleteall", other_disbursement.DDeleteAllHandler)
	mux.HandleFunc("/transaction/other_disbursement/details/cancel", other_disbursement.DCancelHandler)
	mux.HandleFunc("/transaction/other_disbursement/details/email", other_disbursement.EmailHandler)
	mux.HandleFunc("/transaction/other_disbursement/details/DownloadLoanBal", other_disbursement.Download_LoanBalHandler)
	mux.HandleFunc("/transaction/other_disbursement/details/Upload", other_disbursement.UploadHandler)
	mux.HandleFunc("/transaction/other_disbursement/details/AjaxPolling", other_disbursement.AjaxPollingeHandler)
	mux.HandleFunc("/transaction/other_disbursement/details/duplicate", other_disbursement.Duplicate)

	//////////////////////////CWCC_LOAN_PAYMENT_ENd//////////////////////////

	//////////////////////////CWCC_RESIGN_PROCESS_START//////////////////////////
	mux.HandleFunc("/transaction/resign_process", resign_process.HListHandler)
	mux.HandleFunc("/transaction/resign_process/HaddHandler", resign_process.HAddHandler)
	mux.HandleFunc("/transaction/resign_process/HEditHandler", resign_process.HEditHandler)
	mux.HandleFunc("/transaction/resign_process/HDeleteHandler", resign_process.HDeleteHandler)
	mux.HandleFunc("/transaction/resign_process/HViewHandler", resign_process.HViewHandler)
	mux.HandleFunc("/transaction/resign_process/details", resign_process.DListHandler)
	mux.HandleFunc("/transaction/resign_process/details/DaddHandler", resign_process.DAddHandler)
	mux.HandleFunc("/transaction/resign_process/details/DEditHandler", resign_process.DEditHandler)
	mux.HandleFunc("/transaction/resign_process/details/DDeleteHandler", resign_process.DDeleteHandler)
	mux.HandleFunc("/transaction/resign_process/details/post", resign_process.DPostHandler)
	mux.HandleFunc("/transaction/resign_process/details/cancel", resign_process.DCancelHandler)
	mux.HandleFunc("/transaction/resign_process/details/email", resign_process.EmailHandler)

	mux.HandleFunc("/transaction/resign_process/details/download", resign_process.DDownloadHandler)

	//////////////////////////CWCC_RESIGN_PROCESS_ENd//////////////////////////

	///////////////////////////CWCC END////////////////////////////////////////////////

	mux.HandleFunc("/member", member.AllAccountHAndler)
	mux.HandleFunc("/Member/CreateMember", member.CreateMemberHAndler)
	mux.HandleFunc("/Member/EditMember", member.EditMemberHAndler)
	mux.HandleFunc("/Member/DoEditMember", member.DoEditMemberHAndler)
	mux.HandleFunc("/Member/DeleteMember", member.DeleteMemberHAndler)

	mux.HandleFunc("/Member/MemberAssignment", member.MemberAssigmentHAndler)
	mux.HandleFunc("/Member/MemberAssignmentCreate", member.MemberAssigmentCreateHAndler)
	mux.HandleFunc("/Member/MemberAssignmentDelete", member.MemberAssigmentDeleteHAndler)
	mux.HandleFunc("/Member/MemberAssignmentEdit", member.MemberAssigmentEditHAndler)

	mux.HandleFunc("/assignment", assignment.AssignmentHAndler)
	mux.HandleFunc("/assignment/CreateAssignment", assignment.CreateAssignmentHAndler)
	mux.HandleFunc("/assignment/DeleteAssignment", assignment.DeleteAssignmentHAndler)
	mux.HandleFunc("/assignment/EditAssignment", assignment.EditAssignmentHAndler)
	mux.HandleFunc("/assignment/DuplicateAssignment", assignment.DuplicateAssignmentHAndler)

	mux.HandleFunc("/assignment/AssignmentMembers", assignment.AssignmentMembersHAndler)
	mux.HandleFunc("/assignment/AssignmentMembersCreate", assignment.AssignmentMembersCreateHAndler)
	mux.HandleFunc("/assignment/AssignmentMembersEdit", assignment.AssignmentMembersEditHAndler)
	mux.HandleFunc("/assignment/AssignmentMembersDelete", assignment.AssignmentMembersDeleteHAndler)

	mux.HandleFunc("/get_duplicate_dev_id", Get_duplicate_dev_id)

	mux.HandleFunc("/administrator/role", admin.RoleHAndler)
	mux.HandleFunc("/administrator/createrole", admin.CreateRoleHandler)
	mux.HandleFunc("/administrator/accountrole", account_role.AccountRoleHAndler)

	mux.HandleFunc("/administrator/account_org_role", account_org_role.HListHandler)
	mux.HandleFunc("/administrator/account_org_role/HaddHandler", account_org_role.HAddHandler)
	mux.HandleFunc("/administrator/account_org_role/HEditHandler", account_org_role.HEditHandler)
	mux.HandleFunc("/administrator/account_org_role/HDeleteHandler", account_org_role.HDeleteHandler)

	mux.HandleFunc("/timekeeping", timekeeping.TransactionHandler)
	mux.HandleFunc("/timekeeping/CreateTransaction", timekeeping.CreateTransactionHAndler)
	mux.HandleFunc("/timekeeping/CreateTransactionDetails", timekeeping.CreateTransactionDetailsHAndler)
	mux.HandleFunc("/timekeeping/EditTransactionDetails", timekeeping.EditTransactionDetailsHAndler)
	mux.HandleFunc("/timekeeping/CreateTransactionDelete", timekeeping.CreateTransactionDeleteHAndler)
	mux.HandleFunc("/timekeeping/CreateTransactionDetailsAddnew", timekeeping.CreateTransactionDetailsAddnewHAndler)
	mux.HandleFunc("/timekeeping/CreateTransactionDetailsEdit", timekeeping.CreateTransactionDetailsEditHAndler)
	mux.HandleFunc("/timekeeping/CreateTransactionDetails/ajaxgetmemberbyassigment", timekeeping.AjaxgetmemberbyassigmentHAndler)
	mux.HandleFunc("/timekeeping/CreateTransactionDetailsDelete", timekeeping.CreateTransactionDetailsDeleteHAndler)
	mux.HandleFunc("/timekeeping/transactionDetails/post", timekeeping.CreateTransactionDetailsPostHAndler)
	mux.HandleFunc("/timekeeping/transactionDetails/interpret", timekeeping.CreateTransactionDetailsInterpretHAndler)
	mux.HandleFunc("/timekeeping/overtime_logs", timekeeping_overtime_logs.LBR_OTHdr_ListHandler)
	mux.HandleFunc("/timekeeping/overtime_logs/OvertimeLogsHeaderCreate", timekeeping_overtime_logs.LBR_OTHdr_CreateHandler)
	mux.HandleFunc("/timekeeping/overtime_logs/OvertimeLogsHeaderEdit", timekeeping_overtime_logs.LBR_OTHdr_EditHandler)
	mux.HandleFunc("/timekeeping/overtime_logs/OvertimeLogsHeaderDelete", timekeeping_overtime_logs.LBR_OTHdr_DeleteHandler)

	mux.HandleFunc("/timekeeping/assignmentandlogs", assignment_and_logs.Nav_assignmentandlogs)
	mux.HandleFunc("/timekeeping/assignmentandlogs/dragsave", assignment_and_logs.Dragsave)
	mux.HandleFunc("/Ajax/Get_lbr_lograw_get", assignment_and_logs.Ajax_lbr_lograw_get)

	mux.HandleFunc("/timekeeping/time_keep", time_keep.Nav_assignmentandlogs)
	mux.HandleFunc("/timekeeping/time_keep/dragsave", time_keep.Dragsave)
	mux.HandleFunc("/Ajax/time_keep/Get_lbr_lograw_get", time_keep.Ajax_lbr_lograw_get)

	//timekeep v2
	mux.HandleFunc("/timekeeping/time_keep_v2", time_keep_v2.Nav_assignmentandlogs)

	mux.HandleFunc("/timekeeping/time_keep_uploader", time_keep_uploader.Nav)
	mux.HandleFunc("/timekeeping/time_keep_uploader_do_upload", time_keep_uploader.Time_keep_upload)

	mux.HandleFunc("/timekeeping/holidays", holidays.Holiday_listHandler)
	mux.HandleFunc("/timekeeping/holidays/addnew", holidays.HolidayAddnewHandler)
	mux.HandleFunc("/timekeeping/holidays/edit", holidays.HolidayEditHandler)
	mux.HandleFunc("/timekeeping/holidays/delete", holidays.HolidayDeleteHandler)

	mux.HandleFunc("/timekeeping/overtime_logs/details", timekeeping_overtime_logs.LBR_OTDtl_ListHandler)
	mux.HandleFunc("/timekeeping/overtime_logs/details/OvertimeLogsDetailCreate", timekeeping_overtime_logs.LBR_OTDtl_CreateHandler)
	mux.HandleFunc("/timekeeping/overtime_logs/details/OvertimeLogsDetailEdit", timekeeping_overtime_logs.LBR_OTDtl_EditHandler)
	mux.HandleFunc("/timekeeping/overtime_logs/details/OvertimeLogsDetailDelete", timekeeping_overtime_logs.LBR_OTDtl_DeleteHandler)
	mux.HandleFunc("/timekeeping/overtime_logs/details/OvertimeLogsDetailPost", timekeeping_overtime_logs.LBR_OTDtl_PostHandler)
	mux.HandleFunc("/timekeeping/overtime_logs/details/OvertimeLogsDetailCancel", timekeeping_overtime_logs.LBR_OTDtl_CancelHandler)

	mux.HandleFunc("/timekeeping/rptlogreport", report_translog.Nav_Report_member_log)
	mux.HandleFunc("/timekeeping/report/member_log", report_translog.Report_member_log)
	mux.HandleFunc("/timekeeping/report/member_log2", report_translog.Report_member_log2)
	mux.HandleFunc("/timekeeping/report/member_log3", report_translog.Report_member_log3)
	mux.HandleFunc("/timekeeping/report/Lbr_rep_logpair_rpt", report_translog.Lbr_rep_logpair_rpt)
	mux.HandleFunc("/timekeeping/report/Lbr_rep_lograw_rpt", report_translog.Func_lbr_rep_lograw)

	mux.HandleFunc("/timekeeping/rawreports", rawreports.Rawreports_nav)
	mux.HandleFunc("/timekeeping/rawreports/do_print/", rawreports.Rawreports_doprint)
	mux.HandleFunc("/inventory/reports/do_print", standard_report.Rawreports_doprint)

	mux.HandleFunc("/timekeeping/manningfulfillment", mfr.Nav_Report_mfr)
	mux.HandleFunc("/timekeeping/manningfulfillment/report", mfr.Report)
	mux.HandleFunc("/timekeeping/manningfulfillment/report_main", mfr.Report_main)

	mux.HandleFunc("/inventory/item_class", inventory_itemclass.HListHandler)
	mux.HandleFunc("/inventory/item_class/HaddHandler", inventory_itemclass.HAddHandler)
	mux.HandleFunc("/inventory/item_class/HEditHandler", inventory_itemclass.HEditHandler)
	mux.HandleFunc("/inventory/item_class/HDeleteHandler", inventory_itemclass.HDeleteHandler)

	mux.HandleFunc("/inventory/supplier", inventory_supplier.HListHandler)
	mux.HandleFunc("/inventory/supplier/HaddHandler", inventory_supplier.HAddHandler)
	mux.HandleFunc("/inventory/supplier/HEditHandler", inventory_supplier.HEditHandler)
	mux.HandleFunc("/inventory/supplier/HDeleteHandler", inventory_supplier.HDeleteHandler)
	mux.HandleFunc("/inventory/supplier/HaddTagHandler", inventory_supplier.HAddTagHandler)
	mux.HandleFunc("/inventory/supplier/items", inventory_supplier.DListHandler)
	mux.HandleFunc("/inventory/supplier/items/DaddHandler", inventory_supplier.DaddHandler)
	mux.HandleFunc("/inventory/supplier/items/DEditHandler", inventory_supplier.DEditHandler)
	mux.HandleFunc("/inventory/supplier/items/DDeleteHandler", inventory_supplier.DDeleteHandler)

	mux.HandleFunc("/inventory/courier", inventory_courier.HListHandler)
	mux.HandleFunc("/inventory/courier", inventory_courier.HListHandler)
	mux.HandleFunc("/inventory/courier/HaddHandler", inventory_courier.HAddHandler)
	mux.HandleFunc("/inventory/courier/HEditHandler", inventory_courier.HEditHandler)
	mux.HandleFunc("/inventory/courier/HDeleteHandler", inventory_courier.HDeleteHandler)
	mux.HandleFunc("/inventory/courier/HaddTagHandler", inventory_courier.HAddTagHandler)

	mux.HandleFunc("/inventory/supplier_receipt", inventory_supplierreceipt.HListHandler)
	mux.HandleFunc("/inventory/supplier_receipt/HaddHandler", inventory_supplierreceipt.HAddHandler)
	mux.HandleFunc("/inventory/supplier_receipt/HEditHandler", inventory_supplierreceipt.HEditHandler)
	mux.HandleFunc("/inventory/supplier_receipt/HDeleteHandler", inventory_supplierreceipt.HDeleteHandler)
	mux.HandleFunc("/inventory/supplier_receipt/HaddTagHandler", inventory_supplierreceipt.HAddTagHandler)
	mux.HandleFunc("/inventory/supplier_receipt/details", inventory_supplierreceipt.DListHandler)
	mux.HandleFunc("/inventory/supplier_receipt/details/DaddHandler", inventory_supplierreceipt.DAddHandler)
	mux.HandleFunc("/inventory/supplier_receipt/details/DEditHandler", inventory_supplierreceipt.DEditHandler)
	mux.HandleFunc("/inventory/supplier_receipt/details/DDeleteHandler", inventory_supplierreceipt.DDeleteHandler)
	mux.HandleFunc("/inventory/supplier_receipt/details/Ajax/Ajax_getItem", inventory_supplierreceipt.Ajax_getItem)
	mux.HandleFunc("/inventory/supplier_receipt/details/post", inventory_supplierreceipt.DPostHandler)
	mux.HandleFunc("/inventory/supplier_receipt/details/cancel", inventory_supplierreceipt.DCancelHandler)

	mux.HandleFunc("/inventory/supplier_receipt/details/Show_previous_price", inventory_supplierreceipt.Show_prev_price)
	mux.HandleFunc("/inventory/supplier_receipt/details/Show_balance", inventory_supplierreceipt.Show_balance)
	mux.HandleFunc("/inventory/supplier_receipt/details/DownloadSupOrd", inventory_supplierreceipt.Download_SupOrdHandler)

	mux.HandleFunc("/inventory/supplier_return", inventory_supplierreturn.HListHandler)
	mux.HandleFunc("/inventory/supplier_return/HaddHandler", inventory_supplierreturn.HAddHandler)
	mux.HandleFunc("/inventory/supplier_return/HEditHandler", inventory_supplierreturn.HEditHandler)
	mux.HandleFunc("/inventory/supplier_return/HDeleteHandler", inventory_supplierreturn.HDeleteHandler)
	mux.HandleFunc("/inventory/supplier_return/HaddTagHandler", inventory_supplierreturn.HAddTagHandler)
	mux.HandleFunc("/inventory/supplier_return/details", inventory_supplierreturn.DListHandler)
	mux.HandleFunc("/inventory/supplier_return/details/DaddHandler", inventory_supplierreturn.DAddHandler)
	mux.HandleFunc("/inventory/supplier_return/details/DEditHandler", inventory_supplierreturn.DEditHandler)
	mux.HandleFunc("/inventory/supplier_return/details/DDeleteHandler", inventory_supplierreturn.DDeleteHandler)
	mux.HandleFunc("/inventory/supplier_return/details/Ajax/Ajax_getItem", inventory_supplierreturn.Ajax_getItem)
	mux.HandleFunc("/inventory/supplier_return/details/post", inventory_supplierreturn.DPostHandler)
	mux.HandleFunc("/inventory/supplier_return/details/cancel", inventory_supplierreturn.DCancelHandler)
	mux.HandleFunc("/inventory/supplier_return/details/DownloadCustRet", inventory_supplierreturn.Download_CustRetHandler)
	mux.HandleFunc("/inventory/supplier_return/details/DownloadLocRec", inventory_supplierreturn.Download_LocRecHandler)
	mux.HandleFunc("/inventory/supplier_return/details/email", inventory_supplierreturn.EmailHandler)

	mux.HandleFunc("/inventory/supplier_replacement", inventory_supplierreplacement.HListHandler)
	mux.HandleFunc("/inventory/supplier_replacement/HaddHandler", inventory_supplierreplacement.HAddHandler)
	mux.HandleFunc("/inventory/supplier_replacement/HEditHandler", inventory_supplierreplacement.HEditHandler)
	mux.HandleFunc("/inventory/supplier_replacement/HDeleteHandler", inventory_supplierreplacement.HDeleteHandler)
	mux.HandleFunc("/inventory/supplier_replacement/HaddTagHandler", inventory_supplierreplacement.HAddTagHandler)
	mux.HandleFunc("/inventory/supplier_replacement/details", inventory_supplierreplacement.DListHandler)
	mux.HandleFunc("/inventory/supplier_replacement/details/DaddHandler", inventory_supplierreplacement.DAddHandler)
	mux.HandleFunc("/inventory/supplier_replacement/details/DEditHandler", inventory_supplierreplacement.DEditHandler)
	mux.HandleFunc("/inventory/supplier_replacement/details/DDeleteHandler", inventory_supplierreplacement.DDeleteHandler)
	mux.HandleFunc("/inventory/supplier_replacement/details/Ajax/Ajax_getItem", inventory_supplierreplacement.Ajax_getItem)
	mux.HandleFunc("/inventory/supplier_replacement/details/post", inventory_supplierreplacement.DPostHandler)
	mux.HandleFunc("/inventory/supplier_replacement/details/cancel", inventory_supplierreplacement.DCancelHandler)
	mux.HandleFunc("/inventory/supplier_replacement/details/DownloadSupRet", inventory_supplierreplacement.Download_SupRetHandler)
	mux.HandleFunc("/inventory/supplier_replacement/details/Show_prev_price", inventory_supplierreplacement.Show_prev_price)

	mux.HandleFunc("/inventory/adjustment_entry", inventory_adjustmententry.HListHandler)
	mux.HandleFunc("/inventory/adjustment_entry/HaddHandler", inventory_adjustmententry.HAddHandler)
	mux.HandleFunc("/inventory/adjustment_entry/HEditHandler", inventory_adjustmententry.HEditHandler)
	mux.HandleFunc("/inventory/adjustment_entry/HDeleteHandler", inventory_adjustmententry.HDeleteHandler)
	mux.HandleFunc("/inventory/adjustment_entry/HaddTagHandler", inventory_adjustmententry.HAddTagHandler)
	mux.HandleFunc("/inventory/adjustment_entry/details", inventory_adjustmententry.DListHandler)
	mux.HandleFunc("/inventory/adjustment_entry/details/DaddHandler", inventory_adjustmententry.DAddHandler)
	mux.HandleFunc("/inventory/adjustment_entry/details/DEditHandler", inventory_adjustmententry.DEditHandler)
	mux.HandleFunc("/inventory/adjustment_entry/details/DDeleteHandler", inventory_adjustmententry.DDeleteHandler)
	mux.HandleFunc("/inventory/adjustment_entry/details/Ajax/Ajax_getItem", inventory_adjustmententry.Ajax_getItem)
	mux.HandleFunc("/inventory/adjustment_entry/details/post", inventory_adjustmententry.DPostHandler)
	mux.HandleFunc("/inventory/adjustment_entry/details/cancel", inventory_adjustmententry.DCancelHandler)

	mux.HandleFunc("/inventory/customer_sales", inventory_customersales.HListHandler)
	mux.HandleFunc("/inventory/customer_sales/HaddHandler", inventory_customersales.HAddHandler)
	mux.HandleFunc("/inventory/customer_sales/HEditHandler", inventory_customersales.HEditHandler)
	mux.HandleFunc("/inventory/customer_sales/HDeleteHandler", inventory_customersales.HDeleteHandler)
	mux.HandleFunc("/inventory/customer_sales/HaddTagHandler", inventory_customersales.HAddTagHandler)
	mux.HandleFunc("/inventory/customer_sales/details", inventory_customersales.DListHandler)
	mux.HandleFunc("/inventory/customer_sales/details/DaddHandler", inventory_customersales.DAddHandler)
	mux.HandleFunc("/inventory/customer_sales/details/DEditHandler", inventory_customersales.DEditHandler)
	mux.HandleFunc("/inventory/customer_sales/details/DDeleteHandler", inventory_customersales.DDeleteHandler)
	mux.HandleFunc("/inventory/customer_sales/details/Ajax/Ajax_getItem", inventory_customersales.Ajax_getItem)
	mux.HandleFunc("/inventory/customer_sales/details/Ajax/Ajax_getItem_isbarcode", inventory_customersales.Ajax_getItem_isbarcode)
	mux.HandleFunc("/inventory/customer_sales/details/post", inventory_customersales.DPostHandler)
	mux.HandleFunc("/inventory/customer_sales/details/cancel", inventory_customersales.DCancelHandler)
	mux.HandleFunc("/inventory/customer_sales/details/add_reference", inventory_customersales.DAdd_reference)
	mux.HandleFunc("/inventory/customer_sales/details/delete_reference", inventory_customersales.Ddelete_reference)
	mux.HandleFunc("/inventory/customer_sales/details/Ajax/Ajax_getRef", inventory_customersales.Ajax_get_reference)
	mux.HandleFunc("/inventory/customer_sales/details/duplicate_sales", inventory_customersales.Dduplicate_sales)
	mux.HandleFunc("/inventory/customer_sales/details/Show_prev_price", inventory_customersales.Show_prev_price)
	mux.HandleFunc("/inventory/customer_sales/details/DownloadCusOrd", inventory_customersales.Download_CusOrdHandler)
	mux.HandleFunc("/inventory/customer_sales/details/email", inventory_customersales.EmailHandler)

	mux.HandleFunc("/inventory/customer_return", inventory_customer_return.HListHandler)
	mux.HandleFunc("/inventory/customer_return/HaddHandler", inventory_customer_return.HAddHandler)
	mux.HandleFunc("/inventory/customer_return/HEditHandler", inventory_customer_return.HEditHandler)
	mux.HandleFunc("/inventory/customer_return/HDeleteHandler", inventory_customer_return.HDeleteHandler)
	mux.HandleFunc("/inventory/customer_return/HaddTagHandler", inventory_customer_return.HAddTagHandler)
	mux.HandleFunc("/inventory/customer_return/details", inventory_customer_return.DListHandler)
	mux.HandleFunc("/inventory/customer_return/details/DaddHandler", inventory_customer_return.DAddHandler)
	mux.HandleFunc("/inventory/customer_return/details/DEditHandler", inventory_customer_return.DEditHandler)
	mux.HandleFunc("/inventory/customer_return/details/DDeleteHandler", inventory_customer_return.DDeleteHandler)
	mux.HandleFunc("/inventory/customer_return/details/Ajax/Ajax_getItem", inventory_customer_return.Ajax_getItem)
	mux.HandleFunc("/inventory/customer_return/details/Ajax/Ajax_getItem_isbarcode", inventory_customer_return.Ajax_getItem_isbarcode)
	mux.HandleFunc("/inventory/customer_return/details/post", inventory_customer_return.DPostHandler)
	mux.HandleFunc("/inventory/customer_return/details/cancel", inventory_customer_return.DCancelHandler)
	mux.HandleFunc("/inventory/customer_return/details/add_reference", inventory_customer_return.DAdd_reference)
	mux.HandleFunc("/inventory/customer_return/details/delete_reference", inventory_customer_return.Ddelete_reference)
	mux.HandleFunc("/inventory/customer_return/details/Ajax/Ajax_getRef", inventory_customer_return.Ajax_get_reference)
	mux.HandleFunc("/inventory/customer_return/details/Show_prev_price", inventory_customer_return.Show_prev_price)

	mux.HandleFunc("/inventory/customer_adjustment", inventory_customer_adjustment.HListHandler)
	mux.HandleFunc("/inventory/customer_adjustment/HaddHandler", inventory_customer_adjustment.HAddHandler)
	mux.HandleFunc("/inventory/customer_adjustment/HEditHandler", inventory_customer_adjustment.HEditHandler)
	mux.HandleFunc("/inventory/customer_adjustment/HDeleteHandler", inventory_customer_adjustment.HDeleteHandler)
	mux.HandleFunc("/inventory/customer_adjustment/HPostHandler", inventory_customer_adjustment.HPostHandler)
	mux.HandleFunc("/inventory/customer_adjustment/HaddTagHandler", inventory_customer_adjustment.HAddTagHandler)
	mux.HandleFunc("/inventory/customer_adjustment/details", inventory_customer_adjustment.DListHandler)
	mux.HandleFunc("/inventory/customer_adjustment/details/DaddHandler", inventory_customer_adjustment.DAddHandler)
	mux.HandleFunc("/inventory/customer_adjustment/details/DEditHandler", inventory_customer_adjustment.DEditHandler)
	mux.HandleFunc("/inventory/customer_adjustment/details/DDeleteHandler", inventory_customer_adjustment.DDeleteHandler)
	mux.HandleFunc("/inventory/customer_adjustment/details/Ajax/Ajax_getItem", inventory_customer_adjustment.Ajax_getItem)
	mux.HandleFunc("/inventory/customer_adjustment/details/Ajax/Ajax_getItem_isbarcode", inventory_customer_adjustment.Ajax_getItem_isbarcode)
	mux.HandleFunc("/inventory/customer_adjustment/details/post", inventory_customer_adjustment.DPostHandler)
	mux.HandleFunc("/inventory/customer_adjustment/details/cancel", inventory_customer_adjustment.DCancelHandler)
	mux.HandleFunc("/inventory/customer_adjustment/details/add_reference", inventory_customer_adjustment.DAdd_reference)
	mux.HandleFunc("/inventory/customer_adjustment/details/delete_reference", inventory_customer_adjustment.Ddelete_reference)
	mux.HandleFunc("/inventory/customer_adjustment/details/Ajax/Ajax_getRef", inventory_customer_adjustment.Ajax_get_reference)

	mux.HandleFunc("/inventory/customer_payment", inventory_customer_payment.HListHandler)
	mux.HandleFunc("/inventory/customer_payment/HaddHandler", inventory_customer_payment.HAddHandler)
	mux.HandleFunc("/inventory/customer_payment/HEditHandler", inventory_customer_payment.HEditHandler)
	mux.HandleFunc("/inventory/customer_payment/HDeleteHandler", inventory_customer_payment.HDeleteHandler)
	mux.HandleFunc("/inventory/customer_payment/HPostHandler", inventory_customer_payment.HPostHandler)
	mux.HandleFunc("/inventory/customer_payment/HaddTagHandler", inventory_customer_payment.HAddTagHandler)
	mux.HandleFunc("/inventory/customer_payment/details", inventory_customer_payment.DListHandler)
	mux.HandleFunc("/inventory/customer_payment/details/DaddHandler", inventory_customer_payment.DAddHandler)
	mux.HandleFunc("/inventory/customer_payment/details/DEditHandler", inventory_customer_payment.DEditHandler)
	mux.HandleFunc("/inventory/customer_payment/details/DDeleteHandler", inventory_customer_payment.DDeleteHandler)
	mux.HandleFunc("/inventory/customer_payment/details/Ajax/Ajax_getItem", inventory_customer_payment.Ajax_getItem)
	mux.HandleFunc("/inventory/customer_payment/details/Ajax/Ajax_getItem_isbarcode", inventory_customer_payment.Ajax_getItem_isbarcode)
	mux.HandleFunc("/inventory/customer_payment/details/post", inventory_customer_payment.DPostHandler)
	mux.HandleFunc("/inventory/customer_payment/details/cancel", inventory_customer_payment.DCancelHandler)
	mux.HandleFunc("/inventory/customer_payment/details/add_reference", inventory_customer_payment.DAdd_reference)
	mux.HandleFunc("/inventory/customer_payment/details/delete_reference", inventory_customer_payment.Ddelete_reference)
	mux.HandleFunc("/inventory/customer_payment/details/Ajax/Ajax_getRef", inventory_customer_payment.Ajax_get_reference)

	mux.HandleFunc("/inventory/customer_replacement", inventory_customer_replacement.HListHandler)
	mux.HandleFunc("/inventory/customer_replacement/HaddHandler", inventory_customer_replacement.HAddHandler)
	mux.HandleFunc("/inventory/customer_replacement/HEditHandler", inventory_customer_replacement.HEditHandler)
	mux.HandleFunc("/inventory/customer_replacement/HDeleteHandler", inventory_customer_replacement.HDeleteHandler)
	mux.HandleFunc("/inventory/customer_replacement/HaddTagHandler", inventory_customer_replacement.HAddTagHandler)
	mux.HandleFunc("/inventory/customer_replacement/details", inventory_customer_replacement.DListHandler)
	mux.HandleFunc("/inventory/customer_replacement/details/DaddHandler", inventory_customer_replacement.DAddHandler)
	mux.HandleFunc("/inventory/customer_replacement/details/DEditHandler", inventory_customer_replacement.DEditHandler)
	mux.HandleFunc("/inventory/customer_replacement/details/DDeleteHandler", inventory_customer_replacement.DDeleteHandler)
	mux.HandleFunc("/inventory/customer_replacement/details/Ajax/Ajax_getItem", inventory_customer_replacement.Ajax_getItem)
	mux.HandleFunc("/inventory/customer_replacement/details/Ajax/Ajax_getItem_isbarcode", inventory_customer_replacement.Ajax_getItem_isbarcode)
	mux.HandleFunc("/inventory/customer_replacement/details/post", inventory_customer_replacement.DPostHandler)
	mux.HandleFunc("/inventory/customer_replacement/details/cancel", inventory_customer_replacement.DCancelHandler)
	mux.HandleFunc("/inventory/customer_replacement/details/add_reference", inventory_customer_replacement.DAdd_reference)
	mux.HandleFunc("/inventory/customer_replacement/details/delete_reference", inventory_customer_replacement.Ddelete_reference)
	mux.HandleFunc("/inventory/customer_replacement/details/Ajax/Ajax_getRef", inventory_customer_replacement.Ajax_get_reference)
	mux.HandleFunc("/inventory/customer_replacement/details/Show_prev_price", inventory_customer_replacement.Show_prev_price)
	mux.HandleFunc("/inventory/customer_replacement/details/DownloadCusRet", inventory_customer_replacement.Download_CusRetHandler)
	mux.HandleFunc("/inventory/customer_replacement/details/DownloadSupRep", inventory_customer_replacement.Download_SupRepHandler)
	mux.HandleFunc("/inventory/customer_replacement/details/Show_prev_price", inventory_customer_replacement.Show_prev_price)

	mux.HandleFunc("/inventory/location_reciept", inventory_location_reciept.HListHandler)
	mux.HandleFunc("/inventory/location_reciept/HaddHandler", inventory_location_reciept.HAddHandler)
	mux.HandleFunc("/inventory/location_reciept/HEditHandler", inventory_location_reciept.HEditHandler)
	mux.HandleFunc("/inventory/location_reciept/HDeleteHandler", inventory_location_reciept.HDeleteHandler)
	mux.HandleFunc("/inventory/location_reciept/HaddTagHandler", inventory_location_reciept.HAddTagHandler)
	mux.HandleFunc("/inventory/location_reciept/details", inventory_location_reciept.DListHandler)
	mux.HandleFunc("/inventory/location_reciept/details/DaddHandler", inventory_location_reciept.DAddHandler)
	mux.HandleFunc("/inventory/location_reciept/details/DEditHandler", inventory_location_reciept.DEditHandler)
	mux.HandleFunc("/inventory/location_reciept/details/DDeleteHandler", inventory_location_reciept.DDeleteHandler)
	mux.HandleFunc("/inventory/location_reciept/details/Ajax/Ajax_getItem", inventory_location_reciept.Ajax_getItem)
	mux.HandleFunc("/inventory/location_reciept/details/post", inventory_location_reciept.DPostHandler)
	mux.HandleFunc("/inventory/location_reciept/details/cancel", inventory_location_reciept.DCancelHandler)
	mux.HandleFunc("/inventory/location_reciept/details/add_reference", inventory_location_reciept.DAdd_reference)
	mux.HandleFunc("/inventory/location_reciept/details/delete_reference", inventory_location_reciept.Ddelete_reference)
	mux.HandleFunc("/inventory/location_reciept/details/Ajax/Ajax_getRef", inventory_location_reciept.Ajax_get_reference)
	mux.HandleFunc("/inventory/location_reciept/details/DownloadLocRec", inventory_location_reciept.Download_LocRecHandler)

	mux.HandleFunc("/inventory/location_transfer", inventory_location_transfer.HListHandler)
	mux.HandleFunc("/inventory/location_transfer/HaddHandler", inventory_location_transfer.HAddHandler)
	mux.HandleFunc("/inventory/location_transfer/HEditHandler", inventory_location_transfer.HEditHandler)
	mux.HandleFunc("/inventory/location_transfer/HDeleteHandler", inventory_location_transfer.HDeleteHandler)
	mux.HandleFunc("/inventory/location_transfer/HaddTagHandler", inventory_location_transfer.HAddTagHandler)
	mux.HandleFunc("/inventory/location_transfer/details", inventory_location_transfer.DListHandler)
	mux.HandleFunc("/inventory/location_transfer/details/DaddHandler", inventory_location_transfer.DAddHandler)
	mux.HandleFunc("/inventory/location_transfer/details/DEditHandler", inventory_location_transfer.DEditHandler)
	mux.HandleFunc("/inventory/location_transfer/details/DDeleteHandler", inventory_location_transfer.DDeleteHandler)
	mux.HandleFunc("/inventory/location_transfer/details/Ajax/Ajax_getItem", inventory_location_transfer.Ajax_getItem)
	mux.HandleFunc("/inventory/location_transfer/details/post", inventory_location_transfer.DPostHandler)
	mux.HandleFunc("/inventory/location_transfer/details/cancel", inventory_location_transfer.DCancelHandler)
	mux.HandleFunc("/inventory/location_transfer/details/add_reference", inventory_location_transfer.DAdd_reference)
	mux.HandleFunc("/inventory/location_transfer/details/delete_reference", inventory_location_transfer.Ddelete_reference)
	mux.HandleFunc("/inventory/location_transfer/details/Ajax/Ajax_getRef", inventory_location_transfer.Ajax_get_reference)
	mux.HandleFunc("/inventory/location_transfer/details/DownloadLocOrd", inventory_location_transfer.Download_LocOrdHandler)
	mux.HandleFunc("/inventory/location_transfer/details/DownloadCusRet", inventory_location_transfer.Download_CusRetHandler)
	mux.HandleFunc("/inventory/location_transfer/details/DownloadLocRec", inventory_location_transfer.Download_LocRecHandler)
	mux.HandleFunc("/inventory/location_transfer/details/DownloadLocOrd", inventory_location_transfer.Download_LocOrdHandler)
	mux.HandleFunc("/inventory/location_transfer/details/DownloadCusRet", inventory_location_transfer.Download_CusRetHandler)
	mux.HandleFunc("/inventory/location_transfer/details/email", inventory_location_transfer.EmailHandler)

	mux.HandleFunc("/inventory/location_transfer/details/duplicate_loctra", inventory_location_transfer.Dduplicate_loctra)

	mux.HandleFunc("/inventory/location_order", inventory_location_order.HListHandler)
	mux.HandleFunc("/inventory/location_order/HaddHandler", inventory_location_order.HAddHandler)
	mux.HandleFunc("/inventory/location_order/HEditHandler", inventory_location_order.HEditHandler)
	mux.HandleFunc("/inventory/location_order/HDeleteHandler", inventory_location_order.HDeleteHandler)
	mux.HandleFunc("/inventory/location_order/HaddTagHandler", inventory_location_order.HAddTagHandler)
	mux.HandleFunc("/inventory/location_order/details", inventory_location_order.DListHandler)
	mux.HandleFunc("/inventory/location_order/details/DaddHandler", inventory_location_order.DAddHandler)
	mux.HandleFunc("/inventory/location_order/details/DEditHandler", inventory_location_order.DEditHandler)
	mux.HandleFunc("/inventory/location_order/details/DDeleteHandler", inventory_location_order.DDeleteHandler)
	mux.HandleFunc("/inventory/location_order/details/Ajax/Ajax_getItem", inventory_location_order.Ajax_getItem)
	mux.HandleFunc("/inventory/location_order/details/post", inventory_location_order.DPostHandler)
	mux.HandleFunc("/inventory/location_order/details/cancel", inventory_location_order.DCancelHandler)
	mux.HandleFunc("/inventory/location_order/details/add_reference", inventory_location_order.DAdd_reference)
	mux.HandleFunc("/inventory/location_order/details/delete_reference", inventory_location_order.Ddelete_reference)
	mux.HandleFunc("/inventory/location_order/details/Ajax/Ajax_getRef", inventory_location_order.Ajax_get_reference)
	mux.HandleFunc("/inventory/location_order/details/DVoidHandler", inventory_location_order.DVoidHandler)

	mux.HandleFunc("/inventory/supplier_order", inventory_supplier_order.HListHandler)
	mux.HandleFunc("/inventory/supplier_order/HaddHandler", inventory_supplier_order.HAddHandler)
	mux.HandleFunc("/inventory/supplier_order/HEditHandler", inventory_supplier_order.HEditHandler)
	mux.HandleFunc("/inventory/supplier_order/HDeleteHandler", inventory_supplier_order.HDeleteHandler)
	mux.HandleFunc("/inventory/supplier_order/HaddTagHandler", inventory_supplier_order.HAddTagHandler)
	mux.HandleFunc("/inventory/supplier_order/details", inventory_supplier_order.DListHandler)
	mux.HandleFunc("/inventory/supplier_order/details/DaddHandler", inventory_supplier_order.DAddHandler)
	mux.HandleFunc("/inventory/supplier_order/details/DEditHandler", inventory_supplier_order.DEditHandler)
	mux.HandleFunc("/inventory/supplier_order/details/DDeleteHandler", inventory_supplier_order.DDeleteHandler)
	mux.HandleFunc("/inventory/supplier_order/details/Ajax/Ajax_getItem", inventory_supplier_order.Ajax_getItem)
	mux.HandleFunc("/inventory/supplier_order/details/post", inventory_supplier_order.DPostHandler)
	mux.HandleFunc("/inventory/supplier_order/details/cancel", inventory_supplier_order.DCancelHandler)
	//mux.HandleFunc("/inventory/supplier_order/details/Ajax/Ajax_getRef",inventory_supplier_order.Ajax_get_reference)
	mux.HandleFunc("/inventory/supplier_order/details/email", inventory_supplier_order.EmailHandler)
	mux.HandleFunc("/inventory/supplier_order/details/Show_prev_price", inventory_supplier_order.Show_prev_price)
	mux.HandleFunc("/inventory/supplier_order/details/DVoidHandler", inventory_supplier_order.DVoidHandler)
	mux.HandleFunc("/inventory/supplier_order/details/Ajax/Ajax_getRef", inventory_location_reciept.Ajax_get_reference)
	//mux.HandleFunc("/inventory/supplier_order/details/Ajax/Ajax_getRef",inventory_supplier_order.Ajax_get_reference)
	mux.HandleFunc("/inventory/supplier_order/details/email", inventory_supplier_order.EmailHandler)
	mux.HandleFunc("/inventory/supplier_order/details/Show_prev_price", inventory_supplier_order.Show_prev_price)

	mux.HandleFunc("/inventory/item", inventory_item.HListHandler)
	mux.HandleFunc("/inventory/item/HaddHandler", inventory_item.HAddHandler)
	mux.HandleFunc("/inventory/item/HEditHandler", inventory_item.HEditHandler)
	mux.HandleFunc("/inventory/item/HDeleteHandler", inventory_item.HDeleteHandler)
	mux.HandleFunc("/inventory/item/HaddTagHandler", inventory_item.HAddTagHandler)
	mux.HandleFunc("/inventory/item/HDuplicateHandler", inventory_item.HDuplicateHandler)
	mux.HandleFunc("/inventory/item/supplier", inventory_item.DListHandler)
	mux.HandleFunc("/inventory/item/supplier/DaddHandler", inventory_item.DaddHandler)
	mux.HandleFunc("/inventory/item/supplier/DEditHandler", inventory_item.DEditHandler)
	mux.HandleFunc("/inventory/item/supplier/DDeleteHandler", inventory_item.DDeleteHandler)

	mux.HandleFunc("/inventory/location_items", inventory_location_item.HListHandler)
	mux.HandleFunc("/inventory/location_items/HaddHandler", inventory_location_item.HAddHandler)
	mux.HandleFunc("/inventory/location_items/HEditHandler", inventory_location_item.HEditHandler)
	mux.HandleFunc("/inventory/location_items/HDeleteHandler", inventory_location_item.HDeleteHandler)
	mux.HandleFunc("/inventory/location_items/HaddTagHandler", inventory_location_item.HAddTagHandler)
	mux.HandleFunc("/inventory/location_items/HDuplicateHandler", inventory_location_item.HDuplicateHandler)
	mux.HandleFunc("/inventory/location_items/Ajax_getItem", inventory_location_item.Ajax_getItem)

	mux.HandleFunc("/inventory/banks", inventory_banks.HListHandler)
	mux.HandleFunc("/inventory/banks/HaddHandler", inventory_banks.HAddHandler)
	mux.HandleFunc("/inventory/banks/HEditHandler", inventory_banks.HEditHandler)
	mux.HandleFunc("/inventory/banks/HDeleteHandler", inventory_banks.HDeleteHandler)
	mux.HandleFunc("/inventory/banks/HaddTagHandler", inventory_banks.HAddTagHandler)
	mux.HandleFunc("/inventory/banks/HDuplicateHandler", inventory_banks.HDuplicateHandler)
	mux.HandleFunc("/inventory/banks/Ajax_getItem", inventory_banks.Ajax_getItem)

	mux.HandleFunc("/inventory/location", inventory_location.HListHandler)
	mux.HandleFunc("/inventory/location/HaddHandler", inventory_location.HAddHandler)
	mux.HandleFunc("/inventory/location/HEditHandler", inventory_location.HEditHandler)
	mux.HandleFunc("/inventory/location/HDeleteHandler", inventory_location.HDeleteHandler)
	mux.HandleFunc("/inventory/location/HaddTagHandler", inventory_location.HaddTagHandler)
	mux.HandleFunc("/inventory/location/customer", inventory_location.DListHandler)
	mux.HandleFunc("/inventory/location/customer/DaddHandler", inventory_location.DaddHandler)
	mux.HandleFunc("/inventory/location/customer/DEditHandler", inventory_location.DEditHandler)
	mux.HandleFunc("/inventory/location/customer/DDeleteHandler", inventory_location.DDeleteHandler)
	mux.HandleFunc("/inventory/location/HAddChildHandler", inventory_location.HAddChildHandler)

	mux.HandleFunc("/inventory/internal", inventory_internal.HListHandler)
	mux.HandleFunc("/inventory/internal/HaddHandler", inventory_internal.HAddHandler)
	mux.HandleFunc("/inventory/internal/HEditHandler", inventory_internal.HEditHandler)
	mux.HandleFunc("/inventory/internal/HDeleteHandler", inventory_internal.HDeleteHandler)
	mux.HandleFunc("/inventory/internal/HaddTagHandler", inventory_internal.HAddTagHandler)

	mux.HandleFunc("/inventory/customer", inventory_customer.HListHandler)
	mux.HandleFunc("/inventory/customer/HaddHandler", inventory_customer.HAddHandler)
	mux.HandleFunc("/inventory/customer/HEditHandler", inventory_customer.HEditHandler)
	mux.HandleFunc("/inventory/customer/HDeleteHandler", inventory_customer.HDeleteHandler)
	mux.HandleFunc("/inventory/customer/HaddTagHandler", inventory_customer.HAddTagHandler)
	mux.HandleFunc("/inventory/customer/location", inventory_customer.DListHandler)
	mux.HandleFunc("/inventory/customer/location/DaddHandler", inventory_customer.DaddHandler)
	mux.HandleFunc("/inventory/customer/location/DEditHandler", inventory_customer.DEditHandler)
	mux.HandleFunc("/inventory/customer/location/DDeleteHandler", inventory_customer.DDeleteHandler)

	mux.HandleFunc("/inventory/reports", inventory_reports.Rawreports_nav)
	mux.HandleFunc("/inventory/reports/do_print/", inventory_reports.Rawreports_doprint)

	mux.HandleFunc("/inventory/custom_report", inventory_custom_report.Rawreports_nav)
	mux.HandleFunc("/inventory/custom_report/do_print", inventory_custom_report.Rawreports_doprint)

	mux.HandleFunc("/inventory/special_rpt", inventory_reports.Special_rpt)

	mux.HandleFunc("/inventory/ledger", ledger.HListHandler)
	mux.HandleFunc("/inventory/ledger/viewer", ledger.ViewerHandler)
	mux.HandleFunc("/inventory/ledger/ajaxviewer", ledger.AjaxViewerHandler)
	mux.HandleFunc("/inventory/ledger/ajaxviewer_details", ledger.AjaxViewerDetailsHandler)
	mux.HandleFunc("/inventory/ledger/Ajax_getItem", ledger.Ajax_getItem)

	mux.HandleFunc("/inventory/reports/non_selling_aging", non_selling_aging.HListHandler)
	mux.HandleFunc("/inventory/reports/non_selling_aging/viewer", non_selling_aging.ViewerHandler)
	mux.HandleFunc("/inventory/reports/non_selling_aging/ajaxviewer", non_selling_aging.AjaxViewerHandler)
	mux.HandleFunc("/inventory/reports/non_selling_aging/ajaxviewer_details", non_selling_aging.AjaxViewerDetailsHandler)
	mux.HandleFunc("/inventory/reports/item_profit", item_profit.HListHandler)
	mux.HandleFunc("/inventory/reports/item_profit/viewer", item_profit.ViewerHandler)
	mux.HandleFunc("/inventory/reports/item_profit/ajaxviewer", item_profit.AjaxViewerHandler)
	mux.HandleFunc("/inventory/reports/item_profit/ajaxviewer_details", item_profit.AjaxViewerDetailsHandler)

	mux.HandleFunc("/dashboard", dashboard.HListHandler)

	mux.HandleFunc("/report/product_movement", inventory_product_movement.HListHandler)
	mux.HandleFunc("/inventory/report/product_movement_per_item", inventory_product_movement.ViewitemHandler)

	mux.HandleFunc("/inventory/reports/product_movement/viewer", inventory_product_movement.ViewerHandler)
	mux.HandleFunc("/inventory/reports/product_movement/ajaxviewer", inventory_product_movement.AjaxViewerHandler)
	mux.HandleFunc("/inventory/reports/product_movement/ajaxviewer_details", inventory_product_movement.AjaxViewerDetailsHandler)

	mux.HandleFunc("/inventory/brochure", inventory_brochure.HListHandler)
	mux.HandleFunc("/inventory/brochure/HaddHandler", inventory_brochure.HAddHandler)
	mux.HandleFunc("/inventory/brochure/HEditHandler", inventory_brochure.HEditHandler)
	mux.HandleFunc("/inventory/brochure/HDeleteHandler", inventory_brochure.HDeleteHandler)
	mux.HandleFunc("/inventory/brochure/HaddTagHandler", inventory_brochure.HAddTagHandler)
	mux.HandleFunc("/inventory/brochure/details", inventory_brochure.DListHandler)
	mux.HandleFunc("/inventory/brochure/details/DaddHandler", inventory_brochure.DAddHandler)
	mux.HandleFunc("/inventory/brochure/details/DEditHandler", inventory_brochure.DEditHandler)
	mux.HandleFunc("/inventory/brochure/details/DDeleteHandler", inventory_brochure.DDeleteHandler)
	mux.HandleFunc("/inventory/brochure/details/Ajax/Ajax_getItem", inventory_brochure.Ajax_getItem)
	mux.HandleFunc("/inventory/brochure/details/Ajax/Ajax_getItem_isbarcode", inventory_brochure.Ajax_getItem_isbarcode)
	mux.HandleFunc("/inventory/brochure/details/post", inventory_brochure.DPostHandler)
	mux.HandleFunc("/inventory/brochure/details/cancel", inventory_brochure.DCancelHandler)

	mux.HandleFunc("/inventory/actual_count", inventory_actualcount.HListHandler)
	mux.HandleFunc("/inventory/actual_count/HaddHandler", inventory_actualcount.HAddHandler)
	mux.HandleFunc("/inventory/actual_count/HEditHandler", inventory_actualcount.HEditHandler)
	mux.HandleFunc("/inventory/actual_count/HDeleteHandler", inventory_actualcount.HDeleteHandler)
	mux.HandleFunc("/inventory/actual_count/HaddTagHandler", inventory_actualcount.HAddTagHandler)
	mux.HandleFunc("/inventory/actual_count/details", inventory_actualcount.DListHandler)
	mux.HandleFunc("/inventory/actual_count/details/DaddHandler", inventory_actualcount.DAddHandler)
	mux.HandleFunc("/inventory/actual_count/details/DEditHandler", inventory_actualcount.DEditHandler)
	mux.HandleFunc("/inventory/actual_count/details/DDeleteHandler", inventory_actualcount.DDeleteHandler)
	mux.HandleFunc("/inventory/actual_count/details/Ajax/Ajax_getItem", inventory_actualcount.Ajax_getItem)
	mux.HandleFunc("/inventory/actual_count/details/Ajax/Ajax_getItem_isbarcode", inventory_actualcount.Ajax_getItem_isbarcode)
	mux.HandleFunc("/inventory/actual_count/details/post", inventory_actualcount.DPostHandler)
	mux.HandleFunc("/inventory/actual_count/details/cancel", inventory_actualcount.DCancelHandler)

	mux.HandleFunc("/inventory/customer_order", inventory_customer_order.HListHandler)
	mux.HandleFunc("/inventory/customer_order/HaddHandler", inventory_customer_order.HAddHandler)
	mux.HandleFunc("/inventory/customer_order/HEditHandler", inventory_customer_order.HEditHandler)
	mux.HandleFunc("/inventory/customer_order/HDeleteHandler", inventory_customer_order.HDeleteHandler)
	mux.HandleFunc("/inventory/customer_order/HaddTagHandler", inventory_customer_order.HAddTagHandler)
	mux.HandleFunc("/inventory/customer_order/details", inventory_customer_order.DListHandler)
	mux.HandleFunc("/inventory/customer_order/details/DaddHandler", inventory_customer_order.DAddHandler)
	mux.HandleFunc("/inventory/customer_order/details/DEditHandler", inventory_customer_order.DEditHandler)
	mux.HandleFunc("/inventory/customer_order/details/DDeleteHandler", inventory_customer_order.DDeleteHandler)
	mux.HandleFunc("/inventory/customer_order/details/Ajax/Ajax_getItem", inventory_customer_order.Ajax_getItem)
	mux.HandleFunc("/inventory/customer_order/details/Ajax/Ajax_getItem_isbarcode", inventory_customer_order.Ajax_getItem_isbarcode)
	mux.HandleFunc("/inventory/customer_order/details/post", inventory_customer_order.DPostHandler)
	mux.HandleFunc("/inventory/customer_order/details/cancel", inventory_customer_order.DCancelHandler)
	mux.HandleFunc("/inventory/customer_order/details/Show_prev_price", inventory_customer_order.Show_prev_price)
	mux.HandleFunc("/inventory/customer_order/details/email", inventory_customer_order.EmailHandler)
	mux.HandleFunc("/inventory/customer_order/details/DVoidHandler", inventory_customer_order.DVoidHandler)

	mux.HandleFunc("/inventory/customer_bill", inventory_customer_bill.HListHandler)
	mux.HandleFunc("/inventory/customer_bill/HaddHandler", inventory_customer_bill.HAddHandler)
	mux.HandleFunc("/inventory/customer_bill/HEditHandler", inventory_customer_bill.HEditHandler)
	mux.HandleFunc("/inventory/customer_bill/HDeleteHandler", inventory_customer_bill.HDeleteHandler)
	mux.HandleFunc("/inventory/customer_bill/HaddTagHandler", inventory_customer_order.HAddTagHandler)
	mux.HandleFunc("/inventory/customer_bill/details", inventory_customer_bill.DListHandler)
	mux.HandleFunc("/inventory/customer_bill/details/DaddHandler", inventory_customer_bill.DAddHandler)
	mux.HandleFunc("/inventory/customer_bill/details/DEditHandler", inventory_customer_bill.DEditHandler)
	mux.HandleFunc("/inventory/customer_bill/details/DDeleteHandler", inventory_customer_bill.DDeleteHandler)
	mux.HandleFunc("/inventory/customer_bill/details/Ajax/Ajax_getItem", inventory_customer_bill.Ajax_getItem)
	mux.HandleFunc("/inventory/customer_bill/details/Ajax/Ajax_getItem_isbarcode", inventory_customer_bill.Ajax_getItem_isbarcode)
	mux.HandleFunc("/inventory/customer_bill/details/post", inventory_customer_bill.DPostHandler)
	mux.HandleFunc("/inventory/customer_bill/details/cancel", inventory_customer_bill.DCancelHandler)

	mux.HandleFunc("/inventory/internal_order", inventory_internal_order.HListHandler)
	mux.HandleFunc("/inventory/internal_order/HaddHandler", inventory_internal_order.HAddHandler)
	mux.HandleFunc("/inventory/internal_order/HEditHandler", inventory_internal_order.HEditHandler)
	mux.HandleFunc("/inventory/internal_order/HDeleteHandler", inventory_internal_order.HDeleteHandler)
	mux.HandleFunc("/inventory/internal_order/HaddTagHandler", inventory_internal_order.HAddTagHandler)
	mux.HandleFunc("/inventory/internal_order/details", inventory_internal_order.DListHandler)
	mux.HandleFunc("/inventory/internal_order/details/DaddHandler", inventory_internal_order.DAddHandler)
	mux.HandleFunc("/inventory/internal_order/details/DEditHandler", inventory_internal_order.DEditHandler)
	mux.HandleFunc("/inventory/internal_order/details/DDeleteHandler", inventory_internal_order.DDeleteHandler)
	mux.HandleFunc("/inventory/internal_order/details/Ajax/Ajax_getItem", inventory_internal_order.Ajax_getItem)
	mux.HandleFunc("/inventory/internal_order/details/Ajax/Ajax_getItem_isbarcode", inventory_internal_order.Ajax_getItem_isbarcode)
	mux.HandleFunc("/inventory/internal_order/details/post", inventory_internal_order.DPostHandler)
	mux.HandleFunc("/inventory/internal_order/details/cancel", inventory_internal_order.DCancelHandler)
	mux.HandleFunc("/inventory/internal_order/details/allocation", inventory_internal_order.AListHandler)
	mux.HandleFunc("/inventory/internal_order/details/allocation/DaddHandler", inventory_internal_order.AAddHandler)
	mux.HandleFunc("/inventory/internal_order/details/allocation/DEditHandler", inventory_internal_order.AEditHandler)
	mux.HandleFunc("/inventory/internal_order/details/allocation/DDeleteHandler", inventory_internal_order.ADeleteHandler)

	mux.HandleFunc("/inventory/internal_issuance", inventory_internal_issuance.HListHandler)
	mux.HandleFunc("/inventory/internal_issuance/HaddHandler", inventory_internal_issuance.HAddHandler)
	mux.HandleFunc("/inventory/internal_issuance/HEditHandler", inventory_internal_issuance.HEditHandler)
	mux.HandleFunc("/inventory/internal_issuance/HDeleteHandler", inventory_internal_issuance.HDeleteHandler)
	mux.HandleFunc("/inventory/internal_issuance/HaddTagHandler", inventory_internal_issuance.HAddTagHandler)
	mux.HandleFunc("/inventory/internal_issuance/details", inventory_internal_issuance.DListHandler)
	mux.HandleFunc("/inventory/internal_issuance/details/DaddHandler", inventory_internal_issuance.DAddHandler)
	mux.HandleFunc("/inventory/internal_issuance/details/DEditHandler", inventory_internal_issuance.DEditHandler)
	mux.HandleFunc("/inventory/internal_issuance/details/DDeleteHandler", inventory_internal_issuance.DDeleteHandler)
	mux.HandleFunc("/inventory/internal_issuance/details/Ajax/Ajax_getItem", inventory_internal_issuance.Ajax_getItem)
	mux.HandleFunc("/inventory/internal_issuance/details/Ajax/Ajax_getItem_isbarcode", inventory_internal_issuance.Ajax_getItem_isbarcode)
	mux.HandleFunc("/inventory/internal_issuance/details/post", inventory_internal_issuance.DPostHandler)
	mux.HandleFunc("/inventory/internal_issuance/details/cancel", inventory_internal_issuance.DCancelHandler)
	mux.HandleFunc("/inventory/internal_issuance/details/DownloadIntOrd", inventory_internal_issuance.Download_IntOrdHandler)
	mux.HandleFunc("/inventory/internal_issuance/details/allocation", inventory_internal_issuance.AListHandler)
	mux.HandleFunc("/inventory/internal_issuance/details/allocation/DaddHandler", inventory_internal_issuance.AAddHandler)
	mux.HandleFunc("/inventory/internal_issuance/details/allocation/DEditHandler", inventory_internal_issuance.AEditHandler)
	mux.HandleFunc("/inventory/internal_issuance/details/allocation/DDeleteHandler", inventory_internal_issuance.ADeleteHandler)

	mux.HandleFunc("/inventory/internal_return", inventory_internal_return.HListHandler)
	mux.HandleFunc("/inventory/internal_return/HaddHandler", inventory_internal_return.HAddHandler)
	mux.HandleFunc("/inventory/internal_return/HEditHandler", inventory_internal_return.HEditHandler)
	mux.HandleFunc("/inventory/internal_return/HDeleteHandler", inventory_internal_return.HDeleteHandler)
	mux.HandleFunc("/inventory/internal_return/HaddTagHandler", inventory_internal_return.HAddTagHandler)
	mux.HandleFunc("/inventory/internal_return/details", inventory_internal_return.DListHandler)
	mux.HandleFunc("/inventory/internal_return/details/DaddHandler", inventory_internal_return.DAddHandler)
	mux.HandleFunc("/inventory/internal_return/details/DEditHandler", inventory_internal_return.DEditHandler)
	mux.HandleFunc("/inventory/internal_return/details/DDeleteHandler", inventory_internal_return.DDeleteHandler)
	mux.HandleFunc("/inventory/internal_return/details/Ajax/Ajax_getItem", inventory_internal_return.Ajax_getItem)
	mux.HandleFunc("/inventory/internal_return/details/post", inventory_internal_return.DPostHandler)
	mux.HandleFunc("/inventory/internal_return/details/cancel", inventory_internal_return.DCancelHandler)

	mux.HandleFunc("/inventory/reports/order_per_area", inventory_order_per_area.HListHandler)
	mux.HandleFunc("/inventory/reports/order_per_area_v2", inventory_order_per_area_v2.HListHandler)

	mux.HandleFunc("/inventory/item_return_pairing", inventory_internal_return.DCancelHandler)

	mux.HandleFunc("/inventory/reports/item_movement/viewer", inventory_item_movement.ViewerHandler)
	mux.HandleFunc("/inventory/reports/item_movement/ajaxviewer", inventory_item_movement.AjaxViewerHandler)
	mux.HandleFunc("/inventory/reports/item_movement/ajaxviewer_details", inventory_item_movement.AjaxViewerDetailsHandler)
	mux.HandleFunc("/inventory/reports/item_movement/Ajax/Ajax_getItem", inventory_item_movement.Ajax_getItem)
	mux.HandleFunc("/report/item_movement", inventory_item_movement.HListHandler)

	mux.HandleFunc("/Show_previous_price", Show_prev_price)
	mux.HandleFunc("/Show_previous_price", Show_prev_price)

	mux.HandleFunc("/api/ajax/inventory/getlatestlocation/{module}/{org}/{key1}", helper.GetLatestLocation)
	mux.HandleFunc("/api/ajax/{module}/{sub_module}/{org}/{key1}", helper.DefaultApiFunc)
	mux.HandleFunc("/api/ajax/api/{transaction}/{version}", helper.DefaultApiFunc_all)

	mux.HandleFunc("/inventory/helper/application", helper_application.DListHandler)
	mux.HandleFunc("/inventory/helper/application/DaddHandler", helper_application.DAddHandler)
	mux.HandleFunc("/inventory/helper/application/DEditHandler", helper_application.DEditHandler)
	mux.HandleFunc("/inventory/helper/application/DDeleteHandler", helper_application.DDeleteHandler)
	mux.HandleFunc("/inventory/helper/application/post", helper_application.DPostHandler)
	mux.HandleFunc("/inventory/helper/application/cancel", helper_application.DCancelHandler)
	mux.HandleFunc("/inventory/helper/application/cancelpay", helper_application.DCancelPayHandler)

	mux.HandleFunc("/json_func", Raw_func)

	mux.HandleFunc("/test_upload", test_uploadHAndler)
	mux.HandleFunc("/do_test_upload", do_test_uploadHAndler)
	mux.HandleFunc("/do_image_upload", do_image_uploadHAndler)

	mux.HandleFunc("/inventory/auto_complete", auto_complete)

	mux.HandleFunc("/sync_server_side", Sync_server_side)

	//file server
	var dir string = `./metronic_3.1.2`
	mux.PathPrefix("/src/metronic_3.1.2/").Handler(http.StripPrefix("/src/metronic_3.1.2/", http.FileServer(http.Dir(dir))))
	var dir_img string = `./img_uploads`
	mux.PathPrefix("/src/img_uploads/").Handler(http.StripPrefix("/src/img_uploads/", http.FileServer(http.Dir(dir_img))))

	//	mux.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("metronic_3.1.2"))))

	//	http.ListenAndServe("",context.ClearHandler(mux)) // listen on all interfaces on port 8080
	str_port := ":6565"
	//str_port := ":8181"
	fmt.Println("starting server on http://localhost" + str_port + "/")
	//exec.Command("open", "http://localhost").Start()
	//exec.Command("cmd", "/c", "start", "http://localhost"+str_port).Start()
	//go config.LogHandler("starting server on http://localhost/"+str_port)

	http.ListenAndServe(str_port, context.ClearHandler(mux)) // listen on all interfaces on port 8080

}

func auto_complete(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("query")
	tag := r.URL.Query().Get("tag")
	org := r.URL.Query().Get("org")
	sql := `sis_itemtags_get 1, ` + org + `, 0,` + tag + ` ,'` + search + `'`
	arr_data := datatables.DataList(sql)
	js, err := json.Marshal(arr_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func Raw_func(w http.ResponseWriter, r *http.Request) {
	var sql string
	if r.Method == "GET" {
		sql = r.URL.Query().Get("qdata")
	} else {
		r.ParseForm()
		for _, val := range r.Form {
			fmt.Println(val[0])
			sql = val[0]
		}

	}
	arr_data := datatables.DataList(sql)
	js, err := json.Marshal(arr_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func Get_duplicate_dev_id(w http.ResponseWriter, r *http.Request) {
	tconf := make(map[string]string)
	tconf["Device_id"] = r.URL.Query().Get("device_id")
	type LBR_MemAss_CheckSameDeviceID struct {
		Assignment interface{}
		Member     interface{}
	}

	fmt.Println(`LBR_MemAss_CheckSameDeviceID 1, 2,'` + tconf["Device_id"] + `'`)

	//LBR_LogHdr_row ,err, _,_ := config.Ap_sql(`lbr_rep_logpair 3, 2, `+tconf["member"]+` , `+tconf["lbr_assign"]+`, '`+tconf["datefrom"]+`', '`+tconf["dateto"] +`'`,1)
	sql_row, err, _, _ := config.Ap_sql(`LBR_MemAss_CheckSameDeviceID 1, 2,'`+tconf["Device_id"]+`'`, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err.Error())
	}
	LBR_MemAss_CheckSameDeviceID_data := []LBR_MemAss_CheckSameDeviceID{}
	for sql_row.Next() {
		var r LBR_MemAss_CheckSameDeviceID
		//err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
		err = sql_row.Scan(&r.Assignment, &r.Member)
		if err != nil {
			panic(err.Error())
		}
		post2 := LBR_MemAss_CheckSameDeviceID{
			r.Assignment,
			r.Member,
		}
		LBR_MemAss_CheckSameDeviceID_data = append(LBR_MemAss_CheckSameDeviceID_data, post2)
	}
	type Htdata struct {
		LBR_MemAss_CheckSameDeviceID_data []LBR_MemAss_CheckSameDeviceID
	}
	var_htdata := Htdata{LBR_MemAss_CheckSameDeviceID_data}

	tmpl := template.New("get_duplicate_device_id.html").Funcs(config.FuncMap)

	//var err error
	//if tmpl, err = tmpl.ParseFiles("hris/assignments/assignment_list.html"); err != nil {
	if tmpl, err = tmpl.ParseFiles("helper/get_duplicate_device_id/get_duplicate_device_id.html"); err != nil {
		fmt.Println(err)
	}
	//err1 := tmpl.Execute(w,&TblConf{header,ajaxURLdata,tconf})
	err1 := tmpl.Execute(w, var_htdata)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}

}

func Show_prev_price(w http.ResponseWriter, r *http.Request) {
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	customer := r.URL.Query().Get("customer")
	location := r.URL.Query().Get("location")
	item := r.URL.Query().Get("item")
	if item == `` {
		item = `0`
	}
	sql := `SIS_CusSal_ShowLastRecords 2, ` + str_OrgID + `, ` + customer + `, ` + item + `, ` + location + ``
	fmt.Println(sql)
	arr_data := datatables.DataList(sql)
	tmpl := template.New("show_prev.html").Funcs(config.FuncMap)
	var err error
	if tmpl, err = tmpl.ParseFiles("helper/inven_show_previous_price/show_prev.html"); err != nil {
		fmt.Println(err)
	}
	err1 := tmpl.Execute(w, arr_data)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func test_uploadHAndler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	tmpl := template.New("test_upload.html")
	var err error
	if tmpl, err = tmpl.ParseFiles("test_upload.html"); err != nil {
		fmt.Println(err)
		go config.LogHandler(fmt.Sprint(err))
	}
	err1 := tmpl.Execute(w, nil)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		go config.LogHandler(fmt.Sprint(err1))
	}
}

func do_test_uploadHAndler(w http.ResponseWriter, r *http.Request) {
	//get the multipart reader for the request.
	login.Session_renew(w, r)
	session_username, _ := login.Get_account_info(r)
	getid := r.URL.Query().Get("hid")
	device := r.URL.Query().Get("device")

	fmt.Println("select CONVERT(VARCHAR(10),date_from,20 ) ,   CONVERT(VARCHAR(10),date_to,20 )  from LBR_LogHdr where id=" + getid)
	arr_from_to := datatables.Data_row("select CONVERT(VARCHAR(10),date_from,20 ) ,   CONVERT(VARCHAR(10),date_to,20 )  from LBR_LogHdr where id=" + getid)
	fmt.Println(arr_from_to[0])

	fmt.Println(arr_from_to[1])
	fmt.Println("convert time to int")

	//fmt.Println(`logtime`,logtime[0:10])

	//trim_logtime := strings.Replace(logtime[0:10], "-", "", -1)

	trim_timefrom_str := strings.Replace(arr_from_to[0][0:10], "-", "", -1)
	trim_timefrom_int, _ := strconv.Atoi(trim_timefrom_str)
	trim_timeto_str := strings.Replace(arr_from_to[1][0:10], "-", "", -1)
	trim_timeto_int, _ := strconv.Atoi(trim_timeto_str)
	fmt.Println(trim_timefrom_str)
	fmt.Println(trim_timeto_str)

	fmt.Println("running to int here")
	fmt.Println(trim_timefrom_int)
	fmt.Println(trim_timeto_int)

	fmt.Println("end of conversion----------")
	/*const longForm = "Jan 2 2006"
	from_time,_:= time.Parse(longForm,arr_from_to[0])
	to_time,_:= time.Parse(longForm,arr_from_to[1])
	fmt.Println("converted time on from to ")
	fmt.Println(from_time)
	fmt.Println(to_time)
	fmt.Println("----------end converted time on from")*/
	fmt.Println("device :", device)
	fmt.Println("hid :", getid)
	go config.LogHandler(fmt.Sprint("hid :", getid))
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type Translog struct {
		Txtline int
		MsgID   int
		Msg     interface{}
	}
	//fmt.Println(translog_arr)
	fmt.Println("Upload successful.")
	go config.LogHandler("Upload successfully")
	t0 := time.Now().Local()
	t0_newFormat := `Start Time : ` + t0.Format("15:04:05")

	translog_arr := []Translog{}
	translog_arr = append(translog_arr, Translog{0, 0, t0_newFormat})
	var txline int = 1
	var wg sync.WaitGroup
	//copy each part to destination.
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		/*if( strings.Contains(part.FileName(), ".dat")){
		 				fmt.Println(".dat here")
		 				continue
					}else{
						fmt.Println("invalid file name here.")
						break
					}*/
		//if part.FileName() is empty, skip this iteration.
		if part.FileName() == "" {
			continue
		}
		if strings.Contains(part.FileName(), ".dat") == false && strings.Contains(part.FileName(), ".DAT") && strings.Contains(part.FileName(), ".txt") && strings.Contains(part.FileName(), ".TXT") == false {
			translog_arr = append(translog_arr, Translog{0, 0, `Error: invalid filename.`})
			fmt.Println(`Error: invalid filename.`)
			continue
		}

		fmt.Println("print file name :", part.FileName())
		go config.LogHandler(fmt.Sprint("print file name :", part.FileName()))
		t := time.Now() //get current date time to make name of folder
		datetimefoldername := t.Format("20060102150405")
		//os.Mkdir("C:/Go/muse/uploads/"+datetimefoldername,0777)
		os.Mkdir("uploads/"+datetimefoldername, 0777)
		dst, err := os.Create("uploads/" + datetimefoldername + "/" + part.FileName())
		defer dst.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, part); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//on parsing data from file
		contents, err := ioutil.ReadFile("uploads/" + datetimefoldername + "/" + part.FileName())
		if err != nil {
			fmt.Println("error on reading file :", err)
			go config.LogHandler(fmt.Sprint("error on reading file :", err))
		}
		//fmt.Println(string(contents));

		lines := strings.Split(string(contents), "\n")
		len_lines := len(lines)
		str_len_lines := strconv.Itoa(len_lines)

		//lines,_:=strings.ReadString('\n')
		//fmt.Println("data sa arr 0", lines[0][10:len(lines[0])-9])
		for key, _ := range lines {
			if device == "C7" {
				if key > 0 && key < (len(lines)-1) {
					//logid:=lines[key][8:len(lines[key])-29]
					arr_line := strings.Split(string(lines[key]), "\t")
					//fmt.Println("wetwew arr_line")
					fmt.Println(`len here:`, len(arr_line))
					fmt.Println(arr_line[2])
					fmt.Println(arr_line[len(arr_line)-1])
					logid := strings.Trim(arr_line[2], " ")
					logtime := strings.TrimSpace(arr_line[len(arr_line)-1])
					if logid == "" {
						break
					}
					wg.Add(1)
					func() {

						fmt.Println("LBR_LogDtl_Save 'UploadBioText','" + session_username + "', 0, " + getid + ", 0, 0, 0, '" + logid + "', '" + logtime + "'")
						go config.LogHandler("LBR_LogDtl_Save 'UploadBioText','" + session_username + "', 0, " + getid + ", 0, 0, 0, '" + logid + "', '" + logtime + "'")
						var rr Translog
						sqlstr := "LBR_LogDtl_Save 'UploadBioText','" + session_username + "', 0, " + getid + ", 0, 0, 0, '" + logid + "', '" + logtime + "'"
						rowret, err, _, _ := config.Ap_sql(sqlstr, 1)
						//_ ,err, _,_ := config.Ap_sql("LBR_LogDtl_Save 'UploadBioText','"+login.Uname+"', 0, "+getid+", 0, 0, 0, '"+logid+"', '"+logtime+"'",3)
						if err != nil {
							//fmt.Println(err.Error)
							//panic(err.Error())
						}
						for rowret.Next() {
							err = rowret.Scan(&rr.MsgID, &rr.Msg)
							if err != nil {
								panic(err.Error())
							}

							translog_arr = append(translog_arr, Translog{txline, rr.MsgID, rr.Msg})
							go upload_info(w, r, txline, str_len_lines, session_username)

						}
						defer wg.Done()
					}()
					wg.Wait()
					txline = txline + 1
					fmt.Println(arr_line[2])
					fmt.Println(arr_line[len(arr_line)-1])
				}
			} else if device == "N1" {
				fmt.Println("n1 here..")
				if key < (len(lines) - 1) {
					arr_line := strings.Split(string(lines[key]), "\t")
					logid := strings.Trim(arr_line[0], " ")
					logtime := strings.TrimSpace(arr_line[1])

					fmt.Println(`logid`, logid)
					fmt.Println(`logtime`, logtime[0:10])

					trim_logtime := strings.Replace(logtime[0:10], "-", "", -1)
					fmt.Println(`trimmed logtime`, trim_logtime)
					trim_txt_log, _ := strconv.Atoi(trim_logtime)
					fform := "20060102"
					con_logs, terr := time.Parse(fform, trim_logtime)
					fmt.Println("converted timelogs on loop ")
					fmt.Println(con_logs)
					fmt.Println(terr)

					fmt.Println("----------end timelogs timelogs on loop")
					fmt.Println("device :", device)
					fmt.Println("hid :", getid)

					if logid == "" {
						break
					}

					//wg.Add(1)
					// func() {

					if trim_timefrom_int <= trim_txt_log && trim_timeto_int >= trim_txt_log {
						fmt.Println("good...")

						fmt.Println("LBR_LogDtl_Save 'UploadBioText','" + session_username + "', 0, " + getid + ", 0, 0, 0, '" + logid + "', '" + logtime + "'")
						//time.Sleep(100 * time.Millisecond)
						arr_data := datatables.Data_row("LBR_LogDtl_Save 'UploadBioText','" + session_username + "', 0, " + getid + ", 0, 0, 0, '" + logid + "', '" + logtime + "'")
						fmt.Println(arr_data)
						arr_dataid_int, _ := strconv.Atoi(arr_data[0])
						translog_arr = append(translog_arr, Translog{txline, arr_dataid_int, arr_data[1]})
					}
					upload_info(w, r, txline, str_len_lines, session_username)

					/*fmt.Println("LBR_LogDtl_Save 'UploadBioText','"+session_username+"', 0, "+getid+", 0, 0, 0, '"+logid+"', '"+logtime+"'")
					go config.LogHandler("LBR_LogDtl_Save 'UploadBioText','"+session_username+"', 0, "+getid+", 0, 0, 0, '"+logid+"', '"+logtime+"'")
					var rr Translog
					sqlstr :="LBR_LogDtl_Save 'UploadBioText','"+session_username+"', 0, "+getid+", 0, 0, 0, '"+logid+"', '"+logtime+"'"
					rowret ,err, _,_ := config.Ap_sql(sqlstr,1)
					//_ ,err, _,_ := config.Ap_sql("LBR_LogDtl_Save 'UploadBioText','"+login.Uname+"', 0, "+getid+", 0, 0, 0, '"+logid+"', '"+logtime+"'",3)
					if err != nil {
					        //fmt.Println(err.Error)
					        //panic(err.Error())
					}
					for rowret.Next() {
					    err = rowret.Scan(&rr.MsgID,&rr.Msg)
					    if err != nil {
					        panic(err.Error())
					    }

					translog_arr = append(translog_arr,Translog{txline,rr.MsgID,rr.Msg})
					go upload_info(w,r,txline,str_len_lines,session_username)

					}*/
					//	defer wg.Done()
					//}()
					//wg.Wait()
					//translog_arr = append(translog_arr,Translog{txline,arr_data[0],arr_data[1]})

					txline = txline + 1
					fmt.Println(arr_line[2])
					fmt.Println(arr_line[len(arr_line)-1])

				}

			} else if device == "VF300" {
				fmt.Println("VF300 DEVICE UPLOADING")
				if key < (len(lines) - 1) {
					arr_line := strings.Split(string(lines[key]), "\t")
					logid := strings.Trim(arr_line[0], " ")
					logtime := strings.TrimSpace(arr_line[1])
					fmt.Println(`logid`, logid)
					fmt.Println(`logtime`, logtime)
					if logid == "" {
						break
					}

					//fmt.Println("LBR_LogDtl_Save 'UploadBioText','"+session_username+"', 0, "+getid+", 0, 0, 0, '"+logid+"', '"+logtime+"'")

					wg.Add(1)
					func() {

						fmt.Println("LBR_LogDtl_Save 'UploadBioText','" + session_username + "', 0, " + getid + ", 0, 0, 0, '" + logid + "', '" + logtime + "'")
						go config.LogHandler("LBR_LogDtl_Save 'UploadBioText','" + session_username + "', 0, " + getid + ", 0, 0, 0, '" + logid + "', '" + logtime + "'")
						var rr Translog
						sqlstr := "LBR_LogDtl_Save 'UploadBioText','" + session_username + "', 0, " + getid + ", 0, 0, 0, '" + logid + "', '" + logtime + "'"
						rowret, err, _, _ := config.Ap_sql(sqlstr, 1)
						//_ ,err, _,_ := config.Ap_sql("LBR_LogDtl_Save 'UploadBioText','"+login.Uname+"', 0, "+getid+", 0, 0, 0, '"+logid+"', '"+logtime+"'",3)
						if err != nil {
							//fmt.Println(err.Error)
							//panic(err.Error())
						}
						for rowret.Next() {
							err = rowret.Scan(&rr.MsgID, &rr.Msg)
							if err != nil {
								panic(err.Error())
							}

							translog_arr = append(translog_arr, Translog{txline, rr.MsgID, rr.Msg})
							go upload_info(w, r, txline, str_len_lines, session_username)

						}
						defer wg.Done()
					}()
					wg.Wait()
					txline = txline + 1
					fmt.Println(arr_line[2])
					fmt.Println(arr_line[len(arr_line)-1])
				}

			} else {
				if len(lines[key]) == 38 {
					logid := lines[key][0 : len(lines[key])-29] //Need to manually add end of string
					logid = strings.Trim(logid, " ")
					logtime := lines[key][10 : len(lines[key])-9]
					//fmt.Println("Log id : " ,logid)
					//fmt.Println("log time : " ,logtime)
					if logid == "" {
						break
					}
					wg.Add(1)
					func() {

						fmt.Println("LBR_LogDtl_Save 'UploadBioText','" + session_username + "', 0, " + getid + ", 0, 0, 0, '" + logid + "', '" + logtime + "'")
						go config.LogHandler("LBR_LogDtl_Save 'UploadBioText','" + session_username + "', 0, " + getid + ", 0, 0, 0, '" + logid + "', '" + logtime + "'")
						var rr Translog
						sqlstr := "LBR_LogDtl_Save 'UploadBioText','" + session_username + "', 0, " + getid + ", 0, 0, 0, '" + logid + "', '" + logtime + "'"
						rowret, err, _, _ := config.Ap_sql(sqlstr, 1)
						//_ ,err, _,_ := config.Ap_sql("LBR_LogDtl_Save 'UploadBioText','"+login.Uname+"', 0, "+getid+", 0, 0, 0, '"+logid+"', '"+logtime+"'",3)
						if err != nil {
							//fmt.Println(err.Error)
							//panic(err.Error())
						}
						for rowret.Next() {
							err = rowret.Scan(&rr.MsgID, &rr.Msg)
							if err != nil {
								panic(err.Error())
							}

							translog_arr = append(translog_arr, Translog{txline, rr.MsgID, rr.Msg})
							go upload_info(w, r, txline, str_len_lines, session_username)

						}
						defer wg.Done()
					}()
					wg.Wait()
				} else {
					translog_arr = append(translog_arr, Translog{txline, 1, `Error: invalid text line.`})
				}
				txline = txline + 1
			}
		}

	}
	//display success message.
	//display(w, "upload", "Upload successful.")

	//fmt.Println(translog_arr)
	fmt.Println("Upload successful.")
	go config.LogHandler("Upload successful")
	t1 := time.Now().Local()
	t1_newFormat := `End Time : ` + t1.Format("15:04:05")

	translog_arr = append(translog_arr, Translog{0, 0, t1_newFormat})
	js, err := json.Marshal(translog_arr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	//wg.Wait()
}

func do_image_uploadHAndler(w http.ResponseWriter, r *http.Request) {

	getid := r.URL.Query().Get("hid")
	device := r.URL.Query().Get("device")
	var tmp_file string
	t := time.Now() //get current date time to make name of folder
	datetimefoldername := t.Format("20060102150405")

	fmt.Println("hid :", getid)
	fmt.Println("device :", device)
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//copy each part to destination.
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if part.FileName() == "" {
			continue
		}

		fmt.Println("print file name :", part.FileName())

		//os.Mkdir("C:/Go/muse/uploads/"+datetimefoldername,0777)
		mkdir_1st_err := os.Mkdir("uploads/"+datetimefoldername, 0777)
		fmt.Println("mkdir1st err:", mkdir_1st_err)
		tmp_file = part.FileName()
		dst, err := os.Create("uploads/" + datetimefoldername + "/" + part.FileName())
		defer dst.Close()

		mkdir_err := os.Mkdir("uploads/rename/"+datetimefoldername, 0777)
		fmt.Println("mkdir err:", mkdir_err)
		/*errr := os.Rename("uploads/"+datetimefoldername+"/" + tmp_file, "uploads/rename/"+datetimefoldername+"/" + `test_rename`)
		fmt.Println("print error here.")*/

		errr := Copy("uploads/"+datetimefoldername+"/"+part.FileName(), "uploads/rename/"+datetimefoldername)
		fmt.Println(errr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, part); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	// open "test.jpg"
	file, err := os.Open("uploads/" + datetimefoldername + "/" + tmp_file)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(171, 0, img, resize.Lanczos3)

	mkdir_err := os.Mkdir("img_uploads/size_171/"+datetimefoldername, 0777)
	fmt.Println("mkdir err:", mkdir_err)
	if getid == `` {
		getid = `temp_img`
	}

	out, err := os.Create("img_uploads/size_171/" + datetimefoldername + "/" + getid + ".jpg")
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	// base64 conversion

	base64_img := helper.Base64_image_encoding("img_uploads/size_171/" + datetimefoldername + "/" + getid + ".jpg")

	//end base64 conversion

	type Return_ajax_data struct {
		Base64_image  string
		Image_ser_url string
	}

	// w.Write(fmt.Sprintf(base64_img)))

	ret_Data := Return_ajax_data{base64_img, "uploads/" + datetimefoldername + "/" + tmp_file}

	js, err := json.Marshal(ret_Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func Copy(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}

type User_upload_list struct {
	Username   string
	Desciption string
}

var User_ups_status_list []User_upload_list //save to memory,
func upload_info(w http.ResponseWriter, r *http.Request, lineno int, line_tot string, Username string) {
	str_line := strconv.Itoa(lineno)
	pollstat := str_line + ` / ` + line_tot
	var temp_var_User_ups_status_list []User_upload_list
	for _, val := range User_ups_status_list {
		if val.Username != Username {
			temp_var_User_ups_status_list = append(User_ups_status_list, val)
		}
	}
	User_ups_status_list = temp_var_User_ups_status_list
	User_ups_status_list = append(User_ups_status_list, User_upload_list{Username, pollstat})
}
func AjaxPollingeHandler(w http.ResponseWriter, r *http.Request) {
	go login.Session_renew(w, r)
	session_username, _ := login.Get_account_info(r)
	var uploading_status string = ``
	for _, val := range User_ups_status_list {
		if val.Username == session_username {
			uploading_status = val.Desciption
		}
		//	fmt.Println(`key:`, key , `val:`,val)
	}

	js, err := json.Marshal(uploading_status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func AjaxValidateHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	r.ParseForm()
	for n, w := range r.Form {
		fmt.Println("name:", n, " value:", w)
		go config.LogHandler(fmt.Sprint("name:", n, " value:", w))
	}
	type Profile struct {
		Message string // error status output or validated data
		Stats   string //pass or fail
		//Hobbies []string
	}

	/*	draw": 2,
		recordsTotal": 57,
		recordsFiltered": 57,

	*/
	var sqstring string
	//var errstring string
	if r.Form["field"][0] == "name" {
		//sqstring ="sp_validateNAme '"+r.Form["lname"][0]+"','"+r.Form["mname"][0]+"','"+r.Form["fname"][0]+"','"+r.Form["bday"][0]+"'"
		sqstring = "exec Member_Validate 'Add', 0, '', '" + r.Form["lname"][0] + "','" + r.Form["mname"][0] + "', '" + r.Form["fname"][0] + "',  '" + r.Form["bday"][0] + "'"

	} else if r.Form["field"][0] == "TIN" {
		sqstring = "exec Member_Validate 'Add', 0, '', '', '', '',  '', '" + r.Form["value"][0] + "'"
	} else if r.Form["field"][0] == "SSS" {
		sqstring = "exec Member_Validate 'Add', 0, '', '', '', '',  '', '','" + r.Form["value"][0] + "'"
	} else if r.Form["field"][0] == "GSIS" {
		sqstring = "exec Member_Validate 'Add', 0, '', '', '', '',  '', '','','" + r.Form["value"][0] + "'"
	} else if r.Form["field"][0] == "Pagibig" {
		sqstring = "exec Member_Validate 'Add', 0, '', '', '', '',  '', '','','','" + r.Form["value"][0] + "'"
	} else if r.Form["field"][0] == "PHealth" {
		sqstring = "exec Member_Validate 'Add', 0, '', '', '', '',  '', '','','','','" + r.Form["value"][0] + "'"
	} else {
		sqstring = "select username from member where " + r.Form["field"][0] + "='" + r.Form["value"][0] + "'"
	}
	rows, err, _, _ := config.Ap_sql(sqstring, 1)
	var rowConter int
	var ErrMsg string
	ErrMsg = ``
	for rows.Next() {
		err = rows.Scan(&ErrMsg)
		if err != nil {
			panic(err.Error())
		}

		rowConter = rowConter + 1
	}
	if ErrMsg != `` {
		profile := Profile{"warning " + r.Form["field"][0] + ` ` + ErrMsg, "1"}
		js, err := json.Marshal(profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		profile := Profile{"", "0"}
		js, err := json.Marshal(profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func AjaxGetPostalCodeHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)
	getq := r.URL.Query().Get("q")
	getid := r.URL.Query().Get("id")
	getq = getq + getid
	type Inner_item struct {
		Id       int    `json:"id"`
		Fullname string `json:"full_name"`
	}
	type JsondataEmpty struct {
		Total_count        int          `json:"total_count"`
		Incomplete_results bool         `json:"incomplete_results"`
		Items              []Inner_item `json:"items"`
	}
	var count int
	var strtype string
	if getid != "" {
		strtype = "1"
	} else {
		strtype = "2"
	}
	rows, err, _, _ := config.Ap_sql("exec postcode_get "+strtype+",  '"+getq+"'", 1)
	if err != nil {
		fmt.Println("db error:", err)
	}
	InItem := []Inner_item{}
	var InItemRow Inner_item
	for rows.Next() {
		var r Inner_item
		err = rows.Scan(&r.Id, &r.Fullname)
		if err != nil {
			panic(err.Error())
		}
		count = count + 1
		post := Inner_item{
			r.Id,
			r.Fullname,
		}
		InItem = append(InItem, post)
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

func AjaxGetMemberHandler(w http.ResponseWriter, r *http.Request) {
	login.Session_renew(w, r)

	getq := r.URL.Query().Get("q")
	getid := r.URL.Query().Get("id")
	isbaseonassignid := r.URL.Query().Get("isbaseonassignid")
	assignmentID := r.URL.Query().Get("assignmentID")
	isbaseOrg := r.URL.Query().Get("isbaseOrg")
	Org_id := login.Get_session_org_id(r)
	str_OrgID := strconv.Itoa(Org_id)
	if isbaseOrg == `1` {
		assignmentID = str_OrgID
	}

	sp := r.URL.Query().Get("sp")

	fmt.Println(`is base sa assignID`, isbaseonassignid)
	fmt.Println(` assignID`, assignmentID)

	getq = getq + getid
	type Inner_item struct {
		Id       int    `json:"id"`
		Fullname string `json:"full_name"`
	}
	type JsondataEmpty struct {
		Total_count        int          `json:"total_count"`
		Incomplete_results bool         `json:"incomplete_results"`
		Items              []Inner_item `json:"items"`
	}
	var count int
	var strtype string
	if getid != "" {
		strtype = "1"
	} else {
		strtype = "2"
	}
	//`exec lbr_memass_search 1, 395 , '' ,'',''`
	var sqlstr string
	if isbaseonassignid == `1` {
		//sqlstr = `exec lbr_memass_search 1, `+assignmentID+` , '' ,'',''`
		sqlstr = `exec LBR_Member_Get 3,'` + getq + `',  ` + assignmentID
		if sp == `orgmem_get` {

			sqlstr = `exec orgmem_get 0,` + str_OrgID + `,0,'` + getq + `'`
		}
		//orgmem_get 0, 3, 0, ''
	} else {
		sqlstr = "exec LBR_Member_Get " + strtype + ",  '" + getq + "'"
		if sp == `orgmem_get` {
			Org_id := login.Get_session_org_id(r)
			str_OrgID := strconv.Itoa(Org_id)
			sqlstr = `exec orgmem_get 0,` + str_OrgID + `,0,'` + getq + `'`
		}
	}
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
		if err != nil {
			panic(err.Error())
		}
		count = count + 1
		post := Inner_item{
			r.Id,
			r.Fullname,
		}
		InItem = append(InItem, post)
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
func Ajax_lbr_lograw_get(w http.ResponseWriter, r *http.Request) {
	member := r.URL.Query().Get("member")
	lbr_assign := r.URL.Query().Get("lbr_assign")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	sp := r.URL.Query().Get("sp")
	fmt.Println("sp no", sp)
	if sp == "sp1" {

		type Struct_lbr_lograw_get struct {
			Id          string
			LogDate     string
			Logtime     string
			Lbr_logpair string
		}
		var strID string
		var varlogdate time.Time
		var varlogtime time.Time
		var strLbr_logpair string
		fmt.Println(`lbr_lograw_get 1, 2, ` + member + `, ` + lbr_assign + `, '` + from + `' , '` + to + `'`)
		sql_row, err, _, _ := config.Ap_sql(`lbr_lograw_get 1, 2, `+member+`, `+lbr_assign+`, '`+from+`' , '`+to+`'`, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err.Error())
		}
		lbr_lograw_get_data := []Struct_lbr_lograw_get{}
		for sql_row.Next() {
			//err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
			err = sql_row.Scan(&strID, &varlogdate, &varlogtime, &strLbr_logpair)
			if err != nil {
				panic(err.Error())
			}
			post2 := Struct_lbr_lograw_get{
				strID,
				varlogdate.Format("20060102"),
				varlogtime.Format("15:04:05"),
				strLbr_logpair,
			}
			lbr_lograw_get_data = append(lbr_lograw_get_data, post2)
		}

		fmt.Println(lbr_lograw_get_data)
		js, err := json.Marshal(lbr_lograw_get_data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else if sp == "sp2" {

		type Struct_lbr_logpair_get struct {
			Id         string
			Logdate    string
			Assignment string
			Login      string
			Logout     string
			Daycode    string
			Hr_reg     interface{}
			Hr_ot      interface{}
			Hr_nd      interface{}
		}
		var varlogdate time.Time
		var vartime time.Time
		var vartimeout time.Time
		//var varlogtime time.Time
		fmt.Println(`lbr_logpair_get 1,0, 2, ` + member + `, ` + lbr_assign + `, '` + from + `' , '` + to + `'`)
		sql_row, err, _, _ := config.Ap_sql(`lbr_logpair_get 1,0, 2, `+member+`, `+lbr_assign+`, '`+from+`' , '`+to+`'`, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err.Error())
		}
		lbr_logpair_get_data := []Struct_lbr_logpair_get{}
		for sql_row.Next() {
			var r Struct_lbr_logpair_get
			//err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
			err = sql_row.Scan(&r.Id, &varlogdate, &r.Assignment, &vartime, &vartimeout, &r.Daycode, &r.Hr_reg, &r.Hr_ot, &r.Hr_nd)
			if err != nil {
				panic(err.Error())
			}
			post2 := Struct_lbr_logpair_get{
				r.Id,
				varlogdate.Format("20060102"),
				r.Assignment,
				vartime.Format("15:04:05"),
				vartimeout.Format("15:04:05"),
				r.Daycode,
				config.Interfacetosting(r.Hr_reg),
				config.Interfacetosting(r.Hr_ot),
				config.Interfacetosting(r.Hr_nd),
			}

			lbr_logpair_get_data = append(lbr_logpair_get_data, post2)
		}

		fmt.Println(lbr_logpair_get_data)
		js, err := json.Marshal(lbr_logpair_get_data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	} else if sp == "sp3" {

		type Struct_lbr_OT_get struct {
			LogDate string
			Hours   string
		}
		var varlogdate time.Time
		var hrs string
		//var varlogtime time.Time
		fmt.Println(`lbr_OT_get 1, ` + member + `, ` + lbr_assign + `, '` + from + `' , '` + to + `'`)
		sql_row, err, _, _ := config.Ap_sql(`lbr_OT_get 1, `+member+`, `+lbr_assign+`, '`+from+`' , '`+to+`'`, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err.Error())
		}
		lbr_OT_get_data := []Struct_lbr_OT_get{}
		for sql_row.Next() {
			//err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)

			err = sql_row.Scan(&varlogdate, &hrs)
			if err != nil {
				panic(err.Error())
			}
			post2 := Struct_lbr_OT_get{
				varlogdate.Format("20060102"),
				hrs,
			}
			lbr_OT_get_data = append(lbr_OT_get_data, post2)
		}

		fmt.Println(lbr_OT_get_data)
		js, err := json.Marshal(lbr_OT_get_data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else if sp == "LBR_Assign_Get" { //LBR_LogPair_UpdateAssignment

		type LBR_LogPair_UpdateAssignment_struct struct {
			Id   string `json:"id"`
			Name string `json:"text"`
		}
		fmt.Println(`LBR_Assign_Get 3, ` + lbr_assign)
		sql_row, err, _, _ := config.Ap_sql(`LBR_Assign_Get 3,`+lbr_assign, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err.Error())
		}
		LBR_LogPair_UpdateAssignment_data := []LBR_LogPair_UpdateAssignment_struct{}
		for sql_row.Next() {
			//err = LBR_LogHdr_row.Scan(&r.Member_name, &r.Device_id,&r.LogDate,&r.Day_name,&r.DayCode,&r.Assignment,&r.In_Sched,&r.Out_Sched,&r.In_Time,&r.Out_Time,&r.Hr_total,&r.Hr_break,&r.Hr_reg,&r.Hr_ot,&r.Hr_nd,&r.Hr_ndot,&r.Min_graceperiod ,&r.Hr_late,&r.Hr_undertime,&r.colormered)
			var r LBR_LogPair_UpdateAssignment_struct
			err = sql_row.Scan(&r.Id, &r.Name)
			if err != nil {
				panic(err.Error())
			}
			post2 := LBR_LogPair_UpdateAssignment_struct{
				r.Id,
				r.Name,
			}
			LBR_LogPair_UpdateAssignment_data = append(LBR_LogPair_UpdateAssignment_data, post2)
		}

		fmt.Println(LBR_LogPair_UpdateAssignment_data)
		js, err := json.Marshal(LBR_LogPair_UpdateAssignment_data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

func Sync_server_side(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	branchID := r.Form["BranchID"][0]
	fmt.Println(`BranchID ni:`, branchID)
	type TEst_test struct {
		ID     int
		Status string
	}
	LBR_LogPair_UpdateAssignment_data := TEst_test{1001, `Success`}
	fmt.Println(LBR_LogPair_UpdateAssignment_data)
	js, err := json.Marshal(LBR_LogPair_UpdateAssignment_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
