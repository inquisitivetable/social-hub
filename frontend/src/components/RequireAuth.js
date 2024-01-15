import { useLocation, Outlet } from "react-router-dom";
import { useEffect, useState } from "react";
import useAuth from "../hooks/useAuth";
import axios from "axios";
import Chat from "./Chat";
import { WS_URL } from "../utils/routes";
import Login from "../pages/LoginPage";
import { Container, Row, Col } from "react-bootstrap";
import GroupSidebar from "../components/GroupSidebar";
import { AUTH_URL } from "../utils/routes";

const RequireAuth = () => {
  const { auth, setAuth } = useAuth();
  const [socketUrl] = useState(WS_URL);
  const [loading, setLoading] = useState(true);
  const [width, setWidth] = useState(window.innerWidth);
  const breakpoint = 768;

  const location = useLocation();

  useEffect(() => {
    const authorisation = async () => {
      try {
        await axios.get(AUTH_URL, {
          withCredentials: true,
        });
        setAuth(true);
      } catch (err) {
        if (!err?.response) {
          setAuth(false);
        } else if (err.response?.status === 401) {
          setAuth(false);
        } else {
          setAuth(false);
        }
      }
      setLoading(false);
    };

    authorisation();
    // eslint-disable-next-line
  }, [location]);

  useEffect(() => {
    const handleWindowResize = () => setWidth(window.innerWidth);
    window.addEventListener("resize", handleWindowResize);

    return () => window.removeEventListener("resize", handleWindowResize);
  }, []);

  return loading ? null : auth ? (
    <Container fluid className="bg-light">
      <Row>
        <Col className="sidebar p-0 d-none d-md-flex" id="group-sidebar" xs="3">
          <GroupSidebar />
        </Col>
        <Col
          xs="12"
          md={{ span: "6", offset: "3" }}
          className="mt-3 mb-3 justify-content-end"
        >
          <Outlet context={{ socketUrl }} />
        </Col>
        {width >= breakpoint && (
          <Col id="chat-sidebar" xs="3" className="sidebar p-0">
            <Chat />
          </Col>
        )}
      </Row>
    </Container>
  ) : (
    <>
      <Login />
    </>
  );
};

export default RequireAuth;
