import React, { useState, useRef, useEffect, useCallback } from "react";
import { makeRequest } from "../services/makeRequest";
import { Container, Alert } from "react-bootstrap";
import Post from "../components/Post.js";

const FeedPosts = ({ url, reload }) => {
  const observer = useRef();
  const [posts, setPosts] = useState([]);
  const [errMsg, setErrMsg] = useState(null);
  const [offset, setOffset] = useState(0);

  const handlePageChange = (postId) => {
    setOffset(postId);
  };

  useEffect(() => {
    setPosts([]);
    setOffset(0);
  }, [reload]);

  useEffect(() => {
    const abortController = new AbortController();
    const loadPosts = async () => {
      try {
        const response = await makeRequest(`${url}/${offset}`, {
          signal: abortController.signal,
        });
        setPosts((prevPosts) => {
          return [...prevPosts, ...response];
        });
      } catch (error) {
        setErrMsg(error.message);
      }
    };
    loadPosts();

    return () => {
      abortController.abort();
    };
    // eslint-disable-next-line
  }, [offset, reload]);

  const lastPostElementRef = useCallback((node) => {
    if (observer.current) {
      observer.current.disconnect();
    }

    observer.current = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting) {
        const postId = node.getAttribute("data-post-id"); // Get the post ID from the element attribute
        handlePageChange(postId);
      }
    });

    if (node) {
      observer.current.observe(node);
    }
  }, []);

  const renderedPosts = posts?.map((post, index) => {
    const isLastPost = index === posts.length - 1;

    return (
      <Post
        key={index}
        post={post}
        isLastPost={isLastPost}
        lastPostElementRef={lastPostElementRef}
      />
    );
  });

  return (
    <Container fluid>
      {errMsg && <Alert variant="danger">{errMsg}</Alert>}
      {renderedPosts}
    </Container>
  );
};

export default FeedPosts;
