# Fixture Test Data

The idea behind the `testdata/` is to mock the repos:

* clean
* committed-unpushed-change
* https-remote
* multiple-remotes
* ssh-remote
* ssh-remote-no-at
* uncommited-change
* unstaged-change

In a way that we do not have submodules in Waypoint. Hence, why you will see a
`from-DOTgit.sh` && a `to-DOTgit.sh`. 

## from-DOTgit.sh

Converts the mocked DOTgit repos, essentially the folders in the `testdata/`,
to functioning as separate git repos. After running `from-DOTgit.sh`, changes should
be reflected by running `git status`.

 ## to-DOTgit.sh

Converts all the .git folders to the mocked DOTgit versions, so that Waypoint does 
not have submodules. This is needed for the tests to run successfully.