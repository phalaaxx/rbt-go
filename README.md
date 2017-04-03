# rbt

RBT stands for Rsync Backup Tool and is a wrapper around rsync to simplify the process of making incremental backups from remote servers to local filesystems with rsync.

Usage
--

The program accepts a single positional argument - a json file name which describes what remote files should be copied and where on the local filesystems to store backups.
A sample json configuration file looks like this - server.hostname.com.json:

```json
{
    "name": "server.hostname.com",
    "backups": 182,
    "target": "/backup/webservers/$name",
    "files": [
        "/opt",
        "/etc",
        "/var",
        "/home",
        "/root",
        "/backup",
    ],
    "exclude": [
        "/var/lib/mysql",
        "/var/lib/postgresql"
    ]
}
```

Configuration file names are provided with -f command line option. It is possible to provide more than one file to sequentially run backups on multiple servers:

	rbt -f server1.hostname.com -f server5.otherhostname.com

Configuration files must be in JSON format and file names must end with .json extension, however it is not necessary to specify file extension with -f option. It will be appended if necessary.
If the provided file name does not exist, RBT will try to also search for it in /etc/rbt directory.

Authentication
--

RBT does not provide any means to set remote user, port and authentication method. This configuration should be present in ~/.ssh/config instead.
