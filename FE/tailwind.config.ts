import type { Config } from 'tailwindcss';

const config: Config = {
  content: [
    './src/app/**/*.{ts,tsx}',
    './src/features/**/*.{ts,tsx}',
    './src/shared/**/*.{ts,tsx}',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: '#FF6B35',
        secondary: '#E63946',
        accent: '#FFB703',
        background: '#FFF8F0',
        'dark-bg': '#1A1A2E',
        'dark-card': '#16213E',
      },
      fontFamily: {
        heading: ['var(--font-be-vietnam-pro)', 'sans-serif'],
        body: ['var(--font-inter)', 'sans-serif'],
      },
    },
  },
  plugins: [],
};

export default config;
