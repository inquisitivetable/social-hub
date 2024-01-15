import React, { useState } from "react";
import GenericGroupList from "../components/GenericGroupList";
import GenericEventList from "../components/GenericEventList";
import {
  USER_CREATED_GROUPS_URL,
  USER_GROUPS_URL,
  ACCEPTED_EVENTS_URL,
} from "../utils/routes";
import CreateGroup from "../components/CreateGroup";
import { Container, ListGroup, Row, Stack } from "react-bootstrap";
import { Scrollbars } from "react-custom-scrollbars-2";
import { PlusCircle } from "react-bootstrap-icons";
import GenericModal from "../components/GenericModal";

const GroupSidebar = () => {
  const [loadNewGroups, setLoadNewGroups] = useState(0);

  const handleGroupUpdate = () => {
    setLoadNewGroups((prevCount) => prevCount + 1);
  };

  return (
    <Scrollbars autoHide>
      <Container>
        <h4>Groups</h4>
        <ListGroup variant="flush">
          <GenericGroupList url={USER_GROUPS_URL} />
        </ListGroup>
        <Row>
          <Stack direction="horizontal">
            <h4>My groups</h4>
            <div>
              <GenericModal
                img={<PlusCircle />}
                variant="flush"
                headerText="Create a group"
              >
                <CreateGroup onGroupCreated={handleGroupUpdate} />
              </GenericModal>
            </div>
          </Stack>
        </Row>

        <ListGroup variant="flush">
          <GenericGroupList
            url={USER_CREATED_GROUPS_URL}
            loadNewGroups={loadNewGroups}
          />
        </ListGroup>

        <h4>Events</h4>
        <ListGroup variant="flush">
          <GenericEventList url={ACCEPTED_EVENTS_URL} />
        </ListGroup>
      </Container>
    </Scrollbars>
  );
};

export default GroupSidebar;
