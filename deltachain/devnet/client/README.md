# devnet deltachain client config

This folder contains config for running `deltachaind` against the devnet (Tilt) deltachain instance.

### examples

transfer `utest` from the account used by `deltachain-0` to the account used by `deltachain-1`, to smoke-test deltachain - make sure we can connect to the RPC port, the accounts exist, and deltachain is producing blocks.

    ./build/deltachaind --home build tx bank send delta1cyyzpxplxdzkeea7kwsydadg87357qna3zg3tq  delta1wqwywkce50mg6077huy4j9y8lt80943ks5udzr  1utest --from deltachain-0  --yes --broadcast-mode block --keyring-backend test
