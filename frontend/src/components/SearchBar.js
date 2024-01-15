import { useState } from "react";
import { makeRequest } from "../services/makeRequest.js";
import { Form, Alert } from "react-bootstrap";

const SearchBar = ({ setSearchResults }) => {
  const [errMsg, setErrMsg] = useState(null);

  const fetchData = async (value) => {
    try {
      const response = await makeRequest(`search/${value}`, {});
      setSearchResults(response);
    } catch (err) {
      if (!err?.response) {
        setErrMsg("No Server Response");
      } else {
        setErrMsg("Internal Server Error");
      }
    }
  };

  const handleChange = (event) => {
    const value = event.target.value;
    if (!value) {
      setSearchResults([]);
    } else {
      fetchData(value);
    }
  };

  return (
    <>
      {errMsg && (
        <Alert variant="danger" className="text-center">
          {errMsg}
        </Alert>
      )}
      <Form>
        <Form.Control
          placeholder="Search for users and groups"
          onChange={handleChange}
        />
      </Form>
    </>
  );
};

export default SearchBar;
