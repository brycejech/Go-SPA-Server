package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func mapKeyValueFile(fileContent string) (map[string]string, error) {
	kvMap := map[string]string{}
	lines := []string{}

	scanner := bufio.NewScanner(strings.NewReader(fileContent))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for lineNo, line := range lines {
		before, after, found := strings.Cut(line, "=")
		if !found {
			return nil, fmt.Errorf("line %v missing '='", lineNo)
		}
		kvMap[before] = after
	}

	return kvMap, nil
}

func swapRuntimeEnv(envMap map[string]string, requireAll bool) map[string]string {
	runtimeEnv := map[string]string{}

	for key, oldVal := range envMap {
		newVal := os.Getenv(key)
		if newVal == "" {
			if requireAll {
				log.Fatal(fmt.Errorf("error: runtime env missing value for '%s'", key))
			}
			fmt.Printf("warning: runtime env missing value for '%s'\n", key)
			continue
		}
		runtimeEnv[oldVal] = newVal
	}

	return runtimeEnv
}

type serverArgs struct {
	port            string
	staticDir       string
	embeddedEnvFile string
	requireAllVars  bool
}

func loadServerArgs() *serverArgs {
	var (
		port,
		staticDir,
		embeddedEnvFile string
		requireAllVars,
		ok bool
	)

	if port, ok = getOSArg("port"); !ok {
		port = "8000"
		fmt.Println("warning: no --port specified, defaulting to '8000'")
	}

	if staticDir, ok = getOSArg("staticDir"); !ok {
		fmt.Println("warning: no --staticDir specified, all 404s will serve index.html")
	}

	if embeddedEnvFile, ok = getOSArg("embeddedEnvFile"); !ok {
		fmt.Println("warning: no --embeddedEnvFile specified, runtime env cannot be swapped in")
	}

	requireAllVarsStr, ok := getOSArg("requireAllVars")
	if !ok {
		fmt.Println("warning: no --requireAllVars argument specified, defaulting to false")
	} else if strings.EqualFold("true", requireAllVarsStr) || requireAllVarsStr == "1" {
		requireAllVars = true
	}

	return &serverArgs{
		port:            port,
		staticDir:       staticDir,
		embeddedEnvFile: embeddedEnvFile,
		requireAllVars:  requireAllVars,
	}
}

func getOSArg(name string) (val string, ok bool) {
	for _, item := range os.Args {
		item := strings.TrimPrefix(item, "--")

		var (
			arg, val string
			ok       bool
		)
		if arg, val, ok = strings.Cut(item, "="); !ok {
			continue
		}

		if arg == name {
			return val, true
		}
	}

	return "", false
}
