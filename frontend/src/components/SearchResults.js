import React from "react";
import { LinkContainer } from "react-router-bootstrap";
import { ListGroup, ListGroupItem } from "react-bootstrap";
import { Scrollbars } from "react-custom-scrollbars-2";

const SearchResults = ({ searchResults, setSearchResults, handleClose }) => {
  const searchResultsMap = searchResults?.map((result, index) => (
    <LinkContainer
      to={
        result.userId === 0
          ? `/groups/${result.groupId}`
          : `/profile/${result.userId}`
      }
      key={index}
      onClick={() => handleClose && handleClose()}
    >
      <ListGroupItem
        action
        onClick={() => {
          setSearchResults([]);
        }}
      >
        <>
          <div>{result.name}</div>
          {result.userId === 0 ? (
            <small className="text-muted">
              <i>group</i>
            </small>
          ) : (
            <small className="text-muted">
              <i>user</i>
            </small>
          )}
        </>
      </ListGroupItem>
    </LinkContainer>
  ));

  return (
    <ListGroup className="search-results position-fixed">
      <Scrollbars autoHide autoHeight>
        {searchResultsMap}
      </Scrollbars>
    </ListGroup>
  );
};

export default SearchResults;
