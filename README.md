# gabi-cli

`gabi-cli` simplifies the configuration and usage of [Gabi](https://github.com/app-sre/gabi).

## Install

### Option 1: Build from source
First clone the repository somewhere in your $GOPATH.

Example:

```
mkdir $GOPATH/src/github.com/app-sre
cd $GOPATH/src/github.com/app-sre
git clone git@github.com:app-sre/gabi-cli.git
```

Next, cd into the `gabi-cli` folder and run `make`. This command will build the `gabi` binary and place it in $GOPATH.

## Configuration

The `gabi-cli` configuration is done via `profiles` stored in the `$HOME/.config/gabi/gabi.json` file. Here is what the file should look like for two profiles:

```
[
  {
    "name": "ams-stage",
    "alias": "stage",
    "url": "https://...",
    "token": "sha256~...",
    "current": true
  },
  {
    "name": "ams-production",
    "alias": "production",
    "url": "https://...",
    "token": "sha256~...",
    "current": false
  }
]
```

Note: This configuration can be bootstraped via the gabi config init command (see details below)

## Usage

The two main commands of gabi-cli are `config` and `exec`.

### Config

Gets or Sets the gabi-cli configs.

`gabi config -h`

```
Gets or Sets the gabi-cli configs

Usage:
  gabi configure [command]

Aliases:
  configure, config

Available Commands:
  allprofiles       Gets all profiles currently configured for gabi-cli
  currentprofile    Gets the current profile
  init              Initializes a gabi-cli config by creating a gabi.json config file under the ~/.config/gabi directory
  setcurrentprofile Sets the current profile
  settoken          Sets the token for current profile
  seturl            Sets the url for current profile
```

Get started by creating the gabi config file:

`gabi config init`

`Gabi init success! Check /$HOME/.config/gabi/gabi.json for details and complete the setup.`

This will create the `gabi.json` in the following format:

```
[
  {
    "name": "default",
    "alias": "default",
    "url": "",
    "token": "",
    "current": true
  }
]
```

As a next step you can configure the URL and Token (either editing the `gabi.json` file or using the `seturl` and `settoken` commands)

`gabi config seturl https://XXX...`

`gabi config settoken sha256~XXX...`

Confirm everything is ready by running the `gabi config currentprofile` command

```
{
  "name": "default",
  "alias": "default",
  "url": "https://XXX...",
  "token": "sha256~XXX...",
  "current": true
}
```

Note: the token attribute will be redacted for security.

### Exec

Executes queries using the gabi-cli.

`gabi exec -h`

```
Executes a gabi query received from a string as argument, from the file contents specified by a file path with an ".sql" extension, or from stdin. When using stdin, press Enter to move to the next line and then CTRL+D to execute the query (or CTRL+C to Cancel)

Usage:
  gabi execute [string] | [sql_file_path] | stdin  [flags]

Aliases:
  execute, exec

Flags:
      --csv              CSV output
  -h, --help             help for execute
      --raw              Raw output
      --show-row-count   Prints out the number of rows returned by your query
```

By default, gabi-cli will format the output as json.

Here are some examples of what output looks like:

#### JSON Output

`gabi exec "select * from cloud_resources where resource_type='compute.node' and cloud_provider='aws' limit 2"`

```
[
  {
    "active": "true",
    "category": "compute_optimized",
    "category_pretty": "Compute optimized",
    "cloud_provider": "aws",
    "cpu_cores": "48",
    "created_at": "0001-01-01T00:00:00Z",
    "deleted_at": "",
    "generic_name": "highcpu-48-5a",
    "id": "c5a.12xlarge",
    "memory": "103079215104",
    "memory_pretty": "96",
    "name_pretty": "c5a.12xlarge - Compute optimized",
    "resource_type": "compute.node",
    "size_pretty": "12xlarge",
    "updated_at": "2022-10-10T13:55:03.159046Z"
  },
  {
    "active": "true",
    "category": "compute_optimized",
    "category_pretty": "Compute optimized",
    "cloud_provider": "aws",
    "cpu_cores": "64",
    "created_at": "0001-01-01T00:00:00Z",
    "deleted_at": "",
    "generic_name": "highcpu-64-5a",
    "id": "c5a.16xlarge",
    "memory": "137438953472",
    "memory_pretty": "128",
    "name_pretty": "c5a.16xlarge - Compute optimized",
    "resource_type": "compute.node",
    "size_pretty": "16xlarge",
    "updated_at": "2022-10-10T13:55:03.162961Z"
  }
]
```

#### CSV Output

Prints a comma-separated output in the terminal

`gabi exec "select * from cloud_resources where resource_type='compute.node' and cloud_provider='aws' limit 2" --csv`

```
id,created_at,updated_at,deleted_at,name_pretty,generic_name,cloud_provider,resource_type,category,category_pretty,cpu_cores,memory,memory_pretty,size_pretty,active
c5a.12xlarge,0001-01-01T00:00:00Z,2022-10-10T13:55:03.159046Z,,c5a.12xlarge - Compute optimized,highcpu-48-5a,aws,compute.node,compute_optimized,Compute optimized,48,103079215104,96,12xlarge,true
c5a.16xlarge,0001-01-01T00:00:00Z,2022-10-10T13:55:03.162961Z,,c5a.16xlarge - Compute optimized,highcpu-64-5a,aws,compute.node,compute_optimized,Compute optimized,64,137438953472,128,16xlarge,true
```

#### RAW Output

The original gabi response without any formatting

`gabi exec "select * from cloud_resources where resource_type='compute.node' and cloud_provider='aws' limit 2" --raw`

```
[
  [
     "id",
     "created_at",
     "updated_at",
     "deleted_at",
     "name_pretty",
     "generic_name",
     "cloud_provider",
     "resource_type",
     "category",
     "category_pretty",
     "cpu_cores",
     "memory",
     "memory_pretty",
     "size_pretty",
     "active"
  ],
  [
     "c5a.12xlarge",
     "0001-01-01T00:00:00Z",
     "2022-10-10T13:55:03.159046Z",
     "",
     "c5a.12xlarge - Compute optimized",
     "highcpu-48-5a",
     "aws",
     "compute.node",
     "compute_optimized",
     "Compute optimized",
     "48",
     "103079215104",
     "96",
     "12xlarge",
     "true"
  ],
  [
     "c5a.16xlarge",
     "0001-01-01T00:00:00Z",
     "2022-10-10T13:55:03.162961Z",
     "",
     "c5a.16xlarge - Compute optimized",
     "highcpu-64-5a",
     "aws",
     "compute.node",
     "compute_optimized",
     "Compute optimized",
     "64",
     "137438953472",
     "128",
     "16xlarge",
     "true"
  ]
]
```


### Multi-line Select Statements

Multi-line Select statements can be performed using HEREDOC or directly via the Standard Input (stdin)

#### Using HEREDOC

```
gabi exec << EOF 
select 
    count(id) as total_resources,
    category as cat
from cloud_resources cr
group by category
having count(id) > 50
EOF
```

```
[
  {
    "cat": "compute_optimized",
    "total_resources": "100"
  },
  {
    "cat": "memory_optimized",
    "total_resources": "200"
  }
    {
    "cat": "general_purpose",
    "total_resources": "300"
  }
]
```


#### Directly from Standard Input (stdin)

1. Type `gabi exec` and press `Enter`
2. Type your query (it can include line breaks and tabs)
3. When you are done typing it, press `Enter` again to move to the next line and then `CTRL+D` to indicate the end of the input - this will trigger the query execution. Alternatively, press `CTRL+C` to Cancel.

```
gabi exec [ENTER]
select count(id) as total_resources,
       cloud_provider as provider
from cloud_resources cr
group by cloud_provider
having count(id) > 10

[Ctrl+D]
```

```
[
  {
    "provider": "aws",
    "total_resources": "50"
  },
  {
    "provider": "gcp",
    "total_resources": "50"
  }
]
```

### History

Gabi-CLI can keep track of the executed queries.

This feature is disabled by default and can be enabled via the `gabi config enablehistory` command.

`gabi history -h`

```
Executes history operations

Usage:
gabi history [command]

Available Commands:
clear       Clears gabi-cli query history
show        Shows gabi-cli query history

Flags:
-h, --help   help for history
```

After a few queries, you should be able to run the `gabi history show` command which will output something like this:

```
5 select * from cloud_resources where resource_type='compute.node' and cloud_provider='aws'
4 select * from cloud_resources where resource_type='compute.node' and cloud_provider='gcp'
3 select count(id) as total_resources, category as cat from cloud_resources cr group by category having count(id) > 50 
2 select count(id) as total_resources, cloud_provider as provider from cloud_resources cr group by cloud_provider having count(id) > 10 
1 select count(*) from cloud_resources
```

Note: the `gabi history show` command returns the last 100 queries. This number can be overridden with the flag `--max-rows`:

`gabi history show -h`

```
...

Flags:
-h, --help           help for show
--max-rows int   Maximum number of rows returned in the show command (default 100)
```


Use the `gabi history clear` command to wipe out all query logs.
