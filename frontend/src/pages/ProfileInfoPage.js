import React, { useState, useEffect } from "react";
import { useParams, useOutletContext, useNavigate } from "react-router-dom";
import axios from "axios";
import useWebSocketConnection from "../hooks/useWebSocketConnection";
import {
  PROFILE_URL,
  USER_FOLLOWING_URL,
  USER_FOLLOWERS_URL,
  USER_POSTS_URL,
} from "../utils/routes";
import ImageHandler from "../utils/imageHandler.js";
import FeedPosts from "../components/FeedPosts.js";
import { Container, Row, Col, Button, Alert } from "react-bootstrap";
import GenericUserList from "../components/GenericUserList";
import GenericModal from "../components/GenericModal";
import { BirthdayConverter, LongDate } from "../utils/datetimeConverters";

const ProfileInfo = () => {
  const [user, setUser] = useState({});
  const [isFollowed, setIsFollowed] = useState(false);
  const { id } = useParams();
  const { socketUrl } = useOutletContext();
  const { sendJsonMessage } = useWebSocketConnection(socketUrl);
  const [errMsg, setErrMsg] = useState("");
  const navigate = useNavigate();

  const handleFollow = () => {
    sendJsonMessage({
      type: "follow_request",
      data: { id: user.id },
    });
  };

  const handleUnfollow = () => {
    sendJsonMessage({
      type: "unfollow",
      data: { id: user.id },
    });
    setIsFollowed(false);
  };

  useEffect(() => {
    const loadUser = async () => {
      try {
        await axios
          .get(PROFILE_URL + `/${id}`, {
            withCredentials: true,
          })
          .then((response) => {
            if (response?.data?.user?.isOwnProfile === true) {
              navigate("/profile", { replace: true });
            } else {
              setUser(response?.data?.user);
              setIsFollowed(response?.data?.user?.isFollowed);
            }
          });
      } catch (err) {
        if (!err?.response) {
          setErrMsg("No Server Response");
        } else {
          setErrMsg("Internal Server Error");
        }
      }
    };
    loadUser();
    // eslint-disable-next-line
  }, [id, isFollowed]);

  const userList = (following) =>
    following ? (
      <GenericUserList url={USER_FOLLOWING_URL + id} />
    ) : (
      <GenericUserList url={USER_FOLLOWERS_URL + id} />
    );

  const image = ImageHandler(user?.imagePath, "defaultuser.jpg", "profile-img");

  return (
    <Container fluid>
      {errMsg && (
        <Alert variant="danger" className="text-center">
          {errMsg}
        </Alert>
      )}
      {user && (
        <>
          <Row className="gap-2">
            <Col sm>
              <div className="profile-img">{image}</div>
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
                  <h1 className="display-5">{user.nickname}</h1>
                </Col>
                <Row className="gap-2">
                  <div className="d-flex gap-2 justify-content-center">
                    <Col
                      lg="5"
                      as={Button}
                      disabled={isFollowed}
                      onClick={handleFollow}
                    >
                      Follow
                    </Col>

                    <Col
                      lg="5"
                      as={Button}
                      disabled={!isFollowed}
                      onClick={handleUnfollow}
                    >
                      Unfollow
                    </Col>
                  </div>
                </Row>
              </Row>
            </Col>
          </Row>

          {(user?.isPublic || user?.isFollowed) && (
            <>
              <Row>
                <Col>
                  <div className="mt-3 mb-3">{user.about}</div>
                </Col>
              </Row>
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
              <Row className="d-grip gap-2">
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
                    <FeedPosts url={USER_POSTS_URL + id} />
                  </GenericModal>
                </Col>
              </Row>
            </>
          )}
        </>
      )}
    </Container>
  );
};

export default ProfileInfo;
