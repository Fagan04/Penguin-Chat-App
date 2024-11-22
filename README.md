
# Penguin Chat App - API Documentation

## **Base URL**

Each service runs on a different port:
- **User Service**: `http://localhost:8080`
- **Notification Service**: `http://localhost:8082`
- **Chat Service**: `http://localhost:8081`

---

## **User Service API**

### **1. POST /login**
- **Description**: Logs in a user and returns a JWT token for session management.
- **Request Body**:
  ```json
  {
    "username": "user_example",
    "password": "password_example"
  }
  ```
- **Responses**:
  - **200 OK**: Login successful. Token is set in the cookie.
  - **400 Bad Request**: Invalid input.
  - **500 Internal Server Error**: Token creation failed.

---

### **2. POST /register**
- **Description**: Registers a new user.
- **Request Body**:
  ```json
  {
    "username": "new_user",
    "email": "new_user@example.com",
    "password": "password_example"
  }
  ```
- **Responses**:
  - **201 Created**: User created successfully.
  - **400 Bad Request**: Invalid input.
  - **500 Internal Server Error**: Failed to create user.

---

## **Notification Service API**

### **1. POST /addNotification**
- **Description**: Adds a new notification for a specific user.
- **Request Body**:
  ```json
  {
    "user_id": "1",
    "message": "Test notification",
    "is_new": true,
    "timestamp": "2024-11-23T01:00:00Z"
  }
  ```
- **Responses**:
  - **200 OK**: Notification added successfully.
  - **400 Bad Request**: Invalid input.
  - **500 Internal Server Error**: Failed to add the notification.

### **2. GET /notifications/{user_id}**
- **Description**: Retrieves new notifications for a specific user.
- **Request Parameters**:
  - `user_id`: The user's unique identifier.

- **Responses**:
  - **200 OK**: List of new notifications.
  - **404 Not Found**: No notifications found for the user.
  - **500 Internal Server Error**: Failed to fetch notifications.

---

## **Chat Service API**

### **1. POST /createChat**
- **Description**: Creates a new chat.
- **Request Body**:
  ```json
  {
    "chatName": "chat_name_example"
  }
  ```
- **Responses**:
  - **201 Created**: Chat created successfully.
  - **400 Bad Request**: Invalid input.
  - **409 Conflict**: Chat with this name already exists.
  - **500 Internal Server Error**: Failed to create chat.

---

### **2. GET /accessChat**
- **Description**: Retrieves all chats that the user is part of.
- **Request Headers**:
  - `User-ID`: User's unique identifier.

- **Responses**:
  - **200 OK**: List of chats the user is a member of.
  - **400 Bad Request**: Invalid user ID.
  - **500 Internal Server Error**: Failed to fetch chats.

---

### **3. POST /sendMessage**
- **Description**: Sends a message to a specific chat and triggers notifications for chat members.
- **Request Body**:
  ```json
  {
    "chatID": 1,
    "message": "Hello, world!"
  }
  ```
- **Request Headers**:
  - `User-ID`: User's unique identifier.

- **Responses**:
  - **200 OK**: Message sent successfully, and notifications triggered.
  - **400 Bad Request**: Invalid input.
  - **500 Internal Server Error**: Failed to send message or notify users.
  - **403 Forbidden**: User is not part of the chat.

---

## **How to Use the API**

1. **User Service**:
  - For logging in or registering users.
  - Logs in a user and issues a token for access to other services.

2. **Chat Service**:
  - Allows users to create chats, access existing chats, and send messages.
  - Requires user authentication via the `User-ID` header for chat-related operations.

3. **Notification Service**:
  - Allows adding new notifications for users and retrieving them.

---

### **Error Codes**

- **200 OK**: The request was successful.
- **201 Created**: Resource was created successfully (e.g., user or chat).
- **400 Bad Request**: The request was malformed or missing required data.
- **403 Forbidden**: The user doesn't have permission to perform the action.
- **404 Not Found**: The requested resource could not be found.
- **409 Conflict**: The resource already exists (e.g., trying to create a duplicate chat).
- **500 Internal Server Error**: An unexpected error occurred on the server.

### Workflow Example:

1. A user **registers** on the User Service by providing a username, email, and password.
2. The user **logs in**, and a JWT token is generated and returned.
3. The user **creates a chat** via the Chat Service, and the chat is registered.
4. The user can **access the chat** and **send messages**.
5. When a new message is sent, the **Notification Service** will create notifications for users in the chat.
6. The user can **fetch new notifications** from the Notification Service.