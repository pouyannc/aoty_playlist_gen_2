import { createSlice } from "@reduxjs/toolkit";
import { getSpotifyUID } from "../services/user";

const initialState = {
  spotifyUID: "",
};

const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    setUID(state, action) {
      return { ...state, spotifyUID: action.payload };
    },
    logout() {
      localStorage.removeItem("access");
      localStorage.removeItem("refresh");
      localStorage.removeItem("expiresAt");
      return initialState;
    },
  },
});

export const { setUID, logout } = userSlice.actions;

export const getAndSetSpotifyUID = () => async (dispatch) => {
  const uid = await getSpotifyUID();
  dispatch(setUID(uid));
};

export default userSlice.reducer;
