package exporter

type Exporter interface {
	Format(relativePath string, content string, language string) string
}
