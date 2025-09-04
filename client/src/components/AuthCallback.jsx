import axios from "axios";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { setTokens } from "../services/user";
import saveSessionExpiry from "../util/saveSessionExpiry";
import { useDispatch } from "react-redux";
import { getUID } from "../reducers/userReducer";

const AuthCallback = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();

  useEffect(() => {
    const getTokens = async () => {
      try {
        const res = await axios.get(
          `${import.meta.env.VITE_SERVER_URL}/auth/tokens`,
          {
            withCredentials: true,
          },
        );

        setTokens(res.data.access_token, res.data.refresh_token);
        saveSessionExpiry(res.data.expires_in);
        dispatch(getUID());

        navigate("/");
      } catch (err) {
        console.error("error during auth:", err);
      }
    };

    getTokens();
  }, [navigate]);

  return <p>Logging you in...</p>;
};

export default AuthCallback;
