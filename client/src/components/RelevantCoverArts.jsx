import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { v4 as uuidv4 } from "uuid";
import { getCoverUrls, setRetrievingTrue } from "../reducers/coverArtReducer";
import {
  Box,
  IconButton,
  ImageListItem,
  ImageListItemBar,
  Link,
  Paper,
  Skeleton,
  Snackbar,
} from "@mui/material";
import { FaSpotify } from "react-icons/fa";

const RelevantCoverArts = () => {
  const currentPlaylistType = useSelector(
    ({ playlistOptions }) => playlistOptions.type
  );
  const currentPlaylistScrapeUrl = useSelector(
    ({ playlistOptions }) => playlistOptions.scrapeUrl
  );
  const coverArtUrls = useSelector(({ coverArtUrls }) => coverArtUrls);
  const dispatch = useDispatch();

  const [openSnackbar, setOpenSnackbar] = useState(false);

  useEffect(() => {
    if (!coverArtUrls.retrieving) {
      console.log(coverArtUrls[currentPlaylistType]);
      if (coverArtUrls[currentPlaylistType] === undefined) {
        setTimeout(() => {
          dispatch(setRetrievingTrue());
          setOpenSnackbar(true);
          dispatch(getCoverUrls(currentPlaylistType, currentPlaylistScrapeUrl));
        }, 1000);
      }
    }
  }, [currentPlaylistType, coverArtUrls.retrieving]);

  return (
    <Box
      sx={{
        p: 2,
        display: "grid",
        gridTemplateColumns: { xs: "repeat(2, 1fr)", sm: "repeat(4, 1fr)" },
        justifyItems: "center",
        maxWidth: { xs: 400, sm: 900, lg: 2000 },
      }}
    >
      {/* <Snackbar
        anchorOrigin={{ vertical: "top", horizontal: "center" }}
        open={open}
        autoHideDuration={4000}
        message="Generating playlist preview"
        onClose={(e, reason) => reason !== "clickaway" && setOpen(false)}
      /> */}
      {(!Array.isArray(coverArtUrls[currentPlaylistType])
        ? Array.from(new Array(8))
        : coverArtUrls[currentPlaylistType]
      ).map((album) => (
        <Paper
          key={uuidv4()}
          elevation={10}
          sx={{ m: 0.8, bgcolor: "#0F1A20" }}
        >
          <ImageListItem sx={{ p: 0.6 }}>
            {album ? (
              <Link
                href={`https://open.spotify.com/album/${album.id}`}
                target="_blank"
                rel="noreferrer"
              >
                <img src={album.image_url} style={{ width: "100%" }} />
                <ImageListItemBar
                  subtitle={album.artist}
                  actionIcon={
                    <IconButton
                      size="small"
                      sx={{ color: "rgba(255, 255, 255, 0.54)" }}
                    >
                      <FaSpotify />
                    </IconButton>
                  }
                  sx={{ m: 0.8, height: "14%" }}
                />
              </Link>
            ) : (
              <Box>
                <Skeleton animation="wave" variant="rounded">
                  <img
                    src="https://i.scdn.co/image/ab67616d00001e02e04b8c0b83df4247f25ac979"
                    style={{ width: "100%" }}
                  />
                </Skeleton>
              </Box>
            )}
          </ImageListItem>
        </Paper>
      ))}
    </Box>
  );
};

export default RelevantCoverArts;
