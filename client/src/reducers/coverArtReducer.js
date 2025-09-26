import { createSlice } from "@reduxjs/toolkit";
import coverArtsService from "../services/coverArts";

const initialState = {};

const coverArtSlice = createSlice({
  name: "coverArt",
  initialState,
  reducers: {
    setCoverUrls(state, action) {
      const { type, coverUrls } = action.payload;
      return { ...state, [type]: coverUrls };
    },
  },
});

export const { setCoverUrls } = coverArtSlice.actions;

export const getCoverUrls =
  (playlistType, playlistScrapeUrl) => async (dispatch) => {
    dispatch(setCoverUrls({ type: playlistType, coverUrls: true }));
    const coverUrls = await coverArtsService.getCoverArts(
      playlistType,
      playlistScrapeUrl
    );
    dispatch(setCoverUrls({ type: playlistType, coverUrls }));
  };

export default coverArtSlice.reducer;
