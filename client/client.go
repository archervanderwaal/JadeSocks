package client

import (
	"fmt"
	"github.com/archervanderwaal/JadeSocks/cipher"
	"github.com/archervanderwaal/JadeSocks/config"
	JadeSocks "github.com/archervanderwaal/JadeSocks/core"
	"github.com/archervanderwaal/JadeSocks/utils"
	"log"
	"net"
)

type Client struct {
	Cipher     *cipher.Cipher
	LocalAddr  *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

func NewClient(config *config.Config) (*Client, error) {
	localAddr, err := net.ResolveTCPAddr("tcp", config.ListenAddr)
	if err != nil {
		return nil, err
	}
	remoteAddr, err := net.ResolveTCPAddr("tcp", config.RemoteAddr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Cipher: &cipher.Cipher{
			Algorithm: &cipher.AesCryptoAlgorithm{
				Key: utils.NewKey(16, config.Password),
			},
		},
		LocalAddr:  localAddr,
		RemoteAddr: remoteAddr,
	}, nil
}

func (client *Client) Listen(listenHandler func(listenerAddr *net.TCPAddr)) error {
	return JadeSocks.ListenEncryptedTCP(client.LocalAddr, client.Cipher, client.handleConn, listenHandler)
}

func (client *Client) handleConn(userConn *JadeSocks.SecureTCPConn) {
	defer userConn.Close()

	proxyServer, err := JadeSocks.DialEncryptedTCP(client.RemoteAddr, client.Cipher)
	if err != nil {
		log.Println(err)
		return
	}
	defer proxyServer.Close()

	fmt.Println("client=>处理连接")

	// 进行转发
	// 从 proxyServer 读取数据发送到 localUser
	go func() {
		err := proxyServer.DecodeCopy(userConn)
		if err != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			userConn.Close()
			proxyServer.Close()
		}
	}()
	// 从 localUser 发送数据发送到 proxyServer，这里因为处在翻墙阶段出现网络错误的概率更大
	_ = userConn.EncodeCopy(proxyServer)
}
