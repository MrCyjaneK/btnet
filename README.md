# BTnet

[![Build Status](https://ci.mrcyjanek.net/badge/build-btnet.svg)](https://ci.mrcyjanek.net/jobs/build-btnet)

BTnet - implementation of peer to peer www network using BitTorrent network.

 - No central server
 - Using BitTorrent network as it is (you can seed your site from your client)
 - Works in any browser
 - Simple and easy to use and implement

NOTE: main development happens at [git.mrcyjanek.net](https://git.mrcyjanek.net/mrcyjanek/btnet).

### TODO

 - Updates (fetch newer version signed by _btnet/pgp.asc)
 - Minimal api to allow dynamic content.


# Setup

 - First you need to download binary, or compile it yourself (github's releases tab have it compiled.)
    - on debian you can run:
    
    Install my APT repo to your system. Make sure to run this command as root
    ```bash
    # wget 'https://static.mrcyjanek.net/laminarci/apt-repository/cyjan_repo/mrcyjanek-repo-latest.deb' && \
      apt install ./mrcyjanek-repo-latest.deb && \
      rm ./mrcyjanek-repo-latest.deb && \
      apt update
    ```
    After that install btnet
    ```bash
    # apt install btnet
    ```
 - Run the binary
 - Open any web browser, go to settings and use `127.0.0.1:8080` as HTTP-only proxy.
 - open http://6626ae3c23b19bf4ba7d17c765be2c83935d51a3.btnet/ in your browser
 - profit.