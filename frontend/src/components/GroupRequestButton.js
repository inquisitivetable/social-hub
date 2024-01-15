import React from "react";
import { useOutletContext } from "react-router-dom";
import useWebSocketConnection from "../hooks/useWebSocketConnection";
import Button from "react-bootstrap/Button";

const GroupRequestButton = ({ groupid }) => {
  const { socketUrl } = useOutletContext();
  const { sendJsonMessage } = useWebSocketConnection(socketUrl);

  const handleGroupRequest = () => {
    sendJsonMessage({
      type: "group_request",
      data: { group_id: groupid },
    });
  };

  return <Button onClick={handleGroupRequest}>Join Group</Button>;
};

export default GroupRequestButton;
