[![Build Status](https://travis-ci.org/dryvercorp/actor.svg?branch=master)](https://travis-ci.org/dryvercorp/actor) [![Coverage Status](https://coveralls.io/repos/dryvercorp/actor/badge.svg?branch=master&service=github&cachebreak)](https://coveralls.io/github/dryvercorp/actor?branch=master) [![GoDoc](https://godoc.org/github.com/dryvercorp/actor?status.svg)](https://godoc.org/github.com/dryvercorp/actor) 

# Actor: a Gherkin-like actor file specification for behaviour-driven development

This is the Go implementation of the first draft of a .actor file specification to complement the [.feature file format](https://github.com/cucumber/gherkin3) used by [Cucumber](https://cucumber.io/) (and, notably for Gophers, [GoDog](http://data-dog.github.io/godog/)) for [BDD](http://inviqa.com/insights/bdd-guide).  The goal of the format is to record identified actors in a human- and machine-parsable format for the purposes of both manual editing and programmatic interaction, in a format that's familiar to anyone who's used Gherkin.

### Example of the format (from the [example file](https://github.com/dryvercorp/actor/blob/master/examples/valid.actor)):

```
@tag1 @tag2
Actor: Valid actor
    Description and blurb... multi line 1 tab indent
    Some other line of blurb

    @tag3 @tag4
    Goals:
        Goal number 1
        Goal number 2
    
    @tag5 @tag6
    Goal: Goal number 3
```
Only one actor can be defined per file. There are three keywords – `Actor`, `Goal` and `Goals` - which must be followed by a colon and an argument. Keywords can be preceeded by 'tags', which take the the same form as Gherkin tags: an at sign followed by some alphanumeric characters. These tags will then be attached to the resultant object when it's parsed. Any other text is treated as a 'Blurb' – a line of text that describes the actor's motivations, or other notes.

## Example Go code

```
package main

import "github.com/dryvercorp/actor"

func main() {
    parser, err1 := actor.NewFileParser("my.actor")
    // handle err1
    
    actor, err2 := parser.Parse()
    // hande err2
    
    actor.Name = "New actor name"
    
    err3 := actor.WriteToFile("my.actor")
    // handle err3
}
```
See the [GoDoc](https://godoc.org/github.com/dryvercorp/actor) for full documentation.


