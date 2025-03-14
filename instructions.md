# YAG

![yag-logo](../../img/yag-logo.png)

```
TIMEBOX:    4-6 hours max. We mean it! Set at timer and hard-stop at 6 hours ⏱
LANGUAGES:  Go, Rust, C/C++
TESTS:      nice to have, but not mandatory (at least one unit test)
DOCS:       nice to have, but not mandatory
```

## Description
`Yet Another Git` is a simplified and generic `Git` implementation. It's not designed to be the fanciest thing ever, with every bell and whistle under the sun, but contains *just* enough functionality to be useful. To be honest, we barely expect it to even work 😉.

Your task is to build the `YAG` client.

## Design

### Architecture
Code must use some form of [Hash chain](https://en.wikipedia.org/wiki/Hash_chain) or [Merkle DAG](https://factory24.org/wp-content/uploads/2019/10/6.-DAGs-the-Merkle-Forest.pdf) structure for the core `Repo` object.

### Features of Git
1. Initialize a repo
2. Add file(s) to a repo
2. Commit file(s) with a message
3. Create branches
4. Checkout between branches.

You are free to implement any other feature you wish, however these are the core required ones. They're sorted in priority. If you are unable to finish these 4 base features, that is not a problem, as we are more concerned about the structure of your code, then the number of features implemented. See the [Assessment][#Assessment].

### Optional feaures
- State diff'ing
- Branch merging
- Conflict resolution (any strategy)

You are free to choose whatever persistence/storage mechanism you wish. This can be `in-memory`,`key-value store`, `filesystem`, `sqlite`, `network database`, `redis`, etc... Anything you wish.

You are explicilty not required to implement any network related features of `git`, this is a (albeit useless) local only implementation.

## Requirements
- Submitted code
- At least one unit test
- `README.md` file containing:
    - A short explanation of what you built
    - How to test/demo/run (if applicable)
        - NOTE: a 'working' example/client is awesome, however it is NOT a hard requirement. We mean it!
    - Any feedback/notes (i.e. if something was hard, confusing, frustrating, etc)
    - Anything else you'd lke us to know about your submission
- `ROADMAP.md` file  with what you would add/change if you had more time. Dream big.

## Assessment
Your code will be assessed using the same goals and requirements that we use day to day at Source.  on its overall design, structure, and readability. **We rather you implement less features, that are well thought out, then many just for the sake of completeness.**

### Readability 
Readability is an important factor, as we try to follow the concepts and goals of *Literate Programming*. 

> *Literate Programming*:
>
> Instead of imagining that our main task is to instruct a computer what to do, let us concentrate rather on explaining to human beings what we want a computer to do.

This is a subjective goal to evaluate, but something we look, as we believe it is important to effectively communicate to your peers and collaborators what the goal of a specific program/function/snippet is. This is even more important in remote/hybrid work forces, where asynchronos work is common, and you don't always have the freedom to casually talk to your coworker beside you.

### Abstractions and Structure
The common buzzword philisophy for Starups is to "Move fast, and break things". In certain environments, this can be an effective stratedgy. However at Source, we take a more measured and slower approach to most things we do.

Our goals and engineering challenges here at Source often take us to the bleeding edge of various technologies, and when walking on new ground, we want to be confident that the steps we are taking are the rights ones, so we may have a better foundation to build on later. 

It's one thing if your ToDo SaaS app goes offline for an hour, but a Database that irevocably destroys or losses data is a non-starter

As such, running code on a machine is usually the last step in the engineering process here at Source, which starts with ideation and design based on clearly defined goals.

## Need Help?
If you need any help, want to ask a question, need clarification, or anything else, feel free to join our Community Discord to chat with our developers. Or, simply open up an issue thread on this repo. Both are suggested.

If you are stuck, and have reached out for help either via [Issues](https://github.com/shinzonetwork/hiring/issues) or [Discord](https://discord.gg/57UNewmXtE), feel free to put the project away untill you have got your answer. As there is a working time limit of 6 hours, and we don't want you to waste it running in circles on a problem or clarification. 

> Hint: Some of this document is intentionally vague

## Background
> Note: You don't have to implement the internals of git exactly, you may use whatever sturcture/design you wish, as long as it complies with the above [Design](#Design) goals.

### How git works resources
 - https://hackernoon.com/understanding-git-fcffd87c15a3
 - https://shalithasuranga.medium.com/how-does-git-work-internally-7c36dcb1f2cf