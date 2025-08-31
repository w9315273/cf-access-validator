package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
)

var (
	teamDomain string
	appMap     map[string][]string
	verifier   *oidc.IDTokenVerifier
)

func parseAppMap() {
	raw := os.Getenv("APP_MAP")
	if raw == "" {
		log.Fatalf("missing env APP_MAP (format: name=aud[,aud];name2=aud;...)")
	}
	appMap = make(map[string][]string)
	for _, pair := range strings.Split(raw, ";") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			log.Fatalf("bad APP_MAP entry: %q", pair)
		}
		name := strings.TrimSpace(kv[0])
		if name == "" {
			log.Fatalf("bad APP_MAP entry (empty name): %q", pair)
		}
		var auds []string
		for _, a := range strings.Split(kv[1], ",") {
			a = strings.TrimSpace(a)
			if a != "" {
				auds = append(auds, a)
			}
		}
		if len(auds) == 0 {
			log.Fatalf("bad APP_MAP entry (no aud): %q", pair)
		}
		appMap[name] = auds
	}
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}

func initVerifier(ctx context.Context) {
	issuer := "https://" + teamDomain

	httpClient := &http.Client{Timeout: 10 * time.Second}
	ctx = oidc.ClientContext(ctx, httpClient)

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatalf("init OIDC provider failed: %v", err)
	}

	verifier = provider.Verifier(&oidc.Config{
		SkipClientIDCheck:    true,
		SupportedSigningAlgs: []string{"RS256"},
	})
}

func getToken(r *http.Request) (string, error) {
	if t := r.Header.Get("Cf-Access-Jwt-Assertion"); t != "" {
		return t, nil
	}
	return "", errors.New("missing Cf-Access-Jwt-Assertion and CF_Authorization")
}

func audienceAllowed(tokenAUD []string, allowed []string) bool {
	if len(tokenAUD) == 0 || len(allowed) == 0 {
		return false
	}
	whitelist := make(map[string]struct{}, len(allowed))
	for _, a := range allowed {
		whitelist[a] = struct{}{}
	}
	for _, v := range tokenAUD {
		if _, ok := whitelist[v]; ok {
			return true
		}
	}
	return false
}

func validate(w http.ResponseWriter, r *http.Request) {
	app := r.Header.Get("X-Required-App")
	if app == "" {
		http.Error(w, "missing X-Required-App", http.StatusForbidden)
		return
	}
	allowed, ok := appMap[app]
	if !ok {
		http.Error(w, "unknown app", http.StatusForbidden)
		return
	}

	raw, err := getToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	idt, err := verifier.Verify(ctx, raw)
	if err != nil {
		http.Error(w, "invalid token: "+err.Error(), http.StatusForbidden)
		return
	}

	if !audienceAllowed(idt.Audience, allowed) {
		http.Error(w, "invalid token: audience not allowed for app", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func healthz(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }

func main() {
	teamDomain = mustEnv("TEAM_DOMAIN")
	parseAppMap()
	initVerifier(context.Background())

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:9000"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/validate", validate)
	mux.HandleFunc("/healthz", healthz)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
