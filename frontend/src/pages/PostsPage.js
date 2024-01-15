import React, { useState } from "react";
import FeedPosts from "../components/FeedPosts";
import CreatePost from "../components/CreatePost";
import { FEEDPOSTS_URL } from "../utils/routes";
import Col from "react-bootstrap/Col";
import Row from "react-bootstrap/Row";
import Container from "react-bootstrap/esm/Container";
import GenericModal from "../components/GenericModal";

const PostsPage = () => {
  const [reload, setReload] = useState(false);

  const handlePostUpdate = () => {
    setReload(!reload);
  };

  return (
    <Container fluid>
      <Row>
        <Col md="8 mx-auto">
          <GenericModal
            buttonText="Write what's on your mind"
            headerText="Make a post"
          >
            <CreatePost onPostsUpdate={handlePostUpdate} />
          </GenericModal>
        </Col>
      </Row>
      <Row>
        <Col>
          <FeedPosts url={FEEDPOSTS_URL} reload={reload} />
        </Col>
      </Row>
    </Container>
  );
};

export default PostsPage;
