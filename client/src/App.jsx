import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { getUID } from "./reducers/userReducer";
import { setTokens } from "./services/user";
import { Route, Routes } from "react-router-dom";
import { Container } from "@mui/material";
import Nav from "./components/Nav";
import GenPage from "./components/GenPage";
import LoginPage from "./components/LoginPage";
import AboutPage from "./components/AboutPage";
import AuthCallback from "./components/AuthCallback";

function App() {
  const uid = useSelector(({ user }) => user.spotifyUID);
  const dispatch = useDispatch();

  useEffect(() => {
    let accessToken = localStorage.getItem("access");
    let refreshToken = localStorage.getItem("refresh");

    if (accessToken && refreshToken) {
      setTokens(accessToken, refreshToken);
    }

    accessToken && dispatch(getUID());
  }, []);

  return (
    <Routes>
      <Route path="/auth/callback" element={<AuthCallback />} />

      <Route
        path="/*"
        element={
          uid === "" ? (
            <>
              <LoginPage />
            </>
          ) : (
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
                <Route path="/" element={<GenPage />} />
                <Route path="/about" element={<AboutPage />} />
              </Routes>
            </Container>
          )
        }
      />
    </Routes>
  );
}

export default App;
