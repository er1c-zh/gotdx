package tdx

import (
	"gotdx/models"
	"time"
)

const (
	_defaultTCPAddress        = "119.147.212.81:7709"
	_defaultRetryTimes        = 3
	_defaultHeartbeatInterval = 15 * time.Second
)

type Options struct {
	TCPAddress        string   // 服务器地址
	TCPAddressPool    []string // 服务器地址池
	MaxRetryTimes     int      // 重试次数
	HeartbeatInterval time.Duration
	MsgCallback       func(models.ProcessInfo)
}

func defaultOptions() *Options {
	return &Options{
		TCPAddress:        _defaultTCPAddress,
		MaxRetryTimes:     _defaultRetryTimes,
		HeartbeatInterval: _defaultHeartbeatInterval,
		MsgCallback: func(pi models.ProcessInfo) {
			// do nothing
		},
	}
}

func applyOptions(opts ...Option) *Options {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

type Option func(options *Options)

var DefaultOption Option = func(options *Options) {}

func (Option Option) WithTCPAddress(tcpAddress string) Option {
	return func(o *Options) {
		o.TCPAddress = tcpAddress
		Option(o)
	}
}

func (Option Option) WithTCPAddressPool(ips ...string) Option {
	return func(o *Options) {
		o.TCPAddressPool = ips
		Option(o)
	}
}

func (Option Option) WithHeartbitInterval(interval time.Duration) Option {
	return func(o *Options) {
		o.HeartbeatInterval = interval
		Option(o)
	}
}

func (Option Option) WithMsgCallback(callback func(models.ProcessInfo)) Option {
	return func(o *Options) {
		o.MsgCallback = callback
		Option(o)
	}
}
