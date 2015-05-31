# go-toggl-cli
Toggl.com console client

Install:

```
go get github.com/iamdeuterium/go-toggl-cli/toggl
```

Set api token and save to ~/.togglrc:

```
toggl token YOUR_API_TOKEN
```

Show api token:

```
toggl token
```

Start last entry

```
toggl start last
```

Select from last entries and start

```
toggl start
```

Start new entry

```
toggl start My new entry
```

Start new entry with project and workspace definition.

```
toggl start My new entry -p AlphaProject -w Job
```

Short version:

```
toggl start My new entry -p Alph -w J
```
