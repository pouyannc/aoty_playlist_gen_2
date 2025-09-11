import { Box, Button, Container, Typography } from "@mui/material";
import { BiSolidSpeaker } from "react-icons/bi";
import { TbVinyl } from "react-icons/tb";
import { refreshToken, setTokens } from "../services/user";
import { useEffect } from "react";
import { useDispatch } from "react-redux";
import { getAndSetSpotifyUID } from "../reducers/userReducer";

const LoginPage = () => {
  const serverUrl = import.meta.env.VITE_SERVER_URL;
  const dispatch = useDispatch();

  const guestLogin = () => {
    console.log("being built...");
  };

  useEffect(() => {
    dispatch(getAndSetSpotifyUID());
  });

  return (
    <div>
      <Container
        sx={{
          display: "flex",
          flexDirection: "column",
          gap: 3,
          textAlign: "center",
          my: "12%",
        }}
      >
        <Typography variant="h2">AOTY Playlist Gen</Typography>
        <Box
          sx={{
            display: "flex",
            justifyContent: "space-evenly",
            alignItems: "center",
          }}
        >
          <BiSolidSpeaker size={80} />
          <TbVinyl size={70} />
          <BiSolidSpeaker size={80} />
        </Box>
        <Typography variant="h6">
          Discover music by generating playlists compiled from the hottest
          albums as per the{" "}
          <a
            style={{ textDecoration: "none", color: "inherit" }}
            href="https://www.albumoftheyear.org/"
          >
            aoty.org
          </a>{" "}
          community.
        </Typography>
        <Typography variant="h6">
          The current version of this app can be used without linking a Spotify
          account. Logging in with Spotify will be reenabled once the app enters
          extended quota.
        </Typography>
        <Button
          href={`${serverUrl}/login`}
          variant="contained"
          sx={{
            bgcolor: "green",
            fontWeight: 700,
            width: "40%",
            alignSelf: "center",
          }}
        >
          Login with Spotify
        </Button>
        <Button
          variant="contained"
          sx={{ fontWeight: 700, width: "40%", alignSelf: "center" }}
          onClick={guestLogin}
        >
          Enter without login
        </Button>
      </Container>
    </div>
  );
};

export default LoginPage;
