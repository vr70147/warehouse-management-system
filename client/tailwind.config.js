/** @type {import('tailwindcss').Config} */
import tailwindcssAnimate from "tailwindcss-animate";

export default {
	content: [
		"./pages/**/*.{js,ts,jsx,tsx,mdx}",
		"./components/**/*.{js,ts,jsx,tsx,mdx}",
		"./src/**/*.{js,ts,jsx,tsx}",
	],
	darkMode: 'class',
	theme: {
		extend: {
			colors: {
				border: "hsl(var(--border))",
				background: "hsl(var(--background))",
				foreground: "hsl(var(--foreground))",
				destructive: {
					DEFAULT: '#DC2626',
					foreground: '#FFFFFF',
				},
				primaryBlue: {
					DEFAULT: "#2563eb",
					foreground: "#ffffff",
				},
			},
			scale: {
				102: '1.02',
			}
		},
	},
	plugins: [tailwindcssAnimate],
}

