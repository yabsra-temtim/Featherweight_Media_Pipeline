/** Line-art feather mark used as the brand signature and processing indicator. */
export default function FeatherMark({ className = 'w-5 h-5', animated = false }) {
    return (
        <svg
            viewBox="0 0 24 24"
            fill="none"
            className={`${className} ${animated ? 'animate-drift' : ''}`}
            aria-hidden="true"
        >
            <path
                d="M18.5 3C11 3 4 9 4 19.5c0 .8.05 1.6.16 2.34.09.6.85.78 1.2.28C8 15.9 12.7 12.6 18.3 11c.5-.14.5-.86 0-1-4.5.9-9 2.9-12.4 6.5C6.6 8.2 12.2 3.9 18.5 3.2c.6-.07.6-1.13 0-1.2z"
                fill="currentColor"
            />
            <path
                d="M11 14l4.5-4.5"
                stroke="currentColor"
                strokeWidth="1"
                strokeLinecap="round"
                opacity="0.45"
            />
        </svg>
    )
}
