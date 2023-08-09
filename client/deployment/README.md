## Deployment Configuration

The folder structure shown below is auto-generated on running deployment via truffle using the `migrate` command.

```
- development
    - deployment.go (generated using cmd: pnpm run deploy_development)
- sapphirelocalnet
    - deployment.go (generated using cmd: pnpm run deploy_sapphire_localnet)
- sapphiremainnet
    - deployment.go (generated using cmd: pnpm run deploy_sapphire_mainnet)
- sapphiretestnet
    - deployment.go (generated using cmd: pnpm run deploy_sapphire_testnet)
```

With the deployment log file stored by default in the build folder, any of those files can be regenerated e.g using commad:

```
$ go run ./deployInfo ./build/deploy_sapphire_testnet.log
2023/08/09 06:29:40 Log deployment information has me written to : client/deployment/sapphiretestnet/deployment.go
```

Therefore no need to manually edit the files.
