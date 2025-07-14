package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"os"
	"strings"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	host "github.com/libp2p/go-libp2p/core/host"
	network "github.com/libp2p/go-libp2p/core/network"
	peer "github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

const chatProto = "/p2pchat/1.0.0"

var activeStreams = make(map[peer.ID]network.Stream)
var peerByID = map[string]peer.AddrInfo{}
var peerByPort = map[string]peer.ID{}

type mdnsNotifee struct {
	h host.Host
}

func (n *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Print("Найден peer через mDNS:", pi.ID.String()+"\n"+"> ")
	n.h.Peerstore().AddAddrs(pi.ID, pi.Addrs, peerstore.PermanentAddrTTL)
	peerByID[pi.ID.String()] = pi

	for _, addr := range pi.Addrs {
		if tcpPort, err := addr.ValueForProtocol(multiaddr.P_TCP); err == nil {
			peerByPort[tcpPort] = pi.ID
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Создаём узел с криптографическим ключом Ed25519
	priv, _, _ := crypto.GenerateEd25519Key(nil)
	h, _ := libp2p.New(
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0",
			"/ip4/0.0.0.0/udp/0/quic",
		),
	)

	// 2. Показываем свои многоадреса
	fmt.Println("=== P2P Chat ===")
	fmt.Println("Ваш peer ID:   ", h.ID())
	for _, a := range h.Addrs() {
		fmt.Printf("Слушаю: %s/p2p/%s\n", a, h.ID())
	}

	h.SetStreamHandler(chatProto, handleStream)

	notifee := &mdnsNotifee{h: h}
	mdnsService := mdns.NewMdnsService(h, "p2pchat-mdns", notifee)
	defer mdnsService.Close()
	_ = mdnsService.Start()
	// 4. Читаем команды с stdin
	stdin := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _ := stdin.ReadString('\n')
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "/conn"):
			target := strings.TrimSpace(strings.TrimPrefix(line, "/conn"))

			if target == "" {
				fmt.Println("Использование: /conn <peerID | port>")
				continue
			}

			// 1. Попытка по порту
			if pid, ok := peerByPort[target]; ok {
				pi := peer.AddrInfo{ID: pid, Addrs: h.Peerstore().Addrs(pid)}
				if err := connectAndChat(ctx, h, pi); err != nil {
					fmt.Println("Ошибка соединения:", err)
				}
				continue
			}

			// 2. Попытка по peerID
			if pi, ok := peerByID[target]; ok {
				if err := connectAndChat(ctx, h, pi); err != nil {
					fmt.Println("Ошибка соединения:", err)
				}
				continue
			}

			fmt.Println("Неизвестный порт или peer ID")
		case line == "/peers":
			fmt.Println("Обнаруженные узлы (из peerstore):")
			for _, pid := range h.Peerstore().Peers() {
				if pid == h.ID() {
					continue // себя не показываем
				}
				addrs := h.Peerstore().Addrs(pid)
				for _, addr := range addrs {
					fmt.Printf("  %s → %s/p2p/%s\n", pid.ShortString(), addr, pid)
				}
			}

		case strings.HasPrefix(line, "/quit"):
			fmt.Println("Выход.")
			return
		default:
			broadcast(h, []byte(line+"\n"))
		}
	}
}

// handleStream принимает входящий поток и читает данные построчно
func handleStream(s network.Stream) {
	activeStreams[s.Conn().RemotePeer()] = s
	fmt.Print("Соединение установлено с " + s.Conn().LocalPeer().String() + "\n" + "> ")
	//defer s.Close()
	r := bufio.NewReader(s)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			return
		}
		fmt.Printf("\n%s: %s", s.Conn().RemotePeer(), strings.TrimSpace(msg)+"\n"+"> ")
	}
}

func connectAndChat(ctx context.Context, h host.Host, pi peer.AddrInfo) error {
	if err := h.Connect(ctx, pi); err != nil {
		return err
	}
	s, err := h.NewStream(ctx, pi.ID, chatProto)
	if err != nil {
		return err
	}
	activeStreams[pi.ID] = s
	go handleStream(s)
	_, _ = s.Write([]byte("Успешное соединение с " + h.ID().String() + "\n" + "> "))
	return nil
}

func broadcast(h host.Host, data []byte) {
	for pid, stream := range activeStreams {
		_, err := stream.Write(data)
		if err != nil {
			fmt.Println("Ошибка отправки в", pid, ":", err)
		}
	}
}
