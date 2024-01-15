import { useState } from "react";
import { GROUP_EVENTS_URL } from "../utils/routes";
import CreateEvent from "./CreateEvent";
import GenericEventList from "./GenericEventList";
import { ListGroup } from "react-bootstrap";
import GenericModal from "./GenericModal";
import { PlusCircle } from "react-bootstrap-icons";

const Events = ({ groupId }) => {
  const [reload, setReload] = useState(false);

  const handleEventUpdate = () => {
    setReload(!reload);
  };

  const createEvent = (
    <GenericModal
      img={<PlusCircle />}
      variant="flush"
      headerText="Create an event"
    >
      <CreateEvent onEventCreated={handleEventUpdate} id={groupId} />
    </GenericModal>
  );

  return (
    <>
      <GenericModal buttonText="Events" headerButton={createEvent}>
        <ListGroup>
          <GenericEventList key={reload} url={GROUP_EVENTS_URL + groupId} />
        </ListGroup>
      </GenericModal>
    </>
  );
};

export default Events;
