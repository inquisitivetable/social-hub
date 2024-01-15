import useWebSocketConnection from "../hooks/useWebSocketConnection";
import NotificationList from "../components/NotificationList.js";
import axios from "axios";
import React, { useState, useEffect } from "react";
import { WS_URL, NOTIFICATIONS_URL } from "../utils/routes";
import NotificationPopup from "../components/NotificationPopup";
import { Badge, Row, Col, Alert } from "react-bootstrap";
import { BellFill } from "react-bootstrap-icons";

const NotificationNavbarItem = () => {
  const [errMsg, setErrMsg] = useState("");
  const [toggle, setToggle] = useState(false);
  const [newNotification, setNewNotification] = useState(null);
  const { lastJsonMessage } = useWebSocketConnection(WS_URL);
  const [notifications, setNotifications] = useState([]);

  useEffect(() => {
    if (lastJsonMessage && lastJsonMessage.type === "notification") {
      setNotifications((prevNotifications) => {
        return [lastJsonMessage?.data, ...prevNotifications];
      });
    }
  }, [lastJsonMessage]);

  useEffect(() => {
    const loadNotifications = async () => {
      try {
        await axios
          .get(NOTIFICATIONS_URL, {
            withCredentials: true,
          })
          .then((response) => {
            setNotifications(response.data);
          });
      } catch (err) {
        if (!err?.response) {
          setErrMsg("No Server Response");
        } else {
          setErrMsg("Internal Server Error");
        }
      }
    };

    loadNotifications();
  }, []);

  useEffect(() => {
    const exceptions = ["message", "chatlist", "message_history"];

    if (!exceptions.includes(lastJsonMessage?.type)) {
      setNewNotification(lastJsonMessage?.data);
    }
  }, [lastJsonMessage]);

  const handleToggle = () => {
    setToggle(!toggle);
  };

  const onPopupClose = () => {
    setNewNotification(null);
  };

  const notificationCount = notifications.length;

  return (
    <>
      <Row>
        <Col>
          <BellFill
            size={30}
            onClick={handleToggle}
            color={notificationCount > 0 ? "red" : "black"}
          />
          {notificationCount > 0 && (
            <span className="position-absolute">
              <Badge pill bg="danger">
                {notificationCount}
              </Badge>
            </span>
          )}
        </Col>
      </Row>

      {newNotification && (
        <NotificationPopup
          notification={newNotification}
          onPopupClose={onPopupClose}
        />
      )}
      {errMsg ? (
        <Alert variant="danger" className="text-center">
          {errMsg}
        </Alert>
      ) : (
        toggle && (
          <NotificationList
            notifications={notifications}
            setNotifications={setNotifications}
            setToggle={setToggle}
          />
        )
      )}
    </>
  );
};

export default NotificationNavbarItem;
