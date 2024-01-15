import List from "./List.js";
import { ListGroup } from "react-bootstrap";
import { LinkContainer } from "react-router-bootstrap";

const GenericGroupList = ({ url, loadNewGroups }) => {
  const mapGenericGroupList = (group, index) => (
    <LinkContainer
      key={index}
      action
      active={false}
      to={`/groups/${group.groupId}`}
    >
      <ListGroup.Item>
        <>{group.groupName}</>
      </ListGroup.Item>
    </LinkContainer>
  );

  return (
    <List
      url={url}
      mapFunction={mapGenericGroupList}
      loadNewGroups={loadNewGroups}
    />
  );
};

export default GenericGroupList;
