package mirai

import (
	"os"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/go-cqhttp/coolq"
	"github.com/Mrs4s/go-cqhttp/global"
	"github.com/Mrs4s/go-cqhttp/global/config"
	log "github.com/sirupsen/logrus"
)

var (
	validCoolqMessage = map[string]struct{}{
		"string": {}, "array": {},
	}
)

// Abstraction layer for mirai core
type Core struct {
	c   *client.QQClient
	cfg *config.Config
	// md5sum of password
	pwHash [16]byte
	// session token
	session []byte
}

func NewCore(cfg *config.Config) (*Core, error) {
	core := &Core{
		c:   client.NewClientEmpty(),
		cfg: cfg,
	}
	return core, core.initialize()
}

func (c *Core) initialize() (err error) {
	todos := []func() error{
		c.setupDeviceInfo,
		c.setupLogForward,
		c.setupServerUpdate,
		c.setupMiscConfig,
	}
	for _, do := range todos {
		if err = do(); err != nil {
			return err
		}
	}
	return
}

func (c *Core) setupDeviceInfo() (err error) {
	if global.PathExists(config.DefaultDeviceJsonFile) {
		log.Info("将使用 device.json 内的设备信息运行Bot.")
		data, err := os.ReadFile(config.DefaultDeviceJsonFile)
		if err != nil {
			return
		}
		if err = client.SystemDeviceInfo.ReadJson(data); err != nil {
			return
		}
		return
	}
	// no device.json found
	log.Warn("虚拟设备信息不存在, 将自动生成随机设备.")
	client.GenRandomDevice()
	err = os.WriteFile(config.DefaultDeviceJsonFile, client.SystemDeviceInfo.ToJson(), 0644)
	if err != nil {
		log.Warnf("无法保存 device.json : %v", err)
		log.Warn("这将导致每次启动都使用完全随机的设备信息")
		// silence the err. do not treat it as fatal
		err = nil
		return
	}
	log.Info("已生成设备信息并保存到 device.json 文件.")
	return
}

func (c *Core) setupLogForward() error {
	c.c.OnLog(func(_ *client.QQClient, e *client.LogEvent) {
		switch e.Type {
		case "INFO":
			log.Info("Protocol -> " + e.Message)
		case "ERROR":
			log.Error("Protocol -> " + e.Message)
		case "DEBUG":
			log.Debug("Protocol -> " + e.Message)
		}
	})
	return nil
}

func (c *Core) setupServerUpdate() error {
	c.c.OnServerUpdated(func(_ *client.QQClient, e *client.ServerUpdatedEvent) bool {
		if !c.cfg.Account.UseSSOAddress {
			log.Info("收到服务器地址更新通知, 根据配置文件已忽略.")
			return false
		}
		log.Info("收到服务器地址更新通知, 将在下一次重连时应用.")
		return true
	})

	// load customized server address if possible
	if !global.PathExists(config.DefaultServerAddressFile) {
		return nil
	}
	log.Info("检测到 address.txt 文件. 将覆盖目标IP.")
	addr := global.ReadAddrFile("address.txt")
	log.Infof("读取到 %d 个自定义地址.", len(addr))
	if len(addr) > 0 {
		c.c.SetCustomServer(addr)
	}

	return nil
}

func (c *Core) setupMiscConfig() error {
	// for global
	global.Proxy = c.cfg.Message.ProxyRewrite

	// for coolq
	coolq.IgnoreInvalidCQCode = c.cfg.Message.IgnoreInvalidCQCode
	coolq.SplitURL = c.cfg.Message.FixURL
	coolq.ForceFragmented = c.cfg.Message.ForceFragment
	coolq.RemoveReplyAt = c.cfg.Message.RemoveReplyAt
	coolq.ExtraReplyData = c.cfg.Message.ExtraReplyData
	if _, ok := validCoolqMessage[c.cfg.Message.PostFormat]; !ok {
		log.Warnf("post-format 配置错误, 将自动使用 string")
		coolq.SetMessageFormat("string")
	} else {
		coolq.SetMessageFormat(c.cfg.Message.PostFormat)
	}
	return nil
}
