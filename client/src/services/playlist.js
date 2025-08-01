import axios from "axios";
import refreshSessionIfNeeded from "../util/checkAndRefreshSession";

const serverUrl = import.meta.env.VITE_SERVER_URL;

const getTracklist = async (q) => {
  await refreshSessionIfNeeded();
  const { scrapeUrl, tracksPerAlbum, nrOfTracks, type, uid, playlistName } = q;
  const res = await axios.post(
    `${serverUrl}/albums/playlist?scrape_url=${scrapeUrl}&nr_tracks=${nrOfTracks}&tracks_per=${tracksPerAlbum}&type=${
      type.split("/")[0]
    }`,
    { uid, playlistName },
    {
      withCredentials: true,
    }
  );
  console.log(res.data);
  return res.data;
};

export default { getTracklist };
