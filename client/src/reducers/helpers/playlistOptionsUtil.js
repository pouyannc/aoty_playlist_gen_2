const d = new Date();
const currentYear = d.getFullYear();
const currentMonthDigit = d.getMonth();
const month = [
  "january",
  "february",
  "march",
  "april",
  "may",
  "june",
  "july",
  "august",
  "september",
  "october",
  "november",
  "december",
];

//redo the genres for better ones
const genres = {
  pop: 15,
  rock: 7,
  hiphop: 3,
  electronic: 6,
  dance: 132,
  metal: 40,
  rb: 22,
  singersongwriter: 37,
  trap: 213,
  indierock: 1,
};

export const generateScrapeURL = (type) => {
  const [time, sort, genre] = type.split("/");

  const baseURL = "https://www.albumoftheyear.org";

  let monthSegment = "";
  if (time === "months") {
    monthSegment = `${month[currentMonthDigit]}-${(currentMonthDigit + 1)
      .toString()
      .padStart(2, "0")}.php`;
  }

  const relativePath = `/${currentYear}/releases/${monthSegment}`;
  const scrapeURL = new URL(relativePath, baseURL);

  scrapeURL.searchParams.append("type", "lp");
  if (sort === "must-hear") {
    scrapeURL.searchParams.append("sort", "user");
    scrapeURL.searchParams.append("reviews", "500");
  }

  if (genre !== "all") {
    const genreParam = genres[genre.toLowerCase()];
    scrapeURL.searchParams.append("genre", genreParam);
    scrapeURL.searchParams.set("review", "100");
  }

  return encodeURIComponent(scrapeURL.href);
};
