import List from "./List.js";
import { ListGroup } from "react-bootstrap";
import { LinkContainer } from "react-router-bootstrap";

const GenericEventList = ({ url }) => {
  const mapGenericEventList = (event, index) => (
    <LinkContainer key={index} action active={false} to={`/event/${event.id}`}>
      <ListGroup.Item>
        <>{event.title}</>
      </ListGroup.Item>
    </LinkContainer>
  );

  return <List url={url} mapFunction={mapGenericEventList} />;
};

export default GenericEventList;
