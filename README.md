# Whaletail

Automatically advertise selected docker networks as routes on your Tailscale
network.

## Download

- **Binary**: Download the latest release for your platfrom from the GitHub
  releases page: [here]().
- **Docker**: (linux only) Simply pull the docker image: `ghcr.io/willfantom/whalescale:latest`

## Usage

Both Docker and Tailscale are expected to be running on the host system in all
usage modes. Tailscale is also expected to be up and logged in.

- **Binary**: Simply run the `whaletail` application (assuming it is in your PATH)
- **Docker**: Create and run a container using the given image, binding your
  docker socket `/var/run/docker.sock` and your tailscale socket to
  `/var/run/tailscale/tailscaled.sock`.

Whaletail will then run, continuously trying to reconnect to docker or
tailscale if they are not found.

For whaletail to advertise the route to a docker network, the network must have the label `whaletail.enable=true`.

## Configuration: Tailscale Control

- Your nodes running whaletail will only be able to access other routes provided
  through your tailnet (including those advertised by other whaletails) if when
  running `tailscale up`, the flag `--accept-routes=true` was provided.

- Tailscale by default requires that any routes advertised by a node are
  approved via the Admin Console. This can make running a whaletail instance
  annoying... To bypass this you must:
  1.  Give your tailscale node a tag. For example, you could make sure all of
      your whaletail nodes get the tag `tag:whaletail`. This might look like
      this in your Tailscale ACL file.
      ```
      "tagOwners": {
		    "tag:whaletail": ["autogroup:admin"],
	    },
      ```
      Then to apply the tag to your node when running the `tailscale up` command
      with the flag: `--advertise-tags tag:whaletail`.

  2.  Modify your ACL configuration to auto-approve routes advertised by nodes
      with that tag. This might look like so:
      ```
      "autoApprovers": {
		    "routes": {
			    "0.0.0.0/0": ["tag:whaletail"],
		    },
	    },
     ```

## Configuration: Whaletail

Whaletail can be provided with a specific configuration via a file that should
be located at `/etc/whaletail/config.toml`. See the example config
[here](./config/config.toml) for more info. If using Docker, be sure to mount
this to the container if you make any changes from the default.

## Troubleshooting

- When creating docker networks to share on the tailnet, be considerate of IP
  address space overlaps and make sure to only use private address spaces. A
  good practice is to pick an address range with no known overlaps in or
  advertised by your tailnet (e.g. this could be something link
  `10.100.0.0/16`), and give each docker network a unique `/24` that fits inside
  this (e.g. this could be `10.100.100.0/24`).
