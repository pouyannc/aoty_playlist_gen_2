import { useSelector } from "react-redux";
import { Box, Tab, Tabs } from "@mui/material";
import { useNavigate } from "react-router-dom";

const SubPageSwitch = () => {
  const currentYear = new Date().getFullYear();
  const playlistInfo = useSelector(({ playlistOptions }) => playlistOptions);
  const playlistTypeArr = playlistInfo.type.split("/");
  const [time, sort, _] = playlistTypeArr;
  const navigate = useNavigate();

  const handleTabSwitch = (switchVal, typeIdx) => {
    playlistTypeArr[typeIdx] = switchVal;

    navigate("recent/" + playlistTypeArr.join("/"));
  };

  return (
    playlistInfo.category === "recent" && (
      <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
        <Tabs value={time} centered>
          <Tab
            value="months"
            label="Months"
            onClick={() => handleTabSwitch("months", 0)}
          />
          <Tab
            value={currentYear.toString()}
            label={currentYear}
            onClick={() => handleTabSwitch(currentYear.toString(), 0)}
          />
          <Tab
            value="years"
            label="Years"
            onClick={() => handleTabSwitch("years", 0)}
          />
        </Tabs>
        <Tabs value={sort} centered>
          <Tab
            value="must-hear"
            label="Must-Hear"
            onClick={() => handleTabSwitch("must-hear", 1)}
          />
          <Tab
            value="popular"
            label="Popular"
            onClick={() => handleTabSwitch("popular", 1)}
          />
        </Tabs>
      </Box>
    )
  );
};

export default SubPageSwitch;
