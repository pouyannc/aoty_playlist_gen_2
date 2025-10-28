import axios from "axios";

const serverUrl = import.meta.env.VITE_SERVER_URL;

const getTracklist = async (q) => {
  const { scrapeKey, tracksPerAlbum, nrOfTracks, uid, playlistName } = q;
  const res = await axios.post(
    `${serverUrl}/albums/playlist?scrape_key=${scrapeKey}&nr_tracks=${nrOfTracks}&tracks_per=${tracksPerAlbum}`,
    { uid, playlistName },
    {
      withCredentials: true,
    }
  );
  console.log(res.data);
  return res.data;
};

export default { getTracklist };
