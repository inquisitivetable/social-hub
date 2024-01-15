import React, { useState } from "react";
import axios from "axios";
import { Form, Stack, Col, InputGroup, Alert } from "react-bootstrap";
import { GROUP_PAGE_URL } from "../utils/routes";
import PostButton from "./PostButton";

const CreateGroupPosts = ({ groupId, onPostsUpdate }) => {
  const initialFormData = {
    content: "",
    imagePath: "",
    privacyType: 0,
    selectedReceivers: [],
  };

  const [formData, setFormData] = useState(initialFormData);
  const [errMsg, setErrMsg] = useState("");

  const handleChange = (event) => {
    setErrMsg("");
    const { name, value } = event.target;

    setFormData((prevFormData) => ({
      ...prevFormData,
      [name]: value,
    }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    try {
      const response = await axios.post(
        `${GROUP_PAGE_URL}${groupId}/post`,
        JSON.stringify(formData),
        { withCredentials: true },
        {
          headers: { "Content-Type": "application/json" },
        }
      );

      setErrMsg(response.data?.message);
      onPostsUpdate();
    } catch (err) {
      if (!err?.response) {
        setErrMsg("No Server Response");
      } else {
        setErrMsg("Internal Server Error");
      }
    }

    setFormData(initialFormData);
  };

  return (
    <>
      {errMsg && (
        <Alert variant="danger" className="text-center">
          {errMsg}
        </Alert>
      )}
      <Form onSubmit={handleSubmit}>
        <Col>
          <Stack direction="horizontal" gap="2">
            <InputGroup>
              <Form.Control
                as="textarea"
                placeholder="Write what's on your mind"
                onChange={handleChange}
                value={formData.content}
                name="content"
              />
            </InputGroup>
            <Col as={PostButton} className="text-center" />
          </Stack>
        </Col>
      </Form>
    </>
  );
};

export default CreateGroupPosts;
