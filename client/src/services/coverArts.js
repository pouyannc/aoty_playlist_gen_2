import axios from "axios";

const serverUrl = import.meta.env.VITE_SERVER_URL;

const getCoverArts = async (scrapeKey) => {
  try {
    const res = await axios.get(
      `${serverUrl}/albums/covers?scrape_key=${scrapeKey}`,
      {
        withCredentials: true,
      }
    );
    return res.data;
  } catch (error) {
    console.log("Couldn't get album covers:", error);
    return;
  }
};

export default { getCoverArts };
