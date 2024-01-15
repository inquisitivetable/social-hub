import React, { useEffect, useState } from "react";
import { Navigate, useLocation } from "react-router-dom";
import Login from "../pages/LoginPage";
import Signup from "../pages/SignupPage";
import useAuth from "../hooks/useAuth";
import axios from "axios";
import { AUTH_URL } from "../utils/routes";

const RequireGuest = () => {
  const { auth, setAuth } = useAuth();
  const location = useLocation();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const authorisation = async () => {
      try {
        await axios.get(AUTH_URL, {
          withCredentials: true,
        });
        setAuth(true);
      } catch (err) {
        if (!err?.response) {
          setAuth(false);
        } else if (err.response?.status === 401) {
          setAuth(false);
        } else {
          setAuth(false);
        }
      }
      setLoading(false);
    };

    authorisation();
    // eslint-disable-next-line
  }, [location]);

  const isSignup = location?.pathname === "/signup";

  return loading ? null : auth ? (
    <Navigate to="/profile" replace="true" />
  ) : isSignup ? (
    <Signup />
  ) : (
    <Login />
  );
};

export default RequireGuest;
