## Contributing to the ABS language

## A little explanation first...

ABS uses version branches to keep track of changes, meaning
that you will see branches such as [1.0.x](https://github.com/abs-lang/abs/tree/1.0.x) and [1.1.x](https://github.com/abs-lang/abs/tree/1.1.x) in the
repository.

Since we follow [semver](https://semver.org/),
when a new feature is released we don't backport it but simply
create a new version branch, such as `1.3.x`. Bugs, instead,
might be backported from `1.1.0` to, for example, `1.0.x` and we
will have a new [release](https://github.com/abs-lang/abs/releases),
say `1.0.1` for the `1.0.x` version branch.

The infamous `master` branch, instead, is used as a base for
the documentation and latest releases: when we publish a new release,
the version branch gets also merged into `master`, making sure
`master` now has all the latest changes we've released to the
public. This is important since [abs-lang.org](https://www.abs-lang.org/)
gets built whenever new commits land on master, so it would
be problematic to merge PRs directly into master -- the website
and documentation would get instantly updated while the changes 
we merged haven't been released yet.

Why is this important? When you want to contribute to ABS you
should try to figure out against what branch you want changes
to be incorporated. You can open a PR directly to `master`, but
it's likely a member of our team will then change your PR's
base branch to the right version branch. If you do that beforehand,
that'll help us getting things reviewed, and merged, faster :)

## Hacking on ABS

The best way to start hacking on ABS is to clone the repository
and run a `make build`: this will build a docker container with all
the necessary dependencies for developing locally (for experienced
Gophers: you might want to skip this altogether as your environment
will probably work perfectly).

With `make run` you can get inside a container built for ABS'
development, and `make test` will run all tests.

We're planning to switch over to `go mod` in the near future to make
it easier to contribute to ABS and will keep updating this page
with additional directions on how to develop on ABS locally.

## Pull Requests status checks

When you send a PR, it will be automatically tested through
[travis-ci](https://travis-ci.com/abs-lang/abs): remember
we strive to support `linux`, `osx` and `windows`, so sometimes
changes that work locally might trigger failures on travis --
don't be afraid of a red build, it happens to everyone! Once
you find the culprit, push another change to your PR and, once
tests are green, we're ready to merge!

In order to get something merged, this is the broad set
of conventions we follow:

* PR from external contributors must be reviewed and approved by a member of the ABS team
* builds for the PR must be green on travis

...and that's it, we're not very formal!

## Issues? Roadmap?

You can have a look at the list of [open issues on GitHub](https://github.com/abs-lang/abs/issues)
in order to get an idea what we'd like some help on: at the same
time, if you want to propose a change to the language itself, feel
free to open an issue and we'd be delighted to have a discussion
around your ideas!

We also plan [milestones](https://github.com/abs-lang/abs/milestones)
ahead of time so you can get an idea of what's going to land in each
minor / major release by looking at the GitHub milestones. We try
to release early and release often, meaning we're constantly reviewing
our plan and might need to delay releasing a particular feature as
it would hold back the release of other interesting features.

A piece of advice: if you `+1` the issues you care about, there's a
higher change we'll have a look at them in our next release ;-)

We also encourage you to start discussions around the language
or submit RFCs: ABS is nothing without a community that shapes
its direction. [This](https://github.com/abs-lang/abs/issues/124)
is an example of the community proposing a change and the ABS
team happily `+1` it.
