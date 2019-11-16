# Go Programming Blueprints: First edition

![Go Blueprints by Mat Ryer book cover](https://raw.githubusercontent.com/matryer/goblueprints/master/artwork/bookcover.jpg)

This is the official source code repository for the book. You are welcome to browse this repository and use the [issues tab](https://github.com/matryer/goblueprints/issues) to report any problems or ask any questions.

  * **Feel free to copy and paste from the repository where appropriate**, although typing the code out will surely do more for the learning experience
  * If you are enjoying the book, please tell others by [reviewing the book on Amazon](http://bit.ly/goblueprints)

## Get the book

  * From [Amazon.com](https://www.amazon.co.uk/Programming-Blueprints-real-world-production-ready-cutting-edge/dp/1786468948/ref=sr_1_4?keywords=go+programming+blueprints+ryer&qid=1573895643&sr=8-4) 
  * Free articles about Go in the real world: [https://medium.com/@matryer](https://medium.com/@matryer)

## Projects

Throughout the book many projects, programs and packages are developed, including:

  * [Chat application](https://github.com/matryer/goblueprints/tree/master/chapter3/chat) web application that lets people have conversations in their browsers.
  * [Tracer package](https://github.com/matryer/goblueprints/tree/master/chapter1/trace) package that provides tracing capabilities for Go programs. Built for illustration purposes, in production environments consider using [log.Logger](http://golang.org/pkg/log/#Logger) from the standard library instead.
  * [Domain finder](https://github.com/matryer/goblueprints/tree/master/chapter4/domainfinder) program that helps you find the perfect domain name for your projects, including whether they're available or not. Depends on a series of [subprograms](https://github.com/matryer/goblueprints/tree/master/chapter4) that do its bidding.
  * [Thesaurus](https://github.com/matryer/goblueprints/tree/master/chapter4/thesaurus) package that provides an interface and an implementation for [Big Huge Thesaurus](http://words.bighugelabs.com/) that allows you to lookup synonyms of words.
  * [Meander](https://github.com/matryer/goblueprints/tree/master/chapter7/meander) package that provides random event recommendations with associated [web application](https://github.com/matryer/goblueprints/tree/master/chapter7/meanderweb)
  * [Backup](https://github.com/matryer/goblueprints/tree/master/chapter8/backup) program for periodically backing up your source code.

## Chapters

Each chapter has its own section which it is recommended that you read _before_ embarking on the chapter itself, as updates, tweaks, bug fixes, additional notes and tips will be outlined here.

* See also [the original source code for the first edition of the book](https://github.com/matryer/goblueprints/tree/cb2078d9aa6b5b7cc51e80912be82cbba4d2f9a1)

### Chapter 1

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter1)
  * If you're getting some kind of `version != 13` error, you may want to read towards the bottom of [this issue](https://github.com/matryer/goblueprints/issues/18)

### Chapter 2

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter2)

Notes:

  * Page 53: `w.Header.Set` should be `w.Header().Set` since `Header` is a function on `http.ResponseWriter`.

### Chapter 3

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter3)

Notes:

  * Page 81: For Gravatar to work, you need to hash the email address, not the user's name: `io.WriteString(m, strings.ToLower(user.Email()))` - Thanks [@lozandier](https://github.com/lozandier)

Other minor things:

  * Page 78: Autocompleted typo: `gravatarAvitar` should be `gravatarAvatar` - you can name your variables anything you like, but it's nice for them to be spelled correctly. - Thanks [@lozandier](https://github.com/lozandier)
  * Page 83: The HTML `<label>` is not properly attached to the associated `<input>` - [View Diff](https://github.com/matryer/goblueprints/commit/afb4285f47a7482a58f6fa5061982f874a3fa11e) - Thanks [@crbrox](https://github.com/crbrox) 

### Chapter 4

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter4)

Notes:

  * BigHuge is mistyped in a few places as BigHugh. I don't know who Big Hugh is, but I'm sure he's very nice. Either way, he's a little unwelcome in Chapter 4, so you should consistently type big **HUGE** - Thanks [@OAGr](https://github.com/OAGr)
  * Page 112: Sometimes `data.Noun` and `data.Verb` are `nil`, which causes a panic. Make your code a little more bulletproof by first checking if they're `== nil` before trying to access the `Syn` field. See [Issue #11](https://github.com/matryer/goblueprints/issues/11) for a solution. Thanks [@OAGr](https://github.com/OAGr)

### Chapter 5

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter5)

Notes:

  * There's a data-race with the way I call `updater.Reset` from within the function. [Read more about it here](https://github.com/matryer/goblueprints/issues/12).

### Chapter 6

  * Struct tags are written like this: `json:"title"`
  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter6)

### Chapter 7

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter7)

### Chapter 8

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/chapter8)

### Appendix A

  * Browse the [Source code](https://github.com/matryer/goblueprints/tree/master/appendixA)
