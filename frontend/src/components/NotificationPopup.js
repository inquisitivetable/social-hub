import React, { useState } from "react";
import Notification from "../components/Notification";
import { Toast, CloseButton } from "react-bootstrap";

const NotificationPopup = ({ notification, onPopupClose }) => {
  const [show, setShow] = useState(true);

  return (
    <Toast
      className="d-none d-md-flex my-auto position-absolute justify-content-center mx-auto top-100"
      bg="info-subtle"
      autohide
      show={show}
      onClose={() => {
        setShow(false);
        onPopupClose();
      }}
    >
      <Toast.Body>
        <Notification notification={notification} popup={true} />
        <span className="end-0 top-0 me-1 mt-1 position-absolute">
          <CloseButton
            onClick={() => {
              setShow(false);
              onPopupClose();
            }}
          />
        </span>
      </Toast.Body>
    </Toast>
  );
};

export default NotificationPopup;
