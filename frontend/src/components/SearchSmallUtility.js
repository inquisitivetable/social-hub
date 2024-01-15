import { useState } from "react";
import SearchBar from "./SearchBar";
import SearchResults from "./SearchResults";
import { Offcanvas, Nav } from "react-bootstrap";

const SearchSmallUtility = () => {
  const [show, setShow] = useState(false);
  const [searchResults, setSearchResults] = useState([]);

  const handleShow = () => setShow(true);
  const handleClose = () => setShow(false);

  return (
    <>
      <Nav.Link onClick={handleShow}>Search</Nav.Link>

      {show && (
        <Offcanvas show={show} onHide={handleClose} responsive="md">
          <Offcanvas.Header className="ms-auto" closeButton />
          <Offcanvas.Body>
            <SearchBar setSearchResults={setSearchResults} />
            {searchResults.length > 0 && (
              <SearchResults
                searchResults={searchResults}
                setSearchResults={setSearchResults}
                handleClose={handleClose}
              />
            )}
          </Offcanvas.Body>
        </Offcanvas>
      )}
    </>
  );
};

export default SearchSmallUtility;
