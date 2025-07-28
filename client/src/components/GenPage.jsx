import PageSwitch from "./PageSwitch";
import SubPageSwitch from "./SubPageSwitch";
import PageContent from "./PageContent";

const GenPage = () => {
  return (
    <div>
      <PageSwitch />
      <SubPageSwitch />
      <PageContent />
    </div>
  );
};

export default GenPage;
