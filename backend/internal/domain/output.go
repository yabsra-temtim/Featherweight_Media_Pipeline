package domain

// Output holds the result metadata for processed media.
type Output struct {
	Format   string
	Path     string // Cloudinary secure URL (https://res.cloudinary.com/...)
	PublicID string // Cloudinary public_id, used for deletion
	Size     int64
	Width    int
	Height   int
}
