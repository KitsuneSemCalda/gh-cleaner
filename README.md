# gh-cleaner

A small command line interface to clean github using [.netrc](http://devfuria.com.br/git/netrc-nao-pedir-senha/), [bayes theorem](https://en.wikipedia.org/wiki/Bayes%27_theorem) and [prompt-ui](github.com/manifoldco/promptui)

## Example

![gh-cleaner example](https://github.com/KitsuneSemCalda/gh-cleaner/blob/master/assets/gh-cleaner-show.gif)

## How Works

> [!CAUTION]
> This code really delete them repositories, read the instructions with caution

> [!WARNING]
> This project is under development. Contributions are welcome

This code uses the [go-github](https://github.com/google/go-github) library to connect with github using the acess token in .netrc

.netrc is parsed and the username with acess token are used to give go-github authorization to list (private and public) repositories and delete him with user option
