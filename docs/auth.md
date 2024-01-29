# Auth API Spec

## Register API

Endpoint :  POST /api/auth/register

Request Body :

```json
{
    "username" : "test",
    "password" : "rahasia",
    "email" : "test@gmail.com"
}
```

Response Body Success :

```json
{
    "code": 200,
    "status": "Ok",
    "message": "Successfully created user!"
}
```

Response Body Error : 

```json
{
    "code": 400,
    "status": "BadRequest",
    "message": "Error Validation",
    "data": {
        "Email": "email wajib diisi",
        "Username": "username wajib diisi"
    }
}
```

## Login API

Endpoint : POST /api/auth/login

Request Body :

```json
{
    "username" : "test",
    "password" : "rahasia"
}
```

Response Body Success : 

```json
{
    "code": 200,
    "status": "Ok",
    "message": "Successfully log in!",
    "data": {
        "token_type": "Bearer",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDU2Mzk1NjUsImlhdCI6MTcwNTYzNTk2NSwibmJmIjoxNzA1NjM1OTY1LCJzdWIiOjR9.jO5b_ophTcHTNPvqZ7J9eRhawldiokb6UE342s3QQhY"
    }
}
```

Response Body Error :

```json
{
    "code": 400,
    "status": "BadRequest",
    "message": "Error Validation",
    "data": {
        "Password": "password wajib diisi",
        "Username": "username wajib diisi"
    }
}
```

## Forgot Password API

Endpoint : POST /api/auth/forgot-password

Request Body :

```json
{
    "email" : "pqm90601.com"
}
```

Response Body Success : 

```json
{
    "code": 200,
    "status": "OK",
    "message": "Forgot Passworsssd successfully",
    "data": 7951 // data otp
} 
```

Response Body Error : 

```json
{
    "code": 400,
    "status": "BadRequest",
    "message": "Error Validation",
    "data": {
        "Email": "email bukan email yang valid"
    }
}
```

## Check Otp API

Endpoint : POST /api/auth/check-otp

Request Body :

```json
{
    "otp" : 7951
}
```

Response Body Success : 

```json
{
    "code": 200,
    "status": "OK",
    "message": "Otp Valid"
}
```

Response Body Error : 

```json
{
    "code": 400,
    "status": "BadRequest",
    "message": "Error Validation",
    "data": {
        "Otp": "otp wajib diisi"
    }
}
```

## Reset Password API

Endpoint : PATCH /api/auth/reset-password?otp=9904

Response Body Success:

```json
{
    "code": 200,
    "status": "OK",
    "message": "Reset Password successfully"
}
```

Response Body Error : 

```json
{
    "code": 400,
    "status": "BadRequest",
    "message": "Error Validation",
    "data": {
        "PasswordConfirmation": "password_confirmation wajib diisi"
    }
}
```

## Logout API

Endpoint : POST /api/auth/logout

Headers :
- Authorization : Bearer token

Response Body Success : 

```json
{
    "code": 200,
    "status": "OK",
    "message": "Logout Successfully"
}
```

Response Body Error : 

```json
{
    "error" : "invalid token"
}
```
Response Body Error : 

```json
{
    "message": "Unauthorized",
    "status": "fail"
}
```