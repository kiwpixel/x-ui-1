package service

import (
    "fmt"
    "math/rand"
    "strconv"
    "time"
    "x-ui/database"
    "x-ui/database/model"
    "x-ui/util/common"
    "x-ui/util/random"
    "x-ui/xray"

    "gorm.io/gorm"
)

// InboundService 入站服务结构体
type InboundService struct {
    db *gorm.DB
}

// NewInboundService 创建一个新的入站服务实例
func NewInboundService() *InboundService {
    return &InboundService{
        db: database.GetDB(),
    }
}

// BatchAddInbounds 批量添加入站的核心方法
func (s *InboundService) BatchAddInbounds(startPort, count int, username, password string) error {
    inbounds := make([]*model.Inbound, 0, count)
    for i := 0; i < count; i++ {
        port := startPort + i
        // 如果用户名未提供，则随机生成
        if username == "" {
            username = random.Seq(8)
        }
        // 如果密码未提供，则随机生成
        if password == "" {
            password = random.Seq(8)
        }
        // 构建入站配置的 settings 字段
        settings := fmt.Sprintf(`{
            "accounts": [
                {
                    "user": "%s",
                    "pass": "%s"
                }
            ]
        }`, username, password)
        // 创建入站配置实例
        inbound := &model.Inbound{
            Up:           0,
            Down:         0,
            Total:        0,
            Remark:       fmt.Sprintf("批量生成 - 端口 %d", port),
            Enable:       true,
            Listen:       "127.0.0.1",
            Port:         port,
            Protocol:     model.Http, // 可根据需要修改协议
            Settings:     settings,
            StreamSettings: "",
            Sniffing:     "",
        }
        inbounds = append(inbounds, inbound)
    }

    // 检查端口是否已存在
    for _, inbound := range inbounds {
        exist, err := s.checkPortExist(inbound.Port, 0)
        if err != nil {
            return err
        }
        if exist {
            return common.NewError("端口已存在:", inbound.Port)
        }
    }

    // 开启数据库事务
    tx := s.db.Begin()
    var err error
    defer func() {
        if err == nil {
            tx.Commit()
        } else {
            tx.Rollback()
        }
    }()

    // 批量保存入站配置到数据库
    for _, inbound := range inbounds {
        err = tx.Save(inbound).Error
        if err != nil {
            return err
        }
    }

    return nil
}

// checkPortExist 检查端口是否已存在
func (s *InboundService) checkPortExist(port int, ignoreID int) (bool, error) {
    var count int64
    query := s.db.Model(&model.Inbound{}).Where("port = ?", port)
    if ignoreID > 0 {
        query = query.Where("id != ?", ignoreID)
    }
    err := query.Count(&count).Error
    return count > 0, err
}

// GetInbounds 获取所有入站配置
func (s *InboundService) GetInbounds() ([]*model.Inbound, error) {
    var inbounds []*model.Inbound
    err := s.db.Find(&inbounds).Error
    return inbounds, err
}
