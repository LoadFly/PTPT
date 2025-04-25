package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/leslieleung/ptpt/internal/interract"
	"github.com/sirupsen/logrus"
)

func InitLogger() error {
	// 获取日志文件路径
	logDir := filepath.Join(interract.GetPTPTDir(), "logs")

	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 设置日志文件
	logFile := filepath.Join(logDir, "ptpt.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	// 设置日志输出到文件
	logrus.SetOutput(file)

	// 设置日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置日志级别
	logrus.SetLevel(logrus.DebugLevel)

	return nil
}
