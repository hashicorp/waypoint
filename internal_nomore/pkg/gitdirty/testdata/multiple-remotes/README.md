# Test git repo

Note that the `origin` folder is a configured origin.

To recreate this setup:

```shell
mkdir remote
cd remote
git init
mkdir origin
cd origin
git init --bare
cd ..
```

`upstream` is also configured as a remote

From there, you can make and push changes:
```shell
echo "contents" > a.txt
git commit -am "first commit"
git push origin main
```
