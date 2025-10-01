import axios from "axios";

const serverURL = import.meta.env.VITE_SERVER_URL;

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

const loginGuest = async () => {
  try {
    await axios.get(`${serverURL}/login/guest`, {
      withCredentials: true,
    });
  } catch (error) {
    console.log(error);
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

export { getSpotifyUID, loginGuest, logout };
