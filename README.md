# Backup
This is a simple programm to create a backup from a directory.

## How it will work
For each backup a new folder with the current timestamp will be created. If there are modified files, since the last backup they will be copied. If a file was not modified a hardlink to the previous backup will be created. The file won't be copied in that case.

---

## How to build
```
go build backup.go
```

---

## How to run
```
backup -s=<your source> -t=<your target>
```
---

## How the backup will look like
```
D:\TEST\TARGET
+---20180814010127
|   |   1.txt
|   |
|   \---Neuer Ordner
|           2.txt
|           foo.txt
|
+---20180814010219
|   |   1.txt
|   |   3.txt
|   |
|   +---Neuer Ordner
|   |       2.txt
|   |       foo.txt
|   |
|   \---Neuer Ordner (2)
|           4.txt
|
+---20180814012134
|   |   1.txt
|   |   3.txt
|   |
|   +---Neuer Ordner
|   |       2.txt
|   |       foo.txt
|   |
|   \---Neuer Ordner (2)
|           4.txt
|
+---20180814013345
|   |   1.txt
|   |   3.txt
|   |
|   +---Neuer Ordner
|   |       2.txt
|   |       foo.txt
|   |
|   \---Neuer Ordner (2)
|           4.txt
|
\---20180814020828
    |   1.txt
    |   3.txt
    |
    +---Neuer Ordner
    |       2.txt
    |       foo.txt
    |
    \---Neuer Ordner (2)
            4.txt
```