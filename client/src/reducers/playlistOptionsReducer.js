import { createSlice } from "@reduxjs/toolkit";
import { generateScrapeKey, tabTitles } from "./helpers/playlistOptionsUtil";

const initialState = {
  category: "new",
  type: "new",
  title: tabTitles["new"].title,
  description: tabTitles["new"].description,
  tracksPerAlbum: 2,
  nrOfTracks: 20,
  scrapeKey: "new",
};

const playlistOptionsSlice = createSlice({
  name: "playlistOptions",
  initialState,
  reducers: {
    initNew() {
      return initialState;
    },
    setPlaylistOptions(state, action) {
      const splitPath = action.payload.split("/");
      const [category, type, tab] = [
        splitPath[1],
        splitPath.slice(2).join("/"),
        splitPath[2],
      ];

      const { title, description } = tabTitles[tab];

      let [tracksPerAlbum, nrOfTracks] = [
        state.tracksPerAlbum,
        state.nrOfTracks,
      ];
      if (state.category !== category) {
        tracksPerAlbum = 1;
        nrOfTracks = 30;
      }

      return {
        ...state,
        category,
        type,
        title,
        description,
        tracksPerAlbum,
        nrOfTracks,
        scrapeKey: generateScrapeKey(type),
      };
    },
    setTracksPerAlbum(state, action) {
      return { ...state, tracksPerAlbum: parseInt(action.payload) || 0 };
    },
    setNrOfTracks(state, action) {
      return { ...state, nrOfTracks: parseInt(action.payload) || 0 };
    },
  },
});

export const { initNew, setPlaylistOptions, setTracksPerAlbum, setNrOfTracks } =
  playlistOptionsSlice.actions;

export default playlistOptionsSlice.reducer;
