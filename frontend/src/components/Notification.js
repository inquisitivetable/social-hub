import { WS_URL } from "../utils/routes";
import useWebSocketConnection from "../hooks/useWebSocketConnection";
import { LinkContainer } from "react-router-bootstrap";
import { Button, Col, Row, Stack } from "react-bootstrap";
import { ShortDatetime } from "../utils/datetimeConverters";
import { CheckLg, XLg } from "react-bootstrap-icons";

const Notification = ({ notification, onClose, popup }) => {
  const { sendJsonMessage } = useWebSocketConnection(WS_URL);

  const handleReject = () => {
    const msg = {
      type: "response",
      data: { id: notification?.notification_id, reaction: false },
    };

    sendJsonMessage(msg);
    onClose(notification?.notification_id);
  };

  const handleAccept = () => {
    const msg = {
      type: "response",
      data: { id: notification?.notification_id, reaction: true },
    };

    sendJsonMessage(msg);
    onClose(notification?.notification_id);
  };

  const acceptButton = (
    <CheckLg as={Button} size={23} color="green" onClick={handleAccept} />
  );

  const rejectButton = (
    <XLg as={Button} size={23} color="red" onClick={handleReject} />
  );

  const buttons = !popup && (
    <>
      <Col xs="auto" className="d-flex align-items-center">
        <Stack direction="horizontal" gap="2">
          {acceptButton}
          {rejectButton}
        </Stack>
      </Col>
    </>
  );

  const followRequestNotification = (
    <>
      <LinkContainer to={`/profile/${notification?.sender_id}`}>
        <span>
          <strong>{notification?.sender_name}</strong>
        </span>
      </LinkContainer>{" "}
      wants to follow you
    </>
  );

  const groupInviteNotification = (
    <>
      <LinkContainer to={`/profile/${notification?.sender_id}`}>
        <span>
          <strong>{notification?.sender_name}</strong>
        </span>
      </LinkContainer>{" "}
      invites you to join the group{" "}
      <LinkContainer to={`/groups/${notification?.group_id}`}>
        <span>
          <strong>{notification?.group_name}</strong>
        </span>
      </LinkContainer>
    </>
  );

  const groupRequestNotification = (
    <>
      <LinkContainer to={`/profile/${notification?.sender_id}`}>
        <span>
          <strong>{notification?.sender_name}</strong>
        </span>
      </LinkContainer>{" "}
      wants to join your group{" "}
      <LinkContainer to={`/groups/${notification?.group_id}`}>
        <span>
          <strong>{notification?.group_name}</strong>
        </span>
      </LinkContainer>
    </>
  );

  const eventNotification = (
    <>
      <LinkContainer to={`/event/${notification?.event_id}`}>
        <span>
          <strong>{notification?.event_name}</strong>
        </span>
      </LinkContainer>{" "}
      is going to take place on {ShortDatetime(notification?.event_datetime)}
    </>
  );

  const notificationTemplate = (content) => {
    return (
      <Row>
        <Col>{content}</Col>
        {buttons}
      </Row>
    );
  };

  const notificationMessage = () => {
    switch (notification?.notification_type) {
      case "follow_request":
        return notificationTemplate(followRequestNotification);
      case "group_invite":
        return notificationTemplate(groupInviteNotification);
      case "group_request":
        return notificationTemplate(groupRequestNotification);
      case "event_invite":
        return notificationTemplate(eventNotification);
      default:
        break;
    }
  };

  return notificationMessage();
};

export default Notification;
