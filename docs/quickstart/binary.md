<!--
order: 2
-->

# uptickd

`uptickd` is the all-in-one command-line interface. It supports wallet management, queries and transaction operations {synopsis}

## Pre-requisite Readings

- [Installation](./installation.md) {prereq}

## Build and Configuration

### Using `uptickd`

After you have obtained the latest `uptickd` binary, run:

```bash
uptickd [command]
```

Check the version you are running using

```bash
uptickd version
```

There is also a `-h`, `--help` command available

```bash
uptickd -h
```

::: tip
You can also enable auto-completion with the `uptickd completion` command. For example, at the start of a bash session, run `. <(uptickd completion)`, and all `uptickd` subcommands will be auto-completed.
:::

### Config and data directory

By default, your config and data are stored in the folder located at the `~/.uptickd` directory.

```bash
.                                   # ~/.uptickd
  ├── data/                           # Contains the databases used by the node.
  └── config/
      ├── app.toml                   # Application-related configuration file.
      ├── config.toml                # Tendermint-related configuration file.
      ├── genesis.json               # The genesis file.
      ├── node_key.json              # Private key to use for node authentication in the p2p protocol.
      └── priv_validator_key.json    # Private key to use as a validator in the consensus protocol.
```

To specify the `uptickd` config and data storage directory; you can update it using the global flag `--home <directory>`

### Configuring the Node

The Cosmos SDK automatically generates two configuration files inside `~/.uptickd/config`:

- `config.toml`: used to configure the Tendermint, learn more on [Tendermint's documentation](https://docs.tendermint.com/master/nodes/configuration.html),
- `app.toml`: generated by the Cosmos SDK, and used to configure your app, such as state pruning strategies, telemetry, gRPC and REST servers configuration, state sync, JSON-RPC, etc.

Both files are heavily commented, please refer to them directly to tweak your node.

One example config to tweak is the `minimum-gas-prices` field inside `app.toml`, which defines the minimum amount the validator node is willing to accept for processing a transaction. It is an anti spam mechanism and it will reject incoming transactions with less than the minimum gas prices.

If it's empty, make sure to edit the field with some value, for example `10token`, or else the node will halt on startup.

```toml
 # The minimum gas prices a validator is willing to accept for processing a
 # transaction. A transaction's fees must meet the minimum of any denomination
 # specified in this config (e.g. 0.25token1;0.0001token2).
 minimum-gas-prices = "0auptick"
```

### Pruning of State

There are four strategies for pruning state. These strategies apply only to state and do not apply to block storage.
To set pruning, adjust the `pruning` parameter in the `~/.uptickd/config/app.toml` file.
The following pruning state settings are available:

- `everything`: Prune all saved states other than the current state.
- `nothing`: Save all states and delete nothing.
- `default`: Save the last 100 states and the state of every 10,000th block.
- `custom`: Specify pruning settings with the `pruning-keep-recent`, `pruning-keep-every`, and `pruning-interval` parameters.

By default, every node is in `default` mode which is the recommended setting for most environments.
If you would like to change your nodes pruning strategy then you must do so when the node is initialized. Passing a flag when starting `uptick` will always override settings in the `app.toml` file, if you would like to change your node to the `everything` mode then you can pass the `---pruning everything` flag when you call `uptickd start`.

::: warning
**IMPORTANT**:
When you are pruning state you will not be able to query the heights that are not in your store.
:::

### Client configuration

We can view the default client config setting by using `uptickd config` command:

```bash
uptickd config
{
 "chain-id": "",
 "keyring-backend": "os",
 "output": "text",
 "node": "tcp://localhost:26657",
 "broadcast-mode": "sync"
}
```

We can make changes to the default settings upon our choices, so it allows users to set the configuration beforehand all at once, so it would be ready with the same config afterward.

For example, the chain identifier can be changed to `uptick_7000-1` from a blank name by using:

```bash
uptickd config "chain-id" uptick_7000-1
uptickd config
{
 "chain-id": "uptick_7000-1",
 "keyring-backend": "os",
 "output": "text",
 "node": "tcp://localhost:26657",
 "broadcast-mode": "sync"
}
```

Other values can be changed in the same way.

Alternatively, we can directly make the changes to the config values in one place at client.toml. It is under the path of `.uptick/config/client.toml` in the folder where we installed uptick:

```toml
############################################################################
### Client Configuration ###

############################################################################

# The network chain ID

chain-id = "uptick_7000-1"

# The keyring's backend, where the keys are stored (os|file|kwallet|pass|test|memory)

keyring-backend = "os"

# CLI output format (text|json)

output = "number"

# <host>:<port> to Tendermint RPC interface for this chain

node = "tcp://localhost:26657"

# Transaction broadcasting mode (sync|async|block)

broadcast-mode = "sync"
```

After the necessary changes are made in the `client.toml`, then save. For example, if we directly change the chain-id from `uptick_{{ $themeConfig.project.testnet_chain_id }}-2` to `upticktest_7000-1`, and output to number, it would change instantly as shown below.

```bash
uptickd config
{
 "chain-id": "upticktest_7000-1",
 "keyring-backend": "os",
 "output": "number",
 "node": "tcp://localhost:26657",
 "broadcast-mode": "sync"
}
```

### Options

A list of commonly used flags of `uptickd` is listed below:

| Option              | Description                   | Type         | Default Value   |
|---------------------|-------------------------------|--------------|-----------------|
| `--chain-id`        | Full Chain ID                 | String       | ---             |
| `--home`            | Directory for config and data | string       | `~/.uptickd`     |
| `--keyring-backend` | Select keyring's backend      | os/file/test | os              |
| `--output`          | Output format                 | string       | "text"          |

## Command list

A list of commonly used `uptickd` commands. You can obtain the full list by using the `uptickd -h` command.

| Command      | Description              | Subcommands (example)                                                     |
|--------------|--------------------------|---------------------------------------------------------------------------|
| `keys`       | Keys management          | `list`, `show`, `add`, `add  --recover`, `delete`                         |
| `tx`         | Transactions subcommands | `bank send`, `ibc-transfer transfer`, `distribution withdraw-all-rewards` |
| `query`      | Query subcommands        | `bank balance`, `staking validators`, `gov proposals`                     |
| `tendermint` | Tendermint subcommands   | `show-address`, `show-node-id`, `version`                                 |
| `config`     | Client configuration     |                                                                           |
| `init`       | Initialize full node     |                                                                           |
| `start`      | Run full node            |                                                                           |
| `version`    | Uptick version            |                                                                           |
