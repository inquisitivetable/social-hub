import React, { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import axios from "axios";
import AvatarUpdater from "../components/AvatarUpdater.js";
import ImageHandler from "../utils/imageHandler.js";
import {
  Container,
  Row,
  Col,
  Button,
  Form,
  FloatingLabel,
  Alert,
} from "react-bootstrap";
import FeedPosts from "../components/FeedPosts.js";
import {
  PROFILE_URL,
  PROFILE_UPDATE_URL,
  PROFILE_POSTS_URL,
  AVATAR_UPDATER_URL,
  FOLLOWERS_URL,
  FOLLOWING_URL,
} from "../utils/routes.js";
import GenericUserList from "../components/GenericUserList.js";
import GenericModal from "../components/GenericModal.js";
import { LongDate, BirthdayConverter } from "../utils/datetimeConverters.js";

const ProfileEditorPage = () => {
  const [user, setUser] = useState({});
  const [errMsg, setErrMsg] = useState("");
  const values = user;
  const {
    register,
    handleSubmit,
    formState: { errors, isDirty },
  } = useForm({
    mode: "onBlur",
    values,
    criteriaMode: "all",
  });

  const image = ImageHandler(user?.imagePath, "defaultuser.jpg", "profile-img");

  const loadUser = async () => {
    await axios
      .get(PROFILE_URL, {
        withCredentials: true,
      })
      .then((response) => {
        setUser(response.data.user);
      });
  };

  useEffect(() => {
    loadUser();
  }, []);

  const onSubmit = async (data) => {
    try {
      await axios.post(PROFILE_UPDATE_URL, JSON.stringify(data), {
        withCredentials: true,
        headers: { "Content-Type": "application/json" },
      });
    } catch (err) {
      if (!err?.response) {
        setErrMsg("No Server Response");
      } else if (err.response?.status > 200) {
        setErrMsg("Internal Server Error");
      }
    }
  };

  const handleAvatarUpdate = () => {
    loadUser();
  };

  const userList = (following) =>
    following ? (
      <GenericUserList url={FOLLOWING_URL} />
    ) : (
      <GenericUserList url={FOLLOWERS_URL} />
    );

  return (
    <Container fluid>
      {errMsg && (
        <Alert variant="danger" className="text-center">
          {errMsg}
        </Alert>
      )}
      {user && (
        <>
          <Form onSubmit={handleSubmit(onSubmit)}>
            <Row className="gap-2">
              <Col sm>
                <div className="profile-img">
                  {image}
                  <GenericModal buttonText="Upload new image">
                    <AvatarUpdater
                      url={AVATAR_UPDATER_URL}
                      onUploadSuccess={handleAvatarUpdate}
                    />
                  </GenericModal>
                </div>
              </Col>
              <Col sm>
                <Row className="d-grid gap-2">
                  <Col xs="12">
                    <h1>
                      {user.firstName} {user.lastName}
                    </h1>
                  </Col>
                  <Col xs="12">
                    <div>also known as </div>
                    <FloatingLabel
                      className="mb-3"
                      controlId="floatingNickname"
                      label="Nickname"
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
                        <Alert variant="danger">
                          {errors.nickname.message}
                        </Alert>
                      )}
                    </FloatingLabel>
                  </Col>
                  <div className="mb-3">
                    <Form.Check
                      type="checkbox"
                      label="Profile is public"
                      {...register("isPublic")}
                    />
                  </div>
                </Row>
              </Col>
            </Row>

            <Row>
              <Col>
                <FloatingLabel
                  className="mt-3 mb-3"
                  controlId="about"
                  label="About you"
                >
                  <Form.Control
                    as="textarea"
                    className="profile-textarea"
                    placeholder="Write something about yourself"
                    {...register("about")}
                  />
                </FloatingLabel>
              </Col>
              <div className="text-center">
                <Row>
                  <Col>
                    <strong>Email address</strong>
                    <p>{user.email}</p>
                  </Col>
                  <Col>
                    <strong>Profile Type</strong>
                    <p>{user.isPublic ? "Public" : "Private"}</p>
                  </Col>
                </Row>
                <Row>
                  <Col>
                    <strong>Born</strong>
                    <p>{BirthdayConverter(user?.birthday)}</p>
                  </Col>
                  <Col>
                    <strong>Joined</strong>
                    <p>{LongDate(user.createdAt)}</p>
                  </Col>
                </Row>
              </div>
            </Row>

            <div className="d-flex justify-content-center mb-3 mt-3">
              <Button type="submit" disabled={!isDirty}>
                Save changes
              </Button>
            </div>
          </Form>
          <Row className="gap-2">
            <Col>
              <GenericModal buttonText="Following">
                {userList(true)}
              </GenericModal>
            </Col>
            <Col>
              <GenericModal buttonText="Followers">
                {userList(false)}
              </GenericModal>
            </Col>
            <Col>
              <GenericModal buttonText="Posts">
                <FeedPosts url={PROFILE_POSTS_URL} />
              </GenericModal>
            </Col>
          </Row>
        </>
      )}
    </Container>
  );
};

export default ProfileEditorPage;
