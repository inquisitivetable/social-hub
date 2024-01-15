import { useState } from "react";
import GroupSidebar from "../components/GroupSidebar";
import { Offcanvas, Nav } from "react-bootstrap";

const NavbarGroupSidebar = () => {
  const [show, setShow] = useState(false);
  const handleShow = () => setShow(true);
  const handleClose = () => setShow(false);

  return (
    <>
      <Nav.Link onClick={handleShow}>Groups and events</Nav.Link>
      {show && (
        <Offcanvas show={show} onHide={handleClose} responsive="md">
          <Offcanvas.Header className="ms-auto" closeButton />
          <Offcanvas.Body>
            <GroupSidebar />
          </Offcanvas.Body>
        </Offcanvas>
      )}
    </>
  );
};

export default NavbarGroupSidebar;
