# prac-go-htmx-tailwind

## Introduction 

This little repo is my attempt at understanding how to implement simple reactive webapp without using React. The idea is to minimise the amount of Javascript I send over the internet. 

I chose to create a contact manager as a way to understand some of tools provided by htmx. 

### Technology Consideration

- [htmx](https://htmx.org/)
    - A small library that adds ability to do AJAX call to server and update the returning partial html without rendering the whole page
    - My reason for choosing this is that it is simple and lightweight with enough functionality to create server-side rendered html with ease
    - Its implementation is very different to React or main stream frameworks in the market today.
- [Templ](https://templ.guide/)
    - A html templating engine for golang.
    - My reason for choosing this is because this is the one I heard about.
- [tailwindcss](https://tailwindcss.com/)
    - A css library that provides a lot of utility classes for faster literation.
    - Reason: I like it.

## Setup

It is assumed that your computer has Make installed. This should come with most Unix based computers.
After cloning the project. Run
```
make setup
```

This should install all the tools required for this project in `bin/` folder. These include:
- golangci-lint - A tool for code formatting and linting
- air - A tool for enabling hot-reloading 
- templ - A tool for generating html components from templ files

You should be able to run to start the local server by executing:
```
make start
```

## Project Structure

- `assets/`: Contains all static assets like `css`, `images` and `js`
- `bin/`: Contains all binary tooling for the project
- `internal/`: Contains all the code related to the server
- `templates/`: Contains all the `templ` templates for html generation
- `Makefile`: Contains all the useful commands for development

The main access point is via `main.go` at the root level.

To build the project, run
```
make build
```
This will install all the necessary packages, then generate all the css to `assets/css/output.css` and `templ` files before compiling golang code to `app/` folder.

## Lessons Learned

- Designing api endpoint for htmx is different than the usual:
    - Using standard golang `net/http` library to serve the content, it doesn't come with ability to parse parameter from url. 
    - Usually for an individual item in a list, the endpoint is of the form `item/{item_id}`.
    - It might better to use query param instead. ie `item?id={item_id}`.
    - This has a few advantages. First, golang `net/http` has a support for parsing query already.
    - Second, it allows for a shorter url and increase flexibility if futher filtering is requiered
    - Third, it works quite well with htmx when returning partial html by appending state, ie. `item/id={item_id}&state=edit`.


