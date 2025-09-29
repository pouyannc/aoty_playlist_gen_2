import { createSlice } from "@reduxjs/toolkit";
import { getSpotifyUID, logout } from "../services/user";

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
    clearState() {
      return initialState;
    },
  },
});

export const { setUID, clearState } = userSlice.actions;

export const getAndSetSpotifyUID = () => async (dispatch) => {
  const uid = await getSpotifyUID();
  dispatch(setUID(uid));
};

export const logoutAndClearUserState = () => async (dispatch) => {
  await logout();
  dispatch(clearState());
};

export default userSlice.reducer;
