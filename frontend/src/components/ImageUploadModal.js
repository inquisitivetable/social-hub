import React, { useState } from "react";
import "../style.css";
import { Button, Alert, InputGroup, Image, Stack } from "react-bootstrap";
import FileUploader from "../components/FileUploader";

const ImageUploadModal = ({ handleClose, onUploadSuccess }) => {
  const [selectedImage, setSelectedImage] = useState(null);
  const [errMsg, setErrMsg] = useState("");

  const handleUpload = () => {
    onUploadSuccess(selectedImage);
    handleClose();
  };

  return (
    <Stack gap="2">
      {selectedImage && (
        <div className="text-center">
          <Image
            src={URL.createObjectURL(selectedImage)}
            fluid
            alt="Selected"
            className="profile-img"
          />
        </div>
      )}
      <InputGroup>
        <FileUploader
          onFileSelectSuccess={(file) => {
            setSelectedImage(file);
            setErrMsg("");
          }}
          onFileSelectError={({ error }) => setErrMsg(error)}
        />
        <button onClick={handleUpload} disabled={!selectedImage || errMsg}>
          Upload image
        </button>
      </InputGroup>

      {errMsg && (
        <Alert variant="danger" className="text-center">
          {errMsg}
        </Alert>
      )}
    </Stack>
  );
};

export default ImageUploadModal;
