package main
import(
	//"log"
	//"fmt"
	"expvar"
	"runtime"
	"time"
	"os"
	"github.com/joho/godotenv"
	 "go.uber.org/zap"
	 "github.com/go-redis/redis/v8"
	"github.com/Prasanthi-Peram/pigee-connect/internal/env"
	"github.com/Prasanthi-Peram/pigee-connect/internal/store"
	"github.com/Prasanthi-Peram/pigee-connect/internal/db"
	"github.com/Prasanthi-Peram/pigee-connect/internal/mailer"
	"github.com/Prasanthi-Peram/pigee-connect/internal/auth"
	"github.com/Prasanthi-Peram/pigee-connect/internal/store/cache"
	"github.com/Prasanthi-Peram/pigee-connect/internal/ratelimiter"
)

const version="0.0.1"

//	@title			PigeeConnect API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main(){
	
	godotenv.Load()
	cfg:=config{
		addr:env.GetString("ADDR",":8080"),
		apiURL:env.GetString("EXTERNAL_URL","localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL","http://localhost:5174"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR","postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS",30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS",30),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME","15m"),
		},
		redisCfg: redisConfig{
			addr: env.GetString("REDIS_ADDR","localhost:6379"),
			pw: env.GetString("REDIS_PW",""),
			db: env.GetInt("REDIS_DB",0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		env: env.GetString("ENV","development"),
		mail:mailConfig{
			exp: time.Hour*24*3,
			fromEmail: os.Getenv("FROM_EMAIL"),

			sendGrid: sendGridConfig{
				apiKey:os.Getenv("SENDGRID_API_KEY"),
				//fromEmail: env.GetString("FROM_EMAIL",""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER","admin"),
				pass: env.GetString("AUTH_BASIC_PASS","admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET","example"),
				exp: time.Hour*24*3,
				iss: "pigeeconnect",
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}
	//fmt.Println("Loaded API Key from config:", cfg.mail.sendGrid.apiKey)
	//fmt.Println("Loaded API Key from config:", cfg.mail.fromEmail)
    
	//Logger
	logger:=zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//Database
	db,err:= db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err!=nil{
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("db connection pool established")

	//Cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)

		logger.Info("redis cache connection established")

		defer rdb.Close()


	}

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	store:= store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)
	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey,cfg.mail.fromEmail)
	jwtAuthenticator:= auth.NewJWTAuthenticator(cfg.auth.token.secret,cfg.auth.token.iss,cfg.auth.token.iss)
	
	app:=&application{
		config:cfg,
		store: store,
		cacheStorage: cacheStorage,
		logger: logger,
		mailer: mailer,
		authenticator: jwtAuthenticator,
		rateLimiter: rateLimiter,
	}

	//Server Metrics
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux:=app.mount()
	logger.Fatal(app.run(mux))
}