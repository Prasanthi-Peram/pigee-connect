package main
import(
	//"log"
	"fmt"
	"time"
	"os"
	"github.com/joho/godotenv"
	 "go.uber.org/zap"
	"github.com/Prasanthi-Peram/pigee-connect/internal/env"
	"github.com/Prasanthi-Peram/pigee-connect/internal/store"
	"github.com/Prasanthi-Peram/pigee-connect/internal/db"
	"github.com/Prasanthi-Peram/pigee-connect/internal/mailer"
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
		env: env.GetString("ENV","development"),
		mail:mailConfig{
			exp: time.Hour*24*3,
			fromEmail: os.Getenv("FROM_EMAIL"),

			sendGrid: sendGridConfig{
				apiKey:os.Getenv("SENDGRID_API_KEY"),
				//fromEmail: env.GetString("FROM_EMAIL",""),
			},
		},
	}
	fmt.Println("Loaded API Key from config:", cfg.mail.sendGrid.apiKey)
	fmt.Println("Loaded API Key from config:", cfg.mail.fromEmail)
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


	store:= store.NewStorage(db)
	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey,cfg.mail.fromEmail)
	app:=&application{
		config:cfg,
		store: store,
		logger: logger,
		mailer: mailer,
	}

	mux:=app.mount()
	logger.Fatal(app.run(mux))
}