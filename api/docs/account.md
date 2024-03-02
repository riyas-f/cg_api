## Account
**Properties**

| name         | type     | description                                                                |
| ------------ | -------- | -------------------------------------------------------------------------- |
| username     | `string` | A 64 length string that identifies a user                                  |
| name         | `string` | User display name                                                          |
| phone_number | `string` | User phone number                                                          |
| email        | `string` | A valid email for authentication                                           |
| password     | `string` | A minimum of 8 string password with combination of character and symbol    |

**Header**

| name          | description                                                       |
| ------------- | ----------------------------------------------------------------- |
| Content-Type  | The content type of the request. Only accepted `application/json` |
| Authorization | Authorization header. Use `Bearer Token`.                                                                  |

### 1. Register
To register a user to the server, send a POST request to the endpoint below:

```
POST /v1/account/register
```
** Header **

| name         | value                                                                                                              |
| ------------ | ------------------------------------------------------------------------------------------------------------------ |
| Content-Type | <span style="padding: 2px 4px; color: rgb(175, 78, 93); background-color: rgb(30, 15, 17)">application/json</span> |

** Body ** 
```json
{
	"username": "required",
	"name": "required",
	"phone_number": "required",
	"email": "required",
	"password": "required",
}
```

#### Response
**Header**

| name         | value              |
| ------------ | ------------------ |
| Content-Type | `application/json` | 

**Properties**

| name                | description                                                    |
| ------------------- | -------------------------------------------------------------- |
| OTP_confirmation_id | A 64 character string for verifying your created account       |
| token               | A jwt token for giving authorization on verifying your account |
| error_type          | a string for telling the error type                            |

**Succesful Response (200 OK)**

```json
{
	"status": "successfull",
	"OTP_Confirmation_id": "your_id_here",
	"access_token" : {
		"token": "your_token_here",
		"expires_in": "300",
	}
}
```

**Error Response**
1. **400 Bad Request**
	There are various 400 error type, below is the details regarding all possible bad request error type:
	
	
	| error_type           | description                                            |
	| -------------------- | ------------------------------------------------------ |
	| missing_parameter   | There are a missing parameters in the request provided |
	| header_value_mismatch | using non-supported content-type header                |
	| username_exists      | username already taken                                 |
	| username_invalid     | username content is invalid                            |
	| email_exists         | email already registered                               |
	| password_weak        | invalid password                                       |
	| invalid_email        | provided an invalid email                                                       |
	
```json
{
	"status": "error",
	"error_type": "your_error_type_here",
	"description": "details_regarding_the_error"
}
```

### 2. Verify OTP
After a successful registration, an OTP will be sent to the email on the request. To verify newly registered user, user must verify their account by sending a POST request to the endpoint below:

```
POST /v1/account/otp/verify
```

-  # **Request**

* **Header**

| name          | value                    |
| ------------- | ------------------------ |
| Content-Type  | `application/json`       |
| Authorization | `Bearer your_token_here` |

* **Body**
```json
{
	"OTP_confirmation_id": "your_id_here",
	"OTP": "your_OTP_here",
}
```

- Response
**Successful Response (200 OK)**

* **Header**

| name         | value              |
| ------------ | ------------------ |
| Content-Type | `application/json` |

* **Body**
```json
{
	"status": "success",
	"message": "your account has been activated successfully"
}
```

**Error Response**
1. **400 Bad Request**
	This error can be caused by several reasons. The list of all possible error is listed below:
	
	| error_type           | description                                |
	| -------------------- | ------------------------------------------ |
	| missing_parameters   | some body parameters is missing            |
	| header_value_mismatch | Content-Type header isn't acceptable       |
	| invalid_otp          | OTP provided doesn't match with the server |

```json
{
	"status": "error",
	"error_type": "your_error_type_here",
	"descriptipn": "details_regarding_your_error here",
}
```

2. **403 Unauthorized** 
	This error mainly caused by your token provided in the header. Below is all possible error:

| error_type          | description                                                      |
| ------------------- | ---------------------------------------------------------------- |
| Invalid_auth_header | The header specified a wrong auth method other than Bearer Token |
| empty_auth_header   | Auth header is empty                                             |
| invalid_token       | Token provided is invalid                                        |
| token_expired       | Token has already expired                                                                 |

3. **404 Not Found**
	This error is caused by the otp id in the request. There only one error type that are possible that is `id_not_found` which caused by the request id cannot be found on the server.
	
```json
{
	"status": "error",
	"error_type": "id_not_found",
	"description": "Confirmation ID provided doesn't exist"
}
```

### 3. Resend OTP
A user are able to re-send confirmation email with a new OTP. Re-send mail can be sent in a interval of <span style="padding: 2px 4px; color: rgb(175, 78, 93); background-color: rgb(30, 15, 17)">60 seconds</span>. To re-send the mail, send a POST request to endpoint below:

```
POST /v1/account/otp/{OTP_confirmation_id}/send
```
- **Request**
- **Header**

| name          | value                    |
| ------------- | ------------------------ |
| Authorization | `Bearer your_token_here` |

* **Parameters**

| name                | description                                   |
| ------------------- | --------------------------------------------- |
| OTP_confirmation_id | ID of the user that want the OTP to be re-send | 

- **Response**

| name         | value              |
| ------------ | ------------------ |
| Content-Type | `application/json` |

**Successful Response (200 OK)**
```json
{
	"status": "successfull",
	"message": "OTP has been re-send to your email."
}
```

**Error Response**
1. **403 Unauthorized** 
	This error mainly caused by your token provided in the header. Below is all possible error:

| error_type          | description                                                      |
| ------------------- | ---------------------------------------------------------------- |
| Invalid_auth_header | The header specified a wrong auth method other than Bearer Token |
| empty_auth_header   | Auth header is empty                                             |
| invalid_token       | Token provided is invalid                                        |
| token_expired       | Token has already expired                                        | 

```json
{
	"status": "error",
	"error_type": "your_error_type_here",
	"descriptipn": "details_regarding_your_error here",
}
```

2. **404 Not Found**
	This error is caused by the otp id in the request. There only one error type that are possible that is `id_not_found` which caused by the request id cannot be found on the server.
	
```json
{
	"status": "error",
	"error_type": "id_not_found",
	"description": "Confirmation ID provided doesn't exist"
}
```

3. **429 Too Many Requests** 
		This error is caused by sending a re-send request too early. There are only one possible error type that is, `otp_resend_interval_not_reached`
```json
{
	"status": "error",
	"error_type": "otp_resend_interval_not_reached",
	"description": "Re send mail has already been sent"
}
```

### 4. Login
To log in user to the server and gain access to all service, sent a POST request to endpoint below:

```
POST /v1/account/login
```

- **Request**
* ** Header **

| name         | value                                                                                                              |
| ------------ | ------------------------------------------------------------------------------------------------------------------ |
| Content-Type | <span style="padding: 2px 4px; color: rgb(175, 78, 93); background-color: rgb(30, 15, 17)">application/json</span> |

* **Body**
```json
{
	"email": "your_email_here",
	"password": "your_password_here"
}
```

- **Response**
**Successful Response (200 OK)**
```json
{
	"status": "success",
	"token": {
		"accesss_token": "your_token_here",
		"refresh_token": "your_refresh_token_here",
	}
}
```

**Error Response**
1. **400 Bad Request**
	This error can be caused by several reasons. The list of all possible error is listed below:
	
	| error_type           | description                                |
	| -------------------- | ------------------------------------------ |
	| missing_parameters   | some body parameters is missing            |
	| header_value_mismatch | Content-Type header isn't acceptable       |

```json
{
	"status": "error",
	"error_type": "your_error_type_here",
	"description": "details_regarding_your_error here",
}
```

2. **401 Unauthenticated**
		There are two error that can caused this status code, that is
		
| error_type           | description                                                             |
| -------------------- | ----------------------------------------------------------------------- |
| user_marked_inactive | the user isn't registered successfully and currently marked as inactive |
| invalid_credentials  | email or password is invalid                                            | 

```json
{
	"status": "error",
	"error_type": "your_error_type_here",
	"descriptipn": "details_regarding_your_error here",
}
```

### 5. Refresh Access Token
### 6. Patch User Info
Update user info partially, send a PATCH request to endpoint below:

```
PATCH /v1/account/{username}
```

- **Request**
* **Header**

| name          | value                    |
| ------------- | ------------------------ |
| Content-Type  | `application/json`       |
| Authorization | `Bearer your_token_here` |

* **Body**
```json
{
	"name": "your_new_name_here (optional)",
}
```

- **Response**
* **Header**

| name         | value              |
| ------------ | ------------------ |
| Content-Type | `application/json` |

* **Successful Response (200 OK)**
```json
{
	"status": "success",
	"message": "Success updating user info"
}
```

* **Error Response**
1. **400 Bad Request**
This error can be caused by several reasons. The list of all possible error is listed below:
	
| error_type           | description                                |
| -------------------- | ------------------------------------------ |
| header_value_mismatch | Content-Type header isn't acceptable       |

```json
{
	"status": "error",
	"error_type": "your_error_type_here",
	"descriptipn": "details_regarding_your_error here",
}
```

3. **403 Forbidden**

| error_type          | description                                                      |
| ------------------- | ---------------------------------------------------------------- |
| Invalid_auth_header | The header specified a wrong auth method other than Bearer Token |
| empty_auth_header   | Auth header is empty                                             |
| invalid_token       | Token provided is invalid                                        |
| token_expired       | Token has already expired                                        | 

```json
{
	"status": "error",
	"error_type": "your_error_type_here",
	"descriptipn": "details_regarding_your_error here",
}
```

4. **404 Not Found**
Username not foundfound
