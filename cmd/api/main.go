package main
import(
	"log"

	"github.com/prashanthi/social/internal/env"
)

func main(){
	cfg:=config{
		addr:env.GetString("ADDR",":8080"),
	}

	store:= store.NewStorage(db)
	app:=&application{
		config:cfg,
		store: store,
	}

	mux:=app.mount()
	log.Fatal(app.run(mux))
}