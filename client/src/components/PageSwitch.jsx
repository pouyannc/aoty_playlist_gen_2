import { useSelector } from "react-redux";
import { Box, Tab, Tabs } from "@mui/material";
import { Link, useNavigate } from "react-router-dom";

const PageSwitch = () => {
  const playlistOptions = useSelector(({ playlistOptions }) => playlistOptions);

  const navigate = useNavigate();

  return (
    <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
      <Tabs value={playlistOptions.category} centered>
        <Tab
          value="new"
          label="New"
          component={Link}
          to="/"
          onClick={() => navigate("/")}
        />
        <Tab
          value="recent"
          label="Recent"
          onClick={() => navigate("/recent/months/must-hear/all")}
        />
      </Tabs>
    </Box>
  );
};

export default PageSwitch;
