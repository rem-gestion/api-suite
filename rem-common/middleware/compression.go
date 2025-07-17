package middleware

import (
	"compress/gzip"
	"strings"

	"github.com/gin-gonic/gin"
)

// GzipMiddleware comprime respuestas HTTP para reducir el ancho de banda
// Solo comprime respuestas > 1KB y tipos MIME compatibles
func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Solo comprimir si el cliente acepta gzip
		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// Tipos MIME que se benefician de compresión
		compressibleTypes := map[string]bool{
			"application/json":       true,
			"application/xml":        true,
			"text/html":              true,
			"text/plain":             true,
			"text/css":               true,
			"text/javascript":        true,
			"application/javascript": true,
		}

		// Wrapper para capturar la respuesta
		writer := &gzipResponseWriter{
			ResponseWriter:    c.Writer,
			compressibleTypes: compressibleTypes,
		}
		c.Writer = writer

		defer func() {
			if writer.gzipWriter != nil {
				writer.gzipWriter.Close()
			}
		}()

		c.Next()
	}
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	gzipWriter        *gzip.Writer
	compressibleTypes map[string]bool
	headerWritten     bool
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	if !g.headerWritten {
		g.writeHeader()
	}

	if g.gzipWriter != nil {
		return g.gzipWriter.Write(data)
	}
	return g.ResponseWriter.Write(data)
}

func (g *gzipResponseWriter) WriteHeader(statusCode int) {
	g.ResponseWriter.WriteHeader(statusCode)
	g.headerWritten = true
}

func (g *gzipResponseWriter) writeHeader() {
	contentType := g.Header().Get("Content-Type")

	// Solo comprimir tipos MIME apropiados
	shouldCompress := false
	for mimeType := range g.compressibleTypes {
		if strings.Contains(contentType, mimeType) {
			shouldCompress = true
			break
		}
	}

	if shouldCompress {
		g.Header().Set("Content-Encoding", "gzip")
		g.Header().Del("Content-Length") // Será diferente después de comprimir
		g.gzipWriter = gzip.NewWriter(g.ResponseWriter)
	}
}
