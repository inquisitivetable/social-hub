import { useState, useEffect } from "react";
import { Messenger } from "react-bootstrap-icons";
import Chat from "../components/Chat";
import { Offcanvas } from "react-bootstrap";

const NavbarChat = () => {
  const [show, setShow] = useState(false);
  const [width, setWidth] = useState(window.innerWidth);
  const breakpoint = 768;
  const [newMessages, setNewMessages] = useState(false);

  const handleShow = () => {
    setShow(true);
    newMessages && setNewMessages(false);
  };
  const handleClose = () => setShow(false);

  useEffect(() => {
    const handleWindowResize = () => setWidth(window.innerWidth);
    window.addEventListener("resize", handleWindowResize);

    return () => window.removeEventListener("resize", handleWindowResize);
  }, []);

  const chat = width < breakpoint && (
    <Chat newMessages={newMessages} setNewMessages={setNewMessages} />
  );

  const colour = newMessages ? "red" : "black";

  return (
    <>
      <Messenger size={30} color={colour} onClick={handleShow} />

      <Offcanvas show={show} onHide={handleClose} responsive="md">
        <Offcanvas.Header className="ms-auto" closeButton />
        <Offcanvas.Body>{chat}</Offcanvas.Body>
      </Offcanvas>
    </>
  );
};

export default NavbarChat;
