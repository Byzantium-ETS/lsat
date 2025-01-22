# L402

An implementation of the [L402](https://docs.lightning.engineering/the-lightning-network/l402) protocol.

## Overview

L402 is a protocol that leverages the capabilities of the Lightning Network for token minting and service authorization to enable the monetization of APIs through Bitcoin.

> [!NOTE]
> Additionally, it offers an implementation of the [phoenixd](https://phoenix.acinq.co/server) API for integration with a real Lightning node.

## Usage

Currently, this project does not offer standalone server or client implementations. However, it provides essential utilities and an example setup to get started.

### Example

The example available in the `./server/` and `./client/` directories demonstrates using a mocked Lightning node to issue and resolve challenges.

To get started, follow these instructions:

1. **Launch the Server**

   Open a terminal and run the following command to start the server:

   ```sh
   go run ./examples/server/server.go
   # 2024/06/09 09:21:20 Server launched at localhost:8080
   ```

2. **Mint a Token and Access the Service**

   In another terminal, run the following command to mint a token and access the service:

   ```sh
   go run ./examples/client/server.go
   # Requesting Token...
   # {"user_id":"...","caveats": "...","signature":"..."}
   # Sending Authorization Request...
   # ...
   ```

## Model

The following diagram illustrates the domain model for the L402 implementation:

![Domain Model](<out/MDD/Domain Model.png>)

## Authorization Flow

The authorization flow for L402 tokens is depicted in the following diagram:

![Authorization Flow](<out/Authorization/Authorization Flow.png>)

## Resources

For more information, refer to the following resources:

- [Lightning Engineering API Documentation](https://lightning.engineering/api-docs/api/lnd/)
- [L402 Protocol Documentation](https://docs.lightning.engineering/the-lightning-network/l402)
- [Multihop Payments Documentation](https://docs.lightning.engineering/the-lightning-network/multihop-payments)

## License

This project is licensed under the terms of the [MIT License](LICENSE).
