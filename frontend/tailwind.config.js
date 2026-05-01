import daisyui from 'daisyui'

export default {
    content: ['./index.html', './src/**/*.{vue,js}'],
    theme: {
        extend: {
            fontFamily: {
                sans: ['"Space Grotesk"', 'system-ui', 'Segoe UI', 'Roboto', 'sans-serif'],
            },
        }
    },
    plugins: [daisyui],
    daisyui: {
        themes: [
            {
                catdiary: {
                    primary: '#B9FF66',
                    'primary-content': '#191A23',
                    secondary: '#191A23',
                    'secondary-content': '#FFFFFF',
                    accent: '#B9FF66',
                    'accent-content': '#191A23',
                    neutral: '#191A23',
                    'neutral-content': '#FFFFFF',
                    'base-100': '#FFFFFF',
                    'base-200': '#F3F3F3',
                    'base-300': '#F3F3F3',
                    'base-content': '#000000',
                    info: '#3ABFF8',
                    success: '#36D399',
                    warning: '#FBBD23',
                    error: '#F87272',
                },
            },
        ],
    },
}
