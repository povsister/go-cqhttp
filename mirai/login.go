package mirai

import (
	"crypto/md5"
	"errors"

	"github.com/Mrs4s/go-cqhttp/global"
	"github.com/Mrs4s/go-cqhttp/global/config"
)

var (
	errNoCredential = errors.New("uin/password not specified")
)

func (c *Core) Login() error {

	return nil
}

// read/write/encrypt password
func (c *Core) prepareCredential() error {
	// 如果不存在session 且 账号未设置或者密码未设置
	// 就当作没credential处理
	if !global.PathExists(config.DefaultSessionFile) &&
		(c.cfg.Account.Uin == 0 ||
			((len(c.cfg.Account.Password) == 0 && !c.cfg.Account.Encrypt) &&
				(len(c.cfg.Account.PasswordEncrypted) == 0) && c.cfg.Account.Encrypt)) {
		return errNoCredential
	}

	if global.PathExists(config.DefaultSessionFile) {

	}

	if len(c.cfg.Account.Password) > 0 {
		c.pwHash = md5.Sum([]byte(c.cfg.Account.Password))
	}
	// 有密码 不加密
	if !c.cfg.Account.Encrypt && len(c.cfg.Account.Password) > 0 {
		return nil
	}
	// 有密码 启用加密
	if c.cfg.Account.Encrypt && len(c.cfg.Account.Password) > 0 {
		// TODO: encrypt pw and remove it from config file
	}

	// should never reach here
	return errNoCredential
}
