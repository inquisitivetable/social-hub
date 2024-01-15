import React, { useState } from "react";
import SearchBar from "../components/SearchBar";
import SearchResults from "../components/SearchResults";

const SearchUtility = () => {
  const [searchResults, setSearchResults] = useState([]);

  return (
    <>
      <SearchBar setSearchResults={setSearchResults} />
      {searchResults.length > 0 && (
        <SearchResults
          searchResults={searchResults}
          setSearchResults={setSearchResults}
        />
      )}
    </>
  );
};

export default SearchUtility;
