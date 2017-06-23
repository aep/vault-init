package main

import (
    vaultApi   "github.com/hashicorp/vault/api"
    consulApi  "github.com/hashicorp/consul/api"
    "fmt"
    "net"
)



func main() {
    consul, err := consulApi.NewClient(consulApi.DefaultConfig())
    if err != nil {
        panic(err)
    }

    css, _, err := consul.Catalog().Service("vault", "", nil)
    if err != nil {
        panic(err)
    }

    fmt.Println("going to init these vault instances:")
    for _,cs := range css {
        fmt.Printf("  %v\n", cs.Address)
    }

    var re *vaultApi.InitResponse

    for i, cs := range css {
        config := vaultApi.DefaultConfig()
        config.Address = "http://" + cs.Address + ":8200"
        vault, err  := vaultApi.NewClient(config)
        if err != nil {
            panic(err)
        }

        if i == 0 {
            re, err = vault.Sys().Init(&vaultApi.InitRequest{
                SecretShares: 1,
                SecretThreshold: 1,

            })
            if err != nil {
                panic(err)
            }
            fmt.Printf("initialize on %v\nroot token: %v\nunseal keys:\n", cs.Address, re.RootToken)

            for _, k := range re.KeysB64 {
                fmt.Printf("    %v\n", k)
            }
        }
        fmt.Printf("unsealing %v\n", cs.Address)
        for _, k := range re.KeysB64 {
            _, err = vault.Sys().Unseal(k)
            if err != nil {
                panic(err)
            }
        }

        //continue everything waiting for a vault key
        //a half-automated setup process can wait on each node with VAULT_TOKEN=$(netcat -ul 8200)

        //TODO: we should send some other key, not the root key i guess
        ServerAddr,err := net.ResolveUDPAddr("udp", cs.Address + ":8200")
        if err != nil {
            panic(err)
        }

        conn, err := net.DialUDP("udp", nil, ServerAddr)
        if err != nil {
            panic(err)
        }
        defer conn.Close()

        _,err = conn.Write([]byte(re.RootToken))
        if err != nil {
            panic(err)
        }
    }
}

