package main

import (
	"embed"
	"fmt"
	"log"
	"sync"

	"github.com/brycejech/go-spa-server/internal"
)

//go:embed all:artifact
var embeddedContent embed.FS

func main() {
	args := loadServerArgs()

	runtimeEnv := map[string]string{}

	if args.embeddedEnvFile != "" {
		embeddedEnvFile, err := embeddedContent.ReadFile(fmt.Sprintf("artifact/%s", args.embeddedEnvFile))
		if err != nil {
			log.Fatal(fmt.Errorf("error opening embeddedEnvFile '%s': %w", args.embeddedEnvFile, err))
		}

		embeddedEnv, err := mapKeyValueFile(string(embeddedEnvFile))
		if err != nil {
			log.Fatal(fmt.Errorf("error parse error in embeddedEnvFile: %w", err))
		}

		runtimeEnv = swapRuntimeEnv(embeddedEnv, args.requireAllVars)
	}

	cache := internal.MustCreateFileCache(embeddedContent, runtimeEnv)
	server := internal.NewSPAServer(
		args.port,
		cache,
		args.staticDir,
	)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			wg.Done()
			log.Fatal(err)
		}
	}()

	fmt.Printf("Listening on localhost:%s\n", args.port)
	wg.Wait()
}
