import PageSwitch from "./PageSwitch";
import SubPageSwitch from "./SubPageSwitch";
import PageContent from "./PageContent";
import { useLocation } from "react-router-dom";
import { useDispatch } from "react-redux";
import {
  initNew,
  setPlaylistOptions,
} from "../reducers/playlistOptionsReducer";

const GenPage = () => {
  const dispatch = useDispatch();
  const location = useLocation();

  const currentPathname = location.pathname;
  if (currentPathname === "/") dispatch(initNew());
  else dispatch(setPlaylistOptions(currentPathname));

  return (
    <div>
      <PageSwitch />
      <SubPageSwitch />
      <PageContent />
    </div>
  );
};

export default GenPage;
