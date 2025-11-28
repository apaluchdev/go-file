package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/apaluchdev/go-file/docs"
)

// @title           File Serving API
// @version         1.0
// @description     This is a sample server for serving files.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /

func main() {
	       router := gin.Default()
	       // CORS middleware for React app on localhost:80
	       router.Use(func(c *gin.Context) {
		       c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		       c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		       c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		       c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		       if c.Request.Method == "OPTIONS" {
			       c.AbortWithStatus(204)
			       return
		       }
		       c.Next()
	       })

	// Ensure storage directory exists
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}
	log.Printf("Using storage path: %s", storagePath)
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		log.Printf("Creating storage directory: %s", storagePath)
		os.MkdirAll(storagePath, 0755)
	}

		router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	       router.GET("/api/files/:pin", func(c *gin.Context) {
		       listFiles(c, storagePath)
	       })
	       router.POST("/api/files/:pin", func(c *gin.Context) {
		       uploadFile(c, storagePath)
	       })
	       router.GET("/api/files/:pin/:filename", func(c *gin.Context) {
		       downloadFile(c, storagePath)
	       })

	log.Println("Starting server on :8080")
		log.Println("Swagger documentation available at http://localhost:8080/api/swagger/index.html")
	router.Run(":8080")
}

// listFiles godoc
// @Summary      List files
// @Description  Get existing files for a PIN
// @Tags         files
// @Accept       json
// @Produce      json
// @Param        pin path string true "PIN (6-8 digits)"
// @Success      200  {array}   string
// @Router       /files/{pin} [get]
func listFiles(c *gin.Context, storagePath string) {
	pin := c.Param("pin")
	log.Printf("[GET /files/%s] Listing files for PIN: %s", pin, pin)
	pinPath := filepath.Join(storagePath, pin)
	
	// Create PIN directory if it doesn't exist
	if _, err := os.Stat(pinPath); os.IsNotExist(err) {
		log.Printf("[GET /files/%s] Creating new directory for PIN: %s", pin, pin)
		os.MkdirAll(pinPath, 0755)
	}
	
	files, err := os.ReadDir(pinPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type FileInfo struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
	}

	var fileInfos []FileInfo
	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				log.Printf("[GET /files/%s] Error getting info for file %s: %v", c.Param("pin"), file.Name(), err)
				continue
			}
			fileInfos = append(fileInfos, FileInfo{Name: file.Name(), Size: info.Size()})
		}
	}

	log.Printf("[GET /files/%s] Found %d files", c.Param("pin"), len(fileInfos))
	c.JSON(http.StatusOK, gin.H{"files": fileInfos})
}

// uploadFile godoc
// @Summary      Upload a file
// @Description  Upload a file to the storage for a specific PIN
// @Tags         files
// @Accept       multipart/form-data
// @Produce      json
// @Param        pin path string true "PIN (6-8 digits)"
// @Param        file formData file true "File to upload"
// @Success      200  {object}  map[string]string
// @Router       /files/{pin} [post]
func uploadFile(c *gin.Context, storagePath string) {
	pin := c.Param("pin")
	log.Printf("[POST /files/%s] Upload request received for PIN: %s", pin, pin)
	pinPath := filepath.Join(storagePath, pin)
	
	// Create PIN directory if it doesn't exist
	if _, err := os.Stat(pinPath); os.IsNotExist(err) {
		log.Printf("[POST /files/%s] Creating new directory for PIN: %s", pin, pin)
		os.MkdirAll(pinPath, 0755)
	}
	
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, filepath.Join(pinPath, filename)); err != nil {
		log.Printf("[POST /files/%s] Error saving file %s: %v", pin, filename, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[POST /files/%s] File uploaded successfully: %s", pin, filename)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("File %s uploaded successfully", filename)})
}

// downloadFile godoc
// @Summary      Download a file
// @Description  Download a file by name for a specific PIN
// @Tags         files
// @Produce      octet-stream
// @Param        pin path string true "PIN (6-8 digits)"
// @Param        filename path string true "Filename"
// @Success      200  {file}    file
// @Router       /files/{pin}/{filename} [get]
func downloadFile(c *gin.Context, storagePath string) {
	pin := c.Param("pin")
	filename := c.Param("filename")
	log.Printf("[GET /files/%s/%s] Download request for file: %s (PIN: %s)", pin, filename, filename, pin)
	filePath := filepath.Join(storagePath, pin, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("[GET /files/%s/%s] File not found: %s", pin, filename, filePath)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	log.Printf("[GET /files/%s/%s] Serving file: %s", pin, filename, filename)
	c.File(filePath)
}
