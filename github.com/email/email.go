package email

import (
    "bytes"
    "fmt"
    "net/smtp"
    "runtime"
    "strings"
    "text/template"
)

/*func main() {
    SendEmail(
        "smtp.gmail.com",
        587,
        "lonmars143",
        "nolram321",
        []string{"marriechairo@gmail.com"},
        "Gwapo ang nag send ",
        "<html><body>Go lang automation email by Mars Cute.</body></html>Exception 1 ")
    }*/
func catchPanic(err *error, functionName string) {
    if r := recover(); r != nil {
        fmt.Printf("%s : PANIC Defered : %v\n", functionName, r)

        // Capture the stack trace
        buf := make([]byte, 10000)
        runtime.Stack(buf, false)

        fmt.Printf("%s : Stack Trace : %s", functionName, string(buf))

        if err != nil {
            *err = fmt.Errorf("%v", r)
        }
    } else if err != nil && *err != nil {
        fmt.Printf("%s : ERROR : %v\n", functionName, *err)

        // Capture the stack trace
        buf := make([]byte, 10000)
        runtime.Stack(buf, false)

        fmt.Printf("%s : Stack Trace : %s", functionName, string(buf))
    }
}

func SendEmail(host string, port int, userName string, password string, to []string, subject string, message string ,mapdata map[string]string) (err error) {
    defer catchPanic(&err, "SendEmail")

    parameters := struct {
        From string
        To string
        Subject string
        Message string
        
    }{
        userName,
        strings.Join([]string(to), ","),
        subject,
        message,
      
    }

    buffer := new(bytes.Buffer)

    template := template.Must(template.New("emailTemplate").Parse(emailScript(mapdata)))
    template.Execute(buffer, &parameters)

    auth := smtp.PlainAuth("", userName, password, host)

    err = smtp.SendMail(
        fmt.Sprintf("%s:%d", host, port),
        auth,
        userName,
        to,
        buffer.Bytes())

    return err
}

func emailScript(mapdata map[string]string) (string) {
    return `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

{{.Message}}

<html><body>
<table rules="all" style="border-color: #666;" cellpadding="10">
<tr style='background: #eee;'>

    
        
</tr>
</table>
</body></html>


`
}