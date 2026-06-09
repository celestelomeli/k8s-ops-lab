# Phase 3 Observations

## Attempt 1: Manual Setup via AWS Console

The first pass at this was entirely manual. Launched two EC2 instances through the console, SSM'd into each one, and ran kubeadm commands by hand.

What that taught me:

How many moving parts kubeadm actually has: swap off, kernel modules, sysctl params, containerd config, then init and join. How easy it is to miss a step or do things out of order and end up with a broken node. Why the control plane and worker have different jobs. The control plane runs etcd, the API server, scheduler, and controller manager as static pods; the worker just runs kubelet and joins. What `kubectl get nodes` and `kubectl get pods -n kube-system` actually show when you built the cluster yourself.

The limitation: manual steps aren't repeatable. `terraform destroy` means starting from scratch with no record of what you did.

## Attempt 2: Terraform + Ansible

Switched to Terraform to provision the EC2 instances and Ansible to configure them. This made the setup reproducible and reviewable.

**Terraform decisions:**

Remote state in S3 with DynamoDB locking. No SSH key pair on the instances. SSM only for access. IAM instance profile with `AmazonSSMManagedInstanceCore` attached so SSM can reach the instances without open ports. An `availability_zone` variable with a default of `us-east-1a` to avoid hardcoding, but this not best practice. Best practice is a custom VPC with explicit subnets, but sufficient for this default VPC project. Dynamic inventory (`aws_ec2.yml`) instead of a static `inventory.ini` — instances get new IDs after every `terraform destroy` and `apply`, so hardcoding IDs would break every time.

**What Ansible taught me:**

Roles enforce a structure (tasks/, handlers/) that Ansible requires. The folder structure is what gives the ability to run different tasks on different machines." `site.yml` is like a table of contents. It delegates work to roles and controls which hosts each role runs on. The join command flows from the control plane play to the worker play through `hostvars`. Idempotency: tasks with `creates:` guards (like `kubeadm init`) can be re-run safely because they check whether the work is already done.

## SSM + Ansible: What Went Wrong

Choosing SSM as the Ansible connection method caused several problems.

**Dynamic inventory required S3.** The `community.aws.aws_ssm` connection plugin doesn't stream output back directly. It writes command output to an S3 bucket and reads it from there. The bucket name has to be set in the inventory, and the instances need an IAM policy allowing them to write to it. This wasn't obvious upfront.

**Privilege escalation doesn't work over SSM.** Ansible can run as root (`become: yes`) but can't switch to an unprivileged user mid-task (`become_user: ubuntu`) because that requires creating a temp file and handing it between users which SSM doesn't support. Worked around it by staying as root and pointing `KUBECONFIG` at `/etc/kubernetes/admin.conf` instead of the ubuntu user's copy. Not something I'd do with production workloads.

**macOS-specific crash.**  Ansible runs tasks by spinning up multiple background processes at once (one per host). macOS expects each background process to start clean with no inherited state from the parent. The SSM plugin inherits an already-open AWS connection from the parent process and tries to use it in the child, and macOS flags this as unsafe and kills the process. Two workarounds were required: OBJC_DISABLE_INITIALIZE_FORK_SAFETY=YES tells macOS to skip that check, and --forks 1 tells Ansible to run one host at a time instead of all at once. Neither of these would be needed on Linux.

**Better alternatives:**

SSH with a key pair is Ansible's native connection method. It has full support for privilege escalation, no S3 relay, and no macOS workaround needed. The tradeoff is managing key files and keeping port 22 restricted to a bastion or VPN. User data bootstrapping puts the setup script directly in Terraform's `user_data` block and runs it on first boot with no Ansible needed. It is simpler with fewer moving parts but you lose visibility because it runs silently in the background and you check `/var/log/cloud-init.log` to debug.

SSM was chosen here to avoid SSH keys and open ports. Next time, will use SSH with access restricted to a VPN or bastion.