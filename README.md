# PEGGLE DATA COLLECTOR

---

## [web server](https://peggle-data-collector.herokuapp.com/) & [api](https://peggle-data-collector.herokuapp.com/api) that collect all the [peggle games](https://www.youtube.com/watch?v=IwvS8ft7DM8&list=PLkjetvDN3k23J8nTmlDOnxiP3ZXDdCIZD) made by [QDSS 2](https://www.youtube.com/channel/UC5GSO2hiHevgZUhSQIJNd2A)

---

<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#built-with">Built With</a></li>
    <li><a href="#prerequisites">Prerequisites</a></li>
    <li><a href="#set-up">Set Up</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

---

### Built With

- hosting
  - [heroku](https://heroku.com)
- database
  - [mongodb-atlas](https://www.mongodb.com/cloud/atlas)
- front-end
  - [Bootstrap 5](https://getbootstrap.com/docs/5.1/getting-started/introduction/)
  - [Bootstrap 4](https://getbootstrap.com/docs/4.6/getting-started/introduction/)
  - [JQuery](https://jquery.com)
  - [Google Charts](https://developers.google.com/chart)
- back-end
  - [gorilla-mux](https://github.com/gorilla/mux)
  - [mongo-driver](https://pkg.go.dev/go.mongodb.org/mongo-driver?utm_source=godoc)
  - [google-api](https://pkg.go.dev/google.golang.org/api)
  - [yaml.v2](https://gopkg.in/yaml.v2)

---

### Prerequisites

[go](https://golang.org/) >= 1.16

### Set Up

To install all modules required to make the webserver work first you need to clone the repo

```
git clone https://github.com/Vano2903/peggle-data-collector
```

## TODO / v2.0

---

things i wanna implement:

- [ ] hash all the stored passwords
- [ ] add jwt for user authentication
- [ ] add OAuth2.0 (github OAuth seems the best option so far)
- [ ] add full search algorithm (mongob atlas offer that but means migrate to atlas)
- [ ] let user query a date like 2021-2-4 and convert it to 2021-02-04
- [ ] check if a structure has default values (means not modified)

todo for this version:

- [x] colors on buttons in add-game
- [x] user section
- [x] stats database
- [ ] add function that let user clean just a section of the game adding function and not only everything
- [ ] section to add a comment while inserting a game
- [ ] dont create a button in commit area if there are none
- [ ] users section
- [ ] working search box
- [ ] public page dedicated to collaborators
- [x] main page
  - [x] home with game's cards, navbar and search box
  - [x] individual game page with all the stats of the game
  - [ ] players stats (use stats databaes) (scatter graph, bar graph, pie graph)
  - [ ] api documentation area
  - [ ] support area (my paypal or buymeacoffe)
