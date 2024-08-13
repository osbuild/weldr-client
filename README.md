# composer-cli

`composer-cli` is a command line utility used with
[osbuild-composer](https://www.osbuild.org) to manage blueprints, build and
upload images, and manage source repositories.

* [Edit a Blueprint](#edit-a-blueprint)
* [Build an image](#build-an-image)
* [Monitor the build status](#monitor-the-build-status)
* [Download the image](#download-the-image)
* [Image Uploads](#image-uploads)
* [Build an image and upload results](#build-an-image-and-upload-results)
* [JSON Output](#json-output)
* [Blueprint Format](#blueprint-format)
* [Package Sources](#package-sources)

## Edit a Blueprint

Start out by listing the available blueprints using `composer-cli blueprints
list`, pick one and save it to the local directory by running `composer-cli
blueprints save http-server`. If there are no blueprints available you can copy
one of the examples from [here][examples].

Edit the file (it will be saved with a .toml extension) and change the
description, add a package to it. Send it back to the server by
running `composer-cli blueprints push http-server.toml`. You can verify that it was
saved by viewing the changelog - `composer-cli blueprints changes http-server`.

See the [Blueprint Format](#blueprint-format) section for the details on how to
create a blueprint.


## Build an image

Build a `qcow2` disk image from this blueprint by running `composer-cli
compose start http-server qcow2`. It will print a UUID that you can use to
keep track of the build. You can also cancel the build if needed.

The available types of images is displayed by `composer-cli compose types`.
Currently this consists of: ami, fedora-iot-commit, openstack, qcow2, vhd, vmdk

You can optionally start an upload of the finished image, see
[Image Uploads](#image-uploads) for more information.


## Monitor the build status

Monitor it using `composer-cli compose info UUID` where UUID is the UUID
returned by the start command. This will show the status of the build. You can
view the build logs once it is in the `RUNNING` state using `composer-cli
compose log UUID`

Once the build is in the `FINISHED` state you can download the image.


## Download the image

Downloading the final image is done with `composer-cli compose image UUID` and
it will save the qcow2 image as `UUID-disk.qcow2` which you can then use to
boot a VM like this:

```
qemu-kvm --name test-image -m 1024 -hda ./UUID-disk.qcow2
```

## Image Uploads

`composer-cli` can upload the images to a number of services, including AWS,
Azure, and VMWare. The upload can be started when the build is finished by
passing the service's profile.toml to the `compose start` command. For example:

```
composer-cli compose start http-server server aws.toml
```

### providers.toml

Each provider requires it's own set of details in order to upload to it, this
usually involves authentication information. Each provider has its own list of
requirements.

#### AWS

A `provider.toml` file for AWS looks like this:

```
[settings]
bucket = "AWS Bucket"
region = "AWS Region"
key = "AWS Key"
accessKeyID = "AWS Access Key"
secretAccessKey = "AWS Secret Key"
```

The access key and secret key can be created by going to the
`IAM->Users->Security Credentials` section and creating a new access key. The
secret key will only be shown when it is first created so make sure to record
it in a secure place. The region should be the region that you want to use the
AMI in, and the bucket can be an existing bucket, or a new one, following the
normal AWS bucket naming rules. It will be created if it doesn't already exist.

When uploading the image it is first uploaded to the s3 bucket, and then
converted to an AMI.  If the conversion is successful the s3 object will be
deleted. If it fails, re-trying after correcting the problem will re-use the
object if you have not deleted it in the meantime, speeding up the process.

#### Azure

For Azure the `provider.toml` looks like:

```
[settings]
storageAccount = "account"
storageAccessKey = "key"
container = "container"
```

#### VMWare

The VMWare `provider.toml` uses this template:

```
[settings]
host =  "Hostname"
username =  "Username"
password = "Password"
datacenter = "Datacenter"
cluster = "Cluster"
datastore = "Datastore"
```

## Build an image and upload results

If you have a profile named `test-uploads`:

```
composer-cli compose start example-http-server ami "http image" aws test-uploads
```

Or if you have the settings stored in a TOML file:

```
composer-cli compose start example-http-server ami "http image" aws-settings.toml
```

It will return the UUID of the image build, and the UUID of the upload. Once
the build has finished successfully it will start the upload process, which you
can monitor with `composer-cli upload info <UPLOAD-UUID>`

You can also view the upload logs from the Ansible playbook with:

```
composer-cli upload log <UPLOAD-UUID>
```

The type of the image must match the type supported by the provider.


## JSON Output

`composer-cli` can output the JSON data returned by the `osbuild-composer` API,
either for debugging or testing purposes. The return format is a JSON 'object'
that contains 4 fields: `method` with the HTTP method used to make the request,
`path` is the API path that was called, `status` is the HTTP return code from
the server, and `body` contains the raw JSON returned by the server.

Some commands send 2 requests to the server in order to retrieve all the
results at once.  The API supports pagination and defaults to 20 items, so you
need to find the total and set the limit to that total in order to get all of
them.

For example, the JSON response from a `composer-cli blueprints list` looks like
this:

```
[{
    "method": "GET",
    "path": "/blueprints/list?limit=0",
    "status": 200,
    "body": {
        "blueprints": [],
        "limit": 0,
        "offset": 0,
        "total": 3
    }
},
{
    "method": "GET",
    "path": "/blueprints/list?limit=23",
    "status": 200,
    "body": {
        "blueprints": [
            "http-server-bp-1"
            "database-bp-1"
            "dev-bp-1"
        ],
        "limit": 3,
        "offset": 0,
        "total": 3
    }
}]
```

NOTE: This output format changed in weldr-client v35.6, it used to be a stream of objectes and is
now a proper JSON list of objects, making it easier to parse.


## Blueprint Format

Blueprints are simple text files in [TOML](https://github.com/toml-lang/toml) format that describe
which packages, and what versions, to install into the image. They can also define a limited set
of customizations to make to the final image.

Example blueprints can be found in [here][examples], with a simple one looking like this:

```
name = "base"
description = "A base system with bash"
version = "0.0.1"

[[packages]]
name = "bash"
version = "4.4.*"
```

The `name` field is the name of the blueprint. It can contain spaces, but they will be converted to `-`
when it is written to disk. It should be short and descriptive.

`description` can be a longer description of the blueprint, it is only used for display purposes.

`version` is a [semver compatible](https://semver.org/) version number. If
a new blueprint is uploaded with the same `version` the server will
automatically bump the PATCH level of the `version`. If the `version`
doesn't match it will be used as is. eg. Uploading a blueprint with `version`
set to `0.1.0` when the existing blueprint `version` is `0.0.1` will
result in the new blueprint being stored as `version 0.1.0`.

### [[packages]] and [[modules]]

These entries describe the package names and matching version glob to be installed into the image.

The names must match the names exactly, and the versions can be an exact match
or a filesystem-like glob of the version using `*` wildcards and `?`
character matching.

NOTE: Currently there is no difference between `packages` and `modules`. In
the future the `modules` list may be used for module support, so it is best to
just use `[[packages]]` for now.

For example, to install `tmux-2.9a` and `openssh-server-8.*`, you would add
this to your blueprint:

```
[[packages]]
name = "tmux"
version = "2.9a"

[[packages]]
name = "openssh-server"
version = "8.*"
```

### [[groups]]

The `groups` entries describe a group of packages to be installed into the image.  Package groups are
defined in the repository metadata.  Each group has a descriptive name used primarily for display
in user interfaces and an ID more commonly used in kickstart files.  Here, the ID is the expected
way of listing a group.

Groups have three different ways of categorizing their packages:  mandatory, default, and optional.
For purposes of blueprints, mandatory and default packages will be installed.  There is no mechanism
for selecting optional packages.

For example, if you want to install the `anaconda-tools` group you would add this to your
blueprint:

```
[[groups]]
name="anaconda-tools"
```

`groups` is a TOML list, so each group needs to be listed separately, like `packages` but with
no version number.


### Customizations

The `[customizations]` section can be used to configure the hostname of the final image. eg.:

```
[customizations]
hostname = "baseimage"
```

This is optional and may be left out to use the defaults.


#### [customizations.kernel]

This allows you to append arguments to the bootloader's kernel commandline. This will not have any
effect on `tar` or `ext4-filesystem` images since they do not include a bootloader.

For example:

```
[customizations.kernel]
append = "nosmt=force"
```

#### [[customizations.sshkey]]

Set an existing user's ssh key in the final image:

```
[[customizations.sshkey]]
user = "root"
key = "PUBLIC SSH KEY"
```

The key will be added to the user's `authorized_keys` file.

---
**WARNING**

`key` expects the entire content of `~/.ssh/id_rsa.pub`, make sure it is the public key.

---


#### [[customizations.user]]

Add a user to the image, and/or set their ssh key.
All fields for this section are optional except for the `name`, here is a complete example:

```
[[customizations.user]]
name = "admin"
description = "Administrator account"
password = "$6$CHO2$3rN8eviE2t50lmVyBYihTgVRHcaecmeCk31L..."
key = "PUBLIC SSH KEY"
home = "/srv/widget/"
shell = "/usr/bin/bash"
groups = ["widget", "users", "wheel"]
uid = 1200
gid = 1200
```

If the password starts with `$6$`, `$5$`, or `$2b$` it will be stored as
an encrypted password. Otherwise it will be treated as a plain text password.

---
**WARNING**

`key` expects the entire content of `~/.ssh/id_rsa.pub`, make sure it is the public key.

---


#### [[customizations.group]]

Add a new group to the image. `name` is required and `gid` is optional:

```
[[customizations.group]]
name = "widget"
gid = 1130
```

#### [customizations.timezone]

Customizing the timezone and the NTP servers to use for the system:

```
[customizations.timezone]
timezone = "US/Eastern"
ntpservers = ["0.north-america.pool.ntp.org", "1.north-america.pool.ntp.org"]
```

The values supported by `timezone` can be listed by running `timedatectl list-timezones`.

If no timezone is setup the system will default to using `UTC`. The ntp servers are also
optional and will default to using the distribution defaults which are fine for most uses.

In some image types there are already NTP servers setup, eg. Google cloud image, and they
cannot be overridden because they are required to boot in the selected environment. But the
timezone will be updated to the one selected in the blueprint.


#### [customizations.locale]

Customize the locale settings for the system:

```
[customizations.locale]
languages = ["en_US.UTF-8"]
keyboard = "us"
```

The values supported by `languages` can be listed by running `localectl list-locales` from
the command line.

The values supported by `keyboard` can be listed by running `localectl list-keymaps` from
the command line.

Multiple languages can be added. The first one becomes the
primary, and the others are added as secondary. One or the other of `languages`
or `keyboard` must be included (or both) in the section.


#### [customizations.firewall]

By default the firewall blocks all access except for services that enable their ports explicitly,
like `sshd`. This command can be used to open other ports or services. Ports are configured using
the port:protocol format:

```
[customizations.firewall]
ports = ["22:tcp", "80:tcp", "imap:tcp", "53:tcp", "53:udp"]
```

Numeric ports, or their names from `/etc/services` can be used in the `ports` enabled/disabled lists.

The blueprint settings extend any existing settings in the image templates, so if `sshd` is
already enabled it will extend the list of ports with the ones listed by the blueprint.

If the distribution uses `firewalld` you can specify services listed by `firewall-cmd --get-services`
in a `customizations.firewall.services` section:

```
[customizations.firewall.services]
enabled = ["ftp", "ntp", "dhcp"]
disabled = ["telnet"]
```

Remember that the `firewall.services` are different from the names in `/etc/services`.

Both are optional, if they are not used leave them out or set them to an empty list `[]`. If you
only want the default firewall setup this section can be omitted from the blueprint.

NOTE: Some compose types disable the firewall, this cannot be overridden by the blueprint.


#### [customizations.services]

This section can be used to control which services are enabled at boot time.
Some image types already have services enabled or disabled in order for the
image to work correctly, and cannot be overridden. eg. `ami` requires
`sshd`, `chronyd`, and `cloud-init`. Without them the image will not
boot. Blueprint services are added to, not replacing, the list already in the
compose type, if any.

The service names are systemd service units. You may specify any systemd unit
file accepted by `systemctl enable` eg. `cockpit.socket`:

```
[customizations.services]
enabled = ["sshd", "cockpit.socket", "httpd"]
disabled = ["postfix", "telnetd"]
```

## Package Sources

By default osbuild-composer uses the host's configured repositories.
These are immutable system
repositories and cannot be deleted or changed. If you want to add additional
repos use the `composer-cli sources` command to create them.

A new source repository can be added by creating a TOML file with the details
of the repository and add it to the server with `composer-cli sources add
newrepo.toml`:

```
name = "custom-source-1"
url = "https://url/path/to/repository/"
type = "yum-baseurl"
proxy = "https://proxy-url/"
check_ssl = true
check_gpg = true
gpgkey_urls = ["https://url/path/to/gpg-key"]
```

The `proxy` and `gpgkey_urls` entries are optional. All of the others are required. The supported
types for the urls are:

* `yum-baseurl` is a URL to a yum repository.
* `yum-mirrorlist` is a URL for a mirrorlist.
* `yum-metalink` is a URL for a metalink.

If `check_ssl` is true the https certificates must be valid. If they are
self-signed you can either set this to false, or add your Certificate Authority
to the host system.

If `check_gpg` is true the GPG key must either be installed on the host system, or `gpgkey_urls`
should point to it.

You can edit an existing source (other than system sources), by using `sources
add` or `sources change` with the new version of the source. It will overwrite
the previous one.

A list of existing sources is available `composer-cli sources list`, and detailed info
on a source can be retrieved with `composer-cli sources info SOURCE-NAME`. Deleting a non-system source is done using `composer-cli sources delete SOURCE-NAME`.

The configured sources are used for all blueprint depsolve operations, and for composing images.
When adding additional sources you must make sure that the packages in the source do not
conflict with any other package sources, otherwise depsolving will fail.

[examples]: https://github.com/osbuild/weldr-client/tree/main/examples
