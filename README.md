# Go Marketplace
An API that simulates backend behavior of a marketplace.

## Setup
Make sure you have GoLang installed  
https://go.dev/doc/install

Clone the repo and install required packages.
```
git clone https://github.com/cyzhang39/go_market.git
go mod tidy
```

## Start
Docker for running database
```
docker-compose up -d
```

Run API
```
// remember to set a secret key
export SECRET_KEY=TOPSECRET
go run main.go
```



## Calling API
Here I'm using postman

### Sign up (POST)
http://localhost:8000/users/signup  
Request Body:
```
{
  "firstName": "Tester",
  "lastName": "Test",
  "email": "tester@mail.com",
  "password": "passtest",
  "phone": "1111111111"
}
```
Returned Body:
```
{
    "dev_code": <Verification Code>,
    "email": "tester@mail.com",
    "message": "A 6-digit verification is sent to your email, please enter the code to verify."
}
```
### Verify (POST)
http://localhost:8000/users/verify  
Request Body:
```
{ 
    "email": "tester@mail.com", 
    "code": <Verification Code> 
}
```
Returned Body:
```
{
    "message": "email verified"
}
```
### Login (POST)
http://localhost:8000/users/login  
Request Body:
```
{
  "email": "tester@mail.com",
  "password": "passtest"
}
```
If log in without verifying email, returned body looks like:
```
{
    "error": "Account not verified, please verify to continue."
}
```
With verified email, returned body:
```
{
    "id": <userID>,
    "firstName": "Tester",
    "lastName": "Test",
    "password": "$2a$12$2yhqqicric8qBTKnjTNS/eCxRqXf2L7ORxwD5KpAHx9wfzvYv8NK.",
    "email": "tester@mail.com",
    "phone": "1111111111",
    "verified": true,
    "code": "$2a$10$yrivhyz1l20gKRg8bENBGe9J/gYnyIAxPbJEyEIjeQXrDY0iunqii",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InRlc3RlckBtYWlsLmNvbSIsIkZpcnN0TmFtZSI6IlRlc3RlciIsIkxhc3ROYW1lIjoiVGVzdCIsIlVJRCI6IjExMTExMTExMTEiLCJleHAiOjE3NTc2MjcyNzh9.Y2b595el5gRWW_2M7BM2WWIC6nxM2dc286ZrE2wLkdI",
    "refresh": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6IiIsIkZpcnN0TmFtZSI6IiIsIkxhc3ROYW1lIjoiIiwiVUlEIjoiIiwiZXhwIjoxNzU4MTQ1Njc4fQ.bjusOJISeyYYbyuSbNHe7We8StpJ7LORrPiV9cDNgLM",
    "createTime": "2025-09-10T20:17:22Z",
    "updateTime": "2025-09-10T21:47:58.539Z",
    "uid": <userID>,
    "cart": [],
    "addressInfo": [],
    "status": []
}
```
Note the ``<userID>`` and ``<token>`` here, they are required in later requests

### List an Item
http://localhost:8000/users/listItem  
Request Body:
```
{
    "name": "textbook",
    "price": 100,
    "rating": 5.0,
    "img": "textbook.png"
}
```
Returned Body:
```
"Item added successfully."
```

### View all market items (GET)
http://localhost:8000/users/view  
No request body.  
Returned Body:
```
[
    {
        "ID": <itemID>,
        "name": "textbook",
        "price": 100,
        "rating validate:": null,
        "img": "textbook.png"
    }
]
```
### Search for item (GET)
http://localhost:8000/users/search?name=textbook  
No request body. 
Returned Body:
```
[
    {
        "ID": <itemID>,
        "name": "textbook",
        "price": 100,
        "rating validate:": null,
        "img": "textbook.png"
    }
]
```

### Add item to cart (GET)
http://localhost:8000/add?id=itemID&userID=userID  
No request body. 
Attach ``<token>`` to request Headers.  
```
token:<token>
```
Returned body:
```
"Item successfully added"
```

### List items in cart (GET)
http://localhost:8000/list?id=userID  
No request body.  
Attach ``<token>`` to request Headers.  
Returned Body:
```
100[
    {
        "ID": "68c20926ed72b2005b9a8ecc",
        "name": "textbook",
        "price": 100,
        "rating": 5,
        "img": "textbook.png"
    }
]
```
The beginning 100 represents the total price.

### Remove item from cart (GET)
http://localhost:8000/remove?id=itemID&userID=userID  
No request body.  
Attach ``<token>`` to request Headers.  
Returned Body:
```
"Item removed successfully"
```

### Add address (POST)
http://localhost:8000/addressadd?id=userID  
Can add at most two addresses here, first a home address second a work address.  
Attach ``<token>`` to request Headers.  
Request Body:
```
{
  "house": "Tester home",
  "street": "test street",
  "city": "Test",
  "postal": "11111"
}
```
Returned Body:
```
"Address added"
```

### Delete address (GET)
http://localhost:8000/addressdel?id=userID  
Delete both home and work addresses.  
No request body.  
Attach ``<token>`` to request Headers.  
Returned Body:
```
"Address deleted"
```
### Edit home address (PUT)
http://localhost:8000/addresshomeedit?id=userID  
Edits home address.  
Attach ``<token>`` to request Headers.  
Request Body:
```
{
  "house": "Tester new home",
  "street": "New tester street",
  "city": "New Test",
  "postal": "33333"
}
```
Returned Body:
```
"Home address updated"
```

### Edit work address (PUT)
http://localhost:8000/addressworkedit?id=userID  
Edits work address.  
Attach ``<token>`` to request Headers.  
Request Body:
```
{
  "house": "Tester new office",
  "street": "New tester street",
  "city": "New Test",
  "postal": "44444"
}
```
Returned Body:
```
"Work address updated"
```

### Cart checkout (GET)
http://localhost:8000/checkout?id=userID  
Will empty cart and update to user's status.  
No request body.  
Attach ``<token>`` to request Headers.  
Returned Body:
```
"Order placed successfully"
```

### Buy item instantly (GET)
http://localhost:8000/buy?id=itemID&userID=userID  
Will directly buy an item without adding it to cart.  
No request body.  
Attach ``<token>`` to request Headers.  
Returned Body:
```
"Order placed successfully"
```

