# rbt

RBT stands for Rsync Backup Tool and is a wrapper around rsync to simplify the process of making incremental backups from remote servers to local filesystems with rsync.

Usage
--

The program accepts a single positional argument - a json file name which describes what remote files should be copied and where on the local filesystems to store backups.
A sample json configuration file looks like this:

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

RBT does not provide any means to set remote user, port and authentication method. This configuration should be present in ~/.ssh/config instead.
