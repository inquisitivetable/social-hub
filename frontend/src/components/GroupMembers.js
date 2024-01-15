import { useState, useEffect } from "react";
import { makeRequest } from "../services/makeRequest";
import GenericModal from "../components/GenericModal";
import GenericUserList from "../components/GenericUserList";
import { GROUP_MEMBERS_URL } from "../utils/routes";
import { ListGroup, Col, Alert } from "react-bootstrap";
import AddGroupMembers from "./AddGroupMembers";
import { PlusCircle } from "react-bootstrap-icons";

const GroupMembers = ({ groupId }) => {
  const [groupMembers, setGroupMembers] = useState([]);
  const [errMsg, setErrMsg] = useState(null);

  useEffect(() => {
    const loadMembers = async () => {
      try {
        const response = await makeRequest(`/groupmembers/${groupId}`);
        if (response !== null) {
          setGroupMembers(response);
        }
      } catch (err) {
        setErrMsg(err);
      }
    };
    loadMembers();
  }, [groupId]);

  const inviteMembers = (
    <GenericModal img={<PlusCircle />} variant="flush" headerText="Add members">
      <AddGroupMembers id={groupId} />
    </GenericModal>
  );

  return (
    <>
      {errMsg ? (
        <Alert variant="danger" className="text-center">
          {errMsg}
        </Alert>
      ) : (
        <Col xs="auto">
          <GenericModal
            buttonText={`${groupMembers?.length} members`}
            headerText={"Group members"}
            headerButton={inviteMembers}
          >
            <ListGroup>
              <GenericUserList
                variant="flush"
                url={GROUP_MEMBERS_URL + groupId}
              />
            </ListGroup>
          </GenericModal>
        </Col>
      )}
    </>
  );
};

export default GroupMembers;
