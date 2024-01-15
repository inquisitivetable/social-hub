import List from "../components/List.js";
import ImageHandler from "../utils/imageHandler.js";
import { ListGroup } from "react-bootstrap";
import { LinkContainer } from "react-router-bootstrap";

const image = (user) =>
  ImageHandler(user?.imagePath, "defaultuser.jpg", "userlist-img");

const GenericUserList = ({ url }) => {
  const mapUsers = (user, index) => {
    return (
      <LinkContainer key={index} to={`/profile/${user.id}`}>
        <ListGroup.Item action>
          <>
            {image(user)}
            {user?.nickname
              ? `${user.nickname}`
              : `${user.firstName} ${user.lastName}`}
          </>
        </ListGroup.Item>
      </LinkContainer>
    );
  };

  return <List url={url} mapFunction={mapUsers} />;
};

export default GenericUserList;
