![GoShield Logo](docs/assets/goshield_logo.png)

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://user-images.githubusercontent.com/71297412/178180562-38f53e67-a31f-4c9f-b1a1-c9221703df4b.png">
  <source media="(prefers-color-scheme: light)" srcset="https://user-images.githubusercontent.com/71297412/178180441-59f1644e-2ab6-4bf0-866f-2c77b2a63433.png">
  <img alt="Hashnode logo" src="https://user-images.githubusercontent.com/71297412/178180441-59f1644e-2ab6-4bf0-866f-2c77b2a63433.png" height="25">
</picture>

---
![GitHub](https://img.shields.io/github/license/acavella/goshield)
![GitHub commit activity (branch)](https://img.shields.io/github/commit-activity/t/acavella/goshield)
![GitHub last commit (branch)](https://img.shields.io/github/last-commit/acavella/goshield/main)
![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/acavella/goshield/main)
![GitHub release (with filter)](https://img.shields.io/github/v/release/acavella/goshield)

## Description
A simple utility to easily encrypt and decrypt files written in Golang. Goshield utilizes the AES-256-GCM symetric encryption algorithm and Argon2id key derivation function to secure files. Files are encrypted and decrypted utilizing a user provided password.

## Process Diagram

```shell
┌────────┐  ┌──────────┐  ┌─────────────────────────────────┐
│  SALT  │  │ Password │  │            Plaintext            │
└────┬───┘  └───────┬──┘  └────────────────┬────────────────┘
     │              │                      │
     │              │                      │
     │              │                      │
     ├──────────┐   │                      │
     │          │   │                      │
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
 │        │                                                 │
 0        32                                               EOF
```

## Credit
Initial script is based on the hard work of Eli Bendersky [https://eli.thegreenplace.net]
