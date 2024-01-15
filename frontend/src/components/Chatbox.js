import React, { useState, useEffect, useCallback, useRef } from "react";
import { WS_URL } from "../utils/routes";
import useWebSocketConnection from "../hooks/useWebSocketConnection";
import ChatMessage from "../components/ChatMessage";
import { LinkContainer } from "react-router-bootstrap";
import ImageHandler from "../utils/imageHandler";
import InfiniteScroll from "react-infinite-scroller";
import Picker from "emoji-picker-react";
import { EmojiSmileFill, Send } from "react-bootstrap-icons";
import { CloseButton, Form, Button, Stack, Card } from "react-bootstrap";

const Chatbox = ({
  toggleChat,
  chat,
  user,
  updateChatlist,
  resetUnreadCount,
}) => {
  const { sendJsonMessage, lastJsonMessage } = useWebSocketConnection(WS_URL);
  const pickerRef = useRef(null);
  const messageboxRef = useRef(null);
  const [hasMoreMessages, setHasMoreMessages] = useState(true);
  const [scrollToBottomNeeded, setScrollToBottomNeeded] = useState(false);
  const [showPicker, setShowPicker] = useState(false);
  const [loading, setLoading] = useState(false);
  const [messageHistory, setMessageHistory] = useState([]);
  const [message, setMessage] = useState({
    type: "message",
    data: {
      body: "",
    },
  });

  useEffect(() => {
    setMessageHistory([]);
    setHasMoreMessages(true);
  }, [chat]);

  useEffect(() => {
    switch (lastJsonMessage?.type) {
      case "message_history":
        if (lastJsonMessage?.data.length > 0) {
          setMessageHistory((prevMessageHistory) => [
            ...lastJsonMessage?.data,
            ...prevMessageHistory,
          ]);
        }

        if (lastJsonMessage?.data.length < 10) {
          setHasMoreMessages(false);
        }

        setLoading(false);
        break;
      case "message":
        if (
          (lastJsonMessage?.data?.sender_id === chat.user_id &&
            lastJsonMessage?.data?.group_id === 0) ||
          lastJsonMessage?.data?.recipient_id === chat.user_id ||
          lastJsonMessage?.data?.group_id === chat.group_id
        ) {
          setMessageHistory((prevMessageHistory) => [
            ...prevMessageHistory,
            lastJsonMessage?.data,
          ]);
        }
        break;
      default:
        break;
    }
    // eslint-disable-next-line
  }, [lastJsonMessage]);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (pickerRef.current && !pickerRef.current.contains(event.target)) {
        setShowPicker(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [pickerRef]);

  useEffect(() => {
    if (scrollToBottomNeeded) {
      scrollToBottom();
      setScrollToBottomNeeded(false);
    }
  }, [scrollToBottomNeeded]);

  const onEmojiClick = (event) => {
    const emoji = event.emoji;
    setMessage((prevMessage) => ({
      ...prevMessage,
      data: { body: prevMessage.data.body + emoji },
    }));
    setShowPicker(false);
  };

  const handleScrolling = () => {
    const lastMessage = messageHistory[messageHistory.length - 1]?.id;

    if (
      messageboxRef?.current?.scrollHeight -
        messageboxRef?.current?.clientHeight <=
      messageboxRef?.current?.scrollTop + 1
    ) {
      sendJsonMessage({
        type: "messages_read",
        data: {
          id: chat.user_id,
          group_id: chat.group_id,
          last_message: lastMessage,
        },
      });

      resetUnreadCount([chat.user_id, chat.group_id]);
    }
  };

  const image =
    chat?.user_id > 0
      ? ImageHandler(chat?.avatar_image, "defaultuser.jpg", "chatbox-img")
      : ImageHandler(chat?.avatar_image, "defaultgroup.png", "chatbox-img");

  const loadMessages = useCallback(async () => {
    if (loading) {
      return;
    }

    setLoading(true);

    const offset = messageHistory.length > 0 ? messageHistory[0].id : 0;

    sendJsonMessage({
      type: "request_message_history",
      data: {
        id: chat.user_id,
        group_id: chat.group_id,
        last_message: offset,
      },
    });
    // eslint-disable-next-line
  }, [loading, hasMoreMessages, messageHistory]);

  const closeChat = () => {
    toggleChat(null);
  };

  const handleChange = (event) => {
    const { value } = event.target;

    setMessage((prevMessage) => {
      return {
        ...prevMessage,
        data: { body: value },
      };
    });
  };

  const renderedMessages = messageHistory?.map((msg, index) => {
    switch (msg.sender_id) {
      case user:
        return <ChatMessage key={index} msg={msg} own={true} />;
      default:
        return (
          <ChatMessage
            key={index}
            msg={{
              ...msg,
              sender_name:
                messageHistory[index - 1]?.sender_id === msg.sender_id
                  ? ""
                  : msg.sender_name,
            }}
          />
        );
    }
  });

  const handleSubmit = (event) => {
    event.preventDefault();
    if (!message?.data?.body) {
      return;
    }
    let msg = {
      ...message,
      data: {
        ...message.data,
        sender_id: user,
        recipient_id: chat?.user_id,
        group_id: chat?.group_id,
      },
    };

    sendJsonMessage(msg);

    setMessageHistory((prevMessageHistory) => [
      ...prevMessageHistory,
      {
        ...msg.data,
        timestamp: new Date().toISOString(),
      },
    ]);

    updateChatlist([
      chat?.user_id ? chat?.user_id : 0,
      chat?.group_id ? chat?.group_id : 0,
    ]);

    setMessage({ ...message, data: { body: "" } });
    setScrollToBottomNeeded(true);

    resetUnreadCount([chat.user_id, chat.group_id]);
  };

  const chatName =
    chat?.user_id > 0 ? (
      <LinkContainer to={`/profile/${chat.user_id}`}>
        <p className="my-auto">{chat.name}</p>
      </LinkContainer>
    ) : (
      <LinkContainer to={`/groups/${chat.group_id}`}>
        <p className="my-auto">{chat.name}</p>
      </LinkContainer>
    );

  const scrollToBottom = () => {
    messageboxRef.current.scrollTop = messageboxRef.current.scrollHeight;
  };

  const chatbox = (
    <Card>
      <Card.Header>
        <Stack direction="horizontal">
          <div className="me-auto">{image}</div>
          {chatName}
          <CloseButton
            className="ms-auto align-self-center"
            onClick={closeChat}
          />
        </Stack>
      </Card.Header>

      <Card.Body
        className="message-history"
        ref={messageboxRef}
        onScroll={handleScrolling}
      >
        <InfiniteScroll
          pageStart={0}
          isReverse={true}
          loadMore={loadMessages}
          hasMore={hasMoreMessages}
          useWindow={false}
        >
          {renderedMessages}
        </InfiniteScroll>
      </Card.Body>

      <Card.Footer>
        <Form onSubmit={handleSubmit}>
          <Stack direction="horizontal" gap={2}>
            <Form.Control
              placeholder="Message"
              onChange={handleChange}
              name="message"
              value={message.data.body}
              autoFocus
            />
            <EmojiSmileFill
              color="blue"
              size={38}
              onClick={() => setShowPicker(true)}
            />

            {showPicker && (
              <div className="picker-container" ref={pickerRef}>
                <Picker
                  //lazyLoad={true}
                  className="EmojiPicker"
                  searchDisabled={true}
                  skinTonesDisabled={true}
                  previewConfig={{ showPreview: false }}
                  categories={["smileys_people"]}
                  width={282}
                  height={351}
                  onEmojiClick={onEmojiClick}
                />
              </div>
            )}

            <Button type="submit">
              <Send />
            </Button>
          </Stack>
        </Form>
      </Card.Footer>
    </Card>
  );

  return <div className="chatbox">{chatbox}</div>;
};

export default Chatbox;
