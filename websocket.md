# JSON structure for Websocket messages

## 1. BACKEND to FRONTEND

### 1.1 notification

```JSON
{
    "type": "notification",
    "data": {
        "notification_type": "follow_request" || "group_invite" || "group_request" || "event_invite",
        "notification_id": 1, // notification id
        "sender_id": 123,
        "sender_name": "something", // either a username (if exists) or firstname and lastname
        "group_id": 123, // 0 if not group
        "group_name": "something", // empty if not group
        "event_id": 123, // 0 if not event
        "event_name": "something", // empty if not event
        "event_datetime": "2006-01-02T15:04:05Z07:00", // empty if not event
    }
}
```

### 1.2 chatlist

```JSON
{
    "type": "chatlist",
    "data": {
        "userid": 1, //own id
        "user_chatlist" : [
            {
                "user_id": 123,
                "name": "username" || "firstname lastname", // username (if exists) or combined full name
                "timestamp": "2006-01-02T15:04:05Z07:00", // date of last message in the chat if any
                "avatar_image": "link",
                "unread_count": 123, // number of unread messages
            }
        ],
        "group_chatlist" : [
            {
                "group_id": 123,
                "name": "group name",
                "timestamp": "2006-01-02T15:04:05Z07:00", // date of last message in the chat if any
                "avatar_image": "link",
            }
        ]
    }
}
```

### 1.3 message history

```JSON
{
    "type": "message_history",
    "data": {
        "messages" : [{
            "id": 1, //message id
            "sender_id": 123, // not 0 if group (message still has sender)
            "sender_name": "username", // either a  username (if exists) or firstname and lastname
            "recipient_id": 1, // 0 if group
            "recipient_name": 1, // either a username (if   exists) or firstname and lastname && empty if     group
            "group_id": 123, // 0 if user
            "group_name": "name", //empty if user
            "body": "message",
            "timestamp": "2006-01-02T15:04:05Z07:00",
        }]
    }
}
```

## 2. DUPLEX

### 2.1 chat message

```JSON
{
    "type": "message",
    "data": {
        "id": 1, // message id
        "sender_id": 1,
        "sender_name" : "sdfs", // username (if exists) or first name last name
        "avatar_image": "link", // empty if no image
        "recipient_id": 123, // 0 if group chat
        "recipient_name": "somename", // username (if exists) or first name last name
        "group_id": 123, // 0 if private chat
        "group_name": "name", //empty if private chat
        "body": "message content",
        "timestamp": "2006-01-02T15:04:05Z07:00" //won't be sending from frontend, but still need to receive it
    }
}
```

## 3. FRONTEND to BACKEND

### 3.1 request chatlist

```JSON
{
    "type": "request_chatlist",
}
```

### 3.2 request message history

```JSON
{
    "type": "request_message_history",
    "data": {
        "id": 123, //0 if group chat
        "group_id": 123, // 0 if private chat
        "last_message": 123 // 0 if from latest message
    }
}
```

### 3.3 follow request

```JSON
{
    "type": "follow_request",
    "data": {
        "id": 123,
    }
}
```

### 3.4 unfollow

```JSON
{
    "type": "unfollow",
    "data": {
        "id": 123,
    }
}
```

### 3.5 group request - someone wants to join a group you created

```JSON
{
    "type": "group_request", // was group_join, but changed to match follow_request
    "data": {
        "id": 123,
    }
}
```

### 3.6 response - a response to any notification

```JSON
{
    "type": "response",
    "data": {
        "id": 1, // notification id
        "reaction": true || false,
    }
}
```

### 3.6 messages read - indication of open chatbox and scrolldown to last message

```JSON
{
    "type": "messages_read",
    "data": {
        "id": 123, // message id
    }
}
```

