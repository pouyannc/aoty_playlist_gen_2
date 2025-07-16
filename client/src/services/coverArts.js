import axios from "axios";
import refreshSessionIfNeeded from "../util/checkAndRefreshSession";

const serverUrl = import.meta.env.VITE_SERVER_URL;

const getCoverArts = async (q) => {
  await refreshSessionIfNeeded();
  const { scrapeUrl, type } = q;

  const res = await axios.get(
    `${serverUrl}/albums/covers?scrape_url=${scrapeUrl}&type=${
      type.split("/")[0]
    }`,
    {
      withCredentials: true,
    }
  );
  return res.data;
};

export default { getCoverArts };
