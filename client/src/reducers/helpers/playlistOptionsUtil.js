const d = new Date();
const currentYear = d.getFullYear();

const genreCodes = {
  all: 0,
  pop: 15,
  rock: 7,
  hiphop: 3,
  electronic: 6,
  dance: 132,
  metal: 40,
  rb: 22,
  singersongwriter: 37,
};

const currentYearKeySegment = "year";

export const generateScrapeKey = (type) => {
  let [time, sort, genre] = type.split("/");

  if (time === currentYear.toString()) time = currentYearKeySegment;
  genre = genreCodes[genre].toString();

  return [time, sort, genre].join("/");
};

export const tabTitles = {
  new: {
    title: "New Releases",
    description:
      "Generate a playlist to sample this weeks most popular releases",
  },
  months: {
    title: "Recent Months",
    description:
      "Generate a playlist to sample the hottest records of the last four months",
  },
  2025: {
    title: "2025",
    description: "Generate a playlist to sample the hottest records of 2025",
  },
  years: {
    title: "Recent Years",
    description:
      "Generate a playlist to sample the hottest records of the last three years",
  },
};
