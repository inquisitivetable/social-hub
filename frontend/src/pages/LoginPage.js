import React, { useState, useEffect } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import axios from "axios";
import useAuth from "../hooks/useAuth";
import {
  Button,
  Row,
  Col,
  Form,
  Container,
  FloatingLabel,
  Alert,
} from "react-bootstrap";
import { LinkContainer } from "react-router-bootstrap";
import { LOGIN_URL } from "../utils/routes";

const Login = () => {
  const { auth, setAuth } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const from =
    location.state?.from?.pathname !== "/logout"
      ? location.state?.from?.pathname || "/profile"
      : "/profile";

  const [formData, setFormData] = useState({
    username: "",
    password: "",
  });
  const [errMsg, setErrMsg] = useState("");

  useEffect(() => {
    if (auth) {
      navigate("/profile", { replace: true });
    }
  }, [auth, navigate]);

  const handleChange = (event) => {
    const { name, value } = event.target;

    setFormData((prevFormData) => {
      return {
        ...prevFormData,
        [name]: value,
      };
    });
  };

  useEffect(() => {
    setErrMsg("");
  }, [formData]);

  const handleSubmit = async (event) => {
    event.preventDefault();

    try {
      await axios.post(LOGIN_URL, JSON.stringify(formData), {
        headers: { "Content-Type": "application/json" },
        withCredentials: true,
      });

      setAuth(true);
      setFormData({
        username: "",
        password: "",
      });

      navigate(from, { replace: true });
    } catch (err) {
      if (!err?.response) {
        setErrMsg("No Server Response");
      } else if (err.response?.status === 400) {
        setErrMsg("Missing username or password");
      } else if (err.response?.status === 401) {
        setErrMsg("Wrong username or password");
      } else {
        setErrMsg("Login Failed");
      }
    }
  };

  return (
    <Container>
      <Row className="justify-content-center">
        <Col sm="6" className="text-center">
          {errMsg && <Alert variant="danger">{errMsg}</Alert>}
        </Col>
      </Row>

      <Row className="justify-content-center">
        <Col sm="6" className="border rounded p-3">
          <Form onSubmit={handleSubmit}>
            <FloatingLabel
              className="mb-3"
              controlId="floatingEmail"
              label="Email address or username"
            >
              <Form.Control
                type="email"
                placeholder="Email address"
                onChange={handleChange}
                name="username"
                value={formData.username}
                required
                autoFocus
              />
            </FloatingLabel>
            <FloatingLabel
              controlId="floatingPassword"
              className="mb-3"
              label="Password"
            >
              <Form.Control
                type="password"
                placeholder="Password"
                onChange={handleChange}
                name="password"
                value={formData.password}
                required
              />
            </FloatingLabel>
            <Col as={Button} xs="12" type="submit">
              Sign In
            </Col>
          </Form>
        </Col>
      </Row>

      <Row className="justify-content-center">
        <Col sm="6" className="text-center mt-3">
          <LinkContainer className="mx-auto" to={`/signup`}>
            <Col as={Button} xs="12" variant="success">
              Create new account
            </Col>
          </LinkContainer>
        </Col>
      </Row>
    </Container>
  );
};

export default Login;
