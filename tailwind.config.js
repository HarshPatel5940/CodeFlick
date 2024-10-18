import { blue, yellow, gray } from "tailwindcss/colors";

/** @type {import('tailwindcss').Config} */
export const content = ["internal/pages/*.templ"];
export const theme = {
  container: {
    center: true,
    padding: {
      DEFAULT: "1rem",
      mobile: "2rem",
      tablet: "4rem",
      desktop: "5rem",
    },
  },
  extend: {
    colors: {
      primary: blue,
      secondary: yellow,
      neutral: gray,
    },
  },
};
export const plugins = [
  require("@tailwindcss/forms"),
  require("@tailwindcss/typography"),
];
