# Go Marketplace
An API that simulates backend behavior of a marketplace.

<a name="api-endpoints-overview"></a>

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
You can view the collection here.  
https://www.postman.com/czhang35-b1391359-892575/workspace/go-market-workspace/collection/47550840-78e3d3a8-f6b2-48c0-bcc5-52f213f78c19?action=share&creator=47550840  

## 📌 API Endpoints Overview

| **Feature**              | **Method** | **Endpoint**                           | **Description**                                  |
|--------------------------|------------|----------------------------------------|--------------------------------------------------|
| **User Authentication**  |            |                                        |                                                  |
| Sign Up                  | `POST`     | [/users/signup](#sign-up-post)         | Register a new user                              |
| Verify                   | `POST`     | [/users/verify](#verify-post)          | Verify account using code                        |
| Login                    | `POST`     | [/users/login](#login-post)            | Log in and get token                             |
| **Marketplace**          |            |                                        |                                                  |
| List Item                | `POST`     | [/users/listItem](#list-an-item)       | Add a new product                                |
| View All Items           | `GET`      | [/users/view](#view-all-market-items-get) | Fetch all available items                     |
| Search Item              | `GET`      | [/users/search?name=](#search-for-item-get) | Search items by name                        |
| **Cart Management**      |            |                                        |                                                  |
| Add to Cart              | `GET`      | [/add](#add-item-to-cart-get)          | Add item to user’s cart                          |
| List Cart                | `GET`      | [/list](#list-items-in-cart-get)       | Get user’s cart items                            |
| Remove from Cart         | `GET`      | [/remove](#remove-item-from-cart-get)  | Remove item from cart                            |
| Checkout Cart            | `GET`      | [/checkout](#cart-checkout-get)        | Checkout all items in cart                       |
| Instant Buy              | `GET`      | [/buy](#buy-item-instantly-get)        | Buy item instantly without adding to cart        |
| **Address Management**   |            |                                        |                                                  |
| Add Address              | `POST`     | [/addressadd](#add-address-post)       | Add home or work address                         |
| Delete Address           | `GET`      | [/addressdel](#delete-address-get)     | Delete all addresses                             |
| Edit Home Address        | `PUT`      | [/addresshomeedit](#edit-home-address-put) | Update home address                          |
| Edit Work Address        | `PUT`      | [/addressworkedit](#edit-work-address-put) | Update work address                          |
| **Chat & Messaging**     |            |                                        |                                                  |
| Start Chat               | `POST`     | [/chats](#start-chat-post)             | Start chat with another user                     |
| List Chats               | `GET`      | [/chats](#list-all-chats-get)          | Get all chats for user                           |
| Send Message             | `POST`     | [/chats/:chatID/messages](#send-message-post) | Send message in a chat                    |
| List Messages            | `GET`      | [/chats/:chatID/messages](#list-messages-get) | List messages in a chat                   |
| Mark Messages Read       | `POST`     | [/chats/:chatID/read](#read-message-post) | Mark messages as read                         |
| **Product Reviews**      |            |                                        |                                                  |
| List Reviews             | `GET `     | [/products/:productID/reviews](#list-reviews-get) | List reviews of a product             |
| make review              | `POST`     | [/products/:productID/reviews](#make-review-post) | Make review for a product             |

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
    "name": "pen",
    "price": 9.99,
    "img": "pencil.png",
    "description": "black pen 0.5mm with replacable ink"
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
        "ID": "68c34222df9bb0af3283176a",
        "name": "pencil",
        "price": 5,
        "img": "pencil.png",
        "description": null,
        "ratingAvg": 0,
        "ratingCnt": 0,
        "ratingSum": 0
    },
    {
        "ID": "68c342b0df9bb0af3283176c",
        "name": "pen",
        "price": 9.99,
        "img": "pencil.png",
        "description": "black pen 0.5mm with replacable ink",
        "ratingAvg": 0,
        "ratingCnt": 0,
        "ratingSum": 0
    }
]
```
### Search for item (GET)
http://localhost:8000/users/search?name=pen  
No request body. 
Returned Body:
```
[
    {
        "ID": "68c34222df9bb0af3283176a",
        "name": "pencil",
        "price": 5,
        "img": "pencil.png",
        "description": null,
        "ratingAvg": 0,
        "ratingCnt": 0,
        "ratingSum": 0
    },
    {
        "ID": "68c342b0df9bb0af3283176c",
        "name": "pen",
        "price": 9.99,
        "img": "pencil.png",
        "description": "black pen 0.5mm with replacable ink",
        "ratingAvg": 0,
        "ratingCnt": 0,
        "ratingSum": 0
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

### Start chat (POST)
http://localhost:8000/chats?userID=userID  
Start chat with another user.  
Attach ``<token>`` to request Headers.  
Request Body:
```
{ 
    "peerId": <Other UserID> 
}
```
Returned Body:
```
{
    "id": <Chat ID>,
    "members": [
        "User1 ID",
        "User2 ID"
    ],
    "createdAt": "2025-09-11T00:51:03.506648081-04:00",
    "updatedAt": "2025-09-11T00:51:03.506648081-04:00",
    "lastMessage": null,
    "unreadBy": {
        "User1 ID": 0,
        "User2 ID": 0
    }
}
```

### List all chats (GET)
http://localhost:8000/chats?userID=userID  
No request body.  
Attach ``<token>`` to request Headers.  
Returned Body:
```
[
    {
        "id": <Chat ID>,
        "members": [
            "User1 ID",
            "User2 ID"
        ],
        "createdAt": "2025-09-11T04:51:03.506Z",
        "updatedAt": "2025-09-11T04:51:03.506Z",
        "lastMessage": null,
        "unreadBy": {
            "User1 ID": 0,
            "User2 ID": 0
        }
    }
]
```

### Send message (POST)
http://localhost:8000/chats/ChatID/messages?userID=UserID  
Attach ``<token>`` to request Headers.  
Request Body:
```
{ 
    "text": "hello" 
}
```
Returned Body:
```
{
    "id": <Message ID>,
    "chatId": <Chat ID>,
    "senderId": <User1 ID>,
    "text": "hello",
    "createdAt": "2025-09-11T00:58:12.721520423-04:00",
    "readBy": [
        <User2 ID>
    ]
}
```
If call list chats after sending a message, you will see the receiver's unreadby increase by 1 since there is now a newly sent message.  

### List messages (GET)
http://localhost:8000/chats/chatID/messages?userID=userID&limit=10  
Displays the last 10 messages in a chat.  
No request body.  
Attach ``<token>`` to request Headers.  
Returned Body:
```
[
    {
        "id": <Message ID>,
        "chatId": <Chat ID>,
        "senderId": <Sender ID>,
        "text": "hello",
        "createdAt": "2025-09-11T04:58:12.721Z",
        "readBy": [
            <User1 ID>
        ]
    }
]
```

### Read message (POST)
http://localhost:8000/chats/chatID/read?userID=user2ID  
Here we simulate user2, the receiver, reads the message.  
Attach user2's ``<token>`` to request Headers.  
Request Body:
```
{ 
    "peerId": <User1 ID>
}
```
Returned Body:
```
{
    "status": "ok"
}
```
Now if you run list messages, you will see user2's id being appended to "readBy".  
If you run list all chats, you will see user2's unreadby decrease to 0.  

### List reviews (GET)
http://localhost:8000/products/productID/reviews?limit=5  
Displays 5 reviews.  
No request body.  
Attach ``<token>`` to request Headers.  
Returned Body:
If no review, 
```
null
```
If with review:
```
[
    {
        "id": <review id>,
        "pid": <product id>,
        "uid": <user id>,
        "rating": 4.5,
        "review": "Solid pen, smooth to write and steady build",
        "createdAt": "2025-09-11T22:05:20.659Z",
        "updatedAt": "2025-09-11T22:05:20.659Z"
    }
]
```

### make review (POST)
http://localhost:8000/products/productID/reviews?userID=userID   
Attach ``<token>`` to request Headers.  
Request Body:
```
{
  "rating": 4.5,
  "review": "Solid pen, smooth to write and steady build"
}
```
Returned Body:
```
{
    "status": "created"
}
```
The newly added review will be returned from list reviews requetsed, and the prodcut's ratings will be updated accordingly.  