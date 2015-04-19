# Go Programming Blueprints

![Go Blueprints by Mat Ryer book cover](https://raw.githubusercontent.com/matryer/goblueprints/master/artwork/bookcover.jpg)

This is the official source code repository for the book. You are welcome to browse this repository and use the [issues tab](https://github.com/matryer/goblueprints/issues) to report any problems or ask any questions.

  * **Feel free to copy and paste from the repository where appropriate**, althrough typing the code out will surely do more for the learning experience
  * If you are enjoying the book, please tell others by [reviewing the book on Amazon](http://bit.ly/goblueprints)

## Get the book

  * From [Amazon.com](http://bit.ly/goblueprints) 

## Projects

Throughout the book many projects, programs and packages are developed, including:

  * [Chat application](https://github.com/matryer/goblueprints/tree/master/chapter3/chat) web application that lets people have conversations in their browsers.
  * [Tracer package](https://github.com/matryer/goblueprints/tree/master/chapter1/trace) package that provides tracing capabilities for Go programs. Built for illustration purposes, in production environments consider using [log.Logger](http://golang.org/pkg/log/#Logger) from the standard library instead.
  * [Domain finder](https://github.com/matryer/goblueprints/tree/master/chapter4/domainfinder) program that helps you find the perfect domain name for your projects, including whether they're available or not. Depends on a series of [subprograms](https://github.com/matryer/goblueprints/tree/master/chapter4) that do its bidding.
  * [Thesaurus](https://github.com/matryer/goblueprints/tree/master/chapter4/thesaurus) package that provides an interface and an implementation for [Big Hugh Thesaurus](http://words.bighugelabs.com/) that allows you to lookup synonyms of words.
  * [Meander](https://github.com/matryer/goblueprints/tree/master/chapter7/meander) package that provides random event recommendations with associated [web application](https://github.com/matryer/goblueprints/tree/master/chapter7/meanderweb)
  * [Backup](https://github.com/matryer/goblueprints/tree/master/chapter8/backup) program for periodically backing up your source code.

## Chapters

Each chapter has its own section which it is recommended that you read _before_ embarking on the chapter itself, as updates, tweaks, bug fixes, additional notes and tips will be outlined here.

### Chapter 1

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter1)

### Chapter 2

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter2)

Notes:

  * Page 53: `w.Header.SetSet` should just be `w.Header.Set` - your compiler will help you spot this one. - Thanks [@lozandier](https://github.com/lozandier)

### Chapter 3

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter3)

Issues:

  * Page 81: For Gravatar to work, you need to hash the email address, not the user's name: `io.WriteString(m, strings.ToLower(user.Email()))` - Thanks [@lozandier](https://github.com/lozandier)

Other minor things:

  * Page 78: Autocompleted typo: `gravatarAvitar` should be `gravatarAvatar` - you can name your variables anything you like, but it's nice for them to be spelled correctly. - Thanks [@lozandier](https://github.com/lozandier)
  * Page 83: The HTML `<label>` is not properly attached to the associated `<input>` - [View Diff](https://github.com/matryer/goblueprints/commit/afb4285f47a7482a58f6fa5061982f874a3fa11e) - Thanks [@crbrox](https://github.com/crbrox) 

### Chapter 4

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter4)

Issues:

  * BigHuge is mistyped in a few places as BigHugh. I don't konw who Big Hugh is, but I'm sure he's very nice. Either way, you should consistently type big **HUGE** - Thanks [@OAGr](https://github.com/OAGr)

### Chapter 5

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter5)

### Chapter 6

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter6)

### Chapter 7

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter7)

### Chapter 8

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter8)

### Appendix A

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/appendixA)
