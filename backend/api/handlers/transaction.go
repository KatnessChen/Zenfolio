package handlers

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/transaction-tracker/backend/config"
	"github.com/transaction-tracker/backend/internal/ai"
	"github.com/transaction-tracker/backend/internal/constants"
	"github.com/transaction-tracker/backend/internal/types"
)

// ExtractTransactionsHandler handles the image upload and transaction extraction
func ExtractTransactionsHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form: " + err.Error()})
			return
		}
		files := form.File["images"] // Expecting files under the field name "images"

		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No images uploaded. Please upload images under the 'images' field."})
			return
		}

		var imageInputs []types.ImageInput
		var filesToClose []multipart.File

		// Ensure all opened files are closed after processing
		defer func() {
			for _, f := range filesToClose {
				_ = f.Close()
			}
		}()

		for _, fileHeader := range files {
			// Log the received file
			// log.Printf("Received file: %s, size: %d, header: %#v", fileHeader.Filename, fileHeader.Size, fileHeader.Header)

			src, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded image " + fileHeader.Filename + ": " + err.Error()})
				return
			}

			// Collect files for deferred closure outside the loop
			filesToClose = append(filesToClose, src)

			imageInputs = append(imageInputs, types.ImageInput{
				Data:     src,
				Filename: fileHeader.Filename,
				MimeType: fileHeader.Header.Get(constants.ContentTypeHeader),
			})
		}

		// Initialize AI client
		// Note: In a real application, you'd likely initialize the AI client once and reuse it,
		// or use a dependency injection framework. For simplicity here, we create it on each request.
		// Pass the main application config directly to the factory
		aiClient, err := ai.NewClient(cfg) // Using the factory from backend/internal/ai/factory.go
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize AI client: " + err.Error()})
			return
		}
		defer aiClient.Close()

		// Log imageInput details
		// log.Printf("ImageInputs: %d files", len(imageInputs))

		extractResp, err := aiClient.ExtractTransactions(c.Request.Context(), imageInputs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract transactions: " + err.Error()})
			return
		}

		if !extractResp.Success {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction extraction unsuccessful: " + extractResp.Message})
			return
		}

		c.JSON(http.StatusOK, extractResp)
	}
}
