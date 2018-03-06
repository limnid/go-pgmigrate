# Go PGMigrate
A version migration tool, written in Golang

**Example**

```$xslt
./migration --config=config.json [command]
      --init (Init and create migration table)
      --create=test (Create migration)
      --up=[steps]
      --down=[steps]

./migration --config=config.json --init
./migration --config=config.json --create=my_first_migration
./migration --config=config.json --up=1
./migration --config=config.json --down=1
```