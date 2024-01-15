import React, { useEffect, useState } from "react";
import { GROUPFEED_URL, GROUP_PAGE_URL } from "../utils/routes";
import axios from "axios";
import { useParams } from "react-router-dom";
import GroupMembers from "../components/GroupMembers";
import ImageHandler from "../utils/imageHandler";
import FeedPosts from "../components/FeedPosts.js";
import AvatarUpdater from "../components/AvatarUpdater";
import Events from "../components/Events";
import GroupRequestButton from "../components/GroupRequestButton.js";
import CreateGroupPosts from "../components/CreateGroupPosts.js";
import GenericModal from "../components/GenericModal";
import { Alert, Container, Col, Row } from "react-bootstrap";

const GroupPage = () => {
  const [group, setGroup] = useState({});
  const { id } = useParams();
  const [errMsg, setErrMsg] = useState(null);
  const [reload, setReload] = useState(false);

  const handlePostUpdate = () => {
    setReload(!reload);
  };

  const loadGroup = async () => {
    try {
      await axios
        .get(GROUP_PAGE_URL + id, {
          withCredentials: true,
        })
        .then((response) => {
          setGroup(response.data);
        });
    } catch (err) {
      setErrMsg(err.message);
    }
  };

  const handleAvatarUpdate = () => {
    loadGroup();
  };

  useEffect(() => {
    loadGroup();
    //eslint-disable-next-line
  }, [id]);

  const image = ImageHandler(group.imagePath, "defaultgroup.png", "group-img");

  return (
    <>
      {errMsg ? (
        <Alert variant="danger" className="text-center">
          {errMsg}
        </Alert>
      ) : (
        <Container fluid>
          <Col className="text-center">
            <div className="profile-img">
              {image}
              {group?.isCreator && (
                <GenericModal buttonText="Upload new image">
                  <AvatarUpdater
                    url={`${GROUP_PAGE_URL}${id}/avatar`}
                    onUploadSuccess={handleAvatarUpdate}
                  />
                </GenericModal>
              )}
            </div>
            <Row>
              <h1>{group.title}</h1>
            </Row>
          </Col>

          <Row xs="auto" className="justify-content-center">
            {group?.isMember ? (
              <>
                <Col>
                  <Events groupId={+id} />
                </Col>
                <Col>
                  <GroupMembers groupId={+id} />
                </Col>
              </>
            ) : (
              <GroupRequestButton groupid={+id} />
            )}
          </Row>

          <Row className="mt-3 mb-3">
            <Col>{group.description}</Col>
          </Row>
          {group.isMember && (
            <>
              <CreateGroupPosts groupId={id} onPostsUpdate={handlePostUpdate} />
              <FeedPosts url={GROUPFEED_URL + id} key={id} reload={reload} />
            </>
          )}
        </Container>
      )}
    </>
  );
};

export default GroupPage;
