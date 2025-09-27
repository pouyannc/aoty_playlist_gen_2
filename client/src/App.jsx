import { useDispatch, useSelector } from "react-redux";
import { Route, Routes } from "react-router-dom";
import { Container } from "@mui/material";
import Nav from "./components/Nav";
import GenPage from "./components/GenPage";
import LoginPage from "./components/LoginPage";
import AboutPage from "./components/AboutPage";
import AuthCallback from "./components/AuthCallback";
import { useEffect } from "react";
import { getAndSetSpotifyUID } from "./reducers/userReducer";

function App() {
  const uid = useSelector(({ user }) => user.spotifyUID);
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(getAndSetSpotifyUID());
  }, [uid]);

  return (
    <Routes>
      <Route path="/auth/callback" element={<AuthCallback />} />
      <Route
        path="/*"
        element={
          uid ? (
            <Container
              sx={{
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                textAlign: "center",
              }}
            >
              <Nav />
              <Routes>
                <Route path="/*" element={<GenPage />} />
                <Route path="/about" element={<AboutPage />} />
              </Routes>
            </Container>
          ) : (
            <LoginPage />
          )
        }
      />
    </Routes>
  );
}

export default App;
