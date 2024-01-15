import React, { useState, useEffect } from "react";
import Select from "react-select";
import axios from "axios";
import { ADD_GROUP_MEMBERS_URL } from "../utils/routes";
import { Button, Stack, Alert } from "react-bootstrap";

const AddGroupMembers = ({ id, handleClose }) => {
  const [errMsg, setErrMsg] = useState("");
  const [followers, setFollowers] = useState([]);
  const [formData, setFormData] = useState([]);

  useEffect(() => {
    const fetchFollowers = async () => {
      try {
        const response = await axios.get(ADD_GROUP_MEMBERS_URL + `/${id}`, {
          withCredentials: true,
        });
        setFollowers(response.data);
      } catch (err) {
        if (!err?.response) {
          setErrMsg("No Server Response");
        } else {
          setErrMsg("Internal Server Error");
        }
      }
    };
    fetchFollowers();
    // eslint-disable-next-line
  }, []);

  const handleSelectChange = (selectedOptions) => {
    const selectedValues = selectedOptions.map((option) => option.value);
    setFormData(selectedValues);
  };

  const userOptions = followers.map((follower) => ({
    value: follower.id,
    label: `${follower.firstName} ${follower.lastName}`,
  }));

  const handleSubmit = async () => {
    try {
      await axios.post(
        ADD_GROUP_MEMBERS_URL,
        JSON.stringify({ groupId: +id, userIds: formData }),
        {
          withCredentials: true,
        }
      );
      setFormData([]);
      handleClose();
    } catch (err) {
      if (!err?.response) {
        setErrMsg("No Server Response");
      } else {
        setErrMsg("Internal Server Error");
      }
    }
  };

  return (
    <>
      {errMsg && (
        <Alert variant="danger" className="text-center">
          {errMsg}
        </Alert>
      )}
      <Stack direction="horizontal">
        <div className="add-members">
          <Select options={userOptions} isMulti onChange={handleSelectChange} />
        </div>
        <Button onClick={handleSubmit}>Invite</Button>
      </Stack>
    </>
  );
};

export default AddGroupMembers;
