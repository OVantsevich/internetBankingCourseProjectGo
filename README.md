# internetBankingCourseProjectGo
userService api

Request:	
		URL: POST localhost:12345\createUser

		body: "user_name"required:   	 string			
		"surname"required: 	 string
		"user_login"required: 	 string{unique}
		"user_email"required: 	 string{name@domainName}
		"user_password"required: string

Responses: 	information line: string

Request:	URL: GET localhost:12345\signIn

		body: "user_login"required: 	 string{unique}
		"user_password"required: string

Responses: 	token: string

Request:	URL: PUT localhost:12345\updateUser
		Header: "Authorization"required: Bearer “token”

		body: "user_name"optional:   	 string			
		"surname"optional: 		 string
		"user_email"optional: 	 string{name@domainName}
		"user_password"optional: 	 string

Responses: 	information line: string

Request:	URL: DEL localhost:12345\deleteUser
		Header: "Authorization"required: Bearer “token”

Responses: 	information line: string

Account service api

Request:	URL: POST localhost:12344\createAccount
		Header: "Authorization"required: Bearer “token”

		body: "account_name"required:    string	

Responses: 	information line: string

Request:	URL: GET localhost:12344\getUserAccounts
		Header: "Authorization"required: Bearer “token”

Responses: 	array of structs
		[
			"account_name":string
			"amount":int
			"account_number":string
			"creation_date":date-time
		]


Request:	URL: POST localhost:12344\createTransaction
		Header: "Authorization"required: Bearer “token”

		body: "account_sender_number"required:   	 string			
		"account_receiver_number"required: 	 string
		"amount"required: 	 int

Responses: 	information line: string

Request:	URL: GET localhost:12344\getAccountTransactions
		Header: "Authorization"required: Bearer “token”
		Header: "Account-Number"required: string

Responses: 	array of structs
		[
			"account_sender_number":string
			"account_receiver_number":int
			"amount":string
			"creation_date":date-time
		]
