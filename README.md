# qn-marketplace-cli

The `qn-marketplace-cli` is a command line interface for the [QuickNode Marketplace](https://www.quicknode.com/marketplace).

It provides:

- A set of commands to test an add-on's provisioning implementation
- A command to test your an add-on's SSO implementation
- A command to test your an add-on's RPC methods

The CLI is designed to allow you to test your QuickNode Marketplace add-ons directly on your localhost environment as your a developing them.

## Getting Started & Installation

To install `qn-marketplace-cli` to your machine, you can download a pre-built binary for your operating system from the [bin](./bin) directory of this repo.

You can also download this repository and build the code on your local machine. See the [Development](#development) section below for more information on how to do that.

## Commands

### Help

To get more information about how to use the CLI, use the following command:

```
qn-marketplace-cli help
```

Please note that if you are debugging and want to see more information, you can use the `--verbose` flag for every command.


### PUDD Testing

QuickNode Marketplace add-ons [must implement four provisioning API endpoints](https://www.quicknode.com/guides/quicknode-products/marketplace/how-provisioning-works-for-marketplace-partners/):

- **Provision (POST)**: called when a QuickNode customers installs the add-on on an endpoint.
- **Update (PUT)**: called when a previously provisioned endpoint gets updated.
- **Deactivate Endpoint (DELETE)**: called when a previously provisioned endpoint gets discarded.
- **Deprovision (DELETE)**: called when the add-on is uninstalled for a customer account (for all endpoints)

The `qn-marketplace-cli` has 4 different commands that allows you to test each one of these actions in isolation:

PROVISION:

```sh
qn-marketplace-cli provision --url=http://localhost:3000/provision --basic-auth=q24rqaergser --chain=ethereum --network=mainnet --plan=your-plan-slug --quicknode-id=abcdef --endpoint-id=foobar
```

UPDATE:

```sh
qn-marketplace-cli update --url=http://localhost:3000/update --basic-auth=q24rqaergser --chain=ethereum --network=mainnet --plan=your-plan-slug --quicknode-id=abcdef --endpoint-id=foobar
```

DEACTIVATE ENDPOINT:

```sh
qn-marketplace-cli deactivate --url=http://localhost:3000/deactivate_endpoint --basic-auth=q24rqaergser --endpoint-id=foobar
```

DEPROVISION:

```sh
qn-marketplace-cli deprovision --url=http://localhost:3000/deprovision --basic-auth=q24rqaergser --quicknode-id=abcdef
```

It also has one command that allows you to test all four actions at once:

```sh
qn-marketplace-cli pudd --base-url=http://localhost:3000/ --basic-auth=q24rqaergser --chain=ethereum --network=mainnet --plan=your-plan-slug
```

### JSON-RPC Testing

QuickNode Marketplace add-ons extends our capabilities by adding new JSON-RPC methods to QuickNode's existing endpoints.
Please read [this guide](https://www.quicknode.com/guides/quicknode-products/marketplace/how-to-create-an-rpc-add-on-for-marketplace/) for more information.

If your add-on has RPC methods, the `qn-marketplace-cli` allows you to test your implementation by making some JSON-RPC calls to your application.

```sh
qn-marketplace-cli rpc  --url=http://localhost:3000/rpc --method=your_addOnMethod --rpc-params='[9, "f"]' --chain=solana --network=mainnet
```

## Development

`qn-marketplace-cli` is developed using [Go](https://go.dev/) and [Cobra](https://github.com/spf13/cobra) and released under an [MIT License](./LICENSE.txt).

We welcome contributions to this repository to help us improve the CLI.

To fetch, build and install from the Github source:

```
go install github.com/quicknode-labs/qn-marketplace-cli@latest
```

To build the project:

```
go build -o bin/qn-marketplace-cli
```

Then:

```
bin/qn-marketplace-cli
```

## Testing on your localhost/development environment

You can use the following commands to test locally if you have an application that is hosted at http://localhost:3030.

For example, if you are using our [qn-js-add-on](https://github.com/quiknode-labs/qn-js-add-on) repo, which is a
sample add-on built with node.js, then you can use the following commands, assuming you have set the username and
password to `username` and `password` respectively.

```sh
go build
```

### Testing Provisioning

You can test the 4 provision actions:

```sh
./qn-marketplace-cli provision --url http://localhost:3030/provisioning/provision --chain ethereum --network mainnet --plan test --quicknode-id foobar --endpoint-id bazbaz --basic-auth dXNlcm5hbWU6cGFzc3dvcmQ=

./qn-marketplace-cli update --url http://localhost:3030/provisioning/update --chain ethereum --network mainnet --plan test --quicknode-id foobar --endpoint-id bazbaz --basic-auth dXNlcm5hbWU6cGFzc3dvcmQ=

./qn-marketplace-cli deactivate --url http://localhost:3030/provisioning/deactivate_endpoint  --quicknode-id foobar --endpoint-id bazbaz --chain ethereum --network mainnet --basic-auth dXNlcm5hbWU6cGFzc3dvcmQ=

./qn-marketplace-cli deprovision --url http://localhost:3030/provisioning/deprovision  --quicknode-id foobar --endpoint-id bazbaz --chain ethereum --network mainnet --basic-auth dXNlcm5hbWU6cGFzc3dvcmQ=
 ```

 Or you can test them all at once with PUDD:

 ```sh
 ./qn-marketplace-cli pudd --base-url http://localhost:3030/provisioning  --quicknode-id foobar --endpoint-id bazbaz --chain ethereum --network mainnet --basic-auth dXNlcm5hbWU6cGFzc3dvcmQ=
 ```

### Testing Single Sign On (SSO)

QuickNode Marketplace add-ons can provide a user-interface or dashboard that QuickNode customers can access from a link on their quicknode.com account. In order to seamlessly these customers from quicknode.com to Marketplace add-ons, an add-on can implement SSO. You can read [this guide](https://www.quicknode.com/guides/quicknode-products/marketplace/how-sso-works-for-marketplace-partners/) for more information on how SSO works with the QuickNode Marketplace.


Provide the provision url which should return a `dashboard-url` in the response. This will open a browser pointing to the dashboard url with the provisioned information.

 ```sh
 ./qn-marketplace-cli sso --url http://localhost:3030/provisioning/provision --email luc@example.com --name Luc --org QN --jwt-secret jwt-secret --basic-auth dXNlcm5hbWU6cGFzc3dvcmQ=
 ```


### Testing RPC calls

  ```sh
 ./qn-marketplace-cli rpc --url http://localhost:3030/provisioning/provision --rpc-url http://localhost:3030/rpc --rpc-method qn_fetchStuff --rpc-params "[\"abc\",123,\"zoo\"]" --basic-auth dXNlcm5hbWU6cGFzc3dvcmQ=
 ```

 ### Testing Healthcheck URL

  ```sh
 ./qn-marketplace-cli healthcheck --url http://localhost:3030/healthcheck
 ```
