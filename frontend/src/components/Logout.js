import { useEffect } from "react";
import useAuth from "../hooks/useAuth";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { LOGOUT_URL } from "../utils/routes";

const Logout = () => {
  const { setAuth } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    const authorisation = async () => {
      await axios.get(LOGOUT_URL, {
        withCredentials: true,
      });

      setAuth(false);
      navigate("/login", { replace: true });
    };

    authorisation();
    // eslint-disable-next-line
  }, []);
};

export default Logout;
