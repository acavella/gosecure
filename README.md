![GoShield Logo](docs/assets/goshield_logo_60.png)
---
![GitHub](https://img.shields.io/github/license/acavella/goshield)
![GitHub commit activity (branch)](https://img.shields.io/github/commit-activity/t/acavella/goshield)
![GitHub last commit (branch)](https://img.shields.io/github/last-commit/acavella/goshield/main)
![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/acavella/goshield/main)
![GitHub release (with filter)](https://img.shields.io/github/v/release/acavella/goshield)

## Description
A simple utility to easily encrypt and decrypt files written in Golang. Goshield utilizes the AES-256-GCM symetric encryption algorithm and Argon2id key derivation function to secure files. Files are encrypted and decrypted utilizing a user provided password.

## Instructions
### Installation
1. Navigate to the [latest release](https://github.com/acavella/goshield/releases/latest)
2. Download the binary appropriate for your operating system and architecture (Windows-AMD64 or Linux-AMD64)
3. Additionally, you can download the latest source and build an appropriate binary for your architecture
*Note: We are only able to support official builds*

## Process Diagrams
### Encryption
```shell
┌────────┐  ┌──────────┐  ┌─────────────────────────────────┐
│  SALT  │  │ Password │  │            Plaintext            │
└────┬───┘  └───────┬──┘  └────────────────┬────────────────┘
     │              │                      │
     │              │                      │
     ├──────────┐   │                      │
     │          │   │                      │
     │          ▼   ▼                      ▼
     │    ┌────────────────┐      ┌──────────────────┐
     │    │                │      │                  │
     │    │  Argon2id KDF  ├─────►│  GCM Encryption  │
     │    │                │      │                  │
     │    └────────────────┘      └────────┬─────────┘
     │                                     │
     │                                     │
     ▼                                     ▼
 ┌────────┬─────────────────────────────────────────────────┐
 │  SALT  │                     Ciphertext                  │
 ├────────┼─────────────────────────────────────────────────┤
 0        32                                               EOF
```
### Decryption
```shell
               0        32                                               EOF
 ┌──────────┐  ┌────────┬─────────────────────────────────────────────────┐
 │ Password │  │  SALT  │                     Ciphertext                  │
 └─────┬────┘  └───┬────┴──────────────────────────┬──────────────────────┘
       │           │                               │
       │           │                               │
       │           │                               │
       │           │                               │
       │           │                               │
       ▼           ▼                               ▼
     ┌────────────────┐                  ┌──────────────────┐
     │                │                  │                  │
     │  Argon2id KDF  ├─────────────────►│  GCM Encryption  │
     │                │                  │                  │
     └────────────────┘                  └─────────┬────────┘
                                                   │
                                                   │
                                                   ▼
                         ┌─────────────────────────────────────────────────┐
                         │                     Plaintext                   │
                         └─────────────────────────────────────────────────┘
```

## Credit
Initial development and work is based off of the hard work of the following folks:
- Eli Bendersky [https://eli.thegreenplace.net]
- Mert Kimyonşen [https://github.com/mrtkmynsndev/]
