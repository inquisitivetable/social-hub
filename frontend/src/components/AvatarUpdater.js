import AvatarEditor from "react-avatar-editor";
import { useState, useRef } from "react";
import FileUploader from "./FileUploader";
import axios from "axios";
import { Button, Stack, Alert, InputGroup } from "react-bootstrap";

const AvatarUpdater = ({ url, handleClose, onUploadSuccess }) => {
  const editorRef = useRef();
  const [selectedImage, setSelectedImage] = useState(null);
  const [errMsg, setErrMsg] = useState("");

  const handleClick = async () => {
    if (!selectedImage) {
      setErrMsg("You have to select an image first");
      return;
    }

    const canvas = editorRef.current.getImage();

    canvas.toBlob(async (blob) => {
      const formData = new FormData();
      formData.append("image", blob, "avatar.jpg");

      try {
        await axios.post(url, formData, { withCredentials: true });
        handleClose();
        onUploadSuccess();
      } catch (err) {
        if (!err?.response) {
          setErrMsg("No Server Response");
        } else if (err.response?.status > 200) {
          setErrMsg("Internal Server Error");
        }
      }
    }, "image/jpeg");
  };

  return (
    <Stack gap={2}>
      <AvatarEditor
        ref={editorRef}
        image={
          selectedImage
            ? selectedImage
            : `${process.env.PUBLIC_URL}/defaultuser.jpg`
        }
        className="mx-auto"
        width={250}
        height={250}
        color={[255, 255, 255, 0.6]}
        scale={1.2}
        rotate={0}
      />
      <InputGroup>
        <FileUploader
          onFileSelectSuccess={(file) => {
            setSelectedImage(file);
            setErrMsg("");
          }}
          onFileSelectError={({ error }) => setErrMsg(error)}
        />
        <Button onClick={handleClick} disabled={!selectedImage || errMsg}>
          Save image
        </Button>
      </InputGroup>

      {errMsg && (
        <Alert variant="danger" className="text-center">
          {errMsg}
        </Alert>
      )}
    </Stack>
  );
};

export default AvatarUpdater;
