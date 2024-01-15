import React, { useState, useEffect } from "react";
import axios from "axios";
import CreateComment from "./CreateComment";
import Comment from "../components/Comment";
import { COMMENTS_URL } from "../utils/routes";
import { Container, Alert, Button } from "react-bootstrap";

const Comments = ({ postId, setCommentNumber }) => {
  const [comments, setComments] = useState([]);
  const [errMsg, setErrMsg] = useState(null);
  const [offset, setOffset] = useState(0);
  const [loading, setLoading] = useState(true);
  const [hasMoreComments, setHasMoreComments] = useState(true);

  const handleCommentsUpdate = () => {
    setOffset(0);
    setLoading(!loading);
  };

  useEffect(() => {
    const loadComments = async () => {
      try {
        await axios
          .get(`${COMMENTS_URL}${postId}/${offset}`, {
            withCredentials: true,
          })
          .then((response) => {
            response.data.length < 5 && setHasMoreComments(false);
            response.data.length > 0 &&
              setCommentNumber(response.data[0].commentCount);
            setComments((prevComments) => {
              const commentIds = new Set(
                prevComments.map((comment) => comment.id)
              );
              const newComments = response.data.filter(
                (comment) => !commentIds.has(comment.id)
              );
              const updatedComments = [...newComments, ...prevComments];
              const sortedComments = updatedComments.sort(
                (a, b) => new Date(b.createdAt) - new Date(a.createdAt)
              );
              return sortedComments;
            });
          });
      } catch (err) {
        if (err.response?.status === 404) {
          setErrMsg(err.message);
        }
      }
    };

    loadComments();
    // eslint-disable-next-line
  }, [offset, loading]);

  const showMoreComments = () => {
    setOffset(offset + 1);
  };

  const renderedComments = comments.map((comment, index) => (
    <Comment comment={comment} key={index} />
  ));

  return (
    <Container>
      {errMsg ? (
        <Alert variant="danger" className="text-center">
          {errMsg}
        </Alert>
      ) : (
        <div className="my-auto mt-2">
          <CreateComment
            postId={postId}
            onCommentsUpdate={handleCommentsUpdate}
          />
          <>{renderedComments}</>
          {hasMoreComments && (
            <div className="d-flex justify-content-center">
              <Button onClick={showMoreComments}>View more comments</Button>
            </div>
          )}
        </div>
      )}
    </Container>
  );
};

export default Comments;
