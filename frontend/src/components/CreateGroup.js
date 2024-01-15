import { useState } from "react";
import axios from "axios";
import { Form, Button, FloatingLabel, Alert } from "react-bootstrap";
import { useForm } from "react-hook-form";
import { CREATE_GROUP_URL } from "../utils/routes";

const CreateGroup = ({ onGroupCreated, handleClose }) => {
  const [errMsg, setErrMsg] = useState("");
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    mode: "onBlur",
    criteriaMode: "all",
  });

  const onSubmit = async (data) => {
    try {
      await axios.post(
        CREATE_GROUP_URL,
        JSON.stringify(data),
        { withCredentials: true },
        {
          headers: { "Content-Type": "application/json" },
        }
      );
      onGroupCreated();
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
      <Form onSubmit={handleSubmit(onSubmit)}>
        <FloatingLabel
          className="mb-3"
          controlId="floatingTitle"
          label="Group name"
        >
          <Form.Control
            placeholder="Title"
            autoFocus
            {...register("title", {
              required: "Please enter a name for your group",
            })}
          />
          {errors.title && (
            <Alert variant="danger">{errors.title.message}</Alert>
          )}
        </FloatingLabel>
        <FloatingLabel
          className="mb-3"
          controlId="floatingDescription"
          label="Description"
        >
          <Form.Control
            as="textarea"
            placeholder="Description"
            {...register("description", {
              required: "Please enter a description for your group",
            })}
          />
          {errors.description && (
            <Alert variant="danger">{errors.description.message}</Alert>
          )}
        </FloatingLabel>
        <Button type="submit">Create</Button>
      </Form>
    </>
  );
};

export default CreateGroup;
