import ImageHandler from "../utils/imageHandler";

const SingleChatlistItem = ({ chat }) => {
  const image =
    chat?.user_id > 0
      ? ImageHandler(chat?.avatar_image, "defaultuser.jpg", "chatbox-img")
      : ImageHandler(chat?.avatar_image, "defaultgroup.png", "chatbox-img");

  const listItem = (
    <span className="me-1">
      {image} {chat.name}
    </span>
  );

  return listItem;
};

export default SingleChatlistItem;
