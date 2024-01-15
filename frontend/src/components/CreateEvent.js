import React, { useState } from "react";
import axios from "axios";
import { Form, Button, Alert, FloatingLabel } from "react-bootstrap";
import { useForm } from "react-hook-form";
import { CREATE_GROUP_EVENT_URL } from "../utils/routes";

const CreateEvent = ({ onEventCreated, id, handleClose }) => {
  const [errMsg, setErrMsg] = useState("");
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm({
    mode: "onBlur",
    criteriaMode: "all",
  });

  const onSubmit = async (data) => {
    try {
      await axios.post(
        CREATE_GROUP_EVENT_URL,
        JSON.stringify({
          ...data,
          startTime: new Date(data.startTime).toISOString(),
          endTime: new Date(data.endTime).toISOString(),
          group_id: +id,
        }),
        { withCredentials: true },
        {
          headers: { "Content-Type": "application/json" },
        }
      );
      onEventCreated();
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
        <FloatingLabel className="mb-3" controlId="floatingTitle" label="Name">
          <Form.Control
            placeholder="Event name"
            autoFocus
            {...register("title", {
              required: "Please enter a name for the event",
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
              required: "Please enter a description for the event",
            })}
          />
          {errors.description && (
            <Alert variant="danger">{errors.description.message}</Alert>
          )}
        </FloatingLabel>
        <FloatingLabel
          className="mb-3"
          controlId="floatingStartTime"
          label="Start time"
        >
          <Form.Control
            type="datetime-local"
            placeholder="Start time"
            {...register("startTime", {
              required: "Please choose a start time",
              validate: (value) =>
                new Date(value) > new Date() ||
                "Event's start cannot be in the past",
            })}
          />
          {errors.startTime && (
            <Alert variant="danger">{errors.startTime.message}</Alert>
          )}
        </FloatingLabel>
        <FloatingLabel
          className="mb-3"
          controlId="floatingEndTime"
          label="End time"
        >
          <Form.Control
            type="datetime-local"
            placeholder="End time"
            {...register("endTime", {
              required: "Please choose a start time",
              validate: (value) =>
                new Date(value) > new Date(watch("startTime")) ||
                "Event cannot end before it starts",
            })}
          />
          {errors.endTime && (
            <Alert variant="danger">{errors.endTime.message}</Alert>
          )}
        </FloatingLabel>
        <Button type="submit">Create</Button>
      </Form>
    </>
  );
};

export default CreateEvent;
