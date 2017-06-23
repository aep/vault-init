package main

import (
    vaultApi   "github.com/hashicorp/vault/api"
    consulApi  "github.com/hashicorp/consul/api"
    "fmt"
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
    }
}

