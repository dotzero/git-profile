# Git Profile switcher

[![Build Status](https://travis-ci.org/dotzero/git-profile.svg?branch=master)](https://travis-ci.org/dotzero/git-profile)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotzero/git-profile)](https://goreportcard.com/report/github.com/dotzero/git-profile)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dotzero/git-profile/blob/master/LICENSE)

Allows you to switch between multiple user profiles in git repositories

## Installing

 ```bash
go get -u github.com/gesquive/git-user
```

## Usage

```
git profile set home user.email="mail@dotzero.ru"
git profile set home user.name=dotzero
git profile use home
```

## License

http://www.opensource.org/licenses/mit-license.php
