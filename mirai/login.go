package mirai

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
)

var (
	errNoCredential = errors.New("uin/password not specified")
)

func (c *Core) Login() error {

	return nil
}

// read/write/encrypt password
func (c *Core) prepareCredential() error {
	if c.cfg.Account.Uin == 0 ||
		(len(c.cfg.Account.Password) == 0 && len(c.cfg.Account.PasswordEncrypted) == 0) {
		return errNoCredential
	}

	// encrypt password
	if c.cfg.Account.Encrypt && len(c.cfg.Account.Password) > 0 {
		c.pwHash = md5.Sum([]byte(c.cfg.Account.Password))
		log.Infof("密码加密已启用, 请输入Key对密码进行加密: (Enter 提交)")
		byteKey, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("can not read input password: %v", err)
		}
	}

	return nil
}
