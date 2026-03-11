package main

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	zgubagov "github.com/kacperfilipiuk/zguba-gov"
	"github.com/kacperfilipiuk/zguba-gov/internal/config"
	"github.com/kacperfilipiuk/zguba-gov/internal/database"
	"github.com/kacperfilipiuk/zguba-gov/internal/handler"
	"github.com/kacperfilipiuk/zguba-gov/internal/municipality"
	"github.com/kacperfilipiuk/zguba-gov/internal/repository"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("database:", err)
	}
	defer func() { _ = db.Close() }()

	munSvc := municipality.NewService(db)

	repo := repository.NewFoundItemRepo(db)
	apiH := handler.NewAPIHandler(repo)
	odataH := handler.NewODataHandler(repo)
	metaH := handler.NewMetadataHandler()

	templateFS, err := fs.Sub(zgubagov.WebFS, "web")
	if err != nil {
		log.Fatal("template fs:", err)
	}

	pagesH, err := handler.NewPagesHandler(templateFS, repo, munSvc)
	if err != nil {
		log.Fatal("pages handler:", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(_ string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	staticFS, _ := fs.Sub(zgubagov.WebFS, "web/static")
	r.StaticFS("/static", http.FS(staticFS))

	// HTML pages
	r.GET("/", pagesH.Index)
	r.POST("/steps/next", pagesH.NextStep)
	r.POST("/steps/prev", pagesH.PrevStep)
	r.GET("/autocomplete", pagesH.Autocomplete)
	r.POST("/submit", pagesH.Submit)
	r.GET("/export/json", pagesH.ExportJSON)
	r.GET("/export/csv", pagesH.ExportCSV)

	// REST API
	r.GET("/api/found-items", apiH.ListItems)
	r.POST("/api/found-items", apiH.CreateItem)
	r.GET("/api/found-items/categories/list", apiH.CategoriesList)
	r.GET("/api/found-items/:id", apiH.GetItem)
	r.PUT("/api/found-items/:id", apiH.UpdateItem)
	r.DELETE("/api/found-items/:id", apiH.DeleteItem)
	r.GET("/api/stats", apiH.Stats)
	r.GET("/health", apiH.Health)

	// OData
	r.GET("/odata/FoundItems", odataH.Query)
	r.GET("/odata/$metadata", odataH.Metadata)

	// Metadata
	r.GET("/metadata", metaH.Catalog)
	r.GET("/metadata/distribution/:id", metaH.Distribution)

	log.Printf("Starting server on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("server:", err)
	}
}
