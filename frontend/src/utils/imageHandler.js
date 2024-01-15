import Image from "react-bootstrap/Image";

const ImageSource = (path, defaultImage) =>
  path
    ? `${process.env.PUBLIC_URL}/images/${path}`
    : `${process.env.PUBLIC_URL}/${defaultImage}`;

const ImageHandler = (path, defaultImage, className) => {
  return (
    <Image
      fluid
      thumbnail
      className={`${className} p-0 m-0`}
      src={ImageSource(path, defaultImage)}
    />
  );
};

export default ImageHandler;
