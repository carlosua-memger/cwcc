<%
varBookingID = Request.QueryString("bid")
varBookingID = FixNumericHack(varBookingID)

if varBookingID=0 Then
	varBookingID=""
End If


if LEN(varBookingID)<>0 AND ((Session("suname")="") OR (ISNULL(Session("suname"))) OR (LEN(Session("suname"))=0)) Then
	varWalkIn = True
Else
	varWalkIn = False
End If

' keyword logico
varKeywords = GetMasterSEOKeywords()
%>

<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=windows-1252" />
<title>European Accommodation - Euroscape Travel hotels apartments tours sightseeing car hire cruises</title>
<meta name="description" content="euroscape travel planning holiday itinerary for independent travellers, accommodation in hotels, apartments, villas, farmhouses, citytours, escorted tours, transfers, product catalogue in 25 european countries" />
<meta name="keywords" content="<%=varKeywords%>" />
<!--#include file="lib/functions.inc"-->
<!--#include file="inc_header.asp"-->

<!--#include file="ajax_captcha_init.asp" --> 
<script language="javascript" src="ajax_captcha.js"></script> 
<script type="text/javascript" src="http://api.recaptcha.net/js/recaptcha_ajax.js"></script> 

<!--bootrap dependency-->
<link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet" type="text/css"/>

<link href="https://cdnjs.cloudflare.com/ajax/libs/bootstrap-validator/0.5.3/css/bootstrapValidator.min.css" rel="stylesheet" type="text/css"/>
<style type="text/css">
	.control-label .required{
		color: #e02222;
	    font-size: 12px;
	    padding-left: 2px;
	}
</style>

<script src="https://code.jquery.com/jquery-3.1.1.min.js" integrity="sha256-hVVnYaiADRTO2PzUGmuLJr8BLUSjGIZsDYGmIJLv2b8=" crossorigin="anonymous"></script>
<script type="text/javascript" src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script> 
<script src='http://cdnjs.cloudflare.com/ajax/libs/bootstrap-validator/0.4.5/js/bootstrapvalidator.min.js'></script>
<!-- end boot strap dependency -->


</head>
<body onload="overtab(0);">
<div id="topbar">
<p><img border="0" src="images/loading_big.gif" alt="Loading..." /><br /><br />Processing... please wait.</p>
</div>
<div id="container_page" >
	<!--#include file="inc_menu.asp"-->
	
	<div id="page_label" style="background-color: #5375b8;">&nbsp;&nbsp;<span class="where">Where 
		Am I?</span>&nbsp;&nbsp;Home &gt; <%if NOT (varWalkIn) Then%>Travel Agent / 
		Consultant<%Else%>Traveller<%End if%> Registration<p onclick="parent.history.back();">
		&nbsp;&nbsp;&nbsp;&nbsp;GO BACK</p></div>
	
	<!--#include file="inc_boxes.asp"-->
	<!--#include file="inc_search.asp"-->	
	<div id="pageTitle"><h4>&nbsp;<!--<%if NOT (varWalkIn) Then%>Travel Agent / 
		Consultant<%Else%>Traveller<%End if%>-->User Registration</h4></div>
	<a name="sel"><div id="body_container"></a>
	<%
		 	Dim BookingString	
			BookingString = "" 
			Set ConnString = Server.CreateObject("ADODB.Connection")
			ConnString.Open Application("DSN_SQLOLEDB")
			Set rstProd      = Server.CreateObject("ADODB.Recordset")
			Set rstProd.ActiveConnection= connString
				rstProd.Open "Select mgt_tbl_text from msystxt0 where  mgt_tbl_code = 'Booking Selection2' "
				If not (rstProd.EOF and rstProd.BOF) then
					While not rstProd.EOF
						If (BookingSelection) = "" then
							BookingSelection = rstProd("mgt_tbl_text")
						Else	
							BookingSelection = BookingSelection&rstProd("mgt_tbl_text")
						End If
						rstProd.MoveNext
					Wend 
				End If
			rstProd.close()
			ConnString.Close()
	%>
	<table id="registration_content" border="0" width="768" style="border-collapse: collapse" bordercolor="#C0C0C0" cellpadding="0">
	<tr>
		<td valign="top" style="width:240px">
			<p class="desc"><br /><%=BookingSelection%></p>			
		</td>
		<td valign="top" id="body_bg">
		<div id="registration" <%if (varWalkIn) Then%>style="display: none;"<%Else%>style="width:520px;"<%End If%>>
			<form method="post" name="register" action="">
			<table border="0" cellpadding="2" style="border-collapse: collapse; width:520px;">			
				<tr>
					<td height="40px" class="desc"></td>
				</tr>
				<tr align="left">	
						
					<!--old
					<td class="desc"><a href="memb	erregistration.asp#m" class="register_links"><input type="radio"  name="R1" size="20"  value="membe" 	checked align="bottom" ><u>Create a Member Account</u></a><br><br></td>
					end old-->	
					<td class="desc"><a data-toggle="modal" href="#myModal" class="register_links"><input type="radio"  name="R1" size="20"  value="membe" 	checked align="bottom" ><u>Create a Member Account</u></a><br><br></td>
					
						
						
				</tr>		 
				<tr align="left">
						<td class="desc"><a href="register.asp#a"  class="register_links" ><input type="radio"  name="R1" size="20 "   value="agency" valign="bottom" ><u>Travel Agency</u></a><br><br><td>
				</tr>
				<tr align="left">
						<td class="desc"><a href="parentregistration.asp"  class="register_links"><input type="radio" name="R1" size="20" value="partner" valign="bottom"><u>Partners</u></a><br></td>
				</tr>
			</table>
			</form>
		</div>
		</td>
	</tr>
	</table>
	</div>
	<!--#include file="inc_footer.asp"-->
</div>

<!-- Modal form -->
					<div class="modal fade" id="myModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">
						<div class="modal-dialog" role="document">
						    <div class="modal-content">
						    	<div class="modal-header">
							        <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
							        <h4 class="modal-title" id="myModalLabel">User Registration</h4>
						      	</div>
						      <form method="POST" action="save_registration.asp" on id="reg_form" >
						      <div class="modal-body">
						    	
						          <div class="form-group">
						            <label for="recipient-name" class="control-label">Salutation:</label>
						            <select  name="salutation" class="form-control">
						          		<option>Mr</option>
						          		<option>Mrs</option>	
						          		<option>Ms</option>	
						          		<option>Mst</option>
						          		<option>Dr</option>	
						          		<option>Prof</option>	
									</select>
						          </div>

						          <div class="form-group">
							        <label class="	control-label">First/Given Name: <span class="required" aria-required="true">* </span></label>
							        <div class="inputGroupContainer">
							          <div class="input-group"> <span class="input-group-addon"><i class="glyphicon glyphicon-user"></i></span>
							            <input  name="first_name" placeholder="First Name" class="form-control"  type="text">
							          </div>
							        </div>
							      </div>
						          
								  <div class="form-group">
						            <label for="recipient-name" class="control-label">Last Name/Family Name:<span class="required" aria-required="true">* </span></label>
						            <div class="input-group"> 
						            	<span class="input-group-addon">
						            		<i class="glyphicon glyphicon-user"></i>
						            	</span>
						            	<input name="last_name" placeholder="Last Name" class="form-control" type="text">
						          	</div>
						          </div>

						          <div class="form-group">
						            <label for="recipient-name" class="control-label">Address Line1:<span class="required" aria-required="true">* </span></label>
						            <div class="input-group"> 
						            	<span class="input-group-addon">
						            		<i class="glyphicon glyphicon-home"></i>
						            	</span>
						            	<input name="address1" placeholder="Address1" class="form-control" type="text">
						          	</div>
						          </div>

						          <div class="form-group">
						            <label for="recipient-name" class="control-label">Address Line2:</label>
						            <div class="input-group"> 
						            	<span class="input-group-addon">
						            		<i class="glyphicon glyphicon-home"></i>
						            	</span>
						            	<input name="address2" placeholder="Address2" class="form-control" type="text">
						          	</div>
						          </div>
						          
						          <%
										varIPAddress = Request.ServerVariables("REMOTE_ADDR")
										myArray = Split(varIPAddress, ".")

										IPNo = CINT(myArray(0)) * 16777216
										IPNo = IPNo + (CINT(myArray(1)) * 65536 ) 
										IPNo = IPNo + (CINT(myArray(2)) * 256 ) 
										IPNo = IPNo + (CINT(myArray(3))) 

										Set cmdStoredProc = Server.CreateObject("ADODB.Command")
										cmdStoredProc.ActiveConnection = Application("DSN_SQLOLEDB")
										cmdStoredProc.CommandType = 4 'Stored Procedures
										cmdStoredProc.CommandText = "sp_IP2Country"
										cmdStoredProc.Parameters("@mode") = 1
										cmdStoredProc.Parameters("@IP") = CStr(IPNo)
										Set RS_ps= cmdStoredProc.Execute()
										if RS_ps.BOF OR RS_ps.EOF Then
											varCountryCode2 = "AU"
										Else
											varCountryCode2 = RS_ps("country_code2")
										End If
										RS_ps.close
										Set RS_ps=Nothing
										
										if LEFT(varIPAddress,7)="192.168" Then
											varCountryCode2 = "AU"
										End if 
										%>
						          <div class="form-group">
						            <label for="recipient-name" class="control-label" value="<%=varCountryCode2%>">Country:</label>
						            <select class="form-control" size="1" name="country" onchange="//getState(this.form,this,this.form.state,this.form.city,1);">
										<%
										Set cmdStoredProc = Server.CreateObject("ADODB.Command")
										cmdStoredProc.ActiveConnection = Application("DSN_SQLOLEDB")
										cmdStoredProc.CommandType = 4 'Stored Procedures
										cmdStoredProc.CommandText = "usp_GetGeoNameList"
										cmdStoredProc.Parameters("@tiMode") = 1
										Set RS_ps= cmdStoredProc.Execute()
										if RS_ps.BOF OR RS_ps.EOF Then
										Else
											Do While NOT RS_ps.EOF
												varCountryName = TRIM(RS_ps("geoName"))
												varCountryCode = TRIM(RS_ps("geoID"))
												varCountryName = Replace(varCountryName, "'", "\'")
												%><option value="<%=varCountryCode%>"<%if varCountryCode=varCountryCode2 then%> selected="selected"<%End If%>><%=varCountryName%></option>
												<%RS_ps.MoveNext
											Loop
										End If
										RS_ps.close
										Set RS_ps=Nothing%>
									</select><%'=varCountryCode2%><%'=varIPAddress%>
										
								<%'Comments (start)7/6/2011 %>		
						          </div>
						          
						          <div class="form-group">
						            <label for="recipient-name" class="control-label">City/Town:<span class="required" aria-required="true">* </span></label>
						            <div class="input-group">
						            	<span class="input-group-addon">
						            		<i class="glyphicon glyphicon-home"></i>
						            	</span>
						            	<input name="city" placeholder="City / Town" class="form-control" type="text">
						          	</div>
						          </div>

						          <div class="form-group">
							        <label class="control-label">Zipcode:</label>
							        <div class=" inputGroupContainer">
							          <div class="input-group"> <span class="input-group-addon"><i class="glyphicon glyphicon-home"></i></span>
							            <input name="zipcode" class="form-control" type="text">
							          </div>
							        </div>
							      </div>
						          
						          <div class="form-group">
							        <label class="control-label">Telephone:</label>
							        <div class=" inputGroupContainer">
							          <div class="input-group"> <span class="input-group-addon"><i class="glyphicon glyphicon-earphone"></i></span>
							            <input name="phone" placeholder="(845)555-1212" class="form-control" type="text">
							          </div>
							        </div>
							      </div>

						          <div class="form-group">
						            <label for="recipient-name" class="control-label">Mobile Phone:</label>
						            <input type="text" class="form-control" id="recipient-name" name="mobile">
						          </div>

						          <div class="form-group">
							        <label for="recipient-name" class="control-label">Preferred Currency:<span class="required" aria-required="true">* </span></label>
							        <div class="selectContainer">
							          <div class="input-group"> <span class="input-group-addon"><i class="glyphicon glyphicon-list"></i></span>
							            <select name="currency" class="form-control selectpicker" >
							              <option value=" " >Please select your Currency</option>
							              <option value="AU">Australia Dollars</option>
							              <option value="EUR">Euro</option>
							              <option value="GBP">United Kingdom Pounds</option>
							              <option value="NZD">New Zealand Dollars</option>
							              <option value="USD">United States Dollars</option>
										</select>
							          </div>
							        </div>
							      </div>

						       	<h1>Login Information</h1>	
						       	<div class="form-group">
						        	<label for="recipient-name" class="control-label">Email:<span class="required" aria-required="true">* </span></label>
						        	<div class="input-group">
						        		<span class="input-group-addon"><i class="glyphicon glyphicon-envelope"></i></span>
						            	<input name="email" placeholder="E-Mail Address" class="form-control"  type="text">
						          </div>
						        </div>
						        <div class="form-group">
						        	<label for="recipient-name" class="control-label">Retype Email:<span class="required" aria-required="true">* </span></label>
						            <div class="input-group">
						            	<span class="input-group-addon"><i class="glyphicon glyphicon-envelope"></i></span>
						            	<input name="confirmEmail" placeholder="E-Mail Address" class="form-control"  type="text">
						          	</div>

						        </div>
						        <div class="form-group">
						        	<label for="recipient-name" class="control-label">Select Password:<span class="required" aria-required="true">* </span></label>
						            <div class="input-group"> 
						            	<span class="input-group-addon">
						            		<i class="glyphicon glyphicon-pencil"></i>
						            	</span>
						            	<input name="password" placeholder="Password" class="form-control" type="password" >
						          	</div>
						        </div>
						        <div class="form-group">
						        	<label for="recipient-name" class="control-label">Retype Password:<span class="required" aria-required="true">* </span></label>
						            <div class="input-group"> 
						            	<span class="input-group-addon">
						            		<i class="glyphicon glyphicon-pencil"></i>
						            	</span>
						            	<input name="confirmPassword" placeholder="Confirm Password" class="form-control" type="password">
						          	</div>
						        </div>	 	  	   
						        
						       
						      </div>
						      <div class="modal-footer">
						        <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
       							<button type="submit" class="btn btn-primary">Register</button>
						      </div>
						    </form> 
						    </div>
						    

					  	</div>
</div>

						<!--  end modal -->
						
<!-- BEGIN Invitation Positioning  -->
<script language="javascript" type="text/javascript">
var lpPosY = 100;
</script>
<!-- END Invitation Positioning  -->
<!-- BEGIN HumanTag Monitor. DO NOT MOVE! MUST BE PLACED JUST BEFORE THE /BODY TAG -->
<!--<script type="text/javascript" language='javascript' src='http://server.iad.liveperson.net/hc/65260606/x.js?cmd=file&amp;file=chatScript3&amp;site=65260606&amp;&amp;imageUrl=http://www.euroscape-travel.com/images/hotels/LivePerson'> </script>-->
<!-- END HumanTag Monitor. DO NOT MOVE! MUST BE PLACED JUST BEFORE THE /BODY TAG -->

<!-- bootstrap validation script -->
<script type="text/javascript">
	$(document).ready(function() {
	    $('#reg_form').bootstrapValidator({
	        // To use feedback icons, ensure that you use Bootstrap v3.1.0 or later
	        feedbackIcons: {
	            valid: 'glyphicon glyphicon-ok',
	            invalid: 'glyphicon glyphicon-remove',
	            validating: 'glyphicon glyphicon-refresh'
	        },
	        fields: {
	            first_name: {
	                validators: {
	                        stringLength: {
	                        min: 2,
	                    },
	                        notEmpty: {
	                        message: 'Please supply your first name'
	                    }
	                }
	            },
	             last_name: {
	                validators: {
	                     stringLength: {
	                        min: 2,
	                    },
	                    notEmpty: {
	                        message: 'Please supply your last name'
	                    }
	                }
	            },
	           
	            phone: {
	                validators: {
	                    /*notEmpty: {
	                        message: 'Please supply your phone number'
	                    },*/
	                    phone: {
	                        country: 'US',
	                        message: 'Please supply a vaild phone number with area code'
	                    }
	                }
	            },
	            address1: {
	                validators: {
	                     stringLength: {
	                        min: 8,
	                    },
	                    notEmpty: {
	                        message: 'Please supply your street address'
	                    }
	                }
	            },
	            city: {
	                validators: {
	                     stringLength: {
	                        min: 4,
	                    },
	                    notEmpty: {
	                        message: 'Please supply your city'
	                    }
	                }
	            },
	            currency: {
	                validators: {
	                    notEmpty: {
	                        message: 'Please select your currency'
	                    }
	                }
	            },
	            zip: {
	                validators: {
	                    notEmpty: {
	                        message: 'Please supply your zip code'
	                    },
	                    zipCode: {
	                        country: 'US',
	                        message: 'Please supply a vaild zip code'
	                    }
	                }
	            },
			comment: {
	                validators: {
	                      stringLength: {
	                        min: 10,
	                        max: 200,
	                        message:'Please enter at least 10 characters and no more than 200'
	                    },
	                    notEmpty: {
	                        message: 'Please supply a description about yourself'
	                    }
	                    }
	                 },	
		 email: {
	                validators: {
	                    identical: {
		                    field: 'confirmEmail',
		                    message: 'The Email and its confirm are not the same'
	                	},
	                    emailAddress: {
	                        message: 'Please supply a valid email address'
	                    },
	                    notEmpty: {
                        message: 'Please supply your email'
                    }
	                }
	            },
			confirmEmail: {
	            validators: {
	                identical: {
	                    field: 'email',
	                    message: 'The Email and its confirm are not the same'
	                },
	                notEmpty: {
                        message: 'Please cofirm email'
                    }
	            }
	         },			
		password: {
	            validators: {
	                identical: {
	                    field: 'confirmPassword',
	                    message: 'Confirm your password below - type same password please'
	                },
	                notEmpty: {
                        message: 'Please supply your password'
                    }
	            }
	        },
	        confirmPassword: {
	            validators: {
	                identical: {
	                    field: 'password',
	                    message: 'The password and its confirm are not the same'
	                },
	                notEmpty: {
                        message: 'Please confirm password'
                    }
	                
	            }
	         },
				
	            
	            }
	        })
			
	 	
	        .on('success.form.bv', function(e) {
				//console.log($("[name=country]").val())
				$('#success_message').slideDown({ opacity: "show" }, "slow") // Do something ...
	               $('#reg_form').data('bootstrapValidator').resetForm();

	            // Prevent form submission
	            e.preventDefault();

	            // Get the form instance
	            var $form = $(e.target);

	            // Get the BootstrapValidator instance
	            var bv = $form.data('bootstrapValidator');

	            // Use Ajax to submit form data
	            $.post("save_registration.asp", $form.serialize(), function(result) {
	                console.log(result);
					alert ("Success!, you may Login");
					return false;
	            }, 'json');
	        });
	});
	
	///
	//function Member_save () {
		//console.log($("[name=country]").val())
		//$.post("save_registration.asp",
        	//$("#reg_form").serialize()
	//	,
    //    function(data,status){
        	//console.log(data)
        //	alert ("Success!, you may Login")
      //  });		

	//}	

</script>
<!-- end bootstrap validation script -->

</body>
</html>