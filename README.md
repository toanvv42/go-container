# 🐳 go-container: Build Your Own Docker

A simplified container runtime written in Go for learning how Linux containers work under the hood. This project incrementally builds up from raw Linux primitives (namespaces, cgroups, chroot) to a minimal but functional container runtime.

## Why?

Docker and containers feel like magic — until you build one yourself. This project demystifies containers by implementing the core building blocks step by step:

- **Namespaces** for process isolation
- **Cgroups** for resource limits
- **OverlayFS** for layered filesystems
- **Veth/bridge** for container networking

By the end, you'll have a CLI tool that can pull a rootfs, run an isolated process with resource limits, and connect it to the network — a mini Docker.

## Prerequisites

- **Linux machine or VM** (namespaces/cgroups are Linux kernel features; won't work on macOS natively)
- **Root access** (`sudo`) for namespace/cgroup/network operations
- **Go 1.21+**

---

## Project Structure

```
go-container/
├── cmd/             # CLI entry points (run, ps, exec)
├── pkg/
│   ├── container/   # Core: namespaces, chroot, process lifecycle
│   ├── cgroup/      # Cgroup setup & teardown
│   ├── network/     # Veth, bridge, iptables
│   ├── image/       # Image download, extract, overlay
│   └── fs/          # Mount, pivot_root, overlayfs
├── rootfs/          # Downloaded root filesystems (gitignored)
├── go.mod
└── main.go
```

---

## Roadmap

### Phase 1 — Foundations (Run a command in isolation)

Implement the basic `run` command that executes a process inside Linux namespaces with its own root filesystem.

| Step | Description |
|------|-------------|
| 1.1 | Basic `run` command — execute a child process using `os/exec` with `syscall.SysProcAttr` |
| 1.2 | UTS namespace (`CLONE_NEWUTS`) — container gets its own hostname |
| 1.3 | PID namespace (`CLONE_NEWPID`) — container process sees itself as PID 1 |
| 1.4 | Mount namespace (`CLONE_NEWNS`) — private mount table |
| 1.5 | Chroot / `pivot_root` — give the container its own root filesystem (Alpine rootfs) |
| 1.6 | Mount `/proc` — so `ps` and `/proc` work correctly inside the container |

**✅ Success Criteria:**
- [ ] `sudo go-container run /bin/sh` drops you into a shell
- [ ] `hostname` inside the container shows a different hostname than the host
- [ ] `ps aux` inside the container only shows the container's own processes (PID 1 is the shell)
- [ ] The container's filesystem is isolated — changes inside don't affect the host
- [ ] `ls /` inside the container shows the Alpine rootfs, not the host root

---

### Phase 2 — Resource Limits (Cgroups)

Use cgroups v2 to constrain the container's resource usage.

| Step | Description |
|------|-------------|
| 2.1 | Memory limit — write to `/sys/fs/cgroup/.../memory.max` |
| 2.2 | CPU limit — write to `/sys/fs/cgroup/.../cpu.max` |
| 2.3 | PID limit — write to `/sys/fs/cgroup/.../pids.max` |
| 2.4 | Cleanup — remove the cgroup directory on container exit |

**✅ Success Criteria:**
- [ ] A container with `--memory 50m` gets OOM-killed when exceeding 50 MB
- [ ] A container with `--cpus 0.5` is throttled to 50% of one CPU core
- [ ] A container with `--pids 20` cannot fork more than 20 processes (fork bomb protection)
- [ ] After the container exits, its cgroup directory under `/sys/fs/cgroup/` is cleaned up

---

### Phase 3 — Filesystem & Images

Implement copy-on-write layering and basic image management.

| Step | Description |
|------|-------------|
| 3.1 | OverlayFS — layer a writable upper dir on top of a read-only rootfs |
| 3.2 | Image management — download/extract rootfs tarballs, store in a local image directory |
| 3.3 | `ps` / `list` command — track running containers via a state file |

**✅ Success Criteria:**
- [ ] `go-container image pull alpine` downloads and extracts an Alpine rootfs
- [ ] `go-container images` lists available local images
- [ ] Running a container uses OverlayFS — the base image stays unmodified after container exits
- [ ] Two containers can run from the same image simultaneously without conflicts
- [ ] `go-container ps` lists currently running containers with their IDs and commands

---

### Phase 4 — Networking

Give containers network connectivity using veth pairs and a bridge.

| Step | Description |
|------|-------------|
| 4.1 | Network namespace (`CLONE_NEWNET`) — container gets its own network stack |
| 4.2 | Veth pair — create a virtual ethernet pair, move one end into the container |
| 4.3 | Bridge — create a `go-container0` bridge, attach the host-side veth |
| 4.4 | IP assignment — assign IP addresses to bridge and container veth |
| 4.5 | NAT — set up iptables masquerade so containers can reach the internet |

**✅ Success Criteria:**
- [ ] Container has its own `eth0` interface with an assigned IP (e.g., `10.0.0.2/24`)
- [ ] Container can `ping` the host bridge IP (e.g., `10.0.0.1`)
- [ ] Container can reach the internet (`ping 8.8.8.8` and `ping google.com` if DNS is configured)
- [ ] Two containers can ping each other via the bridge
- [ ] Network resources (veth, bridge) are cleaned up on container exit

---

### Phase 5 — CLI & UX

Polish the user experience with a proper CLI and container lifecycle management.

| Step | Description |
|------|-------------|
| 5.1 | Cobra CLI — structured subcommands: `run`, `ps`, `exec`, `images` |
| 5.2 | Logging — capture and display container stdout/stderr |
| 5.3 | `exec` command — enter a running container's namespaces using `setns` |
| 5.4 | Flags — `--name`, `--memory`, `--cpus`, `--pids`, `--detach` |

**✅ Success Criteria:**
- [ ] `go-container run --name mybox --memory 100m alpine /bin/sh` runs a named container with limits
- [ ] `go-container exec mybox /bin/sh` attaches to a running container
- [ ] `go-container ps` shows running containers with name, PID, image, and status
- [ ] `go-container run --detach alpine sleep 3600` runs a container in the background
- [ ] Container logs are retrievable after the container exits

---

## 📚 Recommended Reading

Read these before (or alongside) implementation. Ordered from essential to supplementary.

### Core Concepts

| Resource | What You'll Learn |
|----------|-------------------|
| [Linux Namespaces (man7.org)](https://man7.org/linux/man-pages/man7/namespaces.7.html) | All 8 namespace types, how `clone()`, `unshare()`, `setns()` work |
| [Cgroups v2 (kernel.org)](https://docs.kernel.org/admin-guide/cgroup-v2.html) | Cgroup hierarchy, controllers (memory, cpu, pids), delegation |
| [pivot_root(2)](https://man7.org/linux/man-pages/man2/pivot_root.2.html) | How to safely swap the root filesystem (preferred over chroot) |
| [OverlayFS (kernel.org)](https://docs.kernel.org/filesystems/overlayfs.html) | Union mount: lower/upper/work/merged dirs, copy-up behavior |

### Hands-On Guides

| Resource | What You'll Learn |
|----------|-------------------|
| [Containers From Scratch — Liz Rice (YouTube)](https://www.youtube.com/watch?v=8fi7uSYlOdc) | ~35 min live-coding a container in Go — closest thing to this project |
| [Containers From Scratch — Liz Rice (GitHub)](https://github.com/lizrice/containers-from-scratch) | Source code for the talk above |
| [Build Your Own Container Using Less than 100 Lines of Go](https://www.infoq.com/articles/build-a-container-golang/) | Step-by-step walkthrough of namespaces + chroot in Go |
| [Linux Containers in 500 Lines of Code](https://blog.lizzie.io/linux-containers-in-500-loc.html) | C implementation — great for understanding the syscalls directly |

### Networking

| Resource | What You'll Learn |
|----------|-------------------|
| [Linux Network Namespaces (man7.org)](https://man7.org/linux/man-pages/man7/network_namespaces.7.html) | How network namespaces isolate the network stack |
| [Container Networking From Scratch](https://labs.iximiuz.com/tutorials/container-networking-from-scratch) | Veth pairs, bridges, iptables NAT — exactly what Phase 4 needs |
| [Introduction to Linux interfaces for virtual networking](https://developers.redhat.com/blog/2018/10/22/introduction-to-linux-interfaces-for-virtual-networking) | Bridge, veth, macvlan, and other virtual networking primitives |

### Go-Specific

| Resource | What You'll Learn |
|----------|-------------------|
| [Go `syscall` package](https://pkg.go.dev/syscall) | `SysProcAttr`, `CLONE_*` flags, `Mount`, `PivotRoot` |
| [Go `os/exec` package](https://pkg.go.dev/os/exec) | Running child processes, setting up stdin/stdout/stderr |
| [`github.com/spf13/cobra`](https://github.com/spf13/cobra) | CLI framework used in Docker, Kubernetes, and this project |

### Deep Dives (Optional)

| Resource | What You'll Learn |
|----------|-------------------|
| [OCI Runtime Spec](https://github.com/opencontainers/runtime-spec) | The standard that runc (Docker's runtime) implements |
| [runc source code](https://github.com/opencontainers/runc) | Production container runtime in Go — reference implementation |
| [Docker source code — containerd](https://github.com/containerd/containerd) | How Docker manages containers above the runtime layer |

---

## License

MIT
