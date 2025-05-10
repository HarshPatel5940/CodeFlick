<div align="center">
  <h1>CodeFlick</h1>
  <p><em>Flick it, share it! Open-source Gists for the dev community.</em></p>

  [![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)
</div>

## 🚀 Overview

CodeFlick is a modern, open-source platform designed for developers who want to share, discover, and collaborate on code snippets. Built as an alternative to GitHub Gists, our platform emphasizes speed, simplicity, and social coding.

## ✨ Features

- **Instant Sharing**: Share code snippets with just one click
- **Syntax Highlighting**: Support for 100+ programming languages
- **Version Control**: Track changes to your snippets over time
- **Collaboration**: Comment, fork, and improve others' snippets
- **Embeddable**: Easily embed your snippets in blogs and documentation
- **Privacy Controls**: Public, private, and password-protected snippets
- **Markdown Support**: Rich text formatting for better documentation

## 🛠️ Tech Stack

- **Backend**: Go with Echo framework and PostgreSQL database
- **Frontend**: Nuxt.js with Vue 3 and Nuxt UI
- **Storage**: MinIO for object storage
- **Authentication**: OAuth with multiple providers

## 🚀 Installation

### Prerequisites

- Go (v1.24+)
- Node.js (v16+)
- PostgreSQL
- MinIO (or compatible S3 storage)
- pnpm or npm

### Backend Setup

```bash
# Navigate to the server directory
cd CodeFlick/server

# Copy environment variables
cp .env.example .env

# Install Go dependencies
go mod download

# Run database migrations
make db-migrate

# Start the server
make run
```

### Frontend Setup

```bash
# Navigate to the client directory
cd CodeFlick/client

# Copy environment variables
cp .env.example .env

# Install dependencies
pnpm install

# Start the development server
pnpm dev
```

## 🔍 Usage

1. Create an account or sign in via OAuth
2. Click "New Flick" to create a code snippet
3. Add code, select language, and add description
4. Click "Flick It" to share
5. Copy the unique URL to share with others

## 💡 Contributing

We welcome contributions from developers of all skill levels!

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add some amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📜 License

This project is licensed under the Mozilla Public License 2.0 (MPL-2.0) - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- All our amazing contributors
- The Go and Nuxt.js communities
- The open-source community for inspiration

---

<div align="center">
    <div>
      Made with ❤️ by <a href="https://harshnpatel.in">harshnpatel.in</a>
    </div>
</div>
