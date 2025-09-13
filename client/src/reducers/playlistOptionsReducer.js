import { createSlice } from "@reduxjs/toolkit";
import { generateScrapeURL } from "./helpers/playlistOptionsUtil";

const initialState = {
  category: "new",
  type: "new",
  title: "New Releases",
  description: "Generate a playlist to sample this weeks most popular releases",
  tracksPerAlbum: 2,
  nrOfTracks: 20,
  scrapeUrl: encodeURIComponent(
    "https://www.albumoftheyear.org/releases/this-week/"
  ),
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
      const [category, type] = [splitPath[1], splitPath.slice(2).join("/")];

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
        tracksPerAlbum,
        nrOfTracks,
        scrapeUrl: generateScrapeURL(type),
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
