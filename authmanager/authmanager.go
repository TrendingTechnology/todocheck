package authmanager

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/preslavmihaylov/todocheck/authmanager/authstore"
	"github.com/preslavmihaylov/todocheck/config"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	githubAPITokenPage = "https://github.com/settings/tokens"
	gitlabAPITokenPage = "https://gitlab.com/profile/personal_access_tokens"
)

// AcquireToken stores the issue tracker's auth token based on the auth type specified
func AcquireToken(cfg *config.Local) error {
	switch cfg.Auth.Type {
	case config.AuthTypeNone:
		return nil
	case config.AuthTypeAPIToken:
		return acquireAPIToken(cfg)
	case config.AuthTypeOffline:
		return acquireOfflineToken(cfg.Auth)
	default:
		panic("unknown auth token type")
	}
}

func acquireAPIToken(cfg *config.Local) error {
	return acquireToken(cfg.Auth, cfg.Origin, func() ([]byte, error) {
		var targetPage string
		if cfg.IssueTracker == config.IssueTrackerGithub {
			targetPage = githubAPITokenPage
		} else {
			targetPage = gitlabAPITokenPage
		}

		fmt.Printf("Please go to %v, create a read-only access token & paste it here:\nToken: ", targetPage)
		return readPassword()
	})
}

func acquireOfflineToken(a *config.Auth) error {
	return acquireToken(a, a.OfflineURL, func() ([]byte, error) {
		fmt.Printf("Please go to %v and paste the offline token below:\nToken: ", a.OfflineURL)
		return readPassword()
	})
}

func acquireToken(authCfg *config.Auth, tokenKey string, promptCallback func() ([]byte, error)) error {
	store, err := authstore.CreateIfNotExists(authCfg.TokensCache, authstore.DefaultConfigPermissions)
	if err != nil {
		return fmt.Errorf("couldn't read auth tokens config: %w", err)
	}

	if store.Tokens[tokenKey] != "" {
		authCfg.Token = store.Tokens[tokenKey]
		return nil
	}

	tokenBs, err := promptCallback()
	if err != nil {
		return fmt.Errorf("couldn't acquire token: %w", err)
	}

	authCfg.Token = strings.TrimSpace(string(tokenBs))
	store.Tokens[tokenKey] = authCfg.Token
	store.Save(authCfg.TokensCache)

	return nil
}

// Make token input scriptable, while preserving the hidden prompt behavior for users
// https://github.com/golang/go/issues/19909#issuecomment-399409958
func readPassword() ([]byte, error) {
	if terminal.IsTerminal(syscall.Stdin) {
		return terminal.ReadPassword(syscall.Stdin)
	}

	reader := bufio.NewReader(os.Stdin)
	return reader.ReadBytes('\n')
}
