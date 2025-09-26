import { useEffect } from "react";
import { useDispatch } from "react-redux";
import { getAndSetSpotifyUID } from "../reducers/userReducer";

const AuthCallback = () => {
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(getAndSetSpotifyUID());
  }, []);

  return <p>Logging you in...</p>;
};

export default AuthCallback;
