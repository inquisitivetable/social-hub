import React, { useState, useEffect } from "react";
import axios from "axios";
import Select from "react-select";
import ImageUploadModal from "./ImageUploadModal";
import {
  Form,
  Image,
  InputGroup,
  Alert,
  Button,
  Col,
  Stack,
} from "react-bootstrap";
import { FOLLOWERS_URL, CREATE_POST_URL } from "../utils/routes";
import GenericModal from "../components/GenericModal";
import { ImageFill } from "react-bootstrap-icons";

const CreatePost = ({ onPostsUpdate, handleClose }) => {
  const [followers, setFollowers] = useState([]);
  const [errMsg, setErrMsg] = useState("");
  const initialFormData = {
    content: "",
    image: null,
    privacyType: 1,
    selectedReceivers: [],
  };
  const [formData, setFormData] = useState(initialFormData);

  useEffect(() => {
    const fetchFollowers = async () => {
      try {
        const response = await axios.get(FOLLOWERS_URL, {
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
    if (formData.privacyType === 3) {
      fetchFollowers();
    }
  }, [formData.privacyType]);

  const handleImageUpload = (image) => {
    setFormData((prevFormData) => ({
      ...prevFormData,
      image: image,
    }));
  };

  const handleChange = (event) => {
    const { name, value, type } = event.target;

    setErrMsg();
    setFormData((prevFormData) => ({
      ...prevFormData,
      [name]: type === "radio" ? parseInt(value) : value,
    }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (
      formData?.privacyType === 0 ||
      (formData?.content === "" && formData?.image === null)
    ) {
      setErrMsg("Enter your message and select a privacy type for your post");
      return;
    }

    const formDataWithImage = new FormData();
    formDataWithImage.append("content", formData.content);
    formDataWithImage.append("privacyType", formData.privacyType);
    formDataWithImage.append("selectedReceivers", formData.selectedReceivers);

    if (formData?.image) {
      formDataWithImage.append("image", formData?.image);
    }

    try {
      await axios.post(CREATE_POST_URL, formDataWithImage, {
        withCredentials: true,
        headers: { "Content-Type": "multipart/form-data" },
      });
      onPostsUpdate();
      handleClose();
    } catch (err) {
      if (!err?.response) {
        setErrMsg("No Server Response");
      } else {
        setErrMsg("Internal Server Error");
      }
    }

    setFormData(initialFormData);
  };

  const handleSelectChange = (selectedOptions) => {
    const selectedValues = selectedOptions.map((option) =>
      option.value.toString()
    );
    setFormData((prevFormData) => ({
      ...prevFormData,
      selectedReceivers: selectedValues,
    }));
  };

  const followersOptions = followers.map((follower) => ({
    value: follower.id,
    label: `${follower.firstName} ${follower.lastName}`,
  }));

  return (
    <>
      {errMsg && (
        <Alert variant="danger" className="text-center">
          {errMsg}
        </Alert>
      )}

      {formData?.image && (
        <div className="post-img mb-3">
          <Image
            src={URL.createObjectURL(formData?.image)}
            fluid
            alt="Selected"
          />
        </div>
      )}

      <Form onSubmit={handleSubmit}>
        <Stack direction="horizontal">
          <Col>
            <Stack direction="horizontal">
              <InputGroup className="me-2">
                <Form.Control
                  as="textarea"
                  placeholder="Write what's on your mind"
                  onChange={handleChange}
                  value={formData.content}
                  name="content"
                />
              </InputGroup>
              <div>
                <GenericModal
                  variant="flush"
                  img={<ImageFill />}
                  buttonText="Add an image"
                >
                  <ImageUploadModal onUploadSuccess={handleImageUpload} />
                </GenericModal>
              </div>
              <div>
                <Button type="submit">Post</Button>
              </div>
            </Stack>

            <Col className="mb-3">
              <Form.Check
                inline
                label="Public"
                name="privacyType"
                type="radio"
                id="public"
                value={1}
                onChange={handleChange}
              />
              <Form.Check
                inline
                label="Private"
                name="privacyType"
                type="radio"
                id="private"
                value={2}
                onChange={handleChange}
              />
              <Form.Check
                inline
                label="Sub-private"
                name="privacyType"
                type="radio"
                id="subPrivate"
                value={3}
                onChange={handleChange}
              />
            </Col>
            <Col>
              {formData.privacyType === 3 && (
                <Select
                  options={followersOptions}
                  isMulti
                  onChange={handleSelectChange}
                />
              )}
            </Col>
          </Col>
        </Stack>
      </Form>
    </>
  );
};

export default CreatePost;
