<!--
order: 5
-->

# Interacting with the Node

There are multiple ways to interact with a node: using the CLI, using gRPC or using the REST endpoints. {synopsis}

## Using the CLI

Now that your very own node is running, it is time to try sending tokens from the first account you created to a second account. In a new terminal window, start by running the following query command:

```bash
uptickd query bank balances $MY_VALIDATOR_ADDRESS --chain-id=uptick_7000-1
```

You should see the current balance of the account you created, equal to the original balance of tokens you granted it minus the amount you delegated via the `gentx`. Now, create a second account:

```bash
uptickd keys add recipient --keyring-backend=file

# Put the generated address in a variable for later use.
RECIPIENT=$(uptickd keys show recipient -a --keyring-backend=file)
```

The command above creates a local key-pair that is not yet registered on the chain. An account is created the first time it receives tokens from another account. Now, run the following command to send tokens to the `recipient` account:

```bash
uptickd tx bank send $MY_VALIDATOR_ADDRESS $RECIPIENT 1000000auptick --chain-id=uptick_7000-1 --keyring-backend=file

# Check that the recipient account did receive the tokens.
uptickd query bank balances $RECIPIENT --chain-id=uptick_7000-1
```

Finally, delegate some of the stake tokens sent to the `recipient` account to the validator:

```bash
uptickd tx staking delegate $(uptickd keys show my_validator --bech val -a --keyring-backend=file) 500auptick --from=recipient --chain-id=uptick_7000-1 --keyring-backend=file

# Query the total delegations to `validator`.
uptickd query staking delegations-to $(uptickd keys show my_validator --bech val -a --keyring-backend=file) --chain-id=uptick_7000-1
```

You should see two delegations, the first one made from the `gentx`, and the second one you just performed from the `recipient` account.

## Using gRPC

The Protobuf ecosystem developed tools for different use cases, including code-generation from `*.proto` files into various languages. These tools allow the building of clients easily. Often, the client connection (i.e. the transport) can be plugged and replaced very easily. Let's explore one of the most popular transport: gRPC.

Since the code generation library largely depends on your own tech stack, we will only present three alternatives:

- `grpcurl` for generic debugging and testing
- programmatically via Go
- CosmJS for JavaScript/TypeScript developers

### grpcurl

[grpcurl](https://github.com/fullstorydev/grpcurl) is like `curl` but for gRPC. It is also available as a Go library, but we will use it only as a CLI command for debugging and testing purposes. Follow the instructions in the previous link to install it.

Assuming you have a local node running (either a localnet, or connected a live network), you should be able to run the following command to list the Protobuf services available (you can replace `localhost:9000` by the gRPC server endpoint of another node, which is configured under the `grpc.address` field inside [`app.toml`](/run-node.md#configuring-the-node-using-apptoml)):

```bash
grpcurl -plaintext localhost:9090 list
```

You should see a list of gRPC services, like `cosmos.bank.v1beta1.Query`. This is called reflection, which is a Protobuf endpoint returning a description of all available endpoints. Each of these represents a different Protobuf service, and each service exposes multiple RPC methods you can query against.

In order to get a description of the service you can run the following command:

```bash
# Service we want to inspect
grpcurl \
    localhost:9090 \
    describe cosmos.bank.v1beta1.Query                  
```

It's also possible to execute an RPC call to query the node for information:

```bash
grpcurl \
    -plaintext
    -d '{"address":"$MY_VALIDATOR"}' \
    localhost:9090 \
    cosmos.bank.v1beta1.Query/AllBalances
```

#### Query for historical state using grpcurl

You may also query for historical data by passing some [gRPC metadata](https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md) to the query: the `x-cosmos-block-height` metadata should contain the block to query. Using grpcurl as above, the command looks like:

```bash
grpcurl \
    -plaintext \
    -H "x-cosmos-block-height: 279256" \
    -d '{"address":"$MY_VALIDATOR"}' \
    localhost:9090 \
    cosmos.bank.v1beta1.Query/AllBalances
```

Assuming the state at that block has not yet been pruned by the node, this query should return a non-empty response.

### Programmatically via Go

The following snippet shows how to query the state using gRPC inside a Go program. The idea is to create a gRPC connection, and use the Protobuf-generated client code to query the gRPC server.

```go
import (
    "context"
    "fmt"

  "google.golang.org/grpc"

    sdk "github.com/cosmos/cosmos-sdk/types"
  "github.com/cosmos/cosmos-sdk/types/tx"
)

func queryState() error {
    myAddress, err := sdk.AccAddressFromBech32("eth...")
    if err != nil {
        return err
    }

    // Create a connection to the gRPC server.
    grpcConn := grpc.Dial(
        "127.0.0.1:9090", // your gRPC server address.
        grpc.WithInsecure(), // The SDK doesn't support any transport security mechanism.
    )
    defer grpcConn.Close()

    // This creates a gRPC client to query the x/bank service.
    bankClient := banktypes.NewQueryClient(grpcConn)
    bankRes, err := bankClient.Balance(
        context.Background(),
        &banktypes.QueryBalanceRequest{Address: myAddress, Denom: "eth"},
    )
    if err != nil {
        return err
    }

    fmt.Println(bankRes.GetBalance()) // Prints the account balance

    return nil
}
```

#### Query for historical state using Go

Querying for historical blocks is done by adding the block height metadata in the gRPC request.

```go
import (
  "context"
  "fmt"

  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"

  grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
  "github.com/cosmos/cosmos-sdk/types/tx"
)

func queryState() error {
    // --snip--

    var header metadata.MD
    bankRes, err = bankClient.Balance(
        metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, "12"), // Add metadata to request
        &banktypes.QueryBalanceRequest{Address: myAddress, Denom: denom},
        grpc.Header(&header), // Retrieve header from response
    )
    if err != nil {
        return err
    }
    blockHeight = header.Get(grpctypes.GRPCBlockHeightHeader)

    fmt.Println(blockHeight) // Prints the block height (12)

    return nil
}
```

### CosmJS

CosmJS documentation can be found at [https://cosmos.github.io/cosmjs](https://cosmos.github.io/cosmjs). As of January 2021, CosmJS documentation is still work in progress.

## Using the REST Endpoints

All gRPC services on the Cosmos SDK are made available for more convenient REST-based queries through gRPC-gateway. The format of the URL path is based on the Protobuf service method's full-qualified name, but may contain small customizations so that final URLs look more idiomatic. For example, the REST endpoint for the `cosmos.bank.v1beta1.Query/AllBalances` method is `GET /cosmos/bank/v1beta1/balances/{address}`. Request arguments are passed as query parameters.

As a concrete example, the `curl` command to make balances request is:

```bash
curl \
    -X GET \
    -H "Content-Type: application/json" \
    http://localhost:1317/cosmos/bank/v1beta1/balances/$MY_VALIDATOR
```

Make sure to replace `localhost:1317` with the REST endpoint of your node, configured under the `api.address` field.

The list of all available REST endpoints is available as a Swagger specification file, it can be viewed at `localhost:1317/swagger`. Make sure that the `api.swagger` field is set to true in your [`app.toml`](/run-node.md#configuring-the-node-using-apptoml) file.

### Query for historical state using REST

Querying for historical state is done using the HTTP header `x-cosmos-block-height`. For example, a curl command would look like:

```bash
curl \
    -X GET \
    -H "Content-Type: application/json" \
    -H "x-cosmos-block-height: 279256"
    http://localhost:1317/cosmos/bank/v1beta1/balances/$MY_VALIDATOR
```

Assuming the state at that block has not yet been pruned by the node, this query should return a non-empty response.

### Cross-Origin Resource Sharing (CORS)

[CORS policies](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) are not enabled by default to help with security. If you would like to use the rest-server in a public environment we recommend you provide a reverse proxy, this can be done with [nginx](https://www.nginx.com/). For testing and development purposes there is an `enabled-unsafe-cors` field inside [`app.toml`](/run-node.md#configuring-the-node-using-apptoml).
