package app

import (
	"errors"
	"flag"
	"os"
	"time"

	"be-file-uploader/internal/database"
	"be-file-uploader/internal/handler/auth"
	"be-file-uploader/internal/handler/developer"
	"be-file-uploader/internal/handler/image"
	"be-file-uploader/internal/handler/invite"
	"be-file-uploader/internal/handler/user"
	"be-file-uploader/internal/repository/mysql"
	authSrv "be-file-uploader/internal/service/auth"
	storageSrv "be-file-uploader/internal/service/image"
	inviteSrv "be-file-uploader/internal/service/invite"
	userSrv "be-file-uploader/internal/service/user"
	r2 "be-file-uploader/pkg/storage"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gookit/color"
	"github.com/gookit/slog"

	"github.com/bytedance/sonic"
)

func CreateApp() (app *fiber.App, db *database.Database, err error) {
	setupLogger()

	debug := flag.Bool("debug", false, "debug app with display of incoming requests")
	flag.Parse()

	db, err = database.CreateDatabase()
	if err != nil {
		return nil, nil, err
	}

	webDB := mysql.NewRepository(db.Web)

	app = fiber.New(fiber.Config{
		ServerHeader:  "",
		StrictRouting: true,
		CaseSensitive: true,
		Immutable:     true,
		Concurrency:   256 * 1024,
		ReadTimeout:   10 * time.Second,
		WriteTimeout:  10 * time.Second,
		IdleTimeout:   30 * time.Second,
		ProxyHeader:   fiber.HeaderXForwardedFor,
		GETOnly:       false,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			var e *fiber.Error
			if errors.As(err, &e) {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
					"code":  e.Code,
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": err.Error(),
			})
		},
		DisableKeepalive:             false,
		DisableDefaultDate:           false,
		DisableDefaultContentType:    false,
		DisableHeaderNormalizing:     false,
		AppName:                      "FileUploader",
		StreamRequestBody:            true,
		DisablePreParseMultipartForm: true,
		JSONEncoder:                  sonic.Marshal,
		JSONDecoder:                  sonic.Unmarshal,
		ColorScheme:                  fiber.Colors{},
	})

	if *debug {
		slog.Info("Debug started. Incoming requests will be displayed in app log")

		app.Use(func(c fiber.Ctx) error {
			start := time.Now()
			err := c.Next()

			stop := time.Since(start)

			slog.WithData(slog.M{
				"method":   c.Method(),
				"path":     c.Path(),
				"status":   c.Response().StatusCode(),
				"latency":  stop.String(),
				"ip":       c.IP(),
				"ua":       c.Get("X-User-Agent"),
				"body":     string(c.Body()),
				"query":    c.Queries(),
				"response": c.Response(),
			}).Info("Inbound request")

			return err
		})
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080", "https://uploader.dontkillme.lol"},
		AllowHeaders:     []string{"Origin, Content-Type, Accept, Authorization, Cache-Control, X-Request-Fingerprint, X-User-Agent"},
		AllowMethods:     []string{"GET, POST, PUT, DELETE, PATCH, OPTIONS"},
		AllowCredentials: true,
	}))

	storage, err := r2.NewStorage(os.Getenv("R2_ACCESS_KEY"), os.Getenv("R2_SECRET_KEY"), os.Getenv("R2_BUCKET"), os.Getenv("R2_PUBLIC_URL"))

	authService := authSrv.NewService(webDB)
	authHandler := auth.NewHandler(authService)

	userService := userSrv.NewService(webDB)
	userHandler := user.NewHandler(userService)

	inviteService := inviteSrv.NewService(webDB)
	inviteHandler := invite.NewHandler(inviteService, webDB)

	storageService := storageSrv.NewService(webDB, storage)
	storageHandler := image.NewHandler(storageService, webDB)

	developerHandler := developer.NewHandler(webDB)

	api := app.Group("/v1/api")
	public := api.Group("/public")
	private := api.Group("/private").Use(auth.Middleware(authService, webDB))

	authHandler.RegisterPublicRoutes(public)
	authHandler.RegisterPrivateRoutes(private)
	inviteHandler.RegisterPrivateRoutes(private)
	userHandler.RegisterPrivateRoutes(private)
	storageHandler.RegisterPrivateRoutes(private)
	developerHandler.RegisterPublicRoutes(public)

	if *debug {
		slog.Info("Registering routes...")
		for _, route := range app.GetRoutes() {
			if route.Method != "HEAD" {
				slog.Infof("Mapped [%s] -> %s", route.Method, route.Path)
			}
		}
	}

	return app, db, nil
}

func setupLogger() {
	f := slog.NewTextFormatter()
	f.EnableColor = true
	f.TimeFormat = "02/01/2006 15:04:05.000"

	f.ColorTheme = map[slog.Level]color.Color{
		slog.DebugLevel: color.FgBlue,
		slog.InfoLevel:  color.FgCyan,
		slog.WarnLevel:  color.FgYellow,
		slog.ErrorLevel: color.FgRed,
		slog.FatalLevel: color.FgMagenta,
	}

	slog.SetFormatter(f)
}
