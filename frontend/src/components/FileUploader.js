import React from "react";
import { Form } from "react-bootstrap";

const FileUploader = ({ onFileSelectSuccess, onFileSelectError }) => {
  const handleFileInput = (event) => {
    const allowedTypes = ["image/gif", "image/jpeg", "image/png"];
    const file = event.target.files[0];

    if (!allowedTypes.includes(file?.type)) {
      onFileSelectError({
        error: "You can only upload PNG, JPEG or GIF image files",
      });
    } else if (file?.size > 5242880) {
      onFileSelectError({ error: "File size cannot exceed 5MB" });
    } else onFileSelectSuccess(file);
  };

  return (
    <Form.Group controlId="formFile">
      <Form.Control
        type="file"
        accept="image/gif,image/jpeg,image/png"
        onChange={handleFileInput}
      />
    </Form.Group>
  );
};

export default FileUploader;
