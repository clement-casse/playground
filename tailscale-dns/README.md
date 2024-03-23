# Tailscale Custom DNS Instance Deployed in Fly.io

A very clunky attempt to ue [Fly.io](https://fly.io/) to host a [Tailscale](https://tailscale.com/) + [Blocky](https://0xerr0r.github.io/blocky/) instance to serve as DNS part of my Tailscale Network.
The initial motivation was running a Pi-Hole part of my tailnet instead of my NextDNS account, Tailscale already provides a good documentation how to do it [[3]].
However I do not intend to leverage all of the Pi-Hole functions like DHCP, and Blocky was more appealing as a solution because it is a lightweight process.
Despite the advice in [[4]] to not run a DNS Server on the Cloud some peoples made an attempt of leveraging Fly.io to host a Blocky service [[1]].

Here I try to boostrap Blocky with Tailscale to use Fly.io to have my always-on DNS server in the Cloud while doing my best to keep it closed to the outside world.
Instead of writing a entrypoint.sh script like in [[1], [2]], I used Go to include both the tailscale client ant the blocky binary part of a single binary.
Tailscale daemon is installed seprately but also part of the container.

## References

1. [`devusb/blocky-tailscale` GitHub Repo][1]
2. [Tailscale Knowledge Base entry for adding a Fly.io app to a tailnet][2]
3. [Tailscale Knowledge Base entry for adding containers to a tailnet][3]
4. [Tailscale Knowledge Base entry for adding a pi-hole to a tailnet][4]
5. [riesinger blog post on his Highly Available DNS setup with Blocky][5]

[1]: https://github.com/devusb/blocky-tailscale
[2]: https://tailscale.com/kb/1132/flydotio
[3]: https://tailscale.com/kb/1282/docker
[4]: https://tailscale.com/kb/1114/pi-hole
[5]: https://riesinger.dev/posts/ha-dns-adblocking-with-blocky/
