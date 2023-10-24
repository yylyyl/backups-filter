# backups-filter

This tool helps you remove old daily backups. For example:
- Keep 7 daily backups, and
- Keep 3 weekly backups, and
- Keep 5 monthly backups.

There will be 15 backups to be kept in total.
Other backups should be deleted.

```text
<-Past                                                          Now->
...-----B-----------------B-----------------B-----B-----B-----BBBBBBB
        ^                                   ^                    ^
5 monthly backups                   3 weekly backups      7 daily backups
```

## How to use

Arguments:

```text
Usage of ./backups-filter:
  -backup-daily uint
        number of daily backups (default 7)
  -backup-monthly uint
        number of monthly backups (default 5)
  -backup-weekly uint
        number of weekly backups (default 3)
  -descending
        descending input and output
  -keep
        print the items which should be kept, instead of which should be deleted
  -layout string
        layout of datetime, but this tool only takes the dates into consideration (default "20060102_150405")
```
To use this tool correctly, you should do these things in your backing-up scripts:
1. Do the ordinary backup procedures, normally uploading your files
2. Write the current date into a file, e.g. `date +%Y%m%d_%H%M%S >> backups.txt`
3. Invoke this tool, e.g. `TO_DELETE=$(cat backups.txt | backups-filter)`
4. Delete your old backups, and remove the date line from backups.txt respectively

## Example back up script

```shell
#!/bin/bash
set -e

# upload files to the cloud
DATETIME="$(date +%Y%m%d_%H%M%S)"
s3cmd put xxx.zip "s3://my-bucket/backups/${DATETIME}.zip"

# write the current datetime into a record file
date +%Y%m%d_%H%M%S >> backups.txt

# get the items that should be removed
TO_DELETE=$(cat backups.txt | backups-filter)

# remove items from the cloud and from the record file
for DT in ${TO_DELETE}
do
  s3cmd del "s3://my-bucket/backups/${DT}.zip"
  sed -i '' '/^'"${DT}"'$/d' backups.txt
done
```
