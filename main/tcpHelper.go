package main

import (
	"net"
	"time"
)

type handleData func(from []byte) []byte


/*
获取连接
 */
func tryConnectOrLoop(ip string, d time.Duration) (*net.TCPConn) {
	var conn *net.TCPConn;
	tcpAddr, err := net.ResolveTCPAddr("tcp", ip)
	for conn == nil {
		conn, err = net.DialTCP("tcp", nil, tcpAddr);
		if err != nil {
			logE("connect %s error for %s", ip, err)
			time.Sleep(d)
		}
	}
	conn.SetKeepAlive(true)
	return conn;
}

/*
自动重连发送
 */
func autoSend(ip string, in chan []byte, d time.Duration) {
	var conn *net.TCPConn;
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	for data := range in {
		for {
			if conn == nil {
				conn = tryConnectOrLoop(ip, d)
			}
			num, err := conn.Write(data)
			if err == nil && num == len(data) {
				break
			}
			logE("send data '%s' ip= %s error for %s", string(data[:]), ip, err)
			conn.Close()
			conn = nil
		}
	}
}

/*
自动重连接收
 */
func autoGet(ip string, out chan []byte, d time.Duration) {
	var conn *net.TCPConn;
	var buffer = make([]byte, 1024)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	for {
		if conn == nil {
			conn = tryConnectOrLoop(ip, d)
		}
		for {
			num, err := conn.Read(buffer)
			if err == nil {
				data := buffer[:num]
				out <- data
				continue
			}
			logE("get data error ip= %s error for %s", ip, err)
			conn.Close()
			conn = nil
			break
		}
	}
}

/**
转发
 */
func transfer(from string, to string, d time.Duration, handleData handleData) {
	var sendData = make(chan []byte, 10)
	var receivedData = make(chan []byte, 10)
	go autoGet(from, receivedData, d)
	go autoSend(to, sendData, d)
	for data := range receivedData {
		sendData <- handleData(data)
	}
}

