/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        "empower-blue": "#1e3a8a",
        "empower-dark": "#1e40af",
        "chart-blue": "#3b82f6",
        "success-green": "#10b981",
        "warning-orange": "#f59e0b",
      },
      fontFamily: {
        sans: ["Inter", "system-ui", "sans-serif"],
      },
    },
  },
  plugins: [],
};
