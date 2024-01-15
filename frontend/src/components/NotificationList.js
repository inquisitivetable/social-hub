import React, { useEffect, useRef } from "react";
import Notification from "../components/Notification";
import { ListGroup } from "react-bootstrap";

const NotificationList = ({ notifications, setToggle, setNotifications }) => {
  const ref = useRef(null);

  const handleNotificationClose = (id) => {
    setNotifications((prevNotifications) =>
      prevNotifications.filter((notification) => {
        return notification?.notification_id !== id;
      })
    );
  };

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (ref.current && !ref.current.contains(event.target)) {
        setToggle(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
    // eslint-disable-next-line
  }, [ref]);

  const renderedNotifications = notifications.map((notification, index) => (
    <ListGroup.Item variant="flush" action key={index}>
      <Notification
        notification={notification}
        onClose={handleNotificationClose}
      />
    </ListGroup.Item>
  ));

  return (
    <>
      {notifications.length > 0 && (
        <ListGroup ref={ref} className="scroll position-fixed">
          {renderedNotifications}
        </ListGroup>
      )}
    </>
  );
};

export default NotificationList;
