# Shared Internal Packages

This folder contains Go packages that are meant for internal use only.
They are in this folder rather than "internal" since they are also shared
with the core project.

If you're writing a plugin or are any other type of external user, please
do not use these packages directly. We will NOT be maintaining normal
semver-style compatibility on these packages and no compatibility is
guaranteed at all.
