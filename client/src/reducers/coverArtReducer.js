import { createSlice } from "@reduxjs/toolkit";
import coverArtsService from "../services/coverArts";

const initialState = { retrieving: false };

const coverArtSlice = createSlice({
  name: "coverArt",
  initialState,
  reducers: {
    setCoverUrls(state, action) {
      const { type, coverUrls } = action.payload;
      return { ...state, [type]: coverUrls };
    },
    setRetrievingTrue(state) {
      return { ...state, retrieving: true };
    },
    setRetrievingFalse(state) {
      return { ...state, retrieving: false };
    },
  },
});

export const { setCoverUrls, setRetrievingFalse, setRetrievingTrue } =
  coverArtSlice.actions;

export const getCoverUrls =
  (playlistType, playlistScrapeUrl) => async (dispatch) => {
    dispatch(setCoverUrls({ type: playlistType, coverUrls: true }));
    const coverUrls = await coverArtsService.getCoverArts(
      playlistType,
      playlistScrapeUrl
    );
    dispatch(setCoverUrls({ type: playlistType, coverUrls }));
    dispatch(setRetrievingFalse());
  };

export default coverArtSlice.reducer;
