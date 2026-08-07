/** @type {import('tailwindcss').Config} */
export default {
    darkMode: 'class',
    content: ['./index.html', './src/**/*.{js,jsx}'],
    theme: {
        extend: {
            colors: {
                // ink/parchment/slate all resolve through CSS variables (see src/index.css)
                // so the same class names automatically re-theme in light vs. dark mode.
                ink: {
                    950: '#0B0F16',
                    900: 'rgb(var(--bg-0) / <alpha-value>)',
                    800: 'rgb(var(--bg-1) / <alpha-value>)',
                    700: 'rgb(var(--bg-2) / <alpha-value>)',
                    600: 'rgb(var(--border-0) / <alpha-value>)',
                    500: 'rgb(var(--border-1) / <alpha-value>)',
                },
                parchment: {
                    100: 'rgb(var(--text-0) / <alpha-value>)',
                    300: 'rgb(var(--text-2) / <alpha-value>)',
                },
                slate: {
                    400: 'rgb(var(--text-1) / <alpha-value>)',
                    500: 'rgb(var(--text-2) / <alpha-value>)',
                },
                quill: {
                    DEFAULT: '#E4A93B',
                    soft: '#F0C877',
                    dim: '#8A6A28',
                },
                moss: {
                    DEFAULT: '#5FB6A2',
                    dim: '#2F5C52',
                },
                clay: {
                    DEFAULT: '#E2685A',
                    dim: '#6B2E28',
                },
            },
            fontFamily: {
                display: ['"Fraunces"', 'ui-serif', 'Georgia', 'serif'],
                sans: ['"Inter"', 'ui-sans-serif', 'system-ui', 'sans-serif'],
                mono: ['"IBM Plex Mono"', 'ui-monospace', 'SFMono-Regular', 'monospace'],
            },
            keyframes: {
                drift: {
                    '0%': { transform: 'translateY(0) rotate(0deg)', opacity: '0.9' },
                    '50%': { transform: 'translateY(-6px) rotate(-4deg)', opacity: '1' },
                    '100%': { transform: 'translateY(0) rotate(0deg)', opacity: '0.9' },
                },
                shimmer: {
                    '0%': { backgroundPosition: '-200% 0' },
                    '100%': { backgroundPosition: '200% 0' },
                },
                'fade-up': {
                    '0%': { opacity: '0', transform: 'translateY(6px)' },
                    '100%': { opacity: '1', transform: 'translateY(0)' },
                },
            },
            animation: {
                drift: 'drift 2.4s ease-in-out infinite',
                shimmer: 'shimmer 1.8s linear infinite',
                'fade-up': 'fade-up 0.4s ease-out both',
            },
        },
    },
    plugins: [],
}

