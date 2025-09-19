import axios from "axios";
import saveSessionExpiry from "../util/saveSessionExpiry";
import { setUID } from "../reducers/userReducer";

const serverURL = import.meta.env.VITE_SERVER_URL;

let tokens;

const setTokens = (access, refresh) => {
  tokens = {
    accessToken: `Bearer ${access}`,
    refreshToken: refresh,
  };
  localStorage.setItem("access", access);
  localStorage.setItem("refresh", refresh);
};

const refreshToken = async () => {
  const res = await axios.get(`${serverURL}/login/refresh}`, {
    withCredentials: true,
  });
  setTokens(res.data.access_token, tokens.refreshToken);
  saveSessionExpiry(res.data.expires_in);
};

const getSpotifyUID = async () => {
  try {
    const res = await axios.get(`${serverURL}/auth/tokens`, {
      withCredentials: true,
    });
    return res.data.spotify_uid;
  } catch (error) {
    console.log(error);
    return "";
  }
};

const logout = async () => {
  try {
    await axios.delete(`${serverURL}/logout`, {
      withCredentials: true,
    });
  } catch (error) {
    console.log(error);
  }
};

export { setTokens, refreshToken, getSpotifyUID, logout };
