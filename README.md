this is a hack. don't use this before reading the code!

assuming you have all vault instances registered to consul, run on any node:

    docker run aaep/vault-init

this will act similar to vault init, except it will unseal every node in the cluster too
