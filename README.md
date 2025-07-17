# P2PChat

P2PChat is a decentralized, peer-to-peer chat application built with Go, designed for secure and efficient real-time communication without relying on central servers. Leveraging Go's concurrency model and a robust P2P networking library, P2PChat enables users to connect directly, exchange messages, and enjoy a seamless chat experience.

## Features

- **Decentralized Communication**: Connect directly with peers using a P2P network, no central server required.
- **Real-Time Messaging**: Send and receive text messages instantly with low latency.
- **Secure Connections**: Messages are encrypted to ensure privacy and security.
- **Lightweight and Fast**: Built with Go for high performance and minimal resource usage.
- **Cross-Platform**: Runs on Windows, macOS, and Linux.
- **Customizable**: Easily extendable with additional features like file sharing or group chats.

## Realeses

- Download realese file <a href="https://github.com/abdulatif-abdumannopov/p2p_chat/releases/tag/v1.0.0">here</a>

## Prerequisites

- Go 1.18 or higher
- A working internet connection
- (Optional) Git for cloning the repository

## Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/username/p2pchat.git
   cd p2pchat
   ```

2. **Install Dependencies**:

   ```bash
   go mod tidy
   ```

3. **Build the Application**:

   ```bash
   go build -o p2pchat
   ```

## Usage

1. **Run the Application**:

   ```bash
   ./p2pchat
   ```

2. **Connect to a Peer**:

   - On startup, the app generates a unique peer ID and displays it.

   - Share your peer ID with another user, or enter their peer ID to connect.

   - Example:

     ```bash
       /conn <peer-id>
     ```

3. **Start Chatting**:

   - Once connected, type your messages and press Enter to send.
   - Type `/quit` to exit the application.

Example session:

```bash
$ ./p2pchat
Your Peer ID: 12D3KooW...
/conn 12D3KooW...
Connected to peer!
> Hello, friend!
< Hey, great to connect!
> /quit
```

## Contributing

We welcome contributions! To get started:

1. Fork the repository.
2. Create a new branch: `git checkout -b feature-name`.
3. Make your changes and commit: `git commit -m "Add feature"`.
4. Push to your fork: `git push origin feature-name`.
5. Open a pull request.

Please ensure your code follows the Go Code Review Comments and includes tests where applicable.

## License

This project is licensed under the MIT License. See the <a href="[github.com/abdulatif-abdumannopov/p2p_chat/blob/master/LICENSE.md](https://github.com/abdulatif-abdumannopov/p2p_chat/blob/master/LICENSE.md)">LICENCE</a> file for details.

## Acknowledgments

- Built with libp2p for P2P networking.
- Inspired by the Go community's focus on simplicity and performance.

Happy chatting! 🚀
