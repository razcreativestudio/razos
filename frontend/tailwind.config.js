/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        razDark: '#0f172a',
        razCard: 'rgba(15, 23, 42, 0.75)',
        razAccent: '#38bdf8',
        razGold: '#c5a059',
      }
    },
  },
  plugins: [],
}
