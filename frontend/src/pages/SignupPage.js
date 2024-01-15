import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import { useForm } from "react-hook-form";
import Container from "react-bootstrap/Container";
import Button from "react-bootstrap/Button";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Form from "react-bootstrap/Form";
import FloatingLabel from "react-bootstrap/FloatingLabel";
import Alert from "react-bootstrap/Alert";
import { LinkContainer } from "react-router-bootstrap";
import useAuth from "../hooks/useAuth";
import { SIGNUP_URL } from "../utils/routes";

const Signup = () => {
  const { auth } = useAuth;
  const [errMsg, setErrMsg] = useState("");
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm({
    mode: "onBlur",
    defaultValues: {
      dateOfBirth: new Date(
        new Date().getFullYear() - 13,
        new Date().getMonth(),
        new Date().getDate()
      )
        .toISOString()
        .split("T")[0],
    },
    criteriaMode: "all",
  });

  useEffect(() => {
    if (auth) {
      navigate("/profile", { replace: true });
    }
  }, [auth, navigate]);

  const onSubmit = async (data) => {
    try {
      await axios.post(SIGNUP_URL, JSON.stringify(data), {
        headers: { "Content-Type": "application/json" },
        withCredentials: true,
      });

      navigate("/profile", { replace: true });
    } catch (err) {
      if (!err?.response) {
        setErrMsg("No Server Response");
      } else if (err.response?.status === 400) {
        const data = err.response.data.slice(0, -1);
        if (data === "nickname") {
          setErrMsg("The nickname has already been taken");
        } else if (data === "email") {
          setErrMsg("Please use another email address");
        } else if (data === "password") {
          setErrMsg(
            "Your password should have at least one lowercase and one uppercase letter, a number and a symbol"
          );
        }
      } else {
        setErrMsg("Internal Server Error");
      }
    }
  };

  return (
    <Container fluid="md">
      {errMsg && (
        <Alert variant="danger" className="text-center">
          {errMsg}
        </Alert>
      )}
      <Form onSubmit={handleSubmit(onSubmit)}>
        <h1 className="text-center">Sign up for FREE to start networking</h1>
        <Row>
          <Col>
            <FloatingLabel
              className="mb-3"
              controlId="floatingEmail"
              label="Email address"
            >
              <Form.Control
                type="email"
                placeholder="Please enter your email address"
                {...register("email", {
                  required: "Please enter your email address",
                  pattern: {
                    value:
                      /^[A-Z0-9][A-Z0-9._%+-]{0,63}@(?:[A-Z0-9-]{1,63}\.){1,15}[A-Z]{2,63}$/i,
                    message:
                      "The email address should be in form of example@example.com",
                  },
                })}
              />
              {errors.email && (
                <Alert variant="danger">{errors.email.message}</Alert>
              )}
            </FloatingLabel>
          </Col>
        </Row>

        <Row>
          <Col xs="12" md>
            <FloatingLabel
              className="mb-3"
              controlId="floatingPassword"
              label="Password"
            >
              <Form.Control
                type="password"
                placeholder="Enter your password"
                {...register("password", {
                  required: "Please enter your password",
                  minLength: {
                    value: 8,
                    message:
                      "The password should be at least 8 characters long",
                  },
                  pattern: {
                    value: /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])/,
                    message:
                      "The password should have at least one lowercase and one uppercase letter, a number and a symbol",
                  },
                })}
              />
              {errors.password && (
                <Alert variant="danger">{errors.password.message}</Alert>
              )}
            </FloatingLabel>
          </Col>
          <Col xs="12" md>
            <FloatingLabel
              className="mb-3"
              controlId="floatingConfirmPassword"
              label="Confirm password"
            >
              <Form.Control
                type="password"
                placeholder="Confirm password"
                {...register("confirmPassword", {
                  exclude: true,
                  required: "Please enter your password again",
                  validate: (value) =>
                    value === watch("password") || "The passwords do not match",
                })}
              />
              {errors.confirmPassword && (
                <Alert variant="danger">{errors.confirmPassword.message}</Alert>
              )}
            </FloatingLabel>
          </Col>
        </Row>

        <Row>
          <Col xs="12" md>
            <FloatingLabel
              className="mb-3"
              controlId="floatingFirstName"
              label="First name"
            >
              <Form.Control
                placeholder="Enter your first name"
                {...register("firstName", {
                  required: "Please enter your first name",
                })}
              />
              {errors.firstName && (
                <Alert variant="danger">{errors.firstName.message}</Alert>
              )}
            </FloatingLabel>
          </Col>
          <Col xs="12" md>
            <FloatingLabel
              className="mb-3"
              controlId="floatingLastName"
              label="Last name"
            >
              <Form.Control
                placeholder="Enter your last name"
                {...register("lastName", {
                  required: "Please enter your last name",
                })}
              />
              {errors.lastName && (
                <Alert variant="danger">{errors.lastName.message}</Alert>
              )}
            </FloatingLabel>
          </Col>
        </Row>

        <Row>
          <Col xs="12" md>
            <FloatingLabel
              className="mb-3"
              controlId="floatingDateOfBirth"
              label="Date of birth"
            >
              <Form.Control
                type="date"
                {...register("dateOfBirth", {
                  required: "Please enter your birth date",
                  validate: (value) =>
                    new Date(value) <
                      new Date(
                        new Date().getFullYear() - 13,
                        new Date().getMonth(),
                        new Date().getDate()
                      ) || "You must be 13 years of age or older to sign up",
                })}
              />
              {errors.dateOfBirth && (
                <Alert variant="danger">{errors.dateOfBirth.message}</Alert>
              )}
            </FloatingLabel>
          </Col>
          <Col xs="12" md>
            <FloatingLabel
              className="mb-3"
              controlId="floatingNickname"
              label="Nickname (optional)"
            >
              <Form.Control
                placeholder="Enter your nickname"
                {...register("nickname", {
                  maxLength: {
                    value: 32,
                    message:
                      "A nickname should not be longer than 32 characters long",
                  },
                  pattern: {
                    value: /^[a-zA-Z0-9._ ]{0,32}$/,
                    message:
                      "A nickname can only contain letters, numbers, spaces, dots (.) and underscores (_)",
                  },
                })}
              />
              {errors.nickname && (
                <Alert variant="danger">{errors.nickname.message}</Alert>
              )}
            </FloatingLabel>
          </Col>
        </Row>
        <Row>
          <Col>
            <FloatingLabel
              className="mb-3"
              controlId="floatingAbout"
              label="About you (optional)"
            >
              <Form.Control
                as="textarea"
                placeholder="Write something about yourself"
                {...register("about")}
              />
            </FloatingLabel>
          </Col>
        </Row>
        <Row>
          <Col>
            <Col as={Button} xs="12" variant="primary" type="submit">
              Sign Up
            </Col>
          </Col>
        </Row>
      </Form>

      <Row className="justify-content-center">
        <Col xs="12" md="3" className="text-center mt-3">
          Already have an account?
          <LinkContainer className="mx-auto" to={`/login`}>
            <Col as={Button} xs="12" variant="success">
              Sign in
            </Col>
          </LinkContainer>
        </Col>
      </Row>
    </Container>
  );
};

export default Signup;
