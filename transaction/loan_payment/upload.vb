Imports System.io
Imports System.Data.SqlClient


Public Class trx_lp_upload
    Inherits System.Web.UI.Page

#Region " Web Form Designer Generated Code "

    'This call is required by the Web Form Designer.
    <System.Diagnostics.DebuggerStepThrough()> Private Sub InitializeComponent()

    End Sub
    Protected WithEvents MyFile As System.Web.UI.HtmlControls.HtmlInputFile
    Protected WithEvents Submit1 As System.Web.UI.HtmlControls.HtmlInputButton
    Protected WithEvents UploadDetails As System.Web.UI.HtmlControls.HtmlGenericControl
    Protected WithEvents FileName As System.Web.UI.HtmlControls.HtmlGenericControl
    Protected WithEvents FileContent As System.Web.UI.HtmlControls.HtmlGenericControl
    Protected WithEvents FileSize As System.Web.UI.HtmlControls.HtmlGenericControl
    Protected WithEvents Span1 As System.Web.UI.HtmlControls.HtmlGenericControl
    Protected WithEvents Span2 As System.Web.UI.HtmlControls.HtmlGenericControl
    Protected WithEvents lblTextFile As System.Web.UI.WebControls.Label
    Protected WithEvents hlReturnToDetail As System.Web.UI.WebControls.HyperLink
    Protected WithEvents lblModule As System.Web.UI.WebControls.Label
    Protected WithEvents lblID As System.Web.UI.WebControls.Label

    'NOTE: The following placeholder declaration is required by the Web Form Designer.
    'Do not delete or move it.
    Private designerPlaceholderDeclaration As System.Object

    Private Sub Page_Init(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles MyBase.Init
        'CODEGEN: This method call is required by the Web Form Designer
        'Do not modify it using the code editor.
        InitializeComponent()
    End Sub

#End Region
    Dim strServerPath As String
    Dim strServerFile As String
    Dim intID As Integer
    Dim intOK As Integer = 0
    Dim intNotOK As Integer = 0


    Private Sub Page_Load(ByVal sender As System.Object, ByVal e As System.EventArgs) Handles MyBase.Load

        Submit1.Attributes("onclick") = "return confirm('Remove detail and upload text file?');"

        strServerPath = Server.MapPath(".\TextFileUploads")                      ' Set the path.
        intID = Request("id")
        'lblTextFile.Text = intInvCnt
        hlReturnToDetail.NavigateUrl = "trx_lpdtl_list.aspx?id=" & intID
        lblID.Text = intID
        If Dir(strServerPath, vbDirectory) = "" Then   ' The folder is not there & to be created
            MkDir(strServerPath)        'Folder created
            Span2.InnerHtml = " The destination folder (" & strServerPath & ") does not exist and was  created."
        End If

    End Sub

    Private Sub Submit1_ServerClick(ByVal sender As System.Object, ByVal e As System.EventArgs) _
    Handles Submit1.ServerClick
        ' Display properties of the uploaded file
        FileName.InnerHtml = MyFile.PostedFile.FileName
        FileContent.InnerHtml = MyFile.PostedFile.ContentType
        FileSize.InnerHtml = MyFile.PostedFile.ContentLength
        UploadDetails.Visible = True

        ' Let us recover only the file name from its fully qualified path at client 
        Dim strFileName As String
        strFileName = MyFile.PostedFile.FileName
        Dim c As String = System.IO.Path.GetFileName(strFileName) ' only the attched file name not its path
        strServerFile = strServerPath & "\" & c

        If MyFile.PostedFile.ContentType <> "text/plain" Or _
         Right(MyFile.PostedFile.FileName, 4) <> ".txt" Then
            Span1.InnerHtml = "Please upload a valid text file. "
            UploadDetails.Visible = False
            Span2.Visible = False
            Exit Sub
        End If

        Try
            MyFile.PostedFile.SaveAs(strServerFile)
            Span1.InnerHtml = "File was copied sucessfully at server as : " & strServerFile
            UploadToDB()
        Catch Exp As Exception
            Span1.InnerHtml = "An error occured. Please check the attached  file"
            UploadDetails.Visible = False
            Span2.Visible = False
        End Try

    End Sub



    Private Sub UploadToDB()

        ' Record read
        Dim intPartner As Integer
        Dim datTranDate As Date

        Dim objConnection As New SqlConnection(ConfigurationSettings.AppSettings("appConn"))
        Dim objCommand As New SqlCommand("trx_LPDTL_Upload_ZapDetail", objConnection)
        objCommand.CommandType = CommandType.StoredProcedure
        objCommand.Parameters.Add("@hdr", intID)
        objConnection.Open()
        objCommand.ExecuteNonQuery()

        objCommand = New SqlCommand("trx_lphdr_get", objConnection)
        objCommand.CommandType = CommandType.StoredProcedure
        objCommand.Parameters.Add("@ID", intID)
        Dim objReader As SqlDataReader = objCommand.ExecuteReader
        objReader.Read()
        intPartner = objReader("partner")
        datTranDate = objReader("trandate")

        objReader.Close()
        objConnection.Close()
        ' File read
        Dim objStreamReader As StreamReader
        Dim strLine As String
        Dim bolHeaderLocationMatch As Boolean = True

        objStreamReader = File.OpenText(strServerFile)

        While objStreamReader.Peek() <> -1 And bolHeaderLocationMatch ' do while thier is still file line 
            strLine = objStreamReader.ReadLine()

            If Left(strLine, 1) = "^" Then
                'HEADER HERE
                If intPartner = strLine.Substring(1, 4) And _
                    datTranDate.Month = strLine.Substring(5, 2) And _
                    datTranDate.Day = strLine.Substring(7, 2) And _
                    datTranDate.Year = strLine.Substring(9, 4) Then

                    bolHeaderLocationMatch = True
                    lblTextFile.Text = lblTextFile.Text & "</br>Processing a match header line :" & strLine
                Else
                    lblTextFile.Text = lblTextFile.Text & "</br>Header line not match:" & strLine
                    bolHeaderLocationMatch = False
                End If

            Else
                'LINE HERE
                'only upload to table if below is true 
                If bolHeaderLocationMatch Then
                    lblTextFile.Text = lblTextFile.Text & "</br>Processing detail line " & strLine
                    UploadThisLine(strLine, intPartner)
                End If
            End If
        End While

        objStreamReader.Close()

    End Sub


    Private Sub UploadThisLine(ByVal strLine As String, ByVal intPartner As Integer)
        ' code 128 16 char sample : 0012260010149975

        'character definitions:
        'lrid   
        Dim intLridStart As Integer = 0
        Dim intLridLen As Integer = 8 ' 8 
        Dim strLrid As String = ""

        'member
        Dim intMemberStart As Integer = 18
        Dim intMemberLen As Integer = 8     ' 99,999,999
        Dim strMember As String = ""

        'transaction type 
        Dim intTypeStart As Integer = 26
        Dim intTypeLen As Integer = 4         ' 99 
        Dim strType As String = ""

        'amount 
        Dim intAmtStart As Integer = 30 '32
        Dim intAmtLen As Integer = 9           ' 9,999,999.99
        Dim strAmt As String = ""





        Dim bolValidLine As Boolean = True


        'Check here if line is valid base on length of total expected characters
        If strLine.Length < (intMemberLen + intTypeLen + intAmtLen + intLridLen) Then
            bolValidLine = False
            lblTextFile.Text = lblTextFile.Text & " Line too short. "
        End If

        'get and validate member and trantype 
        If bolValidLine Then

            strMember = strLine.Substring(intMemberStart, intMemberLen).TrimStart("0")
            strType = strLine.Substring(intTypeStart, intTypeLen)
            strAmt = strLine.Substring(intAmtStart, intAmtLen).TrimStart("0")
            strLrid = strLine.Substring(intLridStart, intLridLen).TrimStart("0")

            '            lblTextFile.Text = lblTextFile.Text + " this line strmember " + strMember + _
            '            " strType " + strType + _
            '" strAmt " + strAmt + _
            '" strLrid " + strLrid


            Dim objConnection As New SqlConnection(ConfigurationSettings.AppSettings("appConn"))
            Dim objCommand As New SqlCommand("trx_lpDTL_ValidateUpload", objConnection)
            objCommand.CommandType = CommandType.StoredProcedure
            objCommand.Parameters.Add("@partner", intPartner)
            objCommand.Parameters.Add("@member", strMember)
            objCommand.Parameters.Add("@loantype", strType)
            objCommand.Parameters.Add("@lrid", strLrid)



            objConnection.Open()
            Dim objReader As SqlDataReader = objCommand.ExecuteReader

            If objReader.HasRows Then
                bolValidLine = False
                While objReader.Read()
                    lblTextFile.Text = lblTextFile.Text & objReader("errormessage")
                End While
            End If

            objReader.Close()
            objConnection.Close()

        End If


        If bolValidLine = False Then
            lblTextFile.Text = lblTextFile.Text & " Invalid line. "
            intNotOK = intNotOK + 1

        Else
            lblTextFile.Text = lblTextFile.Text & " Valid line. "

            '27Jan2011 - Check price variance as per Danny, similar with delivery uploading
            'CheckThisLine(intProduct, dblPrice)

            'the connection
            Dim objConnection As New SqlConnection(ConfigurationSettings.AppSettings("appConn"))

            'the stored procedure and parameters
            Dim objCommand As New SqlCommand("Trx_lpdtl_Upload", objConnection)
            objCommand.CommandType = CommandType.StoredProcedure
            'objCommand.Parameters.Add("@FormAction", "Upload")
            objCommand.Parameters.Add("@hdr", intID)
            objCommand.Parameters.Add("@ID", 0)
            objCommand.Parameters.Add("@member", strMember)
            objCommand.Parameters.Add("@loantype", strType)
            objCommand.Parameters.Add("@amt", strAmt)
            objCommand.Parameters.Add("@remarks", strLine)
            objCommand.Parameters.Add("@lrid", strLrid)

            'the returned new code 
            Dim objNewID As New SqlParameter("@newID", SqlDbType.Int)
            objCommand.Parameters.Add(objNewID)
            objNewID.Direction = ParameterDirection.Output

            Try
                'open the connection and execute the query 
                objConnection.Open()
                objCommand.ExecuteNonQuery()

                'save to audit trail 
                Sys_Class.SaveToAuditTrail(Session("username"), _
                                            "Upload", _
                                            "trx_lpdtl", _
                                            objNewID.Value, _
                                            "")

                lblTextFile.Text = lblTextFile.Text & " Line saved."
                intOK = intOK + 1
            Catch ex As Exception
                'in case of error
                'lblTextFile.Text = lblTextFile.Text & " An error occured while saving record - " & ex.ToString
                lblTextFile.Text = lblTextFile.Text & " Line NOT saved. Error:  " & ex.ToString
                '& strLine
                intNotOK = intNotOK + 1
            Finally
                objConnection.Close()
            End Try

        End If
    End Sub


End Class
