---
title: "Welcome to k3s"
date: 2020-10-20T20:21:40+01:00
draft: false
---
## Intro

![k3s logo](https://cdn.shortpixel.ai/client/q_glossy,ret_img/https://www.worksonarm.com/wp-content/uploads/2020/06/Rancher_K3s2-750x422.png)

`k3s` runs faster than upstream kubeadm on devices like the [Raspberry Pi](https://www.raspberrypi.org)

There's also the [k3sup project](https://k3sup.dev) which you may like for installing Kubernetes over SSH.

```bash
k3sup install --ip $SERVER
k3sup join --ip $AGENT1 --server-ip $SERVER
```
